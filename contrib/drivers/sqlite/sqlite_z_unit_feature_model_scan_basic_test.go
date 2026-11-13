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

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
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

// Test_Model_Scan_BasicTypeFieldsExJoin verifies that FieldsEx validation in join
// statements counts the columns of all joined tables, so a multi-column result is
// rejected before the query runs instead of failing on the scanned record.
func Test_Model_Scan_BasicTypeFieldsExJoin(t *testing.T) {
	table := createScanBasicTypeTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// FieldsEx leaves one column of the main table, but the joined table
		// contributes its columns too, so the result has more than one column.
		var name string
		err := db.Model(table+" a").LeftJoin(table+" b", "b.id=a.id").
			FieldsEx("a.id", "a.age", "a.bin").
			Where("a.id", 1).Scan(&name)
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "FieldsEx leaves"), true)
	})
	gtest.C(t, func(t *gtest.T) {
		// FieldsEx leaves exactly one column across all joined tables.
		var bin []byte
		err := db.Model(table+" a").LeftJoin(table+" b", "b.id=a.id").
			FieldsEx(
				"a.id", "a.name", "a.age",
				"b.id", "b.name", "b.age", "b.bin",
			).
			Where("a.id", 1).Scan(&bin)
		t.AssertNil(err)
		t.Assert(bin, []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x02, 0x03, 0x04})
	})
}

// Test_Model_Scan_BasicTypeFieldsExTableAlias verifies that FieldsEx works on a single
// table with an alias, which is still a single table operation, instead of panicking
// with "FieldsEx supports only single table operations".
func Test_Model_Scan_BasicTypeFieldsExTableAlias(t *testing.T) {
	table := createScanBasicTypeTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// "table a" is a single table with an alias, so FieldsEx must not panic.
		var age int
		err := db.Model(table+" a").FieldsEx("id", "name", "bin").Where("a.id", 1).Scan(&age)
		t.AssertNil(err)
		t.Assert(age, 18)
	})
	gtest.C(t, func(t *gtest.T) {
		// The same with an explicit AS alias.
		var name string
		err := db.Model(table+" AS a").FieldsEx("id", "age", "bin").Where("a.id", 1).Scan(&name)
		t.AssertNil(err)
		t.Assert(name, "alice")
	})
	gtest.C(t, func(t *gtest.T) {
		// FieldsEx leaving more than one column is still rejected before the query runs.
		var age int
		err := db.Model(table+" a").FieldsEx("id").Where("a.id", 1).Scan(&age)
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "FieldsEx leaves"), true)
	})
}

// Test_Model_Scan_BasicTypeSlicePointerNull verifies that a NULL column value scanned
// into a slice of pointers to basic types yields a nil element, matching the single
// value path where a NULL leaves a nil pointer.
func Test_Model_Scan_BasicTypeSlicePointerNull(t *testing.T) {
	table := createScanBasicTypeTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var names []*string
		err := db.Model(table).Fields("name").Order("id").Scan(&names)
		t.AssertNil(err)
		t.Assert(len(names), 2)
		t.AssertNE(names[0], nil)
		t.Assert(*names[0], "alice")
		// The second row has a NULL name, so its element must stay a nil pointer.
		t.AssertNil(names[1])
	})
	gtest.C(t, func(t *gtest.T) {
		// The single value path reports the same thing for the same NULL column.
		var name *string
		err := db.Model(table).Fields("name").Where("id", 2).Scan(&name)
		t.AssertNil(err)
		t.AssertNil(name)
	})
	gtest.C(t, func(t *gtest.T) {
		// Non-pointer slice elements keep the zero value for a NULL column.
		var names []string
		err := db.Model(table).Fields("name").Order("id").Scan(&names)
		t.AssertNil(err)
		t.Assert(names, []string{"alice", ""})
	})
}

// Test_Model_Scan_BasicTypeFieldsExSameField verifies that Fields and FieldsEx naming
// the same field is rejected with a clear error instead of producing an empty field
// list and invalid SQL like "SELECT  FROM ...".
func Test_Model_Scan_BasicTypeFieldsExSameField(t *testing.T) {
	table := createScanBasicTypeTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var name string
		err := db.Model(table).Fields("name").FieldsEx("name").Where("id", 1).Scan(&name)
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "excluded by FieldsEx"), true)
	})
	gtest.C(t, func(t *gtest.T) {
		// The same for a slice target.
		var names []string
		err := db.Model(table).Fields("name").FieldsEx("name").Scan(&names)
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "excluded by FieldsEx"), true)
	})
	gtest.C(t, func(t *gtest.T) {
		// FieldsEx naming another field keeps the single field valid.
		var name string
		err := db.Model(table).Fields("name").FieldsEx("age").Where("id", 1).Scan(&name)
		t.AssertNil(err)
		t.Assert(name, "alice")
	})
}

// Test_Model_Scan_BasicTypeSliceMultiColumn verifies that a projection that the
// pre-query validation cannot recognize as multi-column is still rejected after the
// query, instead of silently keeping only the first column.
func Test_Model_Scan_BasicTypeSliceMultiColumn(t *testing.T) {
	table := createScanBasicTypeTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// A quoted identifier containing a comma is a single column, so it is accepted.
		var names []string
		err := db.Model(table).Fields("`name`").Order("id").Scan(&names)
		t.AssertNil(err)
		t.Assert(names, []string{"alice", ""})
	})
	gtest.C(t, func(t *gtest.T) {
		// A raw expression that expands to two columns is rejected.
		var names []string
		err := db.Model(table).Fields(gdb.Raw("name, age")).Order("id").Scan(&names)
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "expanding field"), true)
	})
	gtest.C(t, func(t *gtest.T) {
		// A sub-query join contributes one "prefix.*" entry to the pre-query column
		// count, so the validation passes, but the actual result has several columns.
		// The post-query column-count check must reject it instead of keeping only the
		// first column.
		var names []string
		err := db.Model(table+" a").
			LeftJoin(fmt.Sprintf("(SELECT id, name, age FROM `%s`)", table), "b", "b.id=a.id").
			FieldsEx("a.id", "a.name", "a.age", "a.bin").
			Order("a.id").Scan(&names)
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "but expect one"), true)
	})
}
