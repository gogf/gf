// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Issue4699 tests negative values for Limit/Page/Offset should be treated as zero.
// See https://github.com/gogf/gf/issues/4699
func Test_Issue4699(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a base model for testing
		m := &Model{}

		// Test Limit with single negative parameter
		m1 := m.Limit(-1)
		t.AssertEQ(m1.limit, 0)

		// Test Limit with two parameters (start, limit) where both are negative
		m2 := m.Limit(-10, -5)
		t.AssertEQ(m2.start, 0)
		t.AssertEQ(m2.limit, 0)

		// Test Limit with mixed parameters (negative start, positive limit)
		m3 := m.Limit(-10, 5)
		t.AssertEQ(m3.start, 0)
		t.AssertEQ(m3.limit, 5)

		// Test Page with negative limit
		m4 := m.Page(1, -10)
		t.AssertEQ(m4.start, 0)
		t.AssertEQ(m4.limit, 0)

		// Test Page with negative limit on page 2
		m5 := m.Page(2, -10)
		t.AssertEQ(m5.start, 0) // (2-1) * 0 = 0
		t.AssertEQ(m5.limit, 0)

		// Test Offset with negative value
		m6 := m.Offset(-5)
		t.AssertEQ(m6.offset, 0)

		// Test Offset with positive value (sanity check)
		m7 := m.Offset(10)
		t.AssertEQ(m7.offset, 10)
	})
}
