package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type CellIDCancelledEUTRAItem struct {
	EUTRACGI           EUTRACGI `aper:"valueExt"`
	NumberOfBroadcasts NumberOfBroadcasts
	IEExtensions       *ProtocolExtensionContainerCellIDCancelledEUTRAItemExtIEs `aper:"optional"`
}
