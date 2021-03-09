package simulator_ngap

import (
	"fmt"
	"net"
	"os/exec"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"

	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
	"github.com/sirupsen/logrus"
)

var ngapLog *logrus.Entry

func init() {
	ngapLog = logger.NgapLog
}

func SendNGSetupRequest(ran *simulator_context.RanContext) {
	logger.NgapLog.Info("Send NG Setup Request")
	// send NGSetupRequest Msg
	pkt, err := BuildNGSetupRequest(ran)
	if err != nil {
		ngapLog.Errorf("Build NGSetUp failed : %s", err.Error())
		return
	}
	SendToAmf(ran, pkt)
}

func SendInitailUeMessage_RegistraionRequest(ran *simulator_context.RanContext, ue *simulator_context.UeContext) {
	ngapLog.Info("[RAN] Initail Ue Message (Initail Registration Request)")
	pkt, err := BuildInitialUEMessage(ue, nasMessage.RegistrationType5GSInitialRegistration, "")
	if err != nil {
		ngapLog.Errorf("Build InitialUEMessage failed : %s", err.Error())
		return
	}
	SendToAmf(ran, pkt)
}

func SendUplinkNasTransport(ran *simulator_context.RanContext, ue *simulator_context.UeContext, nasPdu []byte) {

	ngapLog.Info("[RAN] Send Uplink Nas Transport")

	pkt, err := BuildUplinkNasTransport(ue, nasPdu)
	if err != nil {
		ngapLog.Errorf("Build Uplink Nas Transport failed : %s", err.Error())
		return
	}
	SendToAmf(ran, pkt)
}

func SendIntialContextSetupResponse(ran *simulator_context.RanContext, ue *simulator_context.UeContext, pduSessionIds []string) {

	ngapLog.Info("[RAN] Send Intial Context Setup Response")

	pkt, err := BuildInitialContextSetupResponse(ue, pduSessionIds, nil)
	if err != nil {
		ngapLog.Errorf("Build Uplink Nas Transport failed : %s", err.Error())
		return
	}
	SendToAmf(ran, pkt)

}

func SendUeContextReleaseComplete(ran *simulator_context.RanContext, ue *simulator_context.UeContext) {

	ngapLog.Info("[RAN] Send Ue Context Release Complete")

	pkt, err := BuildUEContextReleaseComplete(ue)
	if err != nil {
		ngapLog.Errorf("Build Ue Context Release Complete failed : %s", err.Error())
		return
	}

	// Reset Ue Context
	ue.AmfUeNgapId = simulator_context.AmfNgapIdUnspecified
	for _, sess := range ue.PduSession {
		sess.Remove()
		if sess.UeIp != "" {
			_, err := exec.Command("ip", "addr", "del", sess.UeIp, "dev", "lo").Output()
			if err != nil {
				ngapLog.Errorln(err)
				return
			}
		}
	}

	SendToAmf(ran, pkt)
	if ue.RegisterState == simulator_context.RegisterStateDeregitered {
		// Complete Deregistration
		ue.SendMsg("[DEREG] SUCCESS\n")
	}
}

func SendPDUSessionResourceSetupResponse(
	ran *simulator_context.RanContext,
	ue *simulator_context.UeContext,
	responseList *ngapType.PDUSessionResourceSetupListSURes,
	failedListSURes *ngapType.PDUSessionResourceFailedToSetupListSURes) {

	ngapLog.Infoln("[RAN] Send PDU Session Resource Setup Response")

	pkt, err := BuildPDUSessionResourceSetupResponse(ue, responseList, failedListSURes)
	if err != nil {
		ngapLog.Errorf("Build PDU Session Resource Setup Response failed : %+v", err)
		return
	}

	SendToAmf(ran, pkt)
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

func SendPDUSessionResourceReleaseResponse(
	ran *simulator_context.RanContext,
	ue *simulator_context.UeContext,
	relList ngapType.PDUSessionResourceReleasedListRelRes,
	diagnostics *ngapType.CriticalityDiagnostics) {

	ngapLog.Infoln("[EAN] Send PDU Session Resource Release Response")

	if len(relList.List) < 1 {
		ngapLog.Errorln("PDUSessionResourceReleasedListRelRes is nil. This message shall contain at least one Item")
		return
	}

	pkt, err := BuildPDUSessionResourceReleaseResponse(ue, relList, diagnostics)
	if err != nil {
		ngapLog.Errorf("Build PDU Session Resource Release Response failed : %+v", err)
		return
	}

	SendToAmf(ran, pkt)
}

func SendErrorIndication(ran *simulator_context.RanContext, amfUeNgapId, ranUeNgapId *int64, cause *ngapType.Cause, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	ngapLog.Info("[AMF] Send Error Indication")

	pkt, err := BuildErrorIndication(amfUeNgapId, ranUeNgapId, cause, criticalityDiagnostics)
	if err != nil {
		ngapLog.Errorf("Build ErrorIndication failed : %s", err.Error())
		return
	}
	SendToAmf(ran, pkt)
}

func SendToAmf(ran *simulator_context.RanContext, message []byte) {
	// TODO: complete one to many interface support
	logger.NgapLog.Warnf("ParseIP: %+v", net.ParseIP("127.0.0.1").To4())
	_, err := ran.SctpConn.SCTPWriteTo(message,
		&sctp.SndRcvInfo{
			PPID: ngap.PPID,
		},
		sctp.SCTPEndpoint{
			IPAddr: net.IPAddr{
				IP: net.ParseIP("127.0.0.1"),
			},
			Port: 38412,
		},
	)
	if err != nil {
		logger.InitLog.Error(err)
	}
}
