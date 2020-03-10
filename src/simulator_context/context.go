package simulator_context

import (
	"encoding/hex"
	"golang.org/x/net/ipv4"
	"net"
	"radio_simulator/lib/ngap/ngapType"
	"radio_simulator/lib/openapi/models"
	"sync"
	"syscall"
)

var simContext = Simulator{}

func init() {
	Simulator_Self().RanPool = make(map[string]*RanContext)
	Simulator_Self().SessPool = make(map[string]*SessionContext)
	Simulator_Self().UeContextPool = make(map[string]*UeContext)
	Simulator_Self().GtpConnPool = make(map[string]*net.UDPConn)
}

type Simulator struct {
	DefaultRanSctpUri string
	RanPool           map[string]*RanContext     // RanSctpUri -> RAN_CONTEXT
	SessPool          map[string]*SessionContext // UeIp -> RAN_CONTEXT
	UeContextPool     map[string]*UeContext      // Supi -> UeTestInfo
	GtpConnPool       map[string]*net.UDPConn    // "ranGtpuri,upfUri" -> conn
	TcpServer         net.Listener
	TunMtx            sync.Mutex
	TunFd             int
	TunSockAddr       syscall.Sockaddr
	ListenRawConn     *ipv4.RawConn
}

type UeDBInfo struct {
	AmDate     models.AccessAndMobilitySubscriptionData
	SmfSelData models.SmfSelectionSubscriptionData
	AmPolicy   models.AmPolicyData
	AuthsSubs  models.AuthenticationSubscription
	PlmnId     string
}

func (s *Simulator) SendToTunDev(msg []byte) {
	s.TunMtx.Lock()
	syscall.Sendto(s.TunFd, msg, 0, s.TunSockAddr)
	s.TunMtx.Unlock()
}

func (s *Simulator) AddRanContext(AmfUri, ranSctpUri, ranName string, ranGtpUri AddrInfo, plmnId ngapType.PLMNIdentity, GnbId string, gnbIdLength int) *RanContext {
	ran := NewRanContext()
	ran.AMFUri = AmfUri
	ran.RanSctpUri = ranSctpUri
	ran.RanGtpUri = ranGtpUri
	ran.Name = ranName
	ran.GnbId.BitLength = uint64(gnbIdLength)
	ran.GnbId.Bytes, _ = hex.DecodeString(GnbId)
	s.RanPool[ranSctpUri] = ran
	return ran
}

// Create new AMF context
func Simulator_Self() *Simulator {
	return &simContext
}
