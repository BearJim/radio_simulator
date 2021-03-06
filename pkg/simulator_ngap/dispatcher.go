package simulator_ngap

import (
	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/BearJim/radio_simulator/pkg/logger"
	"github.com/BearJim/radio_simulator/pkg/simulator_context"

	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
)

type NGController struct {
	ran           RanApp
	nasController NASController
}

func New(ranApp RanApp, nasController NASController) *NGController {
	c := &NGController{
		ran:           ranApp,
		nasController: nasController,
	}
	return c
}

type RanApp interface {
	Context() *simulator_context.RanContext
	NewAMF(*sctp.SCTPAddr)
	UnsetFailRecovering()
	Connect(*sctp.SCTPAddr) error
	SendToAMF(*sctp.SCTPAddr, []byte)
}

type NASController interface {
	NewNASConnection(*simulator_context.UeContext) error
	SendToNAS(ranUeNgapID int64, nasPdu []byte)
	CloseNASConnection(ranUeNgapID int64)
}

func (c *NGController) Dispatch(endpoint *sctp.SCTPAddr, msg []byte) {
	pdu, err := ngap.Decoder(msg)
	if err != nil {
		logger.NgapLog.Errorf("NGAP decode error: %s", err)
		return
	}

	logger.NgapLog.Debugf("read from %s", endpoint.String())

	switch pdu.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := pdu.InitiatingMessage
		if initiatingMessage == nil {
			logger.NgapLog.Error("Initiating Message is nil")
			return
		}
		switch initiatingMessage.ProcedureCode.Value {
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			c.handleAMFConfigurationUpdate(endpoint, pdu)
		case ngapType.ProcedureCodeDownlinkNASTransport:
			c.handleDownlinkNASTransport(endpoint, pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			logger.NgapLog.Infof("Handle Initial Context Setup Request")
			c.handleInitialContextSetupRequest(endpoint, pdu)
		case ngapType.ProcedureCodeUEContextRelease:
			logger.NgapLog.Infof("Handle Ue Context Release Command")
			c.HandleUeContextReleaseCommand(endpoint, pdu)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			logger.NgapLog.Infof("Handle Pdu Session Resource Setup Request")
			c.HandlePduSessionResourceSetupRequest(endpoint, pdu)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			logger.NgapLog.Infof("Handle Pdu Session Resource Release Command")
			c.HandlePduSessionResourceReleaseCommand(endpoint, pdu)
		default:
			logger.NgapLog.Warnf("Not implemented Initiating Message (procedureCode:%d)", initiatingMessage.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := pdu.SuccessfulOutcome
		if successfulOutcome == nil {
			logger.NgapLog.Error("successful Outcome is nil")
			return
		}
		switch successfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			c.handleNGSetupResponse(endpoint, pdu)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			c.handleRanConfigurationUpdateAcknowledge(endpoint, pdu)
		default:
			logger.NgapLog.Warnf("Not implemented SuccessfulOutcome (procedureCode:%d)", successfulOutcome.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := pdu.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			logger.NgapLog.Error("unsuccessful Outcome is nil")
			return
		}
		switch unsuccessfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			c.handleRanConfigurationUpdateFailure(endpoint, pdu)
		default:
			logger.NgapLog.Warnf("Not implemented UnsuccessfulOutcome (procedureCode:%d)", unsuccessfulOutcome.ProcedureCode.Value)
		}
	}
}

func (c *NGController) NewNASConnection(ue *simulator_context.UeContext) {
	if err := c.nasController.NewNASConnection(ue); err != nil {
		logger.NgapLog.Error(err)
	}
}

func (c *NGController) CloseNASConnection(ranUeNgapID int64) {
	c.nasController.CloseNASConnection(ranUeNgapID)
}

func (c *NGController) SendNAS(ranUeNgapID int64, nasPdu []byte) {
	c.nasController.SendToNAS(ranUeNgapID, nasPdu)
}
