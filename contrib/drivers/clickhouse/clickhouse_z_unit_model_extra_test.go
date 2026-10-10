// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

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
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
)

const chACreateTime = "2018-10-24 10:00:00"

// chANode returns a configuration node of the test server with the given changes applied.
func chANode(fn func(node *gdb.ConfigNode)) gdb.ConfigNode {
	node := gdb.ConfigNode{
		Host:  "127.0.0.1",
		Port:  "9000",
		User:  "default",
		Name:  "default",
		Type:  "clickhouse",
		Extra: "mutations_sync=1",
	}
	if fn != nil {
		fn(&node)
	}
	return node
}

// chANewDB creates a database object of a new configuration group connected to the test server with the
// given node changes applied.
func chANewDB(fn func(node *gdb.ConfigNode)) gdb.DB {
	group := chAName("chA")
	gdb.AddConfigNode(group, chANode(fn))
	newDb, err := gdb.NewByGroup(group)
	gtest.AssertNil(err)
	return newDb
}

// chAName returns a unique table or database name with the given prefix, as other test processes share
// the server.
func chAName(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, guid.S())
}

// chACreateTable creates a table shaped like the MySQL test table: nullable columns plus create_date.
func chACreateTable(table ...string) string {
	return chACreateTableWithDb(db, table...)
}

// chACreateInitTable creates a chACreateTable table and inserts TableSize records like the MySQL tests.
func chACreateInitTable(table ...string) string {
	return chACreateInitTableWithDb(db, table...)
}

func chACreateTableWithDb(db gdb.DB, table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = chAName(TableName)
	}
	dropTableWithDb(db, name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id          UInt64,
			passport    Nullable(String),
			password    Nullable(String),
			nickname    Nullable(String),
			create_time Nullable(DateTime64(6)),
			create_date Nullable(Date)
		) ENGINE = MergeTree() ORDER BY id`, name,
	)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

func chACreateInitTableWithDb(db gdb.DB, table ...string) (name string) {
	name = chACreateTableWithDb(db, table...)
	array := garray.New(true)
	for i := 1; i <= TableSize; i++ {
		array.Append(g.Map{
			"id":          i,
			"passport":    fmt.Sprintf(`user_%d`, i),
			"password":    fmt.Sprintf(`pass_%d`, i),
			"nickname":    fmt.Sprintf(`name_%d`, i),
			"create_time": gtime.NewFromStr(chACreateTime).String(),
		})
	}
	_, err := db.Insert(ctx, name, array.Slice())
	gtest.AssertNil(err)
	count, err := db.Model(name).Count()
	gtest.AssertNil(err)
	gtest.Assert(count, TableSize)
	return
}

// chACount returns the number of records in `table` matching the optional where condition.
func chACount(table string, where ...any) int {
	m := db.Model(table)
	if len(where) > 0 {
		m = m.Where(where[0], where[1:]...)
	}
	count, err := m.Count()
	gtest.AssertNil(err)
	return count
}

// Test_Model_Insert_Supplement tests the Model.Insert scenarios of the MySQL suite not covered by Test_Model_Insert.
func Test_Model_Insert_Supplement(t *testing.T) {
	table := chACreateTable()
	defer dropTable(table)

	// Note: ClickHouse reports neither last insert id nor affected rows, so the effects are verified
	// by reading the table back.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table, "id", 1), 1)

		_, err = db.Model(table).Data(g.Map{
			"id":          "2",
			"uid":         "2",
			"passport":    "t2",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table, "id", 2), 1)

		type User struct {
			Id         int         `gconv:"id"`
			Uid        int         `gconv:"uid"`
			Passport   string      `json:"passport"`
			Password   string      `gconv:"password"`
			Nickname   string      `gconv:"nickname"`
			CreateTime *gtime.Time `json:"create_time"`
		}
		_, err = db.Model(table).Data(User{
			Id:       3,
			Uid:      3,
			Passport: "t3",
			Password: "25d55ad283aa400af464c76d713c07ad",
			Nickname: "name_3",
		}).Insert()
		t.AssertNil(err)
		value, err := db.Model(table).Fields("passport").Where("id=3").Value()
		t.AssertNil(err)
		t.Assert(value.String(), "t3")

		_, err = db.Model(table).Data(&User{
			Id:         4,
			Uid:        4,
			Passport:   "t4",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "T4",
			CreateTime: gtime.Now(),
		}).Insert()
		t.AssertNil(err)
		value, err = db.Model(table).Fields("passport").Where("id=4").Value()
		t.AssertNil(err)
		t.Assert(value.String(), "t4")

		_, err = db.Model(table).Where("id>?", 1).Delete()
		t.AssertNil(err)
		t.Assert(chACount(table), 1)
	})
}

// Test_Model_Insert_WithStructAndSliceAttribute tests inserting struct and slice values as JSON strings.
func Test_Model_Insert_WithStructAndSliceAttribute(t *testing.T) {
	table := chACreateTable()
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

// Test_Model_Insert_KeyFieldNameMapping tests inserting a struct whose attributes map to table fields.
func Test_Model_Insert_KeyFieldNameMapping(t *testing.T) {
	table := chACreateTable()
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

// Test_Model_Update_KeyFieldNameMapping tests updating with a struct whose attributes map to table fields.
func Test_Model_Update_KeyFieldNameMapping(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

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
	// Note: ClickHouse cannot UPDATE a sorting key column, and the struct carries the key field id.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(data).WherePri(1).Update()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).FieldsEx("id").Data(data).WherePri(1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data.Passport)
		t.Assert(one["create_time"], data.CreateTime)
		t.Assert(one["nickname"], data.Nickname)
	})
}

// Test_Model_Insert_Time tests inserting a datetime string with fractional seconds.
func Test_Model_Insert_Time(t *testing.T) {
	table := chACreateTable()
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

// Test_Model_BatchInsertWithArrayStruct tests batch inserting a garray of maps.
func Test_Model_BatchInsertWithArrayStruct(t *testing.T) {
	table := chACreateTable()
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

		_, err := user.Data(array).Insert()
		t.AssertNil(err)
		// Note: ClickHouse has no last insert id, so the inserted records are counted instead.
		t.Assert(chACount(table), TableSize)
	})
}

// Test_Model_InsertIgnore_Supplement tests inserting a record whose primary key already exists.
func Test_Model_InsertIgnore_Supplement(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
	// Note: ClickHouse has no unique constraint, so a duplicate primary key is inserted as a new record.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table, "id", 1), 2)
	})
}

// Test_Model_Batch tests batch insert, save and replace.
func Test_Model_Batch(t *testing.T) {
	// batch insert
	gtest.C(t, func(t *gtest.T) {
		table := chACreateTable()
		defer dropTable(table)
		_, err := db.Model(table).Data(g.List{
			{
				"id":          2,
				"uid":         2,
				"passport":    "t2",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_2",
				"create_time": gtime.Now().String(),
			},
			{
				"id":          3,
				"uid":         3,
				"passport":    "t3",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_3",
				"create_time": gtime.Now().String(),
			},
		}).Batch(1).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 2)
	})

	// batch insert with the ids given explicitly, as ClickHouse has no auto-increment.
	gtest.C(t, func(t *gtest.T) {
		table := chACreateTable()
		defer dropTable(table)
		_, err := db.Model(table).Data(g.List{
			{"id": 1, "passport": "t1"},
			{"id": 2, "passport": "t2"},
			{"id": 3, "passport": "t3"},
			{"id": 4, "passport": "t4"},
			{"id": 5, "passport": "t5"},
		}).Batch(2).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 5)
	})

	// Note: Save and Replace insert the records as ClickHouse has no upsert, so the table doubles.
	// batch save
	gtest.C(t, func(t *gtest.T) {
		table := chACreateInitTable()
		defer dropTable(table)
		result, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		for _, v := range result {
			v["nickname"].Set(v["nickname"].String() + v["id"].String())
		}
		_, err = db.Model(table).Data(result).Save()
		t.AssertNil(err)
		t.Assert(chACount(table), TableSize*2)
		t.Assert(chACount(table, "nickname", "name_11"), 1)
	})

	// batch replace
	gtest.C(t, func(t *gtest.T) {
		table := chACreateInitTable()
		defer dropTable(table)
		result, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		for _, v := range result {
			v["nickname"].Set(v["nickname"].String() + v["id"].String())
		}
		_, err = db.Model(table).Data(result).Replace()
		t.AssertNil(err)
		t.Assert(chACount(table), TableSize*2)
		t.Assert(chACount(table, "nickname", "name_11"), 1)
	})
}

// Test_Model_Replace_Supplement tests that Model.Replace writes the record.
func Test_Model_Replace_Supplement(t *testing.T) {
	table := chACreateTable()
	defer dropTable(table)

	// Note: ClickHouse has no upsert, so Replace inserts the record each time.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "t11",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T11",
			"create_time": "2018-10-24 10:00:00",
		}
		_, err := db.Model(table).Data(data).Replace()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "t11")
		t.Assert(one["create_time"], "2018-10-24 10:00:00")

		_, err = db.Model(table).Data(data).Replace()
		t.AssertNil(err)
		t.Assert(chACount(table, "id", 1), 2)
	})
}

// Test_Model_Save tests that Model.Save writes the record.
func Test_Model_Save(t *testing.T) {
	table := chACreateTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "t111",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T111",
			"create_time": "2018-10-24 10:00:00",
		}).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "t111")
		t.Assert(one["nickname"], "T111")
		t.Assert(chACount(table), 1)
	})
}

// Test_Model_Update_Supplement tests the Model.Update scenarios of the MySQL suite not covered by Test_Model_Update.
func Test_Model_Update_Supplement(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
	// UPDATE...LIMIT
	// Note: ClickHouse mutations support neither ORDER BY nor LIMIT, so the statement fails and changes nothing.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data("nickname", "T100").Where(1).Order("id desc").Limit(2).Update()
		t.AssertNE(err, nil)

		v1, err := db.Model(table).Fields("nickname").Where("id", 10).Value()
		t.AssertNil(err)
		t.Assert(v1.String(), "name_10")

		v2, err := db.Model(table).Fields("nickname").Where("id", 8).Value()
		t.AssertNil(err)
		t.Assert(v2.String(), "name_8")
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data("passport", "user_22").Where("passport=?", "user_2").Update()
		t.AssertNil(err)
		t.Assert(chACount(table, "passport", "user_22"), 1)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data("passport", "user_2").Where("passport='user_22'").Update()
		t.AssertNil(err)
		t.Assert(chACount(table, "passport", "user_2"), 1)
		t.Assert(chACount(table, "passport", "user_22"), 0)
	})

	// Update + Data(string)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data("passport='user_33'").Where("passport='user_3'").Update()
		t.AssertNil(err)
		t.Assert(chACount(table, "passport", "user_33"), 1)
	})
	// Update + Fields(string)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Fields("passport").Data(g.Map{
			"passport": "user_44",
			"none":     "none",
		}).Where("passport='user_4'").Update()
		t.AssertNil(err)
		t.Assert(chACount(table, "passport", "user_44"), 1)
	})
}

// Test_Model_UpdateAndGetAffected tests UpdateAndGetAffected.
func Test_Model_UpdateAndGetAffected(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
	// Note: ClickHouse mutations support neither ORDER BY nor LIMIT and report no affected rows, so the
	// two records are selected by id, UpdateAndGetAffected updates them and fails as not supported,
	// and the update is verified by reading the table back.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data("nickname", "T100").
			Where("id", g.Slice{9, 10}).
			UpdateAndGetAffected()
		t.AssertNE(err, nil)
		t.Assert(chACount(table, "nickname", "T100"), 2)
	})
}

// Test_Model_Clone tests that a safe model can be reused for several queries.
func Test_Model_Clone(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		md := db.Model(table).Safe(true).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.AssertNil(err)

		record, err := md.Safe(true).Order("id DESC").One()
		t.AssertNil(err)

		result, err := md.Safe(true).Order("id ASC").All()
		t.AssertNil(err)

		t.Assert(count, int64(2))
		t.Assert(record["id"].Int(), 3)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 3)
	})
}

// Test_Model_Safe tests the chaining safety of Model.Safe.
func Test_Model_Safe(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		md := db.Model(table).Safe(false).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))

		md.Where("id = ?", 1)
		count, err = md.Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
	gtest.C(t, func(t *gtest.T) {
		md := db.Model(table).Safe(true).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))

		md.Where("id = ?", 1)
		count, err = md.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))
	})

	gtest.C(t, func(t *gtest.T) {
		md := db.Model(table).Safe().Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))

		md.Where("id = ?", 1)
		count, err = md.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))
	})
	gtest.C(t, func(t *gtest.T) {
		md1 := db.Model(table).Safe()
		md2 := md1.Where("id in (?)", g.Slice{1, 3})
		count, err := md2.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))

		all, err := md2.All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		all, err = md2.Page(1, 10).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		table := chACreateInitTable()
		defer dropTable(table)

		md1 := db.Model(table).Where("id>", 0).Safe()
		md2 := md1.Where("id in (?)", g.Slice{1, 3})
		md3 := md1.Where("id in (?)", g.Slice{4, 5, 6})

		// 1,3
		count, err := md2.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))

		all, err := md2.Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[1]["id"].Int(), 3)

		all, err = md2.Page(1, 10).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		// 4,5,6
		count, err = md3.Count()
		t.AssertNil(err)
		t.Assert(count, int64(3))

		all, err = md3.Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"].Int(), 4)
		t.Assert(all[1]["id"].Int(), 5)
		t.Assert(all[2]["id"].Int(), 6)

		all, err = md3.Page(1, 10).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
	})
}

// Test_Model_All_Supplement tests Model.All returning nil when no record matches.
func Test_Model_All_Supplement(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id<0").All()
		t.Assert(result, nil)
		t.AssertNil(err)
	})
}

// Test_Model_Fields tests Model.Fields with aliases and join statements.
func Test_Model_Fields(t *testing.T) {
	tableName1 := chACreateInitTable()
	defer dropTable(tableName1)

	tableName2 := chAName("user")
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id   UInt64,
			name Nullable(String),
			age  Nullable(UInt32)
		) ENGINE = MergeTree() ORDER BY id`, tableName2,
	)); err != nil {
		gtest.AssertNil(err)
	}
	defer dropTable(tableName2)

	_, err := db.Insert(ctx, tableName2, g.Map{
		"id":   1,
		"name": "table2_1",
		"age":  18,
	})
	gtest.AssertNil(err)
	gtest.Assert(chACount(tableName2), 1)

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

// Test_Model_One_Supplement tests Model.One with a matching and a missing record.
func Test_Model_One_Supplement(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		record, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(record["nickname"].String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		record, err := db.Model(table).Where("id", 0).One()
		t.AssertNil(err)
		t.Assert(record, nil)
	})
}

// Test_Model_Value tests Model.Value with a matching and a missing record.
func Test_Model_Value(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Fields("nickname").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Fields("nickname").Where("id", 0).Value()
		t.AssertNil(err)
		t.Assert(value, nil)
	})
}

// Test_Model_Count_Supplement tests the Model.Count scenarios of the MySQL suite not covered by Test_Model_Count.
func Test_Model_Count_Supplement(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	// Count with cache, check internal ctx data feature.
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 10; i++ {
			count, err := db.Model(table).Cache(gdb.CacheOption{
				Duration: time.Second * 10,
				Name:     guid.S(),
				Force:    false,
			}).Count()
			t.AssertNil(err)
			t.Assert(count, int64(TableSize))
		}
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).FieldsEx("id").Where("id>8").Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Fields("distinct id,nickname").Where("id>8").Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))
	})
	// COUNT...LIMIT...
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Page(1, 2).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
}

// Test_Model_Value_WithCache tests Model.Value with the query cache.
func Test_Model_Value_WithCache(t *testing.T) {
	table := chACreateTable()
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
		_, err := db.Model(table).Data(g.MapStrAny{
			"id":       1,
			"passport": fmt.Sprintf(`passport_%d`, 1),
			"password": fmt.Sprintf(`password_%d`, 1),
			"nickname": fmt.Sprintf(`nickname_%d`, 1),
		}).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 1)
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

// Test_Model_Count_WithCache tests Model.Count with the query cache and a where condition.
func Test_Model_Count_WithCache(t *testing.T) {
	table := chACreateTable()
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
		_, err := db.Model(table).Data(g.MapStrAny{
			"id":       1,
			"passport": fmt.Sprintf(`passport_%d`, 1),
			"password": fmt.Sprintf(`password_%d`, 1),
			"nickname": fmt.Sprintf(`nickname_%d`, 1),
		}).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 1)
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

// Test_Model_Count_All_WithCache tests Model.Count with the query cache and no where condition.
func Test_Model_Count_All_WithCache(t *testing.T) {
	table := chACreateTable()
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
		_, err := db.Model(table).Data(g.MapStrAny{
			"id":       1,
			"passport": fmt.Sprintf(`passport_%d`, 1),
			"password": fmt.Sprintf(`password_%d`, 1),
			"nickname": fmt.Sprintf(`nickname_%d`, 1),
		}).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 1)
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
		_, err := db.Model(table).Data(g.MapStrAny{
			"id":       2,
			"passport": fmt.Sprintf(`passport_%d`, 2),
			"password": fmt.Sprintf(`password_%d`, 2),
			"nickname": fmt.Sprintf(`nickname_%d`, 2),
		}).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 2)
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

// Test_Model_CountColumn_WithCache tests Model.CountColumn with the query cache.
func Test_Model_CountColumn_WithCache(t *testing.T) {
	table := chACreateTable()
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
		_, err := db.Model(table).Data(g.MapStrAny{
			"id":       1,
			"passport": fmt.Sprintf(`passport_%d`, 1),
			"password": fmt.Sprintf(`password_%d`, 1),
			"nickname": fmt.Sprintf(`nickname_%d`, 1),
		}).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 1)
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

// Test_Model_Select tests scanning all records into a struct slice.
func Test_Model_Select(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	type User struct {
		Id         int
		Passport   string
		Password   string
		NickName   string
		CreateTime gtime.Time
	}
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.Model(table).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
	})
}

// Test_Model_Struct tests scanning one record into a struct.
func Test_Model_Struct(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
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
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
	})
	// Auto creating struct object.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := (*User)(nil)
		err := db.Model(table).Where("id=1").Scan(&user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
	})
	// Just using Scan.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := (*User)(nil)
		err := db.Model(table).Where("id=1").Scan(&user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
	})
	// sql.ErrNoRows
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=-1").Scan(user)
		t.Assert(err, sql.ErrNoRows)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var user *User
		err := db.Model(table).Where("id=-1").Scan(&user)
		t.AssertNil(err)
	})
}

// Test_Model_Struct_CustomType tests scanning into a struct with a custom integer type.
func Test_Model_Struct_CustomType(t *testing.T) {
	table := chACreateInitTable()
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

// Test_Model_Structs tests scanning records into struct slices.
func Test_Model_Structs(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
	})
	// Auto create struct slice.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
	})
	// Just using Scan.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
	})
	// sql.ErrNoRows
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Model(table).Where("id<0").Scan(&users)
		t.AssertNil(err)
	})
}

// Test_Model_StructsWithOrmTag tests the fields selected for structs with orm tags.
func Test_Model_StructsWithOrmTag(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	dbInvalid := chANewDB(func(node *gdb.ConfigNode) {
		node.Port = "1"
	})
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
		_ = dbInvalid.Model(table).Order("id asc").Scan(&users)
		t.Assert(
			gstr.Contains(
				buffer.String(),
				fmt.Sprintf("SELECT id,Passport,password,nick_name,create_time FROM %s", table),
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
		t.Assert(one["nickname"], "name_2")
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
	})
}

// Test_Model_Scan_Supplement tests the Model.Scan scenarios of the MySQL suite not covered by Test_Model_Scan.
func Test_Model_Scan_Supplement(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
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
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
	})
	// sql.ErrNoRows
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var (
			user  = new(User)
			users = new([]*User)
		)
		err1 := db.Model(table).Where("id < 0").Scan(user)
		err2 := db.Model(table).Where("id < 0").Scan(users)
		t.Assert(err1, sql.ErrNoRows)
		t.Assert(err2, nil)
	})
}

// Test_Model_Scan_NilSliceAttrWhenNoRecordsFound tests that scanning no records leaves a slice attribute nil.
func Test_Model_Scan_NilSliceAttrWhenNoRecordsFound(t *testing.T) {
	table := chACreateTable()
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

// Test_Model_OrderBy tests Model.Order with columns and raw expressions.
func Test_Model_OrderBy(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("id DESC").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["nickname"].String(), fmt.Sprintf("name_%d", TableSize))
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order(gdb.Raw("NULL")).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["nickname"].String(), "name_1")
	})

	// Note: ClickHouse has no FIELD function, indexOf over an array literal gives the same ordering.
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order(gdb.Raw("indexOf([10,1,2,3,4,5,6,7,8,9], id)")).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["nickname"].String(), "name_10")
		t.Assert(result[1]["nickname"].String(), "name_1")
		t.Assert(result[2]["nickname"].String(), "name_2")
	})
}

// Test_Model_GroupBy tests Model.Group.
func Test_Model_GroupBy(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	// Note: ClickHouse requires every selected column to be grouped or aggregated, so "SELECT *" with
	// "GROUP BY id" fails, and the grouped query selects the grouped columns in an explicit order.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Group("id").All()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields("id,nickname").Group("id,nickname").Order("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["nickname"].String(), "name_1")
	})
}

// Test_Model_Data tests Model.Data with a string condition, a map slice and a garray.
func Test_Model_Data(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := chACreateInitTable()
		defer dropTable(table)
		_, err := db.Model(table).Data("nickname=?", "test").Where("id=?", 3).Update()
		t.AssertNil(err)
		t.Assert(chACount(table, "nickname", "test"), 1)
		t.Assert(chACount(table, "id=3 AND nickname='test'"), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		table := chACreateTable()
		defer dropTable(table)
		users := make([]g.MapStrAny, 0)
		for i := 1; i <= 10; i++ {
			users = append(users, g.MapStrAny{
				"id":       i,
				"passport": fmt.Sprintf(`passport_%d`, i),
				"password": fmt.Sprintf(`password_%d`, i),
				"nickname": fmt.Sprintf(`nickname_%d`, i),
			})
		}
		_, err := db.Model(table).Data(users).Batch(2).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 10)
	})
	gtest.C(t, func(t *gtest.T) {
		table := chACreateTable()
		defer dropTable(table)
		users := garray.New()
		for i := 1; i <= 10; i++ {
			users.Append(g.MapStrAny{
				"id":       i,
				"passport": fmt.Sprintf(`passport_%d`, i),
				"password": fmt.Sprintf(`password_%d`, i),
				"nickname": fmt.Sprintf(`nickname_%d`, i),
			})
		}
		_, err := db.Model(table).Data(users).Batch(2).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 10)
	})
}

// Test_Model_Delete_Supplement tests the Model.Delete scenarios of the MySQL suite not covered by Test_Model_Delete.
func Test_Model_Delete_Supplement(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	// DELETE...LIMIT
	// Note: ClickHouse mutations do not support LIMIT, so the statement fails and deletes nothing.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Where(1).Limit(2).Delete()
		t.AssertNE(err, nil)
		t.Assert(chACount(table), TableSize)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Where(1).Delete()
		t.AssertNil(err)
		t.Assert(chACount(table), 0)
	})
}

// Test_Model_Offset tests Model.Offset.
func Test_Model_Offset(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Limit(2).Offset(5).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"], 6)
		t.Assert(result[1]["id"], 7)
	})
}

// Test_Model_Page tests Model.Page.
func Test_Model_Page(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(3, 3).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"], 7)
		t.Assert(result[1]["id"], 8)
	})
	gtest.C(t, func(t *gtest.T) {
		model := db.Model(table).Safe().Order("id")
		all, err := model.Page(3, 3).All()
		t.AssertNil(err)
		count, err := model.Count()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"], "7")
		t.Assert(count, int64(TableSize))
	})
}

// Test_Model_Option_Map tests the omit options with map data.
func Test_Model_Option_Map(t *testing.T) {
	// Insert
	gtest.C(t, func(t *gtest.T) {
		table := chACreateTable()
		defer dropTable(table)
		_, err := db.Model(table).Fields("id, passport").Data(g.Map{
			"id":       1,
			"passport": "1",
			"password": "1",
			"nickname": "1",
		}).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 1)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["password"].String(), "1")
		t.AssertNE(one["nickname"].String(), "1")
		t.Assert(one["passport"].String(), "1")
	})

	gtest.C(t, func(t *gtest.T) {
		table := chACreateTable()
		defer dropTable(table)
		_, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 1)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["passport"].String(), "0")
		t.AssertNE(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "1")
	})

	// Replace
	// Note: ClickHouse has no upsert, so Replace inserts a second record with id 1 next to the initial one.
	gtest.C(t, func(t *gtest.T) {
		table := chACreateInitTable()
		defer dropTable(table)
		_, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Replace()
		t.AssertNil(err)
		t.Assert(chACount(table, "id", 1), 2)
		one, err := db.Model(table).Where("id", 1).Where("nickname", "1").One()
		t.AssertNil(err)
		t.AssertNE(one["passport"].String(), "0")
		t.AssertNE(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "1")
	})

	// Save
	gtest.C(t, func(t *gtest.T) {
		table := chACreateTable()
		defer dropTable(table)
		_, err := db.Model(table).Fields("id, passport").Data(g.Map{
			"id":       1,
			"passport": "1",
			"password": "1",
			"nickname": "1",
		}).Save()
		t.AssertNil(err)
		t.Assert(chACount(table), 1)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["password"].String(), "1")
		t.AssertNE(one["nickname"].String(), "1")
		t.Assert(one["passport"].String(), "1")
	})
	// Note: ClickHouse has no upsert, so the second Save inserts a second record with id 1.
	gtest.C(t, func(t *gtest.T) {
		table := chACreateTable()
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
		t.Assert(chACount(table, "id", 1), 2)
		one, err = db.Model(table).Where("id", 1).Where("passport", "0").One()
		t.AssertNil(err)
		t.Assert(one["passport"].String(), "0")
		t.Assert(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "1")
	})

	// Update
	gtest.C(t, func(t *gtest.T) {
		table := chACreateInitTable()
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{"nickname": ""}).Where("id", 1).Update()
		t.AssertNil(err)
		t.Assert(chACount(table, "id=1 AND nickname=''"), 1)

		_, err = db.Model(table).OmitEmptyData().Data(g.Map{"nickname": ""}).Where("id", 2).Update()
		t.AssertNE(err, nil)

		_, err = db.Model(table).OmitEmpty().Data(g.Map{"nickname": "", "password": "123"}).Where("id", 3).Update()
		t.AssertNil(err)
		t.Assert(chACount(table, "id=3 AND password='123' AND nickname='name_3'"), 1)

		_, err = db.Model(table).OmitEmpty().Fields("nickname").Data(g.Map{"nickname": "", "password": "123"}).Where("id", 4).Update()
		t.AssertNE(err, nil)

		_, err = db.Model(table).OmitEmpty().
			Fields("password").Data(g.Map{
			"nickname": "",
			"passport": "123",
			"password": "456",
		}).Where("id", 5).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 5).One()
		t.AssertNil(err)
		t.Assert(one["password"], "456")
		t.AssertNE(one["passport"].String(), "")
		t.AssertNE(one["passport"].String(), "123")
	})
}

// Test_Model_Option_List tests the omit options with list data.
func Test_Model_Option_List(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := chACreateTable()
		defer dropTable(table)
		_, err := db.Model(table).Fields("id, password").Data(g.List{
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
		table := chACreateTable()
		defer dropTable(table)
		_, err := db.Model(table).OmitEmpty().Fields("id, password").Data(g.List{
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

// chACreateNameTable creates a table with a non-nullable name column.
func chACreateNameTable() string {
	table := chAName("table")
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id   UInt64,
			name String
		) ENGINE = MergeTree() ORDER BY id`, table)); err != nil {
		gtest.Error(err)
	}
	return table
}

// Test_Model_OmitEmpty tests OmitEmpty, OmitEmptyData and OmitEmptyWhere with Save.
func Test_Model_OmitEmpty(t *testing.T) {
	table := chACreateNameTable()
	defer dropTable(table)

	// Note: ClickHouse fills an omitted non-nullable column with its default value instead of failing.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitEmpty().Data(g.Map{
			"id":   1,
			"name": "",
		}).Save()
		t.AssertNil(err)
		t.Assert(chACount(table, "id=1 AND name=''"), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"id":   2,
			"name": "",
		}).Save()
		t.AssertNil(err)
		t.Assert(chACount(table, "id=2 AND name=''"), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitEmptyWhere().Data(g.Map{
			"id":   1,
			"name": "",
		}).Save()
		t.AssertNil(err)
	})
}

// Test_Model_OmitNil tests OmitNil and OmitNilWhere with Save.
func Test_Model_OmitNil(t *testing.T) {
	table := chACreateNameTable()
	defer dropTable(table)

	// Note: ClickHouse fills an omitted non-nullable column with its default value instead of failing.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitNil().Data(g.Map{
			"id":   1,
			"name": nil,
		}).Save()
		t.AssertNil(err)
		t.Assert(chACount(table, "id=1 AND name=''"), 1)
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

// Test_Model_FieldsEx tests Model.FieldsEx for select and update.
func Test_Model_FieldsEx(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
	// Select.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).FieldsEx("create_time, id").Where("id in (?)", g.Slice{1, 2}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(len(r[0]), 4)
		t.Assert(r[0]["id"], "")
		t.Assert(r[0]["passport"], "user_1")
		t.Assert(r[0]["password"], "pass_1")
		t.Assert(r[0]["nickname"], "name_1")
		t.Assert(r[0]["create_time"], "")
		t.Assert(r[1]["id"], "")
		t.Assert(r[1]["passport"], "user_2")
		t.Assert(r[1]["password"], "pass_2")
		t.Assert(r[1]["nickname"], "name_2")
		t.Assert(r[1]["create_time"], "")
	})
	// Update.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).FieldsEx("password").Data(g.Map{"nickname": "123", "password": "456"}).Where("id", 3).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "123")
		t.AssertNE(one["password"], "456")
	})
}

// Test_Model_FieldsEx_WithReservedWords tests FieldsEx on a table whose columns are reserved words.
func Test_Model_FieldsEx_WithReservedWords(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := chAName("fieldsex_test_table")
		if _, err := db.Exec(ctx, fmt.Sprintf(`
			CREATE TABLE %s (
				id          UInt64,
				key         Nullable(String),
				category_id UInt32,
				user_id     UInt32,
				title       String,
				content     String,
				sort        Nullable(UInt32) DEFAULT 0,
				brief       Nullable(String),
				thumb       Nullable(String),
				tags        Nullable(String),
				referer     Nullable(String),
				status      Nullable(UInt16) DEFAULT 0,
				view_count  Nullable(UInt32) DEFAULT 0,
				zan_count   Nullable(UInt32),
				cai_count   Nullable(UInt32),
				created_at  Nullable(DateTime),
				updated_at  Nullable(DateTime)
			) ENGINE = MergeTree() ORDER BY id`, table,
		)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table)
		_, err := db.Model(table).FieldsEx("content").One()
		t.AssertNil(err)
	})
}

// Test_Model_Prefix tests table name prefix with aliases and joins.
func Test_Model_Prefix(t *testing.T) {
	const prefix = "gf_"
	db := chANewDB(func(node *gdb.ConfigNode) {
		node.Prefix = prefix
	})
	table := chAName(TableName)
	chACreateInitTableWithDb(db, prefix+table)
	defer dropTable(prefix + table)
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

// chACreateSchemas creates two databases for the schema tests and returns their names.
func chACreateSchemas() (string, string) {
	var (
		schema1 = chAName("test1")
		schema2 = chAName("test2")
	)
	for _, schema := range []string{schema1, schema2} {
		if _, err := db.Exec(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", schema)); err != nil {
			gtest.Fatal(err)
		}
	}
	return schema1, schema2
}

// chADropSchemas drops the databases created by chACreateSchemas.
func chADropSchemas(schemas ...string) {
	for _, schema := range schemas {
		if _, err := db.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", schema)); err != nil {
			gtest.Error(err)
		}
	}
}

// Test_Model_Schema1 tests switching schema through DB.Schema and Model.Schema.
func Test_Model_Schema1(t *testing.T) {
	schema1, schema2 := chACreateSchemas()
	defer chADropSchemas(schema1, schema2)

	db := chANewDB(nil).Schema(schema1)
	table := chAName(TableName)
	chACreateInitTableWithDb(db, table)
	db = db.Schema(schema2)
	chACreateInitTableWithDb(db, table)
	// Method.
	gtest.C(t, func(t *gtest.T) {
		db = db.Schema(schema1)
		_, err := db.Model(table).Update(g.Map{"nickname": "name_100"}, "id=1")
		t.AssertNil(err)

		v, err := db.Model(table).Value("nickname", "id=1")
		t.AssertNil(err)
		t.Assert(v.String(), "name_100")

		db = db.Schema(schema2)
		v, err = db.Model(table).Value("nickname", "id=1")
		t.AssertNil(err)
		t.Assert(v.String(), "name_1")
	})
	// Model.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Schema(schema1).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_2")

		_, err = db.Model(table).Schema(schema1).Update(g.Map{"nickname": "name_200"}, "id=2")
		t.AssertNil(err)

		v, err = db.Model(table).Schema(schema1).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_200")

		v, err = db.Model(table).Schema(schema2).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_2")

		v, err = db.Model(table).Schema(schema1).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_200")
	})
	// Model.
	gtest.C(t, func(t *gtest.T) {
		i := 1000
		_, err := db.Model(table).Schema(schema1).Insert(g.Map{
			"id":               i,
			"passport":         fmt.Sprintf(`user_%d`, i),
			"password":         fmt.Sprintf(`pass_%d`, i),
			"nickname":         fmt.Sprintf(`name_%d`, i),
			"create_time":      gtime.NewFromStr("2018-10-24 10:00:00").String(),
			"none-exist-field": 1,
		})
		t.AssertNil(err)

		v, err := db.Model(table).Schema(schema1).Value("nickname", "id=?", i)
		t.AssertNil(err)
		t.Assert(v.String(), "name_1000")

		v, err = db.Model(table).Schema(schema2).Value("nickname", "id=?", i)
		t.AssertNil(err)
		t.Assert(v.String(), "")
	})
}

// Test_Model_Schema2 tests switching schema through DB.Schema before creating the model.
func Test_Model_Schema2(t *testing.T) {
	schema1, schema2 := chACreateSchemas()
	defer chADropSchemas(schema1, schema2)

	db := chANewDB(nil)
	table := chAName(TableName)
	chACreateInitTableWithDb(db.Schema(schema1), table)
	chACreateInitTableWithDb(db.Schema(schema2), table)
	// Schema.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Schema(schema1).Model(table).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_2")

		_, err = db.Schema(schema1).Model(table).Update(g.Map{"nickname": "name_200"}, "id=2")
		t.AssertNil(err)

		v, err = db.Schema(schema1).Model(table).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_200")

		v, err = db.Schema(schema2).Model(table).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_2")

		v, err = db.Schema(schema1).Model(table).Value("nickname", "id=2")
		t.AssertNil(err)
		t.Assert(v.String(), "name_200")
	})
	// Schema.
	gtest.C(t, func(t *gtest.T) {
		i := 1000
		_, err := db.Schema(schema1).Model(table).Insert(g.Map{
			"id":               i,
			"passport":         fmt.Sprintf(`user_%d`, i),
			"password":         fmt.Sprintf(`pass_%d`, i),
			"nickname":         fmt.Sprintf(`name_%d`, i),
			"create_time":      gtime.NewFromStr("2018-10-24 10:00:00").String(),
			"none-exist-field": 1,
		})
		t.AssertNil(err)

		v, err := db.Schema(schema1).Model(table).Value("nickname", "id=?", i)
		t.AssertNil(err)
		t.Assert(v.String(), "name_1000")

		v, err = db.Schema(schema2).Model(table).Value("nickname", "id=?", i)
		t.AssertNil(err)
		t.Assert(v.String(), "")
	})
}

// Test_Model_FieldsExStruct tests FieldsEx with OmitEmpty on struct data.
func Test_Model_FieldsExStruct(t *testing.T) {
	table := chACreateTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int    `orm:"id"       json:"id"`
			Passport string `orm:"password" json:"pass_port"`
			Password string `orm:"password" json:"password"`
			NickName string `orm:"nickname" json:"nick__name"`
		}
		user := &User{
			Id:       1,
			Passport: "111",
			Password: "222",
			NickName: "333",
		}
		_, err := db.Model(table).FieldsEx("create_time, password").OmitEmpty().Data(user).Insert()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "333")
		t.Assert(one["password"].IsNil(), true)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int    `orm:"id"       json:"id"`
			Passport string `orm:"password" json:"pass_port"`
			Password string `orm:"password" json:"password"`
			NickName string `orm:"nickname" json:"nick__name"`
		}
		users := make([]*User, 0)
		for i := 100; i < 110; i++ {
			users = append(users, &User{
				Id:       i,
				Passport: fmt.Sprintf(`passport_%d`, i),
				Password: fmt.Sprintf(`password_%d`, i),
				NickName: fmt.Sprintf(`nickname_%d`, i),
			})
		}
		_, err := db.Model(table).FieldsEx("create_time, password").
			OmitEmpty().
			Batch(2).
			Data(users).
			Insert()
		t.AssertNil(err)
		t.Assert(chACount(table, "id>=100"), 10)
	})
}

// Test_Model_OmitEmpty_Time tests that OmitEmpty skips a zero time.Time attribute on update.
func Test_Model_OmitEmpty_Time(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
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
	// Note: ClickHouse cannot UPDATE a sorting key column, and the struct carries the key field id.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitEmpty().Data(user).WherePri(1).Update()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitEmpty().FieldsEx("id").Data(user).WherePri(1).Update()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["create_time"], chACreateTime)
	})
}

// Test_Result_Chunk tests Result.Chunk.
func Test_Result_Chunk(t *testing.T) {
	table := chACreateInitTable()
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

// Test_Model_DryRun tests that dry run executes queries but not writes.
func Test_Model_DryRun(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
	db := chANewDB(nil)
	db.SetDryRun(true)
	defer db.SetDryRun(false)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["id"], 1)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data("passport", "port_1").WherePri(1).Update()
		t.AssertNil(err)
		t.Assert(chACount(table, "passport", "user_1"), 1)
		t.Assert(chACount(table, "passport", "port_1"), 0)
	})
}

// Test_Model_Join_SubQuery tests joining a sub query.
func Test_Model_Join_SubQuery(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		subQuery := fmt.Sprintf("select * from `%s`", table)
		r, err := db.Model(table, "t1").Fields("t2.id").LeftJoin(subQuery, "t2", "t2.id=t1.id").Order("t2.id").Array()
		t.AssertNil(err)
		t.Assert(len(r), TableSize)
		t.Assert(r[0], "1")
		t.Assert(r[TableSize-1], TableSize)
	})
}

// Test_Model_Cache tests the query cache with updates and transactions.
func Test_Model_Cache(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)
	var (
		cacheName1 = "chA_test1_" + guid.S()
		cacheName2 = "chA_test2_" + guid.S()
		cacheName3 = "chA_test3_" + guid.S()
		cacheName4 = "chA_test4_" + guid.S()
	)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     cacheName1,
			Force:    false,
		}).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		_, err = db.Model(table).Data("passport", "user_100").WherePri(1).Update()
		t.AssertNil(err)
		t.Assert(chACount(table, "passport", "user_100"), 1)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     cacheName1,
			Force:    false,
		}).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		time.Sleep(time.Second * 2)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     cacheName1,
			Force:    false,
		}).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_100")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     cacheName2,
			Force:    false,
		}).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_2")

		_, err = db.Model(table).Data("passport", "user_200").Cache(gdb.CacheOption{
			Duration: -1,
			Name:     cacheName2,
			Force:    false,
		}).WherePri(2).Update()
		t.AssertNil(err)
		t.Assert(chACount(table, "passport", "user_200"), 1)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     cacheName2,
			Force:    false,
		}).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_200")
	})
	// transaction.
	// Note: ClickHouse has no transactions, so the transaction closures fail without running and the
	// cached records stay in place.
	gtest.C(t, func(t *gtest.T) {
		// make cache for id 3
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     cacheName3,
			Force:    false,
		}).WherePri(3).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_3")

		_, err = db.Model(table).Data("passport", "user_300").Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     cacheName3,
			Force:    false,
		}).WherePri(3).Update()
		t.AssertNil(err)
		t.Assert(chACount(table, "passport", "user_300"), 1)

		err = db.Transaction(context.TODO(), func(ctx context.Context, tx gdb.TX) error {
			return nil
		})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:Transaction")

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     cacheName3,
			Force:    false,
		}).WherePri(3).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_3")
	})
	gtest.C(t, func(t *gtest.T) {
		// make cache for id 4
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     cacheName4,
			Force:    false,
		}).WherePri(4).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_4")

		_, err = db.Model(table).Data("passport", "user_400").Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     cacheName3,
			Force:    false,
		}).WherePri(4).Update()
		t.AssertNil(err)
		t.Assert(chACount(table, "passport", "user_400"), 1)

		err = db.Transaction(context.TODO(), func(ctx context.Context, tx gdb.TX) error {
			return nil
		})
		t.AssertNE(err, nil)
		// Read from cache.
		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     cacheName4,
			Force:    false,
		}).WherePri(4).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_4")
	})
}

// Test_Model_Having tests Model.Having.
func Test_Model_Having(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id > 1").Having("id > 8").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id > 1").Having("id > ?", 8).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id > ?", 1).Having("id > ?", 8).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id > ?", 1).Having("id", 8).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
	})
}

// Test_Model_Distinct tests distinct fields and Model.Distinct.
func Test_Model_Distinct(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table, "t").Fields("distinct t.id").Where("id > 1").Having("id > 8").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id > 1").Distinct().Count()
		t.AssertNil(err)
		t.Assert(count, int64(9))
	})
}

// Test_Model_Min_Max tests min and max expressions in Fields.
func Test_Model_Min_Max(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table, "t").Fields("min(t.id)").Where("id > 1").Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table, "t").Fields("max(t.id)").Where("id > 1").Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 10)
	})
}

// Test_Model_Fields_AutoMapping tests mapping field names case-insensitively in Fields.
func Test_Model_Fields_AutoMapping(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Fields("ID").Where("id", 2).Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Fields("NICK_NAME").Where("id", 2).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "name_2")
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields(g.Map{
			"ID":        1,
			"NICK_NAME": 1,
		}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["id"], 2)
		t.Assert(one["nickname"], "name_2")
	})
	// Struct
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			ID       int
			NICKNAME int
		}
		one, err := db.Model(table).Fields(&T{
			ID:       0,
			NICKNAME: 0,
		}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["id"], 2)
		t.Assert(one["nickname"], "name_2")
	})
}

// Test_Model_FieldsEx_AutoMapping tests mapping field names case-insensitively in FieldsEx.
func Test_Model_FieldsEx_AutoMapping(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).FieldsEx("create_date, Passport, Password, NickName, CreateTime").Where("id", 2).Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).FieldsEx("create_date, ID, Passport, Password, CreateTime").Where("id", 2).Value()
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
		t.Assert(len(one), 3)
		t.Assert(one["id"], 2)
		t.Assert(one["nickname"], "name_2")
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
		t.Assert(len(one), 3)
		t.Assert(one["id"], 2)
		t.Assert(one["nickname"], "name_2")
	})
}

// Test_Model_Fields_Struct tests Fields with structs and embedded structs.
func Test_Model_Fields_Struct(t *testing.T) {
	table := chACreateInitTable()
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
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields(&A{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields(B{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 3)
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
		t.Assert(one["nickname"], "name_2")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields(&B{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 3)
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
		t.Assert(one["nickname"], "name_2")
	})
}

// Test_Model_NullField tests scanning a NULL value into a pointer attribute.
func Test_Model_NullField(t *testing.T) {
	table := chACreateTable()
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
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)

		var user *User
		err = one.Struct(&user)
		t.AssertNil(err)
		t.Assert(user.Id, data["id"])
		t.Assert(user.Passport, data["passport"])
	})
}

// Test_Model_Empty_Slice_Argument tests where conditions with an empty slice.
func Test_Model_Empty_Slice_Argument(t *testing.T) {
	table := chACreateInitTable()
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

// Test_Model_HasTable tests Core.HasTable.
func Test_Model_HasTable(t *testing.T) {
	table := chACreateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		t.AssertNil(db.GetCore().ClearCacheAll(ctx))
		result, err := db.GetCore().HasTable(table)
		t.Assert(result, true)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		t.AssertNil(db.GetCore().ClearCacheAll(ctx))
		result, err := db.GetCore().HasTable("table12321")
		t.Assert(result, false)
		t.AssertNil(err)
	})
}

// Test_Model_HasField tests Model.HasField.
func Test_Model_HasField(t *testing.T) {
	table := chACreateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).HasField("id")
		t.Assert(result, true)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).HasField("id123")
		t.Assert(result, false)
		t.AssertNil(err)
	})
}

// chACreateTableForTimeZoneTest creates a table with the soft time fields.
func chACreateTableForTimeZoneTest() string {
	tableName := chAName("user")
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id         UInt64,
			passport   Nullable(String),
			password   Nullable(String),
			nickname   Nullable(String),
			created_at Nullable(DateTime64(6)),
			updated_at Nullable(DateTime64(6)),
			deleted_at Nullable(DateTime64(6))
		) ENGINE = MergeTree() ORDER BY id`, tableName,
	)); err != nil {
		gtest.Fatal(err)
	}
	return tableName
}

// Test_TimeZoneInsert tests inserting times of another time zone.
// https://github.com/gogf/gf/issues/1012
func Test_TimeZoneInsert(t *testing.T) {
	tableName := chACreateTableForTimeZoneTest()
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

// Test_Model_Fields_Map_Struct tests Fields with maps and structs containing unknown fields.
func Test_Model_Fields_Map_Struct(t *testing.T) {
	table := chACreateInitTable()
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

// Test_Model_Min_Max_Avg_Sum tests the aggregate methods.
func Test_Model_Min_Max_Avg_Sum(t *testing.T) {
	table := chACreateInitTable()
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

// Test_Model_CountColumn tests Model.CountColumn.
func Test_Model_CountColumn(t *testing.T) {
	table := chACreateInitTable()
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

// Test_Model_InsertAndGetId tests Model.InsertAndGetId.
func Test_Model_InsertAndGetId(t *testing.T) {
	table := chACreateTable()
	defer dropTable(table)
	// Note: ClickHouse has no last insert id, so InsertAndGetId inserts the record and fails as
	// not supported.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":       1,
			"passport": "user_1",
			"password": "pass_1",
			"nickname": "name_1",
		}).InsertAndGetId()
		t.AssertNE(err, nil)
		t.Assert(chACount(table, "id", 1), 1)
	})
}

// Test_Model_Increment_Decrement tests Model.Increment and Model.Decrement.
func Test_Model_Increment_Decrement(t *testing.T) {
	// Note: ClickHouse cannot UPDATE a sorting key column, so increasing the key column id fails.
	gtest.C(t, func(t *gtest.T) {
		table := chACreateInitTable()
		defer dropTable(table)
		_, err := db.Model(table).Where("id", 1).Increment("id", 100)
		t.AssertNE(err, nil)
		t.Assert(chACount(table, "id", 1), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		table := chAName("counter")
		_, err := db.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id UInt64, num Int64) ENGINE = MergeTree() ORDER BY id`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)
		_, err = db.Model(table).Data(g.List{{"id": 1, "num": 1}, {"id": 2, "num": 2}}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Where("id", 1).Increment("num", 100)
		t.AssertNil(err)
		t.Assert(chACount(table, "num", 101), 1)

		_, err = db.Model(table).Where("num", 101).Decrement("num", 10)
		t.AssertNil(err)
		t.Assert(chACount(table, "num", 91), 1)
		t.Assert(chACount(table, "num", 2), 1)
	})
}

// Test_Model_OnDuplicate tests Model.OnDuplicate with Save.
func Test_Model_OnDuplicate(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	// Note: ClickHouse has no upsert, so Save inserts the record and OnDuplicate has no effect.
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
		t.Assert(chACount(table, "id", 1), 2)
		t.Assert(chACount(table, "id=1 AND nickname='name_1'"), 1)
		t.Assert(chACount(table, "id=1 AND nickname='n1'"), 1)
	})
}

// Test_Model_OnDuplicateWithCounter tests Model.OnDuplicate with a counter and Save.
func Test_Model_OnDuplicateWithCounter(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	// Note: ClickHouse has no upsert, so Save inserts the record and OnDuplicate has no effect.
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
		t.Assert(chACount(table, "id", 1), 2)
		t.Assert(chACount(table, "id", 1000000), 0)
	})
}

// Test_Model_OnDuplicateEx tests Model.OnDuplicateEx with Save.
func Test_Model_OnDuplicateEx(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	// Note: ClickHouse has no upsert, so Save inserts the record and OnDuplicateEx has no effect.
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
		t.Assert(chACount(table, "id", 1), 2)
		t.Assert(chACount(table, "id=1 AND passport='user_1'"), 1)
		t.Assert(chACount(table, "id=1 AND passport='pp1'"), 1)
	})
}

// Test_Model_Raw_Supplement tests the DB.Raw scenarios of the MySQL suite not covered by Test_Model_Raw.
func Test_Model_Raw_Supplement(t *testing.T) {
	table := chACreateInitTable()
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
		t.Assert(count, int64(3))
	})
}

// Test_Model_Handler tests Model.Handler.
func Test_Model_Handler(t *testing.T) {
	table := chACreateInitTable()
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

// Test_Model_FieldCount tests Model.FieldCount.
func Test_Model_FieldCount(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldCount("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["total"].Int(), 1)
	})
}

// Test_Model_FieldMax tests Model.FieldMax.
func Test_Model_FieldMax(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldMax("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["total"].Int(), 1)
	})
}

// Test_Model_FieldMin tests Model.FieldMin.
func Test_Model_FieldMin(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldMin("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["total"].Int(), 1)
	})
}

// Test_Model_FieldAvg tests Model.FieldAvg.
func Test_Model_FieldAvg(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldAvg("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["total"].Int(), 1)
	})
}

// Test_Model_OmitEmptyWhere tests OmitEmptyWhere with basic, slice, struct and map conditions.
func Test_Model_OmitEmptyWhere(t *testing.T) {
	table := chACreateInitTable()
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
		count, err := db.Model(table).Where(g.Map{
			"id": []int{},
		}).OmitEmptyWhere().Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
}

// Test_Model_GTime_DefaultValue tests inserting and scanning a nil *gtime.Time attribute.
// https://github.com/gogf/gf/issues/1387
func Test_Model_GTime_DefaultValue(t *testing.T) {
	table := chACreateTable()
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
		t.Assert(chACount(table), 2)
	})
}

// Test_Model_Insert_Filter tests that filtering the data does not change the value outside.
func Test_Model_Insert_Filter(t *testing.T) {
	// map
	gtest.C(t, func(t *gtest.T) {
		table := chACreateTable()
		defer dropTable(table)
		data := g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 1)

		t.Assert(data["uid"], 1)
	})
	// slice
	gtest.C(t, func(t *gtest.T) {
		table := chACreateTable()
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

		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 2)

		t.Assert(data[0]["uid"], 1)
		t.Assert(data[1]["uid"], 2)
	})
}

// Test_Model_Embedded_Filter tests inserting a struct with an embedded struct and unknown fields.
func Test_Model_Embedded_Filter(t *testing.T) {
	table := chACreateTable()
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
		_, err := db.Model(table).Data(User{
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
		t.Assert(chACount(table), 1)

		var user *User
		err = db.Model(table).Fields(user).Where("id=100").Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Passport, "john-test")
		t.Assert(user.Id, 100)
	})
}

// Test_Model_Fields_AutoFilterInJoinStatement tests filtering qualified fields in join statements.
func Test_Model_Fields_AutoFilterInJoinStatement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err    error
			table1 = chAName("user")
			table2 = chAName("score")
			table3 = chAName("info")
		)
		if _, err := db.Exec(ctx, fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id   UInt64,
				name String DEFAULT ''
			) ENGINE = MergeTree() ORDER BY id`, table1,
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
				id      UInt64,
				user_id UInt64 DEFAULT 0,
				number  String DEFAULT ''
			) ENGINE = MergeTree() ORDER BY id`, table2,
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
				id          UInt64,
				user_id     UInt64 DEFAULT 0,
				description String DEFAULT ''
			) ENGINE = MergeTree() ORDER BY id`, table3,
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

		// Note: ClickHouse rejects a join condition that does not reference both tables, so the join of
		// table3 on its own columns fails and the following queries join it on table1.
		_, err = db.Model(table1).
			Where(table1+".id", 1).
			Fields(fmt.Sprintf("%s.number,%s.name", table2, table1)).
			LeftJoin(table2, fmt.Sprintf("%s.id=%s.user_id", table1, table2)).
			LeftJoin(table3, fmt.Sprintf("%s.id=%s.user_id", table3, table3)).
			Order(table1 + ".id asc").
			One()
		t.AssertNE(err, nil)

		one, err := db.Model(table1).
			Where(table1+".id", 1).
			Fields(fmt.Sprintf("%s.number,%s.name", table2, table1)).
			LeftJoin(table2, fmt.Sprintf("%s.id=%s.user_id", table1, table2)).
			LeftJoin(table3, fmt.Sprintf("%s.id=%s.user_id", table1, table3)).
			Order(table1 + ".id asc").
			One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["name"].String(), "john")
		t.Assert(one["number"].String(), "n")

		one, err = db.Model(table1).
			LeftJoin(table2, fmt.Sprintf("%s.id=%s.user_id", table1, table2)).
			LeftJoin(table3, fmt.Sprintf("%s.id=%s.user_id", table1, table3)).
			Fields(fmt.Sprintf("%s.number,%s.name", table2, table1)).
			One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["name"].String(), "john")
		t.Assert(one["number"].String(), "n")
	})
}

// Test_ScanList_NoRecreate_PtrAttribute tests that ScanList does not recreate a pointer attribute.
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

// Test_ScanList_NoRecreate_StructAttribute tests that ScanList does not recreate a struct attribute.
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

// Test_ScanList_NoRecreate_SliceAttribute_Ptr tests that ScanList does not recreate a pointer slice attribute.
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

// Test_ScanList_NoRecreate_SliceAttribute_Struct tests that ScanList does not recreate a struct slice attribute.
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

// TestResult_Structs1 tests Result.Structs with an embedded pointer struct.
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

// Test_Builder_OmitEmptyWhere tests OmitEmptyWhere with a where builder.
func Test_Builder_OmitEmptyWhere(t *testing.T) {
	table := chACreateInitTable()
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

// Test_Scan_Nil_Result_Error tests the errors of scanning an empty result.
func Test_Scan_Nil_Result_Error(t *testing.T) {
	table := chACreateInitTable()
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

// Test_Model_FixGdbJoin tests the statement generated for joins on fields with prefixed fields.
func Test_Model_FixGdbJoin(t *testing.T) {
	var (
		commonResource  = chAName("common_resource")
		managedResource = chAName("managed_resource")
		rulesTemplate   = chAName("rules_template")
		resourceMark    = chAName("resource_mark")
	)
	for _, v := range []string{
		fmt.Sprintf(`CREATE TABLE %s (
			id UInt64, app_id Int64, resource_id String, src_instance_id Nullable(String),
			region Nullable(String), zone Nullable(String), database_kind String, source_type String,
			ip Nullable(String), port Nullable(Int32), vpc_id Nullable(String), subnet_id Nullable(String),
			proxy_ip Nullable(String), proxy_port Nullable(Int32), proxy_id Nullable(Int64),
			proxy_snat_ip Nullable(String), lease_at Nullable(DateTime), uin String
		) ENGINE = MergeTree() ORDER BY id`, commonResource),
		fmt.Sprintf(`INSERT INTO %s VALUES
			(1,1,'2','2','2','3','1','1','1',1,'1','1','1',1,1,'1',NULL,''),
			(3,2,'3','3','3','3','3','3','3',3,'3','3','3',3,3,'3',NULL,''),
			(18,1303697168,'dmc-rgnh9qre','vdb-6b6m3u1u','ap-guangzhou','','vdb','cloud','10.0.1.16',80,'vpc-m3dchft7','subnet-9as3a3z2','9.27.72.189',11131,228476,'169.254.128.5, ','2023-11-08 08:13:04',''),
			(20,1303697168,'dmc-4grzi4jg','tdsqlshard-313spncx','ap-guangzhou','','tdsql','cloud','10.255.0.27',3306,'vpc-407k0e8x','subnet-qhkkk3bo','30.86.239.200',24087,0,'',NULL,'')`,
			commonResource),
		fmt.Sprintf(`CREATE TABLE %s (
			id UInt64, instance_id String, resource_id String, resource_name Nullable(String),
			status String DEFAULT 'valid', status_message Nullable(String), user String, password String,
			pay_mode Nullable(Int8) DEFAULT 0, safe_publication Nullable(Bool) DEFAULT false,
			created_at DateTime DEFAULT now(), updated_at DateTime DEFAULT now(), expired_at Nullable(DateTime),
			deleted Int8 DEFAULT 0, resource_mark_id Nullable(Int32), comments Nullable(String),
			rule_template_id String
		) ENGINE = MergeTree() ORDER BY id`, managedResource),
		fmt.Sprintf(`INSERT INTO %s VALUES
			(1,'2','3','1','1','1','1','1',1,true,'2023-11-06 12:14:21','2023-11-06 12:14:21',NULL,1,1,'1',''),
			(2,'3','2','1','1','1','1','1',1,false,'2023-11-06 12:15:07','2023-11-06 12:15:07',NULL,1,2,'1',''),
			(5,'dmcins-jxy0x75m','dmc-rgnh9qre','erichmao-vdb-test','invalid','The Ip field is required','root','2e39af3d',1,true,'2023-11-08 08:13:20','2023-11-09 05:31:07',NULL,0,11,NULL,'12345'),
			(6,'dmcins-erxms6ya','dmc-4grzi4jg','erichmao-vdb-test','invalid','The Ip field is required','leotaowang','641d846c',1,true,'2023-11-08 22:15:17','2023-11-09 05:31:07',NULL,0,11,NULL,'12345')`,
			managedResource),
		fmt.Sprintf(`CREATE TABLE %s (
			id Int64, app_id Nullable(Int64), name String, database_kind Nullable(String),
			is_default Int8 DEFAULT 0, win_rules Nullable(String), inception_rules Nullable(String),
			auto_exec_rules Nullable(String), order_check_step Nullable(String), template_id String DEFAULT '',
			version Int32 DEFAULT 1, deleted Int8 DEFAULT 0, create_at DateTime DEFAULT now(),
			update_at DateTime DEFAULT now(), is_system Int8 DEFAULT 0, uin Nullable(String),
			subAccountUin Nullable(String)
		) ENGINE = MergeTree() ORDER BY id`, rulesTemplate),
		fmt.Sprintf(`CREATE TABLE %s (
			id Int64, app_id Int64, mark_name String, color String, creator String,
			created_at DateTime DEFAULT now(), updated_at DateTime DEFAULT now()
		) ENGINE = MergeTree() ORDER BY id`, resourceMark),
		fmt.Sprintf(`INSERT INTO %s VALUES (10,1,'test','red','1','2023-11-06 02:45:46','2023-11-06 02:45:46')`,
			resourceMark),
	} {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(commonResource, managedResource, rulesTemplate, resourceMark)
	gtest.C(t, func(t *gtest.T) {
		t.AssertNil(db.GetCore().ClearCacheAll(ctx))
		sqlSlice, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			orm := db.Model(managedResource).Ctx(ctx).
				LeftJoinOnField(commonResource, `resource_id`).
				LeftJoinOnFields(resourceMark, `resource_mark_id`, `=`, `id`).
				LeftJoinOnFields(rulesTemplate, `rule_template_id`, `=`, `template_id`).
				FieldsPrefix(
					managedResource,
					"resource_id", "user", "status", "status_message", "safe_publication", "rule_template_id",
					"created_at", "comments", "expired_at", "resource_mark_id", "instance_id", "resource_name",
					"pay_mode").
				FieldsPrefix(resourceMark, "mark_name", "color").
				FieldsPrefix(rulesTemplate, "name").
				FieldsPrefix(commonResource, `src_instance_id`, "database_kind", "source_type", "ip", "port")
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

		var (
			mr       = managedResource
			cr       = commonResource
			rm       = resourceMark
			rt       = rulesTemplate
			expected = fmt.Sprintf(
				"SELECT %[1]s.resource_id,%[1]s.user,%[1]s.status,%[1]s.status_message,%[1]s.safe_publication,"+
					"%[1]s.rule_template_id,%[1]s.created_at,%[1]s.comments,%[1]s.expired_at,%[1]s.resource_mark_id,"+
					"%[1]s.instance_id,%[1]s.resource_name,%[1]s.pay_mode,%[3]s.mark_name,%[3]s.color,%[4]s.name,"+
					"%[2]s.src_instance_id,%[2]s.database_kind,%[2]s.source_type,%[2]s.ip,%[2]s.port "+
					"FROM %[1]s "+
					"LEFT JOIN %[2]s ON (%[1]s.resource_id=%[2]s.resource_id) "+
					"LEFT JOIN %[3]s ON (%[1]s.resource_mark_id = %[3]s.id) "+
					"LEFT JOIN %[4]s ON (%[1]s.rule_template_id = %[4]s.template_id) "+
					"ORDER BY src_instance_id ASC",
				mr, cr, rm, rt,
			)
		)
		t.Assert(expected, sqlSlice[len(sqlSlice)-1])
	})
}

// Test_Model_Year_Date_Time_DateTime_Timestamp tests writing gtime values into date and time columns.
func Test_Model_Year_Date_Time_DateTime_Timestamp(t *testing.T) {
	table := chAName("date_time_example")
	// Note: ClickHouse 24.11 has no YEAR or TIME types, so the year is kept in a UInt16 column and the
	// time of day in a String column.
	_, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id        UInt64,
			year      Nullable(UInt16),
			date      Nullable(Date),
			time      Nullable(String),
			datetime  Nullable(DateTime),
			timestamp Nullable(DateTime64(6))
		) ENGINE = MergeTree() ORDER BY id`, table,
	))
	gtest.AssertNil(err)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// insert.
		var now = gtime.Now()
		_, err := db.Model(table).Insert(g.Map{
			"id":        1,
			"year":      now.Year(),
			"date":      now,
			"time":      now.Format("H:i:s"),
			"datetime":  now,
			"timestamp": now,
		})
		t.AssertNil(err)
		// select.
		one, err := db.Model(table).One()
		t.AssertNil(err)
		t.Assert(one["year"].String(), now.Format("Y"))
		t.Assert(one["date"].String(), now.Format("Y-m-d"))
		t.Assert(one["time"].String(), now.Format("H:i:s"))
		t.AssertLT(one["datetime"].GTime().Sub(now).Seconds(), 5)
		t.AssertLT(one["timestamp"].GTime().Sub(now).Seconds(), 5)
	})
}

// Test_OrderBy_Statement_Generated tests the ORDER BY statement generated for several columns.
func Test_OrderBy_Statement_Generated(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := chAName("employee")
		for _, v := range []string{
			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
				id   UInt64,
				name String,
				age  Int32
			) ENGINE = MergeTree() ORDER BY id`, table),
			fmt.Sprintf(`INSERT INTO %s(id, name, age) VALUES (1, 'John', 30)`, table),
			fmt.Sprintf(`INSERT INTO %s(id, name, age) VALUES (2, 'Mary', 28)`, table),
		} {
			if _, err := db.Exec(ctx, v); err != nil {
				gtest.Error(err)
			}
		}
		defer dropTable(table)
		sqlArray, _ := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err := db.Ctx(ctx).Model(table).Order("name asc", "age desc").All()
			return err
		})
		rawSql := strings.ReplaceAll(sqlArray[len(sqlArray)-1], " ", "")
		expectSql := strings.ReplaceAll(fmt.Sprintf("SELECT * FROM %s ORDER BY name asc, age desc", table), " ", "")
		t.Assert(rawSql, expectSql)
	})
}

// Test_Fields_Raw tests Fields with raw expressions.
func Test_Fields_Raw(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := chACreateInitTable()
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

// Test_ClickHouse_Model_Insert_DifferentKeys tests batch inserting records that have different keys.
func Test_ClickHouse_Model_Insert_DifferentKeys(t *testing.T) {
	table := chACreateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"id": 1, "passport": "user_1"},
			{"id": 2, "passport": "user_2", "nickname": "name_2"},
		}).Insert()
		t.AssertNil(err)
		t.Assert(chACount(table), 2)

		one, err := db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_2")
		t.Assert(one["nickname"], "name_2")
	})
}
