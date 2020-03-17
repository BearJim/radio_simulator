package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type PortNumber struct {
	Value aper.OctetString `aper:"sizeLB:2,sizeUB:2"`
}
