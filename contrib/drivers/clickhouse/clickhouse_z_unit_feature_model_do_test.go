// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

// chDCreateNullableTable creates a table like the MySQL test table, whose columns except id are nullable.
func chDCreateNullableTable(table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
	}
	dropTable(name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		   id UInt64 NOT NULL,
		   passport Nullable(String),
		   password Nullable(String),
		   nickname Nullable(String),
		   create_time Nullable(DateTime)
		) ENGINE = MergeTree()
		ORDER BY id`, name,
	)); err != nil {
		gtest.Fatal(err)
	}
	return
}

// Test_Model_Insert_Data_DO tests inserting a DO struct, whose nil fields are not inserted.
func Test_Model_Insert_Data_DO(t *testing.T) {
	table := chDCreateNullableTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := User{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
		}
		// Note: ClickHouse has no last insert id, so the inserted rows are verified by reading them back.
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, 1)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], `pass_1`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)
	})
}

// Test_Model_Insert_Data_List_DO tests inserting a list of DO structs.
func Test_Model_Insert_Data_List_DO(t *testing.T) {
	table := chDCreateNullableTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := g.Slice{
			User{
				Id:       1,
				Passport: "user_1",
				Password: "pass_1",
			},
			User{
				Id:       2,
				Passport: "user_2",
				Password: "pass_2",
			},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, 2)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], `pass_1`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)

		one, err = db.Model(table).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `2`)
		t.Assert(one[`passport`], `user_2`)
		t.Assert(one[`password`], `pass_2`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)
	})
}

// Test_Model_Update_Data_DO tests updating with a DO struct, whose nil fields are not updated.
func Test_Model_Update_Data_DO(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := User{
			Id:       1,
			Passport: "user_100",
			Password: "pass_100",
		}
		// Note: ClickHouse cannot UPDATE a column of the sorting key, so setting id fails.
		_, err := db.Model(table).Data(data).WherePri(1).Update()
		t.AssertNE(err, nil)

		data.Id = nil
		_, err = db.Model(table).Data(data).WherePri(1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_100`)
		t.Assert(one[`password`], `pass_100`)
		t.Assert(one[`nickname`], `name_1`)

		one, err = db.Model(table).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one[`passport`], `user_2`)
	})
}

// Test_Model_Update_Pointer_Data_DO tests updating with a DO struct holding pointer values.
func Test_Model_Update_Pointer_Data_DO(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type NN string
		type Req struct {
			Id       int
			Passport *string
			Password *string
			Nickname *NN
		}
		type UserDo struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		var (
			nickname = NN("nickname_111")
			req      = Req{
				Password: gconv.PtrString("12345678"),
				Nickname: &nickname,
			}
			data = UserDo{
				Passport: req.Passport,
				Password: req.Password,
				Nickname: req.Nickname,
			}
		)

		_, err := db.Model(table).Data(data).WherePri(1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`password`], `12345678`)
		t.Assert(one[`nickname`], `nickname_111`)
	})
}

// Test_Model_Where_DO tests using a DO struct as the where condition.
func Test_Model_Where_DO(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		where := User{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
		}
		one, err := db.Model(table).Where(where).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], `pass_1`)
		t.Assert(one[`nickname`], `name_1`)
	})
}

// Test_Model_Insert_Data_ForDao tests inserting a dao struct of any fields, whose nil fields are not inserted.
func Test_Model_Insert_Data_ForDao(t *testing.T) {
	table := chDCreateNullableTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type UserForDao struct {
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := UserForDao{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, 1)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], `pass_1`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)
	})
}

// Test_Model_Insert_Data_List_ForDao tests inserting a list of dao structs of any fields.
func Test_Model_Insert_Data_List_ForDao(t *testing.T) {
	table := chDCreateNullableTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type UserForDao struct {
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := g.Slice{
			UserForDao{
				Id:       1,
				Passport: "user_1",
				Password: "pass_1",
			},
			UserForDao{
				Id:       2,
				Passport: "user_2",
				Password: "pass_2",
			},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, 2)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], `pass_1`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)

		one, err = db.Model(table).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `2`)
		t.Assert(one[`passport`], `user_2`)
		t.Assert(one[`password`], `pass_2`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)
	})
}

// Test_Model_Update_Data_ForDao tests updating with a dao struct of any fields.
func Test_Model_Update_Data_ForDao(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type UserForDao struct {
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := UserForDao{
			Id:       1,
			Passport: "user_100",
			Password: "pass_100",
		}
		// Note: ClickHouse cannot UPDATE a column of the sorting key, so setting id fails.
		_, err := db.Model(table).Data(data).WherePri(1).Update()
		t.AssertNE(err, nil)

		data.Id = nil
		_, err = db.Model(table).Data(data).WherePri(1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_100`)
		t.Assert(one[`password`], `pass_100`)
		t.Assert(one[`nickname`], `name_1`)

		one, err = db.Model(table).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one[`passport`], `user_2`)
	})
}

// Test_Model_Where_ForDao tests using a dao struct of any fields as the where condition.
func Test_Model_Where_ForDao(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type UserForDao struct {
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		where := UserForDao{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
		}
		one, err := db.Model(table).Where(where).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], `pass_1`)
		t.Assert(one[`nickname`], `name_1`)
	})
}

// chDCreateInstanceTable creates the table of the MySQL testdata table_with_prefix.sql under a unique name.
func chDCreateInstanceTable() string {
	name := fmt.Sprintf(`instance_%d`, gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		   f_id Int32 NOT NULL,
		   name String DEFAULT ''
		) ENGINE = MergeTree()
		ORDER BY f_id`, name,
	)); err != nil {
		gtest.Fatal(err)
	}
	if _, err := db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s VALUES (1, 'john')`, name)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

// Test_Model_Where_FieldPrefix tests DO where conditions on columns mapped by orm tags.
func Test_Model_Where_FieldPrefix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := chDCreateInstanceTable()
		defer dropTable(table)

		type Instance struct {
			ID   int `orm:"f_id"`
			Name string
		}

		type InstanceDo struct {
			g.Meta `orm:"table:instance, do:true"`
			ID     any `orm:"f_id"`
		}
		var instance *Instance
		err := db.Model(table).Where(InstanceDo{
			ID: 1,
		}).Scan(&instance)
		t.AssertNil(err)
		t.AssertNE(instance, nil)
		t.Assert(instance.ID, 1)
		t.Assert(instance.Name, "john")
	})
	// With omitempty.
	gtest.C(t, func(t *gtest.T) {
		table := chDCreateInstanceTable()
		defer dropTable(table)

		type Instance struct {
			ID   int `orm:"f_id,omitempty"`
			Name string
		}

		type InstanceDo struct {
			g.Meta `orm:"table:instance, do:true"`
			ID     any `orm:"f_id,omitempty"`
		}
		var instance *Instance
		err := db.Model(table).Where(InstanceDo{
			ID: 1,
		}).Scan(&instance)
		t.AssertNil(err)
		t.AssertNE(instance, nil)
		t.Assert(instance.ID, 1)
		t.Assert(instance.Name, "john")
	})
}
