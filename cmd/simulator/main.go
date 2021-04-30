package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/free5gc/version"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/ran"
	"github.com/urfave/cli/v2"
)

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version:\t%+s\n"+
			"build time:\t%s\n"+
			"commit hash:\t%s\n"+
			"go version:\t%s\n",
			version.VERSION,
			version.BUILD_TIME,
			version.COMMIT_HASH,
			runtime.Version(),
		)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Ran UE Simulator"
	app.Usage = "5G NG-RAN and UE Simulator"
	app.Version = version.VERSION
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   "",
			Usage:   "Load configuration from `FILE`",
		},
		&cli.StringFlag{
			Name:    "apiaddr",
			Aliases: []string{"a"},
			Value:   "",
			Usage:   "set API server address (ip:port) to `ADDRESS`",
		},
		cli.VersionFlag,
	}
	app.Action = action

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Simulator run error: %+v\n", err)
		os.Exit(1)
	}
}

func action(c *cli.Context) error {
	if c.Bool("version") {
		cli.ShowVersion(c)
		return nil
	}

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

	logger.AppLog.Infow("Start running RAN",
		"version", version.VERSION, "commit hash", version.COMMIT_HASH, "build time", version.BUILD_TIME)
	ranApp.Run()
	// Raw Socket Server
	// self.ListenRawConn = raw_socket.ListenRawSocket(factory.SimConfig.ListenIp)
	// TCP server for cli test UE
	// tcp_server.StartTcpServer()
	// simulator_util.ClearDB()
	return nil
}
