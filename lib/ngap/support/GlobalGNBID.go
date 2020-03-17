package ngap

// GlobalGNBID Type have extension
type GlobalGNBID struct {
	PLMNIdentity PLMNIdentity
	GNBID        GNBID                                        `aper:"valueUB:1"`
	IEExtensions *ProtocolExtensionContainerGlobalGNBIDExtIEs `aper:"optional"`
}
