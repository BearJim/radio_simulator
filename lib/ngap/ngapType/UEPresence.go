package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

const (
	UEPresencePresentIn      aper.Enumerated = 0
	UEPresencePresentOut     aper.Enumerated = 1
	UEPresencePresentUnknown aper.Enumerated = 2
)

type UEPresence struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:2"`
}
