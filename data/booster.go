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

type cardProbabilityEntry struct {
	cumulativeProbability float64
	card                  *Card
}

type offeringProbabilityList struct {
	// Should always be 100, but atm we have some odd data where it's less.
	// Should solve that core issue one day.
	totalProbability float64
	entries          []*cardProbabilityEntry
}

func (o *offeringProbabilityList) append(card *Card, probability float64) {
	if probability == 0 {
		return
	}

	o.totalProbability += probability
	o.entries = append(o.entries, &cardProbabilityEntry{
		cumulativeProbability: o.totalProbability,
		card:                  card,
	})
}

func (o *offeringProbabilityList) pickRandomCard() *Card {
	num := rand.Float64() * o.totalProbability
	for _, e := range o.entries {
		if num <= e.cumulativeProbability {
			return e.card
		}
	}
	panic(fmt.Sprintf("Invalid algorithm %v num %v total", num, o.totalProbability))
}

type Booster struct {
	name                          string
	cards                         []*Card
	crownExclusiveExpansionNumber ExpansionNumber
	offerings                     iter.Seq[*BoosterCardOffering]
	regularPack1To3List           *offeringProbabilityList
	regularPack4List              *offeringProbabilityList
	regularPack5List              *offeringProbabilityList
	rarePackList                  *offeringProbabilityList
}

func NewBooster(
	name string,
	cards []*Card,
	offeringRates OfferingRatesTable,
	crownExclusiveExpansionNumber ExpansionNumber,
) Booster {
	offerings := make([]*BoosterCardOffering, len(cards))
	regularPack1To3List := offeringProbabilityList{}
	regularPack4List := offeringProbabilityList{}
	regularPack5List := offeringProbabilityList{}
	rarePackList := offeringProbabilityList{}
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

		regularPack1To3List.append(c, offeringRef.first3CardOffering)
		regularPack4List.append(c, offeringRef.fourthCardOffering)
		regularPack5List.append(c, offeringRef.fifthCardOffering)
		rarePackList.append(c, offeringRef.rareOffering)
	}

	return Booster{
		name:                          name,
		cards:                         cards,
		crownExclusiveExpansionNumber: crownExclusiveExpansionNumber,
		offerings:                     slices.Values(offerings),
		regularPack1To3List:           &regularPack1To3List,
		regularPack4List:              &regularPack4List,
		regularPack5List:              &regularPack5List,
		rarePackList:                  &rarePackList,
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
	// Rare pack
	if rand.Float64() < RarePackRate {
		return NewBoosterInstance([5]*Card{
			b.rarePackList.pickRandomCard(),
			b.rarePackList.pickRandomCard(),
			b.rarePackList.pickRandomCard(),
			b.rarePackList.pickRandomCard(),
			b.rarePackList.pickRandomCard(),
		})
	}

	// Regular pack
	return NewBoosterInstance([5]*Card{
		b.regularPack1To3List.pickRandomCard(),
		b.regularPack1To3List.pickRandomCard(),
		b.regularPack1To3List.pickRandomCard(),
		b.regularPack4List.pickRandomCard(),
		b.regularPack5List.pickRandomCard(),
	})
}
