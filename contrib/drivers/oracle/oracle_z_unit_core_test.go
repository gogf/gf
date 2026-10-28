// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/encoding/gxml"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_New(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		node := gdb.ConfigNode{
			Host: TestDbIP,
			Port: TestDbPort,
			User: TestDbUser,
			Pass: TestDbPass,
			Name: TestDbName,
			Type: TestDbType,
		}
		newDb, err := gdb.New(node)
		t.AssertNil(err)
		value, err := newDb.GetValue(ctx, `select 1 from dual`)
		t.AssertNil(err)
		t.Assert(value, `1`)
		t.AssertNil(newDb.Close(ctx))
	})
}

func Test_DB_Prepare(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		st, err := db.Prepare(ctx, "SELECT 100 FROM dual")
		t.AssertNil(err)

		rows, err := st.Query()
		t.AssertNil(err)

		array, err := rows.Columns()
		t.AssertNil(err)
		t.Assert(array[0], "100")

		err = rows.Close()
		t.AssertNil(err)
	})
}

// Fix issue: https://github.com/gogf/gf/issues/819
func Test_DB_Insert_WithStructAndSliceAttribute(t *testing.T) {
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
		_, err := db.Insert(ctx, table, data)
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data["passport"])
		t.Assert(one["CREATE_TIME"], data["create_time"])
		t.Assert(one["NICKNAME"].String(), gjson.New(data["nickname"]).MustToJsonString())
	})

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          2,
			"passport":    "t2",
			"password":    "p2",
			"nickname":    []string{"A", "B", "C"},
			"create_time": gtime.Now().String(),
		}
		_, err := db.Insert(ctx, table, data)
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 2)
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data["passport"])
		t.Assert(one["NICKNAME"].String(), gjson.New(data["nickname"]).MustToJsonString())
	})
}

func Test_DB_Insert_KeyFieldNameMapping(t *testing.T) {
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
		_, err := db.Insert(ctx, table, data)
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data.Passport)
		t.Assert(one["CREATE_TIME"], data.CreateTime)
		t.Assert(one["NICKNAME"], data.Nickname)
	})
}

func Test_DB_Insert_NilGjson(t *testing.T) {
	var tableName = "nil" + gtime.TimestampNanoStr()
	_, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s (
		id NUMBER(10) NOT NULL,
		json_empty_string CLOB DEFAULT NULL,
		json_nil CLOB DEFAULT NULL,
		json_null CLOB DEFAULT NULL,
		PRIMARY KEY (id)
	)
	`, tableName))
	if err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(tableName)

	gtest.C(t, func(t *gtest.T) {
		type Json struct {
			Id              int
			JsonEmptyString *gjson.Json
			JsonNil         *gjson.Json
			JsonNull        *gjson.Json
		}

		data := Json{
			Id:              1,
			JsonEmptyString: gjson.New(""),
			JsonNil:         gjson.New(nil),
			JsonNull:        gjson.New(struct{}{}),
		}

		_, err = db.Insert(ctx, tableName, data)
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", tableName), 1)
		t.AssertNil(err)

		t.AssertEQ(len(one), 4)

		t.Assert(one["JSON_EMPTY_STRING"], nil)
		t.Assert(one["JSON_NIL"], nil)
		t.Assert(one["JSON_NULL"], "null")
	})
}

func Test_DB_Update_KeyFieldNameMapping(t *testing.T) {
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
		_, err := db.Update(ctx, table, data, "id=1")
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data.Passport)
		t.Assert(one["CREATE_TIME"], data.CreateTime)
		t.Assert(one["NICKNAME"], data.Nickname)
	})
}

func Test_DB_InsertIgnore(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Insert(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		t.AssertNE(err, nil)
		t.AssertIN("ORA-00001", err.Error())
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.InsertIgnore(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 0)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"].String(), "user_1")
		t.Assert(one["NICKNAME"].String(), "name_1")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

func Test_DB_Save(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		timeStr := gtime.New("2024-10-01 12:01:01").String()
		_, err := db.Save(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T11",
			"create_time": timeStr,
		})
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 1)
		t.Assert(one["PASSPORT"].String(), "t1")
		t.Assert(one["PASSWORD"].String(), "25d55ad283aa400af464c76d713c07ad")
		t.Assert(one["NICKNAME"].String(), "T11")
		t.Assert(one["CREATE_TIME"].GTime().String(), timeStr)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

func Test_DB_Replace(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		timeStr := gtime.New("2024-10-01 12:01:01").String()
		_, err := db.Replace(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T11",
			"create_time": timeStr,
		})
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 1)
		t.Assert(one["PASSPORT"].String(), "t1")
		t.Assert(one["PASSWORD"].String(), "25d55ad283aa400af464c76d713c07ad")
		t.Assert(one["NICKNAME"].String(), "T11")
		t.Assert(one["CREATE_TIME"].GTime().String(), timeStr)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

func Test_DB_Time(t *testing.T) {
	table := fmt.Sprintf("t_time_%d", gtime.TimestampMicro()%1e9)
	_, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s (
		ID NUMBER(10) NOT NULL,
		PASSPORT VARCHAR2(45) NOT NULL,
		PASSWORD CHAR(32) NOT NULL,
		NICKNAME VARCHAR2(45) NOT NULL,
		CREATE_TIME TIMESTAMP,
		PRIMARY KEY (ID)
	)
	`, table))
	if err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		t1 := time.Now()
		result, err := db.Insert(ctx, table, g.Map{
			"id":          200,
			"passport":    "t200",
			"password":    "123456",
			"nickname":    "T200",
			"create_time": t1,
		})
		if err != nil {
			gtest.Error(err)
		}
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		value, err := db.GetValue(ctx, fmt.Sprintf("select passport from %s where id=?", table), 200)
		t.AssertNil(err)
		t.Assert(value.String(), "t200")
		value, err = db.GetValue(ctx, fmt.Sprintf("select create_time from %s where id=?", table), 200)
		t.AssertNil(err)
		t.Assert(value.GTime().Format("Y-m-d H:i:s"), gtime.New(t1).Format("Y-m-d H:i:s"))
	})

	gtest.C(t, func(t *gtest.T) {
		t1 := time.Now()
		result, err := db.Insert(ctx, table, g.Map{
			"id":          300,
			"passport":    "t300",
			"password":    "123456",
			"nickname":    "T300",
			"create_time": &t1,
		})
		if err != nil {
			gtest.Error(err)
		}
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		value, err := db.GetValue(ctx, fmt.Sprintf("select passport from %s where id=?", table), 300)
		t.AssertNil(err)
		t.Assert(value.String(), "t300")
		value, err = db.GetValue(ctx, fmt.Sprintf("select create_time from %s where id=?", table), 300)
		t.AssertNil(err)
		t.Assert(value.GTime().Format("Y-m-d H:i:s"), gtime.New(t1).Format("Y-m-d H:i:s"))
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Delete(ctx, table, "1=1")
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 2)
	})
}

func Test_DB_ToJson(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields("*").Where("id =? ", 1).All()
		if err != nil {
			gtest.Fatal(err)
		}

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		users := make([]User, 0)

		err = result.Structs(users)
		t.AssertNE(err, nil)

		err = result.Structs(&users)
		if err != nil {
			gtest.Fatal(err)
		}

		// ToJson
		resultJson, err := gjson.LoadContent([]byte(result.Json()))
		if err != nil {
			gtest.Fatal(err)
		}

		t.Assert(users[0].Id, resultJson.Get("0.ID").Int())
		t.Assert(users[0].Passport, resultJson.Get("0.PASSPORT").String())
		t.Assert(users[0].Password, resultJson.Get("0.PASSWORD").String())
		t.Assert(users[0].NickName, resultJson.Get("0.NICKNAME").String())
		t.Assert(users[0].CreateTime, resultJson.Get("0.CREATE_TIME").String())

		result = nil
		t.Assert(result.Structs(&users), sql.ErrNoRows)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields("*").Where("id =? ", 1).One()
		if err != nil {
			gtest.Fatal(err)
		}

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		users := User{}

		err = result.Struct(&users)
		if err != nil {
			gtest.Fatal(err)
		}

		result = nil
		err = result.Struct(&users)
		t.AssertNE(err, nil)
	})
}

func Test_DB_ToXml(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		record, err := db.Model(table).Fields("*").Where("id = ?", 1).One()
		if err != nil {
			gtest.Fatal(err)
		}

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		user := User{}
		err = record.Struct(&user)
		if err != nil {
			gtest.Fatal(err)
		}

		result, err := gxml.Decode([]byte(record.Xml("doc")))
		if err != nil {
			gtest.Fatal(err)
		}

		resultXml := result["doc"].(map[string]any)
		if v, ok := resultXml["ID"]; ok {
			t.Assert(user.Id, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["PASSPORT"]; ok {
			t.Assert(user.Passport, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["PASSWORD"]; ok {
			t.Assert(gstr.TrimRight(user.Password), v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["NICKNAME"]; ok {
			t.Assert(user.NickName, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["CREATE_TIME"]; ok {
			t.Assert(user.CreateTime, v)
		} else {
			gtest.Fatal("FAIL")
		}
	})
}

func Test_DB_ToStringMap(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)
	gtest.C(t, func(t *gtest.T) {
		id := "1"
		result, err := db.Model(table).Fields("*").Where("id = ?", 1).All()
		if err != nil {
			gtest.Fatal(err)
		}

		type t_user struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		t_users := make([]t_user, 0)
		err = result.Structs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultStringMap := result.MapKeyStr("ID")
		t.Assert(t_users[0].Id, resultStringMap[id]["ID"])
		t.Assert(t_users[0].Passport, resultStringMap[id]["PASSPORT"])
		t.Assert(t_users[0].Password, resultStringMap[id]["PASSWORD"])
		t.Assert(t_users[0].NickName, resultStringMap[id]["NICKNAME"])
		t.Assert(t_users[0].CreateTime, resultStringMap[id]["CREATE_TIME"])
	})
}

func Test_DB_ToIntMap(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		id := 1
		result, err := db.Model(table).Fields("*").Where("id = ?", id).All()
		if err != nil {
			gtest.Fatal(err)
		}

		type t_user struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		t_users := make([]t_user, 0)
		err = result.Structs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultIntMap := result.MapKeyInt("ID")
		t.Assert(t_users[0].Id, resultIntMap[id]["ID"])
		t.Assert(t_users[0].Passport, resultIntMap[id]["PASSPORT"])
		t.Assert(t_users[0].Password, resultIntMap[id]["PASSWORD"])
		t.Assert(t_users[0].NickName, resultIntMap[id]["NICKNAME"])
		t.Assert(t_users[0].CreateTime, resultIntMap[id]["CREATE_TIME"])
	})
}

func Test_DB_ToUintMap(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		id := 1
		result, err := db.Model(table).Fields("*").Where("id = ?", id).All()
		if err != nil {
			gtest.Fatal(err)
		}

		type t_user struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		t_users := make([]t_user, 0)
		err = result.Structs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultUintMap := result.MapKeyUint("ID")
		t.Assert(t_users[0].Id, resultUintMap[uint(id)]["ID"])
		t.Assert(t_users[0].Passport, resultUintMap[uint(id)]["PASSPORT"])
		t.Assert(t_users[0].Password, resultUintMap[uint(id)]["PASSWORD"])
		t.Assert(t_users[0].NickName, resultUintMap[uint(id)]["NICKNAME"])
		t.Assert(t_users[0].CreateTime, resultUintMap[uint(id)]["CREATE_TIME"])
	})
}

func Test_DB_ToStringRecord(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		id := 1
		ids := "1"
		result, err := db.Model(table).Fields("*").Where("id = ?", id).All()
		if err != nil {
			gtest.Fatal(err)
		}

		type t_user struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		t_users := make([]t_user, 0)
		err = result.Structs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultStringRecord := result.RecordKeyStr("ID")
		t.Assert(t_users[0].Id, resultStringRecord[ids]["ID"].Int())
		t.Assert(t_users[0].Passport, resultStringRecord[ids]["PASSPORT"].String())
		t.Assert(t_users[0].Password, resultStringRecord[ids]["PASSWORD"].String())
		t.Assert(t_users[0].NickName, resultStringRecord[ids]["NICKNAME"].String())
		t.Assert(t_users[0].CreateTime, resultStringRecord[ids]["CREATE_TIME"].String())
	})
}

func Test_DB_ToIntRecord(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		id := 1
		result, err := db.Model(table).Fields("*").Where("id = ?", id).All()
		if err != nil {
			gtest.Fatal(err)
		}

		type t_user struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		t_users := make([]t_user, 0)
		err = result.Structs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultIntRecord := result.RecordKeyInt("ID")
		t.Assert(t_users[0].Id, resultIntRecord[id]["ID"].Int())
		t.Assert(t_users[0].Passport, resultIntRecord[id]["PASSPORT"].String())
		t.Assert(t_users[0].Password, resultIntRecord[id]["PASSWORD"].String())
		t.Assert(t_users[0].NickName, resultIntRecord[id]["NICKNAME"].String())
		t.Assert(t_users[0].CreateTime, resultIntRecord[id]["CREATE_TIME"].String())
	})
}

func Test_DB_ToUintRecord(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		id := 1
		result, err := db.Model(table).Fields("*").Where("id = ?", id).All()
		if err != nil {
			gtest.Fatal(err)
		}

		type t_user struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		t_users := make([]t_user, 0)
		err = result.Structs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultUintRecord := result.RecordKeyUint("ID")
		t.Assert(t_users[0].Id, resultUintRecord[uint(id)]["ID"].Int())
		t.Assert(t_users[0].Passport, resultUintRecord[uint(id)]["PASSPORT"].String())
		t.Assert(t_users[0].Password, resultUintRecord[uint(id)]["PASSWORD"].String())
		t.Assert(t_users[0].NickName, resultUintRecord[uint(id)]["NICKNAME"].String())
		t.Assert(t_users[0].CreateTime, resultUintRecord[uint(id)]["CREATE_TIME"].String())
	})
}

func Test_DB_TableField(t *testing.T) {
	name := fmt.Sprintf("field_test_%d", gtime.TimestampMicro()%1e9)
	dropTable(name)
	defer dropTable(name)
	_, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		field_tinyint  NUMBER(3) NULL,
		field_int  NUMBER(8) NULL,
		field_integer  NUMBER(8) NULL,
		field_bigint  NUMBER(19) NULL,
		field_bit  NUMBER(1) NULL,
		field_real  NUMBER(8,0) NULL,
		field_double  NUMBER(12,2) NULL,
		field_varchar  VARCHAR2(10) NULL,
		field_varbinary  RAW(255) NULL
	)
	`, name))
	if err != nil {
		gtest.Fatal(err)
	}

	data := gdb.Map{
		"FIELD_TINYINT":   1,
		"FIELD_INT":       2,
		"FIELD_INTEGER":   3,
		"FIELD_BIGINT":    4,
		"FIELD_BIT":       6,
		"FIELD_REAL":      123,
		"FIELD_DOUBLE":    123.25,
		"FIELD_VARCHAR":   "abc",
		"FIELD_VARBINARY": []byte("aaa"),
	}
	gtest.C(t, func(t *gtest.T) {
		res, err := db.Model(name).Data(data).Insert()
		if err != nil {
			t.Fatal(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		} else {
			t.Assert(n, 1)
		}

		result, err := db.Model(name).Fields("*").Where("field_int = ?", 2).All()
		if err != nil {
			t.Fatal(err)
		}
		t.Assert(result[0], data)
	})

}

func Test_DB_Prefix(t *testing.T) {
	db, err := gdb.New(gdb.ConfigNode{
		Host:   TestDbIP,
		Port:   TestDbPort,
		User:   TestDbUser,
		Pass:   TestDbPass,
		Name:   TestDbName,
		Type:   TestDbType,
		Prefix: TableNamePrefix1,
	})
	gtest.AssertNil(err)
	defer db.Close(ctx)
	name := fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampMicro()%1e9)
	table := TableNamePrefix1 + name
	createTableWithDb(db, table)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		result, err := db.Insert(ctx, name, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:00").String(),
		})
		t.AssertNil(err)

		n, e := result.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		result, err := db.Replace(ctx, name, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:01").String(),
		})
		t.AssertNil(err)

		n, e := result.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, 1)

		v, err := db.Model(name).Where("id", id).Value("create_time")
		t.AssertNil(err)
		t.Assert(v.String(), "2018-10-24 10:00:01")
	})

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		result, err := db.Save(ctx, name, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:02").String(),
		})
		t.AssertNil(err)

		n, e := result.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, 1)

		v, err := db.Model(name).Where("id", id).Value("create_time")
		t.AssertNil(err)
		t.Assert(v.String(), "2018-10-24 10:00:02")
	})

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		result, err := db.Update(ctx, name, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:03").String(),
		}, "id=?", id)
		t.AssertNil(err)

		n, e := result.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		result, err := db.Delete(ctx, name, "id=?", id)
		t.AssertNil(err)

		n, e := result.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		for i := 1; i <= TableSize; i++ {
			array.Append(g.Map{
				"id":          i,
				"passport":    fmt.Sprintf(`user_%d`, i),
				"password":    fmt.Sprintf(`pass_%d`, i),
				"nickname":    fmt.Sprintf(`name_%d`, i),
				"create_time": gtime.NewFromStr("2018-10-24 10:00:00").String(),
			})
		}

		result, err := db.Insert(ctx, name, array.Slice())
		t.AssertNil(err)

		n, e := result.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, TableSize)
	})
}

// update counter test.
func Test_DB_UpdateCounter(t *testing.T) {
	tableName := fmt.Sprintf("gf_update_counter_%d", gtime.TimestampMicro()%1e9)
	_, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		id NUMBER(10) NOT NULL,
		views NUMBER(8) DEFAULT 0 NOT NULL,
		updated_time NUMBER(10) DEFAULT 0 NOT NULL
	)
	`, tableName))
	if err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(tableName)

	gtest.C(t, func(t *gtest.T) {
		insertData := g.Map{
			"id":           1,
			"views":        0,
			"updated_time": 0,
		}
		_, err = db.Insert(ctx, tableName, insertData)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		gdbCounter := &gdb.Counter{
			Field: "id",
			Value: 1,
		}
		updateData := g.Map{
			"views": gdbCounter,
		}
		result, err := db.Update(ctx, tableName, updateData, "id", 1)
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(tableName).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 1)
		t.Assert(one["VIEWS"].Int(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		gdbCounter := &gdb.Counter{
			Field: "views",
			Value: -1,
		}
		updateData := g.Map{
			"views":        gdbCounter,
			"updated_time": gtime.Now().Unix(),
		}
		result, err := db.Update(ctx, tableName, updateData, "id", 1)
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(tableName).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 1)
		t.Assert(one["VIEWS"].Int(), 1)
	})
}

func Test_DB_Ctx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// Note: on Windows the statement runs to its end before ORA-01013 is returned, as go-ora v2.9.0
		// cannot interrupt it there (sendOOB in network/net_windows.go does nothing).
		_, err := db.Query(ctx, "BEGIN DBMS_LOCK.SLEEP(3); END;")
		t.AssertNE(err, nil)
		t.AssertIN("ORA-01013", err.Error())
	})
}

func Test_DB_Ctx_Logger(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer db.SetDebug(db.GetDebug())
		db.SetDebug(true)
		ctx := context.WithValue(context.Background(), "Trace-Id", "123456789")
		_, err := db.Query(ctx, "SELECT 1 FROM dual")
		t.AssertNil(err)
	})
}

// All types testing.
func Test_Types(t *testing.T) {
	table := fmt.Sprintf("t_types_%d", gtime.TimestampMicro()%1e9)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE %s (
        id NUMBER(10) NOT NULL,
        col_blob BLOB NOT NULL,
        col_binary RAW(8) NOT NULL,
        col_date DATE NOT NULL,
        col_timestamp TIMESTAMP(6) NOT NULL,
        col_decimal NUMBER(5,2) NOT NULL,
        col_double BINARY_DOUBLE NOT NULL,
        col_bit NUMBER(1) NOT NULL,
        col_tinyint NUMBER(1) NOT NULL,
        col_bool NUMBER(1) NOT NULL,
        PRIMARY KEY (id)
    )
    `, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	data := g.Map{
		"id":            1,
		"col_blob":      []byte("i love gf"),
		"col_binary":    []byte("abcdefgh"),
		"col_date":      gtime.New("1880-10-24"),
		"col_timestamp": gtime.New("2022-02-14 12:00:01.123456"),
		"col_decimal":   -123.456,
		"col_double":    -123.456,
		"col_bit":       2,
		"col_tinyint":   true,
		"col_bool":      false,
	}
	type T struct {
		Id           int
		ColBlob      []byte
		ColBinary    []byte
		ColDate      *gtime.Time
		ColTimestamp *gtime.Time
		ColDecimal   float64
		ColDouble    float64
		ColBit       int8
		ColTinyint   bool
	}
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 1)
		t.Assert(one["COL_BLOB"].String(), data["col_blob"])
		t.Assert(one["COL_BINARY"].String(), data["col_binary"])
		t.Assert(one["COL_DATE"].String(), "1880-10-24 00:00:00")
		t.Assert(one["COL_DECIMAL"].String(), -123.46)
		t.Assert(one["COL_DOUBLE"].String(), data["col_double"])
		t.Assert(one["COL_BIT"].Int(), data["col_bit"])
		t.Assert(one["COL_TINYINT"].Bool(), data["col_tinyint"])
		t.Assert(one["COL_BOOL"].Bool(), data["col_bool"])

		var obj *T
		err = db.Model(table).Scan(&obj)
		t.AssertNil(err)
		t.Assert(obj.Id, 1)
		t.Assert(obj.ColBlob, data["col_blob"])
		t.Assert(obj.ColBinary, data["col_binary"])
		t.Assert(obj.ColDate.Format("Y-m-d"), "1880-10-24")
		t.Assert(obj.ColDecimal, -123.46)
		t.Assert(obj.ColDouble, data["col_double"])
		t.Assert(obj.ColBit, data["col_bit"])
		t.Assert(obj.ColTinyint, data["col_tinyint"])
	})

	gtest.C(t, func(t *gtest.T) {
		const expect = `2022-02-14 12:00:01.123`
		one, err := db.Model(table).One()
		t.AssertNil(err)
		t.Assert(one["COL_TIMESTAMP"].GTime().Format(`Y-m-d H:i:s.u`), expect)

		var obj *T
		err = db.Model(table).Scan(&obj)
		t.AssertNil(err)
		t.Assert(obj.ColTimestamp.Format(`Y-m-d H:i:s.u`), expect)
	})
}

func Test_Core_ClearTableFields(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(len(fields), 6)
	})
	gtest.C(t, func(t *gtest.T) {
		err := db.GetCore().ClearTableFields(ctx, table)
		t.AssertNil(err)
	})
}

func Test_Core_ClearTableFieldsAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.GetCore().ClearTableFieldsAll(ctx)
		t.AssertNil(err)
	})
}

func Test_Core_ClearCache(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.GetCore().ClearCache(ctx, "")
		t.AssertNil(err)
	})
}

func Test_Core_ClearCacheAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.GetCore().ClearCacheAll(ctx)
		t.AssertNil(err)
	})
}

func Test_DB_Insert_ToSQL(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			_, err := db.Ctx(ctx).Model(table).Data(g.Map{
				"id":       1,
				"passport": "t1",
				"password": "p1",
				"nickname": "n1",
			}).Insert()
			return err
		})
		t.AssertNil(err)
		t.Assert(sql, fmt.Sprintf("INSERT INTO %s(ID,NICKNAME,PASSPORT,PASSWORD) VALUES(1,'n1','t1','p1')", table))

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		var id int64
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) (err error) {
			id, err = db.Ctx(ctx).Model(table).Data(g.Map{
				"passport": "t2",
				"password": "p2",
				"nickname": "n2",
			}).InsertAndGetId()
			return err
		})
		t.AssertNil(err)
		t.Assert(id, 0)
		t.Assert(sql, fmt.Sprintf("INSERT INTO %s(NICKNAME,PASSPORT,PASSWORD) VALUES('n2','t2','p2') RETURNING ID INTO :v4", table))

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

func Test_DB_Insert_CatchSQL(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var id int64
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) (err error) {
			id, err = db.Ctx(ctx).Model(table).Data(g.Map{
				"passport": "t1",
				"password": "p1",
				"nickname": "n1",
			}).InsertAndGetId()
			return err
		})
		t.AssertNil(err)
		t.Assert(id, 1)
		t.AssertGT(len(sqlArray), 0)
		t.Assert(
			sqlArray[len(sqlArray)-1],
			fmt.Sprintf("INSERT INTO %s(NICKNAME,PASSPORT,PASSWORD) VALUES('n1','t1','p1') RETURNING ID INTO :v4", table),
		)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "t1")
	})
}

func Test_DB_Save_ToSQL(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			_, err := db.Ctx(ctx).Model(table).Data(g.Map{
				"id":       1,
				"passport": "t1",
				"password": "p1",
				"nickname": "n1",
			}).Save()
			return err
		})
		t.AssertNil(err)
		t.Assert(sql, fmt.Sprintf(
			"MERGE INTO %s T1 USING (SELECT 1 AS ID,'n1' AS NICKNAME,'t1' AS PASSPORT,'p1' AS PASSWORD FROM DUAL) T2 ON (T1.ID = T2.ID) "+
				"WHEN NOT MATCHED THEN INSERT(ID,NICKNAME,PASSPORT,PASSWORD) VALUES (T2.ID,T2.NICKNAME,T2.PASSPORT,T2.PASSWORD)"+
				"WHEN MATCHED THEN UPDATE SET T1.NICKNAME = T2.NICKNAME,T1.PASSPORT = T2.PASSPORT,T1.PASSWORD = T2.PASSWORD",
			table,
		))

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["NICKNAME"], "name_1")
	})
}
