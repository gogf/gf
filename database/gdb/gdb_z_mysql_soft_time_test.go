// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/test/gtest"
)

func Test_SoftCreateUpdateDeleteTime(t *testing.T) {
	table := "time_test_table"
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at datetime DEFAULT NULL,
  update_at datetime DEFAULT NULL,
  delete_at datetime DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
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
		r, err := db.Table(table).Data(dataInsert).Insert()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		oneInsert, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["name"].String(), "name_1")
		t.Assert(oneInsert["delete_at"].String(), "")
		t.AssertGE(oneInsert["create_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		t.AssertGE(oneInsert["update_at"].GTime().Timestamp(), gtime.Timestamp())

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Save
		dataSave := g.Map{
			"id":   1,
			"name": "name_10",
		}
		r, err = db.Table(table).Data(dataSave).Save()
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		oneSave, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["delete_at"].String(), "")
		t.Assert(oneSave["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertNE(oneSave["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())
		t.AssertGE(oneSave["update_at"].GTime().Timestamp(), gtime.Now().Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := g.Map{
			"name": "name_1000",
		}
		r, err = db.Table(table).Data(dataUpdate).WherePri(1).Update()
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		oneUpdate, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["delete_at"].String(), "")
		t.Assert(oneUpdate["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertGE(oneUpdate["update_at"].GTime().Timestamp(), gtime.Now().Timestamp()-2)

		// Replace
		dataReplace := g.Map{
			"id":   1,
			"name": "name_100",
		}
		r, err = db.Table(table).Data(dataReplace).Replace()
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		oneReplace, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(oneReplace["id"].Int(), 1)
		t.Assert(oneReplace["name"].String(), "name_100")
		t.Assert(oneReplace["delete_at"].String(), "")
		t.AssertGE(oneReplace["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertGE(oneReplace["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Delete
		r, err = db.Table(table).Delete("id", 1)
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
		// Delete Select
		one4, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(len(one4), 0)
		one5, err := db.Table(table).Unscoped().FindOne(1)
		t.Assert(err, nil)
		t.Assert(one5["id"].Int(), 1)
		t.AssertGE(one5["delete_at"].GTime().Timestamp(), gtime.Now().Timestamp()-2)
		// Delete Count
		i, err := db.Table(table).FindCount()
		t.Assert(err, nil)
		t.Assert(i, 0)
		i, err = db.Table(table).Unscoped().FindCount()
		t.Assert(err, nil)
		t.Assert(i, 1)

		// Delete Unscoped
		r, err = db.Table(table).Unscoped().Delete("id", 1)
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
		one6, err := db.Table(table).Unscoped().FindOne(1)
		t.Assert(err, nil)
		t.Assert(len(one6), 0)
		i, err = db.Table(table).Unscoped().FindCount()
		t.Assert(err, nil)
		t.Assert(i, 0)
	})
}

func Test_SoftUpdateTime(t *testing.T) {
<<<<<<< HEAD
	table := "time_test_table"
=======
	table := "time_test_table_" + gtime.TimestampNanoStr()
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  num       int(11) DEFAULT NULL,
  create_at datetime DEFAULT NULL,
  update_at datetime DEFAULT NULL,
  delete_at datetime DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
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
		r, err := db.Table(table).Data(dataInsert).Insert()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		oneInsert, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["num"].Int(), 10)

		// Update.
		r, err = db.Table(table).Data("num=num+1").Where("id=?", 1).Update()
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_SoftDelete(t *testing.T) {
	table := "time_test_table"
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at datetime DEFAULT NULL,
  update_at datetime DEFAULT NULL,
  delete_at datetime DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)
	// db.SetDebug(true)
	gtest.C(t, func(t *gtest.T) {
		for i := 1; i <= 10; i++ {
			data := g.Map{
				"id":   i,
				"name": fmt.Sprintf("name_%d", i),
			}
			r, err := db.Table(table).Data(data).Insert()
			t.Assert(err, nil)
			n, _ := r.RowsAffected()
			t.Assert(n, 1)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.AssertNE(one["create_at"].String(), "")
		t.AssertNE(one["update_at"].String(), "")
		t.Assert(one["delete_at"].String(), "")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).FindOne(10)
		t.Assert(err, nil)
		t.AssertNE(one["create_at"].String(), "")
		t.AssertNE(one["update_at"].String(), "")
		t.Assert(one["delete_at"].String(), "")
	})
	gtest.C(t, func(t *gtest.T) {
		ids := g.SliceInt{1, 3, 5}
		r, err := db.Table(table).Where("id", ids).Delete()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 3)

		count, err := db.Table(table).FindCount(ids)
		t.Assert(err, nil)
		t.Assert(count, 0)

		all, err := db.Table(table).Unscoped().FindAll(ids)
		t.Assert(err, nil)
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

func Test_SoftDelete_Join(t *testing.T) {
	table1 := "time_test_table1"
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at datetime DEFAULT NULL,
  update_at datetime DEFAULT NULL,
  delete_at datetime DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, table1)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table1)

	table2 := "time_test_table2"
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  createat datetime DEFAULT NULL,
  updateat datetime DEFAULT NULL,
  deleteat datetime DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, table2)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		//db.SetDebug(true)
		dataInsert1 := g.Map{
			"id":   1,
			"name": "name_1",
		}
		r, err := db.Table(table1).Data(dataInsert1).Insert()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		dataInsert2 := g.Map{
			"id":   1,
			"name": "name_2",
		}
		r, err = db.Table(table2).Data(dataInsert2).Insert()
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Table(table1, "t1").LeftJoin(table2, "t2", "t2.id=t1.id").Fields("t1.name").FindOne()
		t.Assert(err, nil)
		t.Assert(one["name"], "name_1")

		// Soft deleting.
		r, err = db.Table(table1).Delete()
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		one, err = db.Table(table1, "t1").LeftJoin(table2, "t2", "t2.id=t1.id").Fields("t1.name").FindOne()
		t.Assert(err, nil)
		t.Assert(one.IsEmpty(), true)

		one, err = db.Table(table2, "t2").LeftJoin(table1, "t1", "t2.id=t1.id").Fields("t2.name").FindOne()
		t.Assert(err, nil)
		t.Assert(one.IsEmpty(), true)
	})
}

func Test_CreateUpdateTime_Struct(t *testing.T) {
	table := "time_test_table"
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at datetime DEFAULT NULL,
  update_at datetime DEFAULT NULL,
  delete_at datetime DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
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
		r, err := db.Table(table).Data(dataInsert).Insert()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		oneInsert, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["name"].String(), "name_1")
		t.Assert(oneInsert["delete_at"].String(), "")
		t.AssertGE(oneInsert["create_at"].GTime().Timestamp(), gtime.Timestamp()-2)
		t.AssertGE(oneInsert["update_at"].GTime().Timestamp(), gtime.Timestamp())

		time.Sleep(2 * time.Second)

		// Save
		dataSave := &Entity{
			Id:       1,
			Name:     "name_10",
			CreateAt: nil,
			UpdateAt: nil,
			DeleteAt: nil,
		}
		r, err = db.Table(table).Data(dataSave).Save()
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		oneSave, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["delete_at"].String(), "")
		t.Assert(oneSave["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertNE(oneSave["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())
		t.AssertGE(oneSave["update_at"].GTime().Timestamp(), gtime.Now().Timestamp()-2)

		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := &Entity{
			Id:       1,
			Name:     "name_1000",
			CreateAt: nil,
			UpdateAt: nil,
			DeleteAt: nil,
		}
		r, err = db.Table(table).Data(dataUpdate).WherePri(1).Update()
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		oneUpdate, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["delete_at"].String(), "")
		t.Assert(oneUpdate["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertGE(oneUpdate["update_at"].GTime().Timestamp(), gtime.Now().Timestamp()-2)

		// Replace
		dataReplace := &Entity{
			Id:       1,
			Name:     "name_100",
			CreateAt: nil,
			UpdateAt: nil,
			DeleteAt: nil,
		}
		r, err = db.Table(table).Data(dataReplace).Replace()
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		oneReplace, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(oneReplace["id"].Int(), 1)
		t.Assert(oneReplace["name"].String(), "name_100")
		t.Assert(oneReplace["delete_at"].String(), "")
		t.AssertGE(oneReplace["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertGE(oneReplace["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())

		time.Sleep(2 * time.Second)

		// Delete
		r, err = db.Table(table).Delete("id", 1)
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
		// Delete Select
		one4, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(len(one4), 0)
		one5, err := db.Table(table).Unscoped().FindOne(1)
		t.Assert(err, nil)
		t.Assert(one5["id"].Int(), 1)
		t.AssertGE(one5["delete_at"].GTime().Timestamp(), gtime.Now().Timestamp()-2)
		// Delete Count
		i, err := db.Table(table).FindCount()
		t.Assert(err, nil)
		t.Assert(i, 0)
		i, err = db.Table(table).Unscoped().FindCount()
		t.Assert(err, nil)
		t.Assert(i, 1)

		// Delete Unscoped
		r, err = db.Table(table).Unscoped().Delete("id", 1)
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
		one6, err := db.Table(table).Unscoped().FindOne(1)
		t.Assert(err, nil)
		t.Assert(len(one6), 0)
		i, err = db.Table(table).Unscoped().FindCount()
		t.Assert(err, nil)
		t.Assert(i, 0)
	})
}
