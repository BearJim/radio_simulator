package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type QosFlowAcceptedItem struct {
	QosFlowIdentifier QosFlowIdentifier
	IEExtensions      *ProtocolExtensionContainerQosFlowAcceptedItemExtIEs `aper:"optional"`
}
