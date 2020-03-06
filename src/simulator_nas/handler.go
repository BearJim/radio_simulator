package simulator_nas

import (
	"fmt"
	"radio_simulator/lib/nas/nasMessage"
	"radio_simulator/src/simulator_context"
	"radio_simulator/src/simulator_nas/nas_packet"
	"radio_simulator/src/simulator_ngap"
)

func HandleAuthenticationRequest(ue *simulator_context.UeContext, request *nasMessage.AuthenticationRequest) error {

	nasLog.Infof("Ue[%s] Handle Authentication Request", ue.Supi)

	if request == nil {
		return fmt.Errorf("AuthenticationRequest body is nil")
	}
	rand := request.GetRANDValue()
	resStat := ue.DeriveRESstarAndSetKey(rand[:])
	nasPdu := nas_packet.GetAuthenticationResponse(resStat, "")
	simulator_ngap.SendUplinkNasTransport(ue.Ran, ue, nasPdu)
	return nil
}

func HandleSecurityModeCommand(ue *simulator_context.UeContext, request *nasMessage.SecurityModeCommand) error {
	ue.EncAlg = request.SelectedNASSecurityAlgorithms.GetTypeOfCipheringAlgorithm()
	ue.IntAlg = request.SelectedNASSecurityAlgorithms.GetTypeOfIntegrityProtectionAlgorithm()
	nasContent, err := nas_packet.GetRegistrationRequestWith5GMM(ue, nasMessage.RegistrationType5GSInitialRegistration, nil, nil)
	if err != nil {
		return err
	}
	nasPdu, err := nas_packet.GetSecurityModeComplete(ue, nasContent)
	if err != nil {
		return err
	}
	simulator_ngap.SendUplinkNasTransport(ue.Ran, ue, nasPdu)
	return nil
}
