package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type ERABInformationItem struct {
	ERABID       ERABID
	DLForwarding *DLForwarding                                        `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerERABInformationItemExtIEs `aper:"optional"`
}
