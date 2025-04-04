package data

import (
	"fmt"
	"iter"
	"slices"
	"sort"
)

type ExpansionId = string

type Expansion struct {
	id                  ExpansionId
	name                string
	boosters            []*Booster
	cards               []*Card
	totalNonSecretCards uint16
	totalSecretCards    uint16
}

func NewExpansion(
	id ExpansionId,
	name string, boosters []*Booster) Expansion {
	var cards []*Card
	for _, b := range boosters {
		for _, c := range b.cards {
			if !slices.ContainsFunc(cards, func(c2 *Card) bool { return c2.number == c.number }) {
				cards = append(cards, c)
			}
		}
	}
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Number() < cards[j].Number()
	})

	var totalSecretCards uint16 = 0
	for _, c := range cards {
		if c.rarity.isSecret {
			totalSecretCards += 1
		}
	}

	return Expansion{
		id:                  id,
		name:                name,
		boosters:            boosters,
		cards:               cards,
		totalSecretCards:    totalSecretCards,
		totalNonSecretCards: uint16(len(cards)) - totalSecretCards,
	}
}

func (e *Expansion) Id() ExpansionId {
	return e.id
}

func (e *Expansion) Name() string {
	return e.name
}

func (e *Expansion) Cards() iter.Seq[*Card] {
	return slices.Values(e.cards)
}

func (e *Expansion) Boosters() iter.Seq[*Booster] {
	return slices.Values(e.boosters)
}

func (e *Expansion) TotalNonSecretCards() uint16 {
	return e.totalNonSecretCards
}

func (e *Expansion) TotalSecretCards() uint16 {
	return e.totalSecretCards
}

func (e *Expansion) TotalCards() uint16 {
	return uint16(len(e.cards))
}

func (e *Expansion) GetCardByNumber(number ExpansionNumber) (*Card, error) {
	cIndex := slices.IndexFunc(e.cards, func(c *Card) bool {
		return c.number == number
	})
	if cIndex == -1 {
		return nil, fmt.Errorf("no card with number %v", number)
	}
	return e.cards[cIndex], nil
}

func (e *Expansion) GetHighestOfferingBoosterForMissingCards(
	missingCardNumbers []ExpansionNumber,
) (*Booster, error) {
	if len(missingCardNumbers) == 0 {
		return nil, fmt.Errorf("no missing card numbers provided")
	}

	if len(e.boosters) == 1 {
		return e.boosters[0], nil
	}

	var bestBooster *Booster
	var bestBoosterProbability = -1.0
	for b := range e.Boosters() {
		boosterProbability := b.GetInstanceProbabilityForMissing(missingCardNumbers)
		if boosterProbability > bestBoosterProbability {
			bestBoosterProbability = boosterProbability
			bestBooster = b
		}
	}

	if bestBoosterProbability == 0.0 {
		return nil, fmt.Errorf("no booster offering any card number")
	}

	return bestBooster, nil
}
