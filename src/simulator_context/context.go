package simulator_context

import (
	"encoding/hex"
	"radio_simulator/lib/ngap/ngapType"
)

var simContext = Simulator{}

func init() {
	Simulator_Self().RanPool = make(map[string]*RanContext)
}

type Simulator struct {
	RanPool map[string]*RanContext // RanUri -> RAN_CONTEXT
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
