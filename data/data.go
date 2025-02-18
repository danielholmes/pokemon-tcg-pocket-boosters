package data

import (
	"iter"
	"ptcgpocket/ref"
	"slices"
	"sort"
)

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

type Card struct {
	name   string
	number ref.CardSetNumber
	rarity *Rarity
}

func NewCard(
	name string,
	number ref.CardSetNumber,
	rarity *Rarity,
) Card {
	return Card{name: name, number: number, rarity: rarity}
}

func (c *Card) Rarity() *Rarity {
	return c.rarity
}

func (c *Card) Name() string {
	return c.name
}

func (c *Card) Number() ref.CardSetNumber {
	return c.number
}

type BoosterOffering struct {
	card               Card
	first3CardOffering float64
	fourthCardOffering float64
	fifthCardOffering  float64
	packProbability    float64
}

func (b *BoosterOffering) Card() *Card {
	return &b.card
}

func (b *BoosterOffering) First3CardOffering() float64 {
	return b.first3CardOffering
}

func (b *BoosterOffering) FourthCardOffering() float64 {
	return b.fourthCardOffering
}

func (b *BoosterOffering) FifthCardOffering() float64 {
	return b.fifthCardOffering
}

func (b *BoosterOffering) PackProbability() float64 {
	return b.packProbability
}

type Booster struct {
	name      string
	offerings []BoosterOffering
}

func NewBooster(
	name string,
	offerings []BoosterOffering,
) Booster {
	return Booster{name: name, offerings: offerings}
}

func (b *Booster) Name() string {
	return b.name
}

func (b *Booster) Offerings() iter.Seq[BoosterOffering] {
	return slices.Values(b.offerings)
}

type CardSetDetails struct {
	set                 ref.CardSet
	boosters            []Booster
	cards               []Card
	totalNonSecretCards uint16
	totalSecretCards    uint16
}

func (c *CardSetDetails) Set() *ref.CardSet {
	return &c.set
}

func (c *CardSetDetails) Cards() iter.Seq[Card] {
	return slices.Values(c.cards)
}

func (c *CardSetDetails) Boosters() iter.Seq[Booster] {
	return slices.Values(c.boosters)
}

func (c *CardSetDetails) TotalNonSecretCards() uint16 {
	return c.totalNonSecretCards
}

func (c *CardSetDetails) TotalSecretCards() uint16 {
	return c.totalSecretCards
}

func (c *CardSetDetails) TotalCards() uint16 {
	return uint16(len(c.cards))
}

func NewCardSetDetails(set ref.CardSet, boosters []Booster) CardSetDetails {
	var cards []Card
	for _, b := range boosters {
		for _, o := range b.offerings {
			// TODO: More efficient way than this. e.g. card number
			// TODO: Validate that cards with same number are the same
			card := o.Card()
			if !slices.Contains(cards, *card) {
				cards = append(cards, *card)
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

	return CardSetDetails{
		set:                 set,
		boosters:            boosters,
		cards:               cards,
		totalSecretCards:    totalSecretCards,
		totalNonSecretCards: uint16(len(cards)) - totalSecretCards,
	}
}

func NewBoosterOffering(
	card Card,
	first3CardProbability float64,
	fourthCardProbability float64,
	fifthCardProbability float64,
) BoosterOffering {
	return BoosterOffering{
		card:               card,
		first3CardOffering: first3CardProbability,
		fourthCardOffering: fourthCardProbability,
		fifthCardOffering:  fifthCardProbability,
		packProbability:    first3CardProbability*3 + fourthCardProbability + fifthCardProbability,
	}
}
