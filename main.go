package main

import (
	"fmt"
	"slices"
	"sync"

	"ptcgpocket/collection"
	"ptcgpocket/data"
	"ptcgpocket/source"
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

var userCollection collection.UserCollection = collection.NewUserCollection(
	map[data.CardSetId]([]data.CardSetNumber){
		"genetic-apex": {
			3, 4, 7, 10, 13, 22, 36, 39, 41, 47, 50, 56, 61, 69, 73, 76, 80, 84, 89, 95, 98, 101, 107, 117,
			123, 124, 145, 146, 148, 149, 159, 163, 166, 175, 177, 178, 185, 191, 195, 197, 203, 204, 205, 221,
			225, 226,
			228, 229, 230, 231, 232, 233, 236, 237, 238, 240, 241, 242, 243, 244, 246, 248, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286,
		},
		"mythical-island": {
			2, 3, 6, 7, 18, 25, 26, 32, 44, 46, 59, 60, 62,
			71, 73, 75, 76, 79, 80, 81, 82, 83, 84, 85, 86,
		},
		"space-time-smackdown": {
			5, 7, 18, 20, 22, 24, 29, 32, 33, 34, 36, 37, 41, 60, 65, 76, 79, 89, 90, 92, 94, 104, 109, 113,
			117, 120, 123, 129, 147, 153,
			156, 157, 158, 159, 160, 161, 162, 164, 166, 167, 168, 169, 170, 171, 172, 173, 176, 177, 178, 179,
			180, 181, 182, 183, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 200, 201, 202,
			203, 205, 206, 207,
		},
	},
)

func main() {
	var wg sync.WaitGroup

	results := make(chan data.CardSet, len(cardSetDataSources))
	for _, s := range cardSetDataSources {
		wg.Add(1)
		go source.FetchCardSetDetails(s, &wg, results)
	}

	wg.Wait()
	close(results)

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
