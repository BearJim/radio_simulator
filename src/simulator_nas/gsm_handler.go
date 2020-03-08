package simulator_nas

import (
	"fmt"
	"radio_simulator/lib/nas/nasConvert"
	"radio_simulator/lib/nas/nasMessage"
	"radio_simulator/src/simulator_context"
)

func HandlePduSessionEstblishmentAccept(ue *simulator_context.UeContext, request *nasMessage.PDUSessionEstablishmentAccept) error {

	nasLog.Infof("UE[%s] Handle PDU Session Establishment Accept", ue.Supi)

	pduSessionId := int64(request.GetPDUSessionID())
	sess, exist := ue.PduSession[pduSessionId]
	if !exist {
		fmt.Errorf("pduSessionId[%d] is not exist in UE", pduSessionId)
	}
	if request.DNN != nil {
		sess.Dnn = string(request.GetDNN())
	}
	if request.SNSSAI != nil {
		sess.Snssai = nasConvert.SnssaiToModels(request.SNSSAI)
	}
	return nil
}
