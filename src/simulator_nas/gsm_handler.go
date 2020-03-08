package simulator_nas

import (
	"fmt"
	"net"
	"radio_simulator/lib/nas/nasConvert"
	"radio_simulator/lib/nas/nasMessage"
	"radio_simulator/src/simulator_context"
)

func HandlePduSessionEstblishmentAccept(ue *simulator_context.UeContext, request *nasMessage.PDUSessionEstablishmentAccept) error {

	nasLog.Infof("UE[%s] Handle PDU Session Establishment Accept", ue.Supi)

	pduSessionId := int64(request.GetPDUSessionID())
	sess, exist := ue.PduSession[pduSessionId]
	if !exist {
		return fmt.Errorf("pduSessionId[%d] is not exist in UE", pduSessionId)
	}
	if request.DNN != nil {
		sess.Dnn = string(request.GetDNN())
	}
	if request.SNSSAI != nil {
		sess.Snssai = nasConvert.SnssaiToModels(request.SNSSAI)
	}
	if request.PDUAddress != nil {
		ipBytes := request.PDUAddress.GetPDUAddressInformation()
		fmt.Println(ipBytes)
		switch request.PDUAddress.GetPDUSessionTypeValue() {
		case nasMessage.PDUSessionTypeIPv4:
			sess.UeIp = net.IP(ipBytes[:4]).String()
		case nasMessage.PDUSessionTypeIPv6, nasMessage.PDUSessionTypeIPv4IPv6:
			return fmt.Errorf("Ipv6 is not support yet")
		}
	}
	return nil
}
