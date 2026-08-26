// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file verifies PostgreSQL generated primary-key value result behavior.

package pgsql

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

// TestResultLastInsertPrimaryKeyValue verifies UUID values and empty result handling.
func TestResultLastInsertPrimaryKeyValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		const expectedUUID = "4262b9b5-bc41-418a-a668-9201d72b069f"

		var result sql.Result = Result{lastInsertPrimaryKeyValue: gvar.New(expectedUUID)}
		valueResult, ok := result.(gdb.ResultWithPrimaryKeyValue)
		t.Assert(ok, true)
		value, err := valueResult.LastInsertPrimaryKeyValue()
		t.AssertNil(err)
		t.Assert(value.String(), expectedUUID)
	})

	gtest.C(t, func(t *gtest.T) {
		result := Result{}
		value, err := result.LastInsertPrimaryKeyValue()
		t.AssertNil(err)
		t.Assert(value.IsNil(), true)
	})

	gtest.C(t, func(t *gtest.T) {
		expectedError := errors.New("UUID is not an integer ID")
		result := Result{
			lastInsertIdError:         expectedError,
			lastInsertPrimaryKeyValue: gvar.New("f93d9363-5cb8-4751-9b62-ad145f61cc30"),
		}
		_, err := result.LastInsertId()
		t.Assert(errors.Is(err, expectedError), true)

		value, err := result.LastInsertPrimaryKeyValue()
		t.AssertNil(err)
		t.Assert(value.String(), "f93d9363-5cb8-4751-9b62-ad145f61cc30")
	})
}
