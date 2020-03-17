package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type EUTRACGI struct {
	PLMNIdentity      PLMNIdentity
	EUTRACellIdentity EUTRACellIdentity
	IEExtensions      *ProtocolExtensionContainerEUTRACGIExtIEs `aper:"optional"`
}
