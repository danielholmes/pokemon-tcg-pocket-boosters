package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"ptcgpocket/collection"
	"ptcgpocket/data"
	"ptcgpocket/ref"
	"ptcgpocket/source"
)

var userCollection collection.UserCollection = collection.NewUserCollection(
	map[ref.CardSet]([]ref.CardSetNumber){
		ref.CardSetGeneticApex: {
			3, 4, 7, 10, 13, 22, 32, 36, 39, 41, 47, 50, 56, 61, 69, 73, 76, 80, 84, 89, 93, 95, 98, 101, 107, 117,
			123, 124, 145, 146, 148, 149, 159, 163, 166, 175, 177, 178, 185, 191, 195, 197, 202, 203, 204, 205, 221,
			225, 226,
			228, 229, 230, 231, 232, 233, 236, 237, 238, 240, 241, 242, 243, 244, 246, 248, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286,
		},
		ref.CardSetMythicalIsland: {
			2, 3, 6, 7, 18, 25, 26, 32, 44, 46, 59, 60, 62,
			71, 73, 75, 76, 79, 80, 81, 82, 83, 84, 85, 86,
		},
		ref.CardSetSpacetimeSmackdown: {
			5, 7, 18, 20, 22, 24, 29, 32, 33, 34, 36, 37, 41, 60, 65, 76, 79, 89, 90, 92, 94, 103, 104, 109, 113,
			117, 120, 123, 129, 147, 153,
			156, 157, 158, 159, 160, 161, 162, 164, 166, 167, 168, 169, 170, 171, 172, 173, 176, 177, 178, 179,
			180, 181, 182, 183, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 200, 201, 202,
			203, 205, 206, 207,
		},
	},
)

func fetchUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getOnlyDexTable(doc *html.Node) (*html.Node, error) {
	for d := range doc.Descendants() {
		if d.DataAtom == atom.Table {
			for _, a := range d.Attr {
				if a.Key == "class" && a.Val == "dextable" {
					return d, nil
				}
			}
		}
	}
	return nil, errors.New("no dextable found")
}

func getImmediateRows(table *html.Node) []*html.Node {
	rows := []*html.Node{}

	var tbody *html.Node
	for n := range table.ChildNodes() {
		if n.DataAtom == atom.Tbody {
			tbody = n
		}
	}
	if tbody == nil {
		return rows
	}

	i := 0
	for d := range tbody.ChildNodes() {
		if d.DataAtom == atom.Tr {
			if i != 0 {
				rows = append(rows, d)
			}
			i += 1
		}
	}
	return rows
}

func readFileIfExists(filename string) string {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return ""
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}

	return string(data)
}

func fetchBoosterFile(booster source.BoosterDataSource) (string, error) {
	parsed, uErr := url.Parse(booster.SerebiiUrl())
	if uErr != nil {
		return "", fmt.Errorf("error parsing URL: %v", uErr)
	}

	dir, dErr := os.Getwd()
	if dErr != nil {
		return "", dErr
	}

	cacheFilepath := filepath.Join(dir, ".cache", parsed.Hostname(), parsed.Path)
	var fileBody = readFileIfExists(cacheFilepath)
	if fileBody != "" {
		return fileBody, nil
	}

	fmt.Printf("No cached file found for %v, fetching %v\n", booster.Name(), booster.SerebiiUrl())
	var body, err = fetchUrl(booster.SerebiiUrl())
	if err != nil {
		return "", err
	}

	cacheFileDirpath := filepath.Dir(cacheFilepath)
	mDErr := os.MkdirAll(cacheFileDirpath, 0755)
	if mDErr != nil {
		return "", mDErr
	}

	wErr := os.WriteFile(cacheFilepath, []byte(body), 0755)
	if wErr != nil {
		return "", wErr
	}

	return body, nil
}

func fetchBoosterDetails(booster source.BoosterDataSource, wg *sync.WaitGroup, results chan<- data.Booster) {
	defer wg.Done()

	var body, err = fetchBoosterFile(booster)
	// TODO: Find idiomatic way to handle go routine errors
	if err != nil {
		fmt.Println(err)
		return
	}

	var doc, dErr = html.Parse(strings.NewReader(body))
	if dErr != nil {
		fmt.Println(dErr)
		return
	}

	var table, tErr = getOnlyDexTable(doc)
	if tErr != nil {
		fmt.Println(tErr)
		return
	}

	rows := getImmediateRows(table)

	offerings := make([]data.BoosterOffering, len(rows))
	for i, r := range rows {
		cells := []*html.Node{}
		for c := range r.ChildNodes() {
			if c.DataAtom == atom.Td {
				cells = append(cells, c)
			}
		}

		if len(cells) != 7 {
			fmt.Printf("Expected 7 cells in booster row (got %v) %v\n", len(cells), r)
			return
		}

		// Number
		var number ref.CardSetNumber
		numRe := regexp.MustCompile("([0-9]+?) / ([0-9]+)")
		var imageNode *html.Node
		for d := range cells[0].Descendants() {
			if d.DataAtom == atom.Img {
				imageNode = d
			}

			dMatch := numRe.FindStringSubmatch(d.Data)
			if dMatch != nil {
				value, _ := strconv.ParseUint(dMatch[1], 10, 16)
				number = ref.CardSetNumber(value)
			}
		}

		// Name
		var name string
		for d := range cells[2].Descendants() {
			if d.DataAtom == atom.A {
				nameParent := d

				// Look for font
				for c := range d.ChildNodes() {
					if c.DataAtom == atom.Font {
						nameParent = c
					}
				}

				// No font, fallback is root
				if nameParent == nil {
					nameParent = nameParent.FirstChild
				}

				name = nameParent.FirstChild.Data
			}
		}
		if name == "" {
			fmt.Println("No name")
			return
		}

		// rarity
		// diamond1 diamond2 diamond3 diamond4 star1 star2 star3 crown
		if imageNode == nil {
			fmt.Printf("no img found %v) %v\n", number, name)
			return
		}
		srcAttrIndex := slices.IndexFunc(imageNode.Attr, func(a html.Attribute) bool {
			return a.Key == "src"
		})
		if srcAttrIndex == -1 {
			fmt.Printf("no img src found %v) %v\n", number, name)
			return
		}
		srcAttr := imageNode.Attr[srcAttrIndex]
		comps := strings.Split(srcAttr.Val, "/")
		imageName := strings.Split(comps[len(comps)-1], ".")[0]
		imageNameRarities := map[string]*data.Rarity{
			"diamond1": &data.RarityOneDiamond,
			"diamond2": &data.RarityTwoDiamond,
			"diamond3": &data.RarityThreeDiamond,
			"diamond4": &data.RarityFourDiamond,
			"star1":    &data.RarityOneStar,
			"star2":    &data.RarityTwoStar,
			"star3":    &data.RarityThreeStar,
			"crown":    &data.RarityCrown,
		}
		var rarity *data.Rarity = imageNameRarities[imageName]
		if rarity == nil {
			fmt.Printf("no rarity found for image name %v\n", imageName)
			return
		}

		offeringFix, offeringFixExists := booster.OfferingFixForNumber(number)

		// 1-3 probability
		var firstThreeOffering float64
		if offeringFixExists {
			firstThreeOffering = offeringFix[0]
		} else {
			firstThreeCell := cells[4].FirstChild
			if firstThreeCell != nil {
				rawRarity := firstThreeCell.Data
				firstThreeOffering, _ = strconv.ParseFloat(rawRarity[:len(rawRarity)-1], 64)
			}
		}

		// 4 probability
		var fourthOffering float64
		if offeringFixExists {
			fourthOffering = offeringFix[1]
		} else {
			fourthCell := cells[5].FirstChild
			if fourthCell != nil {
				rawRarity := fourthCell.Data
				fourthOffering, _ = strconv.ParseFloat(rawRarity[:len(rawRarity)-1], 64)
			}
		}

		// 5 probability
		var fifthOffering float64
		if offeringFixExists {
			fifthOffering = offeringFix[2]
		} else {
			fifthCell := cells[6].FirstChild
			if fifthCell != nil {
				rawRarity := fifthCell.Data
				fifthOffering, _ = strconv.ParseFloat(rawRarity[:len(rawRarity)-1], 64)
			}
		}

		offerings[i] = data.NewBoosterOffering(
			data.NewCard(
				name,
				number,
				rarity,
			),
			firstThreeOffering,
			fourthOffering,
			fifthOffering,
		)
	}

	results <- data.NewBooster(
		booster.Name(),
		offerings,
	)
}

func fetchCardSetDetails(s source.CardSetDataSource, wg *sync.WaitGroup, results chan<- data.CardSetDetails) {
	defer wg.Done()

	var bwg sync.WaitGroup
	boosterResults := make(chan data.Booster, s.NumBoosterSources())
	for s := range s.BoosterSources() {
		bwg.Add(1)
		go fetchBoosterDetails(s, &bwg, boosterResults)
	}

	bwg.Wait()
	close(boosterResults)

	var collectedResults []data.Booster
	for o := range boosterResults {
		collectedResults = append(collectedResults, o)
	}

	results <- data.NewCardSetDetails(s.Set(), collectedResults)
}

func main() {
	var wg sync.WaitGroup

	results := make(chan data.CardSetDetails, len(source.CardSetDataSources))
	for _, s := range source.CardSetDataSources {
		wg.Add(1)
		go fetchCardSetDetails(s, &wg, results)
	}

	wg.Wait()
	close(results)

	var cardDetails []data.CardSetDetails
	for o := range results {
		cardDetails = append(cardDetails, o)
	}

	// Show check of booster probabilities
	fmt.Println("# Booster gathered data audit")
	for _, c := range cardDetails {
		for b := range c.Boosters() {
			totalOffering := 0.0
			totalFirstToThirdOffering := 0.0
			totalFourthOffering := 0.0
			totalFifthOffering := 0.0
			for o := range b.Offerings() {
				totalOffering += o.PackProbability()
				totalFirstToThirdOffering += o.First3CardOffering()
				totalFourthOffering += o.FourthCardOffering()
				totalFifthOffering += o.FifthCardOffering()
			}
			fmt.Printf(
				" ## %v - %v\n   1-3: %.2f / 100%%\n   4: %.2f / 100%%\n   5: %.2f / 100%%\n   total: %.2f / 500%%\n",
				c.Set().Name(),
				b.Name(),
				totalFirstToThirdOffering,
				totalFourthOffering,
				totalFifthOffering,
				totalOffering,
			)
		}
	}

	// Show collection
	fmt.Println()
	fmt.Println("# Current collection")
	for set, missing := range userCollection.MissingCardNumbers() {
		fmt.Printf(" ## %v\n", set.Name())
		var setDetails *data.CardSetDetails
		for _, d := range cardDetails {
			if *d.Set() == set {
				setDetails = &d
			}
		}
		if setDetails == nil {
			panic(fmt.Sprintf("no set details for %v", set.Name()))
		}

		totalSecretCardsCollected := 0
		totalNonSecretCardsCollected := 0
		for c := range setDetails.Cards() {
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
			setDetails.TotalNonSecretCards(),
			100*totalNonSecretCardsCollected/int(setDetails.TotalNonSecretCards()),
			totalSecretCardsCollected,
			totalCollectedIncludingSecrets,
			setDetails.TotalCards(),
			100*(totalSecretCardsCollected+totalNonSecretCardsCollected)/int(setDetails.TotalCards()),
		)
	}

	// Show booster values
	fmt.Println()
	fmt.Println("# Booster values")
	for set, missing := range userCollection.MissingCardNumbers() {
		fmt.Printf(" ## %v\n", set.Name())
		var setDetails *data.CardSetDetails
		for _, d := range cardDetails {
			if *d.Set() == set {
				setDetails = &d
			}
		}
		if setDetails == nil {
			panic(fmt.Sprintf("no set details for %v", set.Name()))
		}

		for b := range setDetails.Boosters() {
			fmt.Printf("  ### %v\n", b.Name())

			totalOfferingMissing := 0.0
			for o := range b.Offerings() {
				if slices.Contains(missing, o.Card().Number()) {
					totalOfferingMissing += o.PackProbability()
					// fmt.Printf("  Offers missing %v) %v => %v\n", c.card.number, c.card.name, c.offering)
				}
			}
			fmt.Printf("   Total chance of receiving a missing %.2f%%\n", totalOfferingMissing)
		}
	}
}
