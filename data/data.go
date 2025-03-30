package data

import (
	"fmt"
	"iter"
	"slices"
	"sort"
	"strings"
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

type Rarity struct {
	order    uint8
	isSecret bool
	value    string
}

const diamondChar = "♢"
const starChar = "☆"
const shinyChar = "✵"
const crownChar = "♕"

var (
	RarityOneDiamond   = Rarity{0, false, strings.Repeat(diamondChar, 1)}
	RarityTwoDiamond   = Rarity{1, false, strings.Repeat(diamondChar, 2)}
	RarityThreeDiamond = Rarity{2, false, strings.Repeat(diamondChar, 3)}
	RarityFourDiamond  = Rarity{3, false, strings.Repeat(diamondChar, 4)}
	RarityOneStar      = Rarity{4, true, strings.Repeat(starChar, 1)}
	RarityTwoStar      = Rarity{5, true, strings.Repeat(starChar, 2)}
	RarityThreeStar    = Rarity{6, true, strings.Repeat(starChar, 3)}
	RarityOneShiny     = Rarity{7, true, strings.Repeat(shinyChar, 1)}
	RarityTwoShiny     = Rarity{8, true, strings.Repeat(shinyChar, 2)}
	RarityCrown        = Rarity{9, true, strings.Repeat(crownChar, 1)}
)

func (r *Rarity) IsStar() bool {
	return strings.Contains(r.value, starChar)
}

func (r *Rarity) IsCrown() bool {
	return strings.Contains(r.value, crownChar)
}

func (r *Rarity) IsShiny() bool {
	return strings.Contains(r.value, shinyChar)
}

type OfferingRatesTable map[*Rarity]BoosterOffering

type BaseCard struct {
	name   string
	health uint8
	// retreatCost
	// moves
	// type
}

func NewBaseCard(name string, health uint8) *BaseCard {
	return &BaseCard{name: name, health: health}
}

func (c *BaseCard) Name() string {
	return c.name
}

type ExpansionNumber uint16

type Card struct {
	core   *BaseCard
	number ExpansionNumber
	rarity *Rarity
}

func NewCard(
	core *BaseCard,
	number ExpansionNumber,
	rarity *Rarity,
) Card {
	return Card{core: core, number: number, rarity: rarity}
}

func (c *Card) Base() *BaseCard {
	return c.core
}

func (c *Card) Name() string {
	return c.core.name
}

func (c *Card) Rarity() *Rarity {
	return c.rarity
}

func (c *Card) Number() ExpansionNumber {
	return c.number
}

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

type Booster struct {
	name                          string
	cards                         []*Card
	offeringRates                 OfferingRatesTable
	crownExclusiveExpansionNumber ExpansionNumber
}

func NewBooster(
	name string,
	cards []*Card,
	offeringRates OfferingRatesTable,
	crownExclusiveExpansionNumber ExpansionNumber,
) Booster {
	return Booster{name: name, cards: cards, offeringRates: offeringRates, crownExclusiveExpansionNumber: crownExclusiveExpansionNumber}
}

func (b *Booster) Name() string {
	return b.name
}

func (b *Booster) Offerings() iter.Seq[*BoosterCardOffering] {
	offerings := make([]*BoosterCardOffering, len(b.cards))
	for i, c := range b.cards {
		offeringRef, offeringRefExists := b.offeringRates[c.Rarity()]
		if !offeringRefExists {
			m, _ := fmt.Printf("Offering rate not found for %v %v", b.name, c.Rarity().value)
			panic(m)
		}
		rareCardOffering := 0.0
		if c.Rarity() != &RarityCrown || c.number == b.crownExclusiveExpansionNumber {
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
	return slices.Values(offerings)
}

type ExpansionId = string

type Expansion struct {
	id                  ExpansionId
	name                string
	boosters            []*Booster
	cards               []*Card
	totalNonSecretCards uint16
	totalSecretCards    uint16
}

func (s *Expansion) Id() ExpansionId {
	return s.id
}

func (s *Expansion) Name() string {
	return s.name
}

func (c *Expansion) Cards() iter.Seq[*Card] {
	return slices.Values(c.cards)
}

func (c *Expansion) Boosters() iter.Seq[*Booster] {
	return slices.Values(c.boosters)
}

func (c *Expansion) TotalNonSecretCards() uint16 {
	return c.totalNonSecretCards
}

func (c *Expansion) TotalSecretCards() uint16 {
	return c.totalSecretCards
}

func (c *Expansion) TotalCards() uint16 {
	return uint16(len(c.cards))
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
