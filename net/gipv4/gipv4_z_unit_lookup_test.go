// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gipv4_test

import (
	"testing"

	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestGetHostByName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ip, err := gipv4.GetHostByName("localhost")
		t.AssertNil(err)
		t.Assert(ip, "127.0.0.1")
	})
}

func TestGetHostsByName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ips, err := gipv4.GetHostsByName("localhost")
		t.AssertNil(err)
		t.AssertIN("127.0.0.1", ips)
	})
}
