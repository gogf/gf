// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

// Fix issue: https://github.com/gogf/gf/issues/819
func Test_Model_Insert_WithStructAndSliceAttribute(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type Password struct {
			Salt string `json:"salt"`
			Pass string `json:"pass"`
		}
		data := g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    &Password{"123", "456"},
			"nickname":    []string{"A", "B", "C"},
			"create_time": gtime.Now().String(),
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).One("id", 1)
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data["passport"])
		t.Assert(one["CREATE_TIME"], data["create_time"])
		t.Assert(one["NICKNAME"], gjson.New(data["nickname"]).MustToJson())
	})
}

func Test_Model_UpdateAndGetAffected(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		n, err := db.Model(table).Data("nickname", "T100").
			Where("id>?", TableSize-2).
			UpdateAndGetAffected()
		t.AssertNil(err)
		t.Assert(n, 2)

		v1, err := db.Model(table).Fields("nickname").Where("id", 10).Value()
		t.AssertNil(err)
		t.Assert(v1.String(), "T100")

		v2, err := db.Model(table).Fields("nickname").Where("id", 8).Value()
		t.AssertNil(err)
		t.Assert(v2.String(), "name_8")
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data("nickname", "T200").
			Where("1=1").Order("id desc").Limit(2).
			UpdateAndGetAffected()
		t.AssertNE(err, nil)
		t.AssertIN("Incorrect syntax near the keyword 'ORDER'", fmt.Sprint(err))
	})
}

func Test_Model_Fields(t *testing.T) {
	tableName1 := createInitTable()
	defer dropTable(tableName1)

	tableName2 := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        ID   int NOT NULL,
	        NAME varchar(45) NULL,
	        AGE  int NULL,
	        PRIMARY KEY (ID)
	    )`, tableName2,
	)); err != nil {
		gtest.AssertNil(err)
	}
	defer dropTable(tableName2)

	r, err := db.Insert(ctx, tableName2, g.Map{
		"id":   1,
		"name": "table2_1",
		"age":  18,
	})
	gtest.AssertNil(err)
	n, _ := r.RowsAffected()
	gtest.Assert(n, 1)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(tableName1).As("u").Fields("u.passport,u.id").Where("u.id<2").All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(len(all[0]), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(tableName1).As("u1").
			LeftJoin(tableName1, "u2", "u2.id=u1.id").
			Fields("u1.passport,u1.id,u2.id AS u2id").
			Where("u1.id<2").
			All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(len(all[0]), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(tableName1).As("u1").
			LeftJoin(tableName2, "u2", "u2.id=u1.id").
			Fields("u1.passport,u1.id,u2.name,u2.age").
			Where("u1.id<2").
			All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(len(all[0]), 4)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["age"], 18)
		t.Assert(all[0]["name"], "table2_1")
		t.Assert(all[0]["passport"], "user_1")
	})
}

func Test_Model_Value_WithCache(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Where("id", 1).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.MapStrAny{
			"id":       1,
			"passport": fmt.Sprintf(`passport_%d`, 1),
			"password": fmt.Sprintf(`password_%d`, 1),
			"nickname": fmt.Sprintf(`nickname_%d`, 1),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Where("id", 1).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Value("id")
		t.AssertNil(err)
		t.Assert(value.Int(), 1)
	})
}

func Test_Model_Count_WithCache(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 1).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.MapStrAny{
			"id":       1,
			"passport": fmt.Sprintf(`passport_%d`, 1),
			"password": fmt.Sprintf(`password_%d`, 1),
			"nickname": fmt.Sprintf(`nickname_%d`, 1),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 1).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
}

func Test_Model_Count_All_WithCache(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.MapStrAny{
			"id":       1,
			"passport": fmt.Sprintf(`passport_%d`, 1),
			"password": fmt.Sprintf(`password_%d`, 1),
			"nickname": fmt.Sprintf(`nickname_%d`, 1),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.MapStrAny{
			"id":       2,
			"passport": fmt.Sprintf(`passport_%d`, 2),
			"password": fmt.Sprintf(`password_%d`, 2),
			"nickname": fmt.Sprintf(`nickname_%d`, 2),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
}

func Test_Model_CountColumn_WithCache(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 1).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).CountColumn("id")
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.MapStrAny{
			"id":       1,
			"passport": fmt.Sprintf(`passport_%d`, 1),
			"password": fmt.Sprintf(`password_%d`, 1),
			"nickname": fmt.Sprintf(`nickname_%d`, 1),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 1).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).CountColumn("id")
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
}

func Test_Model_StructsWithOrmTag(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	dbErr.SetDebug(true)
	defer dbErr.SetDebug(false)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid      int `orm:"id"`
			Passport string
			Password string     `orm:"password"`
			Name     string     `orm:"nick_name"`
			Time     gtime.Time `orm:"create_time"`
		}
		var (
			users  []User
			buffer = bytes.NewBuffer(nil)
		)
		dbErr.GetLogger().(*glog.Logger).SetWriter(buffer)
		defer dbErr.GetLogger().(*glog.Logger).SetWriter(os.Stdout)
		dbErr.Model(table).Order("id asc").Scan(&users)
		t.Assert(
			gstr.Contains(
				buffer.String(),
				fmt.Sprintf(`SELECT "id","Passport","password","nick_name","create_time" FROM "%s"`, table),
			),
			true,
		)
	})

	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Passport string
			Password string
		}
		type B struct {
			A
			NickName string
		}
		one, err := db.Model(table).Fields(&B{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 3)
		t.Assert(one["NICKNAME"], "name_2")
		t.Assert(one["PASSPORT"], "user_2")
		t.Assert(one["PASSWORD"], "pass_2")
	})
}

func Test_Model_GroupBy(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Group("id").All()
		t.AssertNE(err, nil)
		t.AssertIN("is invalid in the select list", fmt.Sprint(err))
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields("id, nickname").Group("id, nickname").Order("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["NICKNAME"].String(), "name_1")
	})
}

func Test_Model_Page(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(3, 3).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"], 7)
		t.Assert(result[1]["ID"], 8)
	})
	gtest.C(t, func(t *gtest.T) {
		model := db.Model(table).Safe().Order("id")
		all, err := model.Page(3, 3).All()
		count, err := model.Count()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["ID"], "7")
		t.Assert(count, int64(TableSize))
	})
}

func Test_Model_Option_List(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Model(table).Fields("id, password").Data(g.List{
			g.Map{
				"id":       1,
				"passport": "1",
				"password": "1",
				"nickname": "1",
			},
			g.Map{
				"id":       2,
				"passport": "2",
				"password": "2",
				"nickname": "2",
			},
		}).Save()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)
		list, err := db.Model(table).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(list), 2)
		t.Assert(list[0]["ID"].String(), "1")
		t.Assert(list[0]["NICKNAME"].String(), "")
		t.Assert(list[0]["PASSPORT"].String(), "")
		t.Assert(list[0]["PASSWORD"].String(), "1")

		t.Assert(list[1]["ID"].String(), "2")
		t.Assert(list[1]["NICKNAME"].String(), "")
		t.Assert(list[1]["PASSPORT"].String(), "")
		t.Assert(list[1]["PASSWORD"].String(), "2")
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Model(table).OmitEmpty().Fields("id, password").Data(g.List{
			g.Map{
				"id":       1,
				"passport": "1",
				"password": 0,
				"nickname": "1",
			},
			g.Map{
				"id":       2,
				"passport": "2",
				"password": "2",
				"nickname": "2",
			},
		}).Save()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)
		list, err := db.Model(table).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(list), 2)
		t.Assert(list[0]["ID"].String(), "1")
		t.Assert(list[0]["NICKNAME"].String(), "")
		t.Assert(list[0]["PASSPORT"].String(), "")
		t.Assert(list[0]["PASSWORD"].String(), "0")

		t.Assert(list[1]["ID"].String(), "2")
		t.Assert(list[1]["NICKNAME"].String(), "")
		t.Assert(list[1]["PASSPORT"].String(), "")
		t.Assert(list[1]["PASSWORD"].String(), "2")
	})
}

func Test_Model_OmitEmpty(t *testing.T) {
	table := fmt.Sprintf(`table_%s`, gtime.TimestampNanoStr())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE %s (
        ID   int NOT NULL,
        NAME varchar(45) NOT NULL,
        PRIMARY KEY (ID)
    )`, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitEmpty().Data(g.Map{
			"id":   1,
			"name": "",
		}).Save()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"id":   1,
			"name": "",
		}).Save()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitEmptyWhere().Data(g.Map{
			"id":   1,
			"name": "",
		}).Save()
		t.AssertNil(err)
	})
}

func Test_Model_OmitNil(t *testing.T) {
	table := fmt.Sprintf(`table_%s`, gtime.TimestampNanoStr())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE %s (
        ID   int NOT NULL,
        NAME varchar(45) NOT NULL,
        PRIMARY KEY (ID)
    )`, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitNil().Data(g.Map{
			"id":   1,
			"name": nil,
		}).Save()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitNil().Data(g.Map{
			"id":   1,
			"name": "",
		}).Save()
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitNilWhere().Data(g.Map{
			"id":   1,
			"name": "",
		}).Save()
		t.AssertNil(err)
	})
}

func Test_Model_FieldsEx_WithReservedWords(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := "fieldsex_test_table"
		dropTable(table)
		if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		  ID int NOT NULL,
		  [KEY] varchar(45) DEFAULT NULL,
		  CATEGORY_ID int NOT NULL,
		  USER_ID int NOT NULL,
		  TITLE varchar(255) NOT NULL,
		  CONTENT varchar(max) NOT NULL,
		  SORT int DEFAULT 0,
		  BRIEF varchar(255) DEFAULT NULL,
		  THUMB varchar(255) DEFAULT NULL,
		  TAGS varchar(900) DEFAULT NULL,
		  REFERER varchar(255) DEFAULT NULL,
		  STATUS smallint DEFAULT 0,
		  VIEW_COUNT int DEFAULT 0,
		  ZAN_COUNT int DEFAULT NULL,
		  CAI_COUNT int DEFAULT NULL,
		  CREATED_AT datetime DEFAULT NULL,
		  UPDATED_AT datetime DEFAULT NULL,
		  PRIMARY KEY (ID)
		)`, table)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table)
		_, err := db.Model(table).FieldsEx("content").One()
		t.AssertNil(err)
	})
}

func Test_Model_Prefix(t *testing.T) {
	nodePrefix := *db.GetConfig()
	nodePrefix.Prefix = TableNamePrefix1
	dbPrefix, err := gdb.New(nodePrefix)
	gtest.AssertNil(err)
	defer dbPrefix.Close(ctx)

	db := dbPrefix
	table := fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
	createInitTableWithDb(db, TableNamePrefix1+table)
	defer dropTable(TableNamePrefix1 + table)
	// Select.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Where("id in (?)", g.Slice{1, 2}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["ID"], "1")
		t.Assert(r[1]["ID"], "2")
	})
	// Select with alias.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table+" as u").Where("u.id in (?)", g.Slice{1, 2}).Order("u.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["ID"], "1")
		t.Assert(r[1]["ID"], "2")
	})
	// Select with alias to struct.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Passport string
			Password string
			NickName string
		}
		var users []User
		err := db.Model(table+" u").Where("u.id in (?)", g.Slice{1, 5}).Order("u.id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 5)
	})
	// Select with alias and join statement.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table+" as u1").LeftJoin(table+" as u2", "u2.id=u1.id").Where("u1.id in (?)", g.Slice{1, 2}).Order("u1.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["ID"], "1")
		t.Assert(r[1]["ID"], "2")
	})
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).As("u1").LeftJoin(table+" as u2", "u2.id=u1.id").Where("u1.id in (?)", g.Slice{1, 2}).Order("u1.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["ID"], "1")
		t.Assert(r[1]["ID"], "2")
	})
}

func Test_Model_Schema1(t *testing.T) {
	const (
		TestSchema1 = TestSchema
		TestSchema2 = "tempdb"
	)

	db := db.Schema(TestSchema1)
	table := fmt.Sprintf(`%s_%s`, TableName, gtime.TimestampNanoStr())
	createInitTableWithDb(db, table)
	db = db.Schema(TestSchema2)
	createInitTableWithDb(db, table)
	defer func() {
		db = db.Schema(TestSchema1)
		dropTableWithDb(db, table)
		db = db.Schema(TestSchema2)
		dropTableWithDb(db, table)
	}()
	// Method.
	gtest.C(t, func(t *gtest.T) {
		db = db.Schema(TestSchema1)
		r, err := db.Model(table).Update(g.Map{"nickname": "name_100"}, "id=1")
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		v, err := db.Model(table).Value("nickname", "id=1")
		t.AssertNil(err)
		t.Assert(v.String(), "name_100")

		db = db.Schema(TestSchema2)
		v, err = db.Model(table).Value("nickname", "id=1")
		t.AssertNil(err)
		t.Assert(v.String(), "name_1")
	})
	// Model.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Schema(TestSchema1).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_2")

		r, err := db.Model(table).Schema(TestSchema1).Update(g.Map{"nickname": "name_200"}, "id=2")
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		v, err = db.Model(table).Schema(TestSchema1).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_200")

		v, err = db.Model(table).Schema(TestSchema2).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_2")

		v, err = db.Model(table).Schema(TestSchema1).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_200")
	})
	// Model.
	gtest.C(t, func(t *gtest.T) {
		i := 1000
		_, err := db.Model(table).Schema(TestSchema1).Insert(g.Map{
			"id":               i,
			"passport":         fmt.Sprintf(`user_%d`, i),
			"password":         fmt.Sprintf(`pass_%d`, i),
			"nickname":         fmt.Sprintf(`name_%d`, i),
			"create_time":      gtime.NewFromStr("2018-10-24 10:00:00").String(),
			"none-exist-field": 1,
		})
		t.AssertNil(err)

		v, err := db.Model(table).Schema(TestSchema1).Value("nickname", "id=?", i)
		t.AssertNil(err)
		t.Assert(v.String(), "name_1000")

		v, err = db.Model(table).Schema(TestSchema2).Value("nickname", "id=?", i)
		t.AssertNil(err)
		t.Assert(v.String(), "")
	})
}

func Test_Model_Schema2(t *testing.T) {
	const (
		TestSchema1 = TestSchema
		TestSchema2 = "tempdb"
	)

	db := db.Schema(TestSchema1)
	table := fmt.Sprintf(`%s_%s`, TableName, gtime.TimestampNanoStr())
	createInitTableWithDb(db, table)
	db = db.Schema(TestSchema2)
	createInitTableWithDb(db, table)
	defer func() {
		db = db.Schema(TestSchema1)
		dropTableWithDb(db, table)
		db = db.Schema(TestSchema2)
		dropTableWithDb(db, table)
	}()
	// Schema.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Schema(TestSchema1).Model(table).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_2")

		r, err := db.Schema(TestSchema1).Model(table).Update(g.Map{"nickname": "name_200"}, "id=2")
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		v, err = db.Schema(TestSchema1).Model(table).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_200")

		v, err = db.Schema(TestSchema2).Model(table).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_2")

		v, err = db.Schema(TestSchema1).Model(table).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_200")
	})
	// Schema.
	gtest.C(t, func(t *gtest.T) {
		i := 1000
		_, err := db.Schema(TestSchema1).Model(table).Insert(g.Map{
			"id":               i,
			"passport":         fmt.Sprintf(`user_%d`, i),
			"password":         fmt.Sprintf(`pass_%d`, i),
			"nickname":         fmt.Sprintf(`name_%d`, i),
			"create_time":      gtime.NewFromStr("2018-10-24 10:00:00").String(),
			"none-exist-field": 1,
		})
		t.AssertNil(err)

		v, err := db.Schema(TestSchema1).Model(table).Value("nickname", "id=?", i)
		t.AssertNil(err)
		t.Assert(v.String(), "name_1000")

		v, err = db.Schema(TestSchema2).Model(table).Value("nickname", "id=?", i)
		t.AssertNil(err)
		t.Assert(v.String(), "")
	})
}

func Test_Model_Cache(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test1",
			Force:    false,
		}).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")

		r, err := db.Model(table).Data("passport", "user_100").WherePri(1).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test1",
			Force:    false,
		}).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")

		time.Sleep(time.Second * 2)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test1",
			Force:    false,
		}).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_100")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test2",
			Force:    false,
		}).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_2")

		r, err := db.Model(table).Data("passport", "user_200").Cache(gdb.CacheOption{
			Duration: -1,
			Name:     "test2",
			Force:    false,
		}).WherePri(2).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test2",
			Force:    false,
		}).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_200")
	})
	// transaction.
	gtest.C(t, func(t *gtest.T) {
		// make cache for id 3
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test3",
			Force:    false,
		}).WherePri(3).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_3")

		r, err := db.Model(table).Data("passport", "user_300").Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test3",
			Force:    false,
		}).WherePri(3).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		err = db.Transaction(context.TODO(), func(ctx context.Context, tx gdb.TX) error {
			one, err := tx.Model(table).Cache(gdb.CacheOption{
				Duration: time.Second,
				Name:     "test3",
				Force:    false,
			}).WherePri(3).One()
			t.AssertNil(err)
			t.Assert(one["PASSPORT"], "user_300")
			return nil
		})
		t.AssertNil(err)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test3",
			Force:    false,
		}).WherePri(3).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_3")
	})
	gtest.C(t, func(t *gtest.T) {
		// make cache for id 4
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test4",
			Force:    false,
		}).WherePri(4).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_4")

		r, err := db.Model(table).Data("passport", "user_400").Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test3",
			Force:    false,
		}).WherePri(4).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		err = db.Transaction(context.TODO(), func(ctx context.Context, tx gdb.TX) error {
			// Cache feature disabled.
			one, err := tx.Model(table).Cache(gdb.CacheOption{
				Duration: time.Second,
				Name:     "test4",
				Force:    false,
			}).WherePri(4).One()
			t.AssertNil(err)
			t.Assert(one["PASSPORT"], "user_400")
			// Update the cache.
			r, err := tx.Model(table).Data("passport", "user_4000").
				Cache(gdb.CacheOption{
					Duration: -1,
					Name:     "test4",
					Force:    false,
				}).WherePri(4).Update()
			t.AssertNil(err)
			n, err := r.RowsAffected()
			t.AssertNil(err)
			t.Assert(n, 1)
			return nil
		})
		t.AssertNil(err)
		// Read from db.
		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test4",
			Force:    false,
		}).WherePri(4).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_4000")
	})
}

func Test_Model_FieldsEx_AutoMapping(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// "id":          i,
	// "passport":    fmt.Sprintf(`user_%d`, i),
	// "password":    fmt.Sprintf(`pass_%d`, i),
	// "nickname":    fmt.Sprintf(`name_%d`, i),
	// "create_time": gtime.NewFromStr("2018-10-24 10:00:00").String(),

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).FieldsEx("created_at, updated_at, Passport, Password, NickName, CreateTime").Where("id", 2).Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).FieldsEx("created_at, updated_at, ID, Passport, Password, CreateTime").Where("id", 2).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "name_2")
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).FieldsEx(g.Map{
			"Passport":   1,
			"Password":   1,
			"CreateTime": 1,
		}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 4)
		t.Assert(one["ID"], 2)
		t.Assert(one["NICKNAME"], "name_2")
	})
	// Struct
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Passport   int
			Password   int
			CreateTime int
		}
		one, err := db.Model(table).FieldsEx(&T{
			Passport:   0,
			Password:   0,
			CreateTime: 0,
		}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 4)
		t.Assert(one["ID"], 2)
		t.Assert(one["NICKNAME"], "name_2")
	})
}

func Test_Model_Fields_Struct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type A struct {
		Passport string
		Password string
	}
	type B struct {
		A
		NickName string
	}
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields(A{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["PASSPORT"], "user_2")
		t.Assert(one["PASSWORD"], "pass_2")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields(&A{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["PASSPORT"], "user_2")
		t.Assert(one["PASSWORD"], "pass_2")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields(B{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 3)
		t.Assert(one["PASSPORT"], "user_2")
		t.Assert(one["PASSWORD"], "pass_2")
		t.Assert(one["NICKNAME"], "name_2")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields(&B{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 3)
		t.Assert(one["PASSPORT"], "user_2")
		t.Assert(one["PASSWORD"], "pass_2")
		t.Assert(one["NICKNAME"], "name_2")
	})
}

func Test_Model_Empty_Slice_Argument(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(`id`, g.Slice{}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(`id in(?)`, g.Slice{}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

// https://github.com/gogf/gf/issues/1012
func Test_TimeZoneInsert(t *testing.T) {
	tableName := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        ID         int NOT NULL,
	        PASSPORT   varchar(45) NULL,
	        PASSWORD   varchar(32) NULL,
	        NICKNAME   varchar(45) NULL,
	        CREATED_AT datetimeoffset NULL,
	        UPDATED_AT datetimeoffset NULL,
	        DELETED_AT datetimeoffset NULL,
	        PRIMARY KEY (ID)
	    )`, tableName,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(tableName)

	tokyoLoc, err := time.LoadLocation("Asia/Tokyo")
	gtest.AssertNil(err)

	CreateTime := "2020-11-22 12:23:45"
	UpdateTime := "2020-11-22 13:23:46"
	DeleteTime := "2020-11-22 14:23:47"
	type User struct {
		Id        int         `json:"id"`
		CreatedAt *gtime.Time `json:"created_at"`
		UpdatedAt gtime.Time  `json:"updated_at"`
		DeletedAt time.Time   `json:"deleted_at"`
	}
	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", CreateTime, tokyoLoc)
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", UpdateTime, tokyoLoc)
	t3, _ := time.ParseInLocation("2006-01-02 15:04:05", DeleteTime, tokyoLoc)
	u := &User{
		Id:        1,
		CreatedAt: gtime.New(t1.UTC()),
		UpdatedAt: *gtime.New(t2.UTC()),
		DeletedAt: t3.UTC(),
	}

	gtest.C(t, func(t *gtest.T) {
		_, err = db.Model(tableName).Unscoped().Insert(u)
		t.AssertNil(err)
		userEntity := &User{}
		err = db.Model(tableName).Where("id", 1).Unscoped().Scan(&userEntity)
		t.AssertNil(err)
		t.Assert(userEntity.CreatedAt.String(), "2020-11-22 03:23:45")
		t.Assert(userEntity.UpdatedAt.String(), "2020-11-22 04:23:46")
		t.Assert(gtime.NewFromTime(userEntity.DeletedAt).String(), "2020-11-22 05:23:47")
		t.Assert(userEntity.CreatedAt.Time.Equal(t1), true)
		t.Assert(userEntity.UpdatedAt.Time.Equal(t2), true)
		t.Assert(userEntity.DeletedAt.Equal(t3), true)
	})
}

func Test_Model_Fields_Map_Struct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields(g.Map{
			"ID":         1,
			"PASSPORT":   1,
			"NONE_EXIST": 1,
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result["ID"], 1)
		t.Assert(result["PASSPORT"], "user_1")
	})
	// struct
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			ID       int
			PASSPORT string
			XXX_TYPE int
		}
		a := A{}
		err := db.Model(table).Fields(a).Where("id", 1).Scan(&a)
		t.AssertNil(err)
		t.Assert(a.ID, 1)
		t.Assert(a.PASSPORT, "user_1")
		t.Assert(a.XXX_TYPE, 0)
	})
	// *struct
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			ID       int
			PASSPORT string
			XXX_TYPE int
		}
		var a *A
		err := db.Model(table).Fields(a).Where("id", 1).Scan(&a)
		t.AssertNil(err)
		t.Assert(a.ID, 1)
		t.Assert(a.PASSPORT, "user_1")
		t.Assert(a.XXX_TYPE, 0)
	})
	// **struct
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			ID       int
			PASSPORT string
			XXX_TYPE int
		}
		var a *A
		err := db.Model(table).Fields(&a).Where("id", 1).Scan(&a)
		t.AssertNil(err)
		t.Assert(a.ID, 1)
		t.Assert(a.PASSPORT, "user_1")
		t.Assert(a.XXX_TYPE, 0)
	})
}

func Test_Model_InsertAndGetId(t *testing.T) {
	table := fmt.Sprintf(`user_%d`, gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s (
		ID int IDENTITY(1,1) NOT NULL,
		PASSPORT varchar(45) NULL,
		PASSWORD varchar(32) NULL,
		NICKNAME varchar(45) NULL,
		PRIMARY KEY (ID)
	)`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Exec(fmt.Sprintf("SET IDENTITY_INSERT %s ON", table))
			t.AssertNil(err)
			id, err := tx.Model(table).Data(g.Map{
				"id":       1,
				"passport": "user_1",
				"password": "pass_1",
				"nickname": "name_1",
			}).InsertAndGetId()
			t.AssertNil(err)
			t.Assert(id, 1)
			_, err = tx.Exec(fmt.Sprintf("SET IDENTITY_INSERT %s OFF", table))
			return err
		})
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{
			"passport": "user_2",
			"password": "pass_2",
			"nickname": "name_2",
		}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 2)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":       3,
			"passport": "user_3",
			"password": "pass_3",
			"nickname": "name_3",
		}).InsertAndGetId()
		t.AssertNE(err, nil)
		t.AssertIN("IDENTITY_INSERT is set to OFF", fmt.Sprint(err))
	})
}

func Test_Model_Increment_Decrement(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 1).Increment("id", 100)
		t.AssertNil(err)
		rows, _ := result.RowsAffected()
		t.Assert(rows, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 101).Decrement("id", 10)
		t.AssertNil(err)
		rows, _ := result.RowsAffected()
		t.Assert(rows, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 91).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
}

func Test_Model_OnDuplicate(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// string type 1.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicate("passport,password").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data["passport"])
		t.Assert(one["PASSWORD"], data["password"])
		t.Assert(one["NICKNAME"], "name_1")
	})

	// string type 2.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicate("passport", "password").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data["passport"])
		t.Assert(one["PASSWORD"], data["password"])
		t.Assert(one["NICKNAME"], "name_1")
	})

	// slice.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicate(g.Slice{"passport", "password"}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data["passport"])
		t.Assert(one["PASSWORD"], data["password"])
		t.Assert(one["NICKNAME"], "name_1")
	})

	// map.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicate(g.Map{
			"passport": "nickname",
			"password": "nickname",
		}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data["nickname"])
		t.Assert(one["PASSWORD"], data["nickname"])
		t.Assert(one["NICKNAME"], "name_1")
	})

	// map+raw.
	gtest.C(t, func(t *gtest.T) {
		data := g.MapStrStr{
			"id":          "1",
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicate(g.Map{
			"passport": gdb.Raw("CONCAT(T2.PASSPORT, '1')"),
			"password": gdb.Raw("CONCAT(T2.PASSWORD, '2')"),
		}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data["passport"]+"1")
		t.Assert(one["PASSWORD"], data["password"]+"2")
		t.Assert(one["NICKNAME"], "name_1")
	})
}

func Test_Model_OnDuplicateWithCounter(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict("id").OnDuplicate(g.Map{
			"id": gdb.Counter{Field: "id", Value: 999999},
		}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.AssertNil(one)
	})
}

func Test_Model_OnDuplicateEx(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// string type 1.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicateEx("nickname,create_time").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data["passport"])
		t.Assert(one["PASSWORD"], data["password"])
		t.Assert(one["NICKNAME"], "name_1")
	})

	// string type 2.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicateEx("nickname", "create_time").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data["passport"])
		t.Assert(one["PASSWORD"], data["password"])
		t.Assert(one["NICKNAME"], "name_1")
	})

	// slice.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicateEx(g.Slice{"nickname", "create_time"}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data["passport"])
		t.Assert(one["PASSWORD"], data["password"])
		t.Assert(one["NICKNAME"], "name_1")
	})

	// map.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicateEx(g.Map{
			"nickname":    "nickname",
			"create_time": "nickname",
		}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data["passport"])
		t.Assert(one["PASSWORD"], data["password"])
		t.Assert(one["NICKNAME"], "name_1")
	})
}

// https://github.com/gogf/gf/issues/1387
func Test_Model_GTime_DefaultValue(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			Nickname   string
			CreateTime *gtime.Time
		}
		data := User{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
			Nickname: "name_1",
		}
		// Insert
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Select
		var (
			user *User
		)
		err = db.Model(table).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Passport, data.Passport)
		t.Assert(user.Password, data.Password)
		t.Assert(user.CreateTime, data.CreateTime)
		t.Assert(user.Nickname, data.Nickname)

		// Insert
		user.Id = 2
		_, err = db.Model(table).Data(user).Insert()
		t.AssertNil(err)
	})
}

// Using filter does not affect the outside value inside function.
func Test_Model_Insert_Filter(t *testing.T) {
	// map
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		data := g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 0)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "t1")
		t.Assert(one["NICKNAME"], "name_1")

		t.Assert(data["uid"], 1)
	})
	// slice
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		data := g.List{
			g.Map{
				"id":          1,
				"uid":         1,
				"passport":    "t1",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_1",
				"create_time": gtime.Now().String(),
			},
			g.Map{
				"id":          2,
				"uid":         2,
				"passport":    "t1",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_1",
				"create_time": gtime.Now().String(),
			},
		}

		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 0)
		n, _ = result.RowsAffected()
		t.Assert(n, 2)

		all, err := db.Model(table).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["ID"], 1)
		t.Assert(all[1]["ID"], 2)
		t.Assert(all[1]["PASSPORT"], "t1")

		t.Assert(data[0]["uid"], 1)
		t.Assert(data[1]["uid"], 2)
	})
}

func Test_Model_Embedded_Filter(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type Base struct {
			Id         int
			Uid        int
			CreateTime string
			NoneExist  string
		}
		type User struct {
			Base
			Passport string
			Password string
			Nickname string
		}
		result, err := db.Model(table).Data(User{
			Passport: "john-test",
			Password: "123456",
			Nickname: "John",
			Base: Base{
				Id:         100,
				Uid:        100,
				CreateTime: gtime.Now().String(),
			},
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		var user *User
		err = db.Model(table).Fields(user).Where("id=100").Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Passport, "john-test")
		t.Assert(user.Id, 100)
	})
}

func Test_Model_Fields_AutoFilterInJoinStatement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err    error
			suffix = gtime.TimestampNano()
			table1 = fmt.Sprintf("user_%d", suffix)
			table2 = fmt.Sprintf("score_%d", suffix)
			table3 = fmt.Sprintf("info_%d", suffix)
		)
		if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		   ID int NOT NULL,
		   NAME varchar(500) NOT NULL DEFAULT '',
		 PRIMARY KEY (ID)
		)`, table1,
		)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table1)
		_, err = db.Model(table1).Insert(g.Map{
			"id":   1,
			"name": "john",
		})
		t.AssertNil(err)

		if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			ID int NOT NULL,
			USER_ID int NOT NULL DEFAULT 0,
		    NUMBER varchar(500) NOT NULL DEFAULT '',
		 PRIMARY KEY (ID)
		)`, table2,
		)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table2)
		_, err = db.Model(table2).Insert(g.Map{
			"id":      1,
			"user_id": 1,
			"number":  "n",
		})
		t.AssertNil(err)

		if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			ID int NOT NULL,
			USER_ID int NOT NULL DEFAULT 0,
		    DESCRIPTION varchar(500) NOT NULL DEFAULT '',
		 PRIMARY KEY (ID)
		)`, table3,
		)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table3)
		_, err = db.Model(table3).Insert(g.Map{
			"id":          1,
			"user_id":     1,
			"description": "brief",
		})
		t.AssertNil(err)

		one, err := db.Model(table1).
			Where(table1+".id", 1).
			Fields(fmt.Sprintf("%s.number,%s.name", table2, table1)).
			LeftJoin(table2, fmt.Sprintf("%s.id=%s.user_id", table1, table2)).
			LeftJoin(table3, fmt.Sprintf("%s.id=%s.user_id", table3, table3)).
			Order(table1 + ".id asc").
			One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["name"].String(), "john")
		t.Assert(one["number"].String(), "n")

		one, err = db.Model(table1).
			LeftJoin(table2, fmt.Sprintf("%s.id=%s.user_id", table1, table2)).
			LeftJoin(table3, fmt.Sprintf("%s.id=%s.user_id", table3, table3)).
			Fields(fmt.Sprintf("%s.number,%s.name", table2, table1)).
			One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["name"].String(), "john")
		t.Assert(one["number"].String(), "n")
	})
}

// https://github.com/gogf/gf/issues/1159
func Test_ScanList_NoRecreate_PtrAttribute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S1 struct {
			Id    int
			Name  string
			Age   int
			Score int
		}
		type S3 struct {
			One *S1
		}
		var (
			s   []*S3
			err error
		)
		r1 := gdb.Result{
			gdb.Record{
				"id":   gvar.New(1),
				"name": gvar.New("john"),
				"age":  gvar.New(16),
			},
			gdb.Record{
				"id":   gvar.New(2),
				"name": gvar.New("smith"),
				"age":  gvar.New(18),
			},
		}
		err = r1.ScanList(&s, "One")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)

		r2 := gdb.Result{
			gdb.Record{
				"id":  gvar.New(1),
				"age": gvar.New(20),
			},
			gdb.Record{
				"id":  gvar.New(2),
				"age": gvar.New(21),
			},
		}
		err = r2.ScanList(&s, "One", "One", "id:Id")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 20)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 21)
	})
}

// https://github.com/gogf/gf/issues/1159
func Test_ScanList_NoRecreate_StructAttribute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S1 struct {
			Id    int
			Name  string
			Age   int
			Score int
		}
		type S3 struct {
			One S1
		}
		var (
			s   []*S3
			err error
		)
		r1 := gdb.Result{
			gdb.Record{
				"id":   gvar.New(1),
				"name": gvar.New("john"),
				"age":  gvar.New(16),
			},
			gdb.Record{
				"id":   gvar.New(2),
				"name": gvar.New("smith"),
				"age":  gvar.New(18),
			},
		}
		err = r1.ScanList(&s, "One")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)

		r2 := gdb.Result{
			gdb.Record{
				"id":  gvar.New(1),
				"age": gvar.New(20),
			},
			gdb.Record{
				"id":  gvar.New(2),
				"age": gvar.New(21),
			},
		}
		err = r2.ScanList(&s, "One", "One", "id:Id")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 20)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 21)
	})
}

// https://github.com/gogf/gf/issues/1159
func Test_ScanList_NoRecreate_SliceAttribute_Ptr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S1 struct {
			Id    int
			Name  string
			Age   int
			Score int
		}
		type S2 struct {
			Id    int
			Pid   int
			Name  string
			Age   int
			Score int
		}
		type S3 struct {
			One  *S1
			Many []*S2
		}
		var (
			s   []*S3
			err error
		)
		r1 := gdb.Result{
			gdb.Record{
				"id":   gvar.New(1),
				"name": gvar.New("john"),
				"age":  gvar.New(16),
			},
			gdb.Record{
				"id":   gvar.New(2),
				"name": gvar.New("smith"),
				"age":  gvar.New(18),
			},
		}
		err = r1.ScanList(&s, "One")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)

		r2 := gdb.Result{
			gdb.Record{
				"id":   gvar.New(100),
				"pid":  gvar.New(1),
				"age":  gvar.New(30),
				"name": gvar.New("john"),
			},
			gdb.Record{
				"id":   gvar.New(200),
				"pid":  gvar.New(1),
				"age":  gvar.New(31),
				"name": gvar.New("smith"),
			},
		}
		err = r2.ScanList(&s, "Many", "One", "pid:Id")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(len(s[0].Many), 2)
		t.Assert(s[0].Many[0].Name, "john")
		t.Assert(s[0].Many[0].Age, 30)
		t.Assert(s[0].Many[1].Name, "smith")
		t.Assert(s[0].Many[1].Age, 31)

		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)
		t.Assert(len(s[1].Many), 0)

		r3 := gdb.Result{
			gdb.Record{
				"id":  gvar.New(100),
				"pid": gvar.New(1),
				"age": gvar.New(40),
			},
			gdb.Record{
				"id":  gvar.New(200),
				"pid": gvar.New(1),
				"age": gvar.New(41),
			},
		}
		err = r3.ScanList(&s, "Many", "One", "pid:Id")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(len(s[0].Many), 2)
		t.Assert(s[0].Many[0].Name, "john")
		t.Assert(s[0].Many[0].Age, 40)
		t.Assert(s[0].Many[1].Name, "smith")
		t.Assert(s[0].Many[1].Age, 41)

		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)
		t.Assert(len(s[1].Many), 0)
	})
}

// https://github.com/gogf/gf/issues/1159
func Test_ScanList_NoRecreate_SliceAttribute_Struct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S1 struct {
			Id    int
			Name  string
			Age   int
			Score int
		}
		type S2 struct {
			Id    int
			Pid   int
			Name  string
			Age   int
			Score int
		}
		type S3 struct {
			One  S1
			Many []S2
		}
		var (
			s   []S3
			err error
		)
		r1 := gdb.Result{
			gdb.Record{
				"id":   gvar.New(1),
				"name": gvar.New("john"),
				"age":  gvar.New(16),
			},
			gdb.Record{
				"id":   gvar.New(2),
				"name": gvar.New("smith"),
				"age":  gvar.New(18),
			},
		}
		err = r1.ScanList(&s, "One")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)

		r2 := gdb.Result{
			gdb.Record{
				"id":   gvar.New(100),
				"pid":  gvar.New(1),
				"age":  gvar.New(30),
				"name": gvar.New("john"),
			},
			gdb.Record{
				"id":   gvar.New(200),
				"pid":  gvar.New(1),
				"age":  gvar.New(31),
				"name": gvar.New("smith"),
			},
		}
		err = r2.ScanList(&s, "Many", "One", "pid:Id")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(len(s[0].Many), 2)
		t.Assert(s[0].Many[0].Name, "john")
		t.Assert(s[0].Many[0].Age, 30)
		t.Assert(s[0].Many[1].Name, "smith")
		t.Assert(s[0].Many[1].Age, 31)

		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)
		t.Assert(len(s[1].Many), 0)

		r3 := gdb.Result{
			gdb.Record{
				"id":  gvar.New(100),
				"pid": gvar.New(1),
				"age": gvar.New(40),
			},
			gdb.Record{
				"id":  gvar.New(200),
				"pid": gvar.New(1),
				"age": gvar.New(41),
			},
		}
		err = r3.ScanList(&s, "Many", "One", "pid:Id")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(len(s[0].Many), 2)
		t.Assert(s[0].Many[0].Name, "john")
		t.Assert(s[0].Many[0].Age, 40)
		t.Assert(s[0].Many[1].Name, "smith")
		t.Assert(s[0].Many[1].Age, 41)

		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)
		t.Assert(len(s[1].Many), 0)
	})
}

func TestResult_Structs1(t *testing.T) {
	type A struct {
		Id int `orm:"id"`
	}
	type B struct {
		*A
		Name string
	}
	gtest.C(t, func(t *gtest.T) {
		r := gdb.Result{
			gdb.Record{"id": gvar.New(nil), "name": gvar.New("john")},
			gdb.Record{"id": gvar.New(1), "name": gvar.New("smith")},
		}
		array := make([]*B, 2)
		err := r.Structs(&array)
		t.AssertNil(err)
		t.Assert(array[0].Id, 0)
		t.Assert(array[1].Id, 1)
		t.Assert(array[0].Name, "john")
		t.Assert(array[1].Name, "smith")
	})
}

func Test_Builder_OmitEmptyWhere(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 1).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 0).OmitEmptyWhere().Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	gtest.C(t, func(t *gtest.T) {
		builder := db.Model(table).OmitEmptyWhere().Builder()
		count, err := db.Model(table).Where(
			builder.Where("id", 0),
		).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
}

func Test_Scan_Nil_Result_Error(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type S struct {
		Id    int
		Name  string
		Age   int
		Score int
	}
	gtest.C(t, func(t *gtest.T) {
		var s *S
		err := db.Model(table).Where("id", 1).Scan(&s)
		t.AssertNil(err)
		t.Assert(s.Id, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		var s *S
		err := db.Model(table).Where("id", 100).Scan(&s)
		t.AssertNil(err)
		t.Assert(s, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		var s S
		err := db.Model(table).Where("id", 100).Scan(&s)
		t.Assert(err, sql.ErrNoRows)
	})
	gtest.C(t, func(t *gtest.T) {
		var ss []*S
		err := db.Model(table).Scan(&ss)
		t.AssertNil(err)
		t.Assert(len(ss), TableSize)
	})
	// If the result is empty, it returns error.
	gtest.C(t, func(t *gtest.T) {
		var ss = make([]*S, 10)
		err := db.Model(table).WhereGT("id", 100).Scan(&ss)
		t.Assert(err, sql.ErrNoRows)
	})
}

func Test_Model_FixGdbJoin(t *testing.T) {
	array := []string{
		`CREATE TABLE common_resource (
			id bigint NOT NULL,
			app_id bigint NOT NULL,
			resource_id varchar(64) NOT NULL,
			src_instance_id varchar(64) DEFAULT NULL,
			region varchar(36) DEFAULT NULL,
			zone varchar(36) DEFAULT NULL,
			database_kind varchar(20) NOT NULL,
			source_type varchar(64) NOT NULL,
			ip varchar(64) DEFAULT NULL,
			port int DEFAULT NULL,
			vpc_id varchar(20) DEFAULT NULL,
			subnet_id varchar(20) DEFAULT NULL,
			proxy_ip varchar(64) DEFAULT NULL,
			proxy_port int DEFAULT NULL,
			proxy_id bigint DEFAULT NULL,
			proxy_snat_ip varchar(64) DEFAULT NULL,
			lease_at datetime2 NULL DEFAULT NULL,
			uin varchar(32) NOT NULL,
			PRIMARY KEY (id)
		)`,
		`INSERT INTO common_resource (id,app_id,resource_id,src_instance_id,region,zone,database_kind,source_type,ip,port,vpc_id,subnet_id,proxy_ip,proxy_port,proxy_id,proxy_snat_ip,lease_at,uin) VALUES ` +
			`(1,1,'2','2','2','3','1','1','1',1,'1','1','1',1,1,'1',NULL,''),` +
			`(3,2,'3','3','3','3','3','3','3',3,'3','3','3',3,3,'3',NULL,''),` +
			`(18,1303697168,'dmc-rgnh9qre','vdb-6b6m3u1u','ap-guangzhou','','vdb','cloud','10.0.1.16',80,'vpc-m3dchft7','subnet-9as3a3z2','9.27.72.189',11131,228476,'169.254.128.5, ','2023-11-08 08:13:04',''),` +
			`(20,1303697168,'dmc-4grzi4jg','tdsqlshard-313spncx','ap-guangzhou','','tdsql','cloud','10.255.0.27',3306,'vpc-407k0e8x','subnet-qhkkk3bo','30.86.239.200',24087,0,'',NULL,'')`,
		`CREATE TABLE managed_resource (
			id bigint NOT NULL,
			instance_id varchar(64) NOT NULL,
			resource_id varchar(64) NOT NULL,
			resource_name varchar(64) DEFAULT NULL,
			status varchar(36) NOT NULL DEFAULT 'valid',
			status_message varchar(64) DEFAULT NULL,
			user_name varchar(64) NOT NULL,
			password varchar(1024) NOT NULL,
			pay_mode tinyint DEFAULT 0,
			safe_publication bit DEFAULT 0,
			created_at datetime2 NOT NULL DEFAULT GETDATE(),
			updated_at datetime2 NOT NULL DEFAULT GETDATE(),
			expired_at datetime2 NULL DEFAULT NULL,
			deleted tinyint NOT NULL DEFAULT 0,
			resource_mark_id int DEFAULT NULL,
			comments varchar(64) DEFAULT NULL,
			rule_template_id varchar(64) NOT NULL,
			PRIMARY KEY (id)
		)`,
		`INSERT INTO managed_resource (id,instance_id,resource_id,resource_name,status,status_message,user_name,password,pay_mode,safe_publication,created_at,updated_at,expired_at,deleted,resource_mark_id,comments,rule_template_id) VALUES ` +
			`(1,'2','3','1','1','1','1','1',1,1,'2023-11-06 12:14:21','2023-11-06 12:14:21',NULL,1,1,'1',''),` +
			`(2,'3','2','1','1','1','1','1',1,0,'2023-11-06 12:15:07','2023-11-06 12:15:07',NULL,1,2,'1',''),` +
			`(5,'dmcins-jxy0x75m','dmc-rgnh9qre','erichmao-vdb-test','invalid','The Ip field is required','root','2e39af3dd1d447e2',1,1,'2023-11-08 08:13:20','2023-11-09 05:31:07',NULL,0,11,NULL,'12345'),` +
			`(6,'dmcins-erxms6ya','dmc-4grzi4jg','erichmao-vdb-test','invalid','The Ip field is required','leotaowang','641d846cf75bc794',1,1,'2023-11-08 22:15:17','2023-11-09 05:31:07',NULL,0,11,NULL,'12345')`,
		`CREATE TABLE rules_template (
			id bigint NOT NULL,
			app_id bigint DEFAULT NULL,
			name varchar(255) NOT NULL,
			database_kind varchar(64) DEFAULT NULL,
			is_default tinyint NOT NULL DEFAULT 0,
			win_rules varchar(2048) DEFAULT NULL,
			inception_rules varchar(2048) DEFAULT NULL,
			auto_exec_rules varchar(2048) DEFAULT NULL,
			order_check_step varchar(2048) DEFAULT NULL,
			template_id varchar(64) NOT NULL DEFAULT '',
			version int NOT NULL DEFAULT 1,
			deleted tinyint NOT NULL DEFAULT 0,
			create_at datetime2 NOT NULL DEFAULT GETDATE(),
			update_at datetime2 NOT NULL DEFAULT GETDATE(),
			is_system tinyint NOT NULL DEFAULT 0,
			uin varchar(64) DEFAULT NULL,
			subAccountUin varchar(64) DEFAULT NULL,
			PRIMARY KEY (id)
		)`,
		`CREATE TABLE resource_mark (
			id bigint NOT NULL,
			app_id bigint NOT NULL,
			mark_name varchar(64) NOT NULL,
			color varchar(11) NOT NULL,
			creator varchar(32) NOT NULL,
			created_at datetime2 NOT NULL DEFAULT GETDATE(),
			updated_at datetime2 NOT NULL DEFAULT GETDATE(),
			PRIMARY KEY (id)
		)`,
		`INSERT INTO resource_mark (id,app_id,mark_name,color,creator,created_at,updated_at) VALUES (10,1,'test','red','1','2023-11-06 02:45:46','2023-11-06 02:45:46')`,
	}
	dropTable(`common_resource`)
	dropTable(`managed_resource`)
	dropTable(`rules_template`)
	dropTable(`resource_mark`)
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(`common_resource`)
	defer dropTable(`managed_resource`)
	defer dropTable(`rules_template`)
	defer dropTable(`resource_mark`)
	gtest.C(t, func(t *gtest.T) {
		t.AssertNil(db.GetCore().ClearCacheAll(ctx))
		sqlSlice, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			orm := db.Model(`managed_resource`).Ctx(ctx).
				LeftJoinOnField(`common_resource`, `resource_id`).
				LeftJoinOnFields(`resource_mark`, `resource_mark_id`, `=`, `id`).
				LeftJoinOnFields(`rules_template`, `rule_template_id`, `=`, `template_id`).
				FieldsPrefix(
					`managed_resource`,
					"resource_id", "user_name", "status", "status_message", "safe_publication", "rule_template_id",
					"created_at", "comments", "expired_at", "resource_mark_id", "instance_id", "resource_name",
					"pay_mode").
				FieldsPrefix(`resource_mark`, "mark_name", "color").
				FieldsPrefix(`rules_template`, "name").
				FieldsPrefix(`common_resource`, `src_instance_id`, "database_kind", "source_type", "ip", "port")
			all, err := orm.OrderAsc("src_instance_id").All()
			t.Assert(err, nil)
			t.Assert(len(all), 4)
			t.Assert(all[0]["pay_mode"], 1)
			t.Assert(all[0]["src_instance_id"], 2)
			t.Assert(all[3]["instance_id"], "dmcins-jxy0x75m")
			t.Assert(all[3]["src_instance_id"], "vdb-6b6m3u1u")
			t.Assert(all[3]["resource_mark_id"], "11")
			return err
		})
		t.AssertNil(err)

		expect := `SELECT "managed_resource"."resource_id","managed_resource"."user_name","managed_resource"."status",` +
			`"managed_resource"."status_message","managed_resource"."safe_publication","managed_resource"."rule_template_id",` +
			`"managed_resource"."created_at","managed_resource"."comments","managed_resource"."expired_at","managed_resource"."resource_mark_id",` +
			`"managed_resource"."instance_id","managed_resource"."resource_name","managed_resource"."pay_mode",` +
			`"resource_mark"."mark_name","resource_mark"."color","rules_template"."name",` +
			`"common_resource"."src_instance_id","common_resource"."database_kind","common_resource"."source_type","common_resource"."ip","common_resource"."port" ` +
			`FROM "managed_resource" ` +
			`LEFT JOIN "common_resource" ON ("managed_resource"."resource_id"="common_resource"."resource_id") ` +
			`LEFT JOIN "resource_mark" ON ("managed_resource"."resource_mark_id" = "resource_mark"."id") ` +
			`LEFT JOIN "rules_template" ON ("managed_resource"."rule_template_id" = "rules_template"."template_id") ` +
			`ORDER BY "src_instance_id" ASC`
		t.Assert(expect, sqlSlice[len(sqlSlice)-1])
	})
}

func Test_Model_Year_Date_Time_DateTime_Timestamp(t *testing.T) {
	table := fmt.Sprintf("date_time_example_%d", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s (
		id        int IDENTITY(1,1) NOT NULL,
		date      date DEFAULT NULL,
		time      time DEFAULT NULL,
		datetime  datetime DEFAULT NULL,
		timestamp datetime2 NULL DEFAULT NULL,
		PRIMARY KEY (id)
	)`, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// insert.
		var now = gtime.Now()
		_, err := db.Model(table).Insert(g.Map{
			"date":      now,
			"time":      now,
			"datetime":  now,
			"timestamp": now,
		})
		t.AssertNil(err)
		// select.
		one, err := db.Model(table).One()
		t.AssertNil(err)
		t.Assert(one["date"].String(), now.Format("Y-m-d"))
		t.Assert(one["time"].String(), now.Format("H:i:s"))
		t.AssertLT(one["datetime"].GTime().Sub(now).Abs(), 5*time.Second)
		t.AssertLT(one["timestamp"].GTime().Sub(now).Abs(), 5*time.Second)
	})
}

func Test_OrderBy_Statement_Generated(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("employee_%d", gtime.TimestampNano())
		array := []string{
			fmt.Sprintf(`CREATE TABLE %s
			(
				id   BIGINT IDENTITY(1,1) PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				age  INT          NOT NULL
			)`, table),
			fmt.Sprintf(`INSERT INTO %s(name, age) VALUES ('John', 30)`, table),
			fmt.Sprintf(`INSERT INTO %s(name, age) VALUES ('Mary', 28)`, table),
		}
		for _, v := range array {
			if _, err := db.Exec(ctx, v); err != nil {
				gtest.Error(err)
			}
		}
		defer dropTable(table)
		sqlArray, _ := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			g.DB("default").Ctx(ctx).Model(table).Order("name asc", "age desc").All()
			return nil
		})
		rawSql := strings.ReplaceAll(sqlArray[len(sqlArray)-1], " ", "")
		expectSql := strings.ReplaceAll(fmt.Sprintf(`SELECT * FROM "%s" ORDER BY "name" asc, "age" desc`, table), " ", "")
		t.Assert(rawSql, expectSql)
	})
}

func Test_Fields_Raw(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		one, err := db.Model(table).Fields(gdb.Raw("1")).One()
		t.AssertNil(err)
		t.Assert(one[""], 1)

		one, err = db.Model(table).Fields(gdb.Raw("2")).One()
		t.AssertNil(err)
		t.Assert(one[""], 2)

		one, err = db.Model(table).Fields(gdb.Raw("2")).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one[""], 2)

		one, err = db.Model(table).Fields(gdb.Raw("2")).Where("id", 10000000000).One()
		t.AssertNil(err)
		t.Assert(len(one), 0)
	})
}
