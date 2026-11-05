// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// sqlite_z_unit_feature_model_scan_basic_test.go tests scanning a single column result
// into basic types, byte slices and slices of basic types.

package sqlite_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// createScanBasicTypeTable creates a table for basic type scanning tests, where
// row 1 holds a 8-byte BLOB and row 2 holds a NULL name column.
func createScanBasicTypeTable() string {
	table := fmt.Sprintf(`scan_basic_type_%d`, gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"CREATE TABLE `%s`(id INTEGER PRIMARY KEY, name TEXT, age INTEGER, bin BLOB)", table,
	)); err != nil {
		gtest.Error(err)
	}
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"INSERT INTO `%s` VALUES(1, 'alice', 18, ?), (2, NULL, 20, ?)",
		table,
	), []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x02, 0x03, 0x04}, []byte{0x0A, 0x0B}); err != nil {
		gtest.Error(err)
	}
	return table
}

// Test_Model_Scan_ByteSlice verifies that `[]byte` and `[]uint8` targets keep the raw
// column bytes instead of being treated as a slice of scalar values.
func Test_Model_Scan_ByteSlice(t *testing.T) {
	table := createScanBasicTypeTable()
	defer dropTable(table)

	expected := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x02, 0x03, 0x04}

	gtest.C(t, func(t *gtest.T) {
		var bin []byte
		err := db.Model(table).Fields("bin").Where("id", 1).Scan(&bin)
		t.AssertNil(err)
		t.Assert(bin, expected)
	})
	gtest.C(t, func(t *gtest.T) {
		var bin []uint8
		err := db.Model(table).Fields("bin").Where("id", 1).Scan(&bin)
		t.AssertNil(err)
		t.Assert(bin, expected)
	})
	gtest.C(t, func(t *gtest.T) {
		// The result is identical to the recommended Value().Scan usage.
		value, err := db.Model(table).Fields("bin").Where("id", 1).Value()
		t.AssertNil(err)
		var bin []byte
		t.AssertNil(value.Scan(&bin))
		t.Assert(bin, expected)
	})
}

// Test_Model_Scan_BasicTypeNullValue verifies that a matching record with a NULL column
// value is distinguished from no matching record: NULL leaves the target untouched
// without error, while no record returns sql.ErrNoRows.
func Test_Model_Scan_BasicTypeNullValue(t *testing.T) {
	table := createScanBasicTypeTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// The record with id=2 exists, but its name column is NULL.
		name := "keep"
		err := db.Model(table).Fields("name").Where("id", 2).Scan(&name)
		t.AssertNil(err)
		t.Assert(name, "keep")

		// A pointer target ends up as a nil pointer for a NULL column value.
		var namePtr *string
		err = db.Model(table).Fields("name").Where("id", 2).Scan(&namePtr)
		t.AssertNil(err)
		t.AssertNil(namePtr)
	})
	gtest.C(t, func(t *gtest.T) {
		// No record matches at all.
		var name string
		err := db.Model(table).Fields("name").Where("id", 9999).Scan(&name)
		t.Assert(err, sql.ErrNoRows)

		var namePtr *string
		err = db.Model(table).Fields("name").Where("id", 9999).Scan(&namePtr)
		t.AssertNil(err)
		t.AssertNil(namePtr)
	})
}

// Test_Model_Scan_BasicTypeSliceNoRow verifies that scanning into a slice of basic types
// follows Result.Structs when no record matches: an existing slice is kept and
// sql.ErrNoRows is returned, while a nil slice is kept nil.
func Test_Model_Scan_BasicTypeSliceNoRow(t *testing.T) {
	table := createScanBasicTypeTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// A record matches, so the slice is filled and no error is returned.
		var names []string
		err := db.Model(table).Fields("name").Where("id", 1).Scan(&names)
		t.AssertNil(err)
		t.Assert(names, []string{"alice"})
	})
	gtest.C(t, func(t *gtest.T) {
		// No record matches and the target already holds values: kept and error returned,
		// which is identical to Result.Structs behavior.
		names := []string{"keep"}
		err := db.Model(table).Fields("name").Where("id", 9999).Scan(&names)
		t.Assert(err, sql.ErrNoRows)
		t.Assert(names, []string{"keep"})
	})
	gtest.C(t, func(t *gtest.T) {
		// No record matches and the target is a nil slice: it stays nil without error.
		var names []string
		err := db.Model(table).Fields("name").Where("id", 9999).Scan(&names)
		t.AssertNil(err)
		t.AssertNil(names)
	})
	gtest.C(t, func(t *gtest.T) {
		// The struct slice path reports the same thing for the same data.
		type user struct {
			Name string `orm:"name"`
		}
		users := []user{{Name: "keep"}}
		err := db.Model(table).Fields("name").Where("id", 9999).Scan(&users)
		t.Assert(err, sql.ErrNoRows)
		t.Assert(len(users), 1)
		t.Assert(users[0].Name, "keep")
	})
}
