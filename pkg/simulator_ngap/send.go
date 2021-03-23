package simulator_ngap

import (
	"fmt"
	"os/exec"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/jay16213/radio_simulator/pkg/api"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"

	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/ngap/ngapType"
)

func (c *NGController) SendNGSetupRequest(endpoint *sctp.SCTPAddr) {
	logger.NgapLog.Info("Send NG Setup Request")
	// send NGSetupRequest Msg
	pkt, err := BuildNGSetupRequest(c.ran.Context())
	if err != nil {
		logger.NgapLog.Errorf("Build NGSetUp failed : %s", err.Error())
		return
	}
	c.ran.SendToAMF(endpoint, pkt)
}

func (c *NGController) SendInitailUeMessage_RegistraionRequest(endpoint *sctp.SCTPAddr, ue *simulator_context.UeContext) {
	logger.NgapLog.Info("Send Initail Ue Message (Initail Registration Request)")
	pkt, err := BuildInitialUEMessage(ue, nasMessage.RegistrationType5GSInitialRegistration, "")
	if err != nil {
		logger.NgapLog.Errorf("Build InitialUEMessage failed : %s", err.Error())
		return
	}
	c.ran.SendToAMF(endpoint, pkt)
}

func (c *NGController) SendUplinkNasTransport(endpoint *sctp.SCTPAddr, ue *simulator_context.UeContext, nasPdu []byte) {

	logger.NgapLog.Info("Send Uplink NAS Transport")

	pkt, err := BuildUplinkNasTransport(ue, nasPdu)
	if err != nil {
		logger.NgapLog.Errorf("Build Uplink Nas Transport failed : %s", err.Error())
		return
	}
	c.ran.SendToAMF(endpoint, pkt)
}

func (c *NGController) SendIntialContextSetupResponse(endpoint *sctp.SCTPAddr, ue *simulator_context.UeContext, pduSessionIds []string) {

	logger.NgapLog.Info("[RAN] Send Intial Context Setup Response")

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

	// Reset Ue Context
	ue.AmfUeNgapId = simulator_context.AmfNgapIdUnspecified
	for _, sess := range ue.PduSession {
		sess.Remove()
		if sess.UeIp != "" {
			_, err := exec.Command("ip", "addr", "del", sess.UeIp, "dev", "lo").Output()
			if err != nil {
				logger.NgapLog.Errorln(err)
				return
			}
		}
	}

	c.ran.SendToAMF(endpoint, pkt)
	if ue.RmState == simulator_context.RegisterStateDeregitered {
		// Complete Deregistration
		ue.CmState = simulator_context.CmStateIdle
		ue.SendAPINotification(api.StatusCode_OK, simulator_context.MsgDeregisterSuccess)
	}
}

func (c *NGController) SendPDUSessionResourceSetupResponse(
	endpoint *sctp.SCTPAddr,
	ue *simulator_context.UeContext,
	responseList *ngapType.PDUSessionResourceSetupListSURes,
	failedListSURes *ngapType.PDUSessionResourceFailedToSetupListSURes) {

	logger.NgapLog.Infoln("Send PDU Session Resource Setup Response")

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

	logger.NgapLog.Infoln("Send PDU Session Resource Release Response")

	if len(relList.List) < 1 {
		logger.NgapLog.Errorln("PDUSessionResourceReleasedListRelRes is nil. This message shall contain at least one Item")
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
