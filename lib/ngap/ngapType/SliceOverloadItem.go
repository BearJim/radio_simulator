package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type SliceOverloadItem struct {
	SNSSAI       SNSSAI                                             `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerSliceOverloadItemExtIEs `aper:"optional"`
}
