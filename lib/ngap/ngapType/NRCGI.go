package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type NRCGI struct {
	PLMNIdentity   PLMNIdentity
	NRCellIdentity NRCellIdentity
	IEExtensions   *ProtocolExtensionContainerNRCGIExtIEs `aper:"optional"`
}
