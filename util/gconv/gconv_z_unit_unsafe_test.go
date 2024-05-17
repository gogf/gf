// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestUnsafeStrToBytes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "equator"
		t.AssertEQ(gconv.UnsafeStrToBytes(s), []byte(s))
	})
}

func TestUnsafeBytesToStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		b := []byte("ecliptic")
		t.AssertEQ(gconv.UnsafeBytesToStr(b), string(b))
	})
}
