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

package noopbackend

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	rpc "github.com/gogo/googleapis/google/rpc"
	adptModel "istio.io/api/mixer/adapter/model/v1beta1"
	"istio.io/istio/mixer/template/authorization"
)

type (
	// Server is basic server interface
	Server interface {
		Addr() net.Addr
		Close() error
		Run()
	}

	// NoSessionServer models no session adapter backend.
	NoSessionServer struct {
		listener net.Listener
		shutdown chan error
		server   *grpc.Server
	}
)

var _ authorization.HandleAuthorizationServiceServer = &NoSessionServer{}

// HandleAuthorization records authorization entries and responds with the programmed response
func (s *NoSessionServer) HandleAuthorization(c context.Context, r *authorization.HandleAuthorizationRequest) (*adptModel.CheckResult, error) {
	return &adptModel.CheckResult{
		Status: rpc.Status{
			Code: int32(rpc.OK),
		},
		ValidDuration: 1000000000 * time.Second,
		ValidUseCount: 1000000000}, nil
}

// Addr returns the listening address of the server
func (s *NoSessionServer) Addr() net.Addr {
	return s.listener.Addr()
}

// Run starts the server run
func (s *NoSessionServer) Run() {
	s.shutdown = make(chan error, 1)
	go func() {
		err := s.server.Serve(s.listener)

		// notify closer we're done
		s.shutdown <- err
	}()
}

// Wait waits for server to stop
func (s *NoSessionServer) Wait() error {
	if s.shutdown == nil {
		return fmt.Errorf("server not running")
	}

	err := <-s.shutdown
	s.shutdown = nil
	return err
}

// Close gracefully shuts down the server
func (s *NoSessionServer) Close() error {
	if s.shutdown != nil {
		s.server.GracefulStop()
		_ = s.Wait()
	}

	if s.listener != nil {
		_ = s.listener.Close()
	}

	return nil
}

// NewNoSessionServer creates a new no session server from given args.
func NewNoSessionServer(addr string) (Server, error) {
	if addr == "" {
		addr = "0"
	}
	s := &NoSessionServer{}
	var err error

	if s.listener, err = net.Listen("tcp", fmt.Sprintf(":%s", addr)); err != nil {
		_ = s.Close()
		return nil, fmt.Errorf("unable to listen on socket: %v", err)
	}

	fmt.Printf("listening on :%v", s.listener.Addr())

	s.server = grpc.NewServer()
	authorization.RegisterHandleAuthorizationServiceServer(s.server, s)

	return s, nil
}
