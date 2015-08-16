package punkt

type SetString struct {
	items map[string]int
}

func NewSetString(items map[string]int) *SetString {
	return &SetString{items}
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

func (ss *SetString) Array() []string {
	arr := make([]string, 0, len(ss.items))

	for key := range ss.items {
		arr = append(arr, key)
	}

	return arr
}

// Stores data used to perform sentence boundary detection with punkt
type Storage struct {
	AbbrevTypes  *SetString
	Collocations *SetString
	SentStarters *SetString
	OrthoContext *SetString
}

func NewStorage() *Storage {
	return &Storage{
		NewSetString(nil),
		NewSetString(nil),
		NewSetString(nil),
		NewSetString(nil),
	}
}

func (p *Storage) addOrthoContext(typ string, flag int) {
	p.OrthoContext.items[typ] |= flag
}

func (p *Storage) IsAbbr(tokens ...string) bool {
	for _, token := range tokens {
		if p.AbbrevTypes.Has(token) {
			return true
		}
	}
	return false
}
