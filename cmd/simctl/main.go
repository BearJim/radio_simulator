package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jay16213/radio_simulator/pkg/simulator"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"
	"github.com/mohae/deepcopy"
	"github.com/spf13/cobra"
)

// flags
var (
	simulatorDBUrl   string
	numOfUEs         int32
	followOnRequest  bool
	triggerFailCount int32
)

var rootCmd = &cobra.Command{
	Use:     "simctl",
	Short:   "simctl - cli for Radio Simulator",
	Version: "v0.0.1",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&simulatorDBUrl, "db", "mongodb://127.0.0.1:27017", "Database URL for simulator")
	rootCmd.AddCommand(uploadCommand())
	rootCmd.AddCommand(loadCommand())
	rootCmd.AddCommand(getCommand())
	rootCmd.AddCommand(describeCommand())
	rootCmd.AddCommand(registerCommand())
	rootCmd.AddCommand(registerAllCommand())
	rootCmd.AddCommand(serviceRequestCommand())
	rootCmd.AddCommand(deregisterCommand())
	rootCmd.AddCommand(deregisterAllCommand())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func uploadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upload",
		Short:   "upload all UEs to free5gc DB",
		Example: "upload mongodb://127.0.0.1:27017",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("should provide the url of free5gc DB")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			s := initSimulator(simulatorDBUrl)
			s.UploadUEProfile("free5gc", args[0])
		},
	}
	return cmd
}

func getCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "get information of RAN or UEs",
		Example: "get ues",
	}
	getUEs := &cobra.Command{
		Use:   "ues",
		Args:  cobra.NoArgs,
		Short: "get information of all UEs",
		Run: func(cmd *cobra.Command, args []string) {
			s := initSimulator(simulatorDBUrl)
			s.GetUEs()
		},
	}
	getRANs := &cobra.Command{
		Use:   "rans",
		Args:  cobra.NoArgs,
		Short: "get information of all RANs",
		Run: func(cmd *cobra.Command, args []string) {
			s := initSimulator(simulatorDBUrl)
			s.GetRANs()
		},
	}
	cmd.AddCommand(getUEs, getRANs)
	return cmd
}

func describeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe [ue|ran] [<SUPI>|<RanName>]",
		Short:   "describe the detail information of RAN or UE",
		Example: "describe ue imsi-2089300000003",
	}
	describeUE := &cobra.Command{
		Use:   "ue <SUPI>",
		Args:  cobra.ExactArgs(1),
		Short: "get the detail information of UE with SUPI",
		Run: func(cmd *cobra.Command, args []string) {
			s := initSimulator(simulatorDBUrl)
			s.DescribeUE(args[0])
		},
	}
	// getRANs := &cobra.Command{
	// 	Use:   "rans",
	// 	Args:  cobra.NoArgs,
	// 	Short: "get information of all RANs",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		s := initSimulator(simulatorDBUrl)
	// 		s.GetRANs()
	// 	},
	// }
	cmd.AddCommand(describeUE)
	return cmd
}

func registerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "reg <SUPI> <RanName>",
		Short:   "trigger initial registration procedure for UE with SUPI via RanName",
		Example: "reg imsi-2089300000003 ran1",
		Args:    cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			s := initSimulator(simulatorDBUrl)
			s.SingleUeRegister(args[0], args[1], int(triggerFailCount), followOnRequest)
		},
	}
	cmd.PersistentFlags().Int32VarP(&triggerFailCount, "fail", "f", 0, "trigger AMF fail ue count")
	cmd.PersistentFlags().BoolVarP(&followOnRequest, "for", "o", false, "follow-on request pending")
	return cmd
}

func registerAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "regall <RanName>",
		Short:   "trigger initial registration procedure for all UEs via RanName",
		Example: "regall ran1",
		Args:    cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			s := initSimulator(simulatorDBUrl)
			s.AllUeRegister(args[0], int(triggerFailCount), followOnRequest)
		},
	}
	cmd.PersistentFlags().Int32VarP(&triggerFailCount, "fail", "f", 0, "trigger AMF fail ue count")
	cmd.PersistentFlags().BoolVarP(&followOnRequest, "for", "o", false, "follow-on request pending")
	return cmd
}

func serviceRequestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "srvreq <SUPI>",
		Short:   "trigger service request procedure for UE with SUPI",
		Example: "srvreq imsi-2089300000001",
		Args:    cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			s := initSimulator(simulatorDBUrl)
			s.SingleUeServiceRequest(args[0], int(triggerFailCount))
		},
	}
	cmd.PersistentFlags().Int32VarP(&triggerFailCount, "fail", "f", 0, "trigger AMF fail ue count")
	return cmd
}

func deregisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dereg <SUPI>",
		Short:   "trigger deregistration procedure for UE with SUPI",
		Example: "dereg imsi-2089300000003",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			s := initSimulator(simulatorDBUrl)
			s.SingleUeDeregister(args[0])
		},
	}
	return cmd
}

func deregisterAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deregall <RanName>",
		Short:   "trigger deregistration procedure for all UEs in RanName",
		Example: "deregall ran1",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			s := initSimulator(simulatorDBUrl)
			s.AllUeDeregister(args[0])
		},
	}
	return cmd
}

func loadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "load <file path>",
		Short:   "load UE context from FILE",
		Example: "load config/uecfg.yaml",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			s := initSimulator(simulatorDBUrl)
			ueContexts := s.ParseUEData(args[0])
			if cmd.HasPersistentFlags() {
				defaultUE := ueContexts[0]

				for i := 2; i <= int(numOfUEs); i++ {
					ue := deepcopy.Copy(defaultUE).(*simulator_context.UeContext)
					ue.Supi = fmt.Sprintf("imsi-20893%08d", i)
					ueContexts = append(ueContexts, ue)
				}
			}
			s.InsertUEContextToDB(ueContexts)
		},
	}
	cmd.PersistentFlags().Int32VarP(&numOfUEs, "generate", "g", 10, "automatically generate UE profile for input quantity")
	return cmd
}

func initSimulator(dbUrl string) *simulator.Simulator {
	// logger. ("db: %s\n", dbUrl)
	s, err := simulator.New("simulator", dbUrl)
	if err != nil {
		fmt.Printf("Init error: %+v\n", err)
		os.Exit(1)
	}
	return s
}
