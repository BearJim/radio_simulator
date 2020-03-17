package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type UEAssociatedLogicalNGConnectionItem struct {
	AMFUENGAPID  *AMFUENGAPID                                                         `aper:"optional"`
	RANUENGAPID  *RANUENGAPID                                                         `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerUEAssociatedLogicalNGConnectionItemExtIEs `aper:"optional"`
}
