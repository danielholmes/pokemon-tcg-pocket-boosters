package data

import (
	"strings"
)

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

func (r *Rarity) IsSecret() bool {
	return r.isSecret
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
