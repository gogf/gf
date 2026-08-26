// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gaussdb_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

// Test_IssueBytea_RoundTrip verifies binary data survives a bytea round trip.
// The driver used to route []byte through the array-literal rewrite ('[' -> '{',
// ']' -> '}'), silently corrupting payloads containing 0x5B/0x5D.
func Test_IssueBytea_RoundTrip(t *testing.T) {
	table := "issue_bytea_round_trip"
	dropTable(table)
	if _, err := db.Exec(ctx, `CREATE TABLE `+table+` (
		id   bigserial PRIMARY KEY,
		data bytea
	)`); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Payload deliberately contains '[' (0x5B) and ']' (0x5D).
		input := []byte{0xDE, 0xAD, 0x5B, 0x5D, 0xBE, 0xEF}
		_, err := db.Model(table).Data(g.Map{"id": 1, "data": input}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["data"].Bytes(), input)
	})

	gtest.C(t, func(t *gtest.T) {
		// Every byte value must survive unchanged, not just the two above.
		input := make([]byte, 256)
		for i := range input {
			input[i] = byte(i)
		}
		_, err := db.Model(table).Data(g.Map{"id": 2, "data": input}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["data"].Bytes(), input)
	})

	gtest.C(t, func(t *gtest.T) {
		// An empty payload round trips as an empty, non-nil slice.
		_, err := db.Model(table).Data(g.Map{"id": 3, "data": []byte{}}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(len(one["data"].Bytes()), 0)
		t.Assert(one["data"].IsNil(), false)
	})

	gtest.C(t, func(t *gtest.T) {
		// A NULL bytea column reads back as nil, not as an empty payload.
		_, err := db.Model(table).Data(g.Map{"id": 4}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 4).One()
		t.AssertNil(err)
		t.Assert(one["data"].IsNil(), true)
	})
}

// Test_IssueSave_NullTypedColumn verifies Save() round trips a row whose
// non-text column is NULL.
//
// Save() is built as a MERGE, and the NULL placeholder in its USING branch used
// to carry no type. GaussDB inferred text for it, so assigning the value back to
// the real column failed with "column ... is of type X but expression is of type
// text" — breaking any read-modify-write cycle over a nullable non-text column.
func Test_IssueSave_NullTypedColumn(t *testing.T) {
	// Column types that are not implicitly coercible from text, plus types whose
	// TableField.Type carries a modifier — "int4(32)" and friends reject a cast
	// with one, so the modifier has to be stripped before it is emitted.
	for _, columnType := range []string{
		"numeric[]", "text[]", "jsonb", "boolean", "bytea", "uuid",
		"int2", "int4", "int8", "float4", "float8", "numeric(10,2)", "varchar(45)",
	} {
		gtest.C(t, func(t *gtest.T) {
			table := "issue_save_null_" + gtime.TimestampMicroStr()
			if _, err := db.Exec(ctx, fmt.Sprintf(
				`CREATE TABLE %s (id bigserial PRIMARY KEY, nickname varchar(45), val %s)`,
				table, columnType,
			)); err != nil {
				gtest.Fatal(err)
			}
			defer dropTable(table)

			_, err := db.Model(table).Data(g.Map{"id": 1, "nickname": "n1"}).Insert()
			t.AssertNil(err)

			// Read the row back; val comes back as NULL.
			one, err := db.Model(table).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["val"].IsNil(), true)

			// Modify an unrelated column and write the whole row back.
			data := one.Map()
			data["nickname"] = "n1-updated"
			_, err = db.Model(table).Data(data).Save()
			if err != nil {
				t.Errorf(`Save() failed for column type %s: %v`, columnType, err)
				return
			}

			one, err = db.Model(table).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["nickname"].String(), "n1-updated")
			t.Assert(one["val"].IsNil(), true)
		})
	}
}

// Test_IssueSave_NullTypedColumn_Replace covers the same path through Replace(),
// which the driver also routes through doSave.
func Test_IssueSave_NullTypedColumn_Replace(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := "issue_replace_null_" + gtime.TimestampMicroStr()
		if _, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id bigserial PRIMARY KEY, nickname varchar(45), tags text[])`, table,
		)); err != nil {
			gtest.Fatal(err)
		}
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{"id": 1, "nickname": "n1"}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"id": 1, "nickname": "n1-replaced", "tags": nil,
		}).Replace()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "n1-replaced")
		t.Assert(one["tags"].IsNil(), true)
	})
}

// Test_IssueSave_NullTypedColumn_QuotedTypeName covers a column whose type was
// created with a quoted, mixed-case name.
//
// pg_type.typname keeps the original case, so emitting the cast target without
// quotes let the server fold it to lower case and then fail to resolve it with
// `type "myenum" does not exist`.
func Test_IssueSave_NullTypedColumn_QuotedTypeName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			suffix   = gtime.TimestampMicroStr()
			typeName = `"MyEnum` + suffix + `"`
			table    = "issue_save_quoted_type_" + suffix
		)
		if _, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TYPE %s AS ENUM ('a', 'b')`, typeName,
		)); err != nil {
			gtest.Fatal(err)
		}
		defer func() {
			_, err := db.Exec(ctx, fmt.Sprintf(`DROP TYPE IF EXISTS %s CASCADE`, typeName))
			t.AssertNil(err)
		}()

		if _, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id bigserial PRIMARY KEY, nickname varchar(45), val %s)`,
			table, typeName,
		)); err != nil {
			gtest.Fatal(err)
		}
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{"id": 1, "nickname": "n1"}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["val"].IsNil(), true)

		data := one.Map()
		data["nickname"] = "n1-updated"
		_, err = db.Model(table).Data(data).Save()
		t.AssertNil(err)

		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "n1-updated")
		t.Assert(one["val"].IsNil(), true)
	})
}

// Test_IssuePointInterval_NotConvertedToInt verifies that column types whose
// names merely contain an integer keyword are not converted to integers.
//
// gdb's fallback type detection matches "int" as a substring, so "point",
// "interval", "tinterval" and the range types were all classified as integers
// and their values became 0.
func Test_IssuePointInterval_NotConvertedToInt(t *testing.T) {
	table := "issue_point_interval_" + gtime.TimestampMicroStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (
		id     int PRIMARY KEY,
		pt     point,
		span   interval,
		tspan  tinterval,
		r4     int4range,
		r8     int8range,
		pts    point[],
		spans  interval[]
	)`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s VALUES (
			1, '(1.5,2.5)', '2 days',
			'["2024-01-01 00:00:00" "2024-01-02 00:00:00"]',
			'[1,5)', '[10,50)',
			'{"(1,2)","(3,4)"}', '{"1 day","2 days"}'
		)`, table))
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)

		// Scalars keep their real content instead of collapsing to 0.
		t.Assert(one["pt"].String(), "(1.5,2.5)")
		t.Assert(one["span"].String(), "2 days")
		t.Assert(one["r4"].String(), "[1,5)")
		t.Assert(one["r8"].String(), "[10,50)")
		t.AssertNE(one["tspan"].String(), "0")
		t.Assert(gstr.Contains(one["tspan"].String(), "2024-01-01"), true)

		// Array forms are read as their element text representation.
		t.Assert(one["pts"].Strings(), g.SliceStr{"(1,2)", "(3,4)"})
		t.Assert(one["spans"].Strings(), g.SliceStr{"1 day", "2 days"})
	})

	gtest.C(t, func(t *gtest.T) {
		// Control: real integer columns must keep working.
		intTable := "issue_int_control_" + gtime.TimestampMicroStr()
		if _, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id int PRIMARY KEY, a int2, b int4, c int8, d int4[])`, intTable,
		)); err != nil {
			gtest.Fatal(err)
		}
		defer dropTable(intTable)

		_, err := db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s VALUES (1, 1, 2, 3, '{4,5}')`, intTable))
		t.AssertNil(err)

		one, err := db.Model(intTable).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["a"].Int(), 1)
		t.Assert(one["b"].Int(), 2)
		t.Assert(one["c"].Int64(), int64(3))
		t.Assert(one["d"].Ints(), []int{4, 5})
	})
}
