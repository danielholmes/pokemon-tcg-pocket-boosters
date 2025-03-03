package data

import (
	"fmt"
	"iter"
	"slices"
	"sort"
	"strings"
)

type CardSetId = string

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

type Rarity struct {
	order    uint8
	isSecret bool
	value    string
}

var (
	RarityOneDiamond   = Rarity{0, false, "♢"}
	RarityTwoDiamond   = Rarity{1, false, "♢♢"}
	RarityThreeDiamond = Rarity{2, false, "♢♢♢"}
	RarityFourDiamond  = Rarity{3, false, "♢♢♢♢"}
	RarityOneStar      = Rarity{4, true, "☆"}
	RarityTwoStar      = Rarity{5, true, "☆☆"}
	RarityThreeStar    = Rarity{6, true, "☆☆☆"}
	RarityCrown        = Rarity{7, true, "♕"}
)

func (r *Rarity) IsSecret() bool {
	return r.isSecret
}

func (r *Rarity) IsStar() bool {
	return strings.Contains(r.value, "☆")
}

func (r *Rarity) IsCrown() bool {
	return strings.Contains(r.value, "♕")
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

type CardSetNumber uint16

type Card struct {
	core   *BaseCard
	number CardSetNumber
	rarity *Rarity
}

func NewCard(
	core *BaseCard,
	number CardSetNumber,
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

func (c *Card) Number() CardSetNumber {
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
	name                        string
	cards                       []*Card
	offeringRates               OfferingRatesTable
	crownExclusiveCardSetNumber CardSetNumber
}

func NewBooster(
	name string,
	cards []*Card,
	offeringRates OfferingRatesTable,
	crownExclusiveCardSetNumber CardSetNumber,
) Booster {
	return Booster{name: name, cards: cards, offeringRates: offeringRates, crownExclusiveCardSetNumber: crownExclusiveCardSetNumber}
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
		if c.Rarity() != &RarityCrown || c.number == b.crownExclusiveCardSetNumber {
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

type CardSet struct {
	id                  CardSetId
	name                string
	boosters            []*Booster
	cards               []*Card
	totalNonSecretCards uint16
	totalSecretCards    uint16
}

func (s *CardSet) Id() CardSetId {
	return s.id
}

func (s *CardSet) Name() string {
	return s.name
}

func (c *CardSet) Cards() iter.Seq[*Card] {
	return slices.Values(c.cards)
}

func (c *CardSet) Boosters() iter.Seq[*Booster] {
	return slices.Values(c.boosters)
}

func (c *CardSet) TotalNonSecretCards() uint16 {
	return c.totalNonSecretCards
}

func (c *CardSet) TotalSecretCards() uint16 {
	return c.totalSecretCards
}

func (c *CardSet) TotalCards() uint16 {
	return uint16(len(c.cards))
}

func NewCardSet(
	id CardSetId,
	name string, boosters []*Booster) CardSet {
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

	return CardSet{
		id:                  id,
		name:                name,
		boosters:            boosters,
		cards:               cards,
		totalSecretCards:    totalSecretCards,
		totalNonSecretCards: uint16(len(cards)) - totalSecretCards,
	}
}
