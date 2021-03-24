package simulator_nas

import (
	"fmt"
	"strconv"

	"github.com/jay16213/radio_simulator/pkg/api"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
	"github.com/jay16213/radio_simulator/pkg/simulator_nas/nas_packet"

	"github.com/free5gc/nas/nasMessage"
)

func (c *NASController) handleAuthenticationRequest(ue *simulator_context.UeContext, request *nasMessage.AuthenticationRequest) error {

	nasLog.Infof("UE[%s] Handle Authentication Request", ue.Supi)

	if request == nil {
		return fmt.Errorf("AuthenticationRequest body is nil")
	}
	ue.NgKsi = request.GetNasKeySetIdentifiler()
	rand := request.GetRANDValue()
	resStar := ue.DeriveRESstarAndSetKey(rand[:])
	c.SendAuthenticationResponse(ue, resStar)
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
	ue.RmState = simulator_context.RegisterStateRegistered
	num, _ := strconv.ParseInt(ue.AuthData.SQN, 16, 64)
	ue.AuthData.SQN = fmt.Sprintf("%x", num+1)
	ue.SendAPINotification(api.StatusCode_OK, simulator_context.MsgRegisterSuccess)
	return nil
}
func (c *NASController) handleDeregistrationAccept(ue *simulator_context.UeContext, request *nasMessage.DeregistrationAcceptUEOriginatingDeregistration) error {

	nasLog.Infof("UE[%s] Handle Deregistration Accept", ue.Supi)

	ue.RmState = simulator_context.RegisterStateDeregitered
	return nil
}

func (c *NASController) handleDLNASTransport(ue *simulator_context.UeContext, request *nasMessage.DLNASTransport) error {

	nasLog.Infof("UE[%s] Handle DL NAS Transport", ue.Supi)

	switch request.GetPayloadContainerType() {
	case nasMessage.PayloadContainerTypeN1SMInfo:
		c.HandleNAS(ue, request.GetPayloadContainerContents())
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
