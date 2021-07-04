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

package logging

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"istio.io/pkg/log"
)

// StartUDSServer ...
func StartUDSLogServer(sockFileName string, stop <-chan struct{}) error {
	if sockFileName == "" {
		return nil
	}
	log.Info("Start a UDS server for CNI plugin logs")
	os.Remove("/var/run/istio-cni/cni.sock")
	unixListener, err := net.Listen("unix", "/var/run/istio-cni/cni.sock")
	if err != nil {
		return fmt.Errorf("Failed to create UDS listener: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/log", handleLog)
	loggingServer := &http.Server{
		Handler: mux,
	}

	go func() {
		if err := loggingServer.Serve(unixListener); err != nil {
			log.Errorf("Error running CNI plugin log server: %v", err)
		}
	}()

	go func() {
		<-stop
		err := loggingServer.Close()
		log.Debugf("CNI log server terminated: %v", err)
	}()

	return nil
}

func handleLog(w http.ResponseWriter, req *http.Request) {
	var body []byte
	if req.Body != nil {
		if data, err := ioutil.ReadAll(req.Body); err == nil {
			body = data
		} else {
			log.Error("Failed to read log report from cni plugin.")
		}
	}
	type cniLog struct {
		Level string `json:"level"`
		Msg   string `json:"msg"`
	}
	cniLogs := make([]string, 0)
	err := json.Unmarshal(body, &cniLogs)
	if err != nil {
		log.Errorf("Failed to unmarshal CNI plugin logs: %v", err)
		return
	}
	for _, l := range cniLogs {
		var msg cniLog
		if err := json.Unmarshal([]byte(l), &msg); err != nil {
			log.Errorf("Failed to unmarshal CNI plugin log entry: %v", err)
			continue
		}
		tm := strings.TrimSpace(msg.Msg)
		switch msg.Level {
		case "debug":
			log.Debug(tm)
		case "info":
			log.Info(tm)
		}
	}
}
