package simulator_init

import (
	"fmt"
	"net"
	"strings"

	"github.com/free5gc/ngap"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
	"github.com/jay16213/radio_simulator/pkg/simulator_handler"
	"github.com/jay16213/radio_simulator/pkg/simulator_ngap"

	"git.cs.nctu.edu.tw/calee/sctp"
)

func getNgapIp(amfIP, ranIP string, amfPort, ranPort int) (amfAddr, ranAddr *sctp.SCTPAddr, err error) {
	ips := []net.IPAddr{}
	if ip, err1 := net.ResolveIPAddr("ip", "127.0.0.1"); err1 != nil {
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

func ConnectToAmf(amfIP, ranIP string, amfPort, ranPort int) (*sctp.SCTPConn, error) {
	amfAddr, ranAddr, err := getNgapIp(amfIP, ranIP, amfPort, ranPort)
	if err != nil {
		return nil, err
	}
	conn, err := sctp.DialSCTPOneToMany("sctp", ranAddr, amfAddr)
	if err != nil {
		return nil, err
	}
	info, _ := conn.GetDefaultSentParam()
	info.PPID = ngap.PPID
	err = conn.SetDefaultSentParam(info)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func connectToUpf(enbIP, upfIP string, gnbPort, upfPort int) (*net.UDPConn, error) {
	upfAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", upfIP, upfPort))
	if err != nil {
		return nil, err
	}
	gnbAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", enbIP, gnbPort))
	if err != nil {
		return nil, err
	}
	return net.DialUDP("udp", gnbAddr, upfAddr)
}

func RanStart(ran *simulator_context.RanContext) {
	var amfPort, ranPort int
	amfAddr := strings.Split(ran.AMFUri, ":")
	ranAddr := strings.Split(ran.RanSctpUri, ":")
	amfIp, ranIp := amfAddr[0], ranAddr[0]
	fmt.Sscanf(amfAddr[1], "%d", &amfPort)
	fmt.Sscanf(ranAddr[1], "%d", &ranPort)

	var err error
	// RAN connect to UPF
	// for _, upf := range ran.UpfInfoList {
	// upf.GtpConn, err = connectToUpf(ran.RanGtpUri.IP, upf.Addr.IP, ran.RanGtpUri.Port, upf.Addr.Port)
	// check(err)
	// simulator_context.Simulator_Self().GtpConnPool[fmt.Sprintf("%s,%s", ran.RanGtpUri.IP, upf.Addr.IP)] = upf.GtpConn
	// go StartHandleGtp(upf)
	// }
	// RAN connect to AMF
	conn, err := ConnectToAmf(amfIp, ranIp, amfPort, ranPort)
	if err != nil {
		logger.InitLog.Error(err.Error())
		return
	}
	ran.SctpConn = conn

	// Listen NGAP Channel
	err = conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
	if err != nil {
		logger.NgapLog.Errorf("Failed to subscribe SCTP Event: %v", err)
	}
	msgChan := make(chan []byte, 1024)
	go simulator_handler.Handle(ran, msgChan)
	go StartSCTPAssociation(ran.SctpConn, msgChan)
	simulator_ngap.SendNGSetupRequest(ran)
}

func StartSCTPAssociation(conn *sctp.SCTPConn, msgChan chan []byte) {
	defer conn.Close()
	for {
		buffer := make([]byte, 8192)
		n, info, _, err := conn.SCTPRead(buffer)
		if err != nil {
			logger.NgapLog.Debugf("Read Error: %v", err)
			break
		} else if info == nil || info.PPID != ngap.PPID {
			logger.NgapLog.Warnf("Recv SCTP PPID != 60")
			continue
		}
		msgChan <- buffer[:n]
	}
}

// func StartHandleGtp(upf *simulator_context.UpfInfo) {
// 	defer upf.GtpConn.Close()
// 	buffer := make([]byte, 8192)
// 	for {
// 		n, err := upf.GtpConn.Read(buffer)
// 		if err != nil {
// 			logger.GtpLog.Debugf("Error %v", err)
// 			break
// 		}
// 		msg := buffer[8:n] // remove gtp header
// 		// if msg[0] != 0x45 {
// 		// 	msg = msg[4:]
// 		// }
// 		simulator_context.Simulator_Self().SendToTunDev(msg)
// 	}
// }
