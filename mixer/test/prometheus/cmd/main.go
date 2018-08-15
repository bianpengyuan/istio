// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"istio.io/istio/mixer/test/prometheus"
)

// Args represents args consumed by prometheus OOP adapter.
type Args struct {
	// Port to start the grpc adapter on
	AdapterPort uint16

	// Port to use for the prometheus endpoint
	PrometheusPort uint16

	Mtls      bool
	BasicAuth bool
}

func defaultArgs() *Args {
	return &Args{
		AdapterPort:    uint16(8080),
		PrometheusPort: uint16(42422),
		Mtls:           false,
		BasicAuth:      false,
	}
}

// GetCmd returns t he cobra command-tree.
func GetCmd(args []string) *cobra.Command {
	sa := defaultArgs()
	cmd := &cobra.Command{
		Use:   "prometheus",
		Short: "Prometheus out of process adapter.",
		Run: func(cmd *cobra.Command, args []string) {
			runServer(sa)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("'%s' is an invalid argument", args[0])
			}
			return nil
		},
	}

	f := cmd.PersistentFlags()
	f.Uint16VarP(&sa.AdapterPort, "port", "p", sa.AdapterPort,
		"TCP port to use for gRPC Adapter API")
	f.Uint16VarP(&sa.PrometheusPort, "prometheusport", "a", sa.PrometheusPort,
		"TCP port to expose prometheus endpoint on")
	f.BoolVarP(&sa.Mtls, "enable_mtls", "m", false, "")
	f.BoolVarP(&sa.BasicAuth, "enable_basic_auth", "b", false, "")

	return cmd
}

func main() {
	cmd := GetCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func runServer(args *Args) {
	fmt.Printf("here in the args are %v", args)
	s, err := prometheus.NewNoSessionServer(args.AdapterPort, args.PrometheusPort, args.Mtls, args.BasicAuth)
	if err != nil {
		fmt.Printf("unable to start sever: %v", err)
		os.Exit(-1)
	}

	s.Run()
	s.Wait()
}
