package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"ptcgpocket/data"
	"ptcgpocket/sim"
	"ptcgpocket/source"
	"ptcgpocket/userdata"

	"golang.org/x/sync/errgroup"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer = message.NewPrinter(language.English)

var expansionDataSources = [...]*source.ExpansionSerebiiSource{
	source.NewExpansionSerebiiSource(
		"genetic-apex",
		"Genetic Apex",
		[]*source.BoosterSerebiiSource{
			source.NewBoosterSerebiiSource(
				"Pikachu",
				"https://www.serebii.net/tcgpocket/geneticapex/pikachu.shtml",
				data.OfferingRatesTable{
					data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/50.0, 0, 0, 0),
					data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 90.0/35.0, 60.0/35.0, 0),
					data.RarityThreeDiamond: *data.NewBoosterOffering(0, 5.0/14.0, 20.0/14.0, 0),
					data.RarityFourDiamond:  *data.NewBoosterOffering(0, 1.666/5.0, 6.664/5.0, 0),
					data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					data.RarityOneStar:      *data.NewBoosterOffering(0, 2.572/8.0, 10.288/8.0, 40.0/8.0),
					data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.50/10.0, 0.200/10.0, 50.0/10.0),
					data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222/1.0, 0.888/1.0, 5.0/1.0),
					data.RarityCrown:        *data.NewBoosterOffering(0, 0.4/3.0, 0.16/3.0, 5.0/1.0),
				},
				285,
			),
			source.NewBoosterSerebiiSource(
				"MewTwo",
				"https://www.serebii.net/tcgpocket/geneticapex/mewtwo.shtml",
				data.OfferingRatesTable{
					data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/50.0, 0, 0, 0),
					data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 90.0/35.0, 60.0/35.0, 0),
					data.RarityThreeDiamond: *data.NewBoosterOffering(0, 5.0/14.0, 20.0/14.0, 0),
					data.RarityFourDiamond:  *data.NewBoosterOffering(0, 1.666/5.0, 6.664/5.0, 0),
					data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					data.RarityOneStar:      *data.NewBoosterOffering(0, 2.572/8.0, 10.288/8.0, 42.105/8.0),
					data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.50/10.0, 0.200/10.0, 47.368/9.0),
					data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222/1.0, 0.888/1.0, 5.263/1.0),
					data.RarityCrown:        *data.NewBoosterOffering(0, 0.4/3.0, 0.16/3.0, 5.263/1.0),
				},
				286,
			),
			source.NewBoosterSerebiiSource(
				"Charizard",
				"https://www.serebii.net/tcgpocket/geneticapex/charizard.shtml",
				data.OfferingRatesTable{
					data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/50.0, 0, 0, 0),
					data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 90.0/35.0, 60.0/35.0, 0),
					data.RarityThreeDiamond: *data.NewBoosterOffering(0, 5.0/14.0, 20.0/14.0, 0),
					data.RarityFourDiamond:  *data.NewBoosterOffering(0, 1.666/5.0, 6.664/5.0, 0),
					data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					data.RarityOneStar:      *data.NewBoosterOffering(0, 2.572/8.0, 10.288/8.0, 40.0/8.0),
					data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.50/10.0, 0.200/10.0, 50.0/10.0),
					data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222/1.0, 0.888/1.0, 5.0/1.0),
					data.RarityCrown:        *data.NewBoosterOffering(0, 0.4/3.0, 0.16/3.0, 5.0/1.0),
				},
				284,
			),
		},
	),
	source.NewExpansionSerebiiSource(
		"mythical-island",
		"Mythical Island",
		[]*source.BoosterSerebiiSource{
			source.NewBoosterSerebiiSource(
				"Mew",
				"https://www.serebii.net/tcgpocket/mythicalisland/mew.shtml",
				data.OfferingRatesTable{
					data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/32.0, 0, 0, 0),
					data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 90.0/23.0, 60.0/23.0, 0),
					data.RarityThreeDiamond: *data.NewBoosterOffering(0, 5.0/8.0, 20.0/8.0, 0),
					data.RarityFourDiamond:  *data.NewBoosterOffering(0, 1.666/5, 6.664/5.0, 0),
					data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					data.RarityOneStar:      *data.NewBoosterOffering(0, 2.572/6.0, 10.288/6.0, 33.333/6.0),
					data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.500/10.0, 2.000/10.0, 55.555/10.0),
					data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222/1.0, 0.888/1.0, 5.555/1.0),
					data.RarityCrown:        *data.NewBoosterOffering(0, 0.040/1.0, 0.160/1.0, 5.555/1.0),
				},
				86,
			),
		},
	),
	source.NewExpansionSerebiiSource(
		"space-time-smackdown",
		"Space-time Smackdown",
		[]*source.BoosterSerebiiSource{
			source.NewBoosterSerebiiSource(
				"Dialga",
				"https://www.serebii.net/tcgpocket/space-timesmackdown/dialga.shtml",
				data.OfferingRatesTable{
					data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/46.0, 0, 0, 0),
					data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 90.0/34.0, 60.0/34.0, 0),
					data.RarityThreeDiamond: *data.NewBoosterOffering(0, 5.0/14.0, 20.0/14.0, 0),
					data.RarityFourDiamond:  *data.NewBoosterOffering(0, 1.666/5.0, 6.664/5.0, 0),
					data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					data.RarityOneStar:      *data.NewBoosterOffering(0, 2.572/12.0, 10.288/12.0, 46.153/12.0),
					data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.5/12.0, 2.0/12.0, 46.153/12.0),
					data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222/1.0, 0.888/1.0, 3.846/1.0),
					data.RarityCrown:        *data.NewBoosterOffering(0, 0.040/2.0, 0.160/2.0, 3.846/1.0),
				},
				207,
			),
			source.NewBoosterSerebiiSource(
				"Palkia",
				"https://www.serebii.net/tcgpocket/space-timesmackdown/palkia.shtml",
				data.OfferingRatesTable{
					data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/44.0, 0, 0, 0),
					data.RarityTwoDiamond:   *data.NewBoosterOffering(0, 90.0/36.0, 60.0/36.0, 0),
					data.RarityThreeDiamond: *data.NewBoosterOffering(0, 5.0/14.0, 20.0/14.0, 0),
					data.RarityFourDiamond:  *data.NewBoosterOffering(0, 1.666/5.0, 6.664/5.0, 0),
					data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					data.RarityOneStar:      *data.NewBoosterOffering(0, 2.572/12.0, 10.288/12.0, 46.153/12.0),
					data.RarityTwoStar:      *data.NewBoosterOffering(0, 0.5/12.0, 2.0/12.0, 46.153/12.0),
					data.RarityThreeStar:    *data.NewBoosterOffering(0, 0.222/1.0, 0.888/1.0, 3.846/1.0),
					data.RarityCrown:        *data.NewBoosterOffering(0, 0.040/2.0, 0.160/2.0, 3.846/1.0),
				},
				206,
			),
		},
	),
	source.NewExpansionSerebiiSource(
		"triumphant-light",
		"Triumphant Light",
		[]*source.BoosterSerebiiSource{
			source.NewBoosterSerebiiSource(
				"Arceus",
				"https://www.serebii.net/tcgpocket/triumphantlight/arceus.shtml",
				data.OfferingRatesTable{
					data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/31.0, 0, 0, 0),
					data.RarityTwoDiamond:   *data.NewBoosterOffering(0.000, 90.0/26.0, 60.0/26.0, 0.000),
					data.RarityThreeDiamond: *data.NewBoosterOffering(0.000, 5.0/13.0, 20.0/13.0, 0.000),
					data.RarityFourDiamond:  *data.NewBoosterOffering(0.000, 1.666/5.0, 6.664/5.0, 0.000),
					data.RarityOneShiny:     *data.NotPresentBoosterOffering,
					data.RarityTwoShiny:     *data.NotPresentBoosterOffering,
					data.RarityOneStar:      *data.NewBoosterOffering(0.000, 2.572/6.0, 10.288/6.0, 28.571/6.0),
					data.RarityTwoStar:      *data.NewBoosterOffering(0.000, 0.5/13.0, 2.0/13.0, 61.904/13.0),
					data.RarityThreeStar:    *data.NewBoosterOffering(0.000, 0.222/1.0, 0.888/1.0, 4.761/1.0),
					data.RarityCrown:        *data.NewBoosterOffering(0.000, 0.040/1.0, 0.160/1.0, 4.761/1.0),
				},
				96,
			),
		},
	),
	source.NewExpansionSerebiiSource(
		"shining-revelry",
		"Shining Revelry",
		[]*source.BoosterSerebiiSource{
			source.NewBoosterSerebiiSource(
				"Booster",
				"https://www.serebii.net/tcgpocket/shiningrevelry/booster.shtml",
				data.OfferingRatesTable{
					data.RarityOneDiamond:   *data.NewBoosterOffering(100.0/32.0, 0, 0, 0),
					data.RarityTwoDiamond:   *data.NewBoosterOffering(0.000, 89.000/22, 56.000/22, 0.000),
					data.RarityThreeDiamond: *data.NewBoosterOffering(0.000, 4.952/9, 19.810/9, 0.000),
					data.RarityFourDiamond:  *data.NewBoosterOffering(0.000, 1.666/9, 6.664/9, 0.000),
					data.RarityOneStar:      *data.NewBoosterOffering(0.000, 2.572/6, 10.288/6, 15.384/6),
					data.RarityTwoStar:      *data.NewBoosterOffering(0.000, 0.500/17, 2.000/17, 43.589/17),
					data.RarityThreeStar:    *data.NewBoosterOffering(0.000, 0.222/1, 0.888/1, 2.564/1),
					data.RarityOneShiny:     *data.NewBoosterOffering(0.000, 0.714/10, 2.857/10, 25.641/10),
					data.RarityTwoShiny:     *data.NewBoosterOffering(0.000, 0.333/4, 1.333/4, 10.256/4),
					data.RarityCrown:        *data.NewBoosterOffering(0.000, 0.040/1, 0.160, 2.564),
				},
				111,
			),
		},
	),
}

func readUserData(expansions []*data.Expansion) (*userdata.UserData, error) {
	dir, dErr := os.Getwd()
	if dErr != nil {
		return nil, dErr
	}

	dataFilepath := filepath.Join(dir, "data.json")

	return userdata.ReadFromFilepath(dataFilepath, expansions)
}

func printHeading1(heading string) {
	fmt.Printf("\033[1;32m# %v\n\033[0m", heading)
}

func printHeading2(heading string) {
	fmt.Printf("\033[32m  ## %v\n\033[0m", heading)
}

func printFullCardPosition(
	label string,
	amount float64,
	tallies map[*data.Rarity]float64,
) {
	allTallyDescriptions := make([]string, len(tallies))
	i := 0
	for _, r := range data.OrderedRarities {
		t, tFound := tallies[r]
		if tFound && t > 0 {
			allTallyDescriptions[i] = fmt.Sprintf("%v%.3f", r, t)
			i++
		}
	}
	tallyDescriptionRows := slices.Chunk(allTallyDescriptions, 5)

	// Note: Official numbers for Genetic Apex packs don't match up to 100%
	// for 4th or 5th cards
	colour := ""
	colourReset := ""
	if math.Abs(amount-100.0) > 0.1 {
		colour = "\033[0;31m"
		colourReset = "\033[0m"
	}

	fmt.Printf("%s   %v: %.2f / 100%%\n", colour, label, amount)
	for t := range tallyDescriptionRows {
		fmt.Printf("      %v%s\n", strings.Join(t, " "), colourReset)
	}
}

func printBoosterDataAudit(expansions []*data.Expansion) {
	printHeading1("Booster gathered data audit")
	for _, e := range expansions {
		for b := range e.Boosters() {

			offeringTallies := make(map[string]map[*data.Rarity]float64)
			offeringTallies["1-3"] = make(map[*data.Rarity]float64)
			offeringTallies["4"] = make(map[*data.Rarity]float64)
			offeringTallies["5"] = make(map[*data.Rarity]float64)
			offeringTallies["rare"] = make(map[*data.Rarity]float64)

			totalRegularPackOffering := 0.0
			totalRarePackOffering := 0.0
			totalFirstToThirdOffering := 0.0
			totalFourthOffering := 0.0
			totalFifthOffering := 0.0
			totalRareCardOffering := 0.0
			for c := range b.Offerings() {
				offeringTallies["1-3"][c.Card().Rarity()] += c.First3CardOffering()
				offeringTallies["4"][c.Card().Rarity()] += c.FourthCardOffering()
				offeringTallies["5"][c.Card().Rarity()] += c.FifthCardOffering()
				offeringTallies["rare"][c.Card().Rarity()] += c.RareCardOffering()

				totalRegularPackOffering += c.RegularPackOffering()
				totalRarePackOffering += c.RegularPackOffering()
				totalFirstToThirdOffering += c.First3CardOffering()
				totalFourthOffering += c.FourthCardOffering()
				totalFifthOffering += c.FifthCardOffering()
				totalRareCardOffering += c.RareCardOffering()
			}

			printHeading2(fmt.Sprintf("%v - %v", e.Name(), b.Name()))
			printFullCardPosition(
				"1-3",
				totalFirstToThirdOffering,
				offeringTallies["1-3"],
			)
			printFullCardPosition(
				"4",
				totalFourthOffering,
				offeringTallies["4"],
			)
			printFullCardPosition(
				"5",
				totalFifthOffering,
				offeringTallies["5"],
			)
			fmt.Printf("   total regular: %.2f / 500%%\n", totalRegularPackOffering)
			printFullCardPosition(
				"rare",
				totalRareCardOffering,
				offeringTallies["rare"],
			)
			fmt.Printf("   total rare: %.2f / 500%%\n", totalRarePackOffering)
			fmt.Println()
		}
	}
}

func printCurrentCollectionStats(expansions []*data.Expansion, userCollection *userdata.UserCollection) {
	printHeading1("Current collection")
	for _, e := range expansions {
		missing, sExists := userCollection.MissingForExpansion(e.Id())
		if !sExists {
			fmt.Printf("Set id %v not found\n", e.Id())
			return
		}

		printHeading2(e.Name())

		var totalStarSecretCardsCollected uint64
		var totalCrownSecretCardsCollected uint64
		var totalNonSecretCardsCollected uint64
		var totalShinySecretCardsCollected uint64
		for c := range e.Cards() {
			if !slices.Contains(missing, c) {
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
		totalCollectedIncludingSecrets := totalStarSecretCardsCollected +
			totalShinySecretCardsCollected +
			totalCrownSecretCardsCollected +
			totalNonSecretCardsCollected

		rarityTypes := map[rune]uint64{
			data.RarityStarChar:  totalCrownSecretCardsCollected,
			data.RarityCrownChar: totalCrownSecretCardsCollected,
		}
		if e.HasShiny() {
			rarityTypes[data.RarityShinyChar] = totalShinySecretCardsCollected
		}
		var rarityCounts []string
		for r, t := range rarityTypes {
			rarityCounts = append(rarityCounts, fmt.Sprintf("%v: %v", r, t))
		}

		fmt.Printf(
			"    %v / %v (%v%%) %v Inc. secret %v / %v (%v%%)\n",
			totalNonSecretCardsCollected,
			e.TotalNonSecretCards(),
			100*totalNonSecretCardsCollected/uint64(e.TotalNonSecretCards()),
			strings.Join(rarityCounts, " "),
			totalCollectedIncludingSecrets,
			e.TotalCards(),
			100*(totalCollectedIncludingSecrets)/uint64(e.TotalCards()),
		)
	}
}

func printBoosterProbabilities(
	heading string,
	getTargets func(e *data.Expansion) ([]*data.Card, bool),
	expansions []*data.Expansion,
) {
	var allBoosters []boosterWithOrigin
	for _, e := range expansions {
		missing, sExists := getTargets(e)
		if !sExists {
			continue
		}

		for b := range e.Boosters() {
			totalOfferingMissing := b.GetInstanceProbabilityForMissing(missing)
			allBoosters = append(allBoosters, boosterWithOrigin{
				booster:              b,
				totalOfferingMissing: totalOfferingMissing,
				expansion:            e,
			})
		}
	}
	slices.SortFunc(allBoosters, func(a, b boosterWithOrigin) int {
		return int(1000*b.totalOfferingMissing) - int(1000*a.totalOfferingMissing)
	})

	printHeading1(heading)
	for i, b := range allBoosters {
		fmt.Printf("  %v) %.2f%% %v - %v\n", i+1, b.totalOfferingMissing, b.expansion.Name(), b.booster.Name())
	}
}

type expansionSimRunAmounts struct {
	numOpened                      uint64
	totalPackPoints                uint64
	numCardsObtainedFromPackPoints uint64
	numRarePacks                   uint64
}

func runSimulations(
	title string,
	runMode *runOptions,
	expansions []*data.Expansion,
	userCollection *userdata.UserCollection,
	completePredicate sim.ExpansionSimCompletePredicate,
) error {
	printHeading1(printer.Sprintf("%v - pack opening simulations (%d runs)", title, runMode.simulationRuns))
	fmt.Printf("  Seed: %v\n", runMode.randomSeed)
	fmt.Println("  The number of booster openings required to complete the collection.")

	simResults := make(chan *sim.SimRun, runMode.simulationRuns)
	sim.RunAllSimulations(
		expansions,
		userCollection,
		completePredicate,
		runMode.simulationRuns,
		runMode.randomSeed,
		context.Background(),
		simResults,
	)
	close(simResults)

	expansionTotals := make(map[*data.Expansion]*expansionSimRunAmounts)
	var total uint64
	for r := range simResults {
		for e, run := range r.ExpansionRuns() {
			eTotals := expansionTotals[e]
			if eTotals == nil {
				eTotals = &expansionSimRunAmounts{}
				expansionTotals[e] = eTotals
			}

			eTotals.numOpened += run.NumOpened()
			eTotals.totalPackPoints += run.TotalPackPoints()
			eTotals.numCardsObtainedFromPackPoints += run.NumCardsObtainedFromPackPoints()
			eTotals.numRarePacks += run.NumRarePacks()
			total += run.NumOpened()
		}
	}
	expansionAverages := make(map[*data.Expansion]*expansionSimRunAmounts)
	var averagesTotal uint64
	for e, t := range expansionTotals {
		expansionAverages[e] = &expansionSimRunAmounts{
			numOpened:                      t.numOpened / runMode.simulationRuns,
			totalPackPoints:                t.totalPackPoints / runMode.simulationRuns,
			numCardsObtainedFromPackPoints: t.numCardsObtainedFromPackPoints / runMode.simulationRuns,
			numRarePacks:                   t.numRarePacks / runMode.simulationRuns,
		}
		averagesTotal += t.numOpened / runMode.simulationRuns
	}
	printer.Printf("  Total pack openings across all simulations: %d\n", total)
	fmt.Println()
	for e, a := range expansionAverages {
		printHeading2(e.Name())
		printer.Printf("     Packs opened        %v\n", a.numOpened)
		printer.Printf("     Rare packs          %v\n", a.numRarePacks)
		printer.Printf("     Cards from pack pts %v\n", a.numCardsObtainedFromPackPoints)
	}
	printer.Println()
	printHeading2(fmt.Sprintf("Total pack openings %d\n", averagesTotal))

	return nil
}

type runOptions struct {
	simulationRuns uint64
	randomSeed     uint64
}

func readRunOptions() (*runOptions, error) {
	simRunsPointer := flag.Uint64("r", 10, "number of sim runs")
	randomSeedPointer := flag.Uint64("s", rand.Uint64(), "sim random seed")
	flag.Parse()

	return &runOptions{simulationRuns: *simRunsPointer, randomSeed: *randomSeedPointer}, nil
}

func main() {
	runMode, rErr := readRunOptions()
	if rErr != nil {
		panic(rErr)
	}

	// Gather data from sources
	results := make(chan *data.Expansion, len(expansionDataSources))
	g, ctx := errgroup.WithContext(context.Background())
	indexMap := make(map[data.ExpansionId]int)
	for i, s := range expansionDataSources {
		indexMap[s.Id()] = i
		g.Go(func() error {
			return source.FetchExpansionDetails(ctx, s, results)
		})
	}
	err := g.Wait()
	close(results)
	if err != nil {
		panic(err)
	}

	var expansions []*data.Expansion
	for e := range results {
		expansions = append(expansions, e)
	}
	slices.SortFunc(expansions, func(e1, e2 *data.Expansion) int {
		return indexMap[e1.Id()] - indexMap[e2.Id()]
	})

	// Loading collection
	userData, uErr := readUserData(expansions)
	if uErr != nil {
		panic(uErr)
	}

	printBoosterDataAudit(expansions)
	fmt.Println()

	printCurrentCollectionStats(expansions, userData.Collection())
	fmt.Println()

	for w := range userData.Wishlists() {
		printBoosterProbabilities(
			fmt.Sprintf("Wishlist '%v' booster probabilities", w.Name()),
			func(e *data.Expansion) ([]*data.Card, bool) {
				return w.CardsForExpansion(e.Id())
			},
			expansions,
		)
		fmt.Println()
	}

	printBoosterProbabilities(
		"Collection booster probabilities",
		func(e *data.Expansion) ([]*data.Card, bool) {
			return userData.Collection().MissingForExpansion(e.Id())
		},
		expansions,
	)
	fmt.Println()

	runSimulations(
		"Whole collection",
		runMode,
		expansions,
		userData.Collection(),
		func(e *data.Expansion, m []*data.Card) bool {
			return len(m) == 0
		},
	)
	fmt.Println()

	runSimulations(
		"Non-secret cards collection",
		runMode,
		expansions,
		userData.Collection(),
		func(e *data.Expansion, m []*data.Card) bool {
			for _, card := range m {
				if !card.Rarity().IsSecret() {
					return false
				}
			}
			return true
		},
	)
}

type boosterWithOrigin struct {
	expansion            *data.Expansion
	booster              *data.Booster
	totalOfferingMissing float64
}
