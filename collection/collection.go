package collection

import (
	"iter"
	"maps"
	"ptcgpocket/ref"
)

type UserCollection struct {
	missingCardNumbers map[ref.CardSet]([]ref.CardSetNumber)
}

func NewUserCollection(missingCardNumbers map[ref.CardSet]([]ref.CardSetNumber)) UserCollection {
	return UserCollection{missingCardNumbers: missingCardNumbers}
}

func (c *UserCollection) MissingCardNumbers() iter.Seq2[ref.CardSet, []ref.CardSetNumber] {
	return maps.All(c.missingCardNumbers)
}
