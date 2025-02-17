package ref

type CardSetId struct {
	value string
}

type CardSet struct {
	id CardSetId
	name string
}

func (s *CardSet) Name() string {
    return s.name
}

type CardSetNumber uint16

var (
	CardSetIdGeneticApex CardSetId = CardSetId{"genetic-apex"}
	CardSetIdMythicalIsland CardSetId = CardSetId{"mythical-island"}
	CardSetIdSpaceTimeSmackdown CardSetId = CardSetId{"space-time-smackdown"}
)

var (
	CardSetGeneticApex CardSet = CardSet{id: CardSetIdGeneticApex, name: "Genetic Apex"}
	CardSetMythicalIsland CardSet = CardSet{id: CardSetIdMythicalIsland, name: "Mythical Island"}
	CardSetSpacetimeSmackdown CardSet = CardSet{id: CardSetIdSpaceTimeSmackdown, name: "Space-Time Smackdown"}
)
