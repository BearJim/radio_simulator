package ngap_handler

import (
	"radio_simulator/lib/aper"
	"radio_simulator/lib/ngap/ngapType"
	"radio_simulator/src/simulator_context"
	"radio_simulator/src/simulator_nas"
	"radio_simulator/src/simulator_ngap"
)

func HandleDownlinkNASTransport(ran *simulator_context.RanContext, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var oldAMF *ngapType.AMFName
	var rANPagingPriority *ngapType.RANPagingPriority
	var nASPDU *ngapType.NASPDU
	var mobilityRestrictionList *ngapType.MobilityRestrictionList
	var indexToRFSP *ngapType.IndexToRFSP
	var uEAggregateMaximumBitRate *ngapType.UEAggregateMaximumBitRate
	var allowedNSSAI *ngapType.AllowedNSSAI

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if ran == nil {
		ngapLog.Error("RAN Context is nil")
		return
	}

	if message == nil {
		ngapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		ngapLog.Error("InitiatingMessage is nil")
		return
	}

	downlinkNASTransport := initiatingMessage.Value.DownlinkNASTransport
	if downlinkNASTransport == nil {
		ngapLog.Error("downlinkNASTransport is nil")
		return
	}

	for _, ie := range downlinkNASTransport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE AMFUENGAPID")
			aMFUENGAPID = ie.Value.AMFUENGAPID
			if aMFUENGAPID == nil {
				ngapLog.Error("AMFUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			ngapLog.Traceln("[NGAP] Decode IE RANUENGAPID")
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				ngapLog.Error("RANUENGAPID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDOldAMF:
			ngapLog.Traceln("[NGAP] Decode IE OldAMF")
			oldAMF = ie.Value.OldAMF
		case ngapType.ProtocolIEIDRANPagingPriority:
			ngapLog.Traceln("[NGAP] Decode IE RANPagingPriority")
			rANPagingPriority = ie.Value.RANPagingPriority
		case ngapType.ProtocolIEIDNASPDU:
			ngapLog.Traceln("[NGAP] Decode IE NASPDU")
			nASPDU = ie.Value.NASPDU
			if nASPDU == nil {
				ngapLog.Error("NASPDU is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ie.Id.Value, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDMobilityRestrictionList:
			ngapLog.Traceln("[NGAP] Decode IE MobilityRestrictionList")
			mobilityRestrictionList = ie.Value.MobilityRestrictionList
		case ngapType.ProtocolIEIDIndexToRFSP:
			ngapLog.Traceln("[NGAP] Decode IE IndexToRFSP")
			indexToRFSP = ie.Value.IndexToRFSP
		case ngapType.ProtocolIEIDUEAggregateMaximumBitRate:
			ngapLog.Traceln("[NGAP] Decode IE UEAggregateMaximumBitRate")
			uEAggregateMaximumBitRate = ie.Value.UEAggregateMaximumBitRate
		case ngapType.ProtocolIEIDAllowedNSSAI:
			ngapLog.Traceln("[NGAP] Decode IE AllowedNSSAI")
			allowedNSSAI = ie.Value.AllowedNSSAI
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procudureCode := ngapType.ProcedureCodeDownlinkNASTransport
		trigger := ngapType.TriggeringMessagePresentInitiatingMessage
		criticality := ngapType.CriticalityPresentIgnore
		criticalityDiagnostics := buildCriticalityDiagnostics(&procudureCode, &trigger, &criticality, &iesCriticalityDiagnostics)
		simulator_ngap.SendErrorIndication(ran, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	var ue *simulator_context.UeContext
	if rANUENGAPID != nil {
		ue = ran.FindUeByRanUeNgapID(rANUENGAPID.Value)
		if ue == nil {
			ngapLog.Warnf("No UE Context[RanUeNgapID:%d]\n", rANUENGAPID.Value)
			return
		}
	}

	if aMFUENGAPID != nil {
		if ue.AmfUeNgapId == simulator_context.NgapIdUnspecified {
			ngapLog.Tracef("Create new logical UE-associated NG-connection")
			ue.AmfUeNgapId = aMFUENGAPID.Value
			// n3iwfUe.SCTPAddr = amf.SCTPAddr
		} else {
			if ue.AmfUeNgapId != aMFUENGAPID.Value {
				ngapLog.Warn("AMFUENGAPID unmatched")
				return
			}
		}
	}

	if nASPDU != nil {
		simulator_nas.HandleNAS(nASPDU)
		// TODO: Send NAS PDU to UE
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
