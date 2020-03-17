package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type TNLInformationItem struct {
	QosFlowPerTNLInformation QosFlowPerTNLInformation                            `aper:"valueExt"`
	IEExtensions             *ProtocolExtensionContainerTNLInformationItemExtIEs `aper:"optional"`
}
