package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type RATRestrictionsItem struct {
	PLMNIdentity              PLMNIdentity
	RATRestrictionInformation RATRestrictionInformation
	IEExtensions              *ProtocolExtensionContainerRATRestrictionsItemExtIEs `aper:"optional"`
}
