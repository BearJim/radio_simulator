package simulator_nas

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jay16213/radio_simulator/pkg/api"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
	"github.com/jay16213/radio_simulator/pkg/simulator_nas/nas_packet"

	"github.com/free5gc/milenage"
	"github.com/free5gc/nas/nasMessage"
)

func (c *NASController) handleAuthenticationRequest(ue *simulator_context.UeContext, request *nasMessage.AuthenticationRequest) error {
	nasLog.Infow("Handle Authentication Request", "supi", ue.Supi, "id", ue.AmfUeNgapId)

	ue.NgKsi = request.GetNasKeySetIdentifiler()

	// Get RAND & AUTN from Authentication request
	RAND := request.GetRANDValue()
	AUTN := request.GetAUTN()
	SQNxorAK := AUTN[0:6]
	AMF := AUTN[6:8]
	MAC := AUTN[8:]

	authData := ue.AuthData
	servingNetworkName := ue.GetServingNetworkName()
	SQNms, _ := hex.DecodeString(authData.SQN)

	// Run milenage
	XMAC, MAC_S := make([]byte, 8), make([]byte, 8)
	CK, IK := make([]byte, 16), make([]byte, 16)
	RES := make([]byte, 8)
	SQN := make([]byte, 6)
	AK, AKstar := make([]byte, 6), make([]byte, 6)
	OPC, _ := hex.DecodeString(authData.Opc)
	K, _ := hex.DecodeString(authData.K)

	// Generate RES, CK, IK, AK
	if err := milenage.F2345(OPC, K, RAND[:], RES, CK, IK, AK, nil); err != nil {
		logger.NASLog.Error(err)
		return nil
	}

	// Derive SQN
	for i := 0; i < 6; i++ {
		SQN[i] = SQNxorAK[i] ^ AK[i]
	}

	// Generate XMAC
	if err := milenage.F1(OPC, K, RAND[:], SQN, AMF, XMAC, nil); err != nil {
		logger.NASLog.Error(err)
		return nil
	}

	// Verify MAC == XMAC
	if !bytes.Equal(MAC, XMAC) {
		logger.NASLog.Errorf("Authentication Failed: MAC (0x%0x) != XMAC (0x%0x)", MAC, XMAC)
		c.SendAuthenticationFailure(ue, nasMessage.Cause5GMMMACFailure, nil)
		return nil
	}

	// Verify that SQN is in the current range TS 33.102
	// sqn is out of sync -> synchronisation failure -> trigger resync procedure
	if !bytes.Equal(SQN, SQNms) {
		logger.NASLog.Errorf("Authentication Synchronisation Failure: SQN (0x%0x) != SQNms (0x%0x)", SQN, SQNms)
		SQNmsXorAK := make([]byte, 6)

		// TS 33.102 6.3.3: The AMF used to calculate MAC S assumes a dummy value of all zeros so that it does not
		// need to be transmitted in the clear in the re-synch message.
		if err := milenage.F1(OPC, K, RAND[:], SQNms, []byte{0x00, 0x00}, nil, MAC_S); err != nil {
			logger.NASLog.Error(err)
			return nil
		}
		if err := milenage.F2345(OPC, K, RAND[:], nil, nil, nil, nil, AKstar); err != nil {
			logger.NASLog.Error(err)
			return nil
		}
		for i := 0; i < 6; i++ {
			SQNmsXorAK[i] = SQNms[i] ^ AKstar[i]
		}
		AUTS := append(SQNmsXorAK, MAC_S...)
		c.SendAuthenticationFailure(ue, nasMessage.Cause5GMMSynchFailure, AUTS)
		return nil
	}

	// derive RES* and send response
	resStar := ue.DeriveRESstar(CK, IK, servingNetworkName, RAND[:], RES)
	c.SendAuthenticationResponse(ue, resStar)

	// generate keys
	kausf := simulator_context.DerivateKausf(CK, IK, servingNetworkName, SQNxorAK)
	logger.NASLog.Debugf("Kausf: 0x%0x", kausf)
	kseaf := simulator_context.DerivateKseaf(kausf, servingNetworkName)
	logger.NASLog.Debugf("Kseaf: 0x%0x", kseaf)
	ue.DerivateKamf(kseaf, []byte{0x00, 0x00})
	ue.AuthDataSQNAddOne()
	return nil
}

func (c *NASController) handleAuthenticationReject(ue *simulator_context.UeContext, message *nasMessage.AuthenticationReject) error {
	logger.NASLog.Errorw("Receive Authentication Reject", "supi", ue.Supi, "id", ue.AmfUeNgapId)
	return nil
}

func (c *NASController) handleRegistrationReject(ue *simulator_context.UeContext, message *nasMessage.RegistrationReject) error {
	logger.NASLog.Warnw("Handle Registration Reject", "supi", ue.Supi, "id", ue.AmfUeNgapId)
	if message.Cause5GMM.GetCauseValue() == nasMessage.Cause5GMMCongestion {
		logger.NASLog.Warnw("Restart Initial Registration", "supi", ue.Supi, "id", ue.AmfUeNgapId)
		ue.RestartCount++
		ue.RestartTimeStamp = time.Now()
		c.ngMessager.SendInitailUeMessage_RegistraionRequest(ue.AMFEndpoint, ue)
	}
	return nil
}

func (c *NASController) handleSecurityModeCommand(ue *simulator_context.UeContext, request *nasMessage.SecurityModeCommand) error {

	nasLog.Infof("UE[%s] Handle Security Mode Command", ue.Supi)

	nasContent, err := nas_packet.GetRegistrationRequestWith5GMM(ue, nasMessage.RegistrationType5GSInitialRegistration, nil, nil)
	if err != nil {
		return err
	}
	c.SendSecurityModeCommand(ue, nasContent)
	return nil
}

func (c *NASController) handleRegistrationAccept(ue *simulator_context.UeContext, request *nasMessage.RegistrationAccept) error {

	nasLog.Infof("UE[%s] Handle Registration Accept", ue.Supi)

	ue.Guti = request.GUTI5G

	nasPdu, err := nas_packet.GetRegistrationComplete(ue, nil)
	if err != nil {
		return err
	}
	nasLog.Info("Send Registration Complete")
	c.ngMessager.SendUplinkNASTransport(ue.AMFEndpoint, ue, nasPdu)
	ue.RmState = simulator_context.RmStateRegistered
	ue.SendAPINotification(api.StatusCode_OK, simulator_context.MsgRegisterSuccess)
	ue.RestartCount = 0
	return nil
}

func (c *NASController) handleDeregistrationAccept(ue *simulator_context.UeContext, request *nasMessage.DeregistrationAcceptUEOriginatingDeregistration) error {

	nasLog.Infof("UE[%s] Handle Deregistration Accept", ue.Supi)

	ue.RmState = simulator_context.RmStateDeregitered
	ue.SendAPINotification(api.StatusCode_OK, simulator_context.MsgDeregisterSuccess)
	return nil
}

func (c *NASController) handleDLNASTransport(ue *simulator_context.UeContext, request *nasMessage.DLNASTransport) error {

	nasLog.Infof("UE[%s] Handle DL NAS Transport", ue.Supi)

	switch request.GetPayloadContainerType() {
	case nasMessage.PayloadContainerTypeN1SMInfo:
		c.handleGsmMessage(ue, request.GetPayloadContainerContents())
	case nasMessage.PayloadContainerTypeSMS:
		return fmt.Errorf("PayloadContainerTypeSMS has not been implemented yet in DL NAS TRANSPORT")
	case nasMessage.PayloadContainerTypeLPP:
		return fmt.Errorf("PayloadContainerTypeLPP has not been implemented yet in DL NAS TRANSPORT")
	case nasMessage.PayloadContainerTypeSOR:
		return fmt.Errorf("PayloadContainerTypeSOR has not been implemented yet in DL NAS TRANSPORT")
	case nasMessage.PayloadContainerTypeUEPolicy:
		return fmt.Errorf("PayloadContainerTypeUEPolicy has not been implemented yet in DL NAS TRANSPORT")
	case nasMessage.PayloadContainerTypeUEParameterUpdate:
		return fmt.Errorf("PayloadContainerTypeUEParameterUpdate has not been implemented yet in DL NAS TRANSPORT")
	case nasMessage.PayloadContainerTypeMultiplePayload:
		return fmt.Errorf("PayloadContainerTypeMultiplePayload has not been implemented yet in DL NAS TRANSPORT")
	}
	return nil
}
