package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type EPSTAI struct {
	PLMNIdentity PLMNIdentity
	EPSTAC       EPSTAC
	IEExtensions *ProtocolExtensionContainerEPSTAIExtIEs `aper:"optional"`
}
