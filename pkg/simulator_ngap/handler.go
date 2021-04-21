package simulator_ngap

import (
	"net"
	"reflect"
	"time"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"

	"github.com/free5gc/aper"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
)

func (c *NGController) handleNGSetupResponse(endpoint *sctp.SCTPAddr, message *ngapType.NGAPPDU) {
	var (
	// amfName         *ngapType.AMFName
	// servedGuamiList *ngapType.ServedGUAMIList
	)

	logger.NgapLog.Info("Handle NG Setup Response")

	ngSetupResponse := message.SuccessfulOutcome.Value.NGSetupResponse
	if ngSetupResponse == nil {
		logger.NgapLog.Error("NGSetupResponse is nil")
		return
	}

	for _, ie := range ngSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFName:
			// amfName = ie.Value.AMFName
		case ngapType.ProtocolIEIDServedGUAMIList:
			// servedGuamiList = ie.Value.ServedGUAMIList
		case ngapType.ProtocolIEIDRelativeAMFCapacity:
			logger.NgapLog.Debug("Decode IE RelativeAMFCapacity")
		case ngapType.ProtocolIEIDPLMNSupportList:
			logger.NgapLog.Debug("Decode IE PLMNSupportList")
		}
	}

	// amf.Name = amfName.Value
	// for _, item := range servedGuamiList.List {
	// 	plmnID := ngapConvert.PlmnIdToModels(item.GUAMI.PLMNIdentity)
	// 	amfID := ngapConvert.AmfIdToModels(item.GUAMI.AMFRegionID.Value, item.GUAMI.AMFSetID.Value, item.GUAMI.AMFPointer.Value)
	// 	guami := models.Guami{
	// 		PlmnId: &plmnID,
	// 		AmfId:  amfID,
	// 	}
	// 	if item.BackupAMFName != nil {
	// 		amf.ServedGUAMIList = append(amf.ServedGUAMIList, simulator_context.ServedGUAMI{
	// 			Guami:         guami,
	// 			BackupAMFName: item.BackupAMFName.Value,
	// 		})
	// 	} else {
	// 		amf.ServedGUAMIList = append(amf.ServedGUAMIList, simulator_context.ServedGUAMI{
	// 			Guami: guami,
	// 		})
	// 	}
	// }
}

func (c *NGController) handleDownlinkNASTransport(endpoint *sctp.SCTPAddr, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	// var oldAMF *ngapType.AMFName
	// var rANPagingPriority *ngapType.RANPagingPriority
	var nASPDU *ngapType.NASPDU
	// var mobilityRestrictionList *ngapType.MobilityRestrictionList
	// var indexToRFSP *ngapType.IndexToRFSP
	// var uEAggregateMaximumBitRate *ngapType.UEAggregateMaximumBitRate
	// var allowedNSSAI *ngapType.AllowedNSSAI

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	downlinkNASTransport := message.InitiatingMessage.Value.DownlinkNASTransport
	if downlinkNASTransport == nil {
		logger.NgapLog.Error("downlinkNASTransport is nil")
		return
	}

	logger.NgapLog.Infow("Handle Downlink NAS Transport", "amf", endpoint.String())
	for _, ie := range downlinkNASTransport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			logger.NgapLog.Debug("Decode IE AMFUENGAPID")
			aMFUENGAPID = ie.Value.AMFUENGAPID
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AMFUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			logger.NgapLog.Debug("Decode IE RANUENGAPID")
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDOldAMF:
			logger.NgapLog.Debug("Decode IE OldAMF")
			// oldAMF = ie.Value.OldAMF
		case ngapType.ProtocolIEIDRANPagingPriority:
			logger.NgapLog.Debug("Decode IE RANPagingPriority")
			// rANPagingPriority = ie.Value.RANPagingPriority
		case ngapType.ProtocolIEIDNASPDU:
			logger.NgapLog.Debug("Decode IE NASPDU")
			nASPDU = ie.Value.NASPDU
			if nASPDU == nil {
				logger.NgapLog.Error("NASPDU is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDMobilityRestrictionList:
			logger.NgapLog.Debug("Decode IE MobilityRestrictionList")
			// mobilityRestrictionList = ie.Value.MobilityRestrictionList
		case ngapType.ProtocolIEIDIndexToRFSP:
			logger.NgapLog.Debug("Decode IE IndexToRFSP")
			// indexToRFSP = ie.Value.IndexToRFSP
		case ngapType.ProtocolIEIDUEAggregateMaximumBitRate:
			logger.NgapLog.Debug("Decode IE UEAggregateMaximumBitRate")
			// uEAggregateMaximumBitRate = ie.Value.UEAggregateMaximumBitRate
		case ngapType.ProtocolIEIDAllowedNSSAI:
			logger.NgapLog.Debug("Decode IE AllowedNSSAI")
			// allowedNSSAI = ie.Value.AllowedNSSAI
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procudureCode := ngapType.ProcedureCodeDownlinkNASTransport
		trigger := ngapType.TriggeringMessagePresentInitiatingMessage
		criticality := ngapType.CriticalityPresentIgnore
		criticalityDiagnostics := buildCriticalityDiagnostics(&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		c.SendErrorIndication(endpoint, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	var ue *simulator_context.UeContext
	if rANUENGAPID != nil {
		ue = c.ran.Context().FindUeByRanUeNgapID(rANUENGAPID.Value)
		if ue == nil {
			logger.NgapLog.Warnf("No UE Context[RanUeNgapID:%d]", rANUENGAPID.Value)
			return
		}
	}

	if aMFUENGAPID != nil {
		if ue.AmfUeNgapId == simulator_context.AmfNgapIdUnspecified {
			logger.NgapLog.Debug("Create new logical UE-associated NG-connection")
			ue.AmfUeNgapId = aMFUENGAPID.Value
			ue.AMFEndpoint = endpoint
		} else {
			if ue.AmfUeNgapId != aMFUENGAPID.Value {
				logger.NgapLog.Warn("AMFUENGAPID unmatched")
				return
			}
		}
	}

	if !reflect.DeepEqual(ue.AMFEndpoint, endpoint) {
		logger.NgapLog.Warnw("AMF endpoint change", "supi", ue.Supi, "id", ue.AmfUeNgapId,
			"old", ue.AMFEndpoint.String(), "new", endpoint.String())
		ue.AMFEndpoint = endpoint
	}

	if nASPDU != nil {
		logger.NASLog.Infow("Forward Downlink NAS Transport", "amf", endpoint.String(), "supi", ue.Supi)
		c.nasController.HandleNAS(ue, nASPDU.Value)
	}
}

func (c *NGController) handleInitialContextSetupRequest(endpoint *sctp.SCTPAddr, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	// var oldAMF *ngapType.AMFName
	var uEAggregateMaximumBitRate *ngapType.UEAggregateMaximumBitRate
	// var coreNetworkAssistanceInformation *ngapType.CoreNetworkAssistanceInformation
	var gUAMI *ngapType.GUAMI
	// var pDUSessionResourceSetupListCxtReq *ngapType.PDUSessionResourceSetupListCxtReq
	var allowedNSSAI *ngapType.AllowedNSSAI
	var uESecurityCapabilities *ngapType.UESecurityCapabilities
	var securityKey *ngapType.SecurityKey
	// var traceActivation *ngapType.TraceActivation
	// var mobilityRestrictionList *ngapType.MobilityRestrictionList
	// var uERadioCapability *ngapType.UERadioCapability
	// var indexToRFSP *ngapType.IndexToRFSP
	// var maskedIMEISV *ngapType.MaskedIMEISV
	var nASPDU *ngapType.NASPDU
	// var emergencyFallbackIndicator *ngapType.EmergencyFallbackIndicator
	// var rRCInactiveTransitionReportRequest *ngapType.RRCInactiveTransitionReportRequest
	// var uERadioCapabilityForPaging *ngapType.UERadioCapabilityForPaging

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("InitiatingMessage is nil")
		return
	}

	initialContextSetupRequest := initiatingMessage.Value.InitialContextSetupRequest
	if initialContextSetupRequest == nil {
		logger.NgapLog.Error("initialContextSetupRequest is nil")
		return
	}

	for _, ie := range initialContextSetupRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			logger.NgapLog.Debug("Decode IE AMFUENGAPID")
			aMFUENGAPID = ie.Value.AMFUENGAPID
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AMFUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			logger.NgapLog.Debug("Decode IE RANUENGAPID")
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDOldAMF:
			logger.NgapLog.Debug("Decode IE OldAMF")
			// oldAMF = ie.Value.OldAMF
		case ngapType.ProtocolIEIDUEAggregateMaximumBitRate:
			logger.NgapLog.Debug("Decode IE UEAggregateMaximumBitRate")
			uEAggregateMaximumBitRate = ie.Value.UEAggregateMaximumBitRate
			if uEAggregateMaximumBitRate == nil {
				logger.NgapLog.Error("UEAggregateMaximumBitRate is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDCoreNetworkAssistanceInformation:
			logger.NgapLog.Debug("Decode IE CoreNetworkAssistanceInformation")
			// coreNetworkAssistanceInformation = ie.Value.CoreNetworkAssistanceInformation
		case ngapType.ProtocolIEIDGUAMI:
			logger.NgapLog.Debug("Decode IE GUAMI")
			gUAMI = ie.Value.GUAMI
			if gUAMI == nil {
				logger.NgapLog.Error("GUAMI is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDPDUSessionResourceSetupListCxtReq:
			logger.NgapLog.Debug("Decode IE PDUSessionResourceSetupListCxtReq")
			// pDUSessionResourceSetupListCxtReq = ie.Value.PDUSessionResourceSetupListCxtReq
		case ngapType.ProtocolIEIDAllowedNSSAI:
			logger.NgapLog.Debug("Decode IE AllowedNSSAI")
			allowedNSSAI = ie.Value.AllowedNSSAI
			if allowedNSSAI == nil {
				logger.NgapLog.Error("AllowedNSSAI is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDUESecurityCapabilities:
			logger.NgapLog.Debug("Decode IE UESecurityCapabilities")
			uESecurityCapabilities = ie.Value.UESecurityCapabilities
			if uESecurityCapabilities == nil {
				logger.NgapLog.Error("UESecurityCapabilities is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDSecurityKey:
			logger.NgapLog.Debug("Decode IE SecurityKey")
			securityKey = ie.Value.SecurityKey
			if securityKey == nil {
				logger.NgapLog.Error("SecurityKey is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDTraceActivation:
			logger.NgapLog.Debug("Decode IE TraceActivation")
			// traceActivation = ie.Value.TraceActivation
		case ngapType.ProtocolIEIDMobilityRestrictionList:
			logger.NgapLog.Debug("Decode IE MobilityRestrictionList")
			// mobilityRestrictionList = ie.Value.MobilityRestrictionList
		case ngapType.ProtocolIEIDUERadioCapability:
			logger.NgapLog.Debug("Decode IE UERadioCapability")
			// uERadioCapability = ie.Value.UERadioCapability
		case ngapType.ProtocolIEIDIndexToRFSP:
			logger.NgapLog.Debug("Decode IE IndexToRFSP")
			// indexToRFSP = ie.Value.IndexToRFSP
		case ngapType.ProtocolIEIDMaskedIMEISV:
			logger.NgapLog.Debug("Decode IE MaskedIMEISV")
			// maskedIMEISV = ie.Value.MaskedIMEISV
		case ngapType.ProtocolIEIDNASPDU:
			logger.NgapLog.Debug("Decode IE NASPDU")
			nASPDU = ie.Value.NASPDU
		case ngapType.ProtocolIEIDEmergencyFallbackIndicator:
			logger.NgapLog.Debug("Decode IE EmergencyFallbackIndicator")
			// emergencyFallbackIndicator = ie.Value.EmergencyFallbackIndicator
		case ngapType.ProtocolIEIDRRCInactiveTransitionReportRequest:
			logger.NgapLog.Debug("Decode IE RRCInactiveTransitionReportRequest")
			// rRCInactiveTransitionReportRequest = ie.Value.RRCInactiveTransitionReportRequest
		case ngapType.ProtocolIEIDUERadioCapabilityForPaging:
			logger.NgapLog.Debug("Decode IE UERadioCapabilityForPaging")
			// uERadioCapabilityForPaging = ie.Value.UERadioCapabilityForPaging
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procudureCode := ngapType.ProcedureCodeInitialContextSetup
		trigger := ngapType.TriggeringMessagePresentInitiatingMessage
		criticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		c.SendErrorIndication(endpoint, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	ue := c.ran.Context().FindUeByRanUeNgapID(rANUENGAPID.Value)
	if ue == nil {
		logger.NgapLog.Warnf("No UE Context[RanUeNgapID:%d]\n", rANUENGAPID.Value)
		return
	}

	// TODO: Service Request Case
	c.SendIntialContextSetupResponse(endpoint, ue, nil)

	if nASPDU != nil {
		c.nasController.HandleNAS(ue, nASPDU.Value)
	}

}

func (c *NGController) HandleUeContextReleaseCommand(endpoint *sctp.SCTPAddr, message *ngapType.NGAPPDU) {
	var uENGAPIDs *ngapType.UENGAPIDs
	var cause *ngapType.Cause

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("InitiatingMessage is nil")
		return
	}

	uEContextReleaseCommand := initiatingMessage.Value.UEContextReleaseCommand
	if uEContextReleaseCommand == nil {
		logger.NgapLog.Error("uEContextReleaseCommand is nil")
		return
	}

	for _, ie := range uEContextReleaseCommand.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDUENGAPIDs:
			logger.NgapLog.Debug("Decode IE UENGAPIDs")
			uENGAPIDs = ie.Value.UENGAPIDs
			if uENGAPIDs == nil {
				logger.NgapLog.Error("UENGAPIDs is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDCause:
			logger.NgapLog.Debug("Decode IE Cause")
			cause = ie.Value.Cause
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procudureCode := ngapType.ProcedureCodeUEContextRelease
		trigger := ngapType.TriggeringMessagePresentInitiatingMessage
		criticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		c.SendErrorIndication(endpoint, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	var ue *simulator_context.UeContext

	switch uENGAPIDs.Present {
	case ngapType.UENGAPIDsPresentAMFUENGAPID:
		ue = c.ran.Context().FindUeByAmfUeNgapID(uENGAPIDs.AMFUENGAPID.Value)
		if ue == nil {
			logger.NgapLog.Warnf("No UE Context[AmfUeNgapID:%d]", uENGAPIDs.AMFUENGAPID.Value)
			return
		}
	case ngapType.UENGAPIDsPresentUENGAPIDPair:
		pair := uENGAPIDs.UENGAPIDPair
		ue = c.ran.Context().FindUeByRanUeNgapID(pair.RANUENGAPID.Value)
		if ue == nil {
			logger.NgapLog.Warnf("No UE Context[RanUeNgapID:%d]", pair.RANUENGAPID.Value)
			return
		}
	}

	printAndGetCause(cause)
	c.SendUeContextReleaseComplete(endpoint, ue)
}

func (c *NGController) HandlePduSessionResourceSetupRequest(endpoint *sctp.SCTPAddr, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	// var rANPagingPriority *ngapType.RANPagingPriority
	var nASPDU *ngapType.NASPDU
	var pDUSessionResourceSetupListSUReq *ngapType.PDUSessionResourceSetupListSUReq

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("InitiatingMessage is nil")
		return
	}

	pDUSessionResourceSetupRequest := initiatingMessage.Value.PDUSessionResourceSetupRequest
	if pDUSessionResourceSetupRequest == nil {
		logger.NgapLog.Error("pDUSessionResourceSetupRequest is nil")
		return
	}

	for _, ie := range pDUSessionResourceSetupRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			logger.NgapLog.Debug("Decode IE AMFUENGAPID")
			aMFUENGAPID = ie.Value.AMFUENGAPID
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AMFUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			logger.NgapLog.Debug("Decode IE RANUENGAPID")
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANPagingPriority:
			logger.NgapLog.Debug("Decode IE RANPagingPriority")
			// rANPagingPriority = ie.Value.RANPagingPriority
		case ngapType.ProtocolIEIDNASPDU:
			logger.NgapLog.Debug("Decode IE NASPDU")
			nASPDU = ie.Value.NASPDU
		case ngapType.ProtocolIEIDPDUSessionResourceSetupListSUReq:
			logger.NgapLog.Debug("Decode IE PDUSessionResourceSetupListSUReq")
			pDUSessionResourceSetupListSUReq = ie.Value.PDUSessionResourceSetupListSUReq
			if pDUSessionResourceSetupListSUReq == nil {
				logger.NgapLog.Error("PDUSessionResourceSetupListSUReq is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procudureCode := ngapType.ProcedureCodePDUSessionResourceSetup
		trigger := ngapType.TriggeringMessagePresentInitiatingMessage
		criticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		c.SendErrorIndication(endpoint, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	ue := c.ran.Context().FindUeByRanUeNgapID(rANUENGAPID.Value)
	if ue == nil {
		logger.NgapLog.Warnf("No UE Context[RanUeNgapID:%d]\n", rANUENGAPID.Value)
		return
	}

	responseList := new(ngapType.PDUSessionResourceSetupListSURes)
	failedListSURes := new(ngapType.PDUSessionResourceFailedToSetupListSURes)

	for _, pduSession := range pDUSessionResourceSetupListSUReq.List {
		pduSessionId := pduSession.PDUSessionID.Value
		sess, exist := ue.PduSession[pduSessionId]
		if !exist {
			logger.NgapLog.Warnf("No PduSession Context[PduSessionId:%d]\n", pduSessionId)
			continue
		}
		if pduSession.PDUSessionNASPDU != nil {
			// Handle Nas Msg
			c.nasController.HandleNAS(ue, pduSession.PDUSessionNASPDU.Value)
		}
		sess.Mtx.Lock()
		c.ran.Context().AttachSession(sess)
		sess.Mtx.Unlock()
		resTransfer, err := handlePDUSessionResourceSetupRequestTransfer(sess, pduSession.PDUSessionResourceSetupRequestTransfer)
		if err == nil {
			AppendPDUSessionResourceSetupListSURes(responseList, pduSessionId, resTransfer)
			// build ULPDR, ULFAR, DLPDR
			simulator_context.Simulator_Self().AttachSession(sess)
		} else {
			logger.NgapLog.Warnf("Pdu Session Resource Setup Fail: %s", err.Error())
			AppendPDUSessionResourceFailedToSetupListSURes(failedListSURes, pduSessionId, resTransfer)
		}
	}
	c.SendPDUSessionResourceSetupResponse(endpoint, ue, responseList, failedListSURes)
	if nASPDU != nil {
		c.nasController.HandleNAS(ue, nASPDU.Value)
	}
}

func (c *NGController) HandlePduSessionResourceReleaseCommand(endpoint *sctp.SCTPAddr, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	// var rANPagingPriority *ngapType.RANPagingPriority
	var nASPDU *ngapType.NASPDU
	var pDUSessionResourceToReleaseListRelCmd *ngapType.PDUSessionResourceToReleaseListRelCmd

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("InitiatingMessage is nil")
		return
	}

	pDUSessionResourceReleaseCommand := initiatingMessage.Value.PDUSessionResourceReleaseCommand
	if pDUSessionResourceReleaseCommand == nil {
		logger.NgapLog.Error("pDUSessionResourceReleaseCommand is nil")
		return
	}

	for _, ie := range pDUSessionResourceReleaseCommand.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			logger.NgapLog.Debug("Decode IE AMFUENGAPID")
			aMFUENGAPID = ie.Value.AMFUENGAPID
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AMFUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			logger.NgapLog.Debug("Decode IE RANUENGAPID")
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANPagingPriority:
			logger.NgapLog.Debug("Decode IE RANPagingPriority")
			// rANPagingPriority = ie.Value.RANPagingPriority
		case ngapType.ProtocolIEIDNASPDU:
			logger.NgapLog.Debug("Decode IE NASPDU")
			nASPDU = ie.Value.NASPDU
		case ngapType.ProtocolIEIDPDUSessionResourceToReleaseListRelCmd:
			logger.NgapLog.Debug("Decode IE PDUSessionResourceToReleaseListRelCmd")
			pDUSessionResourceToReleaseListRelCmd = ie.Value.PDUSessionResourceToReleaseListRelCmd
			if pDUSessionResourceToReleaseListRelCmd == nil {
				logger.NgapLog.Error("PDUSessionResourceToReleaseListRelCmd is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procudureCode := ngapType.ProcedureCodePDUSessionResourceRelease
		trigger := ngapType.TriggeringMessagePresentInitiatingMessage
		criticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		c.SendErrorIndication(endpoint, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	ue := c.ran.Context().FindUeByRanUeNgapID(rANUENGAPID.Value)
	if ue == nil {
		logger.NgapLog.Warnf("No UE Context[RanUeNgapID:%d]\n", rANUENGAPID.Value)
		return
	}

	responseList := ngapType.PDUSessionResourceReleasedListRelRes{}
	for _, pduSession := range pDUSessionResourceToReleaseListRelCmd.List {
		pduSessionId := pduSession.PDUSessionID.Value
		sess, exist := ue.PduSession[pduSessionId]
		if !exist {
			logger.NgapLog.Warnf("No PduSession Context[PduSessionId:%d]\n", pduSessionId)
			continue
		}
		resTransfer, err := handlePDUSessionResourceReleaseCommandTransfer(sess, pduSession.PDUSessionResourceReleaseCommandTransfer)
		if err != nil {
			logger.NgapLog.Warn(err.Error())
		}
		AppendPDUSessionResourceReleasedListRelRes(&responseList, pduSessionId, resTransfer)
		c.ran.Context().DetachSession(sess)
	}

	c.SendPDUSessionResourceReleaseResponse(endpoint, ue, responseList, nil)

	if nASPDU != nil {
		c.nasController.HandleNAS(ue, nASPDU.Value)
	}
}

func (c *NGController) handleAMFConfigurationUpdate(endpoint *sctp.SCTPAddr, message *ngapType.NGAPPDU) {
	var (
		amfTNLAssociationToAddList *ngapType.AMFTNLAssociationToAddList
	)

	logger.NgapLog.Info("Handle AMF Configuration Update")

	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("InitiatingMessage is nil")
		return
	}

	amfConfigurationUpdate := initiatingMessage.Value.AMFConfigurationUpdate
	if amfConfigurationUpdate == nil {
		logger.NgapLog.Error("pDUSessionResourceReleaseCommand is nil")
		return
	}

	for _, ie := range amfConfigurationUpdate.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFTNLAssociationToAddList:
			amfTNLAssociationToAddList = ie.Value.AMFTNLAssociationToAddList
		case ngapType.ProtocolIEIDAMFTNLAssociationToRemoveList:
		case ngapType.ProtocolIEIDAMFTNLAssociationToUpdateList:
		}
	}

	amfTNLAssociationSetupList := ngapType.AMFTNLAssociationSetupList{}
	for _, item := range amfTNLAssociationToAddList.List {
		ipv4Addr, _ := ngapConvert.IPAddressToString(*item.AMFTNLAssociationAddress.EndpointIPAddress)
		sctpAddr := &sctp.SCTPAddr{
			IPAddrs: []net.IPAddr{
				{IP: net.ParseIP(ipv4Addr)},
			},
			Port: 38412,
		}
		if err := c.ran.Connect(sctpAddr); err != nil {
			logger.NgapLog.Error(err)
		} else {
			logger.NgapLog.Infof("establish additional TNL association with AMF success (addr: %s)", sctpAddr)
			setupItem := ngapType.AMFTNLAssociationSetupItem{
				AMFTNLAssociationAddress: ngapType.CPTransportLayerInformation{
					Present:           ngapType.CPTransportLayerInformationPresentEndpointIPAddress,
					EndpointIPAddress: item.AMFTNLAssociationAddress.EndpointIPAddress,
				},
			}
			amfTNLAssociationSetupList.List = append(amfTNLAssociationSetupList.List, setupItem)
			c.SendRanConfigurationUpdate(sctpAddr)
		}
	}
	c.SendAMFConfigurationUpdateAcknowledge(endpoint, &amfTNLAssociationSetupList)
	time.Sleep(200 * time.Millisecond)
}

func (c *NGController) handleRanConfigurationUpdateAcknowledge(endpoint *sctp.SCTPAddr, message *ngapType.NGAPPDU) {
	logger.NgapLog.Info("Handle RAN Configuration Update Acknowledge")
}

func (c *NGController) handleRanConfigurationUpdateFailure(endpoint *sctp.SCTPAddr, message *ngapType.NGAPPDU) {
	var (
		cause *ngapType.Cause
	)

	logger.NgapLog.Info("Handle RAN Configuration Update Failure")

	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	ranConfigurationUpdateFailure := message.UnsuccessfulOutcome.Value.RANConfigurationUpdateFailure
	if ranConfigurationUpdateFailure == nil {
		logger.NgapLog.Error("ranConfigurationUpdateFailure is nil")
		return
	}

	for _, ie := range ranConfigurationUpdateFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
		}
	}

	if cause != nil {
		printAndGetCause(cause)
	}
}

func buildCriticalityDiagnostics(
	procedureCode *int64,
	triggeringMessage *aper.Enumerated,
	procedureCriticality *aper.Enumerated,
	iesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList) (criticalityDiagnostics ngapType.CriticalityDiagnostics) {

	if procedureCode != nil {
		criticalityDiagnostics.ProcedureCode = new(ngapType.ProcedureCode)
		criticalityDiagnostics.ProcedureCode.Value = *procedureCode
	}

	if triggeringMessage != nil {
		criticalityDiagnostics.TriggeringMessage = new(ngapType.TriggeringMessage)
		criticalityDiagnostics.TriggeringMessage.Value = *triggeringMessage
	}

	if procedureCriticality != nil {
		criticalityDiagnostics.ProcedureCriticality = new(ngapType.Criticality)
		criticalityDiagnostics.ProcedureCriticality.Value = *procedureCriticality
	}

	if iesCriticalityDiagnostics != nil {
		criticalityDiagnostics.IEsCriticalityDiagnostics = iesCriticalityDiagnostics
	}

	return criticalityDiagnostics
}

func buildCriticalityDiagnosticsIEItem(ieCriticality aper.Enumerated, ieID int64, typeOfErr aper.Enumerated) (item ngapType.CriticalityDiagnosticsIEItem) {

	item = ngapType.CriticalityDiagnosticsIEItem{
		IECriticality: ngapType.Criticality{
			Value: ieCriticality,
		},
		IEID: ngapType.ProtocolIEID{
			Value: ieID,
		},
		TypeOfError: ngapType.TypeOfError{
			Value: typeOfErr,
		},
	}

	return item
}

func buildCause(present int, value aper.Enumerated) (cause *ngapType.Cause) {
	cause = new(ngapType.Cause)
	cause.Present = present

	switch present {
	case ngapType.CausePresentRadioNetwork:
		cause.RadioNetwork = new(ngapType.CauseRadioNetwork)
		cause.RadioNetwork.Value = value
	case ngapType.CausePresentTransport:
		cause.Transport = new(ngapType.CauseTransport)
		cause.Transport.Value = value
	case ngapType.CausePresentNas:
		cause.Nas = new(ngapType.CauseNas)
		cause.Nas.Value = value
	case ngapType.CausePresentProtocol:
		cause.Protocol = new(ngapType.CauseProtocol)
		cause.Protocol.Value = value
	case ngapType.CausePresentMisc:
		cause.Misc = new(ngapType.CauseMisc)
		cause.Misc.Value = value
	case ngapType.CausePresentNothing:
	}

	return
}

func printAndGetCause(cause *ngapType.Cause) (present int, value aper.Enumerated) {

	present = cause.Present
	switch cause.Present {
	case ngapType.CausePresentRadioNetwork:
		logger.NgapLog.Warnf("Cause RadioNetwork[%d]", cause.RadioNetwork.Value)
		value = cause.RadioNetwork.Value
	case ngapType.CausePresentTransport:
		logger.NgapLog.Warnf("Cause Transport[%d]", cause.Transport.Value)
		value = cause.Transport.Value
	case ngapType.CausePresentProtocol:
		logger.NgapLog.Warnf("Cause Protocol[%d]", cause.Protocol.Value)
		value = cause.Protocol.Value
	case ngapType.CausePresentNas:
		logger.NgapLog.Warnf("Cause Nas[%d]", cause.Nas.Value)
		value = cause.Nas.Value
	case ngapType.CausePresentMisc:
		logger.NgapLog.Warnf("Cause Misc[%d]", cause.Misc.Value)
		value = cause.Misc.Value
	default:
		logger.NgapLog.Errorf("Invalid Cause group[%d]", cause.Present)
	}
	return
}
