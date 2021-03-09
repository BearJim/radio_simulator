package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/jay16213/radio_simulator/pkg/factory"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
	"github.com/jay16213/radio_simulator/pkg/simulator_init"
	"github.com/jay16213/radio_simulator/pkg/simulator_util"
	"github.com/jay16213/radio_simulator/pkg/tcp_server"
	"github.com/urfave/cli/v2"

	"github.com/free5gc/MongoDBLibrary"
	"github.com/sirupsen/logrus"
)

var self *simulator_context.Simulator = simulator_context.Simulator_Self()

func Terminate() {
	logger.SimulatorLog.Infof("Terminating Simulator...")

	// TODO: Send UE Deregistration to AMF
	logger.SimulatorLog.Infof("Clear UE DB...")

	simulator_util.ClearDB()

	logger.SimulatorLog.Infof("Close SCTP Connection...")

	for _, ran := range self.RanPool {
		logger.SimulatorLog.Infof("Ran[%s] Connection Close", ran.RanSctpUri)
		ran.SctpConn.Close()
	}

	logger.SimulatorLog.Infof("Close TCP Server...")

	if self.TcpServer != nil {
		self.TcpServer.Close()
	}

	logger.SimulatorLog.Infof("Clean Ue IP Addr in IP tables")

	// for key, conn := range self.GtpConnPool {
	// 	logger.InitLog.Infof("GTP[%s] Connection Close", key)
	// 	conn.Close()
	// }
	for _, ue := range self.UeContextPool {
		for _, sess := range ue.PduSession {
			if sess.UeIp != "" {
				_, err := exec.Command("ip", "addr", "del", sess.UeIp, "dev", "lo").Output()
				if err != nil {
					logger.SimulatorLog.Errorf("Delete ue addr failed[%s]", err.Error())
				}
			}
		}
	}

	// logger.SimulatorLog.Infof("Close Raw Socket...")
	// if self.ListenRawConn != nil {
	// 	self.ListenRawConn.Close()
	// }

	logger.SimulatorLog.Infof("Simulator terminated")

}

func main() {
	app := cli.NewApp()
	app.Name = "Radio Simulator"
	app.Usage = "5G NG-RAN and UE Simulator"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   "",
			Usage:   "Load configuration from `FILE`",
		},
	}
	app.Action = action

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Simulator run error: %+v\n", err)
		os.Exit(1)
	}
}

func action(c *cli.Context) error {
	if c.Bool("help") {
		cli.ShowAppHelpAndExit(c, 0)
	}

	ranConfigPath := c.String("config")
	if ranConfigPath == "" {
		ranConfigPath = "./configs/rancfg.conf"
	}

	factory.InitConfigFactory(ranConfigPath)

	if factory.SimConfig.Logger.DebugLevel != "" {
		level, err := logrus.ParseLevel(factory.SimConfig.Logger.DebugLevel)
		if err == nil {
			logger.SetLogLevel(level)
		}
	}
	logger.SetReportCaller(factory.SimConfig.Logger.ReportCaller)
	MongoDBLibrary.SetMongoDB(factory.SimConfig.DBName, factory.SimConfig.DBUrl)

	simulator_util.ParseRanContext()
	simulator_util.ParseTunDev()

	rootPath, err := filepath.Abs(filepath.Dir(ranConfigPath))
	if err != nil {
		logger.SimulatorLog.Errorf(err.Error())
	}
	simulator_util.ParseUeData(rootPath+"/", factory.SimConfig.UeInfoFile)
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
	// Raw Socket Server
	// self.ListenRawConn = raw_socket.ListenRawSocket(factory.SimConfig.ListenIp)
	// TCP server for cli test UE
	tcp_server.StartTcpServer()
	simulator_util.ClearDB()
	return nil
}
