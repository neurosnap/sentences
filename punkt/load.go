package punkt

import (
	"encoding/json"
)

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
