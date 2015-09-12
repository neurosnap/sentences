package punkt

// golang implementation of a set, probably not the best way to do this
// but oh well
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
// This is where all the training data gets stored for future use
type Storage struct {
	AbbrevTypes  *SetString
	Collocations *SetString
	SentStarters *SetString
	OrthoContext *SetString
}

// Creates the default storage container
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

// Detemins if any of the tokens are an abbreviation
func (p *Storage) IsAbbr(tokens ...string) bool {
	for _, token := range tokens {
		if p.AbbrevTypes.Has(token) {
			return true
		}
	}
	return false
}
