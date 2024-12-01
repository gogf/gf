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

const (
	ipv4             string = "192.168.1.1"
	longBigEndian    uint32 = 3232235777
	longLittleEndian uint32 = 16885952
)

func TestIpToLongBigEndian(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var u = gipv4.IpToLongBigEndian(ipv4)
		t.Assert(u, longBigEndian)

		var u2 = gipv4.Ip2long(ipv4)
		t.Assert(u2, longBigEndian)
	})
}

func TestLongToIpBigEndian(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s = gipv4.LongToIpBigEndian(longBigEndian)
		t.Assert(s, ipv4)

		var s2 = gipv4.Long2ip(longBigEndian)
		t.Assert(s2, ipv4)
	})
}

func TestIpToLongLittleEndian(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var u = gipv4.IpToLongLittleEndian(ipv4)
		t.Assert(u, longLittleEndian)
	})
}

func TestLongToIpLittleEndian(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s = gipv4.LongToIpLittleEndian(longLittleEndian)
		t.Assert(s, ipv4)
	})
}
