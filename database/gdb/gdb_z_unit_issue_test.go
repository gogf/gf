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

// Test_ScanValidateSingleFieldSpecified tests that validateSingleFieldSpecified
// and isSingleFieldSpecified correctly accept/reject various field configurations,
// especially the gdb.Raw case which bypasses len(m.fields) checks.
func Test_ScanValidateSingleFieldSpecified(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// No fields at all → reject.
		m0 := &Model{}
		t.AssertNE(m0.validateSingleFieldSpecified(), nil)
		t.Assert(m0.isSingleFieldSpecified(), false)

		// Exactly one normal field → accept.
		m1 := &Model{fields: []any{"name"}}
		t.AssertNil(m1.validateSingleFieldSpecified())
		t.Assert(m1.isSingleFieldSpecified(), true)

		// Two normal fields → reject.
		m2 := &Model{fields: []any{"name", "age"}}
		t.AssertNE(m2.validateSingleFieldSpecified(), nil)
		t.Assert(m2.isSingleFieldSpecified(), false)

		// gdb.Raw("name,age") stored as one field entry → reject,
		// because it may expand to multiple columns at the SQL level.
		m3 := &Model{fields: []any{Raw("name,age")}}
		t.AssertNE(m3.validateSingleFieldSpecified(), nil)
		t.Assert(m3.isSingleFieldSpecified(), false)

		// FieldsEx only → accept (actual column count checked post-execution).
		m4 := &Model{fieldsEx: []any{"id"}}
		t.AssertNil(m4.validateSingleFieldSpecified())

		// gdb.Raw("name") as single field → reject (still Raw, cannot guarantee).
		m5 := &Model{fields: []any{Raw("name")}}
		t.AssertNE(m5.validateSingleFieldSpecified(), nil)
		t.Assert(m5.isSingleFieldSpecified(), false)
	})
}
