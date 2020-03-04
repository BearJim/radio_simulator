package simulator_context

import (
	"git.cs.nctu.edu.tw/calee/sctp"
	"radio_simulator/lib/aper"
	"radio_simulator/lib/ngap/ngapType"
)

type RanContext struct {
	AMFUri        string
	RanUri        string
	Name          string
	GnbId         aper.BitString
	UePool        map[string]*RanUeContext     // Supi
	SupportTAList map[string][]PlmnSupportItem // TAC(hex string) -> PlmnSupportItem
	SctpConn      *sctp.SCTPConn
}

type PlmnSupportItem struct {
	PlmnId     ngapType.PLMNIdentity
	SNssaiList []ngapType.SNSSAI
}

func NewRanContext() *RanContext {
	return &RanContext{
		UePool:        make(map[string]*RanUeContext),
		SupportTAList: make(map[string][]PlmnSupportItem),
	}
}
