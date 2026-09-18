// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mariadb_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

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
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["create_time"], data["create_time"])
		t.Assert(one["nickname"], gjson.New(data["nickname"]).MustToJson())
	})
}

func Test_Model_Insert_KeyFieldNameMapping(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			Nickname   string
			CreateTime string
		}
		data := User{
			Id:         1,
			Passport:   "user_1",
			Password:   "pass_1",
			Nickname:   "name_1",
			CreateTime: "2020-10-10 12:00:01",
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data.Passport)
		t.Assert(one["create_time"], data.CreateTime)
		t.Assert(one["nickname"], data.Nickname)
	})
}

func Test_Model_Update_KeyFieldNameMapping(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			Nickname   string
			CreateTime string
		}
		data := User{
			Id:         1,
			Passport:   "user_10",
			Password:   "pass_10",
			Nickname:   "name_10",
			CreateTime: "2020-10-10 12:00:01",
		}
		_, err := db.Model(table).Data(data).WherePri(1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data.Passport)
		t.Assert(one["create_time"], data.CreateTime)
		t.Assert(one["nickname"], data.Nickname)
	})
}

func Test_Model_Insert_Time(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "p1",
			"nickname":    "n1",
			"create_time": "2020-10-10 20:09:18.334",
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).One("id", 1)
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["create_time"], "2020-10-10 20:09:18")
		t.Assert(one["nickname"], data["nickname"])
	})
}

func Test_Model_BatchInsertWithArrayStruct(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		array := garray.New()
		for i := 1; i <= TableSize; i++ {
			array.Append(g.Map{
				"id":          i,
				"uid":         i,
				"passport":    fmt.Sprintf("t%d", i),
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    fmt.Sprintf("name_%d", i),
				"create_time": gtime.Now().String(),
			})
		}

		result, err := user.Data(array).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, TableSize)
	})
}

func Test_Model_UpdateAndGetAffected(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		n, err := db.Model(table).Data("nickname", "T100").
			Where(1).Order("id desc").Limit(2).
			UpdateAndGetAffected()
		t.AssertNil(err)
		t.Assert(n, 2)
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

func Test_Model_Struct_CustomType(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type MyInt int

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         MyInt
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
	})
}

func Test_Model_StructsWithOrmTag(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	dbInvalid.SetDebug(true)
	defer dbInvalid.SetDebug(false)
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
		dbInvalid.GetLogger().(*glog.Logger).SetWriter(buffer)
		defer dbInvalid.GetLogger().(*glog.Logger).SetWriter(os.Stdout)
		dbInvalid.Model(table).Order("id asc").Scan(&users)
		// fmt.Println(buffer.String())
		t.Assert(
			gstr.Contains(
				buffer.String(),
				"SELECT `id`,`Passport`,`password`,`nick_name`,`create_time` FROM `user",
			),
			true,
		)
	})

	// db.SetDebug(true)
	// defer db.SetDebug(false)
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
		t.Assert(one["nickname"], "name_2")
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
	})
}

func Test_Model_Scan_NilSliceAttrWhenNoRecordsFound(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		type Response struct {
			Users []User `json:"users"`
		}
		var res Response
		err := db.Model(table).Scan(&res.Users)
		t.AssertNil(err)
		t.Assert(res.Users, nil)
	})
}

func Test_Model_Option_Map(t *testing.T) {
	// Insert
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Model(table).Fields("id, passport").Data(g.Map{
			"id":       1,
			"passport": "1",
			"password": "1",
			"nickname": "1",
		}).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["password"].String(), "1")
		t.AssertNE(one["nickname"].String(), "1")
		t.Assert(one["passport"].String(), "1")
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["passport"].String(), "0")
		t.AssertNE(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "1")
	})

	// Replace
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		_, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Replace()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["passport"].String(), "0")
		t.AssertNE(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "1")
	})

	// Save
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Model(table).Fields("id, passport").Data(g.Map{
			"id":       1,
			"passport": "1",
			"password": "1",
			"nickname": "1",
		}).Save()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["password"].String(), "1")
		t.AssertNE(one["nickname"].String(), "1")
		t.Assert(one["passport"].String(), "1")
	})
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		_, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Save()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["passport"].String(), "0")
		t.AssertNE(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "1")

		_, err = db.Model(table).Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Save()
		t.AssertNil(err)
		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"].String(), "0")
		t.Assert(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "1")
	})

	// Update
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		r, err := db.Model(table).Data(g.Map{"nickname": ""}).Where("id", 1).Update()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		_, err = db.Model(table).OmitEmptyData().Data(g.Map{"nickname": ""}).Where("id", 2).Update()
		t.AssertNE(err, nil)

		r, err = db.Model(table).OmitEmpty().Data(g.Map{"nickname": "", "password": "123"}).Where("id", 3).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		_, err = db.Model(table).OmitEmpty().Fields("nickname").Data(g.Map{"nickname": "", "password": "123"}).Where("id", 4).Update()
		t.AssertNE(err, nil)

		r, err = db.Model(table).OmitEmpty().
			Fields("password").Data(g.Map{
			"nickname": "",
			"passport": "123",
			"password": "456",
		}).Where("id", 5).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 5).One()
		t.AssertNil(err)
		t.Assert(one["password"], "456")
		t.AssertNE(one["passport"].String(), "")
		t.AssertNE(one["passport"].String(), "123")
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
		t.Assert(list[0]["id"].String(), "1")
		t.Assert(list[0]["nickname"].String(), "")
		t.Assert(list[0]["passport"].String(), "")
		t.Assert(list[0]["password"].String(), "1")

		t.Assert(list[1]["id"].String(), "2")
		t.Assert(list[1]["nickname"].String(), "")
		t.Assert(list[1]["passport"].String(), "")
		t.Assert(list[1]["password"].String(), "2")
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
		t.Assert(list[0]["id"].String(), "1")
		t.Assert(list[0]["nickname"].String(), "")
		t.Assert(list[0]["passport"].String(), "")
		t.Assert(list[0]["password"].String(), "0")

		t.Assert(list[1]["id"].String(), "2")
		t.Assert(list[1]["nickname"].String(), "")
		t.Assert(list[1]["passport"].String(), "")
		t.Assert(list[1]["password"].String(), "2")
	})
}

func Test_Model_FieldsEx_WithReservedWords(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			table      = "fieldsex_test_table"
			sqlTpcPath = gtest.DataPath("reservedwords_table_tpl.sql")
			sqlContent = gfile.GetContents(sqlTpcPath)
		)
		t.AssertNE(sqlContent, "")
		if _, err := db.Exec(ctx, fmt.Sprintf(sqlContent, table)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table)
		_, err := db.Model(table).FieldsEx("content").One()
		t.AssertNil(err)
	})
}

func Test_Model_Prefix(t *testing.T) {
	db := dbPrefix
	table := fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
	createInitTableWithDb(db, TableNamePrefix1+table)
	defer dropTable(TableNamePrefix1 + table)
	// Select.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Where("id in (?)", g.Slice{1, 2}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
	// Select with alias.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table+" as u").Where("u.id in (?)", g.Slice{1, 2}).Order("u.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
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
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).As("u1").LeftJoin(table+" as u2", "u2.id=u1.id").Where("u1.id in (?)", g.Slice{1, 2}).Order("u1.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
}

func Test_Model_Schema1(t *testing.T) {
	// db.SetDebug(true)

	db = db.Schema(TestSchema1)
	table := fmt.Sprintf(`%s_%s`, TableName, gtime.TimestampNanoStr())
	createInitTableWithDb(db, table)
	db = db.Schema(TestSchema2)
	createInitTableWithDb(db, table)
	defer func() {
		db = db.Schema(TestSchema1)
		dropTableWithDb(db, table)
		db = db.Schema(TestSchema2)
		dropTableWithDb(db, table)
		db = db.Schema(TestSchema1)
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
	// db.SetDebug(true)

	db = db.Schema(TestSchema1)
	table := fmt.Sprintf(`%s_%s`, TableName, gtime.TimestampNanoStr())
	createInitTableWithDb(db, table)
	db = db.Schema(TestSchema2)
	createInitTableWithDb(db, table)
	defer func() {
		db = db.Schema(TestSchema1)
		dropTableWithDb(db, table)
		db = db.Schema(TestSchema2)
		dropTableWithDb(db, table)

		db = db.Schema(TestSchema1)
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

func Test_Model_OmitEmpty_Time(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int       `orm:"id"       json:"id"`
			Passport string    `orm:"password" json:"pass_port"`
			Password string    `orm:"password" json:"password"`
			Time     time.Time `orm:"create_time" `
		}
		user := &User{
			Id:       1,
			Passport: "111",
			Password: "222",
			Time:     time.Time{},
		}
		r, err := db.Model(table).OmitEmpty().Data(user).WherePri(1).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)
	})
}

func Test_Result_Chunk(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Order("id asc").All()
		t.AssertNil(err)
		chunks := r.Chunk(3)
		t.Assert(len(chunks), 4)
		t.Assert(chunks[0][0]["id"].Int(), 1)
		t.Assert(chunks[1][0]["id"].Int(), 4)
		t.Assert(chunks[2][0]["id"].Int(), 7)
		t.Assert(chunks[3][0]["id"].Int(), 10)
	})
}

func Test_Model_DryRun(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	db.SetDryRun(true)
	defer db.SetDryRun(false)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["id"], 1)
	})
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Data("passport", "port_1").WherePri(1).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 0)
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
		t.Assert(one["passport"], "user_1")

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
		t.Assert(one["passport"], "user_1")

		time.Sleep(time.Second * 2)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test1",
			Force:    false,
		}).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_100")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test2",
			Force:    false,
		}).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_2")

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
		t.Assert(one["passport"], "user_200")
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
		t.Assert(one["passport"], "user_3")

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
			t.Assert(one["passport"], "user_300")
			return nil
		})
		t.AssertNil(err)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test3",
			Force:    false,
		}).WherePri(3).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_3")
	})
	gtest.C(t, func(t *gtest.T) {
		// make cache for id 4
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test4",
			Force:    false,
		}).WherePri(4).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_4")

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
			t.Assert(one["passport"], "user_400")
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
		t.Assert(one["passport"], "user_4000")
	})
}

func Test_Model_NullField(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Passport *string
		}
		data := g.Map{
			"id":       1,
			"passport": nil,
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)

		var user *User
		err = one.Struct(&user)
		t.AssertNil(err)
		t.Assert(user.Id, data["id"])
		t.Assert(user.Passport, data["passport"])
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

func createTableForTimeZoneTest() string {
	tableName := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        id          int(10) unsigned NOT NULL AUTO_INCREMENT,
	        passport    varchar(45) NULL,
	        password    char(32) NULL,
	        nickname    varchar(45) NULL,
	        created_at timestamp(6) NULL,
 			updated_at timestamp(6) NULL,
			deleted_at timestamp(6) NULL,
	        PRIMARY KEY (id)
	    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	    `, tableName,
	)); err != nil {
		gtest.Fatal(err)
	}
	return tableName
}

func Test_TimeZoneInsert(t *testing.T) {
	tableName := createTableForTimeZoneTest()
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
		t.Assert(userEntity.CreatedAt.String(), "2020-11-22 11:23:45")
		t.Assert(userEntity.UpdatedAt.String(), "2020-11-22 12:23:46")
		t.Assert(gtime.NewFromTime(userEntity.DeletedAt).String(), "2020-11-22 13:23:47")
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
		t.Assert(result["id"], 1)
		t.Assert(result["passport"], "user_1")
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

func Test_Model_Min_Max_Avg_Sum(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Min("id")
		t.AssertNil(err)
		t.Assert(result, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Max("id")
		t.AssertNil(err)
		t.Assert(result, TableSize)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Avg("id")
		t.AssertNil(err)
		t.Assert(result, 5.5)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Sum("id")
		t.AssertNil(err)
		t.Assert(result, 55)
	})
}

func Test_Model_CountColumn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).CountColumn("id")
		t.AssertNil(err)
		t.Assert(result, TableSize)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).CountColumn("id")
		t.AssertNil(err)
		t.Assert(result, 3)
	})
}

func Test_Model_InsertAndGetId(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{
			"id":       1,
			"passport": "user_1",
			"password": "pass_1",
			"nickname": "name_1",
		}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 1)
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
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
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
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
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
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
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
		t.Assert(one["passport"], data["nickname"])
		t.Assert(one["password"], data["nickname"])
		t.Assert(one["nickname"], "name_1")
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
			"passport": gdb.Raw("CONCAT(VALUES(`passport`), '1')"),
			"password": gdb.Raw("CONCAT(VALUES(`password`), '2')"),
		}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"]+"1")
		t.Assert(one["password"], data["password"]+"2")
		t.Assert(one["nickname"], "name_1")
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
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
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
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
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
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
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
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})
}

func Test_Model_Raw(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.
			Raw(fmt.Sprintf("select * from %s where id in (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("id", 8).
			WhereIn("id", g.Slice{1, 2, 3, 4, 5, 6, 7}).
			OrderDesc("id").
			Limit(2).
			All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"], 7)
		t.Assert(all[1]["id"], 5)
	})

	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("select * from %s where id in (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("id", 8).
			WhereIn("id", g.Slice{1, 2, 3, 4, 5, 6, 7}).
			OrderDesc("id").
			Limit(2).
			Count()
		t.AssertNil(err)
		// Raw SQL selects {1,5,7,8,9,10}, Where filters to id < 8 AND id IN {1,2,3,4,5,6,7}
		// Result: {1,5,7} = 3 records
		t.Assert(count, int64(3))
	})
}

func Test_Model_Handler(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Safe().Handler(
			func(m *gdb.Model) *gdb.Model {
				return m.Page(0, 3)
			},
			func(m *gdb.Model) *gdb.Model {
				return m.Where("id", g.Slice{1, 2, 3, 4, 5, 6})
			},
			func(m *gdb.Model) *gdb.Model {
				return m.OrderDesc("id")
			},
		)
		all, err := m.All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"], 6)
		t.Assert(all[2]["id"], 4)
	})
}

func Test_Model_FieldCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldCount("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["total"].Int(), 1)
	})
}

func Test_Model_FieldMax(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldMax("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["total"].Int(), 1)
	})
}

func Test_Model_FieldMin(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldMin("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["total"].Int(), 1)
	})
}

func Test_Model_FieldAvg(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldAvg("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["total"].Int(), 1)
	})
}

func Test_Model_OmitEmptyWhere(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// Basic type where.
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 0).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).OmitEmptyWhere().Where("id", 0).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).OmitEmptyWhere().Where("id", 0).Where("nickname", "").Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	// Slice where.
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", g.Slice{1, 2, 3}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(3))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", g.Slice{}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).OmitEmptyWhere().Where("id", g.Slice{}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", g.Slice{}).OmitEmptyWhere().Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	// Struct Where.
	gtest.C(t, func(t *gtest.T) {
		type Input struct {
			Id   []int
			Name []string
		}
		count, err := db.Model(table).Where(Input{}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		type Input struct {
			Id   []int
			Name []string
		}
		count, err := db.Model(table).Where(Input{}).OmitEmptyWhere().Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	// Map Where.
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where(g.Map{
			"id":       []int{},
			"nickname": []string{},
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		type Input struct {
			Id []int
		}
		count, err := db.Model(table).Where(g.Map{
			"id": []int{},
		}).OmitEmptyWhere().Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
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
		t.Assert(n, 1)

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
		t.Assert(n, 2)

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

// This is no longer used as the filter feature is automatically enabled from GoFrame v1.16.0.
// func Test_Model_Insert_KeyFieldNameMapping_Error(t *testing.T) {
//	table := createTable()
//	defer dropTable(table)
//
//	gtest.C(t, func(t *gtest.T) {
//		type User struct {
//			Id             int
//			Passport       string
//			Password       string
//			Nickname       string
//			CreateTime     string
//			NoneExistField string
//		}
//		data := User{
//			Id:         1,
//			Passport:   "user_1",
//			Password:   "pass_1",
//			Nickname:   "name_1",
//			CreateTime: "2020-10-10 12:00:01",
//		}
//		_, err := db.Model(table).Data(data).Insert()
//		t.AssertNE(err, nil)
//	})
// }

func Test_Model_Fields_AutoFilterInJoinStatement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var err error
		table1 := "user"
		table2 := "score"
		table3 := "info"
		if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		   id int(11) NOT NULL AUTO_INCREMENT,
		   name varchar(500) NOT NULL DEFAULT '',
		 PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;
	    `, table1,
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
		CREATE TABLE IF NOT EXISTS %s (
			id int(11) NOT NULL AUTO_INCREMENT,
			user_id int(11) NOT NULL DEFAULT 0,
		    number varchar(500) NOT NULL DEFAULT '',
		 PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;
	    `, table2,
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
		CREATE TABLE IF NOT EXISTS %s (
			id int(11) NOT NULL AUTO_INCREMENT,
			user_id int(11) NOT NULL DEFAULT 0,
		    description varchar(500) NOT NULL DEFAULT '',
		 PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;
	    `, table3,
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

		one, err := db.Model("user").
			Where("user.id", 1).
			Fields("score.number,user.name").
			LeftJoin("score", "user.id=score.user_id").
			LeftJoin("info", "info.id=info.user_id").
			Order("user.id asc").
			One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["name"].String(), "john")
		t.Assert(one["number"].String(), "n")

		one, err = db.Model("user").
			LeftJoin("score", "user.id=score.user_id").
			LeftJoin("info", "info.id=info.user_id").
			Fields("score.number,user.name").
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
		// fmt.Printf("%+v", err)
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
		// fmt.Printf("%+v", err)
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
		// fmt.Printf("%+v", err)
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
		// fmt.Printf("%+v", err)
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
	array := gstr.SplitAndTrim(gtest.DataContent(`fix_gdb_join.sql`), ";")
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
					"resource_id", "user", "status", "status_message", "safe_publication", "rule_template_id",
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

		t.Assert(gtest.DataContent(`fix_gdb_join_expect.sql`), sqlSlice[len(sqlSlice)-1])
	})
}

func Test_Model_Year_Date_Time_DateTime_Timestamp(t *testing.T) {
	table := "date_time_example"
	array := gstr.SplitAndTrim(gtest.DataContent(`date_time_example.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// insert.
		var now = gtime.Now()
		_, err := db.Model("date_time_example").Insert(g.Map{
			"year":      now,
			"date":      now,
			"time":      now,
			"datetime":  now,
			"timestamp": now,
		})
		t.AssertNil(err)
		// select.
		one, err := db.Model("date_time_example").One()
		t.AssertNil(err)
		t.Assert(one["year"].String(), now.Format("Y"))
		t.Assert(one["date"].String(), now.Format("Y-m-d"))
		t.Assert(one["time"].String(), now.Format("H:i:s"))
		t.AssertLT(one["datetime"].GTime().Sub(now).Seconds(), 5)
		t.AssertLT(one["timestamp"].GTime().Sub(now).Seconds(), 5)
	})
}

func Test_OrderBy_Statement_Generated(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := gstr.SplitAndTrim(gtest.DataContent(`fix_gdb_order_by.sql`), ";")
		for _, v := range array {
			if _, err := db.Exec(ctx, v); err != nil {
				gtest.Error(err)
			}
		}
		defer dropTable(`employee`)
		sqlArray, _ := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			g.DB("default").Ctx(ctx).Model("employee").Order("name asc", "age desc").All()
			return nil
		})
		rawSql := strings.ReplaceAll(sqlArray[len(sqlArray)-1], " ", "")
		expectSql := strings.ReplaceAll("SELECT * FROM `employee` ORDER BY `name` asc, `age` desc", " ", "")
		t.Assert(rawSql, expectSql)
	})
}

func Test_Fields_Raw(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		one, err := db.Model(table).Fields(gdb.Raw("1")).One()
		t.AssertNil(err)
		t.Assert(one["1"], 1)

		one, err = db.Model(table).Fields(gdb.Raw("2")).One()
		t.AssertNil(err)
		t.Assert(one["2"], 2)

		one, err = db.Model(table).Fields(gdb.Raw("2")).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["2"], 2)

		one, err = db.Model(table).Fields(gdb.Raw("2")).Where("id", 10000000000).One()
		t.AssertNil(err)
		t.Assert(len(one), 0)
	})
}
