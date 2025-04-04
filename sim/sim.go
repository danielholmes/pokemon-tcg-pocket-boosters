package sim

import (
	"fmt"
	"iter"
	"maps"
	"ptcgpocket/collection"
	"ptcgpocket/data"
)

type ExpansionSimRun struct {
	numOpened    uint64
	numRarePacks uint64
}

func (r *ExpansionSimRun) NumOpened() uint64 {
	return r.numOpened
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

			if expansionRuns[e] == nil {
				expansionRuns[e] = &ExpansionSimRun{}
			}
			expansionRuns[e].numOpened++
			if boosterInstance.IsRare() {
				expansionRuns[e].numRarePacks++
			}
		}
	}

	return &SimRun{expansionRuns: expansionRuns}, nil
}
