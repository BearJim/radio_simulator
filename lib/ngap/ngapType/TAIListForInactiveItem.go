package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type TAIListForInactiveItem struct {
	TAI          TAI                                                     `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerTAIListForInactiveItemExtIEs `aper:"optional"`
}
