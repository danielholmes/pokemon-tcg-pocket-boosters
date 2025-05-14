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

const MaxPackPointsPerBooster uint16 = 2_500

const regularPackRate = 0.9995
const rarePackRate = 1.0 - regularPackRate

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

type BoosterInstance struct {
	isRare bool
	cards  iter.Seq[*Card]
}

func NewBoosterInstance(isRare bool, cards [5]*Card) *BoosterInstance {
	return &BoosterInstance{isRare: isRare, cards: slices.Values(cards[:])}
}

func (b *BoosterInstance) IsRare() bool {
	return b.isRare
}

func (b *BoosterInstance) Cards() iter.Seq[*Card] {
	return b.cards
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

func (o *offeringProbabilityList) pickRandomCard(randomGenerator *rand.Rand) *Card {
	num := randomGenerator.Float64() * o.totalProbability
	for _, e := range o.entries {
		if num <= e.cumulativeProbability {
			return e.card
		}
	}
	panic(fmt.Sprintf("Invalid algorithm %v num %v total", num, o.totalProbability))
}

type Booster struct {
	name                string
	cards               []*Card
	offerings           iter.Seq[*BoosterCardOffering]
	regularPack1To3List *offeringProbabilityList
	regularPack4List    *offeringProbabilityList
	regularPack5List    *offeringProbabilityList
	rarePackList        *offeringProbabilityList
}

func NewBooster(
	name string,
	cards []*Card,
	offeringRates OfferingRatesTable,
	rarePackCrownExclusiveExpansionNumber ExpansionCardNumber,
) *Booster {
	offerings := make([]*BoosterCardOffering, len(cards))
	regularPack1To3List := offeringProbabilityList{}
	regularPack4List := offeringProbabilityList{}
	regularPack5List := offeringProbabilityList{}
	rarePackList := offeringProbabilityList{}

	cardsByRarity := make(map[*Rarity]uint16)

	for _, c := range cards {
		cardsByRarity[c.Rarity()] += 1
	}

	for i, c := range cards {
		offeringRef, offeringRefExists := offeringRates[c.Rarity()]
		if !offeringRefExists {
			m, _ := fmt.Printf("Offering rate not found for %v %v", name, c.Rarity().value)
			panic(m)
		}

		rareCardOffering := 0.0
		numOfRarity := float64(cardsByRarity[c.Rarity()])
		if c.Rarity() != RarityCrown || c.number == rarePackCrownExclusiveExpansionNumber {
			rareCardOffering = offeringRef.rareOffering
		} else if c.number == rarePackCrownExclusiveExpansionNumber {
			numOfRarity -= 1
		}

		offerings[i] = &BoosterCardOffering{
			card:               c,
			first3CardOffering: offeringRef.first3CardOffering / numOfRarity,
			fourthCardOffering: offeringRef.fourthCardOffering / numOfRarity,
			fifthCardOffering:  offeringRef.fifthCardOffering / numOfRarity,
			rareCardOffering:   rareCardOffering / numOfRarity,
		}

		regularPack1To3List.append(c, offeringRef.first3CardOffering)
		regularPack4List.append(c, offeringRef.fourthCardOffering)
		regularPack5List.append(c, offeringRef.fifthCardOffering)
		rarePackList.append(c, offeringRef.rareOffering)
	}

	return &Booster{
		name:                name,
		cards:               cards,
		offerings:           slices.Values(offerings),
		regularPack1To3List: &regularPack1To3List,
		regularPack4List:    &regularPack4List,
		regularPack5List:    &regularPack5List,
		rarePackList:        &rarePackList,
	}
}

func (b *Booster) Name() string {
	return b.name
}

func (b *Booster) Offerings() iter.Seq[*BoosterCardOffering] {
	return b.offerings
}

func (b *Booster) GetInstanceProbabilityForMissing(missing []*Card) float64 {
	totalRegularPackOffering := 0.0
	totalRarePackOffering := 0.0
	for o := range b.Offerings() {
		if slices.Contains(missing, o.Card()) {
			totalRegularPackOffering += o.RegularPackOffering()
			totalRarePackOffering += o.RarePackOffering()
		}
	}
	return totalRegularPackOffering*regularPackRate + totalRarePackOffering*rarePackRate
}

func (b *Booster) CreateRandomInstance(randomGenerator *rand.Rand) *BoosterInstance {
	// Rare pack
	rareNum := randomGenerator.Float64()
	if rareNum < rarePackRate {
		return NewBoosterInstance(
			true,
			[5]*Card{
				b.rarePackList.pickRandomCard(randomGenerator),
				b.rarePackList.pickRandomCard(randomGenerator),
				b.rarePackList.pickRandomCard(randomGenerator),
				b.rarePackList.pickRandomCard(randomGenerator),
				b.rarePackList.pickRandomCard(randomGenerator),
			})
	}

	// Regular pack
	return NewBoosterInstance(
		false,
		[5]*Card{
			b.regularPack1To3List.pickRandomCard(randomGenerator),
			b.regularPack1To3List.pickRandomCard(randomGenerator),
			b.regularPack1To3List.pickRandomCard(randomGenerator),
			b.regularPack4List.pickRandomCard(randomGenerator),
			b.regularPack5List.pickRandomCard(randomGenerator),
		})
}
