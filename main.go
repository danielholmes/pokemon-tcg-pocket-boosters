package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"ptcgpocket/collection"
	"ptcgpocket/data"
	"ptcgpocket/source"

	"encoding/json"

	"golang.org/x/sync/errgroup"
)

var cardSetDataSources = [...]*source.CardSetSerebiiSource{
	source.NewCardSetSerebiiSource(
		"genetic-apex",
		"Genetic Apex",
		[]*source.BoosterSerebiiSource{
			source.NewBoosterSerebiiSource(
				"Pikachu",
				"https://www.serebii.net/tcgpocket/geneticapex/pikachu.shtml",
				data.OfferingRatesTable{
					&data.RarityOneDiamond:   *data.NewBoosterOffering(2.000, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 2.571, 1.714, 0),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0, 0.357, 1.428, 0),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0, 0.333, 1.332, 0),
					&data.RarityOneStar:      *data.NewBoosterOffering(0, 0.321, 1.286, 5.0),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.050, 0.200, 5.0),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222, 0.888, 5.0),
					&data.RarityCrown:        *data.NewBoosterOffering(0, 0.013, 0.053, 5.0),
				},
				285,
			),
			source.NewBoosterSerebiiSource(
				"MewTwo",
				"https://www.serebii.net/tcgpocket/geneticapex/mewtwo.shtml",
				data.OfferingRatesTable{
					&data.RarityOneDiamond:   *data.NewBoosterOffering(2.000, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 2.571, 1.714, 0),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0, 0.357, 1.428, 0),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0, 0.333, 1.332, 0),
					&data.RarityOneStar:      *data.NewBoosterOffering(0, 0.321, 1.286, 5.263),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.050, 0.200, 5.263),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222, 0.888, 5.263),
					&data.RarityCrown:        *data.NewBoosterOffering(0, 0.013, 0.053, 5.263),
				},
				286,
			),
			source.NewBoosterSerebiiSource(
				"Charizard",
				"https://www.serebii.net/tcgpocket/geneticapex/charizard.shtml",
				data.OfferingRatesTable{
					&data.RarityOneDiamond:   *data.NewBoosterOffering(2.000, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 2.571, 1.714, 0),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0, 0.357, 1.428, 0),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0, 0.333, 1.332, 0),
					&data.RarityOneStar:      *data.NewBoosterOffering(0, 0.321, 1.286, 5.0),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.050, 0.200, 5.0),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222, 0.888, 5.0),
					&data.RarityCrown:        *data.NewBoosterOffering(0, 0.013, 0.053, 5.0),
				},
				284,
			),
		},
	),
	source.NewCardSetSerebiiSource(
		"mythical-island",
		"Mythical Island",
		[]*source.BoosterSerebiiSource{
			source.NewBoosterSerebiiSource(
				"Mew",
				"https://www.serebii.net/tcgpocket/mythicalisland/mew.shtml",
				data.OfferingRatesTable{
					&data.RarityOneDiamond:   *data.NewBoosterOffering(3.125, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 3.913, 2.608, 0),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0, 0.625, 2.500, 0),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0, 0.333, 1.332, 0),
					&data.RarityOneStar:      *data.NewBoosterOffering(0, 0.428, 1.714, 5.555),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.050, 0.200, 5.555),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222, 0.888, 5.555),
					&data.RarityCrown:        *data.NewBoosterOffering(0, 0.040, 0.160, 5.555),
				},
				86,
			),
		},
	),
	source.NewCardSetSerebiiSource(
		"space-time-smackdown",
		"Space-time Smackdown",
		[]*source.BoosterSerebiiSource{
			source.NewBoosterSerebiiSource(
				"Dialga",
				"https://www.serebii.net/tcgpocket/space-timesmackdown/dialga.shtml",
				data.OfferingRatesTable{
					&data.RarityOneDiamond:   *data.NewBoosterOffering(2.173, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 2.647, 1.764, 0),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0, 0.357, 1.428, 0),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0, 0.333, 1.332, 0),
					&data.RarityOneStar:      *data.NewBoosterOffering(0, 0.214, 0.857, 3.846),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.041, 0.166, 3.846),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222, 0.888, 3.846),
					&data.RarityCrown:        *data.NewBoosterOffering(0, 0.020, 0.080, 3.846),
				},
				207,
			),
			source.NewBoosterSerebiiSource(
				"Palkia",
				"https://www.serebii.net/tcgpocket/space-timesmackdown/palkia.shtml",
				data.OfferingRatesTable{
					&data.RarityOneDiamond:   *data.NewBoosterOffering(2.272, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 2.500, 1.666, 0),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0, 0.357, 1.428, 0),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0, 0.333, 1.332, 0),
					&data.RarityOneStar:      *data.NewBoosterOffering(0, 0.214, 0.857, 3.846),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.041, 0.166, 3.846),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222, 0.888, 3.846),
					&data.RarityCrown:        *data.NewBoosterOffering(0, 0.020, 0.080, 3.846),
				},
				206,
			),
		},
	),
}

func readUserCollection() (*collection.UserCollection, error) {
	dir, dErr := os.Getwd()
	if dErr != nil {
		return nil, dErr
	}

	collectionFilepath := filepath.Join(dir, "collection.json")

	raw, err := os.ReadFile(collectionFilepath)
	if err != nil {
		return nil, err
	}

	var allMissing map[data.CardSetId]([]data.CardSetNumber)
	uErr := json.Unmarshal(raw, &allMissing)
	if uErr != nil {
		return nil, uErr
	}

	userCollection := collection.NewUserCollection(allMissing)
	return &userCollection, nil
}

func main() {
	// Loading collection
	userCollection, uErr := readUserCollection()
	if uErr != nil {
		panic(uErr)
	}

	// Gather data from sources
	results := make(chan data.CardSet, len(cardSetDataSources))
	g, ctx := errgroup.WithContext(context.Background())
	for _, s := range cardSetDataSources {
		g.Go(func() error {
			return source.FetchCardSetDetails(ctx, s, results)
		})
	}
	err := g.Wait()
	close(results)
	if err != nil {
		panic(err)
	}

	var sets []*data.CardSet
	for o := range results {
		sets = append(sets, &o)
	}

	// Show check of booster probabilities
	fmt.Println("# Booster gathered data audit")
	for _, s := range sets {
		for b := range s.Boosters() {
			totalRegularPackOffering := 0.0
			totalRarePackOffering := 0.0
			totalFirstToThirdOffering := 0.0
			totalFourthOffering := 0.0
			totalFifthOffering := 0.0
			totalRareCardOffering := 0.0
			for c := range b.Offerings() {
				totalRegularPackOffering += c.RegularPackOffering()
				totalRarePackOffering += c.RegularPackOffering()
				totalFirstToThirdOffering += c.First3CardOffering()
				totalFourthOffering += c.FourthCardOffering()
				totalFifthOffering += c.FifthCardOffering()
				totalRareCardOffering += c.RareCardOffering()
			}
			fmt.Printf(
				" ## %v - %v\n   1-3: %.2f / 100%%\n   4: %.2f / 100%%\n   5: %.2f / 100%%\n   total regular: %.2f / 500%%\n   rare: %.2f / 100%%\n   total rare: %.2f / 500%%\n   \n",
				s.Name(),
				b.Name(),
				totalFirstToThirdOffering,
				totalFourthOffering,
				totalFifthOffering,
				totalRegularPackOffering,
				totalRareCardOffering,
				totalRarePackOffering,
			)
		}
	}

	// Show collection
	fmt.Println()
	fmt.Println("# Current collection")
	for _, s := range sets {
		missing, sExists := userCollection.MissingForSet(s.Id())
		if !sExists {
			fmt.Printf("Set id %v not found\n", s.Id())
			return
		}

		fmt.Printf(" ## %v\n", s.Name())

		totalSecretCardsCollected := 0
		totalNonSecretCardsCollected := 0
		for c := range s.Cards() {
			if !slices.Contains(missing, c.Number()) {
				if c.Rarity().IsSecret() {
					totalSecretCardsCollected += 1
				} else {
					totalNonSecretCardsCollected += 1
				}
			}
		}
		totalCollectedIncludingSecrets := totalSecretCardsCollected + totalNonSecretCardsCollected
		fmt.Printf(
			"    %v / %v (%v%%) %v★ Inc. secret %v / %v (%v%%)\n",
			totalNonSecretCardsCollected,
			s.TotalNonSecretCards(),
			100*totalNonSecretCardsCollected/int(s.TotalNonSecretCards()),
			totalSecretCardsCollected,
			totalCollectedIncludingSecrets,
			s.TotalCards(),
			100*(totalSecretCardsCollected+totalNonSecretCardsCollected)/int(s.TotalCards()),
		)
	}

	// Show booster values
	fmt.Println()
	fmt.Println("# Booster values")
	for _, s := range sets {
		missing, sExists := userCollection.MissingForSet(s.Id())
		if !sExists {
			fmt.Printf("Set id %v not found\n", s.Id())
			return
		}

		fmt.Printf(" ## %v\n", s.Name())

		for b := range s.Boosters() {
			fmt.Printf("  ### %v\n", b.Name())

			totalOfferingMissing := 0.0
			for o := range b.Offerings() {
				if slices.Contains(missing, o.Card().Number()) {
					totalOfferingMissing += o.RegularPackOffering()*0.9995 + o.RarePackOffering()*0.0005
				}
			}

			fmt.Printf("   Total chance of receiving a missing %.2f%%\n", totalOfferingMissing)
		}
	}
}
