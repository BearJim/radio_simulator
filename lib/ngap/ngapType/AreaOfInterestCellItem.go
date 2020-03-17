package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type AreaOfInterestCellItem struct {
	NGRANCGI     NGRANCGI                                                `aper:"valueLB:0,valueUB:2"`
	IEExtensions *ProtocolExtensionContainerAreaOfInterestCellItemExtIEs `aper:"optional"`
}
