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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"istio.io/istio/cni/pkg/constants"
	"istio.io/pkg/log"
)

func TestUDSLog(t *testing.T) {
	// Start UDS log server
	udsSockDir := t.TempDir()
	udsSock := filepath.Join(udsSockDir, "cni.sock")
	logger := NewUDSLogger()
	stop := make(chan struct{})
	defer func() {
		close(stop)
	}()
	logger.StartUDSLogServer(udsSock, stop)

	// Configure log to tee to UDS server, and print to a tmp file
	r, w, _ := os.Pipe()
	os.Stdout = w
	loggingOptions := log.DefaultOptions()
	loggingOptions.WithTeeToUDS(udsSock, constants.UDSLogPath)
	log.Configure(loggingOptions)
	log.FindScope("default").SetOutputLevel(log.DebugLevel)
	log.Debug("debug log")
	log.Info("info log")
	log.Warn("warn log")
	log.Error("error log")
	log.Sync()

	w.Close()
	out, _ := ioutil.ReadAll(r)

	// For each level, there should be two lines, one from direct log,
	// the other one from UDS server
	wantLevels := []string{"debug", "info", "warn", "error", "debug", "info", "warn", "error"}
	gotLogs := strings.Split(
		strings.TrimSuffix(string(out), "\n"), "\n")
	if want, got := len(wantLevels), len(gotLogs); want != got {
		t.Fatalf("Number of logs want %v, got %v", want, got)
	}

	for i, l := range gotLogs {
		// For each line, there should be two level string, e.g.
		// "2021-07-09T03:26:08.984951Z	debug	debug log"
		if got, want := strings.Count(l, wantLevels[i]), 2; want != got {
			t.Errorf("Number of log level string want %v, got %v", want, got)
		}
	}
}
