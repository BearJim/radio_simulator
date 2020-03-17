package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

const (
	RRCInactiveTransitionReportRequestPresentSubsequentStateTransitionReport aper.Enumerated = 0
	RRCInactiveTransitionReportRequestPresentSingleRrcConnectedStateReport   aper.Enumerated = 1
	RRCInactiveTransitionReportRequestPresentCancelReport                    aper.Enumerated = 2
)

type RRCInactiveTransitionReportRequest struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:2"`
}
