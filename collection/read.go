package collection

import (
	"encoding/json"
	"fmt"
	"os"
	"ptcgpocket/data"
	"slices"
)

type serialisedExpansionCollection struct {
	Missing    []data.ExpansionCardNumber `json:"missing"`
	PackPoints uint16                     `json:"packPoints"`
}

func ReadFromFilepath(filepath string, expansions []*data.Expansion) (*UserCollection, error) {
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
		eIndex := slices.IndexFunc(expansions, func(e *data.Expansion) bool {
			return e.Id() == i
		})
		if eIndex == -1 {
			panic(fmt.Sprintf("expansion id %v unrecognised", i))
		}

		e := expansions[eIndex]
		missingCards := make([]*data.Card, len(s.Missing))
		for i, m := range s.Missing {
			c, cErr := e.GetCardByNumber(m)
			if cErr != nil {
				panic(cErr)
			}
			missingCards[i] = c
		}
		expansionCollections[i] = &ExpansionCollection{
			missingCards: missingCards,
			packPoints:   s.PackPoints,
		}
	}

	return NewUserCollection(expansionCollections), nil
}
