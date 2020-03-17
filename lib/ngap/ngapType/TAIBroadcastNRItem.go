package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type TAIBroadcastNRItem struct {
	TAI                   TAI `aper:"valueExt"`
	CompletedCellsInTAINR CompletedCellsInTAINR
	IEExtensions          *ProtocolExtensionContainerTAIBroadcastNRItemExtIEs `aper:"optional"`
}
