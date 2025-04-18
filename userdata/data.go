package userdata

import (
	"iter"
	"slices"
)

type UserData struct {
	collection *UserCollection
	wishlists  iter.Seq[*Wishlist]
}

func NewUserData(collection *UserCollection, wishlists []*Wishlist) *UserData {
	return &UserData{collection: collection, wishlists: slices.Values(wishlists)}
}

func (u *UserData) Collection() *UserCollection {
	return u.collection
}

func (u *UserData) Wishlists() iter.Seq[*Wishlist] {
	return u.wishlists
}
