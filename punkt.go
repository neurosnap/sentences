package main

// Stores data used to perform sentence boundary detection with punkt
type PunktParameters struct {
	abbrev_types  map[string]string
	collocations  map[string]string
	sent_starters map[string]string
	ortho_context map[string]int
}

func (p *PunktParameters) addOrthoContext(typ string, flag int) {
	p.ortho_context[typ] |= flag
}

// Stores a token of text with annotations produced during sentence boundary detection
type PunktToken struct {
}

// Includes common components of PunkTrainer and PunktSentenceTokenizer
type PunktBase struct{}

// Learns parameters used in Punkt sentence boundary detection
type PunktTrainer struct{}

// A sentence tokenizer which uses an unsupervised algorithm to build a model
// for abbreviation words, collocations, and words that start sentences
// and then uses that model to find sentence boundaries.
type PunktSentenceTokenizer struct{}
