package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type ForbiddenAreaInformationItem struct {
	PLMNIdentity  PLMNIdentity
	ForbiddenTACs ForbiddenTACs
	IEExtensions  *ProtocolExtensionContainerForbiddenAreaInformationItemExtIEs `aper:"optional"`
}
