package punkt

import "encoding/json"

// golang implementation of a set of strings
// probably not the best way to do this but oh well
type SetString map[string]int

func (ss SetString) Add(str string) {
	ss[str] = 1
}

func (ss SetString) Remove(str string) {
	delete(ss, str)
}

func (ss SetString) Has(str string) bool {
	if ss[str] == 0 {
		return false
	} else {
		return true
	}
}

func (ss SetString) Array() []string {
	arr := make([]string, 0, len(ss))

	for key := range ss {
		arr = append(arr, key)
	}

	return arr
}

// Stores data used to perform sentence boundary detection with punkt
// This is where all the training data gets stored for future use
type Storage struct {
	AbbrevTypes  SetString `json:"AbbrevTypes"`
	Collocations SetString `json:"Collocations"`
	SentStarters SetString `json:"SentStarters"`
	OrthoContext SetString `json:"OrthoContext"`
}

// Primary function to load JSON training data.  By default, the sentence tokenizer
// loads in english automatically, but other languages could be loaded into a
// binary file using the `make <lang>` command.
func LoadTraining(data []byte) (*Storage, error) {
	var storage Storage
	err := json.Unmarshal(data, &storage)

	if err != nil {
		return nil, err
	}

	return &storage, nil
}

// Creates the default storage container
func NewStorage() *Storage {
	return &Storage{SetString{}, SetString{}, SetString{}, SetString{}}
}

// Used in the training to add a type to the ortho context
func (p *Storage) addOrthoContext(typ string, flag int) {
	p.OrthoContext[typ] |= flag
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
