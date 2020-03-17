package ngap

// ProtocolIEContainerNGSetupRequestIEs Type
type ProtocolIEContainerNGSetupRequestIEs struct {
	List []NGSetupRequestIEs `aper:"sizeLB:0,sizeUB:65535"`
}
