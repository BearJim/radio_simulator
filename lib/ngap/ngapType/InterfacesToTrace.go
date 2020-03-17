package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type InterfacesToTrace struct {
	Value aper.BitString `aper:"sizeLB:8,sizeUB:8"`
}
