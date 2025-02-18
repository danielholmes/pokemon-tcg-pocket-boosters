package source

import (
	"iter"
	"ptcgpocket/ref"
	"slices"
)

type offeringFix [3]float64

type BoosterDataSource struct {
	name                 string
	serebiiUrl           string
	serebiiOfferingFixes map[ref.CardSetNumber]offeringFix
}

func (b *BoosterDataSource) Name() string {
	return b.name
}

func (b *BoosterDataSource) SerebiiUrl() string {
	return b.serebiiUrl
}

func (b *BoosterDataSource) OfferingFixForNumber(number ref.CardSetNumber) (offeringFix, bool) {
	value, exists := b.serebiiOfferingFixes[number]
	return value, exists
}

type CardSetDataSource struct {
	set            ref.CardSet
	boosterSources []BoosterDataSource
}

func (s *CardSetDataSource) Set() ref.CardSet {
	return s.set
}

func (s *CardSetDataSource) BoosterSources() iter.Seq[BoosterDataSource] {
	return slices.Values(s.boosterSources)
}

func (s *CardSetDataSource) NumBoosterSources() uint8 {
	return uint8(len(s.boosterSources))
}

var CardSetDataSources = [...]CardSetDataSource{
	{
		set: ref.CardSetGeneticApex,
		boosterSources: []BoosterDataSource{
			{name: "Pikachu", serebiiUrl: "https://www.serebii.net/tcgpocket/geneticapex/pikachu.shtml"},
			{name: "MewTwo", serebiiUrl: "https://www.serebii.net/tcgpocket/geneticapex/mewtwo.shtml"},
			{name: "Charizard", serebiiUrl: "https://www.serebii.net/tcgpocket/geneticapex/charizard.shtml"},
		},
	},
	{
		set: ref.CardSetMythicalIsland,
		boosterSources: []BoosterDataSource{
			{name: "Mew", serebiiOfferingFixes: map[ref.CardSetNumber]offeringFix{
				71: {0.0, 0.428, 1.714},
			}, serebiiUrl: "https://www.serebii.net/tcgpocket/mythicalisland/mew.shtml"},
		},
	},
	{
		set: ref.CardSetSpacetimeSmackdown,
		boosterSources: []BoosterDataSource{
			{name: "Dialga", serebiiUrl: "https://www.serebii.net/tcgpocket/space-timesmackdown/dialga.shtml"},
			{name: "Palkia", serebiiUrl: "https://www.serebii.net/tcgpocket/space-timesmackdown/palkia.shtml"},
		},
	},
}
