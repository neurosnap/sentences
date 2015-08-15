package punkt

import (
	"fmt"
	"strings"
)

// A sentence tokenizer which uses an unsupervised algorithm to build a model
// for abbreviation words, collocations, and words that start sentences
// and then uses that model to find sentence boundaries.
type SentenceTokenizer struct {
	*PunktBase
	Punctuation []string
}

func NewSentenceTokenizer(trainedData *PunktParameters) *SentenceTokenizer {
	st := &SentenceTokenizer{
		PunktBase:   NewPunktBase(),
		Punctuation: []string{";", ":", ",", ".", "!", "?"},
	}

	st.PunktParameters = trainedData
	return st
}

func (s *SentenceTokenizer) Tokenize(text string) []string {
	re := s.PunktLanguageVars.RePeriodContext()
	matches := re.FindAllStringSubmatchIndex(text, -1)

	sentences := make([]string, 0, len(matches))
	lastBreak := 0
	for _, match := range matches {
		context := text[match[0]:match[1]]
		nextTok := text[match[4]:match[5]]
		matchStart := match[2]
		matchEnd := match[4]

		/*fmt.Println("context: ", context)
		fmt.Println("next_tok: ", nextTok)
		fmt.Println("start: ", matchStart)
		fmt.Println("end: ", matchEnd)
		fmt.Println("-------")*/

		if s.hasSentBreak(context) {
			noNewline := strings.Replace(text[lastBreak:matchEnd], "\n", "", -1)
			s := strings.Trim(noNewline, " ")
			sentences = append(sentences, s)
			if nextTok != "" {
				lastBreak = matchStart
			} else {
				lastBreak = matchEnd
			}
		}
	}

	return sentences
}

/*
Returns True if the given text includes a sentence break.
*/
func (s *SentenceTokenizer) hasSentBreak(text string) bool {
	found := false
	for _, t := range s.annotateTokens(s.TokenizeWords(text)) {
		if found {
			return true
		}
		if t.SentBreak {
			found = true
		}
	}

	return false
}

/*
Given a set of tokens augmented with markers for line-start and
paragraph-start, returns an iterator through those tokens with full
annotation including predicted sentence breaks.
*/
func (s *SentenceTokenizer) annotateTokens(tokens []*PunktToken) []*PunktToken {
	//Make a preliminary pass through the document, marking likely
	//sentence breaks, abbreviations, and ellipsis tokens.
	tokens = s.annotateFirstPass(tokens)

	tokens = s.annotateSecondPass(tokens)

	return tokens
}

/*
Performs a token-based classification (section 4) over the given
tokens, making use of the orthographic heuristic (4.1.1), collocation
heuristic (4.1.2) and frequent sentence starter heuristic (4.1.3).
*/
func (s *SentenceTokenizer) annotateSecondPass(tokens []*PunktToken) []*PunktToken {
	for _, tokPair := range s.pairIter(tokens) {
		s.secondPassAnnotation(tokPair[0], tokPair[1])

	}
	return tokens
}

func (s *SentenceTokenizer) secondPassAnnotation(tokOne, tokTwo *PunktToken) {
	if tokTwo == (&PunktToken{}) {
		return
	}

	if !tokOne.PeriodFinal {
		return
	}

	typ := tokOne.TypeNoPeriod()
	//nextTok := tokTwo.Tok
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
	if s.PunktParameters.Collocations.items[collocation] != 0 {
		tokOne.SentBreak = false
		tokOne.Abbr = true
		fmt.Println("REASON KNOWN COLLOCATION")
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
		isSentStarter := s.orthoHeuristic(tokTwo)
		if isSentStarter == 1 {
			tokOne.SentBreak = true
			fmt.Println("REASON ABBR WITH ORTHO HEURISTIC")
			return
		}

		/*
			[4.1.3. Frequent Sentence Starter Heruistic] If the
			next word is capitalized, and is a member of the
			frequent-sentence-starters list, then label tok as a
			sentence break.
		*/
		if tokTwo.FirstUpper() && s.PunktParameters.SentStarters.items[nextTyp] != 0 {
			tokOne.SentBreak = true
			fmt.Println("REASON ABBR WITH SENTENCE STARTER")
			return
		}
	}

	/*
		[4.3. Token-Based Detection of Initials and Ordinals]
		Check if any initials or ordinals tokens that are marked
		as sentbreaks should be reclassified as abbreviations.
	*/
	if tokIsInitial || typ == "##number##" {
		isSentStarter := s.orthoHeuristic(tokTwo)

		if isSentStarter == 0 {
			tokOne.SentBreak = false
			tokOne.Abbr = true
			if tokIsInitial {
				fmt.Println("REASON INITIAL WITH ORTHO HEURISTIC")
				return
			} else {
				fmt.Println("REASON NUMBER WITH ORTHO HEURISTIC")
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
			s.PunktParameters.OrthoContext.items[nextTyp]&orthoLc == 0 {

			tokOne.SentBreak = false
			tokOne.Abbr = true
			fmt.Println("REASON INITIAL WITH SPECIAL ORTHO HEURISTIC")
			return
		}
	}
}

/*
Decide whether the given token is the first token in a sentence.
*/
func (s *SentenceTokenizer) orthoHeuristic(token *PunktToken) int {
	if token == nil {
		return 0
	}

	for _, punct := range s.Punctuation {
		if token.Tok == punct {
			return 0
		}
	}

	orthoCtx := s.PunktParameters.OrthoContext.items[token.TypeNoSentPeriod()]

	/*
	   If the word is capitalized, occurs at least once with a
	   lower case first letter, and never occurs with an upper case
	   first letter sentence-internally, then it's a sentence starter.
	*/
	fmt.Println(orthoCtx & orthoLc)
	if token.FirstUpper() && (orthoCtx&orthoLc > 0 || orthoCtx&orthoMidLc == 0) {
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
