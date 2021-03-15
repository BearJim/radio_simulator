package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/jay16213/radio_simulator/pkg/api"
	"google.golang.org/grpc"
)

type Simulator struct {
	cc      *exec.Cmd
	RanPool map[string]api.APIServiceClient // RanSctpUri -> RAN_CONTEXT
	// UeContextPool map[string]*UeContext  // Supi -> UeTestInfo
}

func main() {
	s := &Simulator{
		RanPool: make(map[string]api.APIServiceClient),
	}
	s.StartNewRan()
	time.Sleep(100 * time.Millisecond)
	runCli(s)
}

func runCli(s *Simulator) {
	var cmd string
	for {
		fmt.Printf("> ")
		fmt.Scanln(&cmd)
		s.executor(cmd)
	}
}

func (s *Simulator) StartNewRan() {
	c := exec.Command("./bin/simulator")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}

	if err := c.Start(); err != nil {
		fmt.Printf("run error: %+v\n", err)
	} else {
		fmt.Printf("c.Run err is nil\n")
	}
	s.cc = c
}

func (s *Simulator) NewRANClient(client api.APIServiceClient, ranName string) {
	if _, ok := s.RanPool[ranName]; ok {
		fmt.Printf("duplicate ran name %s\n", ranName)
	} else {
		s.RanPool[ranName] = client
	}
}

func (s *Simulator) executor(command string) {
	if strings.HasPrefix(command, "connect") {
		tokens := strings.Split(command, " ")
		if len(tokens) < 2 {
			fmt.Println("command error")
			return
		}
		tokens = tokens[1:] // cut "connect"
		for _, addr := range tokens {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				fmt.Printf("connect %s error: %+v\n", addr, err)
				return
			}

			client := api.NewAPIServiceClient(conn)
			resp, err := client.DescribeRAN(context.Background(), &api.DescribeRANRequest{})
			if err != nil {
				fmt.Printf("DescribeRAN: %+v\n", err)
				return
			}
			s.NewRANClient(client, resp.Name)
			fmt.Printf("Connect to RAN %s\n", resp.Name)
		}
	}

	if command == "get ues" {
		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
		fmt.Fprintln(writer, "SUPI\tCM-STATE\tRM-STATE\tSERVING-RAN")
		fmt.Fprintln(writer, "imsi-2089300000003\tIDLE\tRegistered\tran-1 (127.0.0.1:38412)")
		fmt.Fprintln(writer, "imsi-2089300000004\tConnected\tRegistered\tran-1 (127.0.0.1:38412)")
		fmt.Fprintln(writer, "imsi-2089300000005\tIDLE\tDeregistered\tran-2 (127.0.0.2:38412)")
		writer.Flush()
	}

	if command == "exit" {
		if err := s.cc.Process.Signal(syscall.SIGTERM); err != nil {
			fmt.Printf("Signal: %+v\n", err)
		}
		os.Exit(0)
	}
}
