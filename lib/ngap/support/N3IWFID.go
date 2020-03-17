package ngap

import (
	"radio_simulator/lib/aper"
)

// N3IWFIDPresent CHOICE value
const (
	N3IWFIDPresentNothing int = iota
	N3IWFIDPresentGNBID
	N3IWFIDPresentChoiceExtensions
)

// N3IWFID CHOICE Type
type N3IWFID struct {
	Present          int
	N3IWFID          *aper.BitString `aper:"sizeLB:16,sizeUB:16"`
	ChoiceExtensions *ProtocolIESingleContainerN3IWFID
}
