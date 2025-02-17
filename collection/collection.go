package collection

import (
	"ptcgpocket/ref"
)

type UserCollection struct {
	MissingCardNumbers map[ref.CardSet]([]ref.CardSetNumber)
}
