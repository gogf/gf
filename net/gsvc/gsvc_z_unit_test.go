// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsvc_test

import (
	"testing"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_NewServiceWithKV_ValidAndMalformed(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		svc, err := gsvc.NewServiceWithKV(
			"/service/default/default/demo/latest/127.0.0.1:8000",
			`{"protocol":"grpc"}`,
		)
		t.AssertNil(err)
		t.Assert(svc.GetName(), "demo")
		t.Assert(svc.GetEndpoints().String(), "127.0.0.1:8000")
	})

	// Extra segment caused by "/" in service name must return an error, not panic.
	gtest.C(t, func(t *gtest.T) {
		svc, err := gsvc.NewServiceWithKV(
			"/service/default/default/gogf-grpcx/demo2/latest/127.0.0.1:8000",
			`{"protocol":"grpc"}`,
		)
		t.AssertNE(err, nil)
		t.AssertNil(svc)
	})

	// Malformed endpoint segment in a 6-segment key must return an error, not panic.
	gtest.C(t, func(t *gtest.T) {
		svc, err := gsvc.NewServiceWithKV(
			"/service/default/default/demo/latest/not_an_endpoint",
			"",
		)
		t.AssertNE(err, nil)
		t.AssertNil(svc)
	})

	// Too few segments must return an error.
	gtest.C(t, func(t *gtest.T) {
		svc, err := gsvc.NewServiceWithKV("/service/default/demo", "")
		t.AssertNE(err, nil)
		t.AssertNil(svc)
	})
}
