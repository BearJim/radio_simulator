package simulator_context

import (
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

func (context *RanContext) FindUeByRanUeNgapID(ranUeNgapID int64) *UeContext {
	if ue, ok := context.UePool[ranUeNgapID]; ok {
		return ue
	} else {
		return nil
	}
}

func NewRanContext() *RanContext {
	return &RanContext{
		RanUeIDGeneator: 0,
		UePool:          make(map[int64]*UeContext),
		SupportTAList:   make(map[string][]PlmnSupportItem),
	}
}
