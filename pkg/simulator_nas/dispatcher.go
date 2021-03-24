package simulator_nas

import (
	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
	"github.com/jay16213/radio_simulator/pkg/simulator_nas/nas_security"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"

	"github.com/sirupsen/logrus"
)

var nasLog *logrus.Entry

func init() {
	nasLog = logger.NASLog
}

func checkMsgError(err error, msg string) {
	if err != nil {
		nasLog.Errorf("Handle %s Error: %s", msg, err.Error())
	}
}

type NASController struct {
	ngMessager NGMessager
}

func NewController() *NASController {
	return &NASController{}
}

func (c *NASController) SetNGMessager(messager NGMessager) {
	c.ngMessager = messager
}

type NGMessager interface {
	SendUplinkNASTransport(*sctp.SCTPAddr, *simulator_context.UeContext, []byte)
}

func (c *NASController) HandleNAS(ue *simulator_context.UeContext, nasPdu []byte) {

	if ue == nil {
		nasLog.Error("Ue is nil")
		return
	}

	if nasPdu == nil {
		nasLog.Error("nasPdu is nil")
		return
	}

	if nas.GetEPD(nasPdu) == nasMessage.Epd5GSSessionManagementMessage {
		// GSM Message
		msg := new(nas.Message)
		err := msg.PlainNasDecode(&nasPdu)
		if err != nil {
			nasLog.Error(err.Error())
			return
		}
		switch msg.GsmMessage.GetMessageType() {
		case nas.MsgTypePDUSessionEstablishmentAccept:
			checkMsgError(c.handlePduSessionEstblishmentAccept(ue, msg.GsmMessage.PDUSessionEstablishmentAccept), "PduSessionEstblishmentAccept")
		case nas.MsgTypePDUSessionReleaseCommand:
			checkMsgError(c.handlePduSessionReleaseCommand(ue, msg.GsmMessage.PDUSessionReleaseCommand), "PduSessionReleaseCommand")
		default:
			nasLog.Errorf("Unknown GsmMessage[%d]\n", msg.GsmMessage.GetMessageType())
		}
		return
	}

	// GMM Message
	msg, err := nas_security.NASDecode(ue, nas.GetSecurityHeaderType(nasPdu)&0x0f, nasPdu)
	if err != nil {
		nasLog.Error(err.Error())
		return
	}

	switch msg.GmmMessage.GetMessageType() {
	case nas.MsgTypeAuthenticationRequest:
		checkMsgError(c.handleAuthenticationRequest(ue, msg.GmmMessage.AuthenticationRequest), "AuthenticationRequest")
	case nas.MsgTypeSecurityModeCommand:
		checkMsgError(c.handleSecurityModeCommand(ue, msg.GmmMessage.SecurityModeCommand), "SecurityModeCommand")
	case nas.MsgTypeRegistrationAccept:
		checkMsgError(c.handleRegistrationAccept(ue, msg.GmmMessage.RegistrationAccept), "RegistrationAccept")
	case nas.MsgTypeDeregistrationAcceptUEOriginatingDeregistration:
		checkMsgError(c.handleDeregistrationAccept(ue, msg.GmmMessage.DeregistrationAcceptUEOriginatingDeregistration), "DeregistraionAccept")
	case nas.MsgTypeDLNASTransport:
		checkMsgError(c.handleDLNASTransport(ue, msg.GmmMessage.DLNASTransport), "DLNASTransport")
	default:
		nasLog.Errorf("Unknown GmmMessage[%d]\n", msg.GmmMessage.GetMessageType())
	}
}
