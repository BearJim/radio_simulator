package simulator_ngap

import (
	"github.com/sirupsen/logrus"
	"radio_simulator/lib/nas/nasMessage"
	"radio_simulator/lib/ngap"
	"radio_simulator/lib/ngap/ngapType"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_context"
)

var ngapLog *logrus.Entry

func init() {
	ngapLog = logger.NgapLog
}

func check(err error) {
	if err != nil {
		logger.InitLog.Error(err.Error())
	}
}

func SendNGSetupRequest(ran *simulator_context.RanContext) {
	var n int
	var recvMsg = make([]byte, 2048)

	// send NGSetupRequest Msg
	pkt, err := BuildNGSetupRequest(ran)
	if err != nil {
		ngapLog.Errorf("Build NGSetUp failed : %s", err.Error())
		return
	}
	SendToAmf(ran, pkt)

	// receive NGSetupResponse Msg
	n, err = ran.SctpConn.Read(recvMsg)
	check(err)
	pdu, err := ngap.Decoder(recvMsg[:n])
	check(err)
	if pdu.Present != ngapType.NGAPPDUPresentSuccessfulOutcome {
		logger.NgapLog.Error("NG SetUp Fail!!!!")
	}
	return
}

func SendInitailUeMessage_RegistraionRequest(ran *simulator_context.RanContext, ue *simulator_context.UeContext) {
	ngapLog.Info("[AMF] Initail Ue Message (Initail Registration Request)")
	pkt, err := BuildInitialUEMessage(ue, nasMessage.RegistrationType5GSInitialRegistration, "")
	if err != nil {
		ngapLog.Errorf("Build InitialUEMessage failed : %s", err.Error())
		return
	}
	SendToAmf(ran, pkt)
}

func SendUplinkNasTransport(ran *simulator_context.RanContext, ue *simulator_context.UeContext, nasPdu []byte) {

	ngapLog.Info("[AMF] Send Uplink Nas Transport")

	pkt, err := BuildUplinkNasTransport(ue, nasPdu)
	if err != nil {
		ngapLog.Errorf("Build Uplink Nas Transport failed : %s", err.Error())
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
	_, err := ran.SctpConn.Write(message)
	check(err)
}
