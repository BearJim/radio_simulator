package simulator_nas

import (
	"context"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
	"github.com/jay16213/radio_simulator/pkg/simulator_nas/nas_security"
	"go.uber.org/zap"

	"github.com/free5gc/nas"
)

var nasLog *zap.SugaredLogger

func init() {
	nasLog = logger.NASLog
}

func checkMsgError(err error, msg string) {
	if err != nil {
		logger.NASLog.Errorf("Handle %s Error: %s", msg, err.Error())
	}
}

type Routine interface {
	Run(context.Context)
}

type NASController struct {
	n1ConnectionsQueue chan Routine
	ngMessager         NGMessager
	cancelCtx          context.CancelFunc
}

func New() *NASController {
	return &NASController{
		n1ConnectionsQueue: make(chan Routine, 512),
	}
}

func (c *NASController) SetNGMessager(messager NGMessager) {
	c.ngMessager = messager
}

type NGMessager interface {
	SendUplinkNASTransport(*sctp.SCTPAddr, *simulator_context.UeContext, []byte)
	SendInitailUeMessage_RegistraionRequest(*sctp.SCTPAddr, *simulator_context.UeContext)
}

func (c *NASController) Run() error {
	return c.dispatch()
}

func (c *NASController) dispatch() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c.cancelCtx = cancel
	for {
		n1Connection := <-c.n1ConnectionsQueue // take a n1 connection
		go func(ctx context.Context) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			n1Connection.Run(ctx)
		}(ctx)
	}
}

func (c *NASController) Stop() {
	c.cancelCtx()
}

func (c *NASController) NewNASConnection(ue *simulator_context.UeContext) chan []byte {
	nasCh := make(chan []byte, 16)
	c.n1ConnectionsQueue <- &NASRoutine{
		rid:           int(ue.RanUeNgapId),
		ue:            ue,
		NASController: c,
		nasPduCh:      nasCh,
	}
	return nasCh
}

func (c *NASController) handleGmmMessage(ue *simulator_context.UeContext, nasPdu []byte) {
	// GMM Message
	msg, err := nas_security.NASDecode(ue, nas.GetSecurityHeaderType(nasPdu)&0x0f, nasPdu)
	if err != nil {
		nasLog.Error(err.Error())
		return
	}

	switch msg.GmmMessage.GetMessageType() {
	case nas.MsgTypeAuthenticationRequest:
		checkMsgError(c.handleAuthenticationRequest(ue, msg.GmmMessage.AuthenticationRequest), "AuthenticationRequest")
	case nas.MsgTypeAuthenticationReject:
		checkMsgError(c.handleAuthenticationReject(ue, msg.GmmMessage.AuthenticationReject), "AuthenticationReject")
	case nas.MsgTypeRegistrationReject:
		checkMsgError(c.handleRegistrationReject(ue, msg.GmmMessage.RegistrationReject), "RegistrationReject")
	case nas.MsgTypeSecurityModeCommand:
		checkMsgError(c.handleSecurityModeCommand(ue, msg.GmmMessage.SecurityModeCommand), "SecurityModeCommand")
	case nas.MsgTypeRegistrationAccept:
		checkMsgError(c.handleRegistrationAccept(ue, msg.GmmMessage.RegistrationAccept), "RegistrationAccept")
	case nas.MsgTypeDeregistrationAcceptUEOriginatingDeregistration:
		checkMsgError(c.handleDeregistrationAccept(ue, msg.GmmMessage.DeregistrationAcceptUEOriginatingDeregistration), "DeregistraionAccept")
	case nas.MsgTypeDLNASTransport:
		checkMsgError(c.handleDLNASTransport(ue, msg.GmmMessage.DLNASTransport), "DLNASTransport")
	default:
		logger.NASLog.Errorf("Unknown GmmMessage[%d]\n", msg.GmmMessage.GetMessageType())
	}
}

func (c *NASController) handleGsmMessage(ue *simulator_context.UeContext, nasPdu []byte) {
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
}
