package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type EmergencyAreaID struct {
	Value aper.OctetString `aper:"sizeLB:3,sizeUB:3"`
}
