package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type PathSwitchRequestSetupFailedTransfer struct {
	Cause        Cause                                                                 `aper:"valueLB:0,valueUB:5"`
	IEExtensions *ProtocolExtensionContainerPathSwitchRequestSetupFailedTransferExtIEs `aper:"optional"`
}
