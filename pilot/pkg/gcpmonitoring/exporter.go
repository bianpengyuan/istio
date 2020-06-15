package gcpmonitoring

import (
	"errors"
	"os"
	"sync"
	"time"

	ocprom "contrib.go.opencensus.io/exporter/prometheus"
	"contrib.go.opencensus.io/exporter/stackdriver"
	ocsd "contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/monitoredresource"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"google.golang.org/api/option"

	"istio.io/istio/pkg/bootstrap/platform"
	"istio.io/pkg/version"
)

type ASMExporter struct {
	PromExporter *ocprom.Exporter
	sdExporter   *ocsd.Exporter
}

func NewASMExporter(pe *ocprom.Exporter) (*ASMExporter, error) {
	labels := &stackdriver.Labels{}
	labels.Set("mesh_uid", "", "ID for Mesh")
	labels.Set("revision", version.Info.Version, "Control plane revision")
	gcpMetadata := platform.NewGCP().Metadata()
	se, err := stackdriver.NewExporter(stackdriver.Options{
		MetricPrefix: "istio.io/control",
		MonitoringClientOptions: []option.ClientOption{
			option.WithEndpoint("staging-monitoring.sandbox.googleapis.com:443"),
			// option.WithTokenSource(&tokenmanager.TokenSource{}),
		},
		GetMetricType: func(view *view.View) string {
			return "istio.io/control/" + view.Name
		},
		MonitoredResource: &monitoredresource.GKEContainer{
			ProjectID:                  gcpMetadata[platform.GCPProject],
			ClusterName:                gcpMetadata[platform.GCPCluster],
			Zone:                       gcpMetadata[platform.GCPLocation],
			NamespaceID:                "istio-system",
			PodID:                      os.Getenv("HOSTNAME"),
			ContainerName:              "discovery",
			LoggingMonitoringV2Enabled: true,
		},
		DefaultMonitoringLabels: labels,
		ReportingInterval:       60 * time.Second,
	})

	if err != nil {
		return nil, errors.New("fail to initialize Stackdriver exporter")
	}

	return &ASMExporter{
		PromExporter: pe,
		sdExporter:   se,
	}, nil
}

func (e *ASMExporter) ExportView(vd *view.Data) {
	if _, ok := viewMap[vd.View.Name]; ok {
		// This indicates that this is a stackdriver view
		e.sdExporter.ExportView(vd)
	} else {
		e.PromExporter.ExportView(vd)
	}
}

type TestExporter struct {
	sync.Mutex

	Rows        map[string][]*view.Row
	invalidTags bool
}

func (t *TestExporter) ExportView(d *view.Data) {
	t.Lock()
	for _, tk := range d.View.TagKeys {
		if len(tk.Name()) < 1 {
			t.invalidTags = true
		}
	}
	t.Rows[d.View.Name] = append(t.Rows[d.View.Name], d.Rows...)
	t.Unlock()
}

func FindTagWithValue(key, value string, tags []tag.Tag) bool {
	for _, t := range tags {
		if t.Key.Name() == key && t.Value == value {
			return true
		}
	}
	return false
}
