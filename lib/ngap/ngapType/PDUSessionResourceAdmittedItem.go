package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type PDUSessionResourceAdmittedItem struct {
	PDUSessionID                       PDUSessionID
	HandoverRequestAcknowledgeTransfer aper.OctetString
	IEExtensions                       *ProtocolExtensionContainerPDUSessionResourceAdmittedItemExtIEs `aper:"optional"`
}
