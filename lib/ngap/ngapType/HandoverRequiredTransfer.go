package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type HandoverRequiredTransfer struct {
	DirectForwardingPathAvailability *DirectForwardingPathAvailability                         `aper:"optional"`
	IEExtensions                     *ProtocolExtensionContainerHandoverRequiredTransferExtIEs `aper:"optional"`
}
