package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

const (
	PrivateIEIDPresentNothing int = iota /* No components present */
	PrivateIEIDPresentLocal
	PrivateIEIDPresentGlobal
)

type PrivateIEID struct {
	Present int
	Local   *int64 `aper:"valueLB:0,valueUB:65535"`
	Global  *aper.ObjectIdentifier
}
