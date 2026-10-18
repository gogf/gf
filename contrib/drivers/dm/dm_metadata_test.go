// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_recordValue(t *testing.T) {
	// Without the columnNameCase option the result set keeps the uppercase column
	// names written in the metadata queries.
	gtest.C(t, func(t *gtest.T) {
		record := gdb.Record{
			"COLUMN_NAME": gvar.New("ACCOUNT_NAME"),
			"DATA_LENGTH": gvar.New(128),
		}
		t.Assert(recordValue(record, "COLUMN_NAME").String(), "ACCOUNT_NAME")
		t.Assert(recordValue(record, "DATA_LENGTH").Int(), 128)
	})
	// With columnNameCase=lower the DM driver lowercases every column name of the
	// result set, so the metadata must still be readable.
	gtest.C(t, func(t *gtest.T) {
		record := gdb.Record{
			"column_name": gvar.New("ACCOUNT_NAME"),
			"data_length": gvar.New(128),
		}
		t.Assert(recordValue(record, "COLUMN_NAME").String(), "ACCOUNT_NAME")
		t.Assert(recordValue(record, "DATA_LENGTH").Int(), 128)
	})
	// A column that is absent from the result set returns a nil value, which is
	// still safe to read from.
	gtest.C(t, func(t *gtest.T) {
		record := gdb.Record{
			"column_name": gvar.New("ACCOUNT_NAME"),
		}
		t.Assert(recordValue(record, "COMMENTS") == nil, true)
		t.Assert(recordValue(record, "COMMENTS").String(), "")
	})
}

func Test_lookupRecordValue(t *testing.T) {
	// A column holding a SQL NULL is present with a nil value, which must be told
	// apart from a column that is absent from the result set.
	gtest.C(t, func(t *gtest.T) {
		record := gdb.Record{
			"iot_name": nil,
		}
		value, ok := lookupRecordValue(record, "IOT_NAME")
		t.Assert(ok, true)
		t.Assert(value == nil, true)

		value, ok = lookupRecordValue(record, "TABLE_NAME")
		t.Assert(ok, false)
		t.Assert(value == nil, true)
	})
}
