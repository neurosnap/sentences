package punkt

type SetString struct {
	items map[string]int
}

func (ss *SetString) Add(str string) {
	ss.items[str] = 1
}

func (ss *SetString) Remove(str string) {
	delete(ss.items, str)
}

func (ss *SetString) Has(str string) bool {
	if ss.items[str] == 0 {
		return false
	} else {
		return true
	}
}

type SetInt struct {
	items map[int]int
}

func (si *SetInt) Add(val int) {
	si.items[val] = 1
}

func (si *SetInt) Remove(val int) {
	delete(si.items, val)
}

// Stores data used to perform sentence boundary detection with punkt
type PunktParameters struct {
	AbbrevTypes  *SetString
	Collocations *SetString
	SentStarters *SetString
	OrthoContext *SetString
}

func (p *PunktParameters) addOrthoContext(typ string, flag int) {
	p.OrthoContext.items[typ] |= flag
}

func (p *PunktParameters) IsAbbr(tokens ...string) bool {
	for _, token := range tokens {
		if p.AbbrevTypes.Has(token) {
			return true
		}
	}
	return false
}
