package punkt

// Stores data used to perform sentence boundary detection with punkt
type PunktParameters struct {
	abbrevTypes  map[string]string
	collocations map[string]string
	sentStarters map[string]string
	orthoContext map[string]int
}

func (p *PunktParameters) addOrthoContext(typ string, flag int) {
	p.orthoContext[typ] |= flag
}

func (p *PunktParameters) IsAbbr(tokens ...string) bool {
	for _, token := range tokens {
		if p.abbrevTypes[token] != "" {
			return true
		}
	}
	return false
}
