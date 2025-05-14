package userdata

import (
	"fmt"
	"iter"
	"maps"
	"ptcgpocket/data"
	"slices"
)

const packPointsPerBooster = 5

type ExpansionCollection struct {
	// Has a max of 2,500
	packPoints   uint16
	missingCards []*data.Card
}

func (c *ExpansionCollection) PackPoints() uint16 {
	return c.packPoints
}

func (c *ExpansionCollection) AcquireCardUsingPackPoints(
	card *data.Card,
) {
	previousCardsLength := len(c.missingCards)
	c.missingCards = slices.DeleteFunc(c.missingCards, func(n *data.Card) bool {
		return n == card
	})
	if previousCardsLength != len(c.missingCards)+1 {
		panic("Card not missing")
	}
	if c.packPoints < card.Rarity().PackPointsToObtain() {
		panic(
			fmt.Sprintf(
				"Not enough pack points to obtain (%v), require (%v)",
				c.packPoints,
				card.Rarity().PackPointsToObtain(),
			),
		)
	}
	c.packPoints -= card.Rarity().PackPointsToObtain()
}

func (c *ExpansionCollection) AcquireCardsFromBooster(
	added iter.Seq[*data.Card],
) {
	c.missingCards = slices.DeleteFunc(c.missingCards, func(c *data.Card) bool {
		for a := range added {
			if a == c {
				return true
			}
		}
		return false
	})
	c.packPoints = min(c.packPoints+packPointsPerBooster, data.MaxPackPointsPerBooster)
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
		packPoints:   c.packPoints,
		missingCards: slices.Clone(c.missingCards),
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
		if len(eC.missingCards) > 0 {
			return eId, nil
		}
	}
	return "", fmt.Errorf("no incomplete expansion")
}

func (c *UserCollection) MissingForExpansion(expansionId data.ExpansionId) ([]*data.Card, bool) {
	v, e := c.expansions[expansionId]
	if v == nil {
		return nil, e
	}
	return v.missingCards, true
}
