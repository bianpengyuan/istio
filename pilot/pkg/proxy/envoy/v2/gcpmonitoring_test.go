package v2

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

func TestGCPMonitoringPilotXDSMetrics(t *testing.T) {
	exp := &gm.TestExporter{Rows: make(map[string][]*view.Row)}
	view.RegisterExporter(exp)
	view.SetReportingPeriod(1 * time.Millisecond)

	var cases = []struct {
		name       string
		m          monitoring.Metric
		increment  bool
		recordVal  float64
		wantMetric string
		wantVal    *view.Row
	}{
		{"cdsPushes", cdsPushes, true, 0, "control/config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "success"}, {Key: tag.MustNewKey("type"), Value: "CDS"}}, Data: &view.SumData{1.0}}},
		{"edsPushes", edsPushes, true, 0, "control/config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "success"}, {Key: tag.MustNewKey("type"), Value: "EDS"}}, Data: &view.SumData{1.0}}},
		{"ldsPushes", ldsPushes, true, 0, "control/config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "success"}, {Key: tag.MustNewKey("type"), Value: "LDS"}}, Data: &view.SumData{1.0}}},
		{"rdsPushes", rdsPushes, true, 0, "control/config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "success"}, {Key: tag.MustNewKey("type"), Value: "RDS"}}, Data: &view.SumData{1.0}}},
		{"apiPushes", apiPushes, true, 0, "control/config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "success"}, {Key: tag.MustNewKey("type"), Value: "API"}}, Data: &view.SumData{1.0}}},

		{"cdsSendErrPushes", cdsSendErrPushes, true, 0, "control/config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "error"}, {Key: tag.MustNewKey("type"), Value: "CDS"}}, Data: &view.SumData{1.0}}},
		{"edsSendErrPushes", edsSendErrPushes, true, 0, "control/config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "error"}, {Key: tag.MustNewKey("type"), Value: "EDS"}}, Data: &view.SumData{1.0}}},
		{"ldsSendErrPushes", ldsSendErrPushes, true, 0, "control/config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "error"}, {Key: tag.MustNewKey("type"), Value: "LDS"}}, Data: &view.SumData{1.0}}},
		{"rdsSendErrPushes", rdsSendErrPushes, true, 0, "control/config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "error"}, {Key: tag.MustNewKey("type"), Value: "RDS"}}, Data: &view.SumData{1.0}}},
		{"apiSendErrPushes", apiSendErrPushes, true, 0, "control/config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("status"), Value: "error"}, {Key: tag.MustNewKey("type"), Value: "API"}}, Data: &view.SumData{1.0}}},

		{"cdsReject", cdsReject, true, 0, "control/rejected_config_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("type"), Value: "CDS"}}, Data: &view.SumData{1.0}}},
		{"edsReject", edsReject, true, 0, "control/rejected_config_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("type"), Value: "EDS"}}, Data: &view.SumData{1.0}}},
		{"ldsReject", ldsReject, true, 0, "control/rejected_config_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("type"), Value: "LDS"}}, Data: &view.SumData{1.0}}},
		{"rdsReject", rdsReject, true, 0, "control/rejected_config_count", &view.Row{
			Tags: []tag.Tag{{Key: tag.MustNewKey("type"), Value: "RDS"}}, Data: &view.SumData{1.0}}},

		{"xdsClients", xdsClients, false, 10, "control/proxy_clients", &view.Row{
			Tags: []tag.Tag{}, Data: &view.LastValueData{Value: 10.0}}},

		{"proxiesConvergeDelay", proxiesConvergeDelay, false, 0.4, "control/config_convergence_latencies", &view.Row{
			Tags: []tag.Tag{},
			Data: &view.DistributionData{
				CountPerBucket: []int64{0, 1, 0, 0, 0, 0, 0, 0, 0},
			},
		}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp.Rows = make(map[string][]*view.Row)
			if tt.increment {
				tt.m.Increment()
			} else {
				tt.m.Record(tt.recordVal)
			}
			if err := retry.UntilSuccess(func() error {
				exp.Lock()
				defer exp.Unlock()
				if len(exp.Rows[tt.wantMetric]) < 1 {
					return fmt.Errorf("wanted metrics %v not received", tt.wantMetric)
				}
				for _, got := range exp.Rows[tt.wantMetric] {
					if len(got.Tags) != len(tt.wantVal.Tags) ||
						(len(tt.wantVal.Tags) != 0 && !reflect.DeepEqual(got.Tags, tt.wantVal.Tags)) {
						continue
					}
					switch v := tt.wantVal.Data.(type) {
					case *view.SumData:
						if int64(v.Value) == int64(got.Data.(*view.SumData).Value) {
							return nil
						}
					case *view.LastValueData:
						if int64(v.Value) == int64(got.Data.(*view.LastValueData).Value) {
							return nil
						}
					case *view.DistributionData:
						gotDist := got.Data.(*view.DistributionData)
						if reflect.DeepEqual(gotDist.CountPerBucket, v.CountPerBucket) {
							return nil
						}
					}
				}
				return fmt.Errorf("metrics %v does not have expecte values, want %+v", tt.m.Name(), tt.wantVal)
			}); err != nil {
				t.Fatalf("failed to get expected metric: %v", err)
			}
		})
	}
}
