package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type CompletedCellsInEAINRItem struct {
	NRCGI        NRCGI                                                      `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerCompletedCellsInEAINRItemExtIEs `aper:"optional"`
}
