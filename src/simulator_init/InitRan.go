package simulator_init

import (
	"fmt"
	"net"
	"radio_simulator/lib/ngap/ngapSctp"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_context"
	"radio_simulator/src/simulator_handler"
	"radio_simulator/src/simulator_handler/simulator_message"
	"radio_simulator/src/simulator_ngap"
	"strings"

	"git.cs.nctu.edu.tw/calee/sctp"
)

func check(err error) {
	if err != nil {
		logger.InitLog.Error(err.Error())
	}
}
func getNgapIp(amfIP, ranIP string, amfPort, ranPort int) (amfAddr, ranAddr *sctp.SCTPAddr, err error) {
	ips := []net.IPAddr{}
	if ip, err1 := net.ResolveIPAddr("ip", amfIP); err1 != nil {
		err = fmt.Errorf("Error resolving address '%s': %v", amfIP, err1)
		return
	} else {
		ips = append(ips, *ip)
	}
	amfAddr = &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    amfPort,
	}
	ips = []net.IPAddr{}
	if ip, err1 := net.ResolveIPAddr("ip", ranIP); err1 != nil {
		err = fmt.Errorf("Error resolving address '%s': %v", ranIP, err1)
		return
	} else {
		ips = append(ips, *ip)
	}
	ranAddr = &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    ranPort,
	}
	return
}

func ConntectToAmf(amfIP, ranIP string, amfPort, ranPort int) (*sctp.SCTPConn, error) {
	amfAddr, ranAddr, err := getNgapIp(amfIP, ranIP, amfPort, ranPort)
	if err != nil {
		return nil, err
	}
	conn, err := sctp.DialSCTP("sctp", ranAddr, amfAddr)
	if err != nil {
		return nil, err
	}
	info, _ := conn.GetDefaultSentParam()
	info.PPID = ngapSctp.NGAP_PPID
	err = conn.SetDefaultSentParam(info)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func RanStart(ran *simulator_context.RanContext) {
	var amfPort, ranPort int
	amfAddr := strings.Split(ran.AMFUri, ":")
	ranAddr := strings.Split(ran.RanUri, ":")
	amfIp, ranIp := amfAddr[0], ranAddr[0]
	fmt.Sscanf(amfAddr[1], "%d", &amfPort)
	fmt.Sscanf(ranAddr[1], "%d", &ranPort)

	// RAN connect to AMF
	conn, err := ConntectToAmf(amfIp, ranIp, amfPort, ranPort)
	check(err)
	ran.SctpConn = conn
	simulator_ngap.SendNGSetupRequest(ran)
	// New NGAP Channel
	simulator_message.AddNgapChannel(ran.RanUri)
	// Listen NGAP Channel
	err = conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
	if err != nil {
		logger.NgapLog.Errorf("Failed to subscribe SCTP Event: %v", err)
	}
	go simulator_handler.Handle(ran.RanUri)
	go StartHandle(ran)
}

func StartHandle(ran *simulator_context.RanContext) {
	defer ran.SctpConn.Close()
	for {
		buffer := make([]byte, 8192)
		n, info, err := ran.SctpConn.SCTPRead(buffer)
		if err != nil {
			logger.NgapLog.Debugf("Error %v", err)
			delete(simulator_context.Simulator_Self().RanPool, ran.RanUri)
			simulator_message.DelNgapChannel(ran.RanUri)
			break
		} else if info == nil || info.PPID != ngapSctp.NGAP_PPID {
			logger.NgapLog.Warnf("Recv SCTP PPID != 60")
			continue
		}
		simulator_message.SendMessage(ran.RanUri, simulator_message.NGAPMessage{Value: buffer[:n]})
	}
}
