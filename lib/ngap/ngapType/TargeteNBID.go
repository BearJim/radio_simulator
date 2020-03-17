package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type TargeteNBID struct {
	GlobalENBID    GlobalNgENBID                                `aper:"valueExt"`
	SelectedEPSTAI EPSTAI                                       `aper:"valueExt"`
	IEExtensions   *ProtocolExtensionContainerTargeteNBIDExtIEs `aper:"optional"`
}
