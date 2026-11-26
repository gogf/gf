// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file tests the partition clause of Model.Partition(): where Core inserts it into the
// table expression, and how a write statement carrying partition names is refused.

package gdb

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

// partitionFakeDB is a dialectFakeDB that emits a partition clause for the names given to it,
// recording the names of each call, or no clause at all without names.
type partitionFakeDB struct {
	*dialectFakeDB
	supported bool
	calls     [][]string
}

func newPartitionFakeDB(supported bool, columns ...string) *partitionFakeDB {
	db := &partitionFakeDB{
		dialectFakeDB: newDialectFakeDB(columns...),
		supported:     supported,
	}
	db.Core.db = db
	return db
}

// FormatPartitionClause renders the names as a driver with partition selection syntax would,
// or returns an empty clause as a driver without it does.
func (d *partitionFakeDB) FormatPartitionClause(partitions []string) string {
	d.calls = append(d.calls, partitions)
	if !d.supported || len(partitions) == 0 {
		return ""
	}
	return "PARTITION (" + gstr.Join(partitions, ",") + ")"
}

// Test_Model_Partition_ClausePlacement asserts that the partition clause is inserted right
// after the first table name, so that an alias or a JOIN that follows it stays valid.
func Test_Model_Partition_ClausePlacement(t *testing.T) {
	var ctx = context.Background()
	gtest.C(t, func(t *gtest.T) {
		for _, c := range []struct {
			supported bool
			model     func(db *partitionFakeDB) *Model
			expect    string
		}{
			{
				true,
				func(db *partitionFakeDB) *Model { return db.Model("user").Partition("p0") },
				`SELECT * FROM "user" PARTITION (p0)`,
			},
			{
				true,
				func(db *partitionFakeDB) *Model {
					return db.Model("user").Partition("p0", "p1")
				},
				`SELECT * FROM "user" PARTITION (p0,p1)`,
			},
			{
				true,
				func(db *partitionFakeDB) *Model { return db.Model("user").Partition("p0").As("u") },
				`SELECT * FROM "user" PARTITION (p0) AS u`,
			},
			{
				true,
				func(db *partitionFakeDB) *Model {
					return db.Model("user").Partition("p0").As("u").LeftJoin("detail", "detail.id=u.id")
				},
				`SELECT * FROM "user" PARTITION (p0) AS u LEFT JOIN "detail" ON (detail.id=u.id)`,
			},
			{
				// The clause stays on the first table of a comma-separated list.
				true,
				func(db *partitionFakeDB) *Model { return db.Model("user,detail").Partition("p0") },
				`SELECT * FROM "user" PARTITION (p0),"detail"`,
			},
			{
				// As() aliases the last table of the list, and the clause stays on the first.
				true,
				func(db *partitionFakeDB) *Model { return db.Model("user,detail").Partition("p0").As("u") },
				`SELECT * FROM "user" PARTITION (p0),"detail" AS u`,
			},
			{
				// A model without partition names emits the table expression unchanged.
				true,
				func(db *partitionFakeDB) *Model { return db.Model("user") },
				`SELECT * FROM "user"`,
			},
			{
				// A driver without partition selection syntax leaves the names unused.
				false,
				func(db *partitionFakeDB) *Model { return db.Model("user").Partition("p0") },
				`SELECT * FROM "user"`,
			},
		} {
			db := newPartitionFakeDB(c.supported, "id", "name")
			sqlWithHolder, _ := c.model(db).getFormattedSqlAndArgs(ctx, SelectTypeDefault, false)
			t.Assert(sqlWithHolder, c.expect)
		}
	})
}

// Test_Model_Partition_UnsupportedDriver asserts that a driver without partition selection
// syntax leaves the names unused instead of emitting a clause its database would misread.
func Test_Model_Partition_UnsupportedDriver(t *testing.T) {
	var ctx = context.Background()
	gtest.C(t, func(t *gtest.T) {
		db := newPartitionFakeDB(false, "id", "name")
		sqlWithHolder, _ := db.Model("user").Partition("p0").getFormattedSqlAndArgs(
			ctx, SelectTypeDefault, false,
		)
		t.Assert(sqlWithHolder, `SELECT * FROM "user"`)
		t.Assert(len(db.calls), 1)
	})
}

// Test_Model_Partition_Count asserts that Count keeps the partition clause, including when the
// rows are counted through a derived table for a grouped model.
func Test_Model_Partition_Count(t *testing.T) {
	var ctx = context.Background()
	gtest.C(t, func(t *gtest.T) {
		for _, c := range []struct {
			model  func(db *partitionFakeDB) *Model
			expect string
		}{
			{
				func(db *partitionFakeDB) *Model { return db.Model("user").Partition("p0") },
				`SELECT COUNT(1) FROM "user" PARTITION (p0)`,
			},
			{
				func(db *partitionFakeDB) *Model {
					return db.Model("user").Partition("p0").Group("name")
				},
				`SELECT COUNT(1) FROM (SELECT COUNT(1) FROM "user" PARTITION (p0) GROUP BY "name") count_alias`,
			},
		} {
			db := newPartitionFakeDB(true, "id", "name")
			sqlWithHolder, _ := c.model(db).getFormattedSqlAndArgs(ctx, SelectTypeCount, false)
			t.Assert(sqlWithHolder, c.expect)
		}
	})
}

// Test_Model_Partition_WriteRefused asserts that a write statement carrying partition names is
// refused instead of executed without the restriction, which would silently affect the rows of
// every partition.
func Test_Model_Partition_WriteRefused(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newPartitionFakeDB(true, "id", "name")
		_, err := db.Model("user").Partition("p0").Where("id", 1).Data("name", "john").Update()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)
		_, err = db.Model("user").Partition("p0").Where("id", 1).Delete()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)
		_, err = db.Model("user").Partition("p0").Delete()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)

		// A model without partition names is not refused.
		t.AssertNil(db.Model("user").Partition().checkPartitionForWrite("DELETE"))
	})
}

// Test_Model_Partition_InsertRefused asserts that every insert statement carrying partition
// names is refused as well: Replace and Save can update or delete rows of a partition the
// caller did not name, which the SELECT-only refusal left open.
func Test_Model_Partition_InsertRefused(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newPartitionFakeDB(true, "id", "name")
		_, err := db.Model("user").Partition("p0").Data("name", "john").Insert()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)
		_, err = db.Model("user").Partition("p0").Data("name", "john").InsertIgnore()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)
		_, err = db.Model("user").Partition("p0").Data("name", "john").InsertAndGetId()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)
		_, err = db.Model("user").Partition("p0").Data("name", "john").Replace()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)
		_, err = db.Model("user").Partition("p0").Data("name", "john").Save()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)

		// The data variants reach the same refusal.
		_, err = db.Model("user").Partition("p0").Insert("name", "john")
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)
	})
}

// Test_Model_Partition_WriteNotRefusedWithoutClause asserts that a driver without partition
// selection syntax is not refused on a write: the names never reach the statement on such a
// driver, for SELECT either, so Partition() is inert rather than half-applied.
func Test_Model_Partition_WriteNotRefusedWithoutClause(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newPartitionFakeDB(false, "id", "name")
		t.AssertNil(db.Model("user").Partition("p0").checkPartitionForWrite("INSERT"))
		t.AssertNil(db.Model("user").Partition("p0").checkPartitionForWrite("UPDATE"))
		t.AssertNil(db.Model("user").Partition("p0").checkPartitionForWrite("DELETE"))

		coreDB := newDialectFakeDB("id", "name")
		t.AssertNil(coreDB.Model("user").Partition("p0").checkPartitionForWrite("DELETE"))
	})
}

// Test_Model_Partition_CoreDefault asserts that the Core default FormatPartitionClause
// returns an empty string, which causes the table expression to be emitted unchanged.
func Test_Model_Partition_CoreDefault(t *testing.T) {
	var ctx = context.Background()
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "name")
		sqlWithHolder, _ := db.Model("user").Partition("p0").getFormattedSqlAndArgs(ctx, SelectTypeDefault, false)
		t.Assert(sqlWithHolder, `SELECT * FROM "user"`)
	})
}

// Test_Model_Partition_CloneDeep asserts that Clone copies partition names into an
// independent slice, so a subsequent Partition call on the clone does not affect the
// original model.
func Test_Model_Partition_CloneDeep(t *testing.T) {
	var ctx = context.Background()
	gtest.C(t, func(t *gtest.T) {
		db := newPartitionFakeDB(true, "id", "name")
		original := db.Model("user").Safe().Partition("p0")
		cloned := original.Clone()
		cloned.partition = []string{"p1"}

		sqlOrig, _ := original.getFormattedSqlAndArgs(ctx, SelectTypeDefault, false)
		sqlCloned, _ := cloned.getFormattedSqlAndArgs(ctx, SelectTypeDefault, false)
		t.Assert(sqlOrig, `SELECT * FROM "user" PARTITION (p0)`)
		t.Assert(sqlCloned, `SELECT * FROM "user" PARTITION (p1)`)
	})
}
