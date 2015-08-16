package main

import (
	"encoding/json"
	"github.com/neurosnap/go-sentences/punkt"
)

type tmpParams struct {
	AbbrevTypes  map[string]int
	Collocations map[string]int
	SentStarters map[string]int
	OrthoContext map[string]int
}

func Load(data []byte) (*punkt.Storage, error) {
	var p tmpParams
	err := json.Unmarshal(data, &p)

	if err != nil {
		return nil, err
	}

	params := &punkt.Storage{
		AbbrevTypes:  punkt.NewSetString(p.AbbrevTypes),
		Collocations: punkt.NewSetString(p.Collocations),
		SentStarters: punkt.NewSetString(p.SentStarters),
		OrthoContext: punkt.NewSetString(p.OrthoContext),
	}

	return params, nil
}
