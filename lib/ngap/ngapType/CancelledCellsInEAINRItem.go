package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type CancelledCellsInEAINRItem struct {
	NRCGI              NRCGI `aper:"valueExt"`
	NumberOfBroadcasts NumberOfBroadcasts
	IEExtensions       *ProtocolExtensionContainerCancelledCellsInEAINRItemExtIEs `aper:"optional"`
}
