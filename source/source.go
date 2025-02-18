package source

import (
	"iter"
	"ptcgpocket/data"
	"slices"
)

type BoosterSerebiiSource struct {
	name          string
	serebiiUrl    string
	offeringRates data.OfferingRatesTable
}

func NewBoosterSerebiiSource(
	name string,
	serebiiUrl string,
	offeringRates data.OfferingRatesTable,
) *BoosterSerebiiSource {
	return &BoosterSerebiiSource{
		name:          name,
		serebiiUrl:    serebiiUrl,
		offeringRates: offeringRates,
	}
}

func (b *BoosterSerebiiSource) Name() string {
	return b.name
}

func (b *BoosterSerebiiSource) SerebiiUrl() string {
	return b.serebiiUrl
}

func (b *BoosterSerebiiSource) OfferingRates() data.OfferingRatesTable {
	return b.offeringRates
}

type CardSetSerebiiSource struct {
	id             string
	name           string
	boosterSources []*BoosterSerebiiSource
}

func NewCardSetSerebiiSource(
	id string,
	name string,
	boosterSources []*BoosterSerebiiSource,
) *CardSetSerebiiSource {
	return &CardSetSerebiiSource{
		id:             id,
		name:           name,
		boosterSources: boosterSources,
	}
}

func (s *CardSetSerebiiSource) Id() string {
	return s.id
}

func (s *CardSetSerebiiSource) Name() string {
	return s.name
}

func (s *CardSetSerebiiSource) BoosterSources() iter.Seq[*BoosterSerebiiSource] {
	return slices.Values(s.boosterSources)
}

func (s *CardSetSerebiiSource) NumBoosterSources() uint8 {
	return uint8(len(s.boosterSources))
}
