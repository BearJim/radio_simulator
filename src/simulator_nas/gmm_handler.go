package simulator_nas

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"radio_simulator/lib/nas/nasMessage"
	"radio_simulator/src/simulator_context"
	"radio_simulator/src/simulator_nas/nas_packet"
	"radio_simulator/src/simulator_ngap"
)

func HandleAuthenticationRequest(ue *simulator_context.UeContext, request *nasMessage.AuthenticationRequest) error {

	nasLog.Infof("UE[%s] Handle Authentication Request", ue.Supi)

	if request == nil {
		return fmt.Errorf("AuthenticationRequest body is nil")
	}
	ue.NgKsi = request.GetNasKeySetIdentifiler()
	rand := request.GetRANDValue()
	resStat := ue.DeriveRESstarAndSetKey(rand[:])
	nasPdu := nas_packet.GetAuthenticationResponse(resStat, "")
	simulator_ngap.SendUplinkNasTransport(ue.Ran, ue, nasPdu)
	return nil
}

func HandleSecurityModeCommand(ue *simulator_context.UeContext, request *nasMessage.SecurityModeCommand) error {

	nasLog.Infof("UE[%s] Handle Security Mode Command", ue.Supi)

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
func HandleRegistrationAccept(ue *simulator_context.UeContext, request *nasMessage.RegistrationAccept) error {

	nasLog.Infof("UE[%s] Handle Registration Accept", ue.Supi)

	ue.Guti = request.GUTI5G

	nasPdu, err := nas_packet.GetRegistrationComplete(ue, nil)
	if err != nil {
		return err
	}
	simulator_ngap.SendUplinkNasTransport(ue.Ran, ue, nasPdu)
	ue.RegisterState = simulator_context.RegisterStateRegistered

	fmt.Println(ue.AuthData.SQN)
	input, _ := ioutil.ReadFile("config/uecfg.conf")
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, "SQN") {
			//noSpace := strings.Replace(lines[i], " ", "", -1)
			//num, _ := strconv.ParseInt(strings.Split(noSpace, ":")[1], 16, 64)
			num, _ := strconv.ParseInt(ue.AuthData.SQN, 16, 64)
			ue.AuthData.SQN = fmt.Sprintf("%x", num+1)
			lines[i] = fmt.Sprintf("  SQN: %s", ue.AuthData.SQN)
		}
	}
	output := strings.Join(lines, "\n")
	ioutil.WriteFile("config/uecfg.conf", []byte(output), 0644)

	ue.SendMsg("[REG] SUCCESS\n")
	return nil
}
func HandleDeregistrationAccept(ue *simulator_context.UeContext, request *nasMessage.DeregistrationAcceptUEOriginatingDeregistration) error {

	nasLog.Infof("UE[%s] Handle Deregistration Accept", ue.Supi)

	ue.RegisterState = simulator_context.RegisterStateDeregitered
	return nil
}

func HandleDLNASTransport(ue *simulator_context.UeContext, request *nasMessage.DLNASTransport) error {

	nasLog.Infof("UE[%s] Handle DL NAS Transport", ue.Supi)

	switch request.GetPayloadContainerType() {
	case nasMessage.PayloadContainerTypeN1SMInfo:
		HandleNAS(ue, request.GetPayloadContainerContents())
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
