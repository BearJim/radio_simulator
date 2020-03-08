package simulator_context

import (
	"encoding/hex"
	"git.cs.nctu.edu.tw/calee/sctp"
	"radio_simulator/lib/aper"
	"radio_simulator/lib/ngap/ngapType"
)

type RanContext struct {
	RanUeIDGeneator int64
	AMFUri          string
	RanUri          string
	Name            string
	GnbId           aper.BitString
	UePool          map[int64]*UeContext // ranUeNgapId
	DefaultTAC      string
	SupportTAList   map[string][]PlmnSupportItem // TAC(hex string) -> PlmnSupportItem
	SctpConn        *sctp.SCTPConn
}

type PlmnSupportItem struct {
	PlmnId     ngapType.PLMNIdentity
	SNssaiList []ngapType.SNSSAI
}

func (ran *RanContext) FindUeByRanUeNgapID(ranUeNgapID int64) *UeContext {
	if ue, ok := ran.UePool[ranUeNgapID]; ok {
		return ue
	} else {
		return nil
	}
}

func (ran *RanContext) FindUeByAmfUeNgapID(amfUeNgapID int64) *UeContext {
	for _, ranUe := range ran.UePool {
		if ranUe.AmfUeNgapId == amfUeNgapID {
			return ranUe
		}
	}
	return nil
}

func (ran *RanContext) GetUserLocation() ngapType.UserLocationInformation {
	userLocationInformation := ngapType.UserLocationInformation{}
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationNR
	userLocationInformation.UserLocationInformationNR = new(ngapType.UserLocationInformationNR)

	userLocationInformationNR := userLocationInformation.UserLocationInformationNR
	userLocationInformationNR.NRCGI.PLMNIdentity = ran.SupportTAList[ran.DefaultTAC][0].PlmnId
	userLocationInformationNR.NRCGI.NRCellIdentity.Value = aper.BitString{
		Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x10},
		BitLength: 36,
	}

	userLocationInformationNR.TAI.PLMNIdentity = ran.SupportTAList[ran.DefaultTAC][0].PlmnId
	userLocationInformationNR.TAI.TAC.Value, _ = hex.DecodeString(ran.DefaultTAC)
	return userLocationInformation
}

func NewRanContext() *RanContext {
	return &RanContext{
		RanUeIDGeneator: 0,
		UePool:          make(map[int64]*UeContext),
		SupportTAList:   make(map[string][]PlmnSupportItem),
	}
}
