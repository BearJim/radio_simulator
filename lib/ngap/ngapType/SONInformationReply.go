package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type SONInformationReply struct {
	XnTNLConfigurationInfo *XnTNLConfigurationInfo                              `aper:"valueExt,optional"`
	IEExtensions           *ProtocolExtensionContainerSONInformationReplyExtIEs `aper:"optional"`
}
