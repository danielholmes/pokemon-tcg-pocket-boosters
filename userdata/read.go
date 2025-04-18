package userdata

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

type serialisedUserData struct {
	Collection map[data.ExpansionId]*serialisedExpansionCollection        `json:"collection"`
	Wishlists  map[string]map[data.ExpansionId][]data.ExpansionCardNumber `json:"wishlists"`
}

func ReadFromFilepath(filepath string, expansions []*data.Expansion) (*UserData, error) {
	raw, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var serialisedUserData serialisedUserData
	uErr := json.Unmarshal(raw, &serialisedUserData)
	if uErr != nil {
		return nil, uErr
	}

	expansionCollections := make(map[data.ExpansionId]*ExpansionCollection, len(serialisedUserData.Collection))
	for i, s := range serialisedUserData.Collection {
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

	wishlists := make([]*Wishlist, len(serialisedUserData.Wishlists))
	i := 0
	for n, s := range serialisedUserData.Wishlists {
		expansionWishlists := make(map[data.ExpansionId]*ExpansionWishlist, len(s))
		for eId, m := range s {
			eIndex := slices.IndexFunc(expansions, func(e *data.Expansion) bool {
				return e.Id() == eId
			})
			if eIndex == -1 {
				panic(fmt.Sprintf("expansion id %v unrecognised", i))
			}

			e := expansions[eIndex]
			cards := make([]*data.Card, len(m))
			for i, cardNumber := range m {
				c, cErr := e.GetCardByNumber(cardNumber)
				if cErr != nil {
					panic(cErr)
				}
				cards[i] = c
			}

			expansionWishlists[e.Id()] = &ExpansionWishlist{
				cards: cards,
			}
		}
		wishlists[i] = &Wishlist{
			name:       n,
			expansions: expansionWishlists,
		}
		i++
	}

	return NewUserData(NewUserCollection(expansionCollections), wishlists), nil
}
