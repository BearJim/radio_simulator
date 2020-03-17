package ngap

import "radio_simulator/lib/aper"

// Criticality ENUMURATED
type Criticality struct {
	Value aper.Enumerated `aper:"valueLB:0,valueUB:2"`
}

// Criticality ENUMURATED Value
const (
	CriticalityReject aper.Enumerated = iota
	CriticalityIgnore
	CriticalityNotify
)
