package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

const (
	OverloadActionPresentRejectNonEmergencyMoDt                                    aper.Enumerated = 0
	OverloadActionPresentRejectRrcCrSignalling                                     aper.Enumerated = 1
	OverloadActionPresentPermitEmergencySessionsAndMobileTerminatedServicesOnly    aper.Enumerated = 2
	OverloadActionPresentPermitHighPrioritySessionsAndMobileTerminatedServicesOnly aper.Enumerated = 3
)

type OverloadAction struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:3"`
}
