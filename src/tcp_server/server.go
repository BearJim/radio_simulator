package tcp_server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"radio_simulator/src/factory"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_context"
	"strings"
)

var self *simulator_context.Simulator = simulator_context.Simulator_Self()

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
	supi := new(string)
	// Make a buffer to hold incoming data.
	for {
		// Read the incoming connection into the buffer.
		err := Read(conn, supi)
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

		if *supi != "" {
			parseCmd(self.UeContextPool[*supi], cmd)
		} else if strings.HasPrefix(cmd, "imsi-") {
			ue := self.UeContextPool[cmd]
			if ue == nil {
				conn.Write([]byte("[ERROR] UE_NOT_EXIST\n"))
			} else {
				ue.TcpConn = conn
				*supi = cmd
				conn.Write([]byte(fmt.Sprintf("Welcom User %s\n", *supi)))
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
