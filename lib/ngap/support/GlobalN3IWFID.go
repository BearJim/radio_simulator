package ngap

// GlobalN3IWFID Type have extension
type GlobalN3IWFID struct {
	PLMNIdentity PLMNIdentity
	N3IWFID      N3IWFID                                        `aper:"valueUB:1"`
	IEExtensions *ProtocolExtensionContainerGlobalN3IWFIDExtIEs `aper:"optional"`
}
