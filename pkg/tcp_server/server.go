package tcp_server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
)

var self *simulator_context.Simulator = simulator_context.Simulator_Self()
var mtx sync.Mutex

func StartApiServer(addr string) net.Listener {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.TcpServerLog.Error(err.Error())
	}
	defer listener.Close()
	logger.TcpServerLog.Infof("TCP server start and listening on %s.", addr)

	// for {
	// 	conn, err := self.TcpServer.Accept()
	// 	if err != nil {
	// 		logger.TcpServerLog.Infof("TCP server closed")
	// 		return nil
	// 	}
	// 	raddr := conn.RemoteAddr().String()
	// 	go handleUeConnection(raddr, conn)
	// }
	return listener
}

func handleUeConnection(raddr string, conn net.Conn) {

	logger.TcpServerLog.Infof("Client connected from: " + raddr)
	conn.Write([]byte("Please Enter Supi:\n"))
	supi := new(string)
	// Make a buffer to hold incoming data.
	for {
		// Read the incoming connection into the buffer.
		err := Read(conn, raddr, supi)
		if err != nil {

			if err == io.EOF {
				logger.TcpServerLog.Infoln("Disconned from ", raddr)
				break
			} else {
				logger.TcpServerLog.Infoln("Error reading:", err.Error())
				break
			}
		}
	}
	// Close the connection when you're done with it.
	// if supi != nil {
	// ue := self.UeContextPool[*supi]
	// if ue != nil {
	// 	delete(ue.TcpConn, raddr)
	// }
	// }
	if conn != nil {
		conn.Close()
	}
}

func Read(conn net.Conn, raddr string, supi *string) error {
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
		if *supi != "" {
			msg := parseCmd(self.UeContextPool[*supi], raddr, cmd)
			if msg != "" {
				if msg[0] != '[' {
					msg = "[ERROR] " + msg + "\n"
				}
				conn.Write([]byte(msg))
			}
		} else if strings.HasPrefix(cmd, "imsi-") {
			ue := self.UeContextPool[cmd]
			if ue == nil {
				conn.Write([]byte("[ERROR] UE_NOT_EXIST\n"))
			} else {
				// mtx.Lock()
				// ue.TcpConn[raddr] = conn
				// mtx.Unlock()
				*supi = cmd
				conn.Write([]byte(fmt.Sprintf("Welcome User %s\n", *supi)))
			}
		} else {
			conn.Write([]byte("Please type Supi first\n"))
		}
		if isPrefix {
			break
		}
	}
	return nil
}
