// Copyright 2017 Istio Authors
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

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/codegen/cmd/mixgeninventory/main.go -f $GOPATH/src/istio.io/istio/mixer/adapter/inventory.yaml -o $GOPATH/src/istio.io/istio/mixer/adapter/inventory.gen.go

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/logentry/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/reportnothing/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/authorization/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/checknothing/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/listentry/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/quota/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/tracespan/tracespan_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/noop --adapter_name noop --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/noop/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/reportnothing/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/checknothing/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/quota/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/bypass --config_package istio.io/istio/mixer/adapter/bypass/config --adapter_name bypass --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/bypass/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/circonus --config_package istio.io/istio/mixer/adapter/circonus/config --adapter_name circonus --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/circonus/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/cloudwatch --config_package istio.io/istio/mixer/adapter/cloudwatch/config --adapter_name cloudwatch --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/cloudwatch/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/checknothing/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/listentry/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/quota/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/denier --adapter_name denier --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/denier/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/dogstatsd --config_package istio.io/istio/mixer/adapter/dogstatsd/config --adapter_name dogstatsd --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/dogstatsd/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/logentry/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/fluentd --config_package istio.io/istio/mixer/adapter/fluentd/config --adapter_name fluentd --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/fluentd/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/listentry/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/list --config_package istio.io/istio/mixer/adapter/list/config --adapter_name listchecker --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/list/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/quota/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/memquota --config_package istio.io/istio/mixer/adapter/memquota/config --adapter_name memquota --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/memquota/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/quota/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/redisquota --config_package istio.io/istio/mixer/adapter/redisquota/config --adapter_name redisquota --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/redisquota/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/authorization/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/opa --config_package istio.io/istio/mixer/adapter/opa/config --adapter_name opa --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/opa/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/prometheus --config_package istio.io/istio/mixer/adapter/prometheus/config --adapter_name prometheus --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/prometheus/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/authorization/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/rbac --config_package istio.io/istio/mixer/adapter/rbac/config --adapter_name rbac --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/rbac/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/apikey/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/quota/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/adapter/servicecontrol/template/servicecontrolreport/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/servicecontrol --config_package istio.io/istio/mixer/adapter/servicecontrol/config --adapter_name servicecontrol --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/servicecontrol/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/tracespan/tracespan_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/signalfx --adapter_name signalfx --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/signalfx/cmd --config_package istio.io/istio/mixer/adapter/signalfx/config

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/logentry/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/solarwinds --adapter_name solarwinds --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/solarwinds/cmd --config_package istio.io/istio/mixer/adapter/solarwinds/config

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/logentry/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/tracespan/tracespan_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/edge/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/stackdriver --adapter_name stackdriver --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/stackdriver/cmd --config_package istio.io/istio/mixer/adapter/stackdriver/config

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/statsd --config_package istio.io/istio/mixer/adapter/statsd/config --adapter_name statsd --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/statsd/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/logentry/template_proto.descriptor_set --adapter_package istio.io/istio/mixer/adapter/stdio --adapter_name stdio --out_dir $GOPATH/src/istio.io/istio/mixer/adapter/stdio/cmd --config_package istio.io/istio/mixer/adapter/stdio/config

// Package adapter contains the inventory for all Mixer adapters that are compiled into
// a specific Mixer binary.
package adapter
