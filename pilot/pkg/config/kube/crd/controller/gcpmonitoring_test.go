package controller

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"go.opencensus.io/stats/view"
	gm "istio.io/istio/pilot/pkg/gcpmonitoring"
	"istio.io/istio/pkg/test/util/retry"
)

func TestGCPMonitoringPilotK8sCfgEvents(t *testing.T) {
	exp := &gm.TestExporter{Rows: make(map[string][]*view.Row)}
	view.RegisterExporter(exp)
	view.SetReportingPeriod(1 * time.Millisecond)

	incrementEvent("foo", "bar")
	if err := retry.UntilSuccess(func() error {
		exp.Lock()
		defer exp.Unlock()
		got := float64(0)
		if len(exp.Rows["control/config_event_count"]) < 1 {
			return errors.New("wanted metrics not received")
		}
		r := exp.Rows["control/config_event_count"][0]
		if gm.FindTagWithValue("operation", "bar", r.Tags) && gm.FindTagWithValue("type", "foo", r.Tags) {
			if sd, ok := r.Data.(*view.SumData); ok {
				got = sd.Value
			}
		}
		if got != 1.0 {
			return fmt.Errorf("bad value for control/config_event_count: %f, want 1.0", got)
		}
		return nil
	}); err != nil {
		t.Fatalf("failed to get expected metric: %v", err)
	}
}
