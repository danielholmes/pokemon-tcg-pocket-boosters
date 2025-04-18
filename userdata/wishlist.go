package userdata

import "ptcgpocket/data"

type ExpansionWishlist struct {
	cards []*data.Card
}

type Wishlist struct {
	name       string
	expansions map[data.ExpansionId]*ExpansionWishlist
}

func (w *Wishlist) Name() string {
	return w.name
}

func (w *Wishlist) CardsForExpansion(expansionId data.ExpansionId) ([]*data.Card, bool) {
	eW, eWFound := w.expansions[expansionId]
	if !eWFound {
		return nil, false
	}
	return eW.cards, true
}
