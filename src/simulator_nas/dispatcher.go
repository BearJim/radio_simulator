package simulator_nas

import (
	"github.com/sirupsen/logrus"
	"radio_simulator/lib/nas"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_context"
	"radio_simulator/src/simulator_nas/nas_security"
)

var nasLog *logrus.Entry

func init() {
	nasLog = logger.NasLog
}

func checkMsgError(err error, msg string) {
	if err != nil {
		nasLog.Errorf("Handle %s Error: %s", msg, err.Error())
	}
}

func HandleNAS(ue *simulator_context.UeContext, nasPdu []byte) {

	if ue == nil {
		nasLog.Error("Ue is nil")
		return
	}

	if nasPdu == nil {
		nasLog.Error("nasPdu is nil")
		return
	}

	var msg *nas.Message

	if ue.RegisterState == simulator_context.RegisterStateRegitered {
		var err error
		msg, err = nas_security.NASDecode(ue, nas.GetSecurityHeaderType(nasPdu)&0x0f, nasPdu)
		if err != nil {
			nasLog.Error(err.Error())
			return
		}
	} else {
		msg = new(nas.Message)
		err := msg.PlainNasDecode(&nasPdu)
		if err != nil {
			nasLog.Error(err.Error())
			return
		}
	}

	if msg.GmmMessage != nil {
		switch msg.GmmMessage.GetMessageType() {
		case nas.MsgTypeAuthenticationRequest:
			checkMsgError(HandleAuthenticationRequest(ue, msg.GmmMessage.AuthenticationRequest), "AuthenticationRequest")
		case nas.MsgTypeSecurityModeCommand:
			checkMsgError(HandleSecurityModeCommand(ue, msg.GmmMessage.SecurityModeCommand), "SecurityModeCommand")
		case nas.MsgTypeRegistrationAccept:
			checkMsgError(HandleRegistrationAccept(ue, msg.GmmMessage.RegistrationAccept), "RegistrationAccept")
		// case nas.MsgTypeULNASTransport:
		// 	return gmm_handler.HandleULNASTransport(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.ULNASTransport)
		// case nas.MsgTypeRegistrationRequest:
		// 	if err := gmm_handler.HandleRegistrationRequest(amfUe, models.AccessType__3_GPP_ACCESS, procedureCode, gmmMessage.RegistrationRequest); err != nil {
		// 		return err
		// 	}
		// case nas.MsgTypeIdentityResponse:
		// 	if err := gmm_handler.HandleIdentityResponse(amfUe, gmmMessage.IdentityResponse); err != nil {
		// 		return err
		// 	}
		// case nas.MsgTypeConfigurationUpdateComplete:
		// 	if err := gmm_handler.HandleConfigurationUpdateComplete(amfUe, gmmMessage.ConfigurationUpdateComplete); err != nil {
		// 		return err
		// 	}
		// case nas.MsgTypeServiceRequest:
		// 	if err := gmm_handler.HandleServiceRequest(amfUe, models.AccessType__3_GPP_ACCESS, procedureCode, gmmMessage.ServiceRequest); err != nil {
		// 		return err
		// 	}
		// case nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration:
		// 	return gmm_handler.HandleDeregistrationRequest(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.DeregistrationRequestUEOriginatingDeregistration)
		// case nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration:
		// 	return gmm_handler.HandleDeregistrationAccept(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.DeregistrationAcceptUETerminatedDeregistration)
		// case nas.MsgTypeStatus5GMM:
		// 	if err := gmm_handler.HandleStatus5GMM(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.Status5GMM); err != nil {
		// 		return err
		// 	}
		default:
			nasLog.Errorf("Unknown GmmMessage[%d]\n", msg.GmmMessage.GetMessageType())
		}

	} else if msg.GsmMessage != nil {
		nasLog.Warn("GSM Message should include in GMM Message")
	} else {
		nasLog.Errorln("Nas Payload is Empty")
	}
	return

}
