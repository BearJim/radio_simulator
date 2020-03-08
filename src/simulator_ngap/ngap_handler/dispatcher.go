package ngap_handler

import (
	"github.com/sirupsen/logrus"
	"radio_simulator/lib/ngap"
	"radio_simulator/lib/ngap/ngapType"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_context"
)

var ngapLog *logrus.Entry

func init() {
	ngapLog = logger.NgapLog
}

func Dispatch(addr string, msg []byte) {
	ran, ok := simulator_context.Simulator_Self().RanPool[addr]
	if !ok {
		ngapLog.Errorf("Cannot find the coressponding RAN Context\n")
		return
	}
	pdu, err := ngap.Decoder(msg)
	if err != nil {
		ngapLog.Errorf("NGAP decode error : %s\n", err)
		return
	}
	switch pdu.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := pdu.InitiatingMessage
		if initiatingMessage == nil {
			ngapLog.Errorln("Initiating Message is nil")
			return
		}
		switch initiatingMessage.ProcedureCode.Value {
		case ngapType.ProcedureCodeDownlinkNASTransport:
			HandleDownlinkNASTransport(ran, pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			HandleInitialContextSetupRequest(ran, pdu)
		case ngapType.ProcedureCodeUEContextRelease:
			HandleUeContextReleaseCommand(ran, pdu)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			HandlePduSessionResourceSetupRequest(ran, pdu)
		default:
			ngapLog.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, initiatingMessage.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := pdu.SuccessfulOutcome
		if successfulOutcome == nil {
			ngapLog.Errorln("successful Outcome is nil")
			return
		}
		switch successfulOutcome.ProcedureCode.Value {

		default:
			ngapLog.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, successfulOutcome.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := pdu.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			ngapLog.Errorln("unsuccessful Outcome is nil")
			return
		}
		switch unsuccessfulOutcome.ProcedureCode.Value {
		// case ngapType.ProcedureCodeAMFConfigurationUpdate:
		// 	ngap_handler.HandleAMFconfigurationUpdateFailure(ran, pdu)
		// case ngapType.ProcedureCodeInitialContextSetup:
		// 	ngap_handler.HandleInitialContextSetupFailure(ran, pdu)
		// case ngapType.ProcedureCodeUEContextModification:
		// 	ngap_handler.HandleUEContextModificationFailure(ran, pdu)
		// case ngapType.ProcedureCodeHandoverResourceAllocation:
		// 	ngap_handler.HandleHandoverFailure(ran, pdu)
		default:
			ngapLog.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, unsuccessfulOutcome.ProcedureCode.Value)
		}

	}

}
