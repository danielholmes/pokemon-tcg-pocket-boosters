package data

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"
)

type BoosterOffering struct {
	first3CardOffering float64
	fourthCardOffering float64
	fifthCardOffering  float64
	rareOffering       float64
}

func NewBoosterOffering(
	first3CardOffering float64,
	fourthCardOffering float64,
	fifthCardOffering float64,
	rareOffering float64,
) *BoosterOffering {
	return &BoosterOffering{
		first3CardOffering: first3CardOffering,
		fourthCardOffering: fourthCardOffering,
		fifthCardOffering:  fifthCardOffering,
		rareOffering:       rareOffering,
	}
}

var NotPresentBoosterOffering = NewBoosterOffering(0, 0, 0, 0)

const RegularPackRate = 0.9995
const RarePackRate = 1.0 - RegularPackRate

type BoosterCardOffering struct {
	card               *Card
	first3CardOffering float64
	fourthCardOffering float64
	fifthCardOffering  float64
	rareCardOffering   float64
}

func (b *BoosterCardOffering) Card() *Card {
	return b.card
}

func (b *BoosterCardOffering) First3CardOffering() float64 {
	return b.first3CardOffering
}

func (b *BoosterCardOffering) FourthCardOffering() float64 {
	return b.fourthCardOffering
}

func (b *BoosterCardOffering) FifthCardOffering() float64 {
	return b.fifthCardOffering
}

func (b *BoosterCardOffering) RareCardOffering() float64 {
	return b.rareCardOffering
}

func (b *BoosterCardOffering) RarePackOffering() float64 {
	return b.rareCardOffering * 5
}

func (b *BoosterCardOffering) RegularPackOffering() float64 {
	return b.first3CardOffering*3 + b.fourthCardOffering + b.fifthCardOffering
}

func (b *BoosterCardOffering) OverallPackOffering() float64 {
	return b.RegularPackOffering()*RegularPackRate + b.RarePackOffering()*RarePackRate
}

type BoosterInstance struct {
	cards [5]*Card
}

func NewBoosterInstance(cards [5]*Card) *BoosterInstance {
	return &BoosterInstance{cards: cards}
}

func (b *BoosterInstance) CardNumbers() [5]ExpansionNumber {
	var numbers [5]ExpansionNumber
	for i, c := range b.cards {
		numbers[i] = c.Number()
	}
	return numbers
}

type OfferingRatesTable map[*Rarity]BoosterOffering

type Booster struct {
	name                          string
	cards                         []*Card
	offeringRates                 OfferingRatesTable
	crownExclusiveExpansionNumber ExpansionNumber
	offerings                     iter.Seq[*BoosterCardOffering]
}

func NewBooster(
	name string,
	cards []*Card,
	offeringRates OfferingRatesTable,
	crownExclusiveExpansionNumber ExpansionNumber,
) Booster {
	offerings := make([]*BoosterCardOffering, len(cards))
	for i, c := range cards {
		offeringRef, offeringRefExists := offeringRates[c.Rarity()]
		if !offeringRefExists {
			m, _ := fmt.Printf("Offering rate not found for %v %v", name, c.Rarity().value)
			panic(m)
		}

		rareCardOffering := 0.0
		if c.Rarity() != &RarityCrown || c.number == crownExclusiveExpansionNumber {
			rareCardOffering = offeringRef.rareOffering
		}

		offerings[i] = &BoosterCardOffering{
			card:               c,
			first3CardOffering: offeringRef.first3CardOffering,
			fourthCardOffering: offeringRef.fourthCardOffering,
			fifthCardOffering:  offeringRef.fifthCardOffering,
			rareCardOffering:   rareCardOffering,
		}
	}

	return Booster{
		name:                          name,
		cards:                         cards,
		offeringRates:                 offeringRates,
		crownExclusiveExpansionNumber: crownExclusiveExpansionNumber,
		offerings:                     slices.Values(offerings),
	}
}

func (b *Booster) Name() string {
	return b.name
}

func (b *Booster) Offerings() iter.Seq[*BoosterCardOffering] {
	return b.offerings
}

func (b *Booster) GetInstanceProbabilityForMissing(missing []ExpansionNumber) float64 {
	totalOfferingMissing := 0.0
	for o := range b.Offerings() {
		if slices.Contains(missing, o.Card().Number()) {
			totalOfferingMissing += o.OverallPackOffering()
		}
	}
	return totalOfferingMissing
}

func (b *Booster) CreateRandomInstance() *BoosterInstance {
	var cards [5]*Card

	// Rare pack
	if rand.Float64() < RarePackRate {
		// TODO:
	}

	// Regular pack
	card1Rand := rand.Float64() * 100.0
	card2Rand := rand.Float64() * 100.0
	card3Rand := rand.Float64() * 100.0
	card4Rand := rand.Float64() * 100.0
	card5Rand := rand.Float64() * 100.0
	for o := range b.Offerings() {
		if cards[0] == nil {
			if card1Rand <= o.first3CardOffering {
				cards[0] = o.card
			} else {
				card1Rand -= o.first3CardOffering
			}
		}
		if cards[1] == nil {
			if card2Rand <= o.first3CardOffering {
				cards[1] = o.card
			} else {
				card2Rand -= o.first3CardOffering
			}
		}
		if cards[2] == nil {
			if card3Rand <= o.first3CardOffering {
				cards[2] = o.card
			} else {
				card3Rand -= o.first3CardOffering
			}
		}
		if cards[3] == nil {
			if card4Rand <= o.fourthCardOffering {
				cards[3] = o.card
			} else {
				card4Rand -= o.fourthCardOffering
			}
		}
		if cards[4] == nil {
			if card5Rand <= o.fifthCardOffering {
				cards[4] = o.card
			} else {
				card5Rand -= o.fifthCardOffering
			}
		}
	}

	// We have some percentages that don't match up to 100 exactly. Leaving some gaps.
	// Ideally we fix by making a;; add up to 100, but for now just shove in the first offering.
	for i, c := range cards {
		if c == nil {
			for o := range b.Offerings() {
				cards[i] = o.card
				break
			}
		}
	}

	return NewBoosterInstance(cards)
}
