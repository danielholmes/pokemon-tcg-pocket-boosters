package sim

import (
	"fmt"
	"iter"
	"maps"
	"ptcgpocket/collection"
	"ptcgpocket/data"
)

type ExpansionSimRun struct {
	numOpened                      uint64
	numPackPoints                  uint64
	totalPackPoints                uint64
	numCardsObtainedFromPackPoints uint64
	numRarePacks                   uint64
}

func NewExpansionSimRun(
	numOpened uint64,
	numPackPoints uint64,
	totalPackPoints uint64,
	numCardsObtainedFromPackPoints uint64,
	numRarePacks uint64,
) *ExpansionSimRun {
	return &ExpansionSimRun{
		numOpened:                      numOpened,
		numPackPoints:                  numPackPoints,
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
) (*SimRun, error) {
	simCollection := userCollection.Clone()
	expansionRuns := make(map[*data.Expansion]*ExpansionSimRun)
	for _, e := range expansions {
		isExpansionComplete := false
		for !isExpansionComplete {
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
			if highestPackPointsCard != nil && eSimRun.numPackPoints >= uint64(highestPackPointsCard.Rarity().PackPointsToObtain()) {
				eSimRun.numPackPoints -= uint64(highestPackPointsCard.Rarity().PackPointsToObtain())
				eSimRun.numCardsObtainedFromPackPoints += 1
				simCollection.AddCard(
					e.Id(),
					highestPackPointsCard.Number(),
				)
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

			boosterInstance := simBooster.CreateRandomInstance()
			simCollection.AddCards(
				e.Id(),
				boosterInstance.CardNumbers(),
			)

			eSimRun.numOpened++
			eSimRun.totalPackPoints += 5
			eSimRun.numPackPoints += 5
			if boosterInstance.IsRare() {
				eSimRun.numRarePacks++
			}
		}
	}

	return &SimRun{expansionRuns: expansionRuns}, nil
}
