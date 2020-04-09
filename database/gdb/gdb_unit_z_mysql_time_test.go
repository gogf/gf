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

func Test_CreateUpdateDeleteTime(t *testing.T) {
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

		time.Sleep(2 * time.Second)

		// Update
		dataUpdate := g.Map{
			"id":   1,
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
