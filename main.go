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
					&data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/50.0, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 90.0/35.0, 60.0/35.0, 0),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0, 5.0/14.0, 20.0/14.0, 0),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0, 1.666/5.0, 6.664/5.0, 0),
					&data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					&data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					&data.RarityOneStar:      *data.NewBoosterOffering(0, 2.572/8.0, 10.288/8.0, 40.0/8.0),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.50/10.0, 0.200/10.0, 50.0/10.0),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222/1.0, 0.888/1.0, 5.0/1.0),
					&data.RarityCrown:        *data.NewBoosterOffering(0, 0.4/3.0, 0.16/3.0, 5.0/1.0),
				},
				285,
			),
			source.NewBoosterSerebiiSource(
				"MewTwo",
				"https://www.serebii.net/tcgpocket/geneticapex/mewtwo.shtml",
				data.OfferingRatesTable{
					&data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/50.0, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 90.0/35.0, 60.0/35.0, 0),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0, 5.0/14.0, 20.0/14.0, 0),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0, 1.666/5.0, 6.664/5.0, 0),
					&data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					&data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					&data.RarityOneStar:      *data.NewBoosterOffering(0, 2.572/8.0, 10.288/8.0, 42.105/8.0),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.50/10.0, 0.200/10.0, 47.368/9.0),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222/1.0, 0.888/1.0, 5.263/1.0),
					&data.RarityCrown:        *data.NewBoosterOffering(0, 0.4/3.0, 0.16/3.0, 5.263/1.0),
				},
				286,
			),
			source.NewBoosterSerebiiSource(
				"Charizard",
				"https://www.serebii.net/tcgpocket/geneticapex/charizard.shtml",
				data.OfferingRatesTable{
					&data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/50.0, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 90.0/35.0, 60.0/35.0, 0),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0, 5.0/14.0, 20.0/14.0, 0),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0, 1.666/5.0, 6.664/5.0, 0),
					&data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					&data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					&data.RarityOneStar:      *data.NewBoosterOffering(0, 2.572/8.0, 10.288/8.0, 40.0/8.0),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.50/10.0, 0.200/10.0, 50.0/10.0),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222/1.0, 0.888/1.0, 5.0/1.0),
					&data.RarityCrown:        *data.NewBoosterOffering(0, 0.4/3.0, 0.16/3.0, 5.0/1.0),
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
					&data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/32.0, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 90.0/23.0, 60.0/23.0, 0),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0, 5.0/8.0, 20.0/8.0, 0),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0, 1.666/5, 6.664/5.0, 0),
					&data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					&data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					&data.RarityOneStar:      *data.NewBoosterOffering(0, 2.572/6.0, 10.288/6.0, 33.333/6.0),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.500/10.0, 2.000/10.0, 55.555/10.0),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222/1.0, 0.888/1.0, 5.555/1.0),
					&data.RarityCrown:        *data.NewBoosterOffering(0, 0.040/1.0, 0.160/1.0, 5.555/1.0),
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
					&data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/46.0, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 90.0/34.0, 60.0/34.0, 0),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0, 5.0/14.0, 20.0/14.0, 0),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0, 1.666/5.0, 6.664/5.0, 0),
					&data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					&data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					&data.RarityOneStar:      *data.NewBoosterOffering(0, 2.572/12.0, 10.288/12.0, 46.153/12.0),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.5/12.0, 2.0/12.0, 46.153/12.0),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222/1.0, 0.888/1.0, 3.846/1.0),
					&data.RarityCrown:        *data.NewBoosterOffering(0, 0.040/2.0, 0.160/2.0, 3.846/1.0),
				},
				207,
			),
			source.NewBoosterSerebiiSource(
				"Palkia",
				"https://www.serebii.net/tcgpocket/space-timesmackdown/palkia.shtml",
				data.OfferingRatesTable{
					&data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/44.0, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 90.0/36.0, 60.0/36.0, 0),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0, 5.0/14.0, 20.0/14.0, 0),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0, 1.666/5.0, 6.664/5.0, 0),
					&data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					&data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					&data.RarityOneStar:      *data.NewBoosterOffering(0, 2.572/12.0, 10.288/12.0, 46.153/12.0),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.5/12.0, 2.0/12.0, 46.153/12.0),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222/1.0, 0.888/1.0, 3.846/1.0),
					&data.RarityCrown:        *data.NewBoosterOffering(0, 0.040/2.0, 0.160/2.0, 3.846/1.0),
				},
				206,
			),
		},
	),
	source.NewCardSetSerebiiSource(
		"triumphant-light",
		"Triumphant Light",
		[]*source.BoosterSerebiiSource{
			source.NewBoosterSerebiiSource(
				"Arceus",
				"https://www.serebii.net/tcgpocket/triumphantlight/arceus.shtml",
				data.OfferingRatesTable{
					&data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/31.0, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0.000, 90.0/26.0, 60.0/26.0, 0.000),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0.000, 5.0/13.0, 20.0/13.0, 0.000),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0.000, 1.666/5.0, 6.664/5.0, 0.000),
					&data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					&data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					&data.RarityOneStar:      *data.NewBoosterOffering(0.000, 2.572/6.0, 10.288/6.0, 28.571/6.0),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0.000, 0.5/13.0, 2.0/13.0, 61.904/13.0),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0.000, 0.222/1.0, 0.888/1.0, 4.761/1.0),
					&data.RarityCrown:        *data.NewBoosterOffering(0.000, 0.040/1.0, 0.160/1.0, 4.761/1.0),
				},
				96,
			),
		},
	),
	source.NewCardSetSerebiiSource(
		"shining-revelry",
		"Shining Revelry",
		[]*source.BoosterSerebiiSource{
			source.NewBoosterSerebiiSource(
				"Booster",
				"https://www.serebii.net/tcgpocket/shiningrevelry/booster.shtml",
				data.OfferingRatesTable{
					&data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/32.0, 0, 0, 0),
					&data.RarityTwoDiamond:   *data.NewBoosterOffering(0.000, 89.000/22, 56.000/22, 0.000),
					&data.RarityThreeDiamond: *data.NewBoosterOffering(0.000, 4.952/9, 19.810/9, 0.000),
					&data.RarityFourDiamond:  *data.NewBoosterOffering(0.000, 1.666/9, 6.664/9, 0.000),
					&data.RarityOneStar:      *data.NewBoosterOffering(0.000, 2.572/6, 10.288/6, 15.384/6),
					&data.RarityTwoStar:      *data.NewBoosterOffering(0.000, 0.500/17, 2.000/17, 43.589/17),
					&data.RarityThreeStar:    *data.NewBoosterOffering(0.000, 0.222/1, 0.888/1, 2.564/1),
					&data.RarityOneShiny:     *data.NewBoosterOffering(0.000, 0.714/10, 2.857/10, 25.641/10),
					&data.RarityTwoShiny:     *data.NewBoosterOffering(0.000, 0.333/4, 1.333/4, 10.256/4),
					&data.RarityCrown:        *data.NewBoosterOffering(0.000, 0.040, 0.160, 2.564),
				},
				111,
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

	// Audit of booster probabilities
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

	// Current collection stats
	fmt.Println()
	fmt.Println("# Current collection")
	for _, s := range sets {
		missing, sExists := userCollection.MissingForSet(s.Id())
		if !sExists {
			fmt.Printf("Set id %v not found\n", s.Id())
			return
		}

		fmt.Printf(" ## %v\n", s.Name())

		totalStarSecretCardsCollected := 0
		totalCrownSecretCardsCollected := 0
		totalNonSecretCardsCollected := 0
		totalShinySecretCardsCollected := 0
		for c := range s.Cards() {
			if !slices.Contains(missing, c.Number()) {
				if c.Rarity().IsStar() {
					totalStarSecretCardsCollected += 1
				} else if c.Rarity().IsCrown() {
					totalCrownSecretCardsCollected += 1
				} else if c.Rarity().IsShiny() {
					totalShinySecretCardsCollected += 1
				} else {
					totalNonSecretCardsCollected += 1
				}
			}
		}
		totalCollectedIncludingSecrets := totalStarSecretCardsCollected + totalShinySecretCardsCollected + totalCrownSecretCardsCollected + totalNonSecretCardsCollected
		fmt.Printf(
			"    %v / %v (%v%%) %v★ %v✵ %v♕ Inc. secret %v / %v (%v%%)\n",
			totalNonSecretCardsCollected,
			s.TotalNonSecretCards(),
			100*totalNonSecretCardsCollected/int(s.TotalNonSecretCards()),
			totalStarSecretCardsCollected,
			totalShinySecretCardsCollected,
			totalCrownSecretCardsCollected,
			totalCollectedIncludingSecrets,
			s.TotalCards(),
			100*(totalCollectedIncludingSecrets)/int(s.TotalCards()),
		)
	}

	// Show booster probabilities
	var allBoosters []boosterWithOrigin
	for _, s := range sets {
		missing, sExists := userCollection.MissingForSet(s.Id())
		if !sExists {
			fmt.Printf("Set id %v not found\n", s.Id())
			return
		}

		for b := range s.Boosters() {
			totalOfferingMissing := 0.0
			for o := range b.Offerings() {
				if slices.Contains(missing, o.Card().Number()) {
					totalOfferingMissing += o.RegularPackOffering()*0.9995 + o.RarePackOffering()*0.0005
				}
			}
			allBoosters = append(allBoosters, boosterWithOrigin{
				booster:              b,
				totalOfferingMissing: totalOfferingMissing,
				set:                  s,
			})
		}
	}
	slices.SortFunc(allBoosters, func(a boosterWithOrigin, b boosterWithOrigin) int {
		return int(1000*b.totalOfferingMissing) - int(1000*a.totalOfferingMissing)
	})

	fmt.Println()
	fmt.Println("# Booster probabilities")
	for i, b := range allBoosters {
		fmt.Printf("  %v) %.2f%% %v - %v\n", i+1, b.totalOfferingMissing, b.set.Name(), b.booster.Name())
	}
}

type boosterWithOrigin struct {
	set                  *data.CardSet
	booster              *data.Booster
	totalOfferingMissing float64
}
