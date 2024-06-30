// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// CreateAt/UpdateAt/DeleteAt.
func Test_SoftTime_CreateUpdateDelete1(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at datetime(6) DEFAULT NULL,
  update_at datetime(6) DEFAULT NULL,
  delete_at datetime(6) DEFAULT NULL,
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
		r, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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
		r, err = db.Model(table).Data(dataSave).Save()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		oneSave, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["delete_at"].String(), "")
		t.Assert(oneSave["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertNE(oneSave["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())
		t.AssertGE(oneSave["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := g.Map{
			"name": "name_1000",
		}
		r, err = db.Model(table).Data(dataUpdate).WherePri(1).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["delete_at"].String(), "")
		t.Assert(oneUpdate["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertGE(oneUpdate["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// Replace
		dataReplace := g.Map{
			"id":   1,
			"name": "name_100",
		}
		r, err = db.Model(table).Data(dataReplace).Replace()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

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
		r, err = db.Model(table).Delete("id", 1)
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
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
		r, err = db.Model(table).Unscoped().Delete("id", 1)
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
		one6, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one6), 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 0)
	})
}

// CreateAt/UpdateAt/DeleteAt.
func Test_SoftTime_CreateUpdateDelete2(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at datetime(0) DEFAULT NULL,
  update_at datetime(0) DEFAULT NULL,
  delete_at datetime(0) DEFAULT NULL,
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
		r, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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
		r, err = db.Model(table).Data(dataSave).Save()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		oneSave, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["delete_at"].String(), "")
		t.Assert(oneSave["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertNE(oneSave["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())
		t.AssertGE(oneSave["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := g.Map{
			"name": "name_1000",
		}
		r, err = db.Model(table).Data(dataUpdate).WherePri(1).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["delete_at"].String(), "")
		t.Assert(oneUpdate["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertGE(oneUpdate["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// Replace
		dataReplace := g.Map{
			"id":   1,
			"name": "name_100",
		}
		r, err = db.Model(table).Data(dataReplace).Replace()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

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
		r, err = db.Model(table).Delete("id", 1)
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
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
		r, err = db.Model(table).Unscoped().Delete("id", 1)
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
		one6, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one6), 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 0)
	})
}

// CreatedAt/UpdatedAt/DeletedAt.
func Test_SoftTime_CreatedUpdatedDeleted_Map(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  created_at datetime(6) DEFAULT NULL,
  updated_at datetime(6) DEFAULT NULL,
  deleted_at datetime(6) DEFAULT NULL,
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
		r, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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
		r, err = db.Model(table).Data(dataSave).Save()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		oneSave, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["deleted_at"].String(), "")
		t.Assert(oneSave["created_at"].GTime().Timestamp(), oneInsert["created_at"].GTime().Timestamp())
		t.AssertNE(oneSave["updated_at"].GTime().Timestamp(), oneInsert["updated_at"].GTime().Timestamp())
		t.AssertGE(oneSave["updated_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := g.Map{
			"name": "name_1000",
		}
		r, err = db.Model(table).Data(dataUpdate).WherePri(1).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["deleted_at"].String(), "")
		t.Assert(oneUpdate["created_at"].GTime().Timestamp(), oneInsert["created_at"].GTime().Timestamp())
		t.AssertGE(oneUpdate["updated_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// Replace
		dataReplace := g.Map{
			"id":   1,
			"name": "name_100",
		}
		r, err = db.Model(table).Data(dataReplace).Replace()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

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
		r, err = db.Model(table).Delete("id", 1)
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
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
		r, err = db.Model(table).Unscoped().Delete("id", 1)
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
		one6, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one6), 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 0)
	})
}

// CreatedAt/UpdatedAt/DeletedAt.
func Test_SoftTime_CreatedUpdatedDeleted_Struct(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  created_at datetime(6) DEFAULT NULL,
  updated_at datetime(6) DEFAULT NULL,
  deleted_at datetime(6) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
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
		r, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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
		r, err = db.Model(table).Data(dataSave).OmitEmpty().Save()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		oneSave, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["deleted_at"].String(), "")
		t.Assert(oneSave["created_at"].GTime().Timestamp(), oneInsert["created_at"].GTime().Timestamp())
		t.AssertNE(oneSave["updated_at"].GTime().Timestamp(), oneInsert["updated_at"].GTime().Timestamp())
		t.AssertGE(oneSave["updated_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := User{
			Name: "name_1000",
		}
		r, err = db.Model(table).Data(dataUpdate).OmitEmpty().WherePri(1).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["deleted_at"].String(), "")
		t.Assert(oneUpdate["created_at"].GTime().Timestamp(), oneInsert["created_at"].GTime().Timestamp())
		t.AssertGE(oneUpdate["updated_at"].GTime().Timestamp(), gtime.Timestamp()-4)

		// Replace
		dataReplace := User{
			Id:   1,
			Name: "name_100",
		}
		r, err = db.Model(table).Data(dataReplace).OmitEmpty().Replace()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

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
		r, err = db.Model(table).Delete("id", 1)
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
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
		r, err = db.Model(table).Unscoped().Delete("id", 1)
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
		one6, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one6), 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 0)
	})
}

func Test_SoftUpdateTime(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  num       int(11) DEFAULT NULL,
  create_at datetime(6) DEFAULT NULL,
  update_at datetime(6) DEFAULT NULL,
  delete_at datetime(6) DEFAULT NULL,
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
		r, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		oneInsert, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["num"].Int(), 10)

		// Update.
		r, err = db.Model(table).Data("num=num+1").Where("id=?", 1).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_SoftUpdateTime_WithDO(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  num       int(11) DEFAULT NULL,
  created_at datetime(6) DEFAULT NULL,
  updated_at datetime(6) DEFAULT NULL,
  deleted_at datetime(6) DEFAULT NULL,
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
		r, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		oneInserted, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInserted["id"].Int(), 1)
		t.Assert(oneInserted["num"].Int(), 10)

		// Update.
		time.Sleep(2 * time.Second)
		type User struct {
			g.Meta    `orm:"do:true"`
			Id        interface{}
			Num       interface{}
			CreatedAt interface{}
			UpdatedAt interface{}
			DeletedAt interface{}
		}
		r, err = db.Model(table).Data(User{
			Num: 100,
		}).Where("id=?", 1).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		oneUpdated, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdated["num"].Int(), 100)
		t.Assert(oneUpdated["created_at"].String(), oneInserted["created_at"].String())
		t.AssertNE(oneUpdated["updated_at"].String(), oneInserted["updated_at"].String())
	})
}

func Test_SoftDelete(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at datetime(6) DEFAULT NULL,
  update_at datetime(6) DEFAULT NULL,
  delete_at datetime(6) DEFAULT NULL,
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
			r, err := db.Model(table).Data(data).Insert()
			t.AssertNil(err)
			n, _ := r.RowsAffected()
			t.Assert(n, 1)
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
		r, err := db.Model(table).Where("id", ids).Delete()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 3)

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

func Test_SoftDelete_Join(t *testing.T) {
	table1 := "time_test_table1"
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at datetime(6) DEFAULT NULL,
  update_at datetime(6) DEFAULT NULL,
  delete_at datetime(6) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, table1)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table1)

	table2 := "time_test_table2"
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  createat datetime(6) DEFAULT NULL,
  updateat datetime(6) DEFAULT NULL,
  deleteat datetime(6) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, table2)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		// db.SetDebug(true)
		dataInsert1 := g.Map{
			"id":   1,
			"name": "name_1",
		}
		r, err := db.Model(table1).Data(dataInsert1).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		dataInsert2 := g.Map{
			"id":   1,
			"name": "name_2",
		}
		r, err = db.Model(table2).Data(dataInsert2).Insert()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table1, "t1").LeftJoin(table2, "t2", "t2.id=t1.id").Fields("t1.name").One()
		t.AssertNil(err)
		t.Assert(one["name"], "name_1")

		// Soft deleting.
		r, err = db.Model(table1).Where(1).Delete()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		one, err = db.Model(table1, "t1").LeftJoin(table2, "t2", "t2.id=t1.id").Fields("t1.name").One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), true)

		one, err = db.Model(table2, "t2").LeftJoin(table1, "t1", "t2.id=t1.id").Fields("t2.name").One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), true)
	})
}

func Test_SoftDelete_WhereAndOr(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at datetime(6) DEFAULT NULL,
  update_at datetime(6) DEFAULT NULL,
  delete_at datetime(6) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)
	// db.SetDebug(true)
	// Add datas.
	gtest.C(t, func(t *gtest.T) {
		for i := 1; i <= 10; i++ {
			data := g.Map{
				"id":   i,
				"name": fmt.Sprintf("name_%d", i),
			}
			r, err := db.Model(table).Data(data).Insert()
			t.AssertNil(err)
			n, _ := r.RowsAffected()
			t.Assert(n, 1)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		ids := g.SliceInt{1, 3, 5}
		r, err := db.Model(table).Where("id", ids).Delete()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 3)

		count, err := db.Model(table).Where("id", 1).WhereOr("id", 3).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

func Test_CreateUpdateTime_Struct(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at datetime(6) DEFAULT NULL,
  update_at datetime(6) DEFAULT NULL,
  delete_at datetime(6) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	// db.SetDebug(true)
	// defer db.SetDebug(false)

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
		r, err := db.Model(table).Data(dataInsert).OmitEmpty().Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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
		r, err = db.Model(table).Data(dataSave).OmitEmpty().Save()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		oneSave, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["delete_at"].String(), "")
		t.Assert(oneSave["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
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
		r, err = db.Model(table).Data(dataUpdate).WherePri(1).OmitEmpty().Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["delete_at"].String(), "")
		t.Assert(oneUpdate["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertGE(oneUpdate["update_at"].GTime().Timestamp(), gtime.Timestamp()-2)

		// Replace
		dataReplace := &Entity{
			Id:       1,
			Name:     "name_100",
			CreateAt: nil,
			UpdateAt: nil,
			DeleteAt: nil,
		}
		r, err = db.Model(table).Data(dataReplace).OmitEmpty().Replace()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		oneReplace, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneReplace["id"].Int(), 1)
		t.Assert(oneReplace["name"].String(), "name_100")
		t.Assert(oneReplace["delete_at"].String(), "")
		t.AssertGE(oneReplace["create_at"].GTime().Timestamp(), oneInsert["create_at"].GTime().Timestamp())
		t.AssertGE(oneReplace["update_at"].GTime().Timestamp(), oneInsert["update_at"].GTime().Timestamp())

		time.Sleep(2 * time.Second)

		// Delete
		r, err = db.Model(table).Delete("id", 1)
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
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
		r, err = db.Model(table).Unscoped().Delete("id", 1)
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
		one6, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(len(one6), 0)
		i, err = db.Model(table).Unscoped().Count()
		t.AssertNil(err)
		t.Assert(i, 0)
	})
}

func Test_SoftTime_CreateUpdateDelete_UnixTimestamp(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at int(11) DEFAULT NULL,
  update_at int(11) DEFAULT NULL,
  delete_at int(11) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
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
		r, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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
		r, err := db.Model(table).Data(dataInsert).WherePri(1).Update()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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
		r, err = db.Model(table).Data("name='name_111'").WherePri(1).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		one, err = db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_111")
		t.Assert(one["create_at"].Int64(), lastCreateTime)
		t.AssertGT(one["update_at"].Int64(), lastUpdateTime)
		t.Assert(one["delete_at"].Int64(), 0)
	})

	// delete
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).WherePri(1).Delete()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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

func Test_SoftTime_CreateUpdateDelete_Bool_Deleted(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at int(11) DEFAULT NULL,
  update_at int(11) DEFAULT NULL,
  delete_at bit(1) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	//db.SetDebug(true)
	// insert
	gtest.C(t, func(t *gtest.T) {
		dataInsert := g.Map{
			"id":   1,
			"name": "name_1",
		}
		r, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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
		r, err := db.Model(table).WherePri(1).Delete()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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

func Test_SoftTime_CreateUpdateDelete_Option_SoftTimeTypeTimestampMilli(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at bigint(19) unsigned DEFAULT NULL,
  update_at bigint(19) unsigned DEFAULT NULL,
  delete_at bit(1) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
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
		r, err := db.Model(table).SoftTime(softTimeOption).Data(dataInsert).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).SoftTime(softTimeOption).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_1")
		t.Assert(len(one["create_at"].String()), 13)
		t.Assert(len(one["update_at"].String()), 13)
		t.Assert(one["delete_at"].Int64(), 0)
	})

	// delete
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).SoftTime(softTimeOption).WherePri(1).Delete()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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

func Test_SoftTime_CreateUpdateDelete_Option_SoftTimeTypeTimestampNano(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at bigint(19) unsigned DEFAULT NULL,
  update_at bigint(19) unsigned DEFAULT NULL,
  delete_at bit(1) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
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
		r, err := db.Model(table).SoftTime(softTimeOption).Data(dataInsert).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).SoftTime(softTimeOption).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "name_1")
		t.Assert(len(one["create_at"].String()), 19)
		t.Assert(len(one["update_at"].String()), 19)
		t.Assert(one["delete_at"].Int64(), 0)
	})

	// delete
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).SoftTime(softTimeOption).WherePri(1).Delete()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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

func Test_SoftTime_CreateUpdateDelete_Specified(t *testing.T) {
	table := "soft_time_test_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at datetime(0) DEFAULT NULL,
  update_at datetime(0) DEFAULT NULL,
  delete_at datetime(0) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
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
		r, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

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
		r, err = db.Model(table).Data(dataSave).Save()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		oneSave, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneSave["id"].Int(), 1)
		t.Assert(oneSave["name"].String(), "name_10")
		t.Assert(oneSave["delete_at"].String(), "")
		t.Assert(oneSave["create_at"].String(), "2024-05-30 20:00:00")
		t.Assert(oneSave["update_at"].String(), "2024-05-30 20:15:00")

		// For time asserting purpose.
		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := g.Map{
			"name":      "name_1000",
			"update_at": gtime.NewFromStr("2024-05-30 20:30:00"),
		}
		r, err = db.Model(table).Data(dataUpdate).WherePri(1).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["delete_at"].String(), "")
		t.Assert(oneUpdate["create_at"].String(), "2024-05-30 20:00:00")
		t.Assert(oneUpdate["update_at"].String(), "2024-05-30 20:30:00")

		// Replace
		dataReplace := g.Map{
			"id":        1,
			"name":      "name_100",
			"create_at": gtime.NewFromStr("2024-05-30 21:00:00"),
			"update_at": gtime.NewFromStr("2024-05-30 21:00:00"),
		}
		r, err = db.Model(table).Data(dataReplace).Replace()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

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
		r, err = db.Model(table).Data(dataInsertDelete).Insert()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

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
