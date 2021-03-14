package ran

import (
	"errors"
	"net"
	"sync"
	"syscall"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/free5gc/ngap"
	"github.com/jay16213/radio_simulator/pkg/api"
	"github.com/jay16213/radio_simulator/pkg/factory"
	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
	"github.com/jay16213/radio_simulator/pkg/simulator_nas"
	"github.com/jay16213/radio_simulator/pkg/simulator_ngap"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var DefaultConfigPath = "configs/rancfg.yaml"

type RanApp struct {
	// config information read from file
	cfg *factory.Config

	// ran context
	primaryAMFEndpoint *sctp.SCTPAddr
	ngController       *simulator_ngap.NGController
	ctx                *simulator_context.RanContext
	sctpConn           *sctp.SCTPConn

	// api server provided by grpc
	grpcServer *grpc.Server
}

func New() *RanApp {
	return &RanApp{}
}

func (r *RanApp) Context() *simulator_context.RanContext {
	return r.ctx
}

func (r *RanApp) Initialize(c *cli.Context) error {
	configPath := c.String("config")

	if configPath == "" {
		configPath = DefaultConfigPath
	}

	if config, err := factory.ReadConfig(configPath); err != nil {
		return err
	} else {
		r.cfg = config
	}

	r.ctx = simulator_context.NewRanContext()
	r.ctx.LoadConfig(*r.cfg)
	r.setLogLevel()
	return nil
}

func (r *RanApp) Run() {
	wg := sync.WaitGroup{}

	logger.InitLog.Info("Start running RAN")

	// RAN connect to UPF
	// for _, upf := range ran.UpfInfoList {
	// upf.GtpConn, err = connectToUpf(ran.RanGtpUri.IP, upf.Addr.IP, ran.RanGtpUri.Port, upf.Addr.Port)
	// check(err)
	// simulator_context.Simulator_Self().GtpConnPool[fmt.Sprintf("%s,%s", ran.RanGtpUri.IP, upf.Addr.IP)] = upf.GtpConn
	// go StartHandleGtp(upf)
	// }

	// RAN connect to AMF
	conn, err := r.connectToAmf()
	if err != nil {
		logger.InitLog.Error(err.Error())
		return
	}
	r.sctpConn = conn
	nasController := simulator_nas.NewController()
	r.ngController = simulator_ngap.NewController(r, nasController)
	nasController.SetNGMessager(r.ngController)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		r.StartSCTPAssociation()
	}(&wg)

	// init API service
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		listener, err := net.Listen("tcp", r.cfg.TcpUri)
		if err != nil {
			logger.InitLog.Fatalf("listen error: %v", err)
		}
		r.grpcServer = grpc.NewServer()
		api.RegisterAPIServiceServer(r.grpcServer, &apiService{ranApp: r})
		if err := r.grpcServer.Serve(listener); err != nil {
			logger.InitLog.Fatalf("api server error: %v", err)
		}
	}(&wg)

	wg.Wait()
}

func (r *RanApp) connectToAmf() (*sctp.SCTPConn, error) {
	amfAddr := &sctp.SCTPAddr{
		Port: r.cfg.AmfSCTPEndpoint.Port,
	}
	for _, ip := range r.cfg.AmfSCTPEndpoint.IPs {
		amfAddr.IPAddrs = append(amfAddr.IPAddrs, net.IPAddr{IP: ip})
	}
	r.primaryAMFEndpoint = amfAddr

	ranAddr := &sctp.SCTPAddr{
		Port: r.cfg.RanSctpEndpoint.Port,
	}
	for _, ip := range r.cfg.AmfSCTPEndpoint.IPs {
		ranAddr.IPAddrs = append(ranAddr.IPAddrs, net.IPAddr{IP: ip})
	}

	sock, err := syscall.Socket(
		syscall.AF_INET,
		syscall.SOCK_SEQPACKET,
		syscall.IPPROTO_SCTP,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			syscall.Close(sock)
		}
	}()

	// conn := sctp.NewSCTPConn(sock, nil)
	// if err = sctp.SCTPBind(sock, ranAddr, sctp.SCTP_BINDX_ADD_ADDR); err != nil {
	// 	return nil, err
	// }
	logger.InitLog.Infof("Connecting to sctp addr: %+v", r.primaryAMFEndpoint)
	conn, err := sctp.DialSCTPOneToMany("sctp", ranAddr, r.primaryAMFEndpoint)
	if err != nil {
		return nil, err
	}
	info, _ := conn.GetDefaultSentParam()
	info.PPID = ngap.PPID
	if err = conn.SetDefaultSentParam(info); err != nil {
		return nil, err
	}
	err = conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
	if err != nil {
		logger.NgapLog.Errorf("Failed to subscribe SCTP Event: %v", err)
	}
	return conn, nil
}

func (r *RanApp) Connect(endpoint *sctp.SCTPAddr) error {
	if r.sctpConn == nil {
		return errors.New("sctp connection is nil")
	}
	return r.sctpConn.Connect(endpoint)
}

// sctp send
func (r *RanApp) SendToAMF(endpoint *sctp.SCTPAddr, pkt []byte) {
	_, err := r.sctpConn.SCTPSendTo(pkt,
		&sctp.SndRcvInfo{
			PPID: ngap.PPID,
		},
		endpoint.ToSockaddr(0),
	)
	if err != nil {
		logger.InitLog.Error(err)
	}
}

func (r *RanApp) StartSCTPAssociation() {
	defer r.sctpConn.Close()

	// trigger the initial NG-C procedure
	r.ngController.SendNGSetupRequest(r.primaryAMFEndpoint)

	for {
		buf := make([]byte, 8192)
		n, info, _, endpoint, err := r.sctpConn.SCTPReadFrom(buf)
		if err != nil {
			logger.NgapLog.Debugf("Read Error: %v", err)
			break
		}
		if info == nil || info.PPID != ngap.PPID {
			logger.NgapLog.Warnf("Recv SCTP PPID != 60")
			if info != nil {
				logger.NgapLog.Warnf("info: %+v", info)
			} else {
				logger.NgapLog.Error("info is nil")
			}
			continue
		}
		r.ngController.Dispatch(sctp.SockaddrToSCTPAddr(endpoint), buf[:n])
	}
}

func (r *RanApp) Terminate() {
	logger.SimulatorLog.Infof("Terminating RAN...")

	// TODO: Send UE Deregistration to AMF
	// logger.SimulatorLog.Infof("Clear UE DB...")

	// simulator_util.ClearDB()

	logger.SimulatorLog.Infof("Close SCTP Connection...")
	if err := r.sctpConn.Close(); err != nil {
		logger.SimulatorLog.Errorf("sctp close error: %+v", err)
	}
	// for _, ran := range r. {
	// 	logger.SimulatorLog.Infof("Ran[%s] Connection Close", ran.RanSctpUri)
	// 	ran.SctpConn.Close()
	// }

	logger.SimulatorLog.Infof("Close gRPC API Server")
	r.grpcServer.Stop()

	logger.SimulatorLog.Infof("Clean Ue IP Addr in IP tables")

	// for key, conn := range self.GtpConnPool {
	// 	logger.InitLog.Infof("GTP[%s] Connection Close", key)
	// 	conn.Close()
	// }
	// for _, ue := range self.UeContextPool {
	// 	for _, sess := range ue.PduSession {
	// 		if sess.UeIp != "" {
	// 			_, err := exec.Command("ip", "addr", "del", sess.UeIp, "dev", "lo").Output()
	// 			if err != nil {
	// 				logger.SimulatorLog.Errorf("Delete ue addr failed[%s]", err.Error())
	// 			}
	// 		}
	// 	}
	// }

	// logger.SimulatorLog.Infof("Close Raw Socket...")
	// if self.ListenRawConn != nil {
	// 	self.ListenRawConn.Close()
	// }

	logger.SimulatorLog.Infof("RAN terminated")
}

func (r *RanApp) setLogLevel() {
	if r.cfg.Logger.DebugLevel != "" {
		level, err := logrus.ParseLevel(r.cfg.Logger.DebugLevel)
		if err == nil {
			logger.SetLogLevel(level)
		}
	}
	logger.SetReportCaller(r.cfg.Logger.ReportCaller)
}
