package sim

import (
	"fmt"
	"iter"
	"maps"
	"math/rand/v2"
	"ptcgpocket/collection"
	"ptcgpocket/data"
)

type ExpansionSimRun struct {
	numOpened                      uint64
	totalPackPoints                uint64
	numCardsObtainedFromPackPoints uint64
	numRarePacks                   uint64
}

func NewExpansionSimRun(
	numOpened uint64,
	totalPackPoints uint64,
	numCardsObtainedFromPackPoints uint64,
	numRarePacks uint64,
) *ExpansionSimRun {
	return &ExpansionSimRun{
		numOpened:                      numOpened,
		totalPackPoints:                totalPackPoints,
		numCardsObtainedFromPackPoints: numCardsObtainedFromPackPoints,
		numRarePacks:                   numRarePacks,
	}
}

func (r *ExpansionSimRun) NumOpened() uint64 {
	return r.numOpened
}

func (r *ExpansionSimRun) TotalPackPoints() uint64 {
	return r.totalPackPoints
}

func (r *ExpansionSimRun) NumCardsObtainedFromPackPoints() uint64 {
	return r.numCardsObtainedFromPackPoints
}

func (r *ExpansionSimRun) NumRarePacks() uint64 {
	return r.numRarePacks
}

type SimRun struct {
	expansionRuns map[*data.Expansion]*ExpansionSimRun
}

func (r *SimRun) TotalPacksOpened() uint64 {
	var total uint64
	for _, n := range r.expansionRuns {
		total += n.numOpened
	}
	return total
}

func (r *SimRun) ExpansionRuns() iter.Seq2[*data.Expansion, *ExpansionSimRun] {
	return maps.All(r.expansionRuns)
}

type ExpansionSimCompletePredicate func(*data.Expansion, []data.ExpansionNumber) bool

func RunSim(
	expansions []*data.Expansion,
	userCollection *collection.UserCollection,
	expansionCompletePredicate ExpansionSimCompletePredicate,
	randomGenerator *rand.Rand,
) (*SimRun, error) {
	simCollection := userCollection.Clone()
	expansionRuns := make(map[*data.Expansion]*ExpansionSimRun)
	for _, e := range expansions {
		isExpansionComplete := false
		for !isExpansionComplete {
			eCollection := simCollection.GetExpansionCollection(e.Id())
			missing, missingFound := simCollection.MissingForExpansion(e.Id())
			if !missingFound {
				panic("No missing found")
			}

			if expansionCompletePredicate(e, missing) {
				isExpansionComplete = true
				continue
			}

			eSimRun := expansionRuns[e]
			if eSimRun == nil {
				eSimRun = &ExpansionSimRun{}
				expansionRuns[e] = eSimRun
			}

			// Decide, should we trade in pack points or pick a booster?
			// TODO: This can be more efficient by ending search early.
			var highestPackPointsCard *data.Card
			for _, missingNumber := range missing {
				card, cErr := e.GetCardByNumber(missingNumber)
				if cErr != nil {
					panic(cErr)
				}
				if highestPackPointsCard == nil || card.Rarity().PackPointsToObtain() > highestPackPointsCard.Rarity().PackPointsToObtain() {
					highestPackPointsCard = card
				}
			}
			if highestPackPointsCard != nil && eCollection.PackPoints() >= highestPackPointsCard.Rarity().PackPointsToObtain() {
				eCollection.AcquireCardUsingPackPoints(highestPackPointsCard)
				eSimRun.numCardsObtainedFromPackPoints += 1
				continue
			}

			// Not enough pack points, now we choose a booster instead.
			simBooster, sErr := e.GetHighestOfferingBoosterForMissingCards(
				missing,
			)
			if sErr != nil {
				fmt.Printf("No missing %v %v\n", e.Id(), missing)
				panic("should be able to find booster for missing number")
			}

			boosterInstance := simBooster.CreateRandomInstance(randomGenerator)
			// fmt.Printf("C %v \n", boosterInstance.CardNumbers())
			eCollection.AddCardsFromBooster(boosterInstance.CardNumbers())

			eSimRun.numOpened++
			eSimRun.totalPackPoints += 5
			if boosterInstance.IsRare() {
				eSimRun.numRarePacks++
			}
		}
	}

	return &SimRun{expansionRuns: expansionRuns}, nil
}
