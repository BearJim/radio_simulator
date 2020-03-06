package simulator_context

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"radio_simulator/lib/UeauCommon"
	"radio_simulator/lib/milenage"
	"radio_simulator/lib/openapi/models"
	"sync"
)

// TS 33501 Annex A.8 Algorithm distinguisher For Knas_int Knas_enc
const (
	N_NAS_ENC_ALG uint8 = 0x01
	N_NAS_INT_ALG uint8 = 0x02
	N_RRC_ENC_ALG uint8 = 0x03
	N_RRC_INT_ALG uint8 = 0x04
	N_UP_ENC_alg  uint8 = 0x05
	N_UP_INT_alg  uint8 = 0x06
)

// TS 33501 Annex D Algorithm identifier values For Knas_int
const (
	ALG_INTEGRITY_128_NIA0 uint8 = 0x00 // NULL
	ALG_INTEGRITY_128_NIA1 uint8 = 0x01 // 128-Snow3G
	ALG_INTEGRITY_128_NIA2 uint8 = 0x02 // 128-AES
	ALG_INTEGRITY_128_NIA3 uint8 = 0x03 // 128-ZUC
)

// TS 33501 Annex D Algorithm identifier values For Knas_enc
const (
	ALG_CIPHERING_128_NEA0 uint8 = 0x00 // NULL
	ALG_CIPHERING_128_NEA1 uint8 = 0x01 // 128-Snow3G
	ALG_CIPHERING_128_NEA2 uint8 = 0x02 // 128-AES
	ALG_CIPHERING_128_NEA3 uint8 = 0x03 // 128-ZUC
)

// 1bit
const (
	SECURITY_DIRECTION_UPLINK   uint8 = 0x00
	SECURITY_DIRECTION_DOWNLINK uint8 = 0x01
)

// 5bits
const (
	SECURITY_ONLY_ONE_BEARER uint8 = 0x00
	SECURITY_BEARER_3GPP     uint8 = 0x01
	SECURITY_BEARER_NON_3GPP uint8 = 0x02
)

// TS 33501 Annex A.0 Access type distinguisher For Kgnb Kn3iwf
const (
	ACCESS_TYPE_3GPP     uint8 = 0x01
	ACCESS_TYPE_NON_3GPP uint8 = 0x02
)

const (
	NgapIdUnspecified int64 = 0xffffffffff
)
const (
	RegisterStateRegitered   = "REGISTERED"
	RegisterStateRegitering  = "REGISTERING"
	RegisterStateDeregitered = "DEREGISTERED"
)

type UeContext struct {
	Supi          string                              `yaml:"supi"`
	Gpsis         []string                            `yaml:"gpsis"`
	Nssai         models.Nssai                        `yaml:"nssai"`
	UeAmbr        UeAmbr                              `yaml:"ueAmbr"`
	SmfSelData    models.SmfSelectionSubscriptionData `yaml:"smfSelData"`
	AuthData      AuthData                            `yaml:"auths"`
	SubscCats     []string                            `yaml:"subscCats,omitempty"`
	ServingPlmnId string                              `yaml:"servingPlmn"`
	RanUeNgapId   int64
	AmfUeNgapId   int64
	// security
	ULCount          uint32
	DLOverflow       uint16
	DLCountSQN       uint8
	CipheringAlgOrig string `yaml:"cipherAlg"`
	IntegrityAlgOrig string `yaml:"integrityAlg"`
	EncAlg           uint8
	IntAlg           uint8
	KnasEnc          []uint8
	KnasInt          []uint8
	Kamf             []uint8
	// PduSession
	PduSession map[int]*SessionContext
	// related Context
	Ran           *RanContext
	RegisterState string
	// For TCP Client
	TcpConn net.Conn // supi -> UeTcpClient
}

type UeAmbr struct {
	Upink    string `yaml:"uplink"`
	DownLink string `yaml:"downlink"`
}

type AuthData struct {
	AuthMethod string `yaml:"authMethod"`
	K          string `yaml:"K"`
	Opc        string `yaml:"Opc,omitempty"`
	Op         string `yaml:"Op,omitempty"`
	AMF        string `yaml:"AMF"`
	SQN        string `yaml:"SQN"`
}

type SessionContext struct {
	Mtx          sync.Mutex
	PduSessionId int64
	ULTEID       int
	DLTEID       int
	UpfUri       string
	RanUri       string
}

func NewUeContext() *UeContext {
	return &UeContext{
		PduSession:    make(map[int]*SessionContext),
		RanUeNgapId:   NgapIdUnspecified,
		RegisterState: RegisterStateDeregitered,
	}
}

func (s *SessionContext) GetTunnelMsg() string {
	s.Mtx.Lock()
	if s.UpfUri == "" {
		return ""
	}
	msg := fmt.Sprintf("ID=%d,ULIP=%s,ULTEID=%d,DLIP=%s,DLTEID=%d", s.PduSessionId, s.UpfUri, s.ULTEID, s.RanUri, s.DLTEID)
	s.Mtx.Unlock()
	return msg
}
func (ue *UeContext) SendSuccessRegister() {
	if ue.TcpConn != nil {
		ue.TcpConn.Write([]byte("[REG] SUCCESS\n"))
	}
}
func (ue *UeContext) SendFailRegister() {
	if ue.TcpConn != nil {
		ue.TcpConn.Write([]byte("[REG] FAIL\n"))
	}
}

func (ue *UeContext) SendSuccessDeregister() {
	if ue.TcpConn != nil {
		ue.TcpConn.Write([]byte("[DEREG] SUCCESS\n"))
	}
}

func (ue *UeContext) AttachRan(ran *RanContext) {
	ue.Ran = ran
	ran.UePool[ran.RanUeIDGeneator] = ue
	ran.RanUeIDGeneator++
}

func (ue *UeContext) DetachRan(ran *RanContext) {
	ue.Ran = nil
	delete(ran.UePool, ue.RanUeNgapId)
}

func (ue *UeContext) GetSecurityULCount() []byte {
	var r = make([]byte, 4)
	binary.BigEndian.PutUint32(r, ue.ULCount&0xffffff)
	return r
}

func (ue *UeContext) GetSecurityDLCount() []byte {
	var r = make([]byte, 4)
	binary.BigEndian.PutUint16(r, ue.DLOverflow)
	r[3] = ue.DLCountSQN
	r[2] = r[1]
	r[1] = r[0]
	r[0] = 0x00
	return r
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
	milenage.F1_Test(OPC, K, RAND, SQN, AMF, MAC_A, MAC_S)

	// Generate RES, CK, IK, AK, AKstar
	milenage.F2345_Test(OPC, K, RAND, RES, CK, IK, AK, AKstar)

	// derive RES*
	key := append(CK, IK...)
	FC := UeauCommon.FC_FOR_RES_STAR_XRES_STAR_DERIVATION
	P0 := []byte(snName)
	P1 := RAND
	P2 := RES

	ue.DerivateKamf(key, snName, SQN, AK)
	ue.DerivateAlgKey()
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

	P0, _ = hex.DecodeString(ue.Supi)
	L0 := UeauCommon.KDFLen(P0)
	P1 = []byte{0x00, 0x00}
	L1 := UeauCommon.KDFLen(P1)

	ue.Kamf = UeauCommon.GetKDFValue(Kseaf, UeauCommon.FC_FOR_KAMF_DERIVATION, P0, L0, P1, L1)
}

// Algorithm key Derivation function defined in TS 33.501 Annex A.9
func (ue *UeContext) DerivateAlgKey() {
	// Security Key
	P0 := []byte{N_NAS_ENC_ALG}
	L0 := UeauCommon.KDFLen(P0)
	P1 := []byte{ue.EncAlg}
	L1 := UeauCommon.KDFLen(P1)

	kenc := UeauCommon.GetKDFValue(ue.Kamf, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	ue.KnasEnc = kenc[16:32]

	// Integrity Key
	P0 = []byte{N_NAS_INT_ALG}
	L0 = UeauCommon.KDFLen(P0)
	P1 = []byte{ue.IntAlg}
	L1 = UeauCommon.KDFLen(P1)

	kint := UeauCommon.GetKDFValue(ue.Kamf, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	ue.KnasInt = kint[16:32]
}
