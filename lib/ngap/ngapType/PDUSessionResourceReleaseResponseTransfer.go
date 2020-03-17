package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type PDUSessionResourceReleaseResponseTransfer struct {
	IEExtensions *ProtocolExtensionContainerPDUSessionResourceReleaseResponseTransferExtIEs `aper:"optional"`
}
