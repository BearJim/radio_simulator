package simulator_context

import (
	"encoding/hex"
	"net"
	"sync"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/BearJim/radio_simulator/pkg/factory"
	"github.com/BearJim/radio_simulator/pkg/simulator_util"
	"github.com/free5gc/aper"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
)

const (
	MaxValueOfTeid uint32 = 0xffffffff
)

type RanContext struct {
	TEIDGenerator    uint32
	RanUeIDGenerator int64
	AmfSctpEndpoint  factory.SCTPEndpoint
	RanSctpEndpoint  factory.SCTPEndpoint
	RanGtpUri        net.UDPAddr
	UpfInfoList      map[string]*UpfInfo // upf ip as key
	PlmnID           ngapType.PLMNIdentity
	Name             string
	GnbId            aper.BitString
	/**/
	uePoolMu sync.RWMutex
	UePool   map[int64]*UeContext // ranUeNgapId
	/**/
	SessPool      map[uint32]*SessionContext
	DefaultTAC    string
	SupportTAList map[string][]PlmnSupportItem // TAC(hex string) -> PlmnSupportItem
	AmfPool       map[*sctp.SCTPAddr]*AMFContext
	SctpConn      *sctp.SCTPConn
}

type UpfInfo struct {
	Addr    net.UDPAddr
	GtpConn *net.UDPConn
}

type PlmnSupportItem struct {
	PlmnId     ngapType.PLMNIdentity
	SNssaiList []ngapType.SNSSAI
}

func (ran *RanContext) NewAMF(addr *sctp.SCTPAddr) {
	ran.AmfPool[addr] = &AMFContext{
		Addr: addr,
	}
}

func (ran *RanContext) AttachSession(sess *SessionContext) {
	sess.DLAddr = ran.RanGtpUri.IP.String()
	sess.DLTEID = ran.TEIDAlloc()
	ran.SessPool[sess.DLTEID] = sess
}

func (ran *RanContext) DetachSession(sess *SessionContext) {
	delete(ran.SessPool, sess.DLTEID)
}

func (ran *RanContext) TEIDAlloc() uint32 {
	ran.TEIDGenerator %= MaxValueOfTeid
	ran.TEIDGenerator++
	for {
		if _, double := ran.SessPool[ran.TEIDGenerator]; double {
			ran.TEIDGenerator++
		} else {
			break
		}
	}
	return ran.TEIDGenerator
}

// FIXME: complete this function at deregistration procedure
func (ran *RanContext) DeleteUE() {

}

func (ran *RanContext) FindUeByRanUeNgapID(ranUeNgapID int64) *UeContext {
	ran.uePoolMu.RLock()
	defer ran.uePoolMu.RUnlock()
	if ue, ok := ran.UePool[ranUeNgapID]; ok {
		return ue
	} else {
		return nil
	}
}

func (ran *RanContext) FindUeByAmfUeNgapID(amfUeNgapID int64) *UeContext {
	ran.uePoolMu.RLock()
	defer ran.uePoolMu.RUnlock()
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
		RanUeIDGenerator: 1,
		TEIDGenerator:    0,
		UePool:           make(map[int64]*UeContext),
		SessPool:         make(map[uint32]*SessionContext),
		SupportTAList:    make(map[string][]PlmnSupportItem),
		UpfInfoList:      make(map[string]*UpfInfo),
		AmfPool:          make(map[*sctp.SCTPAddr]*AMFContext),
	}
}

func (ran *RanContext) NewUE(supi string) *UeContext {
	if ue := ran.FindUEBySupi(supi); ue != nil {
		return ue
	} else {
		ue = NewUeContext()
		ue.Ran = ran
		ran.uePoolMu.Lock()
		ran.UePool[ran.RanUeIDGenerator] = ue
		ue.RanUeNgapId = ran.RanUeIDGenerator
		ran.RanUeIDGenerator++
		ran.uePoolMu.Unlock()
		ue.AMFEndpoint = nil
		ue.CmState = CmStateConnected
		return ue
	}
}

func (ran *RanContext) FindUEBySupi(supi string) *UeContext {
	ran.uePoolMu.RLock()
	defer ran.uePoolMu.RUnlock()
	for _, ue := range ran.UePool {
		if ue.Supi == supi {
			return ue
		}
	}
	return nil
}

func (ran *RanContext) LoadConfig(cfg factory.Config) {
	ran.PlmnID = ngapConvert.PlmnIdToNgap(cfg.GnbId.PlmnId)
	ran.AmfSctpEndpoint = cfg.AmfSCTPEndpoint
	ran.RanSctpEndpoint = cfg.RanSctpEndpoint
	ran.RanGtpUri.IP = cfg.RanGtpUri.IP
	ran.RanGtpUri.Port = cfg.RanGtpUri.Port
	ran.Name = cfg.RanName
	ran.GnbId.BitLength = uint64(cfg.GnbId.BitLength)
	ran.GnbId.Bytes, _ = hex.DecodeString(cfg.GnbId.Value)
	for _, upfUri := range cfg.UpfUriList {
		ran.UpfInfoList[upfUri.IP.String()] = &UpfInfo{
			Addr: upfUri,
		}
	}
	for _, supportItem := range cfg.SupportTAList {
		plmnList := []PlmnSupportItem{}
		for _, item := range supportItem.Plmnlist {
			plmnItem := PlmnSupportItem{}
			plmnItem.PlmnId = ngapConvert.PlmnIdToNgap(item.PlmnId)
			for _, snssai := range item.SNssaiList {
				sNssaiNgap := ngapConvert.SNssaiToNgap(snssai)
				plmnItem.SNssaiList = append(plmnItem.SNssaiList, sNssaiNgap)
			}
			plmnList = append(plmnList, plmnItem)
		}
		tac := simulator_util.TACConfigToHexString(supportItem.Tac)
		ran.DefaultTAC = tac
		ran.SupportTAList[tac] = plmnList
	}
}
