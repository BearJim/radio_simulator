package simulator_nas

import (
	"context"

	"github.com/BearJim/radio_simulator/pkg/logger"
	"github.com/BearJim/radio_simulator/pkg/simulator_context"
)

type NASRoutine struct {
	rid int // routine id, read-only

	*NASController

	// Routine Context
	ue       *simulator_context.UeContext
	nasPduCh chan []byte
}

func (r *NASRoutine) ID() int {
	return r.rid
}

func (r *NASRoutine) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.NASLog.Infow("Cancel NAS Connection", "supi", r.ue.Supi, "id", r.ue.AmfUeNgapId, "amf", r.ue.AMFEndpoint)
			r.ue.CmState = simulator_context.CmStateIdle
			close(r.nasPduCh)
			return
		case nasPdu, ok := <-r.nasPduCh:
			if !ok {
				logger.NASLog.Infow("Close NAS Connection", "supi", r.ue.Supi, "id", r.ue.AmfUeNgapId, "amf", r.ue.AMFEndpoint)
				r.ue.CmState = simulator_context.CmStateIdle
				r.ue.AmfUeNgapId = simulator_context.AmfNgapIdUnspecified
				return
			}
			r.handleGmmMessage(r.ue, nasPdu)
		}
	}
}
