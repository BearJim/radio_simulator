package simulator_ngap

import (
	"encoding/binary"
	"fmt"

	"github.com/BearJim/radio_simulator/pkg/logger"
	"github.com/BearJim/radio_simulator/pkg/simulator_context"

	"github.com/free5gc/aper"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
)

func handlePDUSessionResourceSetupRequestTransfer(sess *simulator_context.SessionContext, b []byte) ([]byte, error) {

	var pduSessionType *ngapType.PDUSessionType
	var ulNGUUPTNLInformation *ngapType.UPTransportLayerInformation
	var qosFlowSetupRequestList *ngapType.QosFlowSetupRequestList
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	transfer := ngapType.PDUSessionResourceSetupRequestTransfer{}

	err := aper.UnmarshalWithParams(b, &transfer, "valueExt")

	if err != nil {
		cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage)
		unsuccessfulTransfer, _ := BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
		return unsuccessfulTransfer, fmt.Errorf("PduSession Transfer IE format Error")
	}

	for _, ie := range transfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDPDUSessionType:
			pduSessionType = ie.Value.PDUSessionType
			if pduSessionType == nil {
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDULNGUUPTNLInformation:
			ulNGUUPTNLInformation = ie.Value.ULNGUUPTNLInformation
			if ulNGUUPTNLInformation == nil {
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDQosFlowSetupRequestList:
			qosFlowSetupRequestList = ie.Value.QosFlowSetupRequestList
			if qosFlowSetupRequestList == nil {
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage)
		criticalityDiagnostics := buildCriticalityDiagnostics(nil, nil, nil, &iesCriticalityDiagnostics)
		unsuccessfulTransfer, _ := BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, &criticalityDiagnostics)
		return unsuccessfulTransfer, fmt.Errorf("PduSession Transfer IE format Error")
	}

	// PDU Session Type
	switch pduSessionType.Value {
	case ngapType.PDUSessionTypePresentIpv4:
	default:
		err = fmt.Errorf("Pdu Session Type has not support for non-ipv4 case yet")
		cause := buildCause(ngapType.CausePresentRadioNetwork, ngapType.CauseRadioNetworkPresentRadioResourcesNotAvailable)
		unsuccessfulTransfer, _ := BuildPDUSessionResourceSetupUnsuccessfulTransfer(*cause, nil)
		return unsuccessfulTransfer, err
	}

	// UL NG-U UP TNL Information
	v4, _ := ngapConvert.IPAddressToString(ulNGUUPTNLInformation.GTPTunnel.TransportLayerAddress)
	sess.Mtx.Lock()
	sess.ULAddr = v4
	sess.ULTEID = binary.BigEndian.Uint32(ulNGUUPTNLInformation.GTPTunnel.GTPTEID.Value)

	// QoS Flow Setup Request List
	for _, item := range qosFlowSetupRequestList.List {
		qosFlow := simulator_context.QosFlow{
			Identifier: item.QosFlowIdentifier.Value,
			Parameters: item.QosFlowLevelQosParameters,
		}
		sess.QosFlows[qosFlow.Identifier] = &qosFlow
	}
	// Flag Set To Zero
	// sess.NewGtpHeader(0, 0, 0)
	sess.Mtx.Unlock()

	successfulTransfer, err := BuildPDUSessionResourceSetupResponseTransfer(sess)
	if err != nil {
		logger.NgapLog.Errorf("Encode PDUSessionResourceSetupResponseTransfer Error: %+v\n", err)
	}
	return successfulTransfer, nil
}

func handlePDUSessionResourceReleaseCommandTransfer(sess *simulator_context.SessionContext, b []byte) ([]byte, error) {

	transfer := ngapType.PDUSessionResourceReleaseCommandTransfer{}

	err := aper.UnmarshalWithParams(b, &transfer, "valueExt")

	successfulTransfer, _ := BuildPDUSessionResourceReleaseResponseTransfer()
	return successfulTransfer, err
}
