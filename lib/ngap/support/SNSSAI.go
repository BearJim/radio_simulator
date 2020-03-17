package ngap

// SNSSAI Type
type SNSSAI struct {
	SST          SST
	SD           *SD                                     `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerSNSSAIExtIEs `aper:"optional"`
}
