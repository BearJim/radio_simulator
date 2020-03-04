package ngap_handler

import (
	"github.com/sirupsen/logrus"
	"radio_simulator/lib/ngap"
	"radio_simulator/lib/ngap/ngapType"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_context"
)

var NgapLog *logrus.Entry

func init() {
	NgapLog = logger.NgapLog
}

func Dispatch(addr string, msg []byte) {
	_, ok := simulator_context.Simulator_Self().RanPool[addr]
	if !ok {
		NgapLog.Errorf("Cannot find the coressponding RAN Context\n")
		return
	}
	pdu, err := ngap.Decoder(msg)
	if err != nil {
		NgapLog.Errorf("NGAP decode error : %s\n", err)
		return
	}
	switch pdu.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := pdu.InitiatingMessage
		if initiatingMessage == nil {
			NgapLog.Errorln("Initiating Message is nil")
			return
		}
		switch initiatingMessage.ProcedureCode.Value {
		// case ngapType.ProcedureCodeNGSetup:
		// 	ngap_handler.HandleNGSetupRequest(ran, pdu)
		// case ngapType.ProcedureCodeInitialUEMessage:
		// 	ngap_handler.HandleInitialUEMessage(ran, pdu)
		// case ngapType.ProcedureCodeUplinkNASTransport:
		// 	ngap_handler.HandleUplinkNasTransport(ran, pdu)
		// case ngapType.ProcedureCodeNGReset:
		// 	ngap_handler.HandleNGReset(ran, pdu)
		// case ngapType.ProcedureCodeHandoverCancel:
		// 	ngap_handler.HandleHandoverCancel(ran, pdu)
		// case ngapType.ProcedureCodeUEContextReleaseRequest:
		// 	ngap_handler.HandleUEContextReleaseRequest(ran, pdu)
		// case ngapType.ProcedureCodeNASNonDeliveryIndication:
		// 	ngap_handler.HandleNasNonDeliveryIndication(ran, pdu)
		// case ngapType.ProcedureCodeLocationReportingFailureIndication:
		// 	ngap_handler.HandleLocationReportingFailureIndication(ran, pdu)
		// case ngapType.ProcedureCodeErrorIndication:
		// 	ngap_handler.HandleErrorIndication(ran, pdu)
		// case ngapType.ProcedureCodeUERadioCapabilityInfoIndication:
		// 	ngap_handler.HandleUERadioCapabilityInfoIndication(ran, pdu)
		// case ngapType.ProcedureCodeHandoverNotification:
		// 	ngap_handler.HandleHandoverNotify(ran, pdu)
		// case ngapType.ProcedureCodeHandoverPreparation:
		// 	ngap_handler.HandleHandoverRequired(ran, pdu)
		// case ngapType.ProcedureCodeRANConfigurationUpdate:
		// 	ngap_handler.HandleRanConfigurationUpdate(ran, pdu)
		// case ngapType.ProcedureCodeRRCInactiveTransitionReport:
		// 	ngap_handler.HandleRRCInactiveTransitionReport(ran, pdu)
		// case ngapType.ProcedureCodePDUSessionResourceNotify:
		// 	ngap_handler.HandlePDUSessionResourceNotify(ran, pdu)
		// case ngapType.ProcedureCodePathSwitchRequest:
		// 	ngap_handler.HandlePathSwitchRequest(ran, pdu)
		// case ngapType.ProcedureCodeLocationReport:
		// 	ngap_handler.HandleLocationReport(ran, pdu)
		// case ngapType.ProcedureCodeUplinkUEAssociatedNRPPaTransport:
		// 	ngap_handler.HandleUplinkUEAssociatedNRPPATransport(ran, pdu)
		// case ngapType.ProcedureCodeUplinkRANConfigurationTransfer:
		// 	ngap_handler.HandleUplinkRanConfigurationTransfer(ran, pdu)
		// case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
		// 	ngap_handler.HandlePDUSessionResourceModifyIndication(ran, pdu)
		// case ngapType.ProcedureCodeCellTrafficTrace:
		// 	ngap_handler.HandleCellTrafficTrace(ran, pdu)
		// case ngapType.ProcedureCodeUplinkRANStatusTransfer:
		// 	ngap_handler.HandleUplinkRanStatusTransfer(ran, pdu)
		// case ngapType.ProcedureCodeUplinkNonUEAssociatedNRPPaTransport:
		// 	ngap_handler.HandleUplinkNonUEAssociatedNRPPATransport(ran, pdu)
		default:
			NgapLog.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, initiatingMessage.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := pdu.SuccessfulOutcome
		if successfulOutcome == nil {
			NgapLog.Errorln("successful Outcome is nil")
			return
		}
		switch successfulOutcome.ProcedureCode.Value {
		// case ngapType.ProcedureCodeNGReset:
		// 	ngap_handler.HandleNGResetAcknowledge(ran, pdu)
		// case ngapType.ProcedureCodeUEContextRelease:
		// 	ngap_handler.HandleUEContextReleaseComplete(ran, pdu)
		// case ngapType.ProcedureCodePDUSessionResourceRelease:
		// 	ngap_handler.HandlePDUSessionResourceReleaseResponse(ran, pdu)
		// case ngapType.ProcedureCodeUERadioCapabilityCheck:
		// 	ngap_handler.HandleUERadioCapabilityCheckResponse(ran, pdu)
		// case ngapType.ProcedureCodeAMFConfigurationUpdate:
		// 	ngap_handler.HandleAMFconfigurationUpdateAcknowledge(ran, pdu)
		// case ngapType.ProcedureCodeInitialContextSetup:
		// 	ngap_handler.HandleInitialContextSetupResponse(ran, pdu)
		// case ngapType.ProcedureCodeUEContextModification:
		// 	ngap_handler.HandleUEContextModificationResponse(ran, pdu)
		// case ngapType.ProcedureCodePDUSessionResourceSetup:
		// 	ngap_handler.HandlePDUSessionResourceSetupResponse(ran, pdu)
		// case ngapType.ProcedureCodePDUSessionResourceModify:
		// 	ngap_handler.HandlePDUSessionResourceModifyResponse(ran, pdu)
		// case ngapType.ProcedureCodeHandoverResourceAllocation:
		// 	ngap_handler.HandleHandoverRequestAcknowledge(ran, pdu)
		default:
			NgapLog.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, successfulOutcome.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := pdu.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			NgapLog.Errorln("unsuccessful Outcome is nil")
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
			NgapLog.Warnf("Not implemented(choice:%d, procedureCode:%d)\n", pdu.Present, unsuccessfulOutcome.ProcedureCode.Value)
		}

	}

}
