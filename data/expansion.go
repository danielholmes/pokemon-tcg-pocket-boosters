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
	boosters            iter.Seq[*Booster]
	cards               iter.Seq[*Card]
	totalNonSecretCards uint16
	totalSecretCards    uint16
}

func NewExpansion(
	id ExpansionId,
	name string,
	boosters []*Booster,
) *Expansion {
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

	return &Expansion{
		id:                  id,
		name:                name,
		boosters:            slices.Values(boosters),
		cards:               slices.Values(cards),
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

func (e *Expansion) HasShiny() bool {
	for c := range e.cards {
		if c.Rarity().IsShiny() {
			return true
		}
	}
	return false
}

func (e *Expansion) Cards() iter.Seq[*Card] {
	return e.cards
}

func (e *Expansion) Boosters() iter.Seq[*Booster] {
	return e.boosters
}

func (e *Expansion) TotalNonSecretCards() uint16 {
	return e.totalNonSecretCards
}

func (e *Expansion) TotalSecretCards() uint16 {
	return e.totalSecretCards
}

func (e *Expansion) TotalCards() uint16 {
	return e.totalNonSecretCards + e.totalSecretCards
}

// TODO this can be more efficient with an array. Need to sort out
// special case mew card first though.
func (e *Expansion) GetCardByNumber(number ExpansionCardNumber) (*Card, error) {
	for c := range e.cards {
		if c.number == number {
			return c, nil
		}
	}
	return nil, fmt.Errorf("no card with number %v", number)
}

func (e *Expansion) GetHighestOfferingBoosterForMissingCards(
	missingCards []*Card,
) (*Booster, error) {
	if len(missingCards) == 0 {
		return nil, fmt.Errorf("no missing card numbers provided")
	}

	// Can't do this optimisation with seq. Worth adding another way to do it?
	// if len(e.boosters) == 1 {
	// 	return e.boosters[0], nil
	// }

	var bestBooster *Booster
	var bestBoosterProbability = -1.0
	for b := range e.Boosters() {
		boosterProbability := b.GetInstanceProbabilityForMissing(missingCards)
		if boosterProbability > bestBoosterProbability {
			bestBoosterProbability = boosterProbability
			bestBooster = b
		}
	}

	if bestBoosterProbability <= 0.0 {
		return nil, fmt.Errorf("no booster offering any card number")
	}

	return bestBooster, nil
}
