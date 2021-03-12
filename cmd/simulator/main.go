package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jay16213/radio_simulator/pkg/ran"
	"github.com/urfave/cli/v2"
)

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

	ranApp := ran.New()
	if err := ranApp.Initialize(c); err != nil {
		return err
	}

	// simulator_util.ParseRanContext()
	// simulator_util.ParseTunDev()

	// TODO: ue util
	// MongoDBLibrary.SetMongoDB(config.DBName, config.DBUrl)
	// rootPath, err := filepath.Abs(filepath.Dir(ranConfigPath))
	// if err != nil {
	// 	logger.SimulatorLog.Errorf(err.Error())
	// }
	// simulator_util.ParseUeData(rootPath+"/", factory.SimConfig.UeInfoFile)
	// simulator_util.InitUeToDB()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		ranApp.Terminate()
		os.Exit(0)
	}()
	ranApp.Run()
	// Raw Socket Server
	// self.ListenRawConn = raw_socket.ListenRawSocket(factory.SimConfig.ListenIp)
	// TCP server for cli test UE
	// tcp_server.StartTcpServer()
	// simulator_util.ClearDB()
	return nil
}
