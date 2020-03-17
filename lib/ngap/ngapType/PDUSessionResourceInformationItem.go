package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type PDUSessionResourceInformationItem struct {
	PDUSessionID              PDUSessionID
	QosFlowInformationList    QosFlowInformationList
	DRBsToQosFlowsMappingList *DRBsToQosFlowsMappingList                                         `aper:"optional"`
	IEExtensions              *ProtocolExtensionContainerPDUSessionResourceInformationItemExtIEs `aper:"optional"`
}
