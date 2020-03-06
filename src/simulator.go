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
	StartTcpServer()
	simulator_util.ClearDB()
}
