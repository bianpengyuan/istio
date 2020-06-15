package gcpmonitoring

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	operationKey, _ = tag.NewKey("operation")
	successKey, _   = tag.NewKey("success")
	typeKey, _      = tag.NewKey("type")

	configEventMeasure        = stats.Int64("config_event_measure", "The number of user configuration event", "1")
	configValidationMeasuare  = stats.Int64("config_validation_measure", "The number of configuration validation event", "1")
	configPushMeasuare        = stats.Int64("config_push_measure", "The number of xds configuration pushes", "1")
	rejectedConfigMeasuare    = stats.Int64("rejected_config_measure", "The number of rejected configurations", "1")
	configConvergenceMeasuare = stats.Float64("config_convergence_measure", "The distribution of latency of user configuration becomes effective", "ms")
	proxyClientsMeasure       = stats.Int64("proxy_client_measure", "The number of proxies connected to this instance", "1")
	sidecarInjectionMeasure   = stats.Int64("sidecar_injection_measure", "The number of sidecar injection event", "1")

	configEventView = &view.View{
		Name:        "config_event_count",
		Measure:     configEventMeasure,
		Description: "number of config events",
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{operationKey, typeKey},
	}
	configValidationView = &view.View{
		Name:        "config_validation_count",
		Measure:     configValidationMeasuare,
		Description: "number of config validation",
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{successKey, typeKey},
	}
	configPushView = &view.View{
		Name:        "config_push_count",
		Measure:     configPushMeasuare,
		Description: "number of config pushes",
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{successKey, typeKey},
	}
	rejectedConfigView = &view.View{
		Name:        "rejected_config_count",
		Measure:     rejectedConfigMeasuare,
		Description: "number of rejected config",
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{typeKey},
	}
	configConvergenceView = &view.View{
		Name:        "config_convergence_latencies",
		Measure:     configConvergenceMeasuare,
		Description: "Latency of config become effective after being applied",
		Aggregation: view.Distribution([]float64{0.001, 0.002, 0.004, 0.008, 0.016, 0.032, 0.064,
			0.128, 0.256, 0.512, 1, 2, 4, 8, 16, 32}...),
		TagKeys: []tag.Key{typeKey},
	}
	proxyClientsView = &view.View{
		Name:        "proxy_clients",
		Measure:     proxyClientsMeasure,
		Description: "Number of proxies connected to this instance",
		Aggregation: view.LastValue(),
	}
	sidecarInjectionView = &view.View{
		Name:        "sidecar_injection_count",
		Measure:     sidecarInjectionMeasure,
		Description: "Number of sidecar injection event",
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{successKey},
	}

	viewMap = make(map[string]bool)
)

func registerView(v *view.View) {
	view.Register(v)
	viewMap[v.Name] = true
}

func init() {
	registerView(configEventView)
	registerView(configValidationView)
	registerView(configPushView)
	registerView(rejectedConfigView)
	registerView(configConvergenceView)
	registerView(proxyClientsView)
	registerView(sidecarInjectionView)

	registerHook()
}
