package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type CellIDBroadcastNRItem struct {
	NRCGI        NRCGI                                                  `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerCellIDBroadcastNRItemExtIEs `aper:"optional"`
}
