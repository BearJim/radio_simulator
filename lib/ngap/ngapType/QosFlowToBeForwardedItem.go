package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type QosFlowToBeForwardedItem struct {
	QosFlowIdentifier QosFlowIdentifier
	IEExtensions      *ProtocolExtensionContainerQosFlowToBeForwardedItemExtIEs `aper:"optional"`
}
