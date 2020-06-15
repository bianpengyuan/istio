package inject

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	gm "istio.io/istio/pilot/pkg/gcpmonitoring"
	"istio.io/istio/pkg/test/util/retry"
	"istio.io/pkg/monitoring"
)

func TestGCPMonitoringSidecarInjection(t *testing.T) {
	exp := &gm.TestExporter{Rows: make(map[string][]*view.Row)}
	view.RegisterExporter(exp)
	view.SetReportingPeriod(1 * time.Millisecond)

	var cases = []struct {
		name    string
		m       monitoring.Metric
		wantVal *view.Row
	}{
		{"totalSuccessfulInjections", totalSuccessfulInjections, &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "success"}}, Data: &view.SumData{1.0}}},
		{"totalFailedInjections", totalFailedInjections, &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "failure"}}, Data: &view.SumData{1.0}}},
		{"totalSkippedInjections", totalSkippedInjections, &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "skip"}}, Data: &view.SumData{1.0}}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp.Rows = make(map[string][]*view.Row)
			tt.m.Increment()
			wantMetric := "control/sidecar_injection_count"
			if err := retry.UntilSuccess(func() error {
				exp.Lock()
				defer exp.Unlock()
				if len(exp.Rows[wantMetric]) < 1 {
					return fmt.Errorf("wanted metrics %v nost received", wantMetric)
				}
				for _, got := range exp.Rows[wantMetric] {
					if !reflect.DeepEqual(got.Tags, tt.wantVal.Tags) {
						continue
					}
					if int64(tt.wantVal.Data.(*view.SumData).Value) == int64(got.Data.(*view.SumData).Value) {
						return nil
					}
				}
				return fmt.Errorf("metrics %v does not have expecte values, want %+v", tt.m.Name(), tt.wantVal)
			}); err != nil {
				t.Fatalf("failed to get expected metric: %v", err)
			}
		})
	}
}
