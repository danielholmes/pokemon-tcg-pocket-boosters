package collection

import (
	"ptcgpocket/data"
	"testing"
)

func TestCreateNewCollectionWithAddedCards(t *testing.T) {
	collection := NewUserCollection(
		map[data.ExpansionId]*ExpansionCollection{
			"genetic-apex":    &ExpansionCollection{missingCardNumbers: []data.ExpansionCardNumber{1, 2, 3}},
			"mythical-island": &ExpansionCollection{missingCardNumbers: []data.ExpansionCardNumber{1, 2, 3}},
		},
	)

	collection.expansions["genetic-apex"].AddCardsFromBooster(
		[5]data.ExpansionCardNumber{1, 3, 99, 100, 101},
	)

	newMissingForGenetic, _ := collection.MissingForExpansion("genetic-apex")
	if len(newMissingForGenetic) != 1 {
		t.Errorf("New missing genetic apex incorrect length = %d; want 1", len(newMissingForGenetic))
	}
	if newMissingForGenetic[0] != 2 {
		t.Errorf("New missing genetic apex incorrect contents = %d; want 2", newMissingForGenetic[0])
	}
	newMissingForMythical, _ := collection.MissingForExpansion("mythical-island")
	if len(newMissingForMythical) != 3 {
		t.Errorf("New missing mythical island incorrect length = %d; want 3", len(newMissingForMythical))
	}
}
