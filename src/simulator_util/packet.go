package simulator_util

import (

	// Nausf_UEAU_Client "radio_simulator/lib/Nausf_UEAuthentication"
	"radio_simulator/lib/ngap"
	"radio_simulator/src/simulator_ngap"
	// "radio_simulator/lib/openapi/models"
)

// func GetNGSetupRequest(gnbId []byte, bitlength uint64, name string) ([]byte, error) {
// 	message, _ := simulator_ngap.BuildNGSetupRequest()
// 	pdu, _ := ngap.Decoder(message)
// 	// GlobalRANNodeID
// 	ie := pdu.InitiatingMessage.Value.NGSetupRequest.ProtocolIEs.List[0]
// 	gnbID := ie.Value.GlobalRANNodeID.GlobalGNBID.GNBID.GNBID
// 	gnbID.Bytes = gnbId
// 	gnbID.BitLength = bitlength
// 	// RANNodeName
// 	ie = pdu.InitiatingMessage.Value.NGSetupRequest.ProtocolIEs.List[1]
// 	ie.Value.RANNodeName.Value = name

// 	return ngap.Encoder(pdu)
// }

// func GetInitialUEMessage(ranUeNgapID int64, nasPdu []byte, fiveGSTmsi string) ([]byte, error) {
// 	message := simulator_ngap.BuildInitialUEMessage(ranUeNgapID, nasPdu, fiveGSTmsi)
// 	return ngap.Encoder(message)
// }

// func GetUplinkNASTransport(amfUeNgapID, ranUeNgapID int64, nasPdu []byte) ([]byte, error) {
// 	message := simulator_ngap.BuildUplinkNasTransport(amfUeNgapID, ranUeNgapID, nasPdu)
// 	return ngap.Encoder(message)
// }

func GetInitialContextSetupResponse(amfUeNgapID int64, ranUeNgapID int64) ([]byte, error) {
	message := simulator_ngap.BuildInitialContextSetupResponseForRegistraionTest(amfUeNgapID, ranUeNgapID)

	return ngap.Encoder(message)
}

// func GetInitialContextSetupResponseForServiceRequest(amfUeNgapID int64, ranUeNgapID int64, ipv4 string) ([]byte, error) {
// 	message := simulator_ngap.BuildInitialContextSetupResponse(amfUeNgapID, ranUeNgapID, ipv4, nil)
// 	return ngap.Encoder(message)
// }

func GetPDUSessionResourceSetupResponse(amfUeNgapID int64, ranUeNgapID int64, ipv4 string) ([]byte, error) {
	message := simulator_ngap.BuildPDUSessionResourceSetupResponseForRegistrationTest(amfUeNgapID, ranUeNgapID, ipv4)
	return ngap.Encoder(message)
}

// func EncodeNasPduWithSecurity(ue *simulator_context.UeContext, pdu []byte) ([]byte, error) {
// 	m := nas.NewMessage()
// 	err := m.PlainNasDecode(&pdu)
// 	if err != nil {
// 		return nil, err
// 	}
// 	m.SecurityHeader = nas.SecurityHeader{
// 		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
// 		SecurityHeaderType:    nas.SecurityHeaderTypeIntegrityProtectedAndCiphered,
// 	}
// 	return NASEncode(ue, m)
// }

// func GetUEContextReleaseComplete(amfUeNgapID int64, ranUeNgapID int64, pduSessionIDList []int64) ([]byte, error) {
// 	message := simulator_ngap.BuildUEContextReleaseComplete(amfUeNgapID, ranUeNgapID, pduSessionIDList)
// 	return ngap.Encoder(message)
// }

func GetUEContextReleaseRequest(amfUeNgapID int64, ranUeNgapID int64, pduSessionIDList []int64) ([]byte, error) {
	message := simulator_ngap.BuildUEContextReleaseRequest(amfUeNgapID, ranUeNgapID, pduSessionIDList)
	return ngap.Encoder(message)
}

func GetPDUSessionResourceReleaseResponse(amfUeNgapID int64, ranUeNgapID int64) ([]byte, error) {
	message := simulator_ngap.BuildPDUSessionResourceReleaseResponseForReleaseTest(amfUeNgapID, ranUeNgapID)
	return ngap.Encoder(message)
}
func GetPathSwitchRequest(amfUeNgapID int64, ranUeNgapID int64) ([]byte, error) {
	message := simulator_ngap.BuildPathSwitchRequest(amfUeNgapID, ranUeNgapID)
	message.InitiatingMessage.Value.PathSwitchRequest.ProtocolIEs.List = message.InitiatingMessage.Value.PathSwitchRequest.ProtocolIEs.List[0:5]
	return ngap.Encoder(message)
}

func GetHandoverRequired(amfUeNgapID int64, ranUeNgapID int64, targetGNBID []byte, targetCellID []byte) ([]byte, error) {
	message := simulator_ngap.BuildHandoverRequired(amfUeNgapID, ranUeNgapID, targetGNBID, targetCellID)
	return ngap.Encoder(message)
}

func GetHandoverRequestAcknowledge(amfUeNgapID int64, ranUeNgapID int64) ([]byte, error) {
	message := simulator_ngap.BuildHandoverRequestAcknowledge(amfUeNgapID, ranUeNgapID)
	return ngap.Encoder(message)
}

func GetHandoverNotify(amfUeNgapID int64, ranUeNgapID int64) ([]byte, error) {
	message := simulator_ngap.BuildHandoverNotify(amfUeNgapID, ranUeNgapID)
	return ngap.Encoder(message)
}

func GetPDUSessionResourceSetupResponseForPaging(amfUeNgapID int64, ranUeNgapID int64, ipv4 string) ([]byte, error) {
	message := simulator_ngap.BuildPDUSessionResourceSetupResponseForPaging(amfUeNgapID, ranUeNgapID, ipv4)
	return ngap.Encoder(message)
}
