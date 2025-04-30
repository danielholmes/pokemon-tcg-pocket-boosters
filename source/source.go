package source

import (
	"iter"
	"ptcgpocket/data"
	"slices"
)

type BoosterSerebiiSource struct {
	name                                  string
	serebiiUrl                            string
	offeringRates                         data.OfferingRatesTable
	rarePackCrownExclusiveExpansionNumber data.ExpansionCardNumber
}

func NewBoosterSerebiiSource(
	name string,
	serebiiUrl string,
	offeringRates data.OfferingRatesTable,
	rarePackCrownExclusiveExpansionNumber data.ExpansionCardNumber,
) *BoosterSerebiiSource {
	return &BoosterSerebiiSource{
		name:                                  name,
		serebiiUrl:                            serebiiUrl,
		offeringRates:                         offeringRates,
		rarePackCrownExclusiveExpansionNumber: rarePackCrownExclusiveExpansionNumber,
	}
}

func (b *BoosterSerebiiSource) Name() string {
	return b.name
}

func (b *BoosterSerebiiSource) RarePackCrownExclusiveExpansionNumber() data.ExpansionCardNumber {
	return b.rarePackCrownExclusiveExpansionNumber
}

func (b *BoosterSerebiiSource) SerebiiUrl() string {
	return b.serebiiUrl
}

func (b *BoosterSerebiiSource) OfferingRates() data.OfferingRatesTable {
	return b.offeringRates
}

type ExpansionSerebiiSource struct {
	id             string
	name           string
	code           string
	boosterSources []*BoosterSerebiiSource
}

func NewExpansionSerebiiSource(
	id string,
	name string,
	code string,
	boosterSources []*BoosterSerebiiSource,
) *ExpansionSerebiiSource {
	return &ExpansionSerebiiSource{
		id:             id,
		name:           name,
		code:           code,
		boosterSources: boosterSources,
	}
}

func (s *ExpansionSerebiiSource) Id() string {
	return s.id
}

func (s *ExpansionSerebiiSource) Name() string {
	return s.name
}

func (s *ExpansionSerebiiSource) Code() string {
	return s.code
}

func (s *ExpansionSerebiiSource) BoosterSources() iter.Seq[*BoosterSerebiiSource] {
	return slices.Values(s.boosterSources)
}

func (s *ExpansionSerebiiSource) NumBoosterSources() uint8 {
	return uint8(len(s.boosterSources))
}
