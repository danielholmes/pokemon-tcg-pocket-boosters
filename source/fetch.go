package source

import (
	"context"
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

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/sync/errgroup"
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

func fetchBoosterDetails(booster *BoosterSerebiiSource, results chan<- *data.Booster) error {
	var body, err = fetchBoosterFile(booster)
	// TODO: Find idiomatic way to handle go routine errors
	if err != nil {
		return err
	}

	var doc, dErr = html.Parse(strings.NewReader(body))
	if dErr != nil {
		return dErr
	}

	var table, tErr = getOnlyDexTable(doc)
	if tErr != nil {
		return tErr
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
			return fmt.Errorf("expected 7 cells in booster row (got %v) %v - row %v", len(cells), booster.name, i)
		}

		// Number
		var number data.ExpansionNumber
		numRe := regexp.MustCompile("([0-9]+?) / ([0-9]+)")
		var imageNode *html.Node
		for d := range cells[0].Descendants() {
			if d.DataAtom == atom.Img {
				imageNode = d
			}

			dMatch := numRe.FindStringSubmatch(d.Data)
			if dMatch != nil {
				value, _ := strconv.ParseUint(dMatch[1], 10, 16)
				number = data.ExpansionNumber(value)
			}
		}

		// Name
		var name string
		for d := range cells[2].Descendants() {
			if d.DataAtom == atom.A {
				var nameComponents []string
				for c := range d.Descendants() {
					if c.DataAtom == 0 {
						comp := strings.TrimSpace(c.Data)
						if comp != "" {
							nameComponents = append(nameComponents, comp)
						}
					}
				}
				name = strings.Join(nameComponents, " ")
			}
		}
		if name == "" {
			return errors.New("no name")
		}

		// Rarity
		if imageNode == nil {
			return fmt.Errorf("no img found %v) %v", number, name)
		}
		srcAttrIndex := slices.IndexFunc(imageNode.Attr, func(a html.Attribute) bool {
			return a.Key == "src"
		})
		if srcAttrIndex == -1 {
			return fmt.Errorf("no img src found %v) %v", number, name)
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
			"shiny1":   &data.RarityOneShiny,
			"shiny2":   &data.RarityTwoShiny,
			"crown":    &data.RarityCrown,
		}
		var rarity *data.Rarity = imageNameRarities[imageName]
		if rarity == nil {
			return fmt.Errorf("no rarity found for image name %v", imageName)
		}

		card := data.NewCard(
			data.NewBaseCard(name, 0),
			number,
			rarity,
		)
		cards[i] = card
	}

	results <- data.NewBooster(
		booster.Name(),
		cards,
		booster.OfferingRates(),
		booster.CrownExclusiveExpansionNumber(),
	)
	return nil
}

func FetchExpansionDetails(ctx context.Context, s *ExpansionSerebiiSource, results chan<- *data.Expansion) error {
	g, _ := errgroup.WithContext(ctx)

	boosterResults := make(chan *data.Booster, s.NumBoosterSources())
	var boosterSources []*BoosterSerebiiSource
	for s := range s.BoosterSources() {
		boosterSources = append(boosterSources, s)
		g.Go(func() error {
			return fetchBoosterDetails(s, boosterResults)
		})
	}
	err := g.Wait()
	if err != nil {
		return err
	}

	close(boosterResults)

	var boosters []*data.Booster
	for b := range boosterResults {
		boosters = append(boosters, b)
	}
	slices.SortFunc(
		boosters,
		func(b1, b2 *data.Booster) int {
			b1Index := -1
			b2Index := -1
			for i, s := range boosterSources {
				if s.name == b1.Name() {
					b1Index = i
				}
				if s.name == b2.Name() {
					b2Index = i
				}
			}
			if b1Index == -1 || b2Index == -1 {
				panic("error sorting fetched boosters")
			}
			return b1Index - b2Index
		},
	)

	results <- data.NewExpansion(s.Id(), s.Name(), boosters)
	return nil
}
