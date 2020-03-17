package ngap

// ProtocolExtensionContainerGlobalGNBIDExtIEs Type
type ProtocolExtensionContainerGlobalGNBIDExtIEs struct {
	List []GlobalGNBIDExtIEs `aper:"sizeLB:0,sizeUB:65535"`
}

// ProtocolExtensionContainerGlobalNgENBIDExtIEs Type
type ProtocolExtensionContainerGlobalNgENBIDExtIEs struct {
	List []GlobalNgENBIDExtIEs `aper:"sizeLB:0,sizeUB:65535"`
}

// ProtocolExtensionContainerGlobalN3IWFIDExtIEs Type
type ProtocolExtensionContainerGlobalN3IWFIDExtIEs struct {
	List []GlobalN3IWFIDExtIEs `aper:"sizeLB:0,sizeUB:65535"`
}

// ProtocolExtensionContainerSupportedTAItemExtIEs Type
type ProtocolExtensionContainerSupportedTAItemExtIEs struct {
	List []SupportedTAItemExtIEs `aper:"sizeLB:0,sizeUB:65535"`
}

// ProtocolExtensionContainerBroadcastPLMNItemExtIEs Type
type ProtocolExtensionContainerBroadcastPLMNItemExtIEs struct {
	List []BroadcastPLMNItemExtIEs `aper:"sizeLB:0,sizeUB:65535"`
}

// ProtocolExtensionContainerSliceSupportItemExtIEs Type
type ProtocolExtensionContainerSliceSupportItemExtIEs struct {
	List []SupportItemExtIEs `aper:"sizeLB:0,sizeUB:65535"`
}

// ProtocolExtensionContainerSNSSAIExtIEs Type
type ProtocolExtensionContainerSNSSAIExtIEs struct {
	List []SNSSAIItemExtIEs `aper:"sizeLB:0,sizeUB:65535"`
}
