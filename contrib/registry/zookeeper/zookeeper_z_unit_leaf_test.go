// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package zookeeper

import (
	"testing"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_serviceInstanceLeaf: unique ZK leaf from endpoints (#4717).
func Test_serviceInstanceLeaf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		svc := &gsvc.LocalService{
			Name:      "hello.svc",
			Version:   "v1.0.0",
			Endpoints: gsvc.NewEndpoints("127.0.0.1:9000"),
		}
		t.Assert(serviceInstanceLeaf(svc), "127.0.0.1-9000")

		svc2 := &gsvc.LocalService{
			Name:      "hello.svc",
			Version:   "v1.0.0",
			Endpoints: gsvc.NewEndpoints("10.0.0.1:8080"),
		}
		t.Assert(serviceInstanceLeaf(svc2), "10.0.0.1-8080")
		// Different instances of the same service get different leaves.
		t.AssertNE(serviceInstanceLeaf(svc), serviceInstanceLeaf(svc2))
	})
}
