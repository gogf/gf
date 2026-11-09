// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_New(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		node := gdb.ConfigNode{
			Host: "127.0.0.1",
			Port: "1433",
			User: TestDbUser,
			Pass: TestDbPass,
			Type: "mssql",
		}
		newDb, err := gdb.New(node)
		t.AssertNil(err)
		value, err := newDb.GetValue(ctx, `select 1`)
		t.AssertNil(err)
		t.Assert(value, `1`)
		t.AssertNil(newDb.Close(ctx))
	})
}

func Test_DB_Prepare(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		st, err := db.Prepare(ctx, "SELECT 100")
		t.AssertNil(err)

		rows, err := st.Query()
		t.AssertNil(err)

		array, err := rows.Columns()
		t.AssertNil(err)
		t.Assert(len(array), 1)
		t.Assert(array[0], "")

		t.Assert(rows.Next(), true)
		var value int
		t.AssertNil(rows.Scan(&value))
		t.Assert(value, 100)

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
		t.Assert(one["NICKNAME"], gjson.New(data["nickname"]).MustToJson())
		t.Assert(one["PASSWORD"], `{"salt":"123","pass":"456"}`)
	})
}

func Test_DB_Insert_NilGjson(t *testing.T) {
	var tableName = "nil" + gtime.TimestampNanoStr()
	_, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s (
		id int NOT NULL,
		json_empty_string nvarchar(max) DEFAULT NULL,
		json_nil nvarchar(max) DEFAULT NULL,
		json_null nvarchar(max) DEFAULT NULL,
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

		t.Assert(one["json_empty_string"], nil)
		t.Assert(one["json_nil"], nil)
		t.Assert(one["json_null"], "null")
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
		t.Assert(one["PASSWORD"].String(), "pass_1")
		t.Assert(one["NICKNAME"].String(), "name_1")
		t.Assert(one["CREATE_TIME"].GTime().String(), CreateTime)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.InsertIgnore(ctx, table, g.Map{
			"id":          TableSize + 1,
			"passport":    "t11",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T11",
			"create_time": CreateTime,
		})
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", TableSize+1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"].String(), "t11")
		t.Assert(one["NICKNAME"].String(), "T11")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize+1)
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

func Test_DB_TableField(t *testing.T) {
	name := "field_test"
	dropTable(name)
	defer dropTable(name)
	_, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		field_tinyint  tinyint NULL ,
		field_int  int NULL ,
		field_integer  integer NULL ,
		field_bigint  bigint NULL ,
		field_bit  bit NULL ,
		field_real  real NULL ,
		field_double  float NULL ,
		field_varchar  varchar(10) NULL ,
		field_varbinary  varbinary(255) NULL
	)
	`, name))
	if err != nil {
		gtest.Fatal(err)
	}

	data := gdb.Map{
		"field_tinyint":   1,
		"field_int":       2,
		"field_integer":   3,
		"field_bigint":    4,
		"field_bit":       true,
		"field_real":      123,
		"field_double":    123.25,
		"field_varchar":   "abc",
		"field_varbinary": []byte("aaa"),
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
	dbPrefix, err := gdb.New(gdb.ConfigNode{
		Host:   "127.0.0.1",
		Port:   "1433",
		User:   TestDbUser,
		Pass:   TestDbPass,
		Name:   TestSchema,
		Type:   "mssql",
		Prefix: TableNamePrefix1,
	})
	if err != nil {
		gtest.Fatal(err)
	}
	defer dbPrefix.Close(ctx)
	db := dbPrefix
	name := fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
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

		value, err := db.Model(name).Where("id", id).Value("create_time")
		t.AssertNil(err)
		t.Assert(value.GTime().String(), "2018-10-24 10:00:01")
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

		value, err := db.Model(name).Where("id", id).Value("create_time")
		t.AssertNil(err)
		t.Assert(value.GTime().String(), "2018-10-24 10:00:02")
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
	tableName := "gf_update_counter_test_" + gtime.TimestampNanoStr()
	_, err := db.Exec(ctx, fmt.Sprintf(`
		IF NOT EXISTS (SELECT * FROM sys.objects WHERE name='%s' and type='U')
		CREATE TABLE %s (
		id int NOT NULL,
		views  int DEFAULT 0  NOT NULL ,
		updated_time int DEFAULT 0 NOT NULL
	)
	`, tableName, tableName))
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
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["views"].Int(), 2)
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
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["views"].Int(), 1)
	})
}

func Test_DB_Ctx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := db.Query(ctx, "WAITFOR DELAY '00:00:10'")
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "deadline"), true)
	})
}

func Test_DB_Ctx_Logger(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer db.SetDebug(db.GetDebug())
		db.SetDebug(true)
		ctx := context.WithValue(context.Background(), "Trace-Id", "123456789")
		_, err := db.Query(ctx, "SELECT 1")
		t.AssertNil(err)
	})
}

// All types testing.
func Test_Types(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := "types_" + gtime.TimestampNanoStr()
		if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE %s (
        id int NOT NULL,
        [blob] varbinary(max) NOT NULL,
        [binary] binary(8) NOT NULL,
        [date] date NOT NULL,
        [time] time NOT NULL,
        [timestamp] datetime2(6) NOT NULL,
        [decimal] decimal(5,2) NOT NULL,
        [float] float NOT NULL,
        [bit] bit NOT NULL,
        [tinyint] tinyint NOT NULL,
        [bool] bit NOT NULL,
        PRIMARY KEY (id)
    )
    `, table)); err != nil {
			gtest.Error(err)
		}
		defer dropTable(table)
		data := g.Map{
			"id":        1,
			"blob":      []byte("i love gf"),
			"binary":    []byte("abcdefgh"),
			"date":      "1880-10-24",
			"time":      "10:00:01",
			"timestamp": "2022-02-14 12:00:01.123456",
			"decimal":   -123.456,
			"float":     -123.456,
			"bit":       2,
			"tinyint":   true,
			"bool":      false,
		}
		r, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["blob"].String(), data["blob"])
		t.Assert(one["binary"].String(), data["binary"])
		t.Assert(one["date"].String(), data["date"])
		t.Assert(one["time"].String(), `10:00:01`)
		t.Assert(one["timestamp"].GTime().Format(`Y-m-d H:i:s.u`), `2022-02-14 12:00:01.123`)
		t.Assert(one["decimal"].String(), -123.46)
		t.Assert(one["float"].String(), data["float"])
		t.Assert(one["bit"].Int(), 1)
		t.Assert(one["tinyint"].Bool(), data["tinyint"])
		t.Assert(one["bool"].Bool(), data["bool"])

		type T struct {
			Id        int
			Blob      []byte
			Binary    []byte
			Date      *gtime.Time
			Time      *gtime.Time
			Timestamp *gtime.Time
			Decimal   float64
			Float     float64
			Bit       int8
			TinyInt   bool
		}
		var obj *T
		err = db.Model(table).Scan(&obj)
		t.AssertNil(err)
		t.Assert(obj.Id, 1)
		t.Assert(obj.Blob, data["blob"])
		t.Assert(obj.Binary, data["binary"])
		t.Assert(obj.Date.Format("Y-m-d"), data["date"])
		t.Assert(obj.Time.String(), `10:00:01`)
		t.Assert(obj.Timestamp.Format(`Y-m-d H:i:s.u`), `2022-02-14 12:00:01.123`)
		t.Assert(obj.Decimal, -123.46)
		t.Assert(obj.Float, data["float"])
		t.Assert(obj.Bit, 1)
		t.Assert(obj.TinyInt, data["tinyint"])
	})
}

func Test_DB_DoubleQuote_Preserved(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable(gtime.TimestampNanoStr() + "_quote")
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{
			"id":       1,
			"passport": "user_1",
		}).Insert()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 1).Value("passport")
		t.AssertNil(err)
		t.Assert(value, "user_1")
	})

	gtest.C(t, func(t *gtest.T) {
		table := "quote_" + gtime.TimestampNanoStr()
		if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE %s (
        id int NOT NULL,
        [double] float NULL,
        PRIMARY KEY (id)
    )
    `, table)); err != nil {
			gtest.Error(err)
		}
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{
			"id":     1,
			"double": -123.456,
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{"double": 1.5}).Where("id", 1).Update()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 1).Value("double")
		t.AssertNil(err)
		t.Assert(value.Float64(), 1.5)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		_, err := db.Exec(ctx, fmt.Sprintf(`UPDATE %s SET nickname='say "hi"' WHERE id=1`, table))
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 1).Value("nickname")
		t.AssertNil(err)
		t.Assert(value.String(), `say "hi"`)
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.GetValue(ctx, `SELECT 'say "hi"'`)
		t.AssertNil(err)
		t.Assert(value.String(), `say "hi"`)
	})
}

func Test_Core_ClearTableFields(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(len(fields), 7)
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
