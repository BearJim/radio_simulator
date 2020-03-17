package ngap

import (
	"radio_simulator/lib/aper"
)

// NgENBIDPresent CHOICE value
const (
	NgENBIDPresentNothing int = iota
	NgENBIDPresentMacroNgENBID
	NgENBIDPresentShortNgENBID
	NgENBIDPresentLongNgENBID
	NgENBIDPresentChoiceExtensions
)

// NgENBID CHOICE Type
type NgENBID struct {
	Present          int
	MacroNgENBID     *aper.BitString `aper:"sizeLB:20,sizeUB:20"`
	ShortNgENBID     *aper.BitString `aper:"sizeLB:18,sizeUB:18"`
	LongNgENBID      *aper.BitString `aper:"sizeLB:21,sizeUB:21"`
	ChoiceExtensions *ProtocolIESingleContainerNgENBID
}
