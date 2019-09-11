// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"testing"
)

func Test_Unsafe(t *testing.T) {
	gtest.Case(t, func() {
		s := "I love 小泽玛利亚"
		gtest.AssertEQ(gconv.UnsafeStrToBytes(s), []byte(s))
	})

	gtest.Case(t, func() {
		b := []byte("I love 小泽玛利亚")
		gtest.AssertEQ(gconv.UnsafeBytesToStr(b), string(b))
	})
}
