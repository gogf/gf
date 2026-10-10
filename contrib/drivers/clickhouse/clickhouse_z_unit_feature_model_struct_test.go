// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

// Test_Model_Embedded_Insert tests inserting a struct with an embedded struct.
func Test_Model_Embedded_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type Base struct {
			Id         int    `json:"id"`
			Uid        int    `json:"uid"`
			CreateTime string `json:"create_time"`
		}
		type User struct {
			Base
			Passport string `json:"passport"`
			Password string `json:"password"`
			Nickname string `json:"nickname"`
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
		n, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, 1)
		value, err := db.Model(table).Fields("passport").Where("id=100").Value()
		t.AssertNil(err)
		t.Assert(value.String(), "john-test")
	})
}

// Test_Model_Embedded_MapToStruct tests converting a record to a struct with nested embedded structs.
func Test_Model_Embedded_MapToStruct(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type Ids struct {
			Id  int `json:"id"`
			Uid int `json:"uid"`
		}
		type Base struct {
			Ids
			CreateTime string `json:"create_time"`
		}
		type User struct {
			Base
			Passport string `json:"passport"`
			Password string `json:"password"`
			Nickname string `json:"nickname"`
		}
		data := g.Map{
			"id":          100,
			"uid":         101,
			"passport":    "t1",
			"password":    "123456",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id=100").One()
		t.AssertNil(err)

		user := new(User)

		t.Assert(one.Struct(user), nil)
		t.Assert(user.Id, data["id"])
		t.Assert(user.Passport, data["passport"])
		t.Assert(user.Password, data["password"])
		t.Assert(user.Nickname, data["nickname"])
		t.Assert(user.CreateTime, data["create_time"])
	})
}

// Test_Struct_Pointer_Attribute tests converting a record to a struct with pointer attributes.
func Test_Struct_Pointer_Attribute(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       *int
		Passport *string
		Password *string
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		user := new(User)
		err = one.Struct(user)
		t.AssertNil(err)
		t.Assert(*user.Id, 1)
		t.Assert(*user.Passport, "user_1")
		t.Assert(*user.Password, "pass_1")
		t.Assert(user.Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		err := db.Model(table).Scan(user, "id=1")
		t.AssertNil(err)
		t.Assert(*user.Id, 1)
		t.Assert(*user.Passport, "user_1")
		t.Assert(*user.Password, "pass_1")
		t.Assert(user.Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(table).Scan(&user, "id=1")
		t.AssertNil(err)
		t.Assert(*user.Id, 1)
		t.Assert(*user.Passport, "user_1")
		t.Assert(*user.Password, "pass_1")
		t.Assert(user.Nickname, "name_1")
	})
}

// Test_Structs_Pointer_Attribute tests converting records to structs with pointer attributes.
func Test_Structs_Pointer_Attribute(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       *int
		Passport *string
		Password *string
		Nickname string
	}
	// All
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).All("id < 3")
		t.AssertNil(err)
		users := make([]User, 0)
		err = one.Structs(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).All("id < 3")
		t.AssertNil(err)
		users := make([]*User, 0)
		err = one.Structs(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		one, err := db.Model(table).All("id < 3")
		t.AssertNil(err)
		err = one.Structs(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		one, err := db.Model(table).All("id < 3")
		t.AssertNil(err)
		err = one.Structs(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	// Structs
	gtest.C(t, func(t *gtest.T) {
		users := make([]User, 0)
		err := db.Model(table).Scan(&users, "id < 3")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		users := make([]*User, 0)
		err := db.Model(table).Scan(&users, "id < 3")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.Model(table).Scan(&users, "id < 3")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		err := db.Model(table).Scan(&users, "id < 3")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
}

// Test_Struct_Empty tests scanning an empty result to a struct.
func Test_Struct_Empty(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport string
		Password string
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		err := db.Model(table).Where("id=100").Scan(user)
		t.Assert(err, sql.ErrNoRows)
		t.AssertNE(user, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where("id=100").One()
		t.AssertNil(err)
		var user *User
		t.Assert(one.Struct(&user), nil)
		t.Assert(user, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(table).Where("id=100").Scan(&user)
		t.AssertNil(err)
		t.Assert(user, nil)
	})
}

// Test_Structs_Empty tests scanning an empty result to struct slices.
func Test_Structs_Empty(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport string
		Password string
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id>100").All()
		t.AssertNil(err)
		users := make([]User, 0)
		t.Assert(all.Structs(&users), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id>100").All()
		t.AssertNil(err)
		users := make([]User, 10)
		t.Assert(all.Structs(&users), sql.ErrNoRows)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id>100").All()
		t.AssertNil(err)
		var users []User
		t.Assert(all.Structs(&users), nil)
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id>100").All()
		t.AssertNil(err)
		users := make([]*User, 0)
		t.Assert(all.Structs(&users), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id>100").All()
		t.AssertNil(err)
		users := make([]*User, 10)
		t.Assert(all.Structs(&users), sql.ErrNoRows)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id>100").All()
		t.AssertNil(err)
		var users []*User
		t.Assert(all.Structs(&users), nil)
	})
}

const chDCreateTime = "2018-10-24 10:00:00"

// chDCreateInitTableWithCreateTime creates a table with TableSize records like the MySQL init table, all
// having the create_time chDCreateTime.
func chDCreateInitTableWithCreateTime() string {
	name := createTable()
	array := garray.New(true)
	for i := 1; i <= TableSize; i++ {
		array.Append(g.Map{
			"id":          uint64(i),
			"passport":    fmt.Sprintf(`user_%d`, i),
			"password":    fmt.Sprintf(`pass_%d`, i),
			"nickname":    fmt.Sprintf(`name_%d`, i),
			"create_time": gtime.NewFromStr(chDCreateTime).String(),
		})
	}
	if _, err := db.Insert(ctx, name, array.Slice()); err != nil {
		gtest.Fatal(err)
	}
	return name
}

type chDMyTime struct {
	gtime.Time
}

type chDMyTimeSt struct {
	CreateTime chDMyTime
}

func (st *chDMyTimeSt) UnmarshalValue(v any) error {
	m := gconv.Map(v)
	t, err := gtime.StrToTime(gconv.String(m["create_time"]))
	if err != nil {
		return err
	}
	st.CreateTime = chDMyTime{*t}
	return nil
}

// Test_Model_Scan_CustomType_Time tests scanning to a struct with a custom time type.
func Test_Model_Scan_CustomType_Time(t *testing.T) {
	table := chDCreateInitTableWithCreateTime()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		st := new(chDMyTimeSt)
		err := db.Model(table).Fields("create_time").Scan(st)
		t.AssertNil(err)
		t.Assert(st.CreateTime.String(), "2018-10-24 10:00:00")
	})
	gtest.C(t, func(t *gtest.T) {
		var stSlice []*chDMyTimeSt
		err := db.Model(table).Fields("create_time").Scan(&stSlice)
		t.AssertNil(err)
		t.Assert(len(stSlice), TableSize)
		t.Assert(stSlice[0].CreateTime.String(), "2018-10-24 10:00:00")
		t.Assert(stSlice[9].CreateTime.String(), "2018-10-24 10:00:00")
	})
}

// Test_Model_Scan_CustomType_String tests scanning to a struct with a custom string type.
func Test_Model_Scan_CustomType_String(t *testing.T) {
	type MyString string

	type MyStringSt struct {
		Passport MyString
	}

	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		st := new(MyStringSt)
		err := db.Model(table).Fields("Passport").WherePri(1).Scan(st)
		t.AssertNil(err)
		t.Assert(st.Passport, "user_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var sts []MyStringSt
		err := db.Model(table).Fields("Passport").Order("id asc").Scan(&sts)
		t.AssertNil(err)
		t.Assert(len(sts), TableSize)
		t.Assert(sts[0].Passport, "user_1")
	})
}

type chDStructUser struct {
	Id         int
	Passport   string
	Password   string
	Nickname   string
	CreateTime *gtime.Time
}

func (user *chDStructUser) UnmarshalValue(value any) error {
	if record, ok := value.(gdb.Record); ok {
		*user = chDStructUser{
			Id:         record["id"].Int(),
			Passport:   record["passport"].String(),
			Password:   "",
			Nickname:   record["nickname"].String(),
			CreateTime: record["create_time"].GTime(),
		}
		return nil
	}
	return gerror.NewCodef(gcode.CodeInvalidParameter, `unsupported value type for UnmarshalValue: %v`, reflect.TypeOf(value))
}

// Test_Model_Scan_UnmarshalValue tests scanning to structs implementing UnmarshalValue.
func Test_Model_Scan_UnmarshalValue(t *testing.T) {
	table := chDCreateInitTableWithCreateTime()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		var users []*chDStructUser
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[0].Passport, "user_1")
		t.Assert(users[0].Password, "")
		t.Assert(users[0].Nickname, "name_1")
		t.Assert(users[0].CreateTime.String(), chDCreateTime)

		t.Assert(users[9].Id, 10)
		t.Assert(users[9].Passport, "user_10")
		t.Assert(users[9].Password, "")
		t.Assert(users[9].Nickname, "name_10")
		t.Assert(users[9].CreateTime.String(), chDCreateTime)
	})
}

// Test_Model_Scan_Map tests scanning records to structs implementing UnmarshalValue.
func Test_Model_Scan_Map(t *testing.T) {
	table := chDCreateInitTableWithCreateTime()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var users []*chDStructUser
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[0].Passport, "user_1")
		t.Assert(users[0].Password, "")
		t.Assert(users[0].Nickname, "name_1")
		t.Assert(users[0].CreateTime.String(), chDCreateTime)

		t.Assert(users[9].Id, 10)
		t.Assert(users[9].Passport, "user_10")
		t.Assert(users[9].Password, "")
		t.Assert(users[9].Nickname, "name_10")
		t.Assert(users[9].CreateTime.String(), chDCreateTime)
	})
}

// Test_Scan_AutoFilteringByStructAttributes tests that Scan selects only the fields of the struct.
func Test_Scan_AutoFilteringByStructAttributes(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport string
	}
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(table).OrderAsc("id").Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.Model(table).OrderAsc("id").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
	})
}
