package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

const (
	SourceOfUEActivityBehaviourInformationPresentSubscriptionInformation aper.Enumerated = 0
	SourceOfUEActivityBehaviourInformationPresentStatistics              aper.Enumerated = 1
)

type SourceOfUEActivityBehaviourInformation struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}
