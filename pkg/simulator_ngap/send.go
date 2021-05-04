package simulator_ngap

import (
	"fmt"
	"os/exec"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
	"github.com/jay16213/radio_simulator/pkg/simulator_nas/nas_packet"

	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/ngap/ngapType"
)

func (c *NGController) SendNGSetupRequest(endpoint *sctp.SCTPAddr) {
	logger.NgapLog.Info("Send NG Setup Request")
	// send NGSetupRequest Msg
	pkt, err := c.BuildNGSetupRequest()
	if err != nil {
		logger.NgapLog.Errorf("Build NGSetUp failed : %s", err.Error())
		return
	}
	c.ran.SendToAMF(endpoint, pkt)
}

func (c *NGController) SendInitailUeMessage_RegistraionRequest(endpoint *sctp.SCTPAddr, ue *simulator_context.UeContext) {
	logger.NgapLog.Infow("Send Initail UE Message (Initail Registration Request)", "rid", ue.RanUeNgapId)

	nasPdu, err := nas_packet.GetRegistrationRequestWith5GMM(ue, nasMessage.RegistrationType5GSInitialRegistration, nil, nil)
	if err != nil {
		logger.NgapLog.Errorf("Build RegistrationRequest failed: %s", err.Error())
		return
	}

	pkt, err := BuildInitialUEMessage(ue, "", nasPdu)
	if err != nil {
		logger.NgapLog.Errorf("Build InitialUEMessage failed: %s", err.Error())
		return
	}
	c.ran.SendToAMF(endpoint, pkt)
}

func (c *NGController) SendInitailUeMessage(endpoint *sctp.SCTPAddr, ue *simulator_context.UeContext, nasPdu []byte) {
	logger.NgapLog.Infow("Send Initail UE Message", "rid", ue.RanUeNgapId)
	pkt, err := BuildInitialUEMessage(ue, "", nasPdu)
	if err != nil {
		logger.NgapLog.Errorf("Build InitialUEMessage failed : %s", err.Error())
		return
	}
	c.ran.SendToAMF(endpoint, pkt)
}

func (c *NGController) SendUplinkNASTransport(endpoint *sctp.SCTPAddr, ue *simulator_context.UeContext, nasPdu []byte) {
	logger.NgapLog.Infow("Send Uplink NAS Transport", "supi", ue.Supi, "id", ue.AmfUeNgapId, "amf", endpoint.String())

	pkt, err := BuildUplinkNasTransport(ue, nasPdu)
	if err != nil {
		logger.NgapLog.Errorf("Build Uplink Nas Transport failed : %s", err.Error())
		return
	}
	c.ran.SendToAMF(endpoint, pkt)
}

func (c *NGController) SendIntialContextSetupResponse(endpoint *sctp.SCTPAddr, ue *simulator_context.UeContext, pduSessionIds []string) {
	logger.NgapLog.Infow("Send Intial Context Setup Response", "id", ue.AmfUeNgapId, "rid", ue.RanUeNgapId)

	pkt, err := BuildInitialContextSetupResponse(ue, pduSessionIds, nil)
	if err != nil {
		logger.NgapLog.Errorf("Build Uplink Nas Transport failed : %s", err.Error())
		return
	}
	c.ran.SendToAMF(endpoint, pkt)
}

func (c *NGController) SendUeContextReleaseComplete(endpoint *sctp.SCTPAddr, ue *simulator_context.UeContext) {
	logger.NgapLog.Info("Send Ue Context Release Complete")

	pkt, err := BuildUEContextReleaseComplete(ue)
	if err != nil {
		logger.NgapLog.Errorf("Build Ue Context Release Complete failed : %s", err.Error())
		return
	}
	c.ran.SendToAMF(endpoint, pkt)

	// Reset Ue Context
	c.CloseNASConnection(ue.RanUeNgapId)
	for _, sess := range ue.PduSession {
		sess.Remove()
		if sess.UeIp != "" {
			_, err := exec.Command("ip", "addr", "del", sess.UeIp, "dev", "lo").Output()
			if err != nil {
				logger.NgapLog.Error(err)
				return
			}
		}
	}
}

func (c *NGController) SendPDUSessionResourceSetupResponse(
	endpoint *sctp.SCTPAddr,
	ue *simulator_context.UeContext,
	responseList *ngapType.PDUSessionResourceSetupListSURes,
	failedListSURes *ngapType.PDUSessionResourceFailedToSetupListSURes) {

	logger.NgapLog.Info("Send PDU Session Resource Setup Response")

	pkt, err := BuildPDUSessionResourceSetupResponse(ue, responseList, failedListSURes)
	if err != nil {
		logger.NgapLog.Errorf("Build PDU Session Resource Setup Response failed : %+v", err)
		return
	}

	c.ran.SendToAMF(endpoint, pkt)
	// Send Callback To Tcp Client
	if responseList != nil {
		for _, item := range responseList.List {
			sess := ue.PduSession[item.PDUSessionID.Value]
			msg := "[SESSION] " + sess.GetTunnelMsg()
			sess.SendMsg(msg)
		}
	}
	if failedListSURes != nil {
		for _, item := range failedListSURes.List {
			sess := ue.PduSession[item.PDUSessionID.Value]
			if sess != nil {
				msg := fmt.Sprintf("[SESSION] ADD %d FAIL\n", item.PDUSessionID.Value)
				sess.SendMsg(msg)
				sess.Remove()
			}
		}
	}
}

func (c *NGController) SendPDUSessionResourceReleaseResponse(
	endpoint *sctp.SCTPAddr,
	ue *simulator_context.UeContext,
	relList ngapType.PDUSessionResourceReleasedListRelRes,
	diagnostics *ngapType.CriticalityDiagnostics) {

	logger.NgapLog.Info("Send PDU Session Resource Release Response")

	if len(relList.List) < 1 {
		logger.NgapLog.Error("PDUSessionResourceReleasedListRelRes is nil. This message shall contain at least one Item")
		return
	}

	pkt, err := BuildPDUSessionResourceReleaseResponse(ue, relList, diagnostics)
	if err != nil {
		logger.NgapLog.Errorf("Build PDU Session Resource Release Response failed : %+v", err)
		return
	}

	c.ran.SendToAMF(endpoint, pkt)
}

func (c *NGController) SendErrorIndication(endpoint *sctp.SCTPAddr, amfUeNgapId, ranUeNgapId *int64, cause *ngapType.Cause, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	logger.NgapLog.Info("Send Error Indication")

	pkt, err := BuildErrorIndication(amfUeNgapId, ranUeNgapId, cause, criticalityDiagnostics)
	if err != nil {
		logger.NgapLog.Errorf("Build ErrorIndication failed : %s", err.Error())
		return
	}
	c.ran.SendToAMF(endpoint, pkt)
}

func (c *NGController) SendAMFConfigurationUpdateAcknowledge(endpoint *sctp.SCTPAddr, setupList *ngapType.AMFTNLAssociationSetupList) {
	logger.NgapLog.Info("Send AMF Configuration Update Acknowledge")

	pkt, err := BuildAMFConfigurationUpdateAcknowledge(setupList)
	if err != nil {
		logger.NgapLog.Errorf("Build AMFConfigurationUpdateAcknowledge failed: %+v", err)
		return
	}
	c.ran.SendToAMF(endpoint, pkt)
}

func (c *NGController) SendRanConfigurationUpdate(endpoint *sctp.SCTPAddr) {
	logger.NgapLog.Info("Send RAN Configuration Update")

	pkt, err := c.BuildRanConfigurationUpdate()
	if err != nil {
		logger.NgapLog.Errorf("Build RanConfigurationUpdate failed : %s", err.Error())
		return
	}
	c.ran.SendToAMF(endpoint, pkt)
}
