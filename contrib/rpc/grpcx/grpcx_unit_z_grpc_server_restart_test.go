// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests restarting a managed gRPC server after it was stopped or rolled back.

package grpcx

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gapp"
	"github.com/gogf/gf/v2/test/gtest"
)

// failFirstStartServer fails its first Start so gapp.App rolls back the servers that did start.
type failFirstStartServer struct {
	starts int32
}

// Start fails on the first call and succeeds afterwards.
func (s *failFirstStartServer) Start() error {
	if atomic.AddInt32(&s.starts, 1) == 1 {
		return errors.New("bind failed")
	}
	return nil
}

// Stop implements gapp.Server.
func (s *failFirstStartServer) Stop(bool) error {
	return nil
}

// assertHealthServing asserts that the health service at address reports SERVING.
func assertHealthServing(t *gtest.T, address string) {
	client, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.AssertNil(err)
	defer func() { t.AssertNil(client.Close()) }()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	response, err := grpc_health_v1.NewHealthClient(client).Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	t.AssertNil(err)
	t.Assert(response.GetStatus(), grpc_health_v1.HealthCheckResponse_SERVING)
}

// Test_StartManaged_RestartAfterStopReplaysRegistrations serves recorded services on a fresh grpc.Server.
func Test_StartManaged_RestartAfterStopReplaysRegistrations(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		graceful bool
	}{
		{name: "graceful", graceful: true},
		{name: "forceful", graceful: false},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var buf safeBuffer
				s := newServeErrorTestServer(&buf)
				grpc_health_v1.RegisterHealthServer(s, health.NewServer())
				adapter := NewGappServerAdapter(s)

				t.AssertNil(adapter.Start())
				assertHealthServing(t, s.GetListenedAddress())
				stoppedServer := s.Server

				t.AssertNil(adapter.Stop(testCase.graceful))
				t.Assert(s.GetListenedPort(), -1)

				t.AssertNil(adapter.Start())
				defer func() { t.AssertNil(adapter.Stop(false)) }()
				t.Assert(s.Server != stoppedServer, true)
				assertHealthServing(t, s.GetListenedAddress())
			})
		})
	}
}

// Test_StartManaged_RestartRejectsDirectRegistrations fails instead of serving without unrecorded services.
func Test_StartManaged_RestartRejectsDirectRegistrations(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var buf safeBuffer
		s := newServeErrorTestServer(&buf)
		grpc_health_v1.RegisterHealthServer(s.Server, health.NewServer())

		t.AssertNil(s.StartManaged())
		s.StopForceful()
		t.Assert(s.GetListenedPort(), -1)

		err := s.StartManaged()
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidOperation)
		t.Assert(strings.Contains(err.Error(), grpc_health_v1.Health_ServiceDesc.ServiceName), true)
		t.Assert(s.GetListenedPort(), -1)
	})
}

// Test_StartManaged_StopBeforeServeKeepsServer leaves a never-served grpc.Server and its direct registrations usable.
func Test_StartManaged_StopBeforeServeKeepsServer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var buf safeBuffer
		s := newServeErrorTestServer(&buf)
		grpc_health_v1.RegisterHealthServer(s.Server, health.NewServer())
		server := s.Server
		s.Stop()
		s.StopForceful()

		t.AssertNil(s.StartManaged())
		defer s.StopForceful()
		t.Assert(s.Server == server, true)
		assertHealthServing(t, s.GetListenedAddress())
	})
}

// Test_Grpcx_App_RestartAfterStartRollback brings gRPC back after App.Start rolls it back for a sibling failure.
func Test_Grpcx_App_RestartAfterStartRollback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var buf safeBuffer
		s := newServeErrorTestServer(&buf)
		grpc_health_v1.RegisterHealthServer(s, health.NewServer())
		app := gapp.New(NewGappServerAdapter(s), &failFirstStartServer{})

		err := app.Start(context.Background())
		t.AssertNE(err, nil)
		t.Assert(strings.Contains(err.Error(), "bind failed"), true)
		t.Assert(s.GetListenedPort(), -1)

		t.AssertNil(app.Start(context.Background()))
		defer func() { t.AssertNil(app.Stop(context.Background(), true)) }()
		assertHealthServing(t, s.GetListenedAddress())
	})
}
