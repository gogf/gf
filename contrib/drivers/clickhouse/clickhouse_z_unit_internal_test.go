// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_LocalTypeCoverage asserts that every type family the server knows has an explicit
// local type. A family missing from localTypeMap reaches the keyword matching of the core,
// which infers the type from substrings of the name and is wrong for any name that merely
// embeds a keyword, such as `Point` embedding "int".
//
// The families are read from the server rather than from a list compiled into the driver,
// so that a server upgrade adding a type is caught as well.
//
// When this test fails, add the reported families to localTypeMap rather than relying on
// the core to guess them.
func Test_LocalTypeCoverage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		conn, err := gdb.New(gdb.ConfigNode{
			Host: "127.0.0.1", Port: "9000", User: "default", Name: "default", Type: "clickhouse",
		})
		t.AssertNil(err)

		var ctx = context.Background()
		all, err := conn.GetAll(ctx, `SELECT name FROM system.data_type_families`)
		t.AssertNil(err)
		t.AssertGT(len(all), 0)

		var missing []string
		for _, record := range all {
			name := strings.ToLower(record["name"].String())
			if _, ok := localTypeMap[name]; !ok {
				missing = append(missing, name)
			}
		}
		t.Assert(missing, nil)
	})
}

// Test_ConvertValueForField_Allocs tests that converting a value for a column allocates nothing for
// the type of the column, as it is done for every value written.
func Test_ConvertValueForField_Allocs(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx    = context.Background()
			driver = &Driver{}
		)
		t.Assert(testing.AllocsPerRun(100, func() {
			_ = fieldTypeFamily("LowCardinality(Nullable(UInt8))")
		}), 0)
		t.Assert(testing.AllocsPerRun(100, func() {
			_, _ = driver.convertValueByFieldType(ctx, "UInt8", 5)
		}), 0)
		t.Assert(testing.AllocsPerRun(100, func() {
			_, _ = driver.ConvertValueForField(ctx, "Int64", true)
		}), 0)
	})
}

// Test_RawValueSql tests the SQL of the values of a record with a raw value: a time is its exact
// instant, or its date in its own time zone for a Date or Date32 column.
func Test_RawValueSql(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value := time.Date(2024, 1, 2, 23, 30, 0, 0, time.FixedZone("", -5*3600))

		sql, args, err := rawValueSql("Date", value)
		t.AssertNil(err)
		t.Assert(sql, "toDate('2024-01-02')")
		t.Assert(len(args), 0)

		sql, _, err = rawValueSql("Date32", &value)
		t.AssertNil(err)
		t.Assert(sql, "toDate32('2024-01-02')")

		sql, _, err = rawValueSql("DateTime", value)
		t.AssertNil(err)
		t.Assert(sql, "fromUnixTimestamp64Milli("+strconv.FormatInt(value.UnixMilli(), 10)+")")

		sql, _, err = rawValueSql("DateTime64", value.Add(time.Microsecond))
		t.AssertNil(err)
		t.Assert(sql, "fromUnixTimestamp64Micro("+strconv.FormatInt(value.UnixMicro()+1, 10)+")")

		sql, _, err = rawValueSql("", value.Add(time.Nanosecond))
		t.AssertNil(err)
		t.Assert(sql, "fromUnixTimestamp64Nano("+strconv.FormatInt(value.UnixNano()+1, 10)+")")
	})
	gtest.C(t, func(t *gtest.T) {
		sql, args, err := rawValueSql("String", []byte{0xff, 0x00})
		t.AssertNil(err)
		t.Assert(sql, "?")
		t.Assert(args, []any{"\xff\x00"})

		sql, args, err = rawValueSql("UInt8", 5)
		t.AssertNil(err)
		t.Assert(sql, "?")
		t.Assert(args, []any{5})

		_, _, err = rawValueSql("Array(String)", []string{"a"})
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)
	})
}

// Test_CatchSQLManager tests that the driver finds the gdb.CatchSQLManager of CatchSQL and ToSQL in the
// context, under the key of the core.
func Test_CatchSQLManager(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertNil(catchSQLManager(context.Background()))

		_, err := gdb.CatchSQL(context.Background(), func(ctx context.Context) error {
			manager := catchSQLManager(ctx)
			t.AssertNE(manager, nil)
			t.Assert(manager.DoCommit, true)
			return nil
		})
		t.AssertNil(err)

		_, err = gdb.ToSQL(context.Background(), func(ctx context.Context) error {
			manager := catchSQLManager(ctx)
			t.AssertNE(manager, nil)
			t.Assert(manager.DoCommit, false)
			return nil
		})
		t.AssertNil(err)
	})
}
