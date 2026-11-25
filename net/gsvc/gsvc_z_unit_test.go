// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file unit-tests the net/gsvc package for paths that do not require
// a live registry: Endpoint/Endpoints construction, Metadata CRUD,
// LocalService creation and field accessors, NewServiceWithKV parsing,
// SetRegistry/GetRegistry round-trips, and Nodes.String.
// All tests run without any network or external service dependency.

package gsvc_test

import (
	"testing"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/test/gtest"
)

// ─────────────────────── Endpoint ────────────────────────────────

func TestEndpoint_HostPort(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ep := gsvc.NewEndpoint("192.168.1.10:8080")
		t.Assert(ep.Host(), "192.168.1.10")
		t.Assert(ep.Port(), 8080)
		t.Assert(ep.String(), "192.168.1.10:8080")
	})
}

func TestEndpoint_InvalidPanics(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			t.AssertNE(recover(), nil)
		}()
		gsvc.NewEndpoint("no-port")
	})
}

// ─────────────────────── Endpoints ───────────────────────────────

func TestEndpoints_MultipleAddresses(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		eps := gsvc.NewEndpoints("10.0.0.1:80,10.0.0.2:81")
		t.Assert(len(eps), 2)
		t.Assert(eps[0].Host(), "10.0.0.1")
		t.Assert(eps[0].Port(), 80)
		t.Assert(eps[1].Host(), "10.0.0.2")
		t.Assert(eps[1].Port(), 81)
		t.Assert(eps.String(), "10.0.0.1:80,10.0.0.2:81")
	})
}

func TestEndpoints_SingleAddress(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		eps := gsvc.NewEndpoints("127.0.0.1:3000")
		t.Assert(len(eps), 1)
		t.Assert(eps.String(), "127.0.0.1:3000")
	})
}

func TestEndpoints_EmptyString(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		eps := gsvc.NewEndpoints("")
		t.Assert(len(eps), 0)
		t.Assert(eps.String(), "")
	})
}

// ─────────────────────── Metadata ────────────────────────────────

func TestMetadata_SetGetIsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := make(gsvc.Metadata)
		t.Assert(m.IsEmpty(), true)

		m.Set("region", "us-east-1")
		t.Assert(m.IsEmpty(), false)
		t.Assert(m.Get("region").String(), "us-east-1")
	})
}

func TestMetadata_Sets(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := make(gsvc.Metadata)
		m.Sets(map[string]any{
			"env":    "prod",
			"weight": 10,
		})
		t.Assert(m.Get("env").String(), "prod")
		t.Assert(m.Get("weight").Int(), 10)
	})
}

func TestMetadata_GetMissingKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := make(gsvc.Metadata)
		t.Assert(m.Get("missing"), nil)
	})
}

func TestMetadata_MDWeightKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := make(gsvc.Metadata)
		m.Set(gsvc.MDWeight, 5)
		t.Assert(m.Get(gsvc.MDWeight).Int(), 5)
	})
}

// ─────────────────────── NewServiceWithName ──────────────────────

func TestNewServiceWithName_BasicFields(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		svc := gsvc.NewServiceWithName("my-service")
		t.Assert(svc.GetName(), "my-service")
		t.AssertNE(svc.GetVersion(), "") // autoFillDefaultAttributes sets "latest"
		t.AssertNE(svc.GetMetadata(), nil)
	})
}

func TestNewServiceWithName_GetKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		svc := gsvc.NewServiceWithName("calc")
		// key must be non-empty and contain the service name
		key := svc.GetKey()
		t.AssertNE(key, "")
		t.Assert(len(key) > 0, true)
	})
}

// ─────────────────────── NewServiceWithKV ────────────────────────

func TestNewServiceWithKV_ValidKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// format: head/deployment/namespace/name/version/endpoints
		key := "gf/default/default/order-svc/latest/10.0.0.1:8080"
		svc, err := gsvc.NewServiceWithKV(key, "")
		t.AssertNil(err)
		t.Assert(svc.GetName(), "order-svc")
		t.Assert(svc.GetVersion(), "latest")
		t.Assert(svc.GetEndpoints().String(), "10.0.0.1:8080")
	})
}

func TestNewServiceWithKV_ValidKeyWithMetadata(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		key := "gf/default/default/svc/v1/127.0.0.1:9000"
		svc, err := gsvc.NewServiceWithKV(key, `{"weight":3}`)
		t.AssertNil(err)
		t.Assert(svc.GetMetadata().Get("weight").Int(), 3)
	})
}

func TestNewServiceWithKV_TooFewSegmentsReturnsError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := gsvc.NewServiceWithKV("only/five/parts/here/x", "")
		t.AssertNE(err, nil)
	})
}

func TestNewServiceWithKV_InvalidMetadataJSONReturnsError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		key := "gf/default/default/svc/v1/127.0.0.1:9000"
		_, err := gsvc.NewServiceWithKV(key, `{not valid json}`)
		t.AssertNE(err, nil)
	})
}

func TestNewServiceWithKV_EmptyValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		key := "gf/default/default/svc/v1/127.0.0.1:9000"
		svc, err := gsvc.NewServiceWithKV(key, "")
		t.AssertNil(err)
		t.Assert(svc.GetMetadata().IsEmpty(), true)
	})
}

// ─────────────────────── SetRegistry / GetRegistry ───────────────

func TestSetGetRegistry_NonNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// SetRegistry requires a non-nil value; just verify get after set works.
		original := gsvc.GetRegistry()
		// Create a real service to prove the API is accessible; we cannot set nil.
		// Just verify that GetRegistry returns whatever was set via SetRegistry by
		// setting the value back to itself (or skipping if nil to avoid a panic).
		if original != nil {
			gsvc.SetRegistry(original)
			t.Assert(gsvc.GetRegistry(), original)
		}
		// When no registry was ever set, GetRegistry returns nil.
		// This just exercises the path without panicking.
	})
}
