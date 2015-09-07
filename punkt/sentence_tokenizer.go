package punkt

import (
	"strings"
)

type SentenceTokenizer interface {
	Tokenize(string) []string
	HasSentBreak(string) bool
	AnnotateTokens([]*DefaultToken) []*DefaultToken
	SecondPassAnnotation(*DefaultToken, *DefaultToken)
	AnnotateSecondPass([]*DefaultToken) []*DefaultToken
	OrthoHeuristic(*DefaultToken) int
}

// A sentence tokenizer which uses an unsupervised algorithm to build a model
// for abbreviation words, collocations, and words that start sentences
// and then uses that model to find sentence boundaries.
type DefaultSentenceTokenizer struct {
	*Base
	SentenceTokenizer
	Punctuation []string
}

var Punctuation = []string{";", ":", ",", ".", "!", "?"}

func NewSentenceTokenizer(trainedData *Storage) *DefaultSentenceTokenizer {
	st := &DefaultSentenceTokenizer{
		Base:        NewBase(),
		Punctuation: Punctuation,
	}

	st.SentenceTokenizer = st
	st.Storage = trainedData
	return st
}

func (s *DefaultSentenceTokenizer) Tokenize(text string) []string {
	re := s.RePeriodContext()
	matches := re.FindAllStringSubmatchIndex(text, -1)

	sentences := make([]string, 0, len(matches))
	lastBreak := 0
	matchEnd := 0
	/*
	 * match = [15, 23, 20, 23, 21, 23]
	 * entire match = 0:1
	 * second token = 2:3
	 * newlines + second token = 4:5
	 */
	for _, match := range matches {
		context := text[match[0]:match[1]]

		nextTok := ""
		if match[2] != -1 && match[3] != -1 {
			nextTok = text[match[2]:match[3]]
		}

		matchEnd = match[1]
		// we want the extra stuff for the actual sentence
		if match[4] >= 0 && (!s.SentenceTokenizer.HasSentBreak(nextTok) || s.SentenceTokenizer.HasSentBreak(text[match[0]:match[4]])) {
			matchEnd = match[4]
		}

		if s.SentenceTokenizer.HasSentBreak(context) {
			noNewline := text[lastBreak:matchEnd]
			s := strings.TrimSpace(noNewline)
			if s == "" {
				continue
			}

			sentences = append(sentences, s)
			lastBreak = matchEnd
		}
	}

	sentences = append(sentences, text[lastBreak:])
	return sentences
}

/*
Returns True if the given text includes a sentence break.
*/
func (s *DefaultSentenceTokenizer) HasSentBreak(text string) bool {
	tokens := s.TokenizeWords(text)

	if len(tokens) == 0 {
		return false
	}

	for _, t := range s.SentenceTokenizer.AnnotateTokens(tokens) {
		if t.SentBreak {
			return true
		}
	}

	return false
}

/*
Given a set of tokens augmented with markers for line-start and
paragraph-start, returns an iterator through those tokens with full
annotation including predicted sentence breaks.
*/
func (s *DefaultSentenceTokenizer) AnnotateTokens(tokens []*DefaultToken) []*DefaultToken {
	//Make a preliminary pass through the document, marking likely
	//sentence breaks, abbreviations, and ellipsis tokens.
	tokens = s.AnnotateFirstPass(tokens)

	/*for _, tok := range tokens {
		logger.Println(tok.Tok, tok.SentBreak)
	}*/

	tokens = s.SentenceTokenizer.AnnotateSecondPass(tokens)

	/*for _, tok := range tokens {
		logger.Println(tok.Tok, tok.SentBreak)
	}*/
	return tokens
}

/*
Performs a token-based classification (section 4) over the given
tokens, making use of the orthographic heuristic (4.1.1), collocation
heuristic (4.1.2) and frequent sentence starter heuristic (4.1.3).
*/
func (s *DefaultSentenceTokenizer) AnnotateSecondPass(tokens []*DefaultToken) []*DefaultToken {
	for _, tokPair := range s.pairIter(tokens) {
		s.SecondPassAnnotation(tokPair[0], tokPair[1])

	}
	return tokens
}

func (s *DefaultSentenceTokenizer) SecondPassAnnotation(tokOne, tokTwo *DefaultToken) {
	if tokTwo == nil {
		return
	}

	if !tokOne.PeriodFinal {
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
	if s.Collocations.items[collocation] != 0 {
		tokOne.SentBreak = false
		tokOne.Abbr = true
		return
	}

	/*
		[4.2. Token-Based Reclassification of Abbreviations] If
		the token is an abbreviation or an ellipsis, then decide
		whether we should *also* classify it as a sentbreak.
	*/
	if (tokOne.Abbr || tokOne.Ellipsis) && !tokIsInitial {
		/*
			[4.1.1. Orthographic Heuristic] Check if there's
			orthogrpahic evidence about whether the next word
			starts a sentence or not.
		*/
		isSentStarter := s.OrthoHeuristic(tokTwo)
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
		if tokTwo.FirstUpper() && s.SentStarters.items[nextTyp] != 0 {
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
		isSentStarter := s.OrthoHeuristic(tokTwo)

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
			s.OrthoContext.items[nextTyp]&orthoLc == 0 {

			tokOne.SentBreak = false
			tokOne.Abbr = true
			return
		}
	}
}

/*
Decide whether the given token is the first token in a sentence.
*/
func (s *DefaultSentenceTokenizer) OrthoHeuristic(token *DefaultToken) int {
	if token == nil {
		return 0
	}

	for _, punct := range s.Punctuation {
		if token.Tok == punct {
			return 0
		}
	}

	orthoCtx := s.OrthoContext.items[token.TypeNoSentPeriod()]

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
