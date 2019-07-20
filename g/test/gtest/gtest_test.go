// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtest_test

import (
	"testing"

	"github.com/gogf/gf/g/test/gtest"
)

func TestCase(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(1, 1)
		gtest.AssertNE(1, 0)
		gtest.AssertEQ(float32(123.456), float32(123.456))
		gtest.AssertEQ(float64(123.456), float64(123.456))
	})
}
