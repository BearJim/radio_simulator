package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

const (
	EmergencyServiceTargetCNPresentFiveGC aper.Enumerated = 0
	EmergencyServiceTargetCNPresentEpc    aper.Enumerated = 1
)

type EmergencyServiceTargetCN struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}