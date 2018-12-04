// Copyright 2018 Istio Authors
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
	"context"
	"fmt"
	"net"
	"sync"

	proto "github.com/gogo/protobuf/types"
	"google.golang.org/grpc"

	adptModel "istio.io/api/mixer/adapter/model/v1beta1"
	"istio.io/api/policy/v1beta1"
	stackdriver "istio.io/istio/mixer/adapter/stackdriver"
	config "istio.io/istio/mixer/adapter/stackdriver/config"
	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/pkg/pool"
	"istio.io/istio/mixer/pkg/runtime/handler"
	"istio.io/istio/mixer/template/edge"
	"istio.io/istio/mixer/template/logentry"
	"istio.io/istio/mixer/template/metric"
	"istio.io/istio/mixer/template/tracespan"
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

		builder     adapter.HandlerBuilder
		env         adapter.Env
		builderLock sync.RWMutex
		handlerMap  map[string]adapter.Handler
	}
)

var _ metric.HandleMetricServiceServer = &NoSession{}
var _ logentry.HandleLogEntryServiceServer = &NoSession{}
var _ tracespan.HandleTraceSpanServiceServer = &NoSession{}
var _ edge.HandleEdgeServiceServer = &NoSession{}

func (s *NoSession) updateHandlers(rawcfg []byte) (adapter.Handler, error) {
	cfg := &config.Params{}

	if err := cfg.Unmarshal(rawcfg); err != nil {
		return nil, err
	}

	s.builderLock.Lock()
	defer s.builderLock.Unlock()
	if handler, ok := s.handlerMap[string(rawcfg)]; ok {
		return handler, nil
	}

	s.env.Logger().Infof("Loaded handler with: %v", cfg)
	s.builder.SetAdapterConfig(cfg)

	if ce := s.builder.Validate(); ce != nil {
		return nil, ce
	}

	h, err := s.builder.Build(context.Background(), s.env)
	if err != nil {
		s.env.Logger().Errorf("could not build: %v", err)
		return nil, err
	}
	s.handlerMap[string(rawcfg)] = h
	return h, nil
}

func (s *NoSession) getMetricHandler(rawcfg []byte) (metric.Handler, error) {
	s.builderLock.RLock()
	if handler, ok := s.handlerMap[string(rawcfg)]; ok {
		h := handler.(metric.Handler)
		s.builderLock.RUnlock()
		return h, nil
	}
	s.builderLock.RUnlock()
	h, err := s.updateHandlers(rawcfg)
	if err != nil {
		return nil, err
	}

	// establish session
	return h.(metric.Handler), nil
}

func (s *NoSession) getLogEntryHandler(rawcfg []byte) (logentry.Handler, error) {
	s.builderLock.RLock()
	if handler, ok := s.handlerMap[string(rawcfg)]; ok {
		h := handler.(logentry.Handler)
		s.builderLock.RUnlock()
		return h, nil
	}
	s.builderLock.RUnlock()
	h, err := s.updateHandlers(rawcfg)
	if err != nil {
		return nil, err
	}

	// establish session
	return h.(logentry.Handler), nil
}

func (s *NoSession) getTraceSpanHandler(rawcfg []byte) (tracespan.Handler, error) {
	s.builderLock.RLock()
	if handler, ok := s.handlerMap[string(rawcfg)]; ok {
		h := handler.(tracespan.Handler)
		s.builderLock.RUnlock()
		return h, nil
	}
	s.builderLock.RUnlock()
	h, err := s.updateHandlers(rawcfg)
	if err != nil {
		return nil, err
	}

	// establish session
	return h.(tracespan.Handler), nil
}

func (s *NoSession) getEdgeHandler(rawcfg []byte) (edge.Handler, error) {
	s.builderLock.RLock()
	if handler, ok := s.handlerMap[string(rawcfg)]; ok {
		h := handler.(edge.Handler)
		s.builderLock.RUnlock()
		return h, nil
	}
	s.builderLock.RUnlock()
	h, err := s.updateHandlers(rawcfg)
	if err != nil {
		return nil, err
	}

	// establish session
	return h.(edge.Handler), nil
}

func metricInstances(in []*metric.InstanceMsg) []*metric.Instance {
	out := make([]*metric.Instance, 0, len(in))

	for _, inst := range in {
		out = append(out, &metric.Instance{
			Name: inst.Name,

			Value:                       transformValue(inst.Value.GetValue()),
			Dimensions:                  transformValueMap(inst.Dimensions),
			MonitoredResourceType:       inst.MonitoredResourceType,
			MonitoredResourceDimensions: transformValueMap(inst.MonitoredResourceDimensions),
		})
	}
	return out
}

func logentryInstances(in []*logentry.InstanceMsg) []*logentry.Instance {
	out := make([]*logentry.Instance, 0, len(in))

	for _, inst := range in {
		tmpTimestamp, err := proto.TimestampFromProto(inst.Timestamp.GetValue())
		if err != nil {
			continue
		}
		out = append(out, &logentry.Instance{
			Name: inst.Name,

			Variables:                   transformValueMap(inst.Variables),
			Timestamp:                   tmpTimestamp,
			Severity:                    inst.Severity,
			MonitoredResourceType:       inst.MonitoredResourceType,
			MonitoredResourceDimensions: transformValueMap(inst.MonitoredResourceDimensions),
		})
	}
	return out
}

func tracespanInstances(in []*tracespan.InstanceMsg) []*tracespan.Instance {
	out := make([]*tracespan.Instance, 0, len(in))

	for _, inst := range in {
		tmpStartTime, err := proto.TimestampFromProto(inst.StartTime.GetValue())
		if err != nil {
			continue
		}
		tmpEndTime, err := proto.TimestampFromProto(inst.EndTime.GetValue())
		if err != nil {
			continue
		}
		out = append(out, &tracespan.Instance{
			Name: inst.Name,

			TraceId:             inst.TraceId,
			SpanId:              inst.SpanId,
			ParentSpanId:        inst.ParentSpanId,
			SpanName:            inst.SpanName,
			StartTime:           tmpStartTime,
			EndTime:             tmpEndTime,
			SpanTags:            transformValueMap(inst.SpanTags),
			HttpStatusCode:      inst.HttpStatusCode,
			ClientSpan:          inst.ClientSpan,
			RewriteClientSpanId: inst.RewriteClientSpanId,
		})
	}
	return out
}

func edgeInstances(in []*edge.InstanceMsg) []*edge.Instance {
	out := make([]*edge.Instance, 0, len(in))

	for _, inst := range in {
		tmpTimestamp, err := proto.TimestampFromProto(inst.Timestamp.GetValue())
		if err != nil {
			continue
		}
		out = append(out, &edge.Instance{
			Name: inst.Name,

			Timestamp:                    tmpTimestamp,
			SourceWorkloadNamespace:      inst.SourceWorkloadNamespace,
			SourceWorkloadName:           inst.SourceWorkloadName,
			SourceOwner:                  inst.SourceOwner,
			SourceUid:                    inst.SourceUid,
			DestinationWorkloadNamespace: inst.DestinationWorkloadNamespace,
			DestinationWorkloadName:      inst.DestinationWorkloadName,
			DestinationOwner:             inst.DestinationOwner,
			DestinationUid:               inst.DestinationUid,
			ContextProtocol:              inst.ContextProtocol,
			ApiProtocol:                  inst.ApiProtocol,
		})
	}
	return out
}

func transformValueMap(in map[string]*v1beta1.Value) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = transformValue(v.GetValue())
	}
	return out
}

func transformValue(in interface{}) interface{} {
	switch t := in.(type) {
	case *v1beta1.Value_StringValue:
		return t.StringValue
	case *v1beta1.Value_Int64Value:
		return t.Int64Value
	case *v1beta1.Value_DoubleValue:
		return t.DoubleValue
	case *v1beta1.Value_BoolValue:
		return t.BoolValue
	case *v1beta1.Value_IpAddressValue:
		return t.IpAddressValue.Value
	case *v1beta1.Value_EmailAddressValue:
		return t.EmailAddressValue.Value
	case *v1beta1.Value_UriValue:
		return t.UriValue.Value
	default:
		return fmt.Sprintf("%v", in)
	}
}

// HandleMetric handles 'Metric' instances.
func (s *NoSession) HandleMetric(ctx context.Context, r *metric.HandleMetricRequest) (*adptModel.ReportResult, error) {
	h, err := s.getMetricHandler(r.AdapterConfig.Value)
	if err != nil {
		return nil, err
	}

	if err = h.HandleMetric(ctx, metricInstances(r.Instances)); err != nil {
		s.env.Logger().Errorf("Could not process: %v", err)
		return nil, err
	}

	return &adptModel.ReportResult{}, nil
}

// HandleLogEntry handles 'LogEntry' instances.
func (s *NoSession) HandleLogEntry(ctx context.Context, r *logentry.HandleLogEntryRequest) (*adptModel.ReportResult, error) {
	h, err := s.getLogEntryHandler(r.AdapterConfig.Value)
	if err != nil {
		return nil, err
	}

	if err = h.HandleLogEntry(ctx, logentryInstances(r.Instances)); err != nil {
		s.env.Logger().Errorf("Could not process: %v", err)
		return nil, err
	}

	return &adptModel.ReportResult{}, nil
}

// HandleTraceSpan handles 'TraceSpan' instances.
func (s *NoSession) HandleTraceSpan(ctx context.Context, r *tracespan.HandleTraceSpanRequest) (*adptModel.ReportResult, error) {
	h, err := s.getTraceSpanHandler(r.AdapterConfig.Value)
	if err != nil {
		return nil, err
	}

	if err = h.HandleTraceSpan(ctx, tracespanInstances(r.Instances)); err != nil {
		s.env.Logger().Errorf("Could not process: %v", err)
		return nil, err
	}

	return &adptModel.ReportResult{}, nil
}

// HandleEdge handles 'Edge' instances.
func (s *NoSession) HandleEdge(ctx context.Context, r *edge.HandleEdgeRequest) (*adptModel.ReportResult, error) {
	h, err := s.getEdgeHandler(r.AdapterConfig.Value)
	if err != nil {
		return nil, err
	}

	if err = h.HandleEdge(ctx, edgeInstances(r.Instances)); err != nil {
		s.env.Logger().Errorf("Could not process: %v", err)
		return nil, err
	}

	return &adptModel.ReportResult{}, nil
}

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

// NewStackdriverNoSessionServer creates a new no session server based on given args.
func NewStackdriverNoSessionServer(addr uint16, poolSize int) (*NoSession, error) {
	saddr := fmt.Sprintf(":%d", addr)

	gp := pool.NewGoroutinePool(poolSize, false)
	inf := stackdriver.GetInfo()
	s := &NoSession{
		builder:    inf.NewBuilder(),
		env:        handler.NewEnv(0, "stackdriver-nosession", gp),
		handlerMap: make(map[string]adapter.Handler),
	}
	var err error
	if s.listener, err = net.Listen("tcp", saddr); err != nil {
		_ = s.Close()
		return nil, fmt.Errorf("unable to listen on socket: %v", err)
	}

	fmt.Printf("listening on :%v\n", s.listener.Addr())
	s.server = grpc.NewServer()

	metric.RegisterHandleMetricServiceServer(s.server, s)
	logentry.RegisterHandleLogEntryServiceServer(s.server, s)
	tracespan.RegisterHandleTraceSpanServiceServer(s.server, s)
	edge.RegisterHandleEdgeServiceServer(s.server, s)

	return s, nil
}
