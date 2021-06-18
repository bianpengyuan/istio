// Copyright Istio Authors
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

package log

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"istio.io/pkg/log"
)

// StartUDSServer ...
func StartUDSLogServer(sockFileName string) error {
	log.Info("uds log server starts!")
	if sockFileName == "" {
		return nil
	}
	os.Remove("/var/run/istio-cni/cni-1.sock")
	unixListener, err := net.Listen("unix", "/var/run/istio-cni/cni-1.sock")
	if err != nil {
		return fmt.Errorf("failed to create uds server for CNI logging: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/log", handleLog)

	go func() {
		if err := http.Serve(unixListener, mux); err != nil {
			log.Error(err)
		}
	}()

	return nil
}

func handleLog(w http.ResponseWriter, req *http.Request) {
	var body []byte
	if req.Body != nil {
		if data, err := ioutil.ReadAll(req.Body); err == nil {
			body = data
		} else {
			log.Error("failed to read log report from cni plugin")
		}
	}
	log.Info(string(body))
}
