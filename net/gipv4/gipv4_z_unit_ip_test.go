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

func TestGetIpArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ips, err := gipv4.GetIpArray()
		t.AssertNil(err)
		t.AssertGT(len(ips), 0)
		for _, ip := range ips {
			t.Assert(gipv4.Validate(ip), true)
		}
	})
}

func TestMustGetIntranetIp(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustGetIntranetIp() panicked: %v", r)
			}
		}()
		ip := gipv4.MustGetIntranetIp()
		t.Assert(gipv4.IsIntranet(ip), true)
	})
}

func TestGetIntranetIp(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ip, err := gipv4.GetIntranetIp()
		t.AssertNil(err)
		t.AssertNE(ip, "")
		t.Assert(gipv4.IsIntranet(ip), true)
	})
}

func TestGetIntranetIpArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ips, err := gipv4.GetIntranetIpArray()
		t.AssertNil(err)
		t.AssertGT(len(ips), 0)
		for _, ip := range ips {
			t.Assert(gipv4.IsIntranet(ip), true)
		}
	})
}

func TestIsIntranet(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"127.0.0.1", true},
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"172.31.255.255", true},
		{"192.168.0.1", true},
		{"192.168.255.255", true},
		{"8.8.8.8", false},
		{"172.32.0.1", false},
		{"256.256.256.256", false},
	}

	gtest.C(t, func(t *gtest.T) {
		for _, test := range tests {
			result := gipv4.IsIntranet(test.ip)
			t.Assert(result, test.expected)
		}
	})
}
