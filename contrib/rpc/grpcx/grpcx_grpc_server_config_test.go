package grpcx

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Grpcx_Grpc_Server_Config(t *testing.T) {
	cfg := Server.NewConfig()
	addr := "10.0.0.29:80"
	cfg.Endpoints = []string{
		addr,
	}
	// cfg set one endpoint
	gtest.C(t, func(t *gtest.T) {
		s := Server.New(cfg)
		s.doServiceRegister()
		for _, svc := range s.services {
			t.Assert(svc.GetEndpoints().String(), addr)
		}
	})
	// cfg set more endpoints
	addr = "10.0.0.29:80,10.0.0.29:81"
	cfg.Endpoints = []string{
		"10.0.0.29:80",
		"10.0.0.29:81",
	}
	gtest.C(t, func(t *gtest.T) {
		s := Server.New(cfg)
		s.doServiceRegister()
		for _, svc := range s.services {
			t.Assert(svc.GetEndpoints().String(), addr)
		}
	})
}
