package punkt

import (
	"encoding/json"
	"github.com/neurosnap/sentences/utils"
)

type tmpStorage struct {
	AbbrevTypes  map[string]int
	Collocations map[string]int
	SentStarters map[string]int
	OrthoContext map[string]int
}

type tmpDist struct {
	TypeDist        map[string]int
	CollocationDist map[string]int
	SentStarterDist map[string]int
}

func LoadStorage(data []byte) (*Storage, error) {
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
