package collection

import (
	"encoding/json"
	"fmt"
	"os"
	"ptcgpocket/data"
)

type serialisedExpansionCollection struct {
	Missing    []data.ExpansionCardNumber `json:"missing"`
	PackPoints uint16                     `json:"packPoints"`
}

func ReadFromFilepath(filepath string) (*UserCollection, error) {
	raw, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var serialisedCollections map[data.ExpansionId]*serialisedExpansionCollection
	uErr := json.Unmarshal(raw, &serialisedCollections)
	if uErr != nil {
		return nil, uErr
	}

	expansionCollections := make(map[data.ExpansionId]*ExpansionCollection, len(serialisedCollections))
	for i, s := range serialisedCollections {
		fmt.Printf("%v: %v\n", i, s)
		expansionCollections[i] = &ExpansionCollection{
			missingCardNumbers: s.Missing,
			packPoints:         s.PackPoints,
		}
	}

	return NewUserCollection(expansionCollections), nil
}
