// +build integ
// Copyright Istio Authors. All Rights Reserved.
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

package google

import (
	"testing"

	"istio.io/istio/pkg/test/framework"
	sd "istio.io/istio/tests/integration/telemetry/stackdriver"
)

func TestMain(m *testing.M) {
	// TODO: uncomment this!!!
	// s := framework.NewSuite(m).
	// 	Label(label.Postsubmit)
	s := framework.NewSuite(m)
	sd.TestMainSetup(s, true /* useRealSD */)
}

func TestStackdriverMonitoring(t *testing.T) {
	sd.TestStackdriverMonitoring(t)
}

func TestStackdriverHTTPAuditLogging(t *testing.T) {
	sd.TestStackdriverHTTPAuditLogging(t)
}

func TestTCPStackdriverMonitoring(t *testing.T) {
	sd.TestTCPStackdriverMonitoring(t)
}
