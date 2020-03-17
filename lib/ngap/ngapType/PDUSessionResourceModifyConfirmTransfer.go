package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type PDUSessionResourceModifyConfirmTransfer struct {
	QosFlowModifyConfirmList  QosFlowModifyConfirmList
	TNLMappingList            *TNLMappingList                                                          `aper:"optional"`
	QosFlowFailedToModifyList *QosFlowList                                                             `aper:"optional"`
	IEExtensions              *ProtocolExtensionContainerPDUSessionResourceModifyConfirmTransferExtIEs `aper:"optional"`
}
