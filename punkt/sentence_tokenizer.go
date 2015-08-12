package punkt

import (
	"fmt"
	//	"strings"
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
	re := s.PunktLanguageVars.RePeriodContext()
	fmt.Println(re)
	matches := re.FindAllStringSubmatchIndex(text, -1)

	for _, match := range matches {
		fmt.Println("context: ", text[match[0]:match[1]])
		fmt.Println("next_tok: ", text[match[4]:match[5]])
		fmt.Println("start: ", match[2])
		fmt.Println("end: ", match[4])
		fmt.Println("-------")
	}
	return []string{}
}
