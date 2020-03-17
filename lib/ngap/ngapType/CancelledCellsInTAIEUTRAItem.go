package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type CancelledCellsInTAIEUTRAItem struct {
	EUTRACGI           EUTRACGI `aper:"valueExt"`
	NumberOfBroadcasts NumberOfBroadcasts
	IEExtensions       *ProtocolExtensionContainerCancelledCellsInTAIEUTRAItemExtIEs `aper:"optional"`
}
