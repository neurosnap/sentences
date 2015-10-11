package punkt

import (
	"strings"
)

type AnnotateTokens interface {
	Annotate([]*Token) []*Token
}

type TypeBasedAnnotation struct {
	*Storage
	PunctStrings
}

func NewTypeBasedAnnotation() *TypeBasedAnnotation {
	return &TypeBasedAnnotation{
		PunctStrings: NewLanguage(),
		Storage:      NewStorage(),
	}
}

/*
Perform the first pass of annotation, which makes decisions
based purely based on the word type of each word:
	- '?', '!', and '.' are marked as sentence breaks.
	- sequences of two or more periods are marked as ellipsis.
	- any word ending in '.' that's a known abbreviation is marked as an abbreviation.
	- any other word ending in '.' is marked as a sentence break.

Return these annotations as a tuple of three sets:
	- sentbreak_toks: The indices of all sentence breaks.
	- abbrev_toks: The indices of all abbreviations.
	- ellipsis_toks: The indices of all ellipsis marks.
*/
func (a *TypeBasedAnnotation) Annotate(tokens []*Token) []*Token {
	for _, augTok := range tokens {
		a.typeAnnotation(augTok)
	}
	return tokens
}

func (a *TypeBasedAnnotation) typeAnnotation(token *Token) {
	chars := []rune(token.Tok)

	if token.HasSentEndChars() {
		token.SentBreak = true
	} else if token.HasPeriodFinal() && !strings.HasSuffix(token.Tok, "..") {
		tokNoPeriod := strings.ToLower(token.Tok[:len(chars)-1])
		tokNoPeriodHypen := strings.Split(tokNoPeriod, "-")
		tokLastHyphEl := string(tokNoPeriodHypen[len(tokNoPeriodHypen)-1])

		if a.IsAbbr(tokNoPeriod, tokLastHyphEl) {
			token.Abbr = true
		} else {
			token.SentBreak = true
		}
	}
}

type TokenBasedAnnotation struct {
	*Storage
	PunctStrings
	TokenGrouper
}

/*
Performs a token-based classification (section 4) over the given
tokens, making use of the orthographic heuristic (4.1.1), collocation
heuristic (4.1.2) and frequent sentence starter heuristic (4.1.3).
*/
func (a *TokenBasedAnnotation) Annotate(tokens []*Token) []*Token {
	for _, tokPair := range a.TokenGrouper.Group(tokens) {
		a.tokenAnnotation(tokPair[0], tokPair[1])
	}

	return tokens
}

func (a *TokenBasedAnnotation) tokenAnnotation(tokOne, tokTwo *Token) {
	if tokTwo == nil {
		return
	}

	if !tokOne.HasPeriodFinal() {
		return
	}

	typ := tokOne.TypeNoPeriod()
	nextTyp := tokTwo.TypeNoSentPeriod()
	tokIsInitial := tokOne.IsInitial()

	/*
	   [4.1.2. Collocation Heuristic] If there's a
	   collocation between the word before and after the
	   period, then label tok as an abbreviation and NOT
	   a sentence break. Note that collocations with
	   frequent sentence starters as their second word are
	   excluded in training.
	*/
	collocation := strings.Join([]string{typ, nextTyp}, ",")
	if a.Collocations[collocation] != 0 {
		tokOne.SentBreak = false
		tokOne.Abbr = true
		return
	}

	/*
		[4.2. Token-Based Reclassification of Abbreviations] If
		the token is an abbreviation or an ellipsis, then decide
		whether we should *also* classify it as a sentbreak.
	*/
	if (tokOne.Abbr || tokOne.IsEllipsis()) && !tokIsInitial {
		/*
			[4.1.1. Orthographic Heuristic] Check if there's
			orthogrpahic evidence about whether the next word
			starts a sentence or not.
		*/
		isSentStarter := a.orthoHeuristic(tokTwo)
		if isSentStarter == 1 {
			tokOne.SentBreak = true
			return
		}

		/*
			[4.1.3. Frequent Sentence Starter Heruistic] If the
			next word is capitalized, and is a member of the
			frequent-sentence-starters list, then label tok as a
			sentence break.
		*/
		if tokTwo.FirstUpper() && a.SentStarters[nextTyp] != 0 {
			tokOne.SentBreak = true
			return
		}
	}

	/*
		[4.3. Token-Based Detection of Initials and Ordinals]
		Check if any initials or ordinals tokens that are marked
		as sentbreaks should be reclassified as abbreviations.
	*/
	if tokIsInitial || typ == "##number##" {
		isSentStarter := a.orthoHeuristic(tokTwo)

		if isSentStarter == 0 {
			tokOne.SentBreak = false
			tokOne.Abbr = true
			if tokIsInitial {
				return
			} else {
				return
			}
		}

		/*
			Special heuristic for initials: if orthogrpahic
			heuristc is unknown, and next word is always
			capitalized, then mark as abbrev (eg: J. Bach).
		*/
		if isSentStarter == -1 &&
			tokIsInitial &&
			tokTwo.FirstUpper() &&
			a.OrthoContext[nextTyp]&orthoLc == 0 {

			tokOne.SentBreak = false
			tokOne.Abbr = true
			return
		}
	}
}

/*
Decide whether the given token is the first token in a sentence.
*/
func (a *TokenBasedAnnotation) orthoHeuristic(token *Token) int {
	if token == nil {
		return 0
	}

	for _, punct := range a.Punctuation() {
		if token.Tok == string(punct) {
			return 0
		}
	}

	orthoCtx := a.OrthoContext[token.TypeNoSentPeriod()]

	/*
	   If the word is capitalized, occurs at least once with a
	   lower case first letter, and never occurs with an upper case
	   first letter sentence-internally, then it's a sentence starter.
	*/
	if token.FirstUpper() && (orthoCtx&orthoLc > 0 && orthoCtx&orthoMidUc == 0) {
		return 1
	}

	/*
		If the word is lower case, and either (a) we've seen it used
		with upper case, or (b) we've never seen it used
		sentence-initially with lower case, then it's not a sentence
		starter.
	*/
	if token.FirstLower() && (orthoCtx&orthoUc > 0 || orthoCtx&orthoBegLc == 0) {
		return 0
	}

	return -1
}
