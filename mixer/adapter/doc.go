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

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/logentry/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/reportnothing/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/authorization/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/checknothing/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/listentry/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/quota/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/tracespan/tracespan_proto.descriptor_set -p istio.io/istio/mixer/template/metric -p istio.io/istio/mixer/template/tracespan -p istio.io/istio/mixer/template/logentry -p istio.io/istio/mixer/template/authorization -p istio.io/istio/mixer/template/checknothing -p istio.io/istio/mixer/template/listentry -p istio.io/istio/mixer/template/quota -p istio.io/istio/mixer/template/reportnothing --adapter_package istio.io/istio/mixer/adapter/noop --adapter_name noop --go_out $GOPATH/src/istio.io/istio/mixer/adapter/noop/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/reportnothing/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/checknothing/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/quota/template_proto.descriptor_set -p istio.io/istio/mixer/template/metric -p istio.io/istio/mixer/template/checknothing -p istio.io/istio/mixer/template/quota -p istio.io/istio/mixer/template/reportnothing --adapter_package istio.io/istio/mixer/adapter/bypass --config_package istio.io/istio/mixer/adapter/bypass/config --adapter_name bypass --go_out $GOPATH/src/istio.io/istio/mixer/adapter/bypass/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -p istio.io/istio/mixer/template/metric --adapter_package istio.io/istio/mixer/adapter/circonus --config_package istio.io/istio/mixer/adapter/circonus/config --adapter_name circonus --go_out $GOPATH/src/istio.io/istio/mixer/adapter/circonus/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -p istio.io/istio/mixer/template/metric --adapter_package istio.io/istio/mixer/adapter/cloudwatch --config_package istio.io/istio/mixer/adapter/cloudwatch/config --adapter_name cloudwatch --go_out $GOPATH/src/istio.io/istio/mixer/adapter/cloudwatch/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/checknothing/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/listentry/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/quota/template_proto.descriptor_set -p istio.io/istio/mixer/template/checknothing -p istio.io/istio/mixer/template/quota -p istio.io/istio/mixer/template/listentry --adapter_package istio.io/istio/mixer/adapter/denier --adapter_name denier --go_out $GOPATH/src/istio.io/istio/mixer/adapter/denier/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -p istio.io/istio/mixer/template/metric --adapter_package istio.io/istio/mixer/adapter/dogstatsd --config_package istio.io/istio/mixer/adapter/dogstatsd/config --adapter_name dogstatsd --go_out $GOPATH/src/istio.io/istio/mixer/adapter/dogstatsd/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/logentry/template_proto.descriptor_set -p istio.io/istio/mixer/template/logentry --adapter_package istio.io/istio/mixer/adapter/fluentd --config_package istio.io/istio/mixer/adapter/fluentd/config --adapter_name fluentd --go_out $GOPATH/src/istio.io/istio/mixer/adapter/fluentd/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/listentry/template_proto.descriptor_set -p istio.io/istio/mixer/template/listentry --adapter_package istio.io/istio/mixer/adapter/list --config_package istio.io/istio/mixer/adapter/list/config --adapter_name listchecker --go_out $GOPATH/src/istio.io/istio/mixer/adapter/list/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/quota/template_proto.descriptor_set -p istio.io/istio/mixer/template/quota --adapter_package istio.io/istio/mixer/adapter/memquota --config_package istio.io/istio/mixer/adapter/memquota/config --adapter_name memquota --go_out $GOPATH/src/istio.io/istio/mixer/adapter/memquota/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/authorization/template_proto.descriptor_set -p istio.io/istio/mixer/template/authorization --adapter_package istio.io/istio/mixer/adapter/opa --config_package istio.io/istio/mixer/adapter/opa/config --adapter_name opa --go_out $GOPATH/src/istio.io/istio/mixer/adapter/opa/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -p istio.io/istio/mixer/template/metric --adapter_package istio.io/istio/mixer/adapter/prometheus --config_package istio.io/istio/mixer/adapter/prometheus/config --adapter_name prometheus --go_out $GOPATH/src/istio.io/istio/mixer/adapter/prometheus/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/authorization/template_proto.descriptor_set -p istio.io/istio/mixer/template/authorization --adapter_package istio.io/istio/mixer/adapter/rbac --config_package istio.io/istio/mixer/adapter/rbac/config --adapter_name rbac --go_out $GOPATH/src/istio.io/istio/mixer/adapter/rbac/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/apikey/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/quota/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/adapter/servicecontrol/template/servicecontrolreport/template_proto.descriptor_set -p istio.io/istio/mixer/template/quota -p istio.io/istio/mixer/template/apikey -p istio.io/istio/mixer/adapter/servicecontrol/template/servicecontrolreport --adapter_package istio.io/istio/mixer/adapter/servicecontrol --config_package istio.io/istio/mixer/adapter/servicecontrol/config --adapter_name servicecontrol --go_out $GOPATH/src/istio.io/istio/mixer/adapter/servicecontrol/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/tracespan/tracespan_proto.descriptor_set -p istio.io/istio/mixer/template/metric -p istio.io/istio/mixer/template/tracespan --adapter_package istio.io/istio/mixer/adapter/signalfx --adapter_name signalfx --go_out $GOPATH/src/istio.io/istio/mixer/adapter/signalfx/cmd --config_package istio.io/istio/mixer/adapter/signalfx/config

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/logentry/template_proto.descriptor_set -p istio.io/istio/mixer/template/metric -p istio.io/istio/mixer/template/logentry --adapter_package istio.io/istio/mixer/adapter/solarwinds --adapter_name solarwinds --go_out $GOPATH/src/istio.io/istio/mixer/adapter/solarwinds/cmd --config_package istio.io/istio/mixer/adapter/solarwinds/config

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/logentry/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/tracespan/tracespan_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/edge/template_proto.descriptor_set -p istio.io/istio/mixer/template/metric -p istio.io/istio/mixer/template/tracespan -p istio.io/istio/mixer/template/logentry -p istio.io/istio/mixer/template/edge --adapter_package istio.io/istio/mixer/adapter/stackdriver --adapter_name stackdriver --go_out $GOPATH/src/istio.io/istio/mixer/adapter/stackdriver/cmd --config_package istio.io/istio/mixer/adapter/stackdriver/config

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -p istio.io/istio/mixer/template/metric --adapter_package istio.io/istio/mixer/adapter/statsd --config_package istio.io/istio/mixer/adapter/statsd/config --adapter_name statsd --go_out $GOPATH/src/istio.io/istio/mixer/adapter/statsd/cmd

//go:generate go run $GOPATH/src/istio.io/istio/mixer/tools/mixgen/main.go server -t $GOPATH/src/istio.io/istio/mixer/template/metric/template_proto.descriptor_set -t $GOPATH/src/istio.io/istio/mixer/template/logentry/template_proto.descriptor_set -p istio.io/istio/mixer/template/metric -p istio.io/istio/mixer/template/logentry --adapter_package istio.io/istio/mixer/adapter/stdio --adapter_name stdio --go_out $GOPATH/src/istio.io/istio/mixer/adapter/stdio/cmd --config_package istio.io/istio/mixer/adapter/stdio/config

// Package adapter contains the inventory for all Mixer adapters that are compiled into
// a specific Mixer binary.
package adapter
