package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type GTPTEID struct {
	Value aper.OctetString `aper:"sizeLB:4,sizeUB:4"`
}