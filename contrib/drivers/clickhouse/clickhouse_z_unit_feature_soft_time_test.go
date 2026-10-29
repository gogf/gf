// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// Note: the effects of writes are verified by reading the data back instead of by the affected rows, and Save
// and Replace insert the records, which the ReplacingMergeTree tables deduplicate when their parts are merged
// by chDOptimizeFinal.

// Test_SoftTime_CreateUpdateDelete1 tests the create_at, update_at and delete_at fields with a map.
func Test_SoftTime_CreateUpdateDelete1(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  create_at Nullable(DateTime64(6)),
  update_at Nullable(DateTime64(6)),
  delete_at Nullable(DateTime64(6))
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert
		dataInsert := g.Map{
			"id":   1,
			"name": "name_1",
		}
		_, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)

		oneInsert, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["name"].String(), "name_1")
		t.Assert(oneInsert["delete_at"].String(), "")
		t.AssertGE(oneInsert["create_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		t.AssertGE(oneInsert["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Save
		dataSave := g.Map{
			"id":   1,
			"name": "name_10",
		}
		_, err = db.Model(table).Data(dataSave).Save()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		oneSave, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["delete_at"].String(), "")
		// Note: Save inserts a new record, which the ReplacingMergeTree table keeps instead of the old one,
		// so the creation time is the time of Save.
		t.AssertGT(oneSave["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertNE(oneSave["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())
		t.AssertGE(oneSave["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := g.Map{
			"name": "name_1000",
		}
		_, err = db.Model(table).Data(dataUpdate).WherePri(1).Update()
		t.AssertNil(err)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["delete_at"].String(), "")
		t.Assert(oneUpdate["create_at"].GTime().Timestamp(), oneSave["create_at"].GTime().Timestamp())
		t.AssertGE(oneUpdate["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// Replace
		dataReplace := g.Map{
			"id":   1,
			"name": "name_100",
		}
		_, err = db.Model(table).Data(dataReplace).Replace()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		oneReplace, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneReplace["id"].Int(), 1)
		t.Assert(oneReplace["name"].String(), "name_100")
		t.Assert(oneReplace["delete_at"].String(), "")
		t.AssertGE(oneReplace["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertGE(oneReplace["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Delete
		_, err = db.Model(table).Delete("id", 1)
		t.AssertNil(err)
		// Delete Select
		one4, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one4), 0)
		one5, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one5["id"].Int(), 1)
		t.AssertGE(one5["delete_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		// Delete Count
		i, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(i, 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 1)

		// Delete Unscoped
		_, err = db.Model(table).Unscoped().Delete("id", 1)
		t.AssertNil(err)
		one6, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one6), 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 0)
	})
}

// Test_SoftTime_CreateUpdateDelete2 tests the create_at, update_at and delete_at fields with a map again.
func Test_SoftTime_CreateUpdateDelete2(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  create_at Nullable(DateTime),
  update_at Nullable(DateTime),
  delete_at Nullable(DateTime)
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert
		dataInsert := g.Map{
			"id":   1,
			"name": "name_1",
		}
		_, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)

		oneInsert, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["name"].String(), "name_1")
		t.Assert(oneInsert["delete_at"].String(), "")
		t.AssertGE(oneInsert["create_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		t.AssertGE(oneInsert["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Save
		dataSave := g.Map{
			"id":   1,
			"name": "name_10",
		}
		_, err = db.Model(table).Data(dataSave).Save()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		oneSave, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["delete_at"].String(), "")
		// Note: Save inserts a new record, which the ReplacingMergeTree table keeps instead of the old one,
		// so the creation time is the time of Save.
		t.AssertGT(oneSave["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertNE(oneSave["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())
		t.AssertGE(oneSave["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := g.Map{
			"name": "name_1000",
		}
		_, err = db.Model(table).Data(dataUpdate).WherePri(1).Update()
		t.AssertNil(err)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["delete_at"].String(), "")
		t.Assert(oneUpdate["create_at"].GTime().Timestamp(), oneSave["create_at"].GTime().Timestamp())
		t.AssertGE(oneUpdate["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// Replace
		dataReplace := g.Map{
			"id":   1,
			"name": "name_100",
		}
		_, err = db.Model(table).Data(dataReplace).Replace()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		oneReplace, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneReplace["id"].Int(), 1)
		t.Assert(oneReplace["name"].String(), "name_100")
		t.Assert(oneReplace["delete_at"].String(), "")
		t.AssertGE(oneReplace["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertGE(oneReplace["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Delete
		_, err = db.Model(table).Delete("id", 1)
		t.AssertNil(err)
		// Delete Select
		one4, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one4), 0)
		one5, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one5["id"].Int(), 1)
		t.AssertGE(one5["delete_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		// Delete Count
		i, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(i, 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 1)

		// Delete Unscoped
		_, err = db.Model(table).Unscoped().Delete("id", 1)
		t.AssertNil(err)
		one6, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one6), 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 0)
	})
}

// Test_SoftTime_CreatedUpdatedDeleted_Map tests the created_at, updated_at and deleted_at fields with a map.
func Test_SoftTime_CreatedUpdatedDeleted_Map(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  created_at Nullable(DateTime64(6)),
  updated_at Nullable(DateTime64(6)),
  deleted_at Nullable(DateTime64(6))
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert
		dataInsert := g.Map{
			"id":   1,
			"name": "name_1",
		}
		_, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)

		oneInsert, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["name"].String(), "name_1")
		t.Assert(oneInsert["deleted_at"].String(), "")
		t.AssertGE(oneInsert["created_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		t.AssertGE(oneInsert["updated_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Save
		dataSave := g.Map{
			"id":   1,
			"name": "name_10",
		}
		_, err = db.Model(table).Data(dataSave).Save()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		oneSave, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["deleted_at"].String(), "")
		// Note: Save inserts a new record, which the ReplacingMergeTree table keeps instead of the old one,
		// so the creation time is the time of Save.
		t.AssertGT(oneSave["created_at"].GTime().Timestamp(), oneInsert["created_at"].GTime().Timestamp())
		t.AssertNE(oneSave["updated_at"].GTime().Timestamp(), oneInsert["updated_at"].GTime().Timestamp())
		t.AssertGE(oneSave["updated_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := g.Map{
			"name": "name_1000",
		}
		_, err = db.Model(table).Data(dataUpdate).WherePri(1).Update()
		t.AssertNil(err)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["deleted_at"].String(), "")
		t.Assert(oneUpdate["created_at"].GTime().Timestamp(), oneSave["created_at"].GTime().Timestamp())
		t.AssertGE(oneUpdate["updated_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// Replace
		dataReplace := g.Map{
			"id":   1,
			"name": "name_100",
		}
		_, err = db.Model(table).Data(dataReplace).Replace()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		oneReplace, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneReplace["id"].Int(), 1)
		t.Assert(oneReplace["name"].String(), "name_100")
		t.Assert(oneReplace["deleted_at"].String(), "")
		t.AssertGE(oneReplace["created_at"].GTime().Timestamp(), oneInsert["created_at"].GTime().Timestamp())
		t.AssertGE(oneReplace["updated_at"].GTime().Timestamp(), oneInsert["updated_at"].GTime().Timestamp())

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Delete
		_, err = db.Model(table).Delete("id", 1)
		t.AssertNil(err)
		// Delete Select
		one4, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one4), 0)
		one5, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one5["id"].Int(), 1)
		t.AssertGE(one5["deleted_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		// Delete Count
		i, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(i, 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 1)

		// Delete Unscoped
		_, err = db.Model(table).Unscoped().Delete("id", 1)
		t.AssertNil(err)
		one6, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one6), 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 0)
	})
}

// Test_SoftTime_CreatedUpdatedDeleted_Struct tests the created_at, updated_at and deleted_at fields with a struct.
func Test_SoftTime_CreatedUpdatedDeleted_Struct(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  created_at Nullable(DateTime64(6)),
  updated_at Nullable(DateTime64(6)),
  deleted_at Nullable(DateTime64(6))
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	type User struct {
		Id        int
		Name      string
		CreatedAT *gtime.Time
		UpdatedAT *gtime.Time
		DeletedAT *gtime.Time
	}
	gtest.C(t, func(t *gtest.T) {
		// Insert
		dataInsert := User{
			Id:   1,
			Name: "name_1",
		}
		_, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)

		oneInsert, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["name"].String(), "name_1")
		t.Assert(oneInsert["deleted_at"].String(), "")
		t.AssertGE(oneInsert["created_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		t.AssertGE(oneInsert["updated_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Save
		dataSave := User{
			Id:   1,
			Name: "name_10",
		}
		_, err = db.Model(table).Data(dataSave).OmitEmpty().Save()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		oneSave, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["deleted_at"].String(), "")
		// Note: Save inserts a new record, which the ReplacingMergeTree table keeps instead of the old one,
		// so the creation time is the time of Save.
		t.AssertGT(oneSave["created_at"].GTime().Timestamp(), oneInsert["created_at"].GTime().Timestamp())
		t.AssertNE(oneSave["updated_at"].GTime().Timestamp(), oneInsert["updated_at"].GTime().Timestamp())
		t.AssertGE(oneSave["updated_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := User{
			Name: "name_1000",
		}
		_, err = db.Model(table).Data(dataUpdate).OmitEmpty().WherePri(1).Update()
		t.AssertNil(err)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["deleted_at"].String(), "")
		t.Assert(oneUpdate["created_at"].GTime().Timestamp(), oneSave["created_at"].GTime().Timestamp())
		t.AssertGE(oneUpdate["updated_at"].GTime().Timestamp(), gtime.Timestamp()-4)

		// Replace
		dataReplace := User{
			Id:   1,
			Name: "name_100",
		}
		_, err = db.Model(table).Data(dataReplace).OmitEmpty().Replace()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		oneReplace, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneReplace["id"].Int(), 1)
		t.Assert(oneReplace["name"].String(), "name_100")
		t.Assert(oneReplace["deleted_at"].String(), "")
		t.AssertGE(oneReplace["created_at"].GTime().Timestamp(), oneInsert["created_at"].GTime().Timestamp())
		t.AssertGE(oneReplace["updated_at"].GTime().Timestamp(), oneInsert["updated_at"].GTime().Timestamp())

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Delete
		_, err = db.Model(table).Delete("id", 1)
		t.AssertNil(err)
		// Delete Select
		one4, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one4), 0)
		one5, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one5["id"].Int(), 1)
		t.AssertGE(one5["deleted_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		// Delete Count
		i, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(i, 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 1)

		// Delete Unscoped
		_, err = db.Model(table).Unscoped().Delete("id", 1)
		t.AssertNil(err)
		one6, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one6), 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 0)
	})
}

// Test_SoftUpdateTime tests updating by a string data on a table with soft time fields.
func Test_SoftUpdateTime(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  num       Nullable(Int32),
  create_at Nullable(DateTime64(6)),
  update_at Nullable(DateTime64(6)),
  delete_at Nullable(DateTime64(6))
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert
		dataInsert := g.Map{
			"id":  1,
			"num": 10,
		}
		_, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)

		oneInsert, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["num"].Int(), 10)

		// Update.
		_, err = db.Model(table).Data("num=num+1").Where("id=?", 1).Update()
		t.AssertNil(err)
	})
}

// Test_SoftUpdateTime_WithDO tests that updating with a DO struct sets the update time only.
func Test_SoftUpdateTime_WithDO(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  num       Nullable(Int32),
  created_at Nullable(DateTime64(6)),
  updated_at Nullable(DateTime64(6)),
  deleted_at Nullable(DateTime64(6))
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert
		dataInsert := g.Map{
			"id":  1,
			"num": 10,
		}
		_, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)

		oneInserted, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInserted["id"].Int(), 1)
		t.Assert(oneInserted["num"].Int(), 10)

		// Update.
		time.Sleep(2 * time.Second)
		type User struct {
			g.Meta    `orm:"do:true"`
			Id        any
			Num       any
			CreatedAt any
			UpdatedAt any
			DeletedAt any
		}
		_, err = db.Model(table).Data(User{
			Num: 100,
		}).Where("id=?", 1).Update()
		t.AssertNil(err)

		oneUpdated, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdated["num"].Int(), 100)
		t.Assert(oneUpdated["created_at"].String(), oneInserted["created_at"].String())
		t.AssertNE(oneUpdated["updated_at"].String(), oneInserted["updated_at"].String())
	})
}

// Test_SoftDelete tests soft deleting records.
func Test_SoftDelete(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  create_at Nullable(DateTime64(6)),
  update_at Nullable(DateTime64(6)),
  delete_at Nullable(DateTime64(6))
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		for i := 1; i <= 10; i++ {
			data := g.Map{
				"id":   i,
				"name": fmt.Sprintf("name_%d", i),
			}
			_, err := db.Model(table).Data(data).Insert()
			t.AssertNil(err)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.AssertNE(one["create_at"].String(), "")
		t.AssertNE(one["update_at"].String(), "")
		t.Assert(one["delete_at"].String(), "")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).WherePri(10).One()
		t.AssertNil(err)
		t.AssertNE(one["create_at"].String(), "")
		t.AssertNE(one["update_at"].String(), "")
		t.Assert(one["delete_at"].String(), "")
	})
	gtest.C(t, func(t *gtest.T) {
		ids := g.SliceInt{1, 3, 5}
		_, err := db.Model(table).Where("id", ids).Delete()
		t.AssertNil(err)

		count, err := db.Model(table).Where("id", ids).Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		all, err := db.Model(table).Unscoped().Where("id", ids).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.AssertNE(all[0]["create_at"].String(), "")
		t.AssertNE(all[0]["update_at"].String(), "")
		t.AssertNE(all[0]["delete_at"].String(), "")
		t.AssertNE(all[1]["create_at"].String(), "")
		t.AssertNE(all[1]["update_at"].String(), "")
		t.AssertNE(all[1]["delete_at"].String(), "")
		t.AssertNE(all[2]["create_at"].String(), "")
		t.AssertNE(all[2]["update_at"].String(), "")
		t.AssertNE(all[2]["delete_at"].String(), "")
	})
}

// Test_SoftDelete_Join tests that soft deleted records are excluded from joins.
func Test_SoftDelete_Join(t *testing.T) {
	table1 := "time_test_table1_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  create_at Nullable(DateTime64(6)),
  update_at Nullable(DateTime64(6)),
  delete_at Nullable(DateTime64(6))
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table1)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table1)

	table2 := "time_test_table2_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  createat Nullable(DateTime64(6)),
  updateat Nullable(DateTime64(6)),
  deleteat Nullable(DateTime64(6))
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table2)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		dataInsert1 := g.Map{
			"id":   1,
			"name": "name_1",
		}
		_, err := db.Model(table1).Data(dataInsert1).Insert()
		t.AssertNil(err)

		dataInsert2 := g.Map{
			"id":   1,
			"name": "name_2",
		}
		_, err = db.Model(table2).Data(dataInsert2).Insert()
		t.AssertNil(err)

		one, err := db.Model(table1, "t1").LeftJoin(table2, "t2", "t2.id=t1.id").Fields("t1.name").One()
		t.AssertNil(err)
		t.Assert(one["name"], "name_1")

		// Soft deleting.
		_, err = db.Model(table1).Where(1).Delete()
		t.AssertNil(err)

		one, err = db.Model(table1, "t1").LeftJoin(table2, "t2", "t2.id=t1.id").Fields("t1.name").One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), true)

		one, err = db.Model(table2, "t2").LeftJoin(table1, "t1", "t2.id=t1.id").Fields("t2.name").One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), true)
	})
}

// Test_SoftDelete_WhereAndOr tests the soft delete condition with Where and WhereOr.
func Test_SoftDelete_WhereAndOr(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  create_at Nullable(DateTime64(6)),
  update_at Nullable(DateTime64(6)),
  delete_at Nullable(DateTime64(6))
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)
	// Add datas.
	gtest.C(t, func(t *gtest.T) {
		for i := 1; i <= 10; i++ {
			data := g.Map{
				"id":   i,
				"name": fmt.Sprintf("name_%d", i),
			}
			_, err := db.Model(table).Data(data).Insert()
			t.AssertNil(err)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		ids := g.SliceInt{1, 3, 5}
		_, err := db.Model(table).Where("id", ids).Delete()
		t.AssertNil(err)

		count, err := db.Model(table).Where("id", 1).WhereOr("id", 3).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_CreateUpdateTime_Struct tests the soft time fields with a struct of pointer time fields.
func Test_CreateUpdateTime_Struct(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  create_at Nullable(DateTime64(6)),
  update_at Nullable(DateTime64(6)),
  delete_at Nullable(DateTime64(6))
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	type Entity struct {
		Id       uint64      `orm:"id,primary" json:"id"`
		Name     string      `orm:"name"       json:"name"`
		CreateAt *gtime.Time `orm:"create_at"  json:"create_at"`
		UpdateAt *gtime.Time `orm:"update_at"  json:"update_at"`
		DeleteAt *gtime.Time `orm:"delete_at"  json:"delete_at"`
	}
	gtest.C(t, func(t *gtest.T) {
		// Insert
		dataInsert := &Entity{
			Id:       1,
			Name:     "name_1",
			CreateAt: nil,
			UpdateAt: nil,
			DeleteAt: nil,
		}
		_, err := db.Model(table).Data(dataInsert).OmitEmpty().Insert()
		t.AssertNil(err)

		oneInsert, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["name"].String(), "name_1")
		t.Assert(oneInsert["delete_at"].String(), "")
		t.AssertGE(oneInsert["create_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		t.AssertGE(oneInsert["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		time.Sleep(2 * time.Second)

		// Save
		dataSave := &Entity{
			Id:       1,
			Name:     "name_10",
			CreateAt: nil,
			UpdateAt: nil,
			DeleteAt: nil,
		}
		_, err = db.Model(table).Data(dataSave).OmitEmpty().Save()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		oneSave, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["delete_at"].String(), "")
		// Note: Save inserts a new record, which the ReplacingMergeTree table keeps instead of the old one,
		// so the creation time is the time of Save.
		t.AssertGT(oneSave["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertNE(oneSave["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())
		t.AssertGE(oneSave["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := &Entity{
			Id:       1,
			Name:     "name_1000",
			CreateAt: nil,
			UpdateAt: nil,
			DeleteAt: nil,
		}
		// Note: ClickHouse cannot UPDATE a column of the sorting key, so setting id fails.
		_, err = db.Model(table).Data(dataUpdate).WherePri(1).OmitEmpty().Update()
		t.AssertNE(err, nil)

		dataUpdate.Id = 0
		_, err = db.Model(table).Data(dataUpdate).WherePri(1).OmitEmpty().Update()
		t.AssertNil(err)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["delete_at"].String(), "")
		t.Assert(oneUpdate["create_at"].GTime().Timestamp(), oneSave["create_at"].GTime().Timestamp())
		t.AssertGE(oneUpdate["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// Replace
		dataReplace := &Entity{
			Id:       1,
			Name:     "name_100",
			CreateAt: nil,
			UpdateAt: nil,
			DeleteAt: nil,
		}
		_, err = db.Model(table).Data(dataReplace).OmitEmpty().Replace()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		oneReplace, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneReplace["id"].Int(), 1)
		t.Assert(oneReplace["name"].String(), "name_100")
		t.Assert(oneReplace["delete_at"].String(), "")
		t.AssertGE(oneReplace["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertGE(oneReplace["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())

		time.Sleep(2 * time.Second)

		// Delete
		_, err = db.Model(table).Delete("id", 1)
		t.AssertNil(err)
		// Delete Select
		one4, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one4), 0)
		one5, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one5["id"].Int(), 1)
		t.AssertGE(one5["delete_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		// Delete Count
		i, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(i, 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 1)

		// Delete Unscoped
		_, err = db.Model(table).Unscoped().Delete("id", 1)
		t.AssertNil(err)
		one6, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one6), 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 0)
	})
}

// Test_SoftTime_CreateUpdateDelete_UnixTimestamp tests soft time fields of integer unix timestamps.
func Test_SoftTime_CreateUpdateDelete_UnixTimestamp(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  create_at Nullable(Int32),
  update_at Nullable(Int32),
  delete_at Nullable(Int32)
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	// insert
	gtest.C(t, func(t *gtest.T) {
		dataInsert := g.Map{
			"id":   1,
			"name": "name_1",
		}
		_, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_1")
		t.AssertGT(one["create_at"].Int64(), 0)
		t.AssertGT(one["update_at"].Int64(), 0)
		t.Assert(one["delete_at"].Int64(), 0)
		t.Assert(len(one["create_at"].String()), 10)
		t.Assert(len(one["update_at"].String()), 10)
	})

	// sleep some seconds to make update time greater than create time.
	time.Sleep(2 * time.Second)

	// update
	gtest.C(t, func(t *gtest.T) {
		// update: map
		dataInsert := g.Map{
			"name": "name_11",
		}
		_, err := db.Model(table).Data(dataInsert).WherePri(1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_11")
		t.AssertGT(one["create_at"].Int64(), 0)
		t.AssertGT(one["update_at"].Int64(), 0)
		t.Assert(one["delete_at"].Int64(), 0)
		t.Assert(len(one["create_at"].String()), 10)
		t.Assert(len(one["update_at"].String()), 10)

		var (
			lastCreateTime = one["create_at"].Int64()
			lastUpdateTime = one["update_at"].Int64()
		)

		time.Sleep(2 * time.Second)

		// update: string
		_, err = db.Model(table).Data("name='name_111'").WherePri(1).Update()
		t.AssertNil(err)

		one, err = db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_111")
		t.Assert(one["create_at"].Int64(), lastCreateTime)
		t.AssertGT(one["update_at"].Int64(), lastUpdateTime)
		t.Assert(one["delete_at"].Int64(), 0)
	})

	// delete
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).WherePri(1).Delete()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one), 0)

		one, err = db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_111")
		t.AssertGT(one["create_at"].Int64(), 0)
		t.AssertGT(one["update_at"].Int64(), 0)
		t.AssertGT(one["delete_at"].Int64(), 0)
	})
}

// Test_SoftTime_CreateUpdateDelete_Bool_Deleted tests a boolean soft delete field.
func Test_SoftTime_CreateUpdateDelete_Bool_Deleted(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  create_at Nullable(Int32),
  update_at Nullable(Int32),
  delete_at Nullable(Bool)
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	// insert
	gtest.C(t, func(t *gtest.T) {
		dataInsert := g.Map{
			"id":   1,
			"name": "name_1",
		}
		_, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_1")
		t.AssertGT(one["create_at"].Int64(), 0)
		t.AssertGT(one["update_at"].Int64(), 0)
		t.Assert(one["delete_at"].Int64(), 0)
		t.Assert(len(one["create_at"].String()), 10)
		t.Assert(len(one["update_at"].String()), 10)
	})

	// delete
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).WherePri(1).Delete()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one), 0)

		one, err = db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_1")
		t.AssertGT(one["create_at"].Int64(), 0)
		t.AssertGT(one["update_at"].Int64(), 0)
		t.Assert(one["delete_at"].Int64(), 1)
	})
}

// Test_SoftTime_CreateUpdateDelete_Option_SoftTimeTypeTimestampMilli tests the SoftTimeTypeTimestampMilli option.
func Test_SoftTime_CreateUpdateDelete_Option_SoftTimeTypeTimestampMilli(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  create_at Nullable(UInt64),
  update_at Nullable(UInt64),
  delete_at Nullable(Bool)
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	var softTimeOption = gdb.SoftTimeOption{
		SoftTimeType: gdb.SoftTimeTypeTimestampMilli,
	}

	// insert
	gtest.C(t, func(t *gtest.T) {
		dataInsert := g.Map{
			"id":   1,
			"name": "name_1",
		}
		_, err := db.Model(table).SoftTime(softTimeOption).Data(dataInsert).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).SoftTime(softTimeOption).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_1")
		t.Assert(len(one["create_at"].String()), 13)
		t.Assert(len(one["update_at"].String()), 13)
		t.Assert(one["delete_at"].Int64(), 0)
	})

	// delete
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).SoftTime(softTimeOption).WherePri(1).Delete()
		t.AssertNil(err)

		one, err := db.Model(table).SoftTime(softTimeOption).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one), 0)

		one, err = db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_1")
		t.AssertGT(one["create_at"].Int64(), 0)
		t.AssertGT(one["update_at"].Int64(), 0)
		t.Assert(one["delete_at"].Int64(), 1)
	})
}

// Test_SoftTime_CreateUpdateDelete_Option_SoftTimeTypeTimestampNano tests the SoftTimeTypeTimestampNano option.
func Test_SoftTime_CreateUpdateDelete_Option_SoftTimeTypeTimestampNano(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  create_at Nullable(UInt64),
  update_at Nullable(UInt64),
  delete_at Nullable(Bool)
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	var softTimeOption = gdb.SoftTimeOption{
		SoftTimeType: gdb.SoftTimeTypeTimestampNano,
	}

	// insert
	gtest.C(t, func(t *gtest.T) {
		dataInsert := g.Map{
			"id":   1,
			"name": "name_1",
		}
		_, err := db.Model(table).SoftTime(softTimeOption).Data(dataInsert).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).SoftTime(softTimeOption).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_1")
		t.Assert(len(one["create_at"].String()), 19)
		t.Assert(len(one["update_at"].String()), 19)
		t.Assert(one["delete_at"].Int64(), 0)
	})

	// delete
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).SoftTime(softTimeOption).WherePri(1).Delete()
		t.AssertNil(err)

		one, err := db.Model(table).SoftTime(softTimeOption).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one), 0)

		one, err = db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_1")
		t.AssertGT(one["create_at"].Int64(), 0)
		t.AssertGT(one["update_at"].Int64(), 0)
		t.Assert(one["delete_at"].Int64(), 1)
	})
}

// Test_SoftTime_CreateUpdateDelete_Specified tests soft time fields given in the data.
func Test_SoftTime_CreateUpdateDelete_Specified(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        Int32 NOT NULL,
  name      Nullable(String),
  create_at Nullable(DateTime),
  update_at Nullable(DateTime),
  delete_at Nullable(DateTime)
) ENGINE = ReplacingMergeTree()
ORDER BY id;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert
		dataInsert := g.Map{
			"id":        1,
			"name":      "name_1",
			"create_at": gtime.NewFromStr("2024-05-30 20:00:00"),
			"update_at": gtime.NewFromStr("2024-05-30 20:00:00"),
		}
		_, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)

		oneInsert, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["name"].String(), "name_1")
		t.Assert(oneInsert["delete_at"].String(), "")
		t.Assert(oneInsert["create_at"].String(), "2024-05-30 20:00:00")
		t.Assert(oneInsert["update_at"].String(), "2024-05-30 20:00:00")

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Save
		dataSave := g.Map{
			"id":        1,
			"name":      "name_10",
			"update_at": gtime.NewFromStr("2024-05-30 20:15:00"),
		}
		_, err = db.Model(table).Data(dataSave).Save()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		oneSave, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["delete_at"].String(), "")
		// Note: Save inserts a new record, which the ReplacingMergeTree table keeps instead of the old one,
		// so the creation time is the time of Save.
		t.AssertGE(oneSave["create_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		t.Assert(oneSave["update_at"].String(), "2024-05-30 20:15:00")

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := g.Map{
			"name":      "name_1000",
			"update_at": gtime.NewFromStr("2024-05-30 20:30:00"),
		}
		_, err = db.Model(table).Data(dataUpdate).WherePri(1).Update()
		t.AssertNil(err)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["delete_at"].String(), "")
		t.Assert(oneUpdate["create_at"].String(), oneSave["create_at"].String())
		t.Assert(oneUpdate["update_at"].String(), "2024-05-30 20:30:00")

		// Replace
		dataReplace := g.Map{
			"id":        1,
			"name":      "name_100",
			"create_at": gtime.NewFromStr("2024-05-30 21:00:00"),
			"update_at": gtime.NewFromStr("2024-05-30 21:00:00"),
		}
		_, err = db.Model(table).Data(dataReplace).Replace()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		oneReplace, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneReplace["id"].Int(), 1)
		t.Assert(oneReplace["name"].String(), "name_100")
		t.Assert(oneReplace["delete_at"].String(), "")
		t.Assert(oneReplace["create_at"].String(), "2024-05-30 21:00:00")
		t.Assert(oneReplace["update_at"].String(), "2024-05-30 21:00:00")

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Insert with delete_at
		dataInsertDelete := g.Map{
			"id":        2,
			"name":      "name_2",
			"create_at": gtime.NewFromStr("2024-05-30 20:00:00"),
			"update_at": gtime.NewFromStr("2024-05-30 20:00:00"),
			"delete_at": gtime.NewFromStr("2024-05-30 20:00:00"),
		}
		_, err = db.Model(table).Data(dataInsertDelete).Insert()
		t.AssertNil(err)

		// Delete Select
		oneDelete, err := db.Model(table).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(len(oneDelete), 0)
		oneDeleteUnscoped, err := db.Model(table).Unscoped().WherePri(2).One()
		t.AssertNil(err)
		t.Assert(oneDeleteUnscoped["id"].Int(), 2)
		t.Assert(oneDeleteUnscoped["name"].String(), "name_2")
		t.Assert(oneDeleteUnscoped["delete_at"].String(), "2024-05-30 20:00:00")
		t.Assert(oneDeleteUnscoped["create_at"].String(), "2024-05-30 20:00:00")
		t.Assert(oneDeleteUnscoped["update_at"].String(), "2024-05-30 20:00:00")
	})
}
