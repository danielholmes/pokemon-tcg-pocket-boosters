package source

import (
	"ptcgpocket/ref"
)

type BoosterDataSource struct {
	Name       string
	SerebiiUrl string
}

type CardSetDataSource struct {
	Set            ref.CardSet
	BoosterSources []BoosterDataSource
}

var CardSetDataSources = [...]CardSetDataSource{
	{
		Set: ref.CardSetGeneticApex,
		BoosterSources: []BoosterDataSource{
			{Name: "Pikachu", SerebiiUrl: "https://www.serebii.net/tcgpocket/geneticapex/pikachu.shtml"},
			{Name: "MewTwo", SerebiiUrl: "https://www.serebii.net/tcgpocket/geneticapex/mewtwo.shtml"},
			{Name: "Charizard", SerebiiUrl: "https://www.serebii.net/tcgpocket/geneticapex/charizard.shtml"},
		},
	},
	{
		Set: ref.CardSetMythicalIsland,
		BoosterSources: []BoosterDataSource{
			{Name: "Mew", SerebiiUrl: "https://www.serebii.net/tcgpocket/mythicalisland/mew.shtml"},
		},
	},
	{
		Set: ref.CardSetSpacetimeSmackdown,
		BoosterSources: []BoosterDataSource{
			{Name: "Dialga", SerebiiUrl: "https://www.serebii.net/tcgpocket/space-timesmackdown/dialga.shtml"},
			{Name: "Palkia", SerebiiUrl: "https://www.serebii.net/tcgpocket/space-timesmackdown/palkia.shtml"},
		},
	},
}
