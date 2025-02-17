package data

import (
	"ptcgpocket/ref"
	"slices"
	"sort"
)

type Rarity struct {
	order uint8
	isSecret bool
	value string
}

var (
	RarityOneDiamond = Rarity{0, false, "♢"}
	RarityTwoDiamond = Rarity{1, false,"♢♢"}
	RarityThreeDiamond = Rarity{2, false,"♢♢♢"}
	RarityFourDiamond = Rarity{3, false,"♢♢♢♢"}
	RarityOneStar = Rarity{4, true,"☆"}
	RarityTwoStar = Rarity{5, true,"☆☆"}
	RarityThreeStar = Rarity{6, true,"☆☆☆"}
	RarityCrown = Rarity{7, true,"♕"}
)

func (r *Rarity) IsSecret() bool {
	return r.isSecret
}

type Card struct {
	Rarity *Rarity
	Name   string
	Number ref.CardSetNumber
}

type BoosterOffering struct {
	Card     Card
	first3CardProbability float64
	fourthCardProbability float64
	fifthCardProbability float64
	PackProbability float64
}

type Booster struct {
	Name       string
	Offerings []BoosterOffering
}

type CardSetDetails struct {
	Set ref.CardSet
	boosters []Booster
	cards []Card
	totalNonSecretCards uint16
	totalSecretCards uint16
}

func (c *CardSetDetails) Cards() []Card {
	return c.cards
}

func (c *CardSetDetails) Boosters() []Booster {
	return c.boosters
}

func (c *CardSetDetails) TotalNonSecretCards() uint16 {
	return c.totalNonSecretCards
}

func (c *CardSetDetails) TotalSecretCards() uint16 {
	return c.totalSecretCards
}

func NewCardSetDetails(set ref.CardSet, boosters []Booster) CardSetDetails {
	var cards []Card
	for _, b := range boosters {
		for _, o := range b.Offerings {	
			// TODO: More efficient way than this. e.g. card number
			// TODO: Validate that cards with same number are the same
			if !slices.Contains(cards, o.Card) {		
				cards = append(cards, o.Card)
			}
		}
	}
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Number < cards[j].Number
	})

	var totalSecretCards uint16 = 0
	for _, c := range cards {
		if c.Rarity.isSecret {
			totalSecretCards += 1
		}
	}

	return CardSetDetails{
		Set: set,
		boosters: boosters,
		cards: cards,
		totalSecretCards: totalSecretCards,
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
		Card: card,
		first3CardProbability: first3CardProbability,
		fourthCardProbability: fourthCardProbability,
		fifthCardProbability: fifthCardProbability,
		PackProbability: first3CardProbability*3 + fourthCardProbability + fifthCardProbability,
	}
}