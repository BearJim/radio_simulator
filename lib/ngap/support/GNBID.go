package ngap

import (
	"radio_simulator/lib/aper"
)

// GNBIDPresent CHOICE value
const (
	GNBIDPresentNothing int = iota
	GNBIDPresentGNBID
	GNBIDPresentChoiceExtensions
)

// GNBID CHOICE Type
type GNBID struct {
	Present          int
	GNBID            *aper.BitString `aper:"sizeLB:22,sizeUB:32"`
	ChoiceExtensions *ProtocolIESingleContainerGNBID
}
