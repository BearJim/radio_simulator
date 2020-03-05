package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"radio_simulator/lib/MongoDBLibrary"
	"radio_simulator/lib/path_util"
	"radio_simulator/src/factory"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_context"
	"radio_simulator/src/simulator_init"
	"radio_simulator/src/simulator_util"
	"syscall"
)

var config string

var self *simulator_context.Simulator = simulator_context.Simulator_Self()

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

	MongoDBLibrary.SetMongoDB(config.DBName, config.DBUrl)
}

func Terminate() {
	logger.InitLog.Infof("Terminating Simulator...")

	// TODO: Send UE Deregistration to AMF
	logger.InitLog.Infof("Clear UE DB...")

	simulator_util.ClearDB()

	logger.InitLog.Infof("Close SCTP Connection...")

	for _, ran := range self.RanPool {
		logger.InitLog.Infof("Ran[%s] Connection Close", ran.RanUri)
		ran.SctpConn.Close()
	}

	logger.InitLog.Infof("Close TCP Connection...")
	if self.TcpConn != nil {
		self.TcpConn.Close()
	}
	if self.TcpServer != nil {
		self.TcpServer.Close()
	}

	logger.InitLog.Infof("Simulator terminated")

}

func startTcpServer() {
	var err error
	srvAddr := factory.SimConfig.TcpUri
	self.TcpServer, err = net.Listen("tcp", srvAddr)
	if err != nil {
		logger.SimulatorLog.Error(err.Error())
	}
	defer self.TcpServer.Close()
	logger.SimulatorLog.Infof("TCP server start and listening on %s.", srvAddr)

	for {
		self.TcpConn, err = self.TcpServer.Accept()
		if err != nil {
			logger.InitLog.Infof("TCP server closed")
			return
		}
		handleConnection(self.TcpConn)
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

func main() {
	Initailize()
	simulator_util.ParseRanContext()

	path, err := filepath.Abs(filepath.Dir(config))
	if err != nil {
		logger.SimulatorLog.Errorf(err.Error())
	}
	simulator_util.ParseUeData(path+"/", factory.SimConfig.UeInfoFile)
	simulator_util.InitUeToDB()

	for _, ran := range self.RanPool {
		simulator_init.RanStart(ran)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		Terminate()
		os.Exit(0)
	}()
	// TCP server for cli test UE
	startTcpServer()
	simulator_util.ClearDB()
}
