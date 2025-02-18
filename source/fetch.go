package source

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"ptcgpocket/data"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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

func fetchBoosterFile(booster *BoosterSerebiiSource) (string, error) {
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

func fetchBoosterDetails(booster *BoosterSerebiiSource, wg *sync.WaitGroup, results chan<- data.Booster) {
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

	cards := make([]*data.Card, len(rows))
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
		var number data.CardSetNumber
		numRe := regexp.MustCompile("([0-9]+?) / ([0-9]+)")
		var imageNode *html.Node
		for d := range cells[0].Descendants() {
			if d.DataAtom == atom.Img {
				imageNode = d
			}

			dMatch := numRe.FindStringSubmatch(d.Data)
			if dMatch != nil {
				value, _ := strconv.ParseUint(dMatch[1], 10, 16)
				number = data.CardSetNumber(value)
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

		card := data.NewCard(
			name,
			number,
			rarity,
		)
		cards[i] = &card
	}

	results <- data.NewBooster(
		booster.Name(),
		cards,
		booster.OfferingRates(),
	)
}

func FetchCardSetDetails(s *CardSetSerebiiSource, wg *sync.WaitGroup, results chan<- data.CardSet) {
	defer wg.Done()

	var bwg sync.WaitGroup
	boosterResults := make(chan data.Booster, s.NumBoosterSources())
	for s := range s.BoosterSources() {
		bwg.Add(1)
		go fetchBoosterDetails(s, &bwg, boosterResults)
	}

	bwg.Wait()
	close(boosterResults)

	var collectedResults []*data.Booster
	for o := range boosterResults {
		collectedResults = append(collectedResults, &o)
	}

	results <- data.NewCardSet(s.Id(), s.Name(), collectedResults)
}
