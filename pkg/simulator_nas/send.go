package simulator_nas

import (
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
	"github.com/jay16213/radio_simulator/pkg/simulator_nas/nas_packet"
)

func (c *NASController) SendAuthenticationResponse(ue *simulator_context.UeContext, resStar []byte) {
	logger.NASLog.Infow("Send Authentication Response", "supi", ue.Supi, "id", ue.AmfUeNgapId)
	nasPdu := nas_packet.BuildAuthenticationResponse(resStar, "")
	c.ngMessager.SendUplinkNASTransport(ue.AMFEndpoint, ue, nasPdu)
}

func (c *NASController) SendAuthenticationFailure(ue *simulator_context.UeContext, cause uint8, authFailParams []byte) {
	logger.NASLog.Infow("Send Authentication Failure", "supi", ue.Supi, "id", ue.AmfUeNgapId)
	nasPdu := nas_packet.BuildAuthenticationFailure(cause, authFailParams)
	c.ngMessager.SendUplinkNASTransport(ue.AMFEndpoint, ue, nasPdu)
}

func (c *NASController) SendSecurityModeComplete(ue *simulator_context.UeContext, nasMsg []byte) {
	logger.NASLog.Infow("Send Security Mode Complete", "supi", ue.Supi, "id", ue.AmfUeNgapId)
	nasPdu, err := nas_packet.GetSecurityModeComplete(ue, nasMsg)
	if err != nil {
		logger.NASLog.Errorf("Build Security Mode Complete error: %+v", err)
		return
	}
	c.ngMessager.SendUplinkNASTransport(ue.AMFEndpoint, ue, nasPdu)
}
