package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jay16213/radio_simulator/pkg/simulator"
)

var sim *simulator.Simulator

func main() {
	if s, err := simulator.New("simulator", "mongodb://localhost:27017"); err != nil {
		fmt.Printf("Init error: %+v\n", err)
		os.Exit(1)
	} else {
		sim = s
	}

	// s.StartNewRan()
	time.Sleep(100 * time.Millisecond)
	runCli(sim)
}

func runCli(s *simulator.Simulator) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		executor(string(line))
	}
}

func executor(command string) {
	if strings.HasPrefix(command, "connect") {
		tokens := tokenize(command)
		if len(tokens) < 1 {
			fmt.Println("command error")
			return
		}

		for _, addr := range tokens {
			if name, err := sim.ConnectToRAN(addr); err != nil {
				fmt.Printf("connect %s error: %+v\n", addr, err)
			} else {
				fmt.Printf("Connect to %s (name: %s)\n", addr, name)
			}
		}
	}

	if strings.HasPrefix(command, "load") {
		rootPath := "./configs/"
		ueContexts := sim.ParseUEData(rootPath, []string{"uecfg.yaml"})
		sim.InsertUEContextToDB(ueContexts)
	}

	if strings.HasPrefix(command, "reg") {
		tokens := tokenize(command)
		if len(tokens) != 2 {
			fmt.Println("command error")
			return
		}
		sim.UeRegister(tokens[0], tokens[1])
	}

	if strings.HasPrefix(command, "dereg") {
		tokens := tokenize(command)
		if len(tokens) != 1 {
			fmt.Println("command error")
			return
		}
		sim.UeDeregister(tokens[0])
	}

	if strings.HasPrefix(command, "upload") {
		tokens := tokenize(command)
		if len(tokens) != 1 {
			fmt.Println("command error")
			return
		}
		sim.UploadUEProfile("free5gc", tokens[0])
	}

	if command == "get" {
		sim.GetUEs()
	}

	if command == "exit" {
		os.Exit(0)
	}
}

func tokenize(cmd string) []string {
	tokens := strings.Split(cmd, " ")
	return tokens[1:]
}
