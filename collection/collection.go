package collection

import (
	"fmt"
	"maps"
	"ptcgpocket/data"
	"slices"
)

const packPointsPerBooster = 5

type ExpansionCollection struct {
	// Has a max of 2,500
	packPoints         uint16
	missingCardNumbers []data.ExpansionCardNumber
}

func (c *ExpansionCollection) PackPoints() uint16 {
	return c.packPoints
}

func (c *ExpansionCollection) AcquireCardUsingPackPoints(
	card *data.Card,
) {
	c.missingCardNumbers = slices.DeleteFunc(c.missingCardNumbers, func(n data.ExpansionCardNumber) bool {
		return n == card.Number()
	})
	c.packPoints -= card.Rarity().PackPointsToObtain()
}

func (c *ExpansionCollection) AddCardsFromBooster(
	addedNumbers [5]data.ExpansionCardNumber,
) {
	c.missingCardNumbers = slices.DeleteFunc(c.missingCardNumbers, func(n data.ExpansionCardNumber) bool {
		for _, a := range addedNumbers {
			if a == n {
				return true
			}
		}
		return false
	})
	c.packPoints += packPointsPerBooster
}

func (c *ExpansionCollection) NumPackPoints() uint16 {
	return c.packPoints
}

func (c *ExpansionCollection) UsePackPoints(amount uint16) {
	if amount > c.packPoints {
		panic("Trying to take out more pack points than available")
	}
	c.packPoints -= amount
}

func (c *ExpansionCollection) Clone() *ExpansionCollection {
	return &ExpansionCollection{
		packPoints:         c.packPoints,
		missingCardNumbers: slices.Clone(c.missingCardNumbers),
	}
}

type UserCollection struct {
	expansions map[data.ExpansionId]*ExpansionCollection
}

func NewUserCollection(expansions map[data.ExpansionId]*ExpansionCollection) *UserCollection {
	return &UserCollection{expansions: maps.Clone(expansions)}
}

func (c *UserCollection) Clone() *UserCollection {
	newExpansions := make(map[data.ExpansionId]*ExpansionCollection, len(c.expansions))
	for eId, c := range c.expansions {
		newExpansions[eId] = c.Clone()
	}

	return &UserCollection{expansions: newExpansions}
}

func (c *UserCollection) GetExpansionCollection(id data.ExpansionId) *ExpansionCollection {
	return c.expansions[id]
}

func (c *UserCollection) FirstIncompleteExpansionId() (data.ExpansionId, error) {
	for eId, eC := range c.expansions {
		if len(eC.missingCardNumbers) > 0 {
			return eId, nil
		}
	}
	return "", fmt.Errorf("no incomplete expansion")
}

func (c *UserCollection) MissingForExpansion(expansionId data.ExpansionId) ([]data.ExpansionCardNumber, bool) {
	v, e := c.expansions[expansionId]
	if v == nil {
		return nil, e
	}
	return v.missingCardNumbers, true
}
