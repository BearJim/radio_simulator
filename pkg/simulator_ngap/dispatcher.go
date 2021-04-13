package simulator_ngap

import (
	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"

	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
	"github.com/sirupsen/logrus"
)

var ngapLog *logrus.Entry

func init() {
	ngapLog = logger.NgapLog
}

type NGController struct {
	ran           RanApp
	nasController NASController

	// handlers map[int64]func(*sctp.SCTPAddr, *ngapType.NGAPPDU)
}

func New(ranApp RanApp, nasController NASController) *NGController {
	c := &NGController{
		ran:           ranApp,
		nasController: nasController,
	}
	// c.registerHandlers()
	return c
}

type RanApp interface {
	Context() *simulator_context.RanContext
	Connect(*sctp.SCTPAddr) error
	SendToAMF(*sctp.SCTPAddr, []byte)
}

type NASController interface {
	HandleNAS(*simulator_context.UeContext, []byte)
}

// func (c *NGController) registerHandlers() {
// 	c.handlers = make(map[int64]func(*sctp.SCTPAddr, *ngapType.NGAPPDU))

// 	// InitiatingMessage
// 	c.handlers[ngapType.ProcedureCodeAMFConfigurationUpdate] = c.handleAMFConfigurationUpdate
// 	c.handlers[ngapType.ProcedureCodeDownlinkNASTransport] = c.handleDownlinkNASTransport
// 	c.handlers[ngapType.ProcedureCodeInitialContextSetup] = c.handleInitialContextSetupRequest
// 	c.handlers[ngapType.ProcedureCodeUEContextRelease] = c.HandleUeContextReleaseCommand
// 	c.handlers[ngapType.ProcedureCodePDUSessionResourceSetup] = c.HandlePduSessionResourceSetupRequest
// 	c.handlers[ngapType.ProcedureCodePDUSessionResourceRelease] = c.HandlePduSessionResourceReleaseCommand

// 	// SuccessfulOutcome
// 	c.handlers[ngapType.ProcedureCodeNGSetup] = c.HandleNGSetupResponse
// 	c.handlers[ngapType.ProcedureCodeRANConfigurationUpdate] = c.handleRanConfigurationUpdateAcknowledge
// }

func (c *NGController) Dispatch(endpoint *sctp.SCTPAddr, msg []byte) {
	pdu, err := ngap.Decoder(msg)
	if err != nil {
		ngapLog.Errorf("NGAP decode error: %s", err)
		return
	}

	logger.NgapLog.Infof("read from %s", endpoint.String())

	switch pdu.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := pdu.InitiatingMessage
		if initiatingMessage == nil {
			ngapLog.Errorln("Initiating Message is nil")
			return
		}
		switch initiatingMessage.ProcedureCode.Value {
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			c.handleAMFConfigurationUpdate(endpoint, pdu)
		case ngapType.ProcedureCodeDownlinkNASTransport:
			ngapLog.Infof("Handle Downlink NAS Transport")
			c.handleDownlinkNASTransport(endpoint, pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			ngapLog.Infof("Handle Initial Context Setup Request")
			c.handleInitialContextSetupRequest(endpoint, pdu)
		case ngapType.ProcedureCodeUEContextRelease:
			ngapLog.Infof("Handle Ue Context Release Command")
			c.HandleUeContextReleaseCommand(endpoint, pdu)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			ngapLog.Infof("Handle Pdu Session Resource Setup Request")
			c.HandlePduSessionResourceSetupRequest(endpoint, pdu)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			ngapLog.Infof("Handle Pdu Session Resource Release Command")
			c.HandlePduSessionResourceReleaseCommand(endpoint, pdu)
		default:
			ngapLog.Warnf("Not implemented Initiating Message (procedureCode:%d)", initiatingMessage.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := pdu.SuccessfulOutcome
		if successfulOutcome == nil {
			ngapLog.Errorln("successful Outcome is nil")
			return
		}
		switch successfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			ngapLog.Info("Handle NG Setup Response")
			c.HandleNGSetupResponse(endpoint, pdu)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			ngapLog.Info("Handle RAN Configuration Update Acknowledge")
			c.handleRanConfigurationUpdateAcknowledge(endpoint, pdu)
		default:
			ngapLog.Warnf("Not implemented SuccessfulOutcome (procedureCode:%d)", successfulOutcome.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := pdu.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			ngapLog.Errorln("unsuccessful Outcome is nil")
			return
		}
		switch unsuccessfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			c.handleRanConfigurationUpdateFailure(endpoint, pdu)
		default:
			ngapLog.Warnf("Not implemented UnsuccessfulOutcome (procedureCode:%d)", unsuccessfulOutcome.ProcedureCode.Value)
		}
	}
}
