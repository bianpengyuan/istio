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
	"bytes"
	"context"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	adptModel "istio.io/api/mixer/adapter/model/v1beta1"
	"istio.io/api/policy/v1beta1"
	rbac "istio.io/istio/mixer/adapter/rbac"
	config "istio.io/istio/mixer/adapter/rbac/config"
	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/pkg/pool"
	"istio.io/istio/mixer/pkg/runtime/handler"
	"istio.io/istio/mixer/template/authorization"
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

		authorizationHandler authorization.Handler
	}
)

var _ authorization.HandleAuthorizationServiceServer = &NoSession{}

func (s *NoSession) updateHandlers(rawcfg []byte) error {
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

	if ce := s.builder.Validate(); ce != nil {
		return ce
	}

	h, err := s.builder.Build(context.Background(), s.env)
	if err != nil {
		s.env.Logger().Errorf("could not build: %v", err)
		return err
	}
	s.rawcfg = rawcfg

	s.authorizationHandler = h.(authorization.Handler)
	return nil
}

func (s *NoSession) getAuthorizationHandler(rawcfg []byte) (authorization.Handler, error) {
	s.builderLock.RLock()
	if 0 == bytes.Compare(rawcfg, s.rawcfg) {
		h := s.authorizationHandler
		s.builderLock.RUnlock()
		return h, nil
	}
	s.builderLock.RUnlock()
	if err := s.updateHandlers(rawcfg); err != nil {
		return nil, err
	}

	// establish session
	return s.authorizationHandler, nil
}

func decodeAuthorizationSubjectMsg(inst *authorization.SubjectMsg) *authorization.Subject {
	return &authorization.Subject{

		User:       inst.User,
		Groups:     inst.Groups,
		Properties: decodeDimensions(inst.Properties),
	}
}
func decodeAuthorizationActionMsg(inst *authorization.ActionMsg) *authorization.Action {
	return &authorization.Action{

		Namespace:  inst.Namespace,
		Service:    inst.Service,
		Method:     inst.Method,
		Path:       inst.Path,
		Properties: decodeDimensions(inst.Properties),
	}
}

func authorizationInstance(inst *authorization.InstanceMsg) *authorization.Instance {
	return &authorization.Instance{
		Name: inst.Name,

		Subject: decodeAuthorizationSubjectMsg(inst.Subject),

		Action: decodeAuthorizationActionMsg(inst.Action),
	}
}

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

// HandleAuthorization ...
func (s *NoSession) HandleAuthorization(ctx context.Context, r *authorization.HandleAuthorizationRequest) (*adptModel.CheckResult, error) {
	h, err := s.getAuthorizationHandler(r.AdapterConfig.Value)
	if err != nil {
		return nil, err
	}
	inst := authorizationInstance(r.Instance)
	if inst == nil {
		return nil, fmt.Errorf("cannot transform instance")
	}
	cr, err := h.HandleAuthorization(ctx, inst)
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

// NewRbacNoSessionServer creates a new no session server from given args.
func NewRbacNoSessionServer(addr uint16, poolSize int) (*NoSession, error) {
	saddr := fmt.Sprintf(":%d", addr)

	gp := pool.NewGoroutinePool(poolSize, false)
	inf := rbac.GetInfo()
	s := &NoSession{
		builder: inf.NewBuilder(),
		env:     handler.NewEnv(0, "rbac-nosession", gp),
		rawcfg:  []byte{0xff, 0xff},
	}
	var err error
	if s.listener, err = net.Listen("tcp", saddr); err != nil {
		_ = s.Close()
		return nil, fmt.Errorf("unable to listen on socket: %v", err)
	}

	fmt.Printf("listening on :%v\n", s.listener.Addr())
	s.server = grpc.NewServer()

	authorization.RegisterHandleAuthorizationServiceServer(s.server, s)

	return s, nil
}
