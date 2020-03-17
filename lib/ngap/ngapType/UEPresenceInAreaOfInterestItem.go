package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type UEPresenceInAreaOfInterestItem struct {
	LocationReportingReferenceID LocationReportingReferenceID
	UEPresence                   UEPresence
	IEExtensions                 *ProtocolExtensionContainerUEPresenceInAreaOfInterestItemExtIEs `aper:"optional"`
}
