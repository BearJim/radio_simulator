package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type PDUSessionResourceSetupItemSUReq struct {
	PDUSessionID                           PDUSessionID
	PDUSessionNASPDU                       *NASPDU `aper:"optional"`
	SNSSAI                                 SNSSAI  `aper:"valueExt"`
	PDUSessionResourceSetupRequestTransfer aper.OctetString
	IEExtensions                           *ProtocolExtensionContainerPDUSessionResourceSetupItemSUReqExtIEs `aper:"optional"`
}
