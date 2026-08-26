// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file verifies native and fallback generated primary-key value results.

package gdb

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/test/gtest"
)

// TestInsertResultSignatureCompatibility verifies that existing insert APIs keep returning sql.Result.
func TestInsertResultSignatureCompatibility(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			dbInsert    func(DB, context.Context, string, any, ...int) (sql.Result, error) = DB.Insert
			txInsert    func(TX, string, any, ...int) (sql.Result, error)                  = TX.Insert
			modelInsert func(*Model, ...any) (sql.Result, error)                           = (*Model).Insert
		)
		t.AssertNE(dbInsert, nil)
		t.AssertNE(txInsert, nil)
		t.AssertNE(modelInsert, nil)
	})
}

// testSQLResult provides a standard integer-only SQL result for fallback tests.
type testSQLResult struct {
	lastInsertID      int64
	lastInsertIDError error
}

// LastInsertId returns the configured integer ID and error.
func (r testSQLResult) LastInsertId() (int64, error) {
	return r.lastInsertID, r.lastInsertIDError
}

// RowsAffected returns one affected row for the test result.
func (r testSQLResult) RowsAffected() (int64, error) {
	return 1, nil
}

// testNativeResultWithPrimaryKeyValue provides a generated primary-key value independently of LastInsertId.
type testNativeResultWithPrimaryKeyValue struct {
	testSQLResult
	lastInsertPrimaryKeyValue Value
}

// LastInsertPrimaryKeyValue returns the configured native generated primary-key value.
func (r testNativeResultWithPrimaryKeyValue) LastInsertPrimaryKeyValue() (Value, error) {
	return r.lastInsertPrimaryKeyValue, nil
}

// TestGetLastInsertPrimaryKeyValue verifies native values, integer fallbacks, and error propagation.
func TestGetLastInsertPrimaryKeyValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		const expectedUUID = "c3e0d2e0-03ec-4db9-a7d5-e8c3ec3ff7f2"
		var result sql.Result = testNativeResultWithPrimaryKeyValue{
			lastInsertPrimaryKeyValue: gvar.New(expectedUUID),
		}
		valueResult, ok := result.(ResultWithPrimaryKeyValue)
		t.Assert(ok, true)
		value, err := valueResult.LastInsertPrimaryKeyValue()
		t.AssertNil(err)
		t.Assert(value.String(), expectedUUID)
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := getLastInsertPrimaryKeyValue(testSQLResult{lastInsertID: 42})
		t.AssertNil(err)
		t.Assert(value.Int64(), int64(42))
	})

	gtest.C(t, func(t *gtest.T) {
		expectedError := errors.New("last insert ID unavailable")
		value, err := getLastInsertPrimaryKeyValue(testSQLResult{lastInsertIDError: expectedError})
		t.Assert(value, nil)
		t.Assert(errors.Is(err, expectedError), true)
	})

	gtest.C(t, func(t *gtest.T) {
		const expectedUUID = "db50b8d6-6e5a-4fc9-b61a-6d21297844e0"

		result := testNativeResultWithPrimaryKeyValue{
			testSQLResult: testSQLResult{
				lastInsertIDError: errors.New("integer conversion must not be used"),
			},
			lastInsertPrimaryKeyValue: gvar.New(expectedUUID),
		}
		value, err := getLastInsertPrimaryKeyValue(result)
		t.AssertNil(err)
		t.Assert(value.String(), expectedUUID)
	})

	gtest.C(t, func(t *gtest.T) {
		const expectedUUID = "94e54331-16b7-4a07-8f45-71d6c6c2cb80"

		result := &SqlResult{
			Result: testNativeResultWithPrimaryKeyValue{
				lastInsertPrimaryKeyValue: gvar.New(expectedUUID),
			},
		}
		value, err := result.LastInsertPrimaryKeyValue()
		t.AssertNil(err)
		t.Assert(value.String(), expectedUUID)
	})

	gtest.C(t, func(t *gtest.T) {
		result := &SqlResult{}
		value, err := result.LastInsertPrimaryKeyValue()
		t.AssertNil(err)
		t.Assert(value.IsNil(), true)
	})
}
