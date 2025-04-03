package data

import "testing"

func TestNewBoosterOfferings(t *testing.T) {
	booster := NewBooster(
		"Test booster",
		[]*Card{
			{
				core: &BaseCard{
					name:   "Pikachu",
					health: 60,
				},
				number: 12,
				rarity: &RarityOneDiamond,
			},
		},
		OfferingRatesTable{
			&RarityOneDiamond: {
				first3CardOffering: 0.5,
				fourthCardOffering: 0.4,
				fifthCardOffering:  0.6,
				rareOffering:       0.0,
			},
		},
		1,
	)

	offeringsSeq := booster.Offerings()
	offerings := make([]*BoosterCardOffering, 0)
	for o := range offeringsSeq {
		offerings = append(offerings, o)
	}
	if len(offerings) != 1 {
		t.Errorf("Booster Offerings incorrect length = %d; want 1", len(offerings))
	}
	offering1 := offerings[0]
	if offering1.first3CardOffering != 0.5 {
		t.Errorf("Booster Offering 1 incorrect first3Card = %v; want 0.5", offering1.first3CardOffering)
	}
	if offering1.fourthCardOffering != 0.4 {
		t.Errorf("Booster Offering 1 incorrect fourthCard = %v; want 0.4", offering1.fourthCardOffering)
	}
	if offering1.fifthCardOffering != 0.6 {
		t.Errorf("Booster Offering 1 incorrect fifthCard = %v; want 0.4", offering1.fifthCardOffering)
	}
	if offering1.rareCardOffering != 0.0 {
		t.Errorf("Booster Offering 1 incorrect rareCard = %v; want 0.0", offering1.rareCardOffering)
	}
	if offering1.RegularPackOffering() != 2.5 {
		t.Errorf("Booster Offering 1 incorrect regular pack offering = %v; want 2.5", offering1.RegularPackOffering())
	}
}
