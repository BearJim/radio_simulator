package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type PLMNSupportItem struct {
	PLMNIdentity     PLMNIdentity
	SliceSupportList SliceSupportList
	IEExtensions     *ProtocolExtensionContainerPLMNSupportItemExtIEs `aper:"optional"`
}
