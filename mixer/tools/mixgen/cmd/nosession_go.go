package cmd

var noSessionServerTempl = `// Copyright 2018 Istio Authors
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

// THIS FILE IS AUTOMATICALLY GENERATED.

package server

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	adptModel "istio.io/api/mixer/adapter/model/v1beta1"
	"istio.io/api/policy/v1beta1"
	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/pkg/pool"
	"istio.io/istio/mixer/pkg/runtime/handler"
	{{- range .TemplatePackages}}
    "{{.}}"
	{{- end}}
	{{.AdapterName}} "{{.AdapterPackage}}"
	{{if ne .AdapterConfigPackage "" -}}
	config "{{.AdapterConfigPackage}}"
	{{end -}}
	$$additional_imports$$
)

type (
	// Server is basic server interface
	Server interface {
		Addr() string
		Close() error
		PromPort() int
		Run()
	}

	// NoSession models nosession adapter backend.
	NoSession struct {
		listener net.Listener
		shutdown chan error
		server   *grpc.Server

		rawcfg      []byte
		builder     adapter.HandlerBuilder
		env         adapter.Env
		builderLock sync.RWMutex
		{{range .Models}}
		{{.GoPackageName}}Handler {{.GoPackageName}}.Handler
		{{- end}}
	}
)
{{range .Models}}
var _ {{.GoPackageName}}.Handle{{.InterfaceName -}}ServiceServer = &NoSession{}
{{- end}}

func (s *NoSession) updateHandlers(rawcfg []byte) error {
	{{if ne .AdapterConfigPackage "" -}}
	cfg := &config.Params{}

	if err := cfg.Unmarshal(rawcfg); err != nil {
		return err
	}

	s.builderLock.Lock()
	defer s.builderLock.Unlock()
	if 0 == bytes.Compare(rawcfg, s.rawcfg) {
		return nil
	}

	s.env.Logger().Infof("Loaded handler with: %v", cfg)
	s.builder.SetAdapterConfig(cfg)
	{{end}}

	if ce := s.builder.Validate(); ce != nil {
		return ce
	}

	h, err := s.builder.Build(context.Background(), s.env)
	if err != nil {
		s.env.Logger().Errorf("could not build: %v", err)
		return err
	}
	s.rawcfg = rawcfg
	{{range .Models}}
	s.{{.GoPackageName}}Handler = h.({{.GoPackageName}}.Handler)
	{{- end}}
	return nil
}
{{range .Models}}
func (s *NoSession) get{{.InterfaceName -}}Handler(rawcfg []byte) ({{.GoPackageName}}.Handler, error) {
	s.builderLock.RLock()
	if 0 == bytes.Compare(rawcfg, s.rawcfg) {
		h := s.{{.GoPackageName}}Handler
		s.builderLock.RUnlock()
		return h, nil
	}
	s.builderLock.RUnlock()
	if err := s.updateHandlers(rawcfg); err != nil {
		return nil, err
	}

	// establish session
	return s.{{.GoPackageName}}Handler, nil
}
{{end}}

{{define "decode"}}
		{{range .Fields -}}
		{{if .ProtoType.IsResourceMessage -}}
		{{$gp := FindInterface $}}
		{{.GoName}}: {{ConstructDecodeFunc $gp .GoType.Name}}Msg(inst.{{.GoName}}),
		{{else if eq .ProtoType.Name "map<string, istio.policy.v1beta1.Value>" -}}
		{{.GoName}}: decodeDimensions(inst.{{.GoName}}),
		{{else if eq .ProtoType.Name "istio.policy.v1beta1.Value" -}}
		{{.GoName}}: decodeValue(inst.{{.GoName}}.GetValue()),
		{{else if eq .ProtoType.Name "istio.policy.v1beta1.TimeStamp" -}}
		{{.GoName}}: tmp{{.GoName}},
		{{else if eq .ProtoType.Name "istio.policy.v1beta1.Duration" -}}
		{{.GoName}}: tmp{{.GoName}},
		{{else -}}
		{{.GoName}}: inst.{{.GoName}},
		{{end -}}
		{{end -}}
{{end -}}

{{range .Models -}}
{{$tn := .InterfaceName -}}
{{$gp := .GoPackageName -}}
{{range .ResourceMessages -}}
func decode{{$tn}}{{.Name}}Msg(inst *{{$gp}}.{{.Name}}Msg) *{{$gp}}.{{.Name}} {
	{{range .Fields -}}
	{{if eq .ProtoType.Name "istio.policy.v1beta1.TimeStamp" -}}
	{{AddProtoToImpt -}}
	tmp{{.GoName}}, err := proto.TimestampFromProto(inst.{{.GoName}}.GetValue())
	if err != nil {
		continue
	}
	{{else if eq .ProtoType.Name "istio.policy.v1beta1.Duration" -}}
	{{AddProtoToImpt -}}
	tmp{{.GoName}}, err := proto.DurationFromProto(inst.{{.GoName}}.GetValue())
	if err != nil {
		continue
	}
	{{end -}}
	{{end -}}
	return &{{$gp}}.{{.Name}}{
		{{template "decode" .}}
	}
}
{{end -}}
{{end -}}

{{range .Models}}
{{if eq .VarietyName "TEMPLATE_VARIETY_REPORT" -}}
func {{.GoPackageName}}Instances(in []*{{.GoPackageName}}.InstanceMsg) []*{{.GoPackageName}}.Instance {
	out := make([]*{{.GoPackageName}}.Instance, 0, len(in))

	for _, inst := range in {
		{{range .TemplateMessage.Fields -}}
		{{if eq .ProtoType.Name "istio.policy.v1beta1.TimeStamp" -}}
		{{AddProtoToImpt -}}
		tmp{{.GoName}}, err := proto.TimestampFromProto(inst.{{.GoName}}.GetValue())
		if err != nil {
			continue
		}
		{{else if eq .ProtoType.Name "istio.policy.v1beta1.Duration" -}}
		{{AddProtoToImpt -}}
		tmp{{.GoName}}, err := proto.DurationFromProto(inst.{{.GoName}}.GetValue())
		if err != nil {
			continue
		}
		{{end -}}
		{{end -}}
		out = append(out, &{{.GoPackageName}}.Instance{
			Name: inst.Name,
			{{template "decode" .TemplateMessage}}
		})
	}
	return out
}
{{else if or (eq .VarietyName "TEMPLATE_VARIETY_CHECK") (eq .VarietyName "TEMPLATE_VARIETY_QUOTA") -}}
func {{.GoPackageName}}Instance(inst *{{.GoPackageName}}.InstanceMsg) *{{.GoPackageName}}.Instance {
	{{range .TemplateMessage.Fields -}}
	{{if eq .ProtoType.Name "istio.policy.v1beta1.TimeStamp" -}}
	{{AddProtoToImpt -}}
	tmp{{.GoName}}, err := proto.TimestampFromProto(inst.{{.GoName}}.GetValue())
	if err != nil {
		return nil
	}
	{{else if eq .ProtoType.Name "istio.policy.v1beta1.Duration" -}}
	{{AddProtoToImpt -}}
	tmp{{.GoName}}, err := proto.DurationFromProto(inst.{{.GoName}}.GetValue())
	if err != nil {
		continue
	}
	{{end -}}
	{{end -}}
	return &{{.GoPackageName}}.Instance{
		Name: inst.Name,
		{{template "decode" .TemplateMessage}}
	}
}
{{end -}}
{{end}}

func decodeDimensions(in map[string]*v1beta1.Value) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = decodeValue(v.GetValue())
	}
	return out
}

func decodeValue(in interface{}) interface{} {
	switch t := in.(type) {
	case *v1beta1.Value_StringValue:
		return t.StringValue
	case *v1beta1.Value_Int64Value:
		return t.Int64Value
	case *v1beta1.Value_DoubleValue:
		return t.DoubleValue
	default:
		return fmt.Sprintf("%v", in)
	}
}

{{range .Models}}
// Handle{{.InterfaceName}} ...
{{if eq .VarietyName "TEMPLATE_VARIETY_REPORT" -}}
func (s *NoSession) Handle{{.InterfaceName -}}(ctx context.Context, r *{{.GoPackageName}}.Handle{{.InterfaceName -}}Request) (*adptModel.ReportResult, error) {
	h, err := s.get{{.InterfaceName -}}Handler(r.AdapterConfig.Value)
	if err != nil {
		return nil, err
	}

	if err = h.Handle{{.InterfaceName -}}(ctx, {{.GoPackageName}}Instances(r.Instances)); err != nil {
		s.env.Logger().Errorf("Could not process: %v", err)
		return nil, err
	}

	return &adptModel.ReportResult{}, nil
}
{{else if eq .VarietyName "TEMPLATE_VARIETY_CHECK" -}}
func (s *NoSession) Handle{{.InterfaceName -}}(ctx context.Context, r *{{.GoPackageName}}.Handle{{.InterfaceName -}}Request) (*adptModel.CheckResult, error) {
	h, err := s.get{{.InterfaceName -}}Handler(r.AdapterConfig.Value)
	if err != nil {
		return nil, err
	}
	inst := {{.GoPackageName}}Instance(r.Instance)
	if inst == nil {
		return nil, fmt.Errorf("cannot transform instance")
	}
	cr, err := h.Handle{{.InterfaceName -}}(ctx, inst)
	if err != nil {
		s.env.Logger().Errorf("Could not process: %v", err)
		return nil, err
	}
	return &adptModel.CheckResult{
		Status:        cr.Status,
		ValidDuration: cr.ValidDuration,
		ValidUseCount: cr.ValidUseCount,
	}, nil
}
{{else if eq .VarietyName "TEMPLATE_VARIETY_QUOTA" -}}
func (s *NoSession) Handle{{.InterfaceName -}}(ctx context.Context, r *{{.GoPackageName}}.Handle{{.InterfaceName -}}Request) (*adptModel.QuotaResult, error) {
	h, err := s.get{{.InterfaceName -}}Handler(r.AdapterConfig.Value)
	if err != nil {
		return nil, err
	}

	qi := {{.GoPackageName}}Instance(r.Instance)
	resp := adptModel.QuotaResult{
		Quotas: make(map[string]adptModel.QuotaResult_Result),
	}
	for qt, p := range r.QuotaRequest.Quotas {
		qa := adapter.QuotaArgs{
			DeduplicationID: r.DedupId,
			QuotaAmount:     p.Amount,
			BestEffort:      p.BestEffort,
		}
		qr, err := h.Handle{{.InterfaceName -}}(ctx, qi, qa)
		if err != nil {
			return nil, err
		}
		resp.Quotas[qt] = adptModel.QuotaResult_Result{
			ValidDuration: qr.ValidDuration,
			GrantedAmount: qr.Amount,
		}
	}
	if err != nil {
		s.env.Logger().Errorf("Could not process: %v", err)
		return nil, err
	}
	return &resp, nil
}
{{end}}
{{end}}

// Addr returns the listening address of the server
func (s *NoSession) Addr() string {
	return s.listener.Addr().String()
}

// Run starts the server run
func (s *NoSession) Run() {
	s.shutdown = make(chan error, 1)
	go func() {
		err := s.server.Serve(s.listener)

		// notify closer we're done
		s.shutdown <- err
	}()
}

// Wait waits for server to stop
func (s *NoSession) Wait() error {
	if s.shutdown == nil {
		return fmt.Errorf("server not running")
	}

	err := <-s.shutdown
	s.shutdown = nil
	return err
}

// Close gracefully shuts down the server
func (s *NoSession) Close() error {
	if s.shutdown != nil {
		s.server.GracefulStop()
		_ = s.Wait()
	}

	if s.listener != nil {
		_ = s.listener.Close()
	}

	return nil
}

// New{{Capitalize .AdapterName}}NoSessionServer creates a new no session server from given args.
func New{{Capitalize .AdapterName}}NoSessionServer(addr uint16, poolSize int) (*NoSession, error) {
	saddr := fmt.Sprintf(":%d", addr)

	gp := pool.NewGoroutinePool(poolSize, false)
	inf := {{.AdapterName}}.GetInfo()
	s := &NoSession{
		builder: inf.NewBuilder(),
		env:     handler.NewEnv(0, "{{.AdapterName}}-nosession", gp),
		rawcfg:  []byte{0xff, 0xff},
	}
	var err error
	if s.listener, err = net.Listen("tcp", saddr); err != nil {
		_ = s.Close()
		return nil, fmt.Errorf("unable to listen on socket: %v", err)
	}

	fmt.Printf("listening on :%v\n", s.listener.Addr())
	s.server = grpc.NewServer()
	{{range .Models}}
	{{.GoPackageName}}.RegisterHandle{{.InterfaceName -}}ServiceServer(s.server, s)
	{{- end}}

	return s, nil
}
`
