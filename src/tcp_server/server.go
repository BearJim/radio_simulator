package tcp_server

import (
	"bufio"
	"io"
	"net"
	"radio_simulator/src/factory"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_context"
	"radio_simulator/src/simulator_ngap"
	"regexp"
	"strconv"
	"strings"
)

var self *simulator_context.Simulator = simulator_context.Simulator_Self()
var stringFormat = regexp.MustCompile(`\S`)

func StartTcpServer() {
	var err error
	srvAddr := factory.SimConfig.TcpUri
	self.TcpServer, err = net.Listen("tcp", srvAddr)
	if err != nil {
		logger.TcpServerLog.Error(err.Error())
	}
	defer self.TcpServer.Close()
	logger.TcpServerLog.Infof("TCP server start and listening on %s.", srvAddr)

	for {
		conn, err := self.TcpServer.Accept()
		if err != nil {
			logger.TcpServerLog.Infof("TCP server closed")
			return
		}
		raddr := conn.RemoteAddr().String()
		go handleUeConnection(raddr, conn)
	}
}

func handleUeConnection(raddr string, conn net.Conn) {

	logger.TcpServerLog.Infof("Client connected from: " + raddr)
	conn.Write([]byte("Please Enter Supi:"))
	var supi *string
	// Make a buffer to hold incoming data.
	for {
		// Read the incoming connection into the buffer.
		err := Read(conn, supi)
		if err != nil {

			if err == io.EOF {
				logger.TcpServerLog.Infof("Disconned from ", raddr)
				break
			} else {
				logger.TcpServerLog.Infof("Error reading:", err.Error())
				break
			}
		}
	}
	// Close the connection when you're done with it.
	if supi != nil {
		ue := self.UeContextPool[*supi]
		if ue != nil {
			ue.TcpConn = nil
		}
	}
	if conn != nil {
		conn.Close()
	}
}

func Read(conn net.Conn, supi *string) error {
	reader := bufio.NewReader(conn)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		cmd := string(line)
		if supi != nil {
			parseCmd(self.UeContextPool[cmd], cmd)
		} else if strings.HasPrefix(cmd, "imsi-") {
			ue := self.UeContextPool[cmd]
			if ue == nil {
				conn.Write([]byte("[ERROR] UE_NOT_EXIST\n"))
			} else {
				ue.TcpConn = conn
				supi = &cmd
			}
		} else {
			conn.Write([]byte("Please type Supi first\n"))
		}
		if !isPrefix {
			break
		}
	}
	return nil
}

// Parse cli Command from Tcp Server
// Usage: [ options ]
//
// options:
//   show [all|${PduSessionId}]
// 		show ue current state
//
//   reg {ip}
// 		register to CoreNetwork
// 		ip is specific Ran Ip to Connect(default={firstRanIp})
//
// 	 dereg
// 		deregister to CoreNetwork
//
//   sess i [add|del]
// 		pduSessionId i add or delete
//
// Output Format
//   show  [all|${PduSessionId}]
//		first line in all case:
//			"[SHOW] REGISTERED\n" or "[SHOW] REGISTERING\n" or "[SHOW] DEREGISTERED\n"
//		sessInfo:
//			"[SHOW] ID=%d,ULIP=%s,ULTEID=%d,DLIP=%s,DLTEID=%d\n"
// 		all means show all Pdu Session Id
//
//   reg {ip}
// 		"[REG] SUCCESS\n" or
// 		"[REG] FAIL\n" or
//
// 	 dereg
// 		"[DEREG] SUCCESS\n" or
// 		"[DEREG] FAIL\n" or
//
//   sess i [add|del]
//		"[SESSION] ID=%d,ULIP=%s,ULTEID=%d,DLIP=%s,DLTEID=%d\n" for add case, "[SESSION] DEL %d\n" for del case or
// 		"[SESSION] ADD/DEL %d FAIL\n"
//
func parseCmd(ue *simulator_context.UeContext, cmd string) {
	params := stringFormat.FindAllString(cmd, -1)
	cnt := len(params)
	if cnt == 0 {
		return
	}
	var msg string
	switch params[0] {
	case "show":
		if cnt == 1 {
			msg = "show missing action[all/{id}]"
		} else {
			switch params[1] {
			case "all":
				msg = "[SHOW] " + ue.RegisterState + "\n"
				for _, sess := range ue.PduSession {
					sessInfo := sess.GetTunnelMsg()
					if sessInfo == "" {
						continue
					}
					msg = msg + "[SHOW] " + sessInfo + "\n"
				}
			default:
				id, err := strconv.Atoi(params[1])
				if err != nil {
					msg = "sess id is not digit"
					break
				}
				sess := ue.PduSession[id]
				if sess == nil {
					msg = "sess " + params[1] + " has not established yet"
					break
				}
				sessInfo := sess.GetTunnelMsg()
				if sessInfo == "" {
					msg = "sess " + params[1] + " is still establishing"
					break
				}
				msg = "[SHOW] " + sessInfo + "\n"
			}
		}
	case "reg":
		switch ue.RegisterState {
		case simulator_context.RegisterStateRegitered:
			ue.SendSuccessRegister()
			return
		case simulator_context.RegisterStateRegitering:
			return
		}
		ran := self.RanPool[self.DefaultRanUri]
		if cnt > 1 {
			ran = self.RanPool[params[1]]
			if ran == nil {
				msg = "ranIp " + params[1] + " does not exist"
				break
			}
			// Use Default RanUri
		}
		ue.AttachRan(ran)
		ue.RegisterState = simulator_context.RegisterStateRegitering
		simulator_ngap.SendInitailUeMessage_RegistraionRequest(ran, ue)
	case "dereg":
		if ue.RegisterState == simulator_context.RegisterStateDeregitered {
			ue.SendSuccessDeregister()
			return
		} else {
			// TODO: Send Degister Request
		}
	case "sess":
		if cnt <= 2 {
			msg = "sess need id and action[add/del]"
			break
		}
		id, err := strconv.Atoi(params[1])
		if err != nil {
			msg = "sess id is not digit"
			break
		}
		sess := ue.PduSession[id]
		switch params[2] {
		case "add":
			if sess == nil {
				// TODO: Send Pdu Session Estblishment
				break
			}
			sessInfo := sess.GetTunnelMsg()
			if sessInfo == "" {
				msg = "sess " + params[1] + " is still establishing"
				break
			}
			msg = "[SESSION] " + sessInfo + "\n"
		case "del":
			if sess == nil {
				msg = "[SESSION] DEL " + params[1] + "\n"
				break
			} else {
				// TODO: Send Pdu Session Release
			}
		default:
			msg = "sess action is not [add/del]"
		}
	}

	if msg != "" {
		if msg[0] != '[' {
			msg = "[ERROR] " + msg + "\n"
		}
		ue.TcpConn.Write([]byte(msg))
	}

}
