// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Model_Save(t *testing.T) {
	table := createTableWithIdentity()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id          int
			AccountName string
			AttrIndex   int
		}
		var (
			user   User
			count  int
			result sql.Result
			err    error
		)

		// First insert: let IDENTITY auto-generate ID - use Insert() instead of Save()
		// because Save() requires a primary key in the data for conflict detection
		result, err = db.Model(table).Data(g.Map{
			"accountName": "ac1",
			"attrIndex":   100,
		}).Insert()

		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		err = db.Model(table).Scan(&user)
		t.AssertNil(err)
		t.AssertGT(user.Id, 0) // ID should be auto-generated
		t.Assert(user.AccountName, "ac1")
		t.Assert(user.AttrIndex, 100)

		// Second save: update the existing record using the generated ID
		_, err = db.Model(table).Data(g.Map{
			"id":          user.Id,
			"accountName": "ac2",
			"attrIndex":   200,
		}).OnConflict("id").Save()
		t.AssertNil(err)

		err = db.Model(table).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.AccountName, "ac2")
		t.Assert(user.AttrIndex, 200)

		_, err = db.Model(table).Data(g.Map{
			"id":          user.Id,
			"accountName": "ac2",
			"attrIndex":   2000,
		}).Save()
		t.AssertNil(err)

		err = db.Model(table).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.AccountName, "ac2")
		t.Assert(user.AttrIndex, 2000)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

func Test_Model_Insert(t *testing.T) {
	// g.Model.insert not lost default not null column
	table := "A_tables"
	createInitTable(table)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		i := 200
		data := User{
			ID:          int64(i),
			AccountName: fmt.Sprintf(`A%dtwo`, i),
			PwdReset:    0,
			AttrIndex:   99,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		}
		result, err := db.Model(table).Insert(&data)
		gtest.AssertNil(err)
		n, err := result.RowsAffected()
		gtest.AssertNil(err)
		gtest.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		i := 201
		data := User{
			ID:          int64(i),
			AccountName: fmt.Sprintf(`A%dtwoONE`, i),
			PwdReset:    1,
			CreatedTime: time.Now(),
			AttrIndex:   98,
			UpdatedTime: time.Now(),
		}
		result, err := db.Model(table).Data(&data).Insert()
		gtest.AssertNil(err)
		n, err := result.RowsAffected()
		gtest.AssertNil(err)
		gtest.Assert(n, 1)
	})
}

func Test_Model_InsertIgnore(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// db.SetDebug(true)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":           1,
			"account_name": fmt.Sprintf(`name_%d`, 777),
			"pwd_reset":    0,
			"attr_index":   777,
			"created_time": gtime.Now(),
		}
		_, err := db.Model(table).Data(data).InsertIgnore()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["ACCOUNT_NAME"].String(), "name_1")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			// "id":           1,
			"account_name": fmt.Sprintf(`name_%d`, 777),
			"pwd_reset":    0,
			"attr_index":   777,
			"created_time": gtime.Now(),
		}
		_, err := db.Model(table).Data(data).InsertIgnore()
		t.AssertNE(err, nil)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

func Test_Model_InsertAndGetId(t *testing.T) {
	table := createTableWithIdentity()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			// "id":          1,
			"account_name": fmt.Sprintf(`name_%d`, 1),
			"pwd_reset":    0,
			"attr_index":   1,
			"created_time": gtime.Now(),
		}
		lastId, err := db.Model(table).Data(data).InsertAndGetId()
		t.AssertNil(err)
		t.AssertGT(lastId, 0)
	})

}
