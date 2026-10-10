// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

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

// chBConfigNode returns the configuration node of the default test database.
func chBConfigNode() gdb.ConfigNode {
	return gdb.ConfigNode{
		Host:  "127.0.0.1",
		Port:  "9000",
		User:  "default",
		Name:  "default",
		Type:  "clickhouse",
		Extra: "mutations_sync=1",
	}
}

// chBCreateNullableTable creates a table shaped like the MySQL test table, whose columns
// except id are nullable and which has the additional create_date column.
func chBCreateNullableTable(table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
	}
	dropTable(name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id          UInt64,
			passport    Nullable(String),
			password    Nullable(String),
			nickname    Nullable(String),
			create_time Nullable(DateTime),
			create_date Nullable(Date)
		) ENGINE = MergeTree()
		ORDER BY id`, name,
	)); err != nil {
		gtest.Fatal(err)
	}
	return
}

// chBCreateNullableInitTable creates a table by chBCreateNullableTable and inserts the same
// records as createInitTable.
func chBCreateNullableInitTable(table ...string) (name string) {
	name = chBCreateNullableTable(table...)
	array := garray.New(true)
	for i := 1; i <= TableSize; i++ {
		array.Append(g.Map{
			"id":          uint64(i),
			"passport":    fmt.Sprintf(`user_%d`, i),
			"password":    fmt.Sprintf(`pass_%d`, i),
			"nickname":    fmt.Sprintf(`name_%d`, i),
			"create_time": gtime.Now(),
		})
	}
	if _, err := db.Insert(ctx, name, array.Slice()); err != nil {
		gtest.Fatal(err)
	}
	return
}

// chBCreateExistsTable creates the subquery table of the WhereExists tests with records uid 1 and 2.
func chBCreateExistsTable() (name string) {
	name = fmt.Sprintf(`table2_%d`, gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id  UInt32,
			uid UInt32 DEFAULT 0
		) ENGINE = MergeTree()
		ORDER BY id`, name,
	)); err != nil {
		gtest.Fatal(err)
	}
	if _, err := db.Model(name).Insert(g.List{
		{"id": uint32(1), "uid": uint32(1)},
		{"id": uint32(2), "uid": uint32(2)},
	}); err != nil {
		gtest.Fatal(err)
	}
	return
}

// Test_DB_Exec_Supplement tests Exec with arguments and with an invalid statement.
func Test_DB_Exec_Supplement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, "SELECT ?", 1)
		t.AssertNil(err)

		_, err = db.Exec(ctx, "ERROR")
		t.AssertNE(err, nil)
	})
}

// Test_DB_Prepare tests Prepare with a SELECT statement.
func Test_DB_Prepare(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Note: clickhouse-go prepares INSERT statements only, as batches.
		_, err := db.Prepare(ctx, "SELECT 100")
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "invalid INSERT query"), true)
	})
}

// Test_DB_Insert_Supplement tests DB.Insert with maps, structs and batches, verifying the stored values.
func Test_DB_Insert_Supplement(t *testing.T) {
	table := chBCreateNullableTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Insert(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
			"create_date": gtime.Date(),
		})
		t.AssertNil(err)

		// normal map
		_, err = db.Insert(ctx, table, g.Map{
			"id":          "2",
			"passport":    "t2",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
			"create_date": gtime.Date(),
		})
		t.AssertNil(err)

		// struct
		type User struct {
			Id         int         `gconv:"id"`
			Passport   string      `json:"passport"`
			Password   string      `gconv:"password"`
			Nickname   string      `gconv:"nickname"`
			CreateTime string      `json:"create_time"`
			CreateDate *gtime.Time `json:"create_date"`
		}
		gTime := gtime.New("2024-10-01 12:01:01")
		timeStr, dateStr := gTime.String(), "2024-10-01 00:00:00"
		_, err = db.Insert(ctx, table, User{
			Id:         3,
			Passport:   "user_3",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "name_3",
			CreateTime: timeStr,
			CreateDate: gTime,
		})
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)

		t.Assert(one["id"].Int(), 3)
		t.Assert(one["passport"].String(), "user_3")
		t.Assert(one["password"].String(), "25d55ad283aa400af464c76d713c07ad")
		t.Assert(one["nickname"].String(), "name_3")
		t.Assert(one["create_time"].GTime().String(), timeStr)
		t.Assert(one["create_date"].GTime().String(), dateStr)

		// *struct
		gTime = gtime.New("2024-10-01 12:01:01")
		timeStr, dateStr = gTime.String(), "2024-10-01 00:00:00"
		_, err = db.Insert(ctx, table, &User{
			Id:         4,
			Passport:   "t4",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "name_4",
			CreateTime: timeStr,
			CreateDate: gTime,
		})
		t.AssertNil(err)

		one, err = db.Model(table).Where("id", 4).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 4)
		t.Assert(one["passport"].String(), "t4")
		t.Assert(one["password"].String(), "25d55ad283aa400af464c76d713c07ad")
		t.Assert(one["nickname"].String(), "name_4")
		t.Assert(one["create_time"].GTime().String(), timeStr)
		t.Assert(one["create_date"].GTime().String(), dateStr)

		// batch with Insert
		gTime = gtime.New("2024-10-01 12:01:01")
		timeStr, dateStr = gTime.String(), "2024-10-01 00:00:00"
		_, err = db.Insert(ctx, table, g.Slice{
			g.Map{
				"id":          200,
				"passport":    "t200",
				"password":    "25d55ad283aa400af464c76d71qw07ad",
				"nickname":    "T200",
				"create_time": timeStr,
				"create_date": gTime,
			},
			g.Map{
				"id":          300,
				"passport":    "t300",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T300",
				"create_time": timeStr,
				"create_date": gTime,
			},
		})
		t.AssertNil(err)

		one, err = db.Model(table).Where("id", 200).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 200)
		t.Assert(one["passport"].String(), "t200")
		t.Assert(one["password"].String(), "25d55ad283aa400af464c76d71qw07ad")
		t.Assert(one["nickname"].String(), "T200")
		t.Assert(one["create_time"].GTime().String(), timeStr)
		t.Assert(one["create_date"].GTime().String(), dateStr)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 6)
	})
}

// Test_DB_Insert_WithStructAndSliceAttribute tests inserting struct and slice values as JSON strings.
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
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["create_time"], data["create_time"])
		t.Assert(one["nickname"], gjson.New(data["nickname"]).MustToJson())
	})
}

// Test_DB_Insert_KeyFieldNameMapping tests mapping struct field names to column names on insert.
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
		t.Assert(one["passport"], data.Passport)
		t.Assert(one["create_time"], data.CreateTime)
		t.Assert(one["nickname"], data.Nickname)
	})
}

// Test_DB_Insert_NilGjson tests inserting empty, nil and empty-object gjson values.
func Test_DB_Insert_NilGjson(t *testing.T) {
	var tableName = "nil" + gtime.TimestampNanoStr()
	_, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id UInt32,
		json_empty_string Nullable(String),
		json_nil Nullable(String),
		json_null Nullable(String)
	) ENGINE = MergeTree()
	ORDER BY id
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

// Test_DB_Update_KeyFieldNameMapping tests mapping struct field names to column names on update.
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
		// Note: ClickHouse cannot update a column of the sorting key, which id is.
		_, err := db.Update(ctx, table, data, "id=1")
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "Cannot UPDATE key column"), true)

		type UserWithoutKey struct {
			Passport   string
			Password   string
			Nickname   string
			CreateTime string
		}
		_, err = db.Update(ctx, table, UserWithoutKey{
			Passport:   data.Passport,
			Password:   data.Password,
			Nickname:   data.Nickname,
			CreateTime: data.CreateTime,
		}, "id=1")
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(one["passport"], data.Passport)
		t.Assert(one["create_time"], data.CreateTime)
		t.Assert(one["nickname"], data.Nickname)
	})
}

// Test_DB_InsertIgnore tests inserting a duplicate primary key and InsertIgnore.
func Test_DB_InsertIgnore(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		// Note: ClickHouse has no unique constraint, so a duplicate primary key is inserted as another record.
		_, err := db.Insert(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)
		count, err := db.Model(table).Where("id", 1).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
	gtest.C(t, func(t *gtest.T) {
		// Note: ClickHouse cannot skip conflicting records, so InsertIgnore is unsupported.
		_, err := db.InsertIgnore(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:InsertIgnore")
	})
}

// Test_DB_BatchInsert tests DB.Insert with lists, slices and records of different fields.
func Test_DB_BatchInsert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		_, err := db.Insert(ctx, table, g.List{
			{
				"id":          2,
				"passport":    "t2",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_2",
				"create_time": gtime.Now().String(),
			},
			{
				"id":          3,
				"passport":    "user_3",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_3",
				"create_time": gtime.Now().String(),
			},
		}, 1)
		t.AssertNil(err)
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		// []any
		_, err := db.Insert(ctx, table, g.Slice{
			g.Map{
				"id":          2,
				"passport":    "t2",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_2",
				"create_time": gtime.Now().String(),
			},
			g.Map{
				"id":          3,
				"passport":    "user_3",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_3",
				"create_time": gtime.Now().String(),
			},
		}, 1)
		t.AssertNil(err)
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})

	// batch insert map
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		_, err := db.Insert(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "p1",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})

	// Batch insert with different fields
	gtest.C(t, func(t *gtest.T) {
		table := chBCreateNullableTable()
		defer dropTable(table)
		_, err := db.Insert(ctx, table, g.List{
			{
				"id":          2,
				"passport":    "t2",
				"password":    "25d55ad283aa400af464c76d713c07ac",
				"create_time": gtime.Now().String(),
			},
			{
				"id":          3,
				"passport":    "user_3",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_3",
				"create_time": gtime.Now().String(),
			},
		}, 1)
		t.AssertNil(err)
		all, err := db.Model(table).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["nickname"], nil)
		t.Assert(all[1]["nickname"], "name_3")
	})
}

// Test_DB_BatchInsert_Struct tests DB.Insert with a struct pointer.
func Test_DB_BatchInsert_Struct(t *testing.T) {
	// batch insert struct
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		type User struct {
			Id         int         `c:"id"`
			Passport   string      `c:"passport"`
			Password   string      `c:"password"`
			NickName   string      `c:"nickname"`
			CreateTime *gtime.Time `c:"create_time"`
		}
		user := &User{
			Id:         1,
			Passport:   "t1",
			Password:   "p1",
			NickName:   "T1",
			CreateTime: gtime.Now(),
		}
		_, err := db.Insert(ctx, table, user)
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "T1")
	})
}

// Test_DB_Save_Supplement tests DB.Save on an existing primary key.
func Test_DB_Save_Supplement(t *testing.T) {
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

		// Note: ClickHouse has no upsert, Save inserts the record besides the existing one.
		count, err := db.Model(table).Where("id", 1).Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		one, err := db.Model(table).Where("id", 1).Where("nickname", "T11").One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["passport"].String(), "t1")
		t.Assert(one["password"].String(), "25d55ad283aa400af464c76d713c07ad")
		t.Assert(one["nickname"].String(), "T11")
		t.Assert(one["create_time"].GTime().String(), timeStr)
	})
}

// Test_DB_Update_Supplement tests DB.Update with string data and condition.
func Test_DB_Update_Supplement(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Update(ctx, table, "password='987654321'", "id=3")
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 3)
		t.Assert(one["passport"].String(), "user_3")
		t.Assert(one["password"].String(), "987654321")
		t.Assert(one["nickname"].String(), "name_3")
	})
}

// Test_DB_GetOne_Supplement tests GetOne with a string condition argument.
func Test_DB_GetOne_Supplement(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		record, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE passport=?", table), "user_1")
		t.AssertNil(err)
		t.Assert(record["nickname"].String(), "name_1")
	})
}

// Test_DB_GetStruct tests GetScan into a struct with value and pointer time fields.
func Test_DB_GetStruct(t *testing.T) {
	table := createInitTable()
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
		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_3")
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
		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_3")
	})
}

// Test_DB_GetStructs tests GetScan into a struct slice with value and pointer time fields.
func Test_DB_GetStructs(t *testing.T) {
	table := createInitTable()
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
		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
		t.AssertNil(err)
		t.Assert(len(users), TableSize-1)
		t.Assert(users[0].Id, 2)
		t.Assert(users[1].Id, 3)
		t.Assert(users[2].Id, 4)
		t.Assert(users[0].NickName, "name_2")
		t.Assert(users[1].NickName, "name_3")
		t.Assert(users[2].NickName, "name_4")
	})

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []User
		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
		t.AssertNil(err)
		t.Assert(len(users), TableSize-1)
		t.Assert(users[0].Id, 2)
		t.Assert(users[1].Id, 3)
		t.Assert(users[2].Id, 4)
		t.Assert(users[0].NickName, "name_2")
		t.Assert(users[1].NickName, "name_3")
		t.Assert(users[2].NickName, "name_4")
	})
}

// Test_DB_GetScan_Supplement tests the GetScan targets not covered by Test_DB_GetScan.
func Test_DB_GetScan_Supplement(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var user *User
		err := db.GetScan(ctx, &user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_3")
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
		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_3")
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
		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
		t.AssertNil(err)
		t.Assert(len(users), TableSize-1)
		t.Assert(users[0].Id, 2)
		t.Assert(users[1].Id, 3)
		t.Assert(users[2].Id, 4)
		t.Assert(users[0].NickName, "name_2")
		t.Assert(users[1].NickName, "name_3")
		t.Assert(users[2].NickName, "name_4")
	})

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []User
		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
		t.AssertNil(err)
		t.Assert(len(users), TableSize-1)
		t.Assert(users[0].Id, 2)
		t.Assert(users[1].Id, 3)
		t.Assert(users[2].Id, 4)
		t.Assert(users[0].NickName, "name_2")
		t.Assert(users[1].NickName, "name_3")
		t.Assert(users[2].NickName, "name_4")
	})
}

// Test_DB_Delete_Supplement tests DB.Delete with a constant condition.
func Test_DB_Delete_Supplement(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Delete(ctx, table, 1)
		t.AssertNil(err)
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_DB_Time tests inserting time.Time and *time.Time values.
func Test_DB_Time(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Insert(ctx, table, g.Map{
			"id":          200,
			"passport":    "t200",
			"password":    "123456",
			"nickname":    "T200",
			"create_time": time.Now(),
		})
		t.AssertNil(err)
		value, err := db.GetValue(ctx, fmt.Sprintf("select `passport` from `%s` where id=?", table), 200)
		t.AssertNil(err)
		t.Assert(value.String(), "t200")
	})

	gtest.C(t, func(t *gtest.T) {
		t1 := time.Now()
		_, err := db.Insert(ctx, table, g.Map{
			"id":          300,
			"passport":    "t300",
			"password":    "123456",
			"nickname":    "T300",
			"create_time": &t1,
		})
		t.AssertNil(err)
		value, err := db.GetValue(ctx, fmt.Sprintf("select `passport` from `%s` where id=?", table), 300)
		t.AssertNil(err)
		t.Assert(value.String(), "t300")
	})

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
		_, err = db.Delete(ctx, table, 1)
		t.AssertNil(err)
		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_DB_ToJson tests converting Result and Record to structs and JSON.
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

		t.Assert(users[0].Id, resultJson.Get("0.id").Int())
		t.Assert(users[0].Passport, resultJson.Get("0.passport").String())
		t.Assert(users[0].Password, resultJson.Get("0.password").String())
		t.Assert(users[0].NickName, resultJson.Get("0.nickname").String())
		t.Assert(users[0].CreateTime, resultJson.Get("0.create_time").String())

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

// Test_DB_ToXml tests converting a Record to XML.
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
		if v, ok := resultXml["id"]; ok {
			t.Assert(user.Id, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["passport"]; ok {
			t.Assert(user.Passport, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["password"]; ok {
			t.Assert(user.Password, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["nickname"]; ok {
			t.Assert(user.NickName, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["create_time"]; ok {
			t.Assert(user.CreateTime, v)
		} else {
			gtest.Fatal("FAIL")
		}
	})
}

// Test_DB_ToStringMap tests Result.MapKeyStr.
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

		resultStringMap := result.MapKeyStr("id")
		t.Assert(t_users[0].Id, resultStringMap[id]["id"])
		t.Assert(t_users[0].Passport, resultStringMap[id]["passport"])
		t.Assert(t_users[0].Password, resultStringMap[id]["password"])
		t.Assert(t_users[0].NickName, resultStringMap[id]["nickname"])
		t.Assert(t_users[0].CreateTime, resultStringMap[id]["create_time"])
	})
}

// Test_DB_ToIntMap tests Result.MapKeyInt.
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

		resultIntMap := result.MapKeyInt("id")
		t.Assert(t_users[0].Id, resultIntMap[id]["id"])
		t.Assert(t_users[0].Passport, resultIntMap[id]["passport"])
		t.Assert(t_users[0].Password, resultIntMap[id]["password"])
		t.Assert(t_users[0].NickName, resultIntMap[id]["nickname"])
		t.Assert(t_users[0].CreateTime, resultIntMap[id]["create_time"])
	})
}

// Test_DB_ToUintMap tests Result.MapKeyUint.
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

		resultUintMap := result.MapKeyUint("id")
		t.Assert(t_users[0].Id, resultUintMap[uint(id)]["id"])
		t.Assert(t_users[0].Passport, resultUintMap[uint(id)]["passport"])
		t.Assert(t_users[0].Password, resultUintMap[uint(id)]["password"])
		t.Assert(t_users[0].NickName, resultUintMap[uint(id)]["nickname"])
		t.Assert(t_users[0].CreateTime, resultUintMap[uint(id)]["create_time"])
	})
}

// Test_DB_ToStringRecord tests Result.RecordKeyStr.
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

		resultStringRecord := result.RecordKeyStr("id")
		t.Assert(t_users[0].Id, resultStringRecord[ids]["id"].Int())
		t.Assert(t_users[0].Passport, resultStringRecord[ids]["passport"].String())
		t.Assert(t_users[0].Password, resultStringRecord[ids]["password"].String())
		t.Assert(t_users[0].NickName, resultStringRecord[ids]["nickname"].String())
		t.Assert(t_users[0].CreateTime, resultStringRecord[ids]["create_time"].String())
	})
}

// Test_DB_ToIntRecord tests Result.RecordKeyInt.
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

		resultIntRecord := result.RecordKeyInt("id")
		t.Assert(t_users[0].Id, resultIntRecord[id]["id"].Int())
		t.Assert(t_users[0].Passport, resultIntRecord[id]["passport"].String())
		t.Assert(t_users[0].Password, resultIntRecord[id]["password"].String())
		t.Assert(t_users[0].NickName, resultIntRecord[id]["nickname"].String())
		t.Assert(t_users[0].CreateTime, resultIntRecord[id]["create_time"].String())
	})
}

// Test_DB_ToUintRecord tests Result.RecordKeyUint.
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

		resultUintRecord := result.RecordKeyUint("id")
		t.Assert(t_users[0].Id, resultUintRecord[uint(id)]["id"].Int())
		t.Assert(t_users[0].Passport, resultUintRecord[uint(id)]["passport"].String())
		t.Assert(t_users[0].Password, resultUintRecord[uint(id)]["password"].String())
		t.Assert(t_users[0].NickName, resultUintRecord[uint(id)]["nickname"].String())
		t.Assert(t_users[0].CreateTime, resultUintRecord[uint(id)]["create_time"].String())
	})
}

// Test_DB_TableField tests inserting and selecting the numeric and string column types.
func Test_DB_TableField(t *testing.T) {
	name := "field_test_" + gtime.TimestampNanoStr()
	dropTable(name)
	defer dropTable(name)
	_, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		field_tinyint   Nullable(Int8),
		field_int       Nullable(Int32),
		field_integer   Nullable(Int32),
		field_bigint    Nullable(Int64),
		field_bit       Nullable(UInt8),
		field_real      Nullable(Float32),
		field_double    Nullable(Float64),
		field_varchar   Nullable(String),
		field_varbinary Nullable(String)
	) ENGINE = MergeTree()
	ORDER BY tuple()
	`, name))
	if err != nil {
		gtest.Fatal(err)
	}

	data := gdb.Map{
		"field_tinyint":   1,
		"field_int":       2,
		"field_integer":   3,
		"field_bigint":    4,
		"field_bit":       6,
		"field_real":      123,
		"field_double":    123.25,
		"field_varchar":   "abc",
		"field_varbinary": "aaa",
	}
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(name).Data(data).Insert()
		t.AssertNil(err)

		result, err := db.Model(name).Fields("*").Where("field_int = ?", 2).All()
		if err != nil {
			t.Fatal(err)
		}
		t.Assert(result[0], data)
	})
}

// Test_DB_Prefix tests the table prefix configuration with insert, replace, save, update and delete.
func Test_DB_Prefix(t *testing.T) {
	const prefix = "gf_"
	node := chBConfigNode()
	node.Prefix = prefix
	db, err := gdb.New(node)
	gtest.AssertNil(err)
	defer db.Close(ctx)

	name := fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
	table := prefix + name
	createTableWithDb(db, table)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		_, err := db.Insert(ctx, name, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:00").String(),
		})
		t.AssertNil(err)

		count, err := db.Model(name).Where("id", id).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		// Note: ClickHouse has no REPLACE statement, so Replace is unsupported.
		_, err := db.Replace(ctx, name, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:01").String(),
		})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:Replace")
	})

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		_, err := db.Save(ctx, name, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:02").String(),
		})
		t.AssertNil(err)

		// Note: ClickHouse has no upsert, Save inserts the record besides the existing one.
		count, err := db.Model(name).Where("id", id).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		_, err := db.Update(ctx, name, g.Map{
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:03").String(),
		}, "id=?", id)
		t.AssertNil(err)

		count, err := db.Model(name).Where("create_time", "2018-10-24 10:00:03").Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		_, err := db.Delete(ctx, name, "id=?", id)
		t.AssertNil(err)

		count, err := db.Model(name).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
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

		_, err := db.Insert(ctx, name, array.Slice())
		t.AssertNil(err)

		count, err := db.Model(name).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

// Test_Model_InnerJoin tests InnerJoin.
func Test_Model_InnerJoin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createInitTable("user1_" + gtime.TimestampNanoStr())
		table2 := createInitTable("user2_" + gtime.TimestampNanoStr())

		defer dropTable(table1)
		defer dropTable(table2)

		_, err := db.Model(table1).Where("id > ?", 5).Delete()
		t.AssertNil(err)

		count, err := db.Model(table1).Count()
		t.AssertNil(err)
		t.Assert(count, 5)

		result, err := db.Model(table1+" u1").InnerJoin(table2+" u2", "u1.id = u2.id").Order("u1.id").All()
		t.AssertNil(err)
		t.Assert(len(result), 5)

		result, err = db.Model(table1+" u1").InnerJoin(table2+" u2", "u1.id = u2.id").Where("u1.id > ?", 1).Order("u1.id").All()
		t.AssertNil(err)
		t.Assert(len(result), 4)
	})
}

// Test_Model_LeftJoin tests LeftJoin.
func Test_Model_LeftJoin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createInitTable("user1_" + gtime.TimestampNanoStr())
		table2 := createInitTable("user2_" + gtime.TimestampNanoStr())

		defer dropTable(table1)
		defer dropTable(table2)

		_, err := db.Model(table2).Where("id > ?", 3).Delete()
		t.AssertNil(err)

		count, err := db.Model(table2).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		result, err := db.Model(table1+" u1").LeftJoin(table2+" u2", "u1.id = u2.id").All()
		t.AssertNil(err)
		t.Assert(len(result), 10)

		result, err = db.Model(table1+" u1").LeftJoin(table2+" u2", "u1.id = u2.id").Where("u1.id > ? ", 2).All()
		t.AssertNil(err)
		t.Assert(len(result), 8)
	})
}

// Test_Model_RightJoin tests RightJoin.
func Test_Model_RightJoin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createInitTable("user1_" + gtime.TimestampNanoStr())
		table2 := createInitTable("user2_" + gtime.TimestampNanoStr())

		defer dropTable(table1)
		defer dropTable(table2)

		_, err := db.Model(table1).Where("id > ?", 3).Delete()
		t.AssertNil(err)

		count, err := db.Model(table1).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		result, err := db.Model(table1+" u1").RightJoin(table2+" u2", "u1.id = u2.id").All()
		t.AssertNil(err)
		t.Assert(len(result), 10)

		result, err = db.Model(table1+" u1").RightJoin(table2+" u2", "u1.id = u2.id").Where("u1.id > 2").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
	})
}

// Test_Empty_Slice_Argument tests an empty slice argument for an IN condition.
func Test_Empty_Slice_Argument(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf(`select * from %s where id in(?)`, table), g.Slice{})
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

// Test_DB_UpdateCounter tests updating with gdb.Counter.
func Test_DB_UpdateCounter(t *testing.T) {
	tableName := "gf_update_counter_test_" + gtime.TimestampNanoStr()
	_, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		id UInt32,
		views UInt32 DEFAULT 0,
		updated_time UInt32 DEFAULT 0
	) ENGINE = MergeTree()
	ORDER BY id
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
		_, err := db.Update(ctx, tableName, updateData, "id", 1)
		t.AssertNil(err)
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
		_, err := db.Update(ctx, tableName, updateData, "id", 1)
		t.AssertNil(err)
		one, err := db.Model(tableName).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["views"].Int(), 1)
	})
}

// Test_DB_Ctx tests that a query is interrupted by the context deadline.
func Test_DB_Ctx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// Note: ClickHouse limits sleep to 3 seconds per block, so it sleeps 3 instead of 10 seconds.
		_, err := db.Query(ctx, "SELECT sleep(3)")
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "deadline"), true)
	})
}

// Test_DB_Ctx_Logger tests querying with debug logging and a context value.
func Test_DB_Ctx_Logger(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer db.SetDebug(db.GetDebug())
		db.SetDebug(true)
		ctx := context.WithValue(context.Background(), "Trace-Id", "123456789")
		_, err := db.Query(ctx, "SELECT 1")
		t.AssertNil(err)
	})
}

// Test_Types tests inserting and reading the common column types.
func Test_Types(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := "types_" + gtime.TimestampNanoStr()
		if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id UInt32,
        %s String,
        %s FixedString(8),
        %s Date32,
        %s String,
        %s DateTime64(6),
        %s Decimal(5,2),
        %s Float64,
        %s UInt8,
        %s Bool,
        %s Bool
    ) ENGINE = MergeTree()
    ORDER BY id
    `,
			table,
			"`blob`",
			"`binary`",
			"`date`",
			"`time`",
			"`timestamp`",
			"`decimal`",
			"`double`",
			"`bit`",
			"`tinyint`",
			"`bool`")); err != nil {
			gtest.Error(err)
		}
		defer dropTable(table)
		// Note: Date32 starts from 1900-01-01 and ClickHouse has no time-of-day type,
		// so the date is moved into range and the time is stored as a string.
		// ClickHouse truncates the digits beyond the scale of a Decimal instead of rounding them.
		data := g.Map{
			"id":        1,
			"blob":      "i love gf",
			"binary":    []byte("abcdefgh"),
			"date":      "1900-10-24",
			"time":      "10:00:01",
			"timestamp": "2022-02-14 12:00:01.123456",
			"decimal":   -123.456,
			"double":    -123.456,
			"bit":       2,
			"tinyint":   true,
			"bool":      false,
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["blob"].String(), data["blob"])
		t.Assert(one["binary"].String(), data["binary"])
		t.Assert(one["date"].String(), data["date"])
		t.Assert(one["time"].String(), `10:00:01`)
		t.Assert(one["timestamp"].GTime().Format(`Y-m-d H:i:s.u`), `2022-02-14 12:00:01.123`)
		t.Assert(one["decimal"].String(), -123.45)
		t.Assert(one["double"].String(), data["double"])
		t.Assert(one["bit"].Int(), data["bit"])
		t.Assert(one["tinyint"].Bool(), data["tinyint"])

		type T struct {
			Id        int
			Blob      []byte
			Binary    []byte
			Date      *gtime.Time
			Time      *gtime.Time
			Timestamp *gtime.Time
			Decimal   float64
			Double    float64
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
		t.Assert(obj.Decimal, -123.45)
		t.Assert(obj.Double, data["double"])
		t.Assert(obj.Bit, data["bit"])
		t.Assert(obj.TinyInt, data["tinyint"])
	})
}

// Test_Core_ClearTableFields tests TableFields and ClearTableFields.
func Test_Core_ClearTableFields(t *testing.T) {
	table := chBCreateNullableTable()
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

// Test_Core_ClearTableFieldsAll tests ClearTableFieldsAll.
func Test_Core_ClearTableFieldsAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.GetCore().ClearTableFieldsAll(ctx)
		t.AssertNil(err)
	})
}

// Test_Core_ClearCache tests ClearCache.
func Test_Core_ClearCache(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.GetCore().ClearCache(ctx, "")
		t.AssertNil(err)
	})
}

// Test_Core_ClearCacheAll tests ClearCacheAll.
func Test_Core_ClearCacheAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.GetCore().ClearCacheAll(ctx)
		t.AssertNil(err)
	})
}
