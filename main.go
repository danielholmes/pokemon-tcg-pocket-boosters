package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"os"
	"net/url"
	"path/filepath"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type booster struct {
	name       string
	serebiiUrl string
}

type cardSet struct {
	name     string
	boosters []booster
}

type card struct {
	name   string
	number uint16
}

type cardOffering struct {
	card     card
	offering float64
}

type boosterOfferings struct {
	booster   booster
	offerings []cardOffering
}

var geneticApexSet = cardSet{
	name: "Genetic Apex",
	boosters: []booster{
		{name: "Pikachu", serebiiUrl: "https://www.serebii.net/tcgpocket/geneticapex/pikachu.shtml"},
		{name: "MewTwo", serebiiUrl: "https://www.serebii.net/tcgpocket/geneticapex/mewtwo.shtml"},
		{name: "Charizard", serebiiUrl: "https://www.serebii.net/tcgpocket/geneticapex/charizard.shtml"},
	},
}

var mythicalIslandSet = cardSet{name: "Mythical Island",
	boosters: []booster{
		{name: "Mew", serebiiUrl: "https://www.serebii.net/tcgpocket/mythicalisland/mew.shtml"},
	}}

var spaceTimeSmackdownSet = cardSet{name: "Space-Time Smackdown",
	boosters: []booster{
		{name: "Dialga", serebiiUrl: "https://www.serebii.net/tcgpocket/space-timesmackdown/dialga.shtml"},
		{name: "Palkia", serebiiUrl: "https://www.serebii.net/tcgpocket/space-timesmackdown/palkia.shtml"},
	}}

var sets = [...]cardSet{
	geneticApexSet,
	mythicalIslandSet,
	spaceTimeSmackdownSet,
}

type missingInSet struct {
	cardSet cardSet
	missing []uint16
}

var missing = [...]missingInSet{
	{cardSet: geneticApexSet, missing: []uint16{
		3, 4, 7, 10, 13, 20, 22, 36, 37, 39, 41, 47, 50, 56, 61, 69, 73, 76, 80, 84, 86, 89, 93, 95, 98, 101, 107, 117, 123, 124, 145, 146, 148, 149, 159, 163, 166, 175, 177, 178, 185, 191, 195, 197, 202, 203, 204, 205, 221, 225, 226,
		228, 229, 230, 231, 232, 233, 236, 237, 238, 240, 241, 242, 243, 244, 246, 248, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286,
	}},
	{cardSet: mythicalIslandSet, missing: []uint16{
		2, 3, 6, 7, 18, 19, 25, 26, 32, 44, 46, 59, 60, 62,
		68, 71, 73, 75, 76, 79, 80, 81, 82, 83, 84, 85, 86,
	}},
	{cardSet: spaceTimeSmackdownSet, missing: []uint16{
		5, 6, 7, 18, 20, 22, 24, 29, 32, 33, 34, 36, 37, 40, 41, 60, 65, 76, 79, 89, 90, 92, 94, 103, 104, 109, 113, 117, 120, 123, 129, 147, 153,
		155, 156, 157, 158, 159, 160, 161, 162, 164, 166, 167, 168, 169, 170, 171, 172, 173, 176, 177, 178, 179, 180, 181, 182, 183, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 200, 201, 202, 203, 205, 206, 207,
	}},
}

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

func fetchBoosterFile(booster booster) (string, error) {
	parsed, uErr := url.Parse(booster.serebiiUrl)
	if uErr != nil {
		return "", fmt.Errorf("error parsing URL: %v", uErr)
	}

	dir, dErr := os.Getwd()
	if dErr != nil {
		return "", dErr
	}

	cacheFilepath := filepath.Join(dir, ".cache", parsed.Path)
	var fileBody = readFileIfExists(cacheFilepath)
	if fileBody != "" {
		return fileBody, nil
	}

	fmt.Printf("No cached file found for %v, fetching %v\n", booster.name, booster.serebiiUrl)
	var body, err = fetchUrl(booster.serebiiUrl)
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

func fetchBoosterOfferings(booster booster, wg *sync.WaitGroup, results chan<- boosterOfferings) {
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

	offerings := make([]cardOffering, len(rows))
	totalOffering := 0.0
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
		var number uint16
		numRe := regexp.MustCompile("([0-9]+?) / ([0-9]+)")
		for d := range cells[0].Descendants() {
			dMatch := numRe.FindStringSubmatch(d.Data)
			if dMatch != nil {
				value, _ := strconv.ParseUint(dMatch[1], 10, 16)
				number = uint16(value)
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

		// 1-3 rarity
		var firstThreeRarity float64
		firstThreeCell := cells[4].FirstChild
		if firstThreeCell != nil {
			rawRarity := firstThreeCell.Data
			firstThreeRarity, _ = strconv.ParseFloat(rawRarity[:len(rawRarity)-1], 64)
		}

		// 4 rarity
		var fourthRarity float64
		fourthCell := cells[5].FirstChild
		if fourthCell != nil {
			rawRarity := fourthCell.Data
			fourthRarity, _ = strconv.ParseFloat(rawRarity[:len(rawRarity)-1], 64)
		}

		// 5 rarity
		var fifthRarity float64
		fifthCell := cells[6].FirstChild
		if fourthCell != nil {
			rawRarity := fifthCell.Data
			fifthRarity, _ = strconv.ParseFloat(rawRarity[:len(rawRarity)-1], 64)
		}

		currentOffering := firstThreeRarity*3 + fourthRarity + fifthRarity
		offerings[i] = cardOffering{
			card: card{
				name:   name,
				number: number,
			},
			offering: currentOffering,
		}

		totalOffering += currentOffering
	}

	fmt.Printf("Booster %s total offerings (should = 500%%) %v\n", booster.name, totalOffering)
	results <- boosterOfferings{
		booster:   booster,
		offerings: offerings,
	}
}

func contains(slice []uint16, value uint16) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func main() {
	var wg sync.WaitGroup

	var allBoosters []booster
	for _, set := range sets {
		allBoosters = append(allBoosters, set.boosters...)
	}
	results := make(chan boosterOfferings, len(allBoosters))

	for _, booster := range allBoosters {
		wg.Add(1)
		go fetchBoosterOfferings(booster, &wg, results)
	}

	wg.Wait()
	close(results)

	var collectedResults []boosterOfferings
	for o := range results {
		collectedResults = append(collectedResults, o)
	}

	for _, m := range missing {
		fmt.Printf("Missing %v\n", m.cardSet.name)
		for _, b := range m.cardSet.boosters {
			fmt.Printf(" Booster %v\n", b.name)
			for _, o := range collectedResults {
				if o.booster == b {
					totalOfferingMissing := 0.0
					for _, c := range o.offerings {
						if contains(m.missing, c.card.number) {
							totalOfferingMissing += c.offering
							// fmt.Printf("  Offers missing %v) %v => %v\n", c.card.number, c.card.name, c.offering)
						}
					}
					fmt.Printf("  Total chance of receiving a missing %v%%\n", totalOfferingMissing)
				}
			}
		}
	}
}
