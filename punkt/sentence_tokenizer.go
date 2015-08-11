package punkt

// A sentence tokenizer which uses an unsupervised algorithm to build a model
// for abbreviation words, collocations, and words that start sentences
// and then uses that model to find sentence boundaries.
type SentenceTokenizer struct {
	*PunktParameters
}

func NewSentenceTokenizer(trainedData *PunktParameters) *SentenceTokenizer {
	return &SentenceTokenizer{
		PunktParameters: trainedData,
	}
}
