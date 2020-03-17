package ngap

// GlobalNgENBID Type have extension
type GlobalNgENBID struct {
	PLMNIdentity PLMNIdentity
	NgENBID      NgENBID                                        `aper:"valueUB:3"`
	IEExtensions *ProtocolExtensionContainerGlobalNgENBIDExtIEs `aper:"optional"`
}
