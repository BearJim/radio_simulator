package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type TargetRANNodeID struct {
	GlobalRANNodeID GlobalRANNodeID                                  `aper:"valueLB:0,valueUB:3"`
	SelectedTAI     TAI                                              `aper:"valueExt"`
	IEExtensions    *ProtocolExtensionContainerTargetRANNodeIDExtIEs `aper:"optional"`
}