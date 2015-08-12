package punkt

import (
	"fmt"
)

// A sentence tokenizer which uses an unsupervised algorithm to build a model
// for abbreviation words, collocations, and words that start sentences
// and then uses that model to find sentence boundaries.
type SentenceTokenizer struct {
	*PunktParameters
	*PunktLanguageVars
}

func NewSentenceTokenizer(trainedData *PunktParameters) *SentenceTokenizer {
	return &SentenceTokenizer{
		PunktParameters:   trainedData,
		PunktLanguageVars: NewPunktLanguageVars(),
	}
}

func (s *SentenceTokenizer) Tokenize(text string) []string {
	//last_break := 0
	matches := s.PunktLanguageVars.PeriodContext(text)
	fmt.Println(len(matches))
	/*for _, match := range matches {
		fmt.Println(match)
	}*/
	return matches
}
