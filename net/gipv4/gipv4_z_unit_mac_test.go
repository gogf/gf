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

func TestGetMac(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		mac, err := gipv4.GetMac()
		t.AssertNil(err)
		t.AssertNE(mac, "")
		// MAC addresses are typically 17 characters in length
		t.Assert(len(mac), 17)
	})
}

func TestGetMacArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		macs, err := gipv4.GetMacArray()
		t.AssertNil(err)
		t.AssertGT(len(macs), 0)
		for _, mac := range macs {
			// MAC addresses are typically 17 characters in length
			t.Assert(len(mac), 17)
		}
	})
}
