package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type DRBStatusUL12 struct {
	ULCOUNTValue              COUNTValueForPDCPSN12                          `aper:"valueExt"`
	ReceiveStatusOfULPDCPSDUs *aper.BitString                                `aper:"sizeLB:1,sizeUB:2048,optional"`
	IEExtension               *ProtocolExtensionContainerDRBStatusUL12ExtIEs `aper:"optional"`
}
