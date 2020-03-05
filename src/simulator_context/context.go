package simulator_context

import (
	"encoding/hex"
	"net"
	"radio_simulator/lib/ngap/ngapType"
	"radio_simulator/lib/openapi/models"
)

var simContext = Simulator{}

func init() {
	Simulator_Self().RanPool = make(map[string]*RanContext)
	Simulator_Self().UeContextPool = make(map[string]*UeContext)
}

type Simulator struct {
	DefaultRanUri string
	RanPool       map[string]*RanContext // RanUri -> RAN_CONTEXT
	UeContextPool map[string]*UeContext  // Supi -> UeTestInfo
	TcpServer     net.Listener
	TcpConn       net.Conn
}

type UeDBInfo struct {
	AmDate     models.AccessAndMobilitySubscriptionData
	SmfSelData models.SmfSelectionSubscriptionData
	AmPolicy   models.AmPolicyData
	AuthsSubs  models.AuthenticationSubscription
	PlmnId     string
}

func (s *Simulator) AddRanContext(AmfUri, ranUri, ranName string, plmnId ngapType.PLMNIdentity, GnbId string, gnbIdLength int) *RanContext {
	ran := NewRanContext()
	ran.AMFUri = AmfUri
	ran.RanUri = ranUri
	ran.Name = ranName
	ran.GnbId.BitLength = uint64(gnbIdLength)
	ran.GnbId.Bytes, _ = hex.DecodeString(GnbId)
	s.RanPool[ranUri] = ran
	return ran
}

// Create new AMF context
func Simulator_Self() *Simulator {
	return &simContext
}
