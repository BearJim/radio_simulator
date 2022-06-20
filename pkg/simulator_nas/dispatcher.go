package simulator_nas

import (
	"context"
	"fmt"
	"sync"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/BearJim/radio_simulator/pkg/logger"
	"github.com/BearJim/radio_simulator/pkg/simulator_context"
	"github.com/BearJim/radio_simulator/pkg/simulator_nas/nas_security"
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

	mu            sync.RWMutex // protect the following fields
	n1Connections map[int64]*NASRoutine
}

func New() *NASController {
	return &NASController{
		n1ConnectionsQueue: make(chan Routine, 512),
		n1Connections:      make(map[int64]*NASRoutine),
	}
}

func (c *NASController) SetNGMessager(messager NGMessager) {
	c.ngMessager = messager
}

type NGMessager interface {
	SendUplinkNASTransport(*sctp.SCTPAddr, *simulator_context.UeContext, []byte)
	SendInitailUeMessage_RegistraionRequest(*simulator_context.UeContext)
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

func (c *NASController) NewNASConnection(ue *simulator_context.UeContext) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.n1Connections[ue.RanUeNgapId]; ok {
		return fmt.Errorf("NAS connection (ran_ue_ngap_id: %d) has exists", ue.RanUeNgapId)
	} else {
		r := &NASRoutine{
			rid:           int(ue.RanUeNgapId),
			ue:            ue,
			NASController: c,
			nasPduCh:      make(chan []byte, 4),
		}
		c.n1Connections[ue.RanUeNgapId] = r
		c.n1ConnectionsQueue <- r
	}
	return nil
}

func (c *NASController) SendToNAS(ranUeNgapID int64, nasPdu []byte) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if routine, ok := c.n1Connections[ranUeNgapID]; ok {
		routine.nasPduCh <- nasPdu
	} else {
		logger.NASLog.Error("Forward NAS message failed (ran_ue_ngap_id: %d)", ranUeNgapID)
	}
}

func (c *NASController) CloseNASConnection(ranUeNgapID int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if routine, ok := c.n1Connections[ranUeNgapID]; ok {
		close(routine.nasPduCh)
		delete(c.n1Connections, ranUeNgapID)
	}
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
	case nas.MsgTypeServiceAccept:
		checkMsgError(c.handleServiceAccept(ue, msg.GmmMessage.ServiceAccept), "ServiceAccept")
	case nas.MsgTypeServiceReject:
		checkMsgError(c.handleServiceReject(ue, msg.GmmMessage.ServiceReject), "ServiceReject")
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
