package data

import (
	"strings"
)

type Rarity struct {
	order              uint8
	isSecret           bool
	value              string
	packPointsToObtain uint16
}

const RarityDiamondChar = '♢'
const RarityStarChar = '☆'
const RarityShinyChar = '✵'
const RarityCrownChar = '♕'

func newRarity(order uint8, isSecret bool, symbol rune, symbolCount uint8, packPointsToObtain uint16) *Rarity {
	value := ""
	for range symbolCount {
		value += string(symbol)
	}
	return &Rarity{
		order:              order,
		isSecret:           isSecret,
		value:              value,
		packPointsToObtain: packPointsToObtain,
	}
}

func (r *Rarity) PackPointsToObtain() uint16 {
	return r.packPointsToObtain
}

func (r *Rarity) IsStar() bool {
	return strings.ContainsRune(r.value, RarityStarChar)
}

func (r *Rarity) IsCrown() bool {
	return strings.ContainsRune(r.value, RarityCrownChar)
}

func (r *Rarity) IsShiny() bool {
	return strings.ContainsRune(r.value, RarityShinyChar)
}

func (r *Rarity) IsSecret() bool {
	return r.isSecret
}

func (r *Rarity) String() string {
	return r.value
}

var (
	RarityOneDiamond   = newRarity(0, false, RarityDiamondChar, 1, 35)
	RarityTwoDiamond   = newRarity(1, false, RarityDiamondChar, 2, 70)
	RarityThreeDiamond = newRarity(2, false, RarityDiamondChar, 3, 150)
	RarityFourDiamond  = newRarity(3, false, RarityDiamondChar, 4, 500)
	RarityOneStar      = newRarity(4, true, RarityStarChar, 1, 400)
	RarityTwoStar      = newRarity(5, true, RarityStarChar, 2, 1_250)
	RarityThreeStar    = newRarity(6, true, RarityStarChar, 3, 1_500)
	RarityOneShiny     = newRarity(7, true, RarityShinyChar, 1, 1_000)
	RarityTwoShiny     = newRarity(8, true, RarityShinyChar, 2, 1_350)
	RarityCrown        = newRarity(9, true, RarityCrownChar, 1, 2_500)
)

var OrderedRarities = []*Rarity{
	RarityOneDiamond,
	RarityTwoDiamond,
	RarityThreeDiamond,
	RarityFourDiamond,
	RarityOneStar,
	RarityTwoStar,
	RarityThreeStar,
	RarityOneShiny,
	RarityTwoShiny,
	RarityCrown,
}

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

type ExpansionCardNumber uint16

type Card struct {
	core   *BaseCard
	number ExpansionCardNumber
	rarity *Rarity
}

func NewCard(
	core *BaseCard,
	number ExpansionCardNumber,
	rarity *Rarity,
) *Card {
	return &Card{core: core, number: number, rarity: rarity}
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

func (c *Card) Number() ExpansionCardNumber {
	return c.number
}
