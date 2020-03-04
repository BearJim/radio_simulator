package simulator_ngap

import (
	"radio_simulator/lib/ngap"
	"radio_simulator/lib/ngap/ngapType"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_context"
)

func check(err error) {
	if err != nil {
		logger.InitLog.Error(err.Error())
	}
}

func SendNGSetupRequest(ran *simulator_context.RanContext) {
	var n int
	var recvMsg = make([]byte, 2048)

	// send NGSetupRequest Msg
	sendMsg, err := BuildNGSetupRequest(ran)
	check(err)
	_, err = ran.SctpConn.Write(sendMsg)
	check(err)

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
