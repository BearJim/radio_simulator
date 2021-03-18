package simulator_context

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"sync"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/jay16213/radio_simulator/pkg/logger"

	"github.com/free5gc/UeauCommon"
	"github.com/free5gc/milenage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
)

const (
	RanNgapIdUnspecified int64 = 0xffffffff
	AmfNgapIdUnspecified int64 = 0xffffffffff
)
const (
	RegisterStateRegistered  = "REGISTERED"
	RegisterStateRegistering = "REGISTERING"
	RegisterStateDeregitered = "DEREGISTERED"
)

type UeContext struct {
	AMFEndpoint *sctp.SCTPAddr

	Supi          string `yaml:"supi" bson:"supi"`
	Guti          *nasType.GUTI5G
	Gpsis         []string                            `yaml:"gpsis" bson:"gpsis"`
	Nssai         models.Nssai                        `yaml:"nssai" bson:"nssai"`
	UeAmbr        UeAmbr                              `yaml:"ueAmbr" bson:"ueAmbr"`
	SmfSelData    models.SmfSelectionSubscriptionData `yaml:"smfSelData" bson:"smfSelectionSubscriptionData"`
	AuthData      AuthData                            `yaml:"auths" bson:"authData"`
	SubscCats     []string                            `yaml:"subscCats,omitempty" bson:"subscCats"`
	ServingPlmnId string                              `yaml:"servingPlmn" bson:"servingPlmn"`
	RanUeNgapId   int64
	AmfUeNgapId   int64
	// security
	ULCount         security.Count `bson:"uplinkCount,omitempty"`
	DLCount         security.Count `bson:"downlinkCount,omitempty"`
	CipheringAlgStr string         `yaml:"cipherAlg" bson:"cipherAlgStr"`
	IntegrityAlgStr string         `yaml:"integrityAlg" bson:"integrityAlgStr"`
	CipheringAlg    uint8          `bson:"cipherAlg"`
	IntegrityAlg    uint8          `bson:"integrityAlg"`
	KnasEnc         [16]uint8
	KnasInt         [16]uint8
	Kamf            []uint8
	NgKsi           uint8
	// PduSession
	PduSession map[int64]*SessionContext
	// related Context
	Ran     *RanContext
	RmState string `bson:"rmState"`
	// For TCP Client
	// TcpChannelMsg map[string]chan string
	// TcpConn       map[string]net.Conn // supi -> UeTcpClient
}

type UeAmbr struct {
	UpLink   string `yaml:"uplink" bson:"uplink"`
	DownLink string `yaml:"downlink" bson:"downlink"`
}

type AuthData struct {
	AuthMethod string `yaml:"authMethod" bson:"authMethod"`
	K          string `yaml:"K" bson:"K"`
	Opc        string `yaml:"Opc,omitempty" bson:"Opc"`
	Op         string `yaml:"Op,omitempty" bson:"Op"`
	AMF        string `yaml:"AMF" bson:"AMF"`
	SQN        string `yaml:"SQN" bson:"SQN"`
}

type SessionContext struct {
	Mtx sync.Mutex
	// GtpHdr       []byte
	// GtpHdrLen    uint16
	PduSessionId int64
	UeIp         string
	ULAddr       string
	ULTEID       uint32
	ULFarID      string
	ULPdrID      string
	DLAddr       string
	DLTEID       uint32
	DLPdrID      string // DLFarID = default far 1(just forward)
	Dnn          string
	Snssai       models.Snssai
	QosFlows     map[int64]*QosFlow // QosFlowIdentifier as key
	Ue           *UeContext
	// Sess Channel To Tcp Client
	SessTcpChannelMsg chan string
}

type QosFlow struct {
	Identifier int64
	Parameters ngapType.QosFlowLevelQosParameters
}

func NewUeContext() *UeContext {
	return &UeContext{
		PduSession:  make(map[int64]*SessionContext),
		AmfUeNgapId: AmfNgapIdUnspecified,
		RanUeNgapId: RanNgapIdUnspecified,
		RmState:     RegisterStateDeregitered,
		// TcpChannelMsg: make(map[string]chan string),
		// TcpConn:       make(map[string]net.Conn),
	}
}

func (ue *UeContext) AddPduSession(pduSessionId uint8, dnn string, snssai models.Snssai) *SessionContext {
	sess := &SessionContext{
		PduSessionId:      int64(pduSessionId),
		Dnn:               dnn,
		Snssai:            snssai,
		QosFlows:          make(map[int64]*QosFlow),
		Ue:                ue,
		SessTcpChannelMsg: make(chan string),
	}
	ue.PduSession[sess.PduSessionId] = sess
	return sess
}

func (s *SessionContext) Remove() {
	if ue := s.Ue; ue != nil {
		if ran := ue.Ran; ran != nil {
			ran.DetachSession(s)
		}
		delete(ue.PduSession, s.PduSessionId)
	}
	Simulator_Self().DetachSession(s)
}

func (s *SessionContext) SendMsg(msg string) {
	if s.SessTcpChannelMsg != nil {
		select {
		case s.SessTcpChannelMsg <- msg:
		default:
			logger.ContextLog.Warnf("Can't send Msg to Tcp client")
		}
	}
}

// func (s *SessionContext) GetGtpConn() (*net.UDPConn, error) {
// 	key := fmt.Sprintf("%s,%s", s.DLAddr, s.ULAddr)
// 	if conn := Simulator_Self().GtpConnPool[key]; conn != nil {
// 		return conn, nil
// 	} else {
// 		return nil, fmt.Errorf("gtp conn is empty, map key [%s]", key)
// 	}
// }

// func (s *SessionContext) NewGtpHeader(extHdrFlag, sqnFlag, numFlag byte) {
// 	extHdrFlag &= 0x1
// 	sqnFlag &= 0x1
// 	numFlag &= 0x1
// 	if extHdrFlag == 0 && sqnFlag == 0 && numFlag == 0 {
// 		s.GtpHdrLen = 8
// 	} else {
// 		s.GtpHdrLen = 12
// 	}
// 	s.GtpHdr = make([]byte, s.GtpHdrLen)
// 	// Version: 3-bit, gtpv1=1
// 	// Protocol type: 1-bit, GTP=1, GTP'=0
// 	// Reserved: 1-bit 0
// 	// E: 1-bit
// 	// S: 1-bit
// 	// PN: 1-bit
// 	s.GtpHdr[0] = 0x01<<5 | 0x01<<4 | extHdrFlag<<2 | sqnFlag<<1 | numFlag
// 	// Message Type: 8-bit reference to 3GPP TS 29.060 subclause 7.1
// 	s.GtpHdr[1] = 0xff
// 	// Total Length: 16-bit not include first 8 bits
// 	// Wait for realData
// 	// TEID: 32-bit
// 	binary.BigEndian.PutUint32(s.GtpHdr[4:8], s.ULTEID)
// 	// Sequence number: 32-bit (optinal, if D is true)
// 	// N-PDU number: 16-bit (optinal, if PN is true)
// 	// Next extension header type: 16-bit (optinal, if E is true)
// }

func (s *SessionContext) GetTunnelMsg() string {
	s.Mtx.Lock()
	if s.ULAddr == "" {
		return ""
	}
	msg := fmt.Sprintf("ID=%d,DNN=%s,SST=%d,SD=%s,UEIP=%s,ULAddr=%s,ULTEID=%d,DLAddr=%s,DLTEID=%d\n",
		s.PduSessionId, s.Dnn, s.Snssai.Sst, s.Snssai.Sd, s.UeIp, s.ULAddr, s.ULTEID, s.DLAddr, s.DLTEID)
	s.Mtx.Unlock()
	return msg
}

func (ue *UeContext) SendMsg(msg string) {
	// for _, channel := range ue.TcpChannelMsg {
	// 	select {
	// 	case channel <- msg:
	// 	default:
	// 		logger.ContextLog.Warnf("Can't send Msg to Tcp client")
	// 	}
	// }
}

func (ue *UeContext) AttachRan(ran *RanContext) {
	ue.Ran = ran
	ran.UePool[ran.RanUeIDGenerator] = ue
	ue.RanUeNgapId = ran.RanUeIDGenerator
	ran.RanUeIDGenerator++
}

func (ue *UeContext) DetachRan(ran *RanContext) {
	ue.Ran = nil
	delete(ran.UePool, ue.RanUeNgapId)
}

func (ue *UeContext) GetServingNetworkName() string {
	mcc := ue.ServingPlmnId[:3]
	mnc := ue.ServingPlmnId[3:]
	if len(mnc) == 2 {
		mnc = "0" + mnc
	}
	return fmt.Sprintf("5G:mnc%s.mcc%s.3gppnetwork.org", mnc, mcc)
}

func (ue *UeContext) DeriveRESstarAndSetKey(RAND []byte) []byte {
	authData := ue.AuthData
	snName := ue.GetServingNetworkName()
	SQN, _ := hex.DecodeString(authData.SQN)

	AMF, _ := hex.DecodeString(authData.AMF)

	// Run milenage
	MAC_A, MAC_S := make([]byte, 8), make([]byte, 8)
	CK, IK := make([]byte, 16), make([]byte, 16)
	RES := make([]byte, 8)
	AK, AKstar := make([]byte, 6), make([]byte, 6)
	OPC, _ := hex.DecodeString(authData.Opc)
	K, _ := hex.DecodeString(authData.K)
	// Generate MAC_A, MAC_S
	if err := milenage.F1(OPC, K, RAND, SQN, AMF, MAC_A, MAC_S); err != nil {
		logger.ContextLog.Errorln(err)
		return nil
	}

	// Generate RES, CK, IK, AK, AKstar
	if err := milenage.F2345(OPC, K, RAND, RES, CK, IK, AK, AKstar); err != nil {
		logger.ContextLog.Errorln(err)
		return nil
	}

	// derive RES*
	key := append(CK, IK...)
	FC := UeauCommon.FC_FOR_RES_STAR_XRES_STAR_DERIVATION
	P0 := []byte(snName)
	P1 := RAND
	P2 := RES

	ue.DerivateKamf(key, snName, SQN, AK)
	kdfVal_for_resStar := UeauCommon.GetKDFValue(key, FC, P0, UeauCommon.KDFLen(P0), P1, UeauCommon.KDFLen(P1), P2, UeauCommon.KDFLen(P2))
	return kdfVal_for_resStar[len(kdfVal_for_resStar)/2:]

}

func (ue *UeContext) DerivateKamf(key []byte, snName string, SQN, AK []byte) {

	FC := UeauCommon.FC_FOR_KAUSF_DERIVATION
	P0 := []byte(snName)
	SQNxorAK := make([]byte, 6)
	for i := 0; i < len(SQN); i++ {
		SQNxorAK[i] = SQN[i] ^ AK[i]
	}
	P1 := SQNxorAK
	Kausf := UeauCommon.GetKDFValue(key, FC, P0, UeauCommon.KDFLen(P0), P1, UeauCommon.KDFLen(P1))
	P0 = []byte(snName)
	Kseaf := UeauCommon.GetKDFValue(Kausf, UeauCommon.FC_FOR_KSEAF_DERIVATION, P0, UeauCommon.KDFLen(P0))

	supiRegexp, _ := regexp.Compile("(?:imsi|supi)-([0-9]{5,15})")
	groups := supiRegexp.FindStringSubmatch(ue.Supi)
	if groups == nil {
		return
	}
	P0 = []byte(groups[1])
	L0 := UeauCommon.KDFLen(P0)
	P1 = []byte{0x00, 0x00}
	L1 := UeauCommon.KDFLen(P1)

	ue.Kamf = UeauCommon.GetKDFValue(Kseaf, UeauCommon.FC_FOR_KAMF_DERIVATION, P0, L0, P1, L1)
}

// Algorithm key Derivation function defined in TS 33.501 Annex A.9
func (ue *UeContext) DerivateAlgKey() {
	// Security Key
	P0 := []byte{security.NNASEncAlg}
	L0 := UeauCommon.KDFLen(P0)
	P1 := []byte{ue.CipheringAlg}
	L1 := UeauCommon.KDFLen(P1)

	kenc := UeauCommon.GetKDFValue(ue.Kamf, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	copy(ue.KnasEnc[:], kenc[16:32])

	// Integrity Key
	P0 = []byte{security.NNASIntAlg}
	L0 = UeauCommon.KDFLen(P0)
	P1 = []byte{ue.IntegrityAlg}
	L1 = UeauCommon.KDFLen(P1)

	kint := UeauCommon.GetKDFValue(ue.Kamf, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	copy(ue.KnasInt[:], kint[16:32])
}
