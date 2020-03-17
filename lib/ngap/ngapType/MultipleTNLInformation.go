package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type MultipleTNLInformation struct {
	TNLInformationList TNLInformationList
	IEExtensions       *ProtocolExtensionContainerMultipleTNLInformationExtIEs `aper:"optional"`
}
