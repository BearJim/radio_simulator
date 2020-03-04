package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"net"
	"radio_simulator/lib/path_util"
	"radio_simulator/src/factory"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_context"
	"radio_simulator/src/simulator_handler"
	"radio_simulator/src/simulator_init"
	"radio_simulator/src/simulator_util"
)

var config string

func Initailize() {
	flag.StringVar(&config, "simcfg", path_util.ModulePath("radio_simulator/config/rancfg.conf"), "ran simulator config file")
	flag.Parse()

	factory.InitConfigFactory(config)
	config := factory.SimConfig
	if config.Logger.DebugLevel != "" {
		level, err := logrus.ParseLevel(config.Logger.DebugLevel)
		if err == nil {
			logger.SetLogLevel(level)
		}
	}
	logger.SetReportCaller(config.Logger.ReportCaller)

}

func main() {
	Initailize()
	self := simulator_context.Simulator_Self()
	simulator_util.ParseRanContext()
	for _, ran := range self.RanPool {
		simulator_init.RanStart(ran)
	}
	go simulator_handler.Handle()
	srvAddr := factory.SimConfig.TcpUri
	listener, err := net.Listen("tcp", srvAddr)
	if err != nil {
		logger.SimulatorLog.Error(err.Error())
	}
	defer listener.Close()
	logger.SimulatorLog.Infof("TCP server start and listening on %s.", srvAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.SimulatorLog.Errorf("Some connection error: %s", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	logger.SimulatorLog.Infof("Client connected from: " + remoteAddr)

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	for {
		// Read the incoming connection into the buffer.
		reqLen, err := conn.Read(buf)
		if err != nil {

			if err.Error() == "EOF" {
				logger.SimulatorLog.Infof("Disconned from ", remoteAddr)
				break
			} else {
				logger.SimulatorLog.Infof("Error reading:", err.Error())
				break
			}
		}
		// Start Client
		logger.SimulatorLog.Infof("len: %d, recv: %s\n", reqLen, string(buf[:reqLen]))
	}
	// Close the connection when you're done with it.
	conn.Close()
}
