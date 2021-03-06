package simulator_context

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/BearJim/radio_simulator/pkg/api"
	"github.com/BearJim/radio_simulator/pkg/logger"

	"github.com/free5gc/UeauCommon"
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
	RmStateRegistered  = "REGISTERED"
	RmStateRegistering = "REGISTERING"
	RmStateDeregitered = "DEREGISTERED"
)

const (
	CmStateConnected = "CONNECTED"
	CmStateIdle      = "IDLE"
)

const (
	MsgRegisterSuccess       = "Registration success"
	MsgRegisterFail          = "Registration fail"
	MsgServiceRequestSuccess = "ServiceRequest success"
	MsgServiceRequestFail    = "ServiceRequest fail"
	MsgDeregisterSuccess     = "Deregistration success"
	MsgDeregisterFail        = "Deregistration fail"
)

type UeContext struct {
	AMFEndpoint *sctp.SCTPAddr

	// registration related
	FollowOnRequest bool
	ServingRan      string `bson:"servingRan"` // serving RAN name
	Supi            string `yaml:"supi" bson:"supi"`
	Guti            *nasType.GUTI5G
	GutiStr         string
	Gpsis           []string                            `yaml:"gpsis" bson:"gpsis"`
	Nssai           models.Nssai                        `yaml:"nssai" bson:"nssai"`
	UeAmbr          UeAmbr                              `yaml:"ueAmbr" bson:"ueAmbr"`
	SmfSelData      models.SmfSelectionSubscriptionData `yaml:"smfSelData" bson:"smfSelectionSubscriptionData"`
	AuthData        AuthData                            `yaml:"auths" bson:"authData"`
	SubscCats       []string                            `yaml:"subscCats,omitempty" bson:"subscCats"`
	ServingPlmnId   string                              `yaml:"servingPlmn" bson:"servingPlmn"`
	RanUeNgapId     int64
	AmfUeNgapId     int64
	// security
	ULCount         security.Count `bson:"nasUplinkCount"`
	DLCount         security.Count `bson:"nasDownlinkCount"`
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
	Ran     *RanContext `bson:"-"`
	RmState string      `bson:"rmState"`
	CmState string      `bson:"cmState"`
	// For API Usage
	RestartCount     int
	RestartTimeStamp time.Time
	ApiNotifyChan    chan ApiNotification `bson:"-"`
}

type ApiNotification struct {
	Status           api.StatusCode
	Message          string
	RestartCount     int
	RestartTimeStamp time.Time
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
	SQN        string `yaml:"SQN" bson:"SQN"` // 48-bit integer in hex format
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
		PduSession:    make(map[int64]*SessionContext),
		AmfUeNgapId:   AmfNgapIdUnspecified,
		RanUeNgapId:   RanNgapIdUnspecified,
		RmState:       RmStateDeregitered,
		CmState:       CmStateIdle,
		ApiNotifyChan: make(chan ApiNotification, 100),
	}
}

func (ue *UeContext) AuthDataSQNAddOne() {
	num, _ := strconv.ParseInt(ue.AuthData.SQN, 16, 48)
	ue.AuthData.SQN = fmt.Sprintf("%x", num+1)
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

func (ue *UeContext) SendAPINotification(status api.StatusCode, msg string) {
	ue.ApiNotifyChan <- ApiNotification{
		Status:           status,
		Message:          msg,
		RestartCount:     ue.RestartCount,
		RestartTimeStamp: ue.RestartTimeStamp,
	}
}

func (ue *UeContext) GetServingNetworkName() string {
	mcc := ue.ServingPlmnId[:3]
	mnc := ue.ServingPlmnId[3:]
	if len(mnc) == 2 {
		mnc = "0" + mnc
	}
	return fmt.Sprintf("5G:mnc%s.mcc%s.3gppnetwork.org", mnc, mcc)
}

// TS 33.501 Annex A.4
func (ue *UeContext) DeriveRESstar(ck []byte, ik []byte, servingNetworkName string, rand []byte, res []byte) []byte {
	inputKey := append(ck, ik...)
	FC := UeauCommon.FC_FOR_RES_STAR_XRES_STAR_DERIVATION
	P0 := []byte(servingNetworkName)
	L0 := UeauCommon.KDFLen(P0)
	P1 := rand
	L1 := UeauCommon.KDFLen(P1)
	P2 := res
	L2 := UeauCommon.KDFLen(P2)
	kdfVal_for_resStar := UeauCommon.GetKDFValue(inputKey, FC, P0, L0, P1, L1, P2, L2)
	return kdfVal_for_resStar[len(kdfVal_for_resStar)/2:]
}

// TS 33.501 Annex A.2
func DerivateKausf(ck []byte, ik []byte, servingNetworkName string, sqnXorAK []byte) []byte {
	inputKey := append(ck, ik...)
	P0 := []byte(servingNetworkName)
	L0 := UeauCommon.KDFLen(P0)
	P1 := sqnXorAK
	L1 := UeauCommon.KDFLen(P1)
	return UeauCommon.GetKDFValue(inputKey, UeauCommon.FC_FOR_KAUSF_DERIVATION, P0, L0, P1, L1)
}

func DerivateKseaf(kausf []byte, servingNetworkName string) []byte {
	P0 := []byte(servingNetworkName)
	L0 := UeauCommon.KDFLen(P0)
	return UeauCommon.GetKDFValue(kausf, UeauCommon.FC_FOR_KSEAF_DERIVATION, P0, L0)
}

func (ue *UeContext) DerivateKamf(kseaf []byte, abba []byte) {
	supiRegexp, _ := regexp.Compile("(?:imsi|supi)-([0-9]{5,15})")
	groups := supiRegexp.FindStringSubmatch(ue.Supi)
	if groups == nil {
		return
	}
	P0 := []byte(groups[1])
	L0 := UeauCommon.KDFLen(P0)
	P1 := abba
	L1 := UeauCommon.KDFLen(P1)

	ue.Kamf = UeauCommon.GetKDFValue(kseaf, UeauCommon.FC_FOR_KAMF_DERIVATION, P0, L0, P1, L1)
	logger.ContextLog.Debugf("Kamf: 0x%0x", ue.Kamf)
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
