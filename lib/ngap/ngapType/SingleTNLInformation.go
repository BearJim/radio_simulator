package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type SingleTNLInformation struct {
	UPTransportLayerInformation UPTransportLayerInformation                           `aper:"valueLB:0,valueUB:1"`
	IEExtensions                *ProtocolExtensionContainerSingleTNLInformationExtIEs `aper:"optional"`
}
