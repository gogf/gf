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
	})

	gtest.C(t, func(t *gtest.T) {
		// A NULL bytea column reads back as empty rather than panicking.
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
	// Column types that are not implicitly coercible from text.
	for _, columnType := range []string{
		"numeric[]", "text[]", "jsonb", "boolean", "bytea", "uuid",
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

// Test_IssuePointInterval_NotConvertedToInt verifies that column types whose
// names merely contain an integer keyword are not converted to integers.
//
// gdb's fallback type detection matches "int" as a substring, so "point" and
// "interval" were both classified as integers and their values became 0.
func Test_IssuePointInterval_NotConvertedToInt(t *testing.T) {
	table := "issue_point_interval_" + gtime.TimestampMicroStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(
		`CREATE TABLE %s (id int PRIMARY KEY, pt point, span interval)`, table,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			`INSERT INTO %s VALUES (1, '(1.5,2.5)', '2 days')`, table,
		))
		t.AssertNil(err)

		// Read through the ORM and compare against an explicit ::text cast,
		// which is unaffected by the type detection.
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)

		want, err := db.Model(table).
			Fields("pt::text as pt, span::text as span").Where("id", 1).One()
		t.AssertNil(err)

		t.Assert(one["pt"].String(), want["pt"].String())
		t.Assert(one["span"].String(), want["span"].String())

		// And the values are the real ones, not zeroes.
		t.Assert(one["pt"].String(), "(1.5,2.5)")
		t.AssertNE(one["span"].String(), "0")
	})
}
