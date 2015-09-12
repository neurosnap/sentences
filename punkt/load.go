package punkt

import (
	"encoding/json"
	"github.com/neurosnap/sentences/utils"
)

// Temporary storage container used to load in training data
type tmpStorage struct {
	AbbrevTypes  map[string]int
	Collocations map[string]int
	SentStarters map[string]int
	OrthoContext map[string]int
}

// Temporary frequency distribution data used to load in training data
type tmpDist struct {
	TypeDist        map[string]int
	CollocationDist map[string]int
	SentStarterDist map[string]int
}

// Primary function to load JSON training data.  By default, the sentence tokenizer
// loads in english automatically, but other languages could be loaded into a
// binary file using the `make <lang>` command.
func LoadTraining(data []byte) (*Storage, error) {
	var s tmpStorage
	err := json.Unmarshal(data, &s)

	if err != nil {
		return nil, err
	}

	storage := &Storage{
		AbbrevTypes:  NewSetString(s.AbbrevTypes),
		Collocations: NewSetString(s.Collocations),
		SentStarters: NewSetString(s.SentStarters),
		OrthoContext: NewSetString(s.OrthoContext),
	}

	return storage, nil
}

// Loads frequency distribution data from JSON
func LoadFreqDist(data []byte) (*Trainer, error) {
	var tn tmpDist
	err := json.Unmarshal(data, &tn)

	if err != nil {
		return nil, err
	}

	trainer := NewTrainer("", nil)
	trainer.TypeDist = utils.NewFreqDist(tn.TypeDist)
	trainer.CollocationDist = utils.NewFreqDist(tn.CollocationDist)
	trainer.SentStarterDist = utils.NewFreqDist(tn.SentStarterDist)

	return trainer, nil
}
