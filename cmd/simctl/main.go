package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/c-bata/go-prompt"
	"google.golang.org/grpc"
)

type Simulator struct {
	cc      *exec.Cmd
	process *os.Process
	gClient []*grpc.ClientConn
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "get ues", Description: "Store the username and age"},
		{Text: "describe ue", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	s := &Simulator{}

	s.StartNewRan()
	time.Sleep(100 * time.Millisecond)
	runCli(s)
}

func runCli(s *Simulator) {
	for {
		cmd := prompt.Input("> ", completer)
		s.executor(cmd)
	}
}

func (s *Simulator) StartNewRan() {
	c := exec.Command("./bin/simulator")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		fmt.Printf("ran error: %+v\n", err)
		return
	} else {
		fmt.Printf("ran run with pid %d\n", c.Process.Pid)
		go func() {
			if err := c.Wait(); err != nil {
				if ee, ok := err.(*exec.ExitError); ok && ee.ProcessState.Success() {
					fmt.Println("ran down")
				} else {
					fmt.Printf("wait error: %+v\n", err)
				}
			}
		}()
	}
	s.cc = c
}

func (s *Simulator) executor(command string) {
	if strings.HasPrefix(command, "connect") {
		addr := "127.0.0.1:9999"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("connect %s error: %+v\n", addr, err)
			return
		}
		defer conn.Close()
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
		if err := s.cc.Process.Signal(os.Interrupt); err != nil {
			fmt.Printf("Signal: %+v\n", err)
		}
		fmt.Printf("Exit simctl\n")
		os.Exit(0)
	}
}
