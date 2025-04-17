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

const diamondChar = "♢"
const starChar = "☆"
const shinyChar = "✵"
const crownChar = "♕"

var (
	RarityOneDiamond   = Rarity{0, false, strings.Repeat(diamondChar, 1), 35}
	RarityTwoDiamond   = Rarity{1, false, strings.Repeat(diamondChar, 2), 70}
	RarityThreeDiamond = Rarity{2, false, strings.Repeat(diamondChar, 3), 150}
	RarityFourDiamond  = Rarity{3, false, strings.Repeat(diamondChar, 4), 500}
	RarityOneStar      = Rarity{4, true, strings.Repeat(starChar, 1), 400}
	RarityTwoStar      = Rarity{5, true, strings.Repeat(starChar, 2), 1_250}
	RarityThreeStar    = Rarity{6, true, strings.Repeat(starChar, 3), 1_500}
	RarityOneShiny     = Rarity{7, true, strings.Repeat(shinyChar, 1), 1_000}
	RarityTwoShiny     = Rarity{8, true, strings.Repeat(shinyChar, 2), 1_350}
	RarityCrown        = Rarity{9, true, strings.Repeat(crownChar, 1), 2_500}
)

var OrderedRarities = []*Rarity{
	&RarityOneDiamond,
	&RarityTwoDiamond,
	&RarityThreeDiamond,
	&RarityFourDiamond,
	&RarityOneStar,
	&RarityTwoStar,
	&RarityThreeStar,
	&RarityOneShiny,
	&RarityTwoShiny,
	&RarityCrown,
}

func (r *Rarity) PackPointsToObtain() uint16 {
	return r.packPointsToObtain
}

func (r *Rarity) IsStar() bool {
	return strings.Contains(r.value, starChar)
}

func (r *Rarity) IsCrown() bool {
	return strings.Contains(r.value, crownChar)
}

func (r *Rarity) IsShiny() bool {
	return strings.Contains(r.value, shinyChar)
}

func (r *Rarity) IsSecret() bool {
	return r.isSecret
}

func (r *Rarity) String() string {
	return r.value
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
