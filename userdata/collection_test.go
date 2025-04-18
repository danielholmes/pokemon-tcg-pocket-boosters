package userdata

import (
	"ptcgpocket/data"
	"slices"
	"testing"
)

func TestCreateNewCollectionWithAddedCards(t *testing.T) {
	ga1 := data.NewCard(
		data.NewBaseCard("Test 1", 100),
		1,
		data.RarityOneDiamond,
	)
	ga2 := data.NewCard(
		data.NewBaseCard("Test 2", 100),
		2,
		data.RarityOneDiamond,
	)
	ga3 := data.NewCard(
		data.NewBaseCard("Test 3", 100),
		3,
		data.RarityOneDiamond,
	)
	ga99 := data.NewCard(
		data.NewBaseCard("Test 99", 100),
		99,
		data.RarityOneDiamond,
	)
	collection := NewUserCollection(
		map[data.ExpansionId]*ExpansionCollection{
			"genetic-apex": &ExpansionCollection{
				missingCards: []*data.Card{
					ga1,
					ga2,
					ga3,
				},
			},
			"mythical-island": &ExpansionCollection{
				missingCards: []*data.Card{
					data.NewCard(
						data.NewBaseCard("Test MI 1", 100),
						1,
						data.RarityOneDiamond,
					),
					data.NewCard(
						data.NewBaseCard("Test MI 2", 100),
						2,
						data.RarityOneDiamond,
					),
					data.NewCard(
						data.NewBaseCard("Test MI 3", 100),
						3,
						data.RarityOneDiamond,
					),
				}},
		},
	)

	collection.expansions["genetic-apex"].AcquireCardsFromBooster(
		slices.Values([]*data.Card{ga1, ga3, ga99, ga99, ga1}[:]),
	)

	newMissingForGenetic, _ := collection.MissingForExpansion("genetic-apex")
	if len(newMissingForGenetic) != 1 {
		t.Errorf("New missing genetic apex incorrect length = %d; want 1", len(newMissingForGenetic))
	}
	if newMissingForGenetic[0] != ga2 {
		t.Errorf("New missing genetic apex incorrect contents = %v; want 2", newMissingForGenetic[0])
	}
	newMissingForMythical, _ := collection.MissingForExpansion("mythical-island")
	if len(newMissingForMythical) != 3 {
		t.Errorf("New missing mythical island incorrect length = %d; want 3", len(newMissingForMythical))
	}
}
