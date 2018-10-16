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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"istio.io/istio/mixer/adapter/rbac/cmd/server"
)

// Args represents args consumed by OOP adapter.
type Args struct {
	// Port to start the grpc adapter on
	AdapterPort uint16

	// Size of adapter API worker pool
	APIWorkerPoolSize int
}

func defaultArgs() *Args {
	return &Args{
		AdapterPort:       uint16(8080),
		APIWorkerPoolSize: 1024,
	}
}

// GetCmd returns the cobra command-tree.
func GetCmd(args []string) *cobra.Command {
	sa := defaultArgs()
	cmd := &cobra.Command{
		Use:   "rbac",
		Short: "rbac out of process adapter.",
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
	f.IntVarP(&sa.APIWorkerPoolSize, "api_pool", "a", sa.APIWorkerPoolSize,
		"Size of adapter API worker pool")
	return cmd
}

func main() {
	cmd := GetCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func runServer(args *Args) {
	s, err := server.NewRbacNoSessionServer(args.AdapterPort, args.APIWorkerPoolSize)
	if err != nil {
		fmt.Printf("unable to start server: %v", err)
		os.Exit(-1)
	}

	s.Run()
	s.Wait()
}
