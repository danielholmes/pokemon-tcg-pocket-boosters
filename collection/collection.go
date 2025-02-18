package collection

import (
	"ptcgpocket/data"
)

type UserCollection struct {
	missingCardNumbers map[data.CardSetId]([]data.CardSetNumber)
}

func NewUserCollection(missingCardNumbers map[data.CardSetId]([]data.CardSetNumber)) UserCollection {
	return UserCollection{missingCardNumbers: missingCardNumbers}
}

func (c *UserCollection) MissingForSet(setId data.CardSetId) ([]data.CardSetNumber, bool) {
	v, e := c.missingCardNumbers[setId]
	return v, e
}
