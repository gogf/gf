// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gaussdb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
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

// Test_Issue4842 tests that types whose names merely contain "int" are not detected as
// integers and read back as 0. The driver lists most of them explicitly, but int2vector
// is only covered by the fallback detection of the core. The daterange and reltime scalars
// were listed as a datetime, which lost everything the parser could not read as a date.
// The int1 and int16 arrays were listed as a scalar integer, which loses the array the
// same way, and the int16 scalar is a 16-byte integer that does not fit an int at all.
// See https://github.com/gogf/gf/issues/4842
func Test_Issue4842(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// The int2vector value must match its own text rendering rather than becoming 0.
		one, err := db.GetOne(ctx, `SELECT indkey AS v, indkey::text AS expect FROM pg_index LIMIT 1`)
		t.AssertNil(err)
		t.Assert(one["v"].String(), one["expect"].String())
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, `SELECT '(1.5,2.5)'::point AS a, '2 days'::interval AS b,
			'[1,10)'::int4range AS c, '[1,10)'::int8range AS d`)
		t.AssertNil(err)
		t.Assert(one["a"].String(), `(1.5,2.5)`)
		t.Assert(one["b"].String(), `2 days`)
		t.Assert(one["c"].String(), `[1,10)`)
		t.Assert(one["d"].String(), `[1,10)`)
	})
	gtest.C(t, func(t *gtest.T) {
		// The daterange and reltime scalars keep their text instead of being parsed as a
		// datetime, which kept the lower bound of a range only and dropped a relative time.
		one, err := db.GetOne(ctx, `SELECT '[2026-01-01,2026-02-01)'::daterange AS a,
			'1 day'::reltime AS b, ARRAY['1 day'::reltime] AS c`)
		t.AssertNil(err)
		t.Assert(one["a"].String(), `[2026-01-01,2026-02-01)`)
		t.Assert(one["b"].String(), `1 day`)
		t.Assert(one["c"].Strings(), g.SliceStr{"1 day"})
	})
	gtest.C(t, func(t *gtest.T) {
		// Genuine integer types must keep being detected as integers.
		one, err := db.GetOne(ctx, `SELECT 41::int2 AS a, 42::int4 AS b, 43::int8 AS c`)
		t.AssertNil(err)
		t.Assert(one["a"].Int(), 41)
		t.Assert(one["b"].Int(), 42)
		t.Assert(one["c"].Int64(), int64(43))
	})
	gtest.C(t, func(t *gtest.T) {
		// int1[] and int16[] are arrays, like the _int2 listed next to them, and were read
		// as a single integer, which turned the whole array literal into 0. int16 is the
		// 16-byte integer of GaussDB rather than a 16-bit one, so reading it as an int
		// truncated any value beyond int64 instead of keeping its decimal text.
		one, err := db.GetOne(ctx, `SELECT '{1,2,255}'::int1[] AS a, 255::int1 AS b,
			'{1,170141183460469231731687303715884105727}'::int16[] AS c,
			170141183460469231731687303715884105727::int16 AS d`)
		t.AssertNil(err)
		t.Assert(one["a"].Ints(), g.SliceInt{1, 2, 255})
		t.Assert(one["b"].Int(), 255)
		t.Assert(one["c"].Strings(), g.SliceStr{"1", "170141183460469231731687303715884105727"})
		t.Assert(one["d"].String(), `170141183460469231731687303715884105727`)
	})
	gtest.C(t, func(t *gtest.T) {
		// int1[] is a plain tinyint array, so a declared column reads back as one too.
		// int16 has no column form: GaussDB rejects it with "not supported to create".
		table := "issue4842_int1_array_" + gtime.TimestampMicroStr()
		if _, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id int PRIMARY KEY, v int1[], w tinyint)`, table,
		)); err != nil {
			gtest.Fatal(err)
		}
		defer dropTable(table)

		_, err := db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s VALUES (1, '{1,2,255}', 7)`, table))
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["v"].Ints(), g.SliceInt{1, 2, 255})
		t.Assert(one["w"].Int(), 7)
	})
}

// Test_Issue4843 tests that Save and Replace work on columns that do not accept a text
// value, such as boolean, json, jsonb, uuid, bytea, timestamptz and the array types.
// Both are built as a MERGE whose USING branch carries every value as an untyped
// parameter, which the server takes as text, so assigning it to such a column failed.
// See https://github.com/gogf/gf/issues/4843
func Test_Issue4843(t *testing.T) {
	type typedCase struct {
		colType string
		value   any
		updated any
		expect  func(v *gvar.Var) any
		want    any
		wantUpd any
	}
	var (
		asString = func(v *gvar.Var) any { return v.String() }
		asBool   = func(v *gvar.Var) any { return v.Bool() }
		asInts   = func(v *gvar.Var) any { return v.Ints() }
		asStrs   = func(v *gvar.Var) any { return v.Strings() }
		asUnix   = func(v *gvar.Var) any { return v.GTime().Timestamp() }
		cases    = []typedCase{
			{"boolean", true, false, asBool, true, false},
			{"int4", 42, 43, asString, "42", "43"},
			{"varchar(20)", "abc", "abd", asString, "abc", "abd"},
			{"numeric(10,2)", "12.34", "56.78", asString, "12.34", "56.78"},
			{"timestamptz", "2024-01-02 03:04:05+00", "2024-02-03 04:05:06+00", asUnix,
				time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC).Unix(),
				time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC).Unix()},
			{"json", `{"a":1}`, `{"a":2}`, asString, `{"a":1}`, `{"a":2}`},
			{"jsonb", `{"a":1}`, `{"a":2}`, asString, `{"a": 1}`, `{"a": 2}`},
			{"uuid", "550e8400-e29b-41d4-a716-446655440000", "550e8400-e29b-41d4-a716-446655440001", asString,
				"550e8400-e29b-41d4-a716-446655440000", "550e8400-e29b-41d4-a716-446655440001"},
			{"bytea", []byte("hello"), []byte("world"), asString, "hello", "world"},
			{"int4[]", []int{1, 2}, []int{3, 4}, asInts, []int{1, 2}, []int{3, 4}},
			{"text[]", []string{"a", "b"}, []string{"c", "d"}, asStrs, []string{"a", "b"}, []string{"c", "d"}},
			{"numeric[]", "{1.5,2.5}", "{3.5,4.5}", asStrs, []string{"1.5", "2.5"}, []string{"3.5", "4.5"}},
		}
	)
	for _, c := range cases {
		c := c
		gtest.C(t, func(t *gtest.T) {
			table := fmt.Sprintf("issue4843_%d", gtime.TimestampNano())
			_, err := db.Exec(ctx, fmt.Sprintf(
				`CREATE TABLE %s (id bigserial PRIMARY KEY, nickname varchar(45), val %s)`, table, c.colType,
			))
			t.AssertNil(err)
			defer dropTable(table)

			// Save inserts a new row.
			_, err = db.Model(table).Data(g.Map{"id": 1, "nickname": "n1", "val": c.value}).Save()
			t.Assert(err, nil)
			one, err := db.Model(table).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["nickname"].String(), "n1")
			t.Assert(c.expect(one["val"]), c.want)

			// Save updates the existing row.
			_, err = db.Model(table).Data(g.Map{"id": 1, "nickname": "n2", "val": c.updated}).Save()
			t.Assert(err, nil)
			one, err = db.Model(table).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["nickname"].String(), "n2")
			t.Assert(c.expect(one["val"]), c.wantUpd)

			// Replace is built the same way.
			_, err = db.Model(table).Data(g.Map{"id": 1, "nickname": "n3", "val": c.value}).Replace()
			t.Assert(err, nil)
			one, err = db.Model(table).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["nickname"].String(), "n3")
			t.Assert(c.expect(one["val"]), c.want)
		})
	}

	// Reading a whole row and writing it back, with a mix of typed and NULL columns.
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("issue4843_rmw_%d", gtime.TimestampNano())
		_, err := db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (
			id bigserial PRIMARY KEY, passport varchar(45) NOT NULL, age int4, score numeric(10,2),
			active boolean, created timestamp, birthday date, meta jsonb, ext uuid, tags text[])`, table))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{
			"id": 1, "passport": "p", "age": 18, "score": "9.50", "active": true,
			"created": "2024-01-02 03:04:05", "birthday": "2000-01-02", "meta": `{"a":1}`,
			"tags": []string{"x", "y"},
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		data := one.Map()
		data["passport"] = "p2"
		_, err = db.Model(table).Data(data).Save()
		t.AssertNil(err)

		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"].String(), "p2")
		t.Assert(one["age"].Int(), 18)
		t.Assert(one["score"].String(), "9.50")
		t.Assert(one["active"].Bool(), true)
		t.Assert(one["meta"].String(), `{"a": 1}`)
		t.Assert(one["ext"].IsNil(), true)
		t.Assert(one["tags"].Strings(), g.SliceStr{"x", "y"})
	})

	// A batch of records, each saved through its own MERGE.
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("issue4843_batch_%d", gtime.TimestampNano())
		_, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id bigserial PRIMARY KEY, nickname varchar(45), val boolean)`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{"id": 1, "nickname": "old", "val": false}).Insert()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.List{
			{"id": 1, "nickname": "a", "val": true},
			{"id": 2, "nickname": "b", "val": true},
		}).Save()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["nickname"].String(), "a")
		t.Assert(all[0]["val"].Bool(), true)
		t.Assert(all[1]["nickname"].String(), "b")
		t.Assert(all[1]["val"].Bool(), true)
	})

	// Custom update columns and transactions go through the same MERGE.
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("issue4843_dup_%d", gtime.TimestampNano())
		_, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id bigserial PRIMARY KEY, nickname varchar(45), val boolean)`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{"id": 1, "nickname": "keep", "val": false}).Insert()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"id": 1, "nickname": "ignored", "val": true}).
			OnDuplicate("val").Save()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "keep")
		t.Assert(one["val"].Bool(), true)

		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).Data(g.Map{"id": 1, "nickname": "tx", "val": false}).Save()
			return err
		})
		t.AssertNil(err)
		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "tx")
		t.Assert(one["val"].Bool(), false)
	})

	// The cast must not carry the type modifier: a value that does not fit the column is
	// rejected as Insert rejects it, and the scale of a numeric column is kept.
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("issue4843_mod_%d", gtime.TimestampNano())
		_, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id bigserial PRIMARY KEY, short varchar(5), amount numeric(10,2), fixed char(5))`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{"id": 1, "short": "abcdefgh"}).Insert()
		t.AssertNE(err, nil)
		_, err = db.Model(table).Data(g.Map{"id": 1, "short": "abcdefgh"}).Save()
		t.AssertNE(err, nil)
		_, err = db.Model(table).Data(g.Map{"id": 1, "short": "abcdefgh"}).Replace()
		t.AssertNE(err, nil)
		n, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, 0)

		_, err = db.Model(table).Data(g.Map{"id": 1, "short": "abc", "amount": "12.345", "fixed": "ab"}).Save()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["short"].String(), "abc")
		t.Assert(one["amount"].String(), "12.35")
		t.Assert(one["fixed"].String(), "ab   ")
	})

	// A column of a user defined type created with a quoted, mixed-case name.
	gtest.C(t, func(t *gtest.T) {
		var (
			suffix   = gtime.TimestampMicroStr()
			typeName = `"MyEnum` + suffix + `"`
			table    = "issue4843_enum_" + suffix
		)
		_, err := db.Exec(ctx, fmt.Sprintf(`CREATE TYPE %s AS ENUM ('a', 'b')`, typeName))
		t.AssertNil(err)
		defer func() {
			_, err := db.Exec(ctx, fmt.Sprintf(`DROP TYPE IF EXISTS %s CASCADE`, typeName))
			t.AssertNil(err)
		}()
		_, err = db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id bigserial PRIMARY KEY, nickname varchar(45), val %s)`, table, typeName,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{"id": 1, "nickname": "n1", "val": "a"}).Save()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"id": 1, "nickname": "n2", "val": "b"}).Save()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "n2")
		t.Assert(one["val"].String(), "b")
	})

	// A column whose type lives in a schema that is not on the search path.
	gtest.C(t, func(t *gtest.T) {
		var (
			suffix = gtime.TimestampMicroStr()
			schema = "issue4843_s_" + suffix
			table  = "issue4843_xs_" + suffix
		)
		_, err := db.Exec(ctx, fmt.Sprintf(`CREATE SCHEMA %s`, schema))
		t.AssertNil(err)
		defer func() {
			_, err := db.Exec(ctx, fmt.Sprintf(`DROP SCHEMA IF EXISTS %s CASCADE`, schema))
			t.AssertNil(err)
		}()
		_, err = db.Exec(ctx, fmt.Sprintf(`CREATE DOMAIN %s.dom_int AS int4`, schema))
		t.AssertNil(err)
		_, err = db.Exec(ctx, fmt.Sprintf(`CREATE TYPE %s.en_t AS ENUM ('a', 'b')`, schema))
		t.AssertNil(err)
		_, err = db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id bigserial PRIMARY KEY, di %s.dom_int, en %s.en_t)`, table, schema, schema,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{"id": 1, "di": 7, "en": "a"}).Save()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"id": 1, "di": 8, "en": "b"}).Save()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"id": 1, "di": nil, "en": nil}).Replace()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"id": 2, "di": 9, "en": "a"}).Replace()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["di"].IsNil(), true)
		t.Assert(one["en"].IsNil(), true)
		one, err = db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["di"].Int(), 9)
		t.Assert(one["en"].String(), "a")
	})

	// The column type is changed after the table was already written to once. The cached
	// table fields do not follow a change by themselves, as everywhere else in gf, and
	// ClearTableFields refreshes them.
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("issue4843_alter_%d", gtime.TimestampNano())
		_, err := db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (id bigserial PRIMARY KEY, val int4)`, table))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{"id": 1, "val": 7}).Save()
		t.AssertNil(err)
		_, err = db.Exec(ctx, fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN val TYPE varchar(20)`, table))
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"id": 2, "val": "hello"}).Save()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["val"].String(), "hello")

		_, err = db.Exec(ctx, fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN val TYPE boolean USING false`, table))
		t.AssertNil(err)
		t.AssertNil(db.GetCore().ClearTableFields(ctx, table))
		_, err = db.Model(table).Data(g.Map{"id": 3, "val": true}).Save()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"id": 3, "val": false}).Save()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"id": 4, "val": true}).Save()
		t.AssertNil(err)

		one, err = db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["val"].Bool(), false)
		one, err = db.Model(table).Where("id", 4).One()
		t.AssertNil(err)
		t.Assert(one["val"].Bool(), true)
	})

	// Writing an empty string to a type that GaussDB converts from text on assignment fails
	// consistently, regardless of whether the table contains an unrelated boolean column.
	// The catalog is queried only when a column type that needs a cast is present, but the
	// types returned are filtered so that such a column is never cast: in the A compatible
	// mode a cast would turn the empty string into NULL, which silently writes NULL instead
	// of failing, and on a NOT NULL column changes the error to a not-null violation.
	gtest.C(t, func(t *gtest.T) {
		cases := []struct {
			name      string
			columnDef string
			errorText string
		}{
			{"tinyint", "v tinyint", "invalid input syntax for integer"},
			{"interval", "v interval", "invalid input syntax for type interval"},
			{"regclass", "v regclass", "invalid name syntax"},
			{"timestamp", "v timestamp", "invalid input syntax for type timestamp"},
			{"int4", "v int4", "invalid input syntax for integer"},
		}
		for _, c := range cases {
			// The column alone: no other column triggers the catalog query.
			t1 := fmt.Sprintf("issue4843_alone_%s_%d", c.name, gtime.TimestampNano())
			_, err := db.Exec(ctx, fmt.Sprintf(
				`CREATE TABLE %s (id bigserial PRIMARY KEY, %s)`, t1, c.columnDef,
			))
			t.AssertNil(err)
			defer dropTable(t1)
			_, err = db.Model(t1).Data(g.Map{"id": 1, "v": ""}).Save()
			t.AssertNE(err, nil)
			t.Assert(gstr.Contains(err.Error(), c.errorText), true)

			// The column beside a boolean, which makes the catalog query run.
			t2 := fmt.Sprintf("issue4843_bool_%s_%d", c.name, gtime.TimestampNano())
			_, err = db.Exec(ctx, fmt.Sprintf(
				`CREATE TABLE %s (id bigserial PRIMARY KEY, %s, flag boolean)`, t2, c.columnDef,
			))
			t.AssertNil(err)
			defer dropTable(t2)
			_, err = db.Model(t2).Data(g.Map{"id": 1, "v": "", "flag": true}).Save()
			t.AssertNE(err, nil)
			t.Assert(gstr.Contains(err.Error(), c.errorText), true)
		}
	})

	// The same types still accept a value that fits them.
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("issue4843_uncast_values_%d", gtime.TimestampNano())
		_, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id bigserial PRIMARY KEY, a tinyint, b interval, c regclass, d timestamp, flag boolean)`,
			table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{
			"id": 1, "a": 7, "b": "1 day", "c": "pg_class", "d": "2024-01-02 03:04:05", "flag": true,
		}).Save()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["a"].Int(), 7)
		t.Assert(one["b"].String(), "1 day")
		t.Assert(one["c"].String(), "pg_class")
		t.Assert(one["d"].String(), "2024-01-02 03:04:05")
		t.Assert(one["flag"].Bool(), true)
	})

	// A string column is not cast even when the table needs casts for other columns: in
	// the A compatible mode a cast turns an empty string into NULL, while the parameter
	// assigned as it is stays an empty string.
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("issue4843_text_%d", gtime.TimestampNano())
		_, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id bigserial PRIMARY KEY, a varchar(10) NOT NULL, b text NOT NULL, c char(3) NOT NULL, d varchar(10), flag boolean)`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{"id": 1, "a": "", "b": "", "c": "", "d": "", "flag": true}).Save()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"id": 1, "a": "", "b": "x", "c": "y", "d": "", "flag": false}).Replace()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["a"].IsNil(), false)
		t.Assert(one["a"].String(), "")
		t.Assert(one["b"].String(), "x")
		t.Assert(one["c"].String(), "y  ")
		t.Assert(one["d"].IsNil(), false)
		t.Assert(one["d"].String(), "")
		t.Assert(one["flag"].Bool(), false)
	})

	// InsertIgnore is not built as a MERGE; it keeps working and keeps ignoring.
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("issue4843_ignore_%d", gtime.TimestampNano())
		_, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id bigserial PRIMARY KEY, nickname varchar(45), val boolean)`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{"id": 1, "nickname": "n1", "val": true}).InsertIgnore()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"id": 1, "nickname": "n2", "val": false}).InsertIgnore()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "n1")
		t.Assert(one["val"].Bool(), true)
	})

	// A table that is gone is reported as such, and not as a failure to read its columns.
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("issue4843_gone_%d", gtime.TimestampNano())
		_, err := db.Model(table).Data(g.Map{"id": 1, "val": true}).Save()
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "does not exist"), true)

		// The cached columns of a table that is dropped afterwards are stale, and the
		// catalog query for them then fails on the same table.
		_, err = db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (id bigserial PRIMARY KEY, val boolean)`, table))
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"id": 1, "val": true}).Save()
		t.AssertNil(err)
		t.AssertNil(db.GetCore().ClearTableFields(ctx, table))
		_, err = db.Model(table).Data(g.Map{"id": 2, "val": true}).Save()
		t.AssertNil(err)
		dropTableWithDb(db, table)
		t.AssertNil(db.GetCore().SetTableFields(ctx, table, map[string]*gdb.TableField{
			"id":  {Name: "id", Type: "int8"},
			"val": {Name: "val", Type: "bool"},
		}))
		_, err = db.Model(table).Data(g.Map{"id": 3, "val": true}).Save()
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "does not exist"), true)
	})

	// A failing record in a batch fails the batch, with its own error.
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("issue4843_batchfail_%d", gtime.TimestampNano())
		_, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id bigserial PRIMARY KEY, nickname varchar(45) NOT NULL, val boolean)`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.List{
			{"id": 1, "nickname": "a", "val": true},
			{"id": 2, "nickname": nil, "val": true},
		}).Save()
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "nickname"), true)

		// The record before the failing one was already written, as it is on master.
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "a")
	})

	// Two columns whose names are identical when lower-cased must each get the correct cast.
	// "value" is int4 (no cast) and "Value" is boolean (needs cast); if the map were keyed
	// by lower-case name both columns would collide and the int4 column would receive the
	// boolean cast, causing Save to fail with a type error.
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("issue4843_casecolumns_%d", gtime.TimestampNano())
		_, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id bigserial PRIMARY KEY, value int4, "Value" boolean)`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{"id": 1, "value": 42, "Value": true}).Save()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["value"].Int(), 42)
		t.Assert(one["Value"].Bool(), true)

		_, err = db.Model(table).Data(g.Map{"id": 1, "value": 99, "Value": false}).Save()
		t.AssertNil(err)
		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["value"].Int(), 99)
		t.Assert(one["Value"].Bool(), false)
	})
}
