package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type TAIListForPagingItem struct {
	TAI          TAI                                                   `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerTAIListForPagingItemExtIEs `aper:"optional"`
}
