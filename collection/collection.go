package collection

import (
	"ptcgpocket/data"
)

type UserCollection struct {
	missingCardNumbers map[data.ExpansionId]([]data.ExpansionNumber)
}

func NewUserCollection(missingCardNumbers map[data.ExpansionId]([]data.ExpansionNumber)) UserCollection {
	return UserCollection{missingCardNumbers: missingCardNumbers}
}

func (c *UserCollection) MissingForSet(setId data.ExpansionId) ([]data.ExpansionNumber, bool) {
	v, e := c.missingCardNumbers[setId]
	return v, e
}
