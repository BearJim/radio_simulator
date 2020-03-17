package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type RANStatusTransferTransparentContainer struct {
	DRBsSubjectToStatusTransferList DRBsSubjectToStatusTransferList
	IEExtensions                    *ProtocolExtensionContainerRANStatusTransferTransparentContainerExtIEs `aper:"optional"`
}
