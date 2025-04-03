package collection

import (
	"fmt"
	"ptcgpocket/data"
	"slices"
)

type UserCollection struct {
	missingCardNumbers map[data.ExpansionId]([]data.ExpansionNumber)
}

func NewUserCollection(missingCardNumbers map[data.ExpansionId]([]data.ExpansionNumber)) UserCollection {
	for eId, m := range missingCardNumbers {
		if slices.Contains(m, 0) {
			panic(fmt.Sprintf("Invalid expansion number %v %v", eId, m))
		}
	}
	return UserCollection{missingCardNumbers: missingCardNumbers}
}

func (c *UserCollection) Clone() *UserCollection {
	newMissingCardNumbers := make(map[data.ExpansionId]([]data.ExpansionNumber))
	for eId, m := range c.missingCardNumbers {
		newMissingCardNumbers[eId] = slices.Clone(m)
	}

	newUserCollection := NewUserCollection(newMissingCardNumbers)
	return &newUserCollection
}

func (c *UserCollection) FirstIncompleteExpansionId() (data.ExpansionId, error) {
	for i, m := range c.missingCardNumbers {
		if len(m) > 0 {
			return i, nil
		}
	}
	return "", fmt.Errorf("no incomplete expansion")
}

func (c *UserCollection) MissingForExpansion(expansionId data.ExpansionId) ([]data.ExpansionNumber, bool) {
	v, e := c.missingCardNumbers[expansionId]
	return v, e
}

func (c *UserCollection) AddCards(
	expansionId data.ExpansionId,
	addedNumbers [5]data.ExpansionNumber,
) {
	c.missingCardNumbers[expansionId] = slices.DeleteFunc(c.missingCardNumbers[expansionId], func(n data.ExpansionNumber) bool {
		for _, a := range addedNumbers {
			if a == n {
				return true
			}
		}
		return false
	})
}
