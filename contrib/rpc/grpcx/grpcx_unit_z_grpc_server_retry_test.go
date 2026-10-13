// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests registration rollback and real gRPC serving after a failed managed start.

package grpcx

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gapp"
	"github.com/gogf/gf/v2/test/gtest"
)

// retryRegistrar fails its selected registration attempt and records rollback.
type retryRegistrar struct {
	failAt        int
	names         []string
	deregistered  []string
	failedAddress string
}

// retryRegisteredService models registry-specific handles returned by Register.
type retryRegisteredService struct {
	gsvc.Service
}

// Register records the service and fails the selected attempt.
func (r *retryRegistrar) Register(_ context.Context, service gsvc.Service) (gsvc.Service, error) {
	r.names = append(r.names, service.GetName())
	if len(r.names) == r.failAt {
		r.failedAddress = service.GetEndpoints().String()
		return nil, errors.New("registration unavailable")
	}
	return &retryRegisteredService{Service: service}, nil
}

// Deregister records services removed during rollback or shutdown.
func (r *retryRegistrar) Deregister(_ context.Context, service gsvc.Service) error {
	// Registries such as Polaris require the handle returned by Register.
	if _, ok := service.(*retryRegisteredService); !ok {
		panic("Deregister requires the registered service")
	}
	r.deregistered = append(r.deregistered, service.GetName())
	return nil
}

// Test_Grpcx_AppRun_RegisterFailureCleanup verifies Run never deregisters raw service definitions.
func Test_Grpcx_AppRun_RegisterFailureCleanup(t *testing.T) {
	for _, testCase := range []struct {
		name   string
		failAt int
	}{
		{name: "first_registration", failAt: 1},
		{name: "partial_registration", failAt: 2},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					buf       safeBuffer
					s         = newServeErrorTestServer(&buf)
					registrar = &retryRegistrar{failAt: testCase.failAt}
				)
				s.registrar = registrar
				s.Service(
					&gsvc.LocalService{Name: "first", Metadata: gsvc.Metadata{}},
					&gsvc.LocalService{Name: "second", Metadata: gsvc.Metadata{}},
				)
				app := gapp.New(NewGappServerAdapter(s))
				err := app.Run(context.Background())
				t.AssertNE(err, nil)
				t.Assert(strings.Contains(err.Error(), "registration unavailable"), true)
				t.Assert(len(registrar.deregistered), testCase.failAt-1)
				t.Assert(s.GetListenedPort(), -1)
			})
		})
	}
}

// Test_Grpcx_Stop_RegisteredServices verifies shutdown uses registration handles exactly once.
func Test_Grpcx_Stop_RegisteredServices(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		graceful bool
	}{
		{name: "graceful", graceful: true},
		{name: "forceful", graceful: false},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			for _, started := range []bool{false, true} {
				name := "before_start"
				if started {
					name = "after_start"
				}
				t.Run(name, func(t *testing.T) {
					gtest.C(t, func(t *gtest.T) {
						var (
							buf       safeBuffer
							s         = newServeErrorTestServer(&buf)
							registrar = &retryRegistrar{}
						)
						s.registrar = registrar
						s.Service(&gsvc.LocalService{Name: "first", Metadata: gsvc.Metadata{}})
						if started {
							t.AssertNil(s.StartManaged())
						}
						adapter := NewGappServerAdapter(s)
						t.AssertNil(adapter.Stop(testCase.graceful))
						t.AssertNil(adapter.Stop(testCase.graceful))
						if started {
							t.Assert(registrar.deregistered, []string{"first"})
						} else {
							t.Assert(len(registrar.deregistered), 0)
						}
					})
				})
			}
		})
	}
}

// Test_StartManaged_RetryAfterRegisterFailure verifies listener rollback and retained RPC bindings.
func Test_StartManaged_RetryAfterRegisterFailure(t *testing.T) {
	for _, testCase := range []struct {
		name   string
		failAt int
	}{
		{name: "first_registration", failAt: 1},
		{name: "partial_registration", failAt: 2},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					buf       safeBuffer
					s         = newServeErrorTestServer(&buf)
					registrar = &retryRegistrar{failAt: testCase.failAt}
				)
				s.registrar = registrar
				s.Service(
					&gsvc.LocalService{Name: "first", Metadata: gsvc.Metadata{}},
					&gsvc.LocalService{Name: "second", Metadata: gsvc.Metadata{}},
				)
				grpc_health_v1.RegisterHealthServer(s.Server, health.NewServer())
				defer s.StopForceful()

				err := s.StartManaged()
				t.AssertNE(err, nil)
				t.Assert(err.Error(), "registration unavailable")
				t.Assert(s.GetListenedPort(), -1)
				t.Assert(len(registrar.deregistered), testCase.failAt-1)
				if testCase.failAt == 2 {
					t.Assert(registrar.deregistered, []string{"first"})
				}
				conn, dialErr := net.DialTimeout("tcp", registrar.failedAddress, time.Second)
				if conn != nil {
					t.AssertNil(conn.Close())
				}
				t.AssertNE(dialErr, nil)

				t.AssertNil(s.StartManaged())
				t.Assert(registrar.names[testCase.failAt:], []string{"first", "second"})
				client, clientErr := grpc.NewClient(s.GetListenedAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
				t.AssertNil(clientErr)
				defer func() { t.AssertNil(client.Close()) }()
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				response, rpcErr := grpc_health_v1.NewHealthClient(client).Check(ctx, &grpc_health_v1.HealthCheckRequest{})
				t.AssertNil(rpcErr)
				t.Assert(response.GetStatus(), grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	}
}
