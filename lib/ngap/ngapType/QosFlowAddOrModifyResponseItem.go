package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type QosFlowAddOrModifyResponseItem struct {
	QosFlowIdentifier QosFlowIdentifier
	IEExtensions      *ProtocolExtensionContainerQosFlowAddOrModifyResponseItemExtIEs `aper:"optional"`
}
