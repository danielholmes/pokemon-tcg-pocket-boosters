package sim

import (
	"fmt"
	"iter"
	"maps"
	"ptcgpocket/collection"
	"ptcgpocket/data"
)

type SimRun struct {
	numberOfPacksOpened map[*data.Expansion]uint64
}

func (r *SimRun) TotalPacksOpened() uint64 {
	var total uint64
	for _, n := range r.numberOfPacksOpened {
		total += n
	}
	return total
}

func (r *SimRun) NumberOfPacksOpened() iter.Seq2[*data.Expansion, uint64] {
	return maps.All(r.numberOfPacksOpened)
}

type ExpansionSimCompletePredicate func(*data.Expansion, []data.ExpansionNumber) bool

func RunSim(
	expansions []*data.Expansion,
	userCollection *collection.UserCollection,
	expansionCompletePredicate ExpansionSimCompletePredicate,
) (*SimRun, error) {
	simCollection := userCollection.Clone()
	numberOfPacksOpened := make(map[*data.Expansion]uint64)
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

			numberOfPacksOpened[e]++
		}
	}

	return &SimRun{numberOfPacksOpened: numberOfPacksOpened}, nil
}
