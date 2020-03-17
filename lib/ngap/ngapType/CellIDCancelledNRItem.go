package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type CellIDCancelledNRItem struct {
	NRCGI              NRCGI `aper:"valueExt"`
	NumberOfBroadcasts NumberOfBroadcasts
	IEExtensions       *ProtocolExtensionContainerCellIDCancelledNRItemExtIEs `aper:"optional"`
}
