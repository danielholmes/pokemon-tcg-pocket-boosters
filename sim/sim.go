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

func RunSim(
	expansions []*data.Expansion,
	userCollection *collection.UserCollection,
) (*SimRun, error) {
	simCollection := userCollection.Clone()
	numberOfPacksOpened := make(map[*data.Expansion]uint64)
	for _, e := range expansions {
		moreMissing := true
		for moreMissing {
			missing, missingFound := simCollection.MissingForExpansion(e.Id())
			if !missingFound {
				panic("No missing found")
			}

			if len(missing) == 0 {
				moreMissing = false
				continue
			}

			simBooster, sErr := e.GetBoosterOfferingCardNumber(
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
