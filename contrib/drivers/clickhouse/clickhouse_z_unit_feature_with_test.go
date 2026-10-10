// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
)

const (
	chDWithTplUser = `
CREATE TABLE IF NOT EXISTS %s (
    id UInt32 NOT NULL,
    name String NOT NULL
) ENGINE = MergeTree()
ORDER BY id`
	chDWithTplUserDetail = `
CREATE TABLE IF NOT EXISTS %s (
    uid UInt32 NOT NULL,
    address String NOT NULL
) ENGINE = MergeTree()
ORDER BY uid`
	chDWithTplUserScores = `
CREATE TABLE IF NOT EXISTS %s (
    id UInt32 NOT NULL,
    uid UInt32 NOT NULL,
    score UInt32 NOT NULL
) ENGINE = MergeTree()
ORDER BY id`
)

// chDWithMultipleDepends is the ClickHouse version of the MySQL testdata with_multiple_depends.sql.
var chDWithMultipleDepends = []string{
	"CREATE TABLE table_a (id Int32 NOT NULL, alias String DEFAULT '') ENGINE = MergeTree() ORDER BY id",
	"INSERT INTO table_a VALUES (1, 'table_a_test1')",
	"INSERT INTO table_a VALUES (2, 'table_a_test2')",
	"CREATE TABLE table_b (id Int32 NOT NULL, table_a_id Int32 NOT NULL, alias String DEFAULT '') ENGINE = MergeTree() ORDER BY id",
	"INSERT INTO table_b VALUES (10, 1, 'table_b_test1')",
	"INSERT INTO table_b VALUES (20, 2, 'table_b_test2')",
	"INSERT INTO table_b VALUES (30, 1, 'table_b_test3')",
	"INSERT INTO table_b VALUES (40, 2, 'table_b_test4')",
	"CREATE TABLE table_c (id Int32 NOT NULL, table_b_id Int32 NOT NULL, alias String DEFAULT '') ENGINE = MergeTree() ORDER BY id",
	"INSERT INTO table_c VALUES (100, 10, 'table_c_test1')",
	"INSERT INTO table_c VALUES (200, 10, 'table_c_test2')",
	"INSERT INTO table_c VALUES (300, 20, 'table_c_test3')",
	"INSERT INTO table_c VALUES (400, 30, 'table_c_test4')",
}

// chDNewSchemaDb creates a database of a unique name and returns a db using it with the database name.
// The tests of relations use it as the table names are fixed by the struct tags.
func chDNewSchemaDb() (gdb.DB, string) {
	schema := chDCreateDatabase("test_with")
	return chDNewDbByGroup().Schema(schema), schema
}

// chDOptimizeSchema merges the parts of every table of the database of `db`, so that the rows of a table
// are read in the order of its sorting key, as MySQL reads them in the order of the primary key.
func chDOptimizeSchema(db gdb.DB) {
	tables, err := db.GetArray(ctx, "SELECT name FROM system.tables WHERE database = currentDatabase()")
	if err != nil {
		gtest.Fatal(err)
	}
	for _, table := range tables {
		if _, err = db.Exec(ctx, fmt.Sprintf("OPTIMIZE TABLE %s FINAL", table.String())); err != nil {
			gtest.Fatal(err)
		}
	}
}

/*
mysql> show tables;
+----------------+
| Tables_in_test |
+----------------+
| user           |
| user_detail    |
| user_score     |
+----------------+
3 rows in set (0.01 sec)

mysql> select * from `user`;
+----+--------+
| id | name   |
+----+--------+
|  1 | name_1 |
|  2 | name_2 |
|  3 | name_3 |
|  4 | name_4 |
|  5 | name_5 |
+----+--------+
5 rows in set (0.01 sec)

mysql> select * from `user_detail`;
+-----+-----------+
| uid | address   |
+-----+-----------+
|   1 | address_1 |
|   2 | address_2 |
|   3 | address_3 |
|   4 | address_4 |
|   5 | address_5 |
+-----+-----------+
5 rows in set (0.00 sec)

mysql> select * from `user_score`;
+----+-----+-------+
| id | uid | score |
+----+-----+-------+
|  1 |   1 |     1 |
|  2 |   1 |     2 |
|  3 |   1 |     3 |
|  4 |   1 |     4 |
|  5 |   1 |     5 |
|  6 |   2 |     1 |
|  7 |   2 |     2 |
|  8 |   2 |     3 |
|  9 |   2 |     4 |
| 10 |   2 |     5 |
| 11 |   3 |     1 |
| 12 |   3 |     2 |
| 13 |   3 |     3 |
| 14 |   3 |     4 |
| 15 |   3 |     5 |
| 16 |   4 |     1 |
| 17 |   4 |     2 |
| 18 |   4 |     3 |
| 19 |   4 |     4 |
| 20 |   4 |     5 |
| 21 |   5 |     1 |
| 22 |   5 |     2 |
| 23 |   5 |     3 |
| 24 |   5 |     4 |
| 25 |   5 |     5 |
+----+-----+-------+
25 rows in set (0.00 sec)
*/

// Test_Table_Relation_With_Scan tests With on a model scanning to a struct.
func Test_Table_Relation_With_Scan(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_score"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUser, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserDetail, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserScores, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScore struct {
		gmeta.Meta `orm:"table:user_score"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user"`
		Id         int          `json:"id"`
		Name       string       `json:"name"`
		UserDetail *UserDetail  `orm:"with:uid=id"`
		UserScores []*UserScore `orm:"with:uid=id"`
	}

	// Note: ClickHouse has no auto increment and reports no last insert id, so InsertAndGetId inserts
	// the record, id 0, and fails, and the ids are given below.
	gtest.C(t, func(t *gtest.T) {
		user := User{
			Name: "name_0",
		}
		_, err := db.Model(user).Data(user).OmitEmpty().InsertAndGetId()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(tableUser).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
		_, err = db.Model(tableUser).Delete("1=1")
		t.AssertNil(err)
	})

	// Initialize the data.
	gtest.C(t, func(t *gtest.T) {
		for i := 1; i <= 5; i++ {
			// User.
			user := User{
				Id:   i,
				Name: fmt.Sprintf(`name_%d`, i),
			}
			_, err := db.Model(user).Data(user).OmitEmpty().Insert()
			t.AssertNil(err)
			// Detail.
			userDetail := UserDetail{
				Uid:     i,
				Address: fmt.Sprintf(`address_%d`, i),
			}
			_, err = db.Model(userDetail).Data(userDetail).OmitEmpty().Insert()
			t.AssertNil(err)
			// Scores.
			for j := 1; j <= 5; j++ {
				userScore := UserScore{
					Id:    (i-1)*5 + j,
					Uid:   i,
					Score: j,
				}
				_, err = db.Model(userScore).Data(userScore).OmitEmpty().Insert()
				t.AssertNil(err)
			}
		}
	})
	for i := 6; i <= 10; i++ {
		// User.
		user := User{
			Id:   i,
			Name: fmt.Sprintf(`name_%d`, i-5),
		}
		_, err := db.Model(user).Data(user).OmitEmpty().Insert()
		gtest.AssertNil(err)
		// Detail.
		userDetail := UserDetail{
			Uid:     i,
			Address: fmt.Sprintf(`address_%d`, i),
		}
		_, err = db.Model(userDetail).Data(userDetail).Insert()
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			userScore := UserScore{
				Id:    (i-1)*5 + j,
				Uid:   i,
				Score: j,
			}
			_, err = db.Model(userScore).Data(userScore).Insert()
			gtest.AssertNil(err)
		}
	}
	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.With(User{}).
			With(User{}.UserDetail).
			With(User{}.UserScores).
			Where("id", 3).
			Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].Uid, 3)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].Uid, 3)
		t.Assert(user.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.With(user).
			With(user.UserDetail).
			With(user.UserScores).
			Where("id", 4).
			Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].Uid, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].Uid, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.With(User{}).
			With(UserDetail{}).
			With(UserScore{}).
			Where("id", 4).
			Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].Uid, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].Uid, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
	// With part attribute: UserDetail.
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.With(user).
			With(user.UserDetail).
			Where("id", 4).
			Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 0)
	})
	// With part attribute: UserScores.
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.With(user).
			With(user.UserScores).
			Where("id", 4).
			Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.Assert(user.UserDetail, nil)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].Uid, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].Uid, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
}

// Test_Table_Relation_With tests With on a model scanning to struct slices.
func Test_Table_Relation_With(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUser, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserDetail, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserScores, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:uid=id"`
		UserScores []*UserScores `orm:"with:uid=id"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"id":    (i-1)*5 + j,
				"uid":   i,
				"score": j,
			})
			gtest.AssertNil(err)
		}
	}

	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		err := db.With(User{}).
			With(User{}.UserDetail).
			With(User{}.UserScores).
			Where("id", []int{3, 4}).
			Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, 3)
		t.Assert(users[0].Name, "name_3")
		t.AssertNE(users[0].UserDetail, nil)
		t.Assert(users[0].UserDetail.Uid, 3)
		t.Assert(users[0].UserDetail.Address, "address_3")
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(users[0].UserScores[0].Uid, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Uid, 3)
		t.Assert(users[0].UserScores[4].Score, 5)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.Uid, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Uid, 4)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.With(User{}).
			With(User{}.UserDetail).
			With(User{}.UserScores).
			Where("id", []int{3, 4}).
			Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, 3)
		t.Assert(users[0].Name, "name_3")
		t.AssertNE(users[0].UserDetail, nil)
		t.Assert(users[0].UserDetail.Uid, 3)
		t.Assert(users[0].UserDetail.Address, "address_3")
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(users[0].UserScores[0].Uid, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Uid, 3)
		t.Assert(users[0].UserScores[4].Score, 5)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.Uid, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Uid, 4)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
	// With part attribute: UserDetail.
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		err := db.With(User{}).
			With(User{}.UserDetail).
			Where("id", []int{3, 4}).
			Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, 3)
		t.Assert(users[0].Name, "name_3")
		t.AssertNE(users[0].UserDetail, nil)
		t.Assert(users[0].UserDetail.Uid, 3)
		t.Assert(users[0].UserDetail.Address, "address_3")
		t.Assert(len(users[0].UserScores), 0)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.Uid, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 0)
	})
	// With part attribute: UserScores.
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		err := db.With(User{}).
			With(User{}.UserScores).
			Where("id", []int{3, 4}).
			Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, 3)
		t.Assert(users[0].Name, "name_3")
		t.Assert(users[0].UserDetail, nil)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(users[0].UserScores[0].Uid, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Uid, 3)
		t.Assert(users[0].UserScores[4].Score, 5)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.Assert(users[1].UserDetail, nil)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Uid, 4)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
}

// Test_Table_Relation_WithAll tests WithAll scanning to a struct.
func Test_Table_Relation_WithAll(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUser, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserDetail, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserScores, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:uid=id"`
		UserScores []*UserScores `orm:"with:uid=id"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"id":    (i-1)*5 + j,
				"uid":   i,
				"score": j,
			})
			gtest.AssertNil(err)
		}
	}
	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].Uid, 3)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].Uid, 3)
		t.Assert(user.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].Uid, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].Uid, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
}

// Test_Table_Relation_WithAll_List tests WithAll scanning to struct slices.
func Test_Table_Relation_WithAll_List(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUser, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserDetail, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserScores, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:uid=id"`
		UserScores []*UserScores `orm:"with:uid=id"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"id":    (i-1)*5 + j,
				"uid":   i,
				"score": j,
			})
			gtest.AssertNil(err)
		}
	}
	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		err := db.Model(tableUser).WithAll().Where("id", []int{3, 4}).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, 3)
		t.Assert(users[0].Name, "name_3")
		t.AssertNE(users[0].UserDetail, nil)
		t.Assert(users[0].UserDetail.Uid, 3)
		t.Assert(users[0].UserDetail.Address, "address_3")
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(users[0].UserScores[0].Uid, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Uid, 3)
		t.Assert(users[0].UserScores[4].Score, 5)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.Uid, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Uid, 4)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.Model(tableUser).WithAll().Where("id", []int{3, 4}).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, 3)
		t.Assert(users[0].Name, "name_3")
		t.AssertNE(users[0].UserDetail, nil)
		t.Assert(users[0].UserDetail.Uid, 3)
		t.Assert(users[0].UserDetail.Address, "address_3")
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(users[0].UserScores[0].Uid, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Uid, 3)
		t.Assert(users[0].UserScores[4].Score, 5)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.Uid, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Uid, 4)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
}

// Test_Table_Relation_WithAllCondition_List tests WithAll with where and order conditions in the with tags.
func Test_Table_Relation_WithAllCondition_List(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUser, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserDetail, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserScores, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:uid=id, where:uid > 3"`
		UserScores []*UserScores `orm:"with:uid=id, where:score>1 and score<5, order:score desc"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"id":    (i-1)*5 + j,
				"uid":   i,
				"score": j,
			})
			gtest.AssertNil(err)
		}
	}

	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		err := db.Model(tableUser).WithAll().Where("id", []int{3, 4}).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, 3)
		t.Assert(users[0].Name, "name_3")
		t.Assert(users[0].UserDetail, nil)
		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.Uid, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 3)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 4)
		t.Assert(users[1].UserScores[2].Uid, 4)
		t.Assert(users[1].UserScores[2].Score, 2)
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.Model(tableUser).WithAll().Where("id", []int{3, 4}).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, 3)
		t.Assert(users[0].Name, "name_3")
		t.Assert(users[0].UserDetail, nil)

		t.Assert(len(users[0].UserScores), 3)
		t.Assert(users[0].UserScores[0].Uid, 3)
		t.Assert(users[0].UserScores[0].Score, 4)
		t.Assert(users[0].UserScores[2].Uid, 3)
		t.Assert(users[0].UserScores[2].Score, 2)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.Uid, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 3)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 4)
		t.Assert(users[1].UserScores[2].Uid, 4)
		t.Assert(users[1].UserScores[2].Score, 2)
	})
}

// Test_Table_Relation_WithAll_Embedded_With_SelfMaintained_Attributes tests WithAll on an embedded struct with its own attributes.
func Test_Table_Relation_WithAll_Embedded_With_SelfMaintained_Attributes(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUser, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserDetail, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserScores, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta  `orm:"table:user"`
		*UserDetail `orm:"with:uid=id"`
		Id          int           `json:"id"`
		Name        string        `json:"name"`
		UserScores  []*UserScores `orm:"with:uid=id"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"id":    (i-1)*5 + j,
				"uid":   i,
				"score": j,
			})
			gtest.AssertNil(err)
		}
	}
	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].Uid, 3)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].Uid, 3)
		t.Assert(user.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].Uid, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].Uid, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
}

// Test_Table_Relation_WithAll_Embedded_Without_SelfMaintained_Attributes tests WithAll on an embedded struct without its own attributes.
func Test_Table_Relation_WithAll_Embedded_Without_SelfMaintained_Attributes(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUser, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserDetail, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserScores, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	// For Test Only
	type UserEmbedded struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	type User struct {
		gmeta.Meta  `orm:"table:user"`
		*UserDetail `orm:"with:uid=id"`
		UserEmbedded
		UserScores []*UserScores `orm:"with:uid=id"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"id":    (i-1)*5 + j,
				"uid":   i,
				"score": j,
			})
			gtest.AssertNil(err)
		}
	}
	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].Uid, 3)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].Uid, 3)
		t.Assert(user.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].Uid, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].Uid, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
}

// Test_Table_Relation_WithAll_Embedded_WithoutMeta tests WithAll on an embedded struct without the table meta.
func Test_Table_Relation_WithAll_Embedded_WithoutMeta(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUser, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserDetail, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserScores, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserScores)

	type UserDetailBase struct {
		Uid     int    `json:"uid"`
		Address string `json:"address"`
	}

	type UserDetail struct {
		UserDetailBase
	}

	type UserScores struct {
		Id    int `json:"id"`
		Uid   int `json:"uid"`
		Score int `json:"score"`
	}

	type User struct {
		*UserDetail `orm:"with:uid=id"`
		Id          int           `json:"id"`
		Name        string        `json:"name"`
		UserScores  []*UserScores `orm:"with:uid=id"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"id":    (i-1)*5 + j,
				"uid":   i,
				"score": j,
			})
			gtest.AssertNil(err)
		}
	}
	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].Uid, 3)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].Uid, 3)
		t.Assert(user.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].Uid, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].Uid, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
}

// Test_Table_Relation_WithAll_AttributeStructAlsoHasWithTag tests WithAll on an attribute struct having with tags itself.
func Test_Table_Relation_WithAll_AttributeStructAlsoHasWithTag(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUser, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserDetail, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserScores, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserScores)

	type UserScores struct {
		gmeta.Meta `orm:"table:user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail"`
		Uid        int           `json:"uid"`
		Address    string        `json:"address"`
		UserScores []*UserScores `orm:"with:uid"`
	}

	type User struct {
		gmeta.Meta  `orm:"table:user"`
		*UserDetail `orm:"with:uid=id"`
		Id          int    `json:"id"`
		Name        string `json:"name"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"id":    (i-1)*5 + j,
				"uid":   i,
				"score": j,
			})
			gtest.AssertNil(err)
		}
	}

	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserDetail.UserScores), 5)
		t.Assert(user.UserDetail.UserScores[0].Uid, 3)
		t.Assert(user.UserDetail.UserScores[0].Score, 1)
		t.Assert(user.UserDetail.UserScores[4].Uid, 3)
		t.Assert(user.UserDetail.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserDetail.UserScores), 5)
		t.Assert(user.UserDetail.UserScores[0].Uid, 4)
		t.Assert(user.UserDetail.UserScores[0].Score, 1)
		t.Assert(user.UserDetail.UserScores[4].Uid, 4)
		t.Assert(user.UserDetail.UserScores[4].Score, 5)
	})
}

// Test_Table_Relation_WithAll_AttributeStructAlsoHasWithTag_MoreDeep tests WithAll on nested attribute structs having with tags.
func Test_Table_Relation_WithAll_AttributeStructAlsoHasWithTag_MoreDeep(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUser, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserDetail, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserScores, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserScores)

	type UserScores struct {
		gmeta.Meta `orm:"table:user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type UserDetail1 struct {
		gmeta.Meta `orm:"table:user_detail"`
		Uid        int           `json:"uid"`
		Address    string        `json:"address"`
		UserScores []*UserScores `orm:"with:uid"`
	}

	type UserDetail2 struct {
		gmeta.Meta  `orm:"table:user_detail"`
		Uid         int           `json:"uid"`
		Address     string        `json:"address"`
		UserDetail1 *UserDetail1  `orm:"with:uid"`
		UserScores  []*UserScores `orm:"with:uid"`
	}

	type UserDetail3 struct {
		gmeta.Meta  `orm:"table:user_detail"`
		Uid         int           `json:"uid"`
		Address     string        `json:"address"`
		UserDetail2 *UserDetail2  `orm:"with:uid"`
		UserScores  []*UserScores `orm:"with:uid"`
	}

	type UserDetail struct {
		gmeta.Meta  `orm:"table:user_detail"`
		Uid         int           `json:"uid"`
		Address     string        `json:"address"`
		UserDetail3 *UserDetail3  `orm:"with:uid"`
		UserScores  []*UserScores `orm:"with:uid"`
	}

	type User struct {
		gmeta.Meta  `orm:"table:user"`
		*UserDetail `orm:"with:uid=id"`
		Id          int    `json:"id"`
		Name        string `json:"name"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"id":    (i-1)*5 + j,
				"uid":   i,
				"score": j,
			})
			gtest.AssertNil(err)
		}
	}

	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 3)
		t.Assert(user.UserDetail.UserDetail3.Uid, 3)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.Uid, 3)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.UserDetail1.Uid, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserDetail.UserScores), 5)
		t.Assert(user.UserDetail.UserScores[0].Uid, 3)
		t.Assert(user.UserDetail.UserScores[0].Score, 1)
		t.Assert(user.UserDetail.UserScores[4].Uid, 3)
		t.Assert(user.UserDetail.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 4)
		t.Assert(user.UserDetail.UserDetail3.Uid, 4)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.Uid, 4)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.UserDetail1.Uid, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserDetail.UserScores), 5)
		t.Assert(user.UserDetail.UserScores[0].Uid, 4)
		t.Assert(user.UserDetail.UserScores[0].Score, 1)
		t.Assert(user.UserDetail.UserScores[4].Uid, 4)
		t.Assert(user.UserDetail.UserScores[4].Score, 5)
	})
}

// Test_Table_Relation_With_AttributeStructAlsoHasWithTag_MoreDeep tests With on nested attribute structs having with tags.
func Test_Table_Relation_With_AttributeStructAlsoHasWithTag_MoreDeep(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUser, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserDetail, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(chDWithTplUserScores, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserScores)

	type UserScores struct {
		gmeta.Meta `orm:"table:user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type UserDetail1 struct {
		gmeta.Meta `orm:"table:user_detail"`
		Uid        int           `json:"uid"`
		Address    string        `json:"address"`
		UserScores []*UserScores `orm:"with:uid"`
	}

	type UserDetail2 struct {
		gmeta.Meta  `orm:"table:user_detail"`
		Uid         int           `json:"uid"`
		Address     string        `json:"address"`
		UserDetail1 *UserDetail1  `orm:"with:uid"`
		UserScores  []*UserScores `orm:"with:uid"`
	}

	type UserDetail3 struct {
		gmeta.Meta  `orm:"table:user_detail"`
		Uid         int           `json:"uid"`
		Address     string        `json:"address"`
		UserDetail2 *UserDetail2  `orm:"with:uid"`
		UserScores  []*UserScores `orm:"with:uid"`
	}

	type UserDetail struct {
		gmeta.Meta  `orm:"table:user_detail"`
		Uid         int           `json:"uid"`
		Address     string        `json:"address"`
		UserDetail3 *UserDetail3  `orm:"with:uid"`
		UserScores  []*UserScores `orm:"with:uid"`
	}

	type User struct {
		gmeta.Meta  `orm:"table:user"`
		*UserDetail `orm:"with:uid=id"`
		Id          int    `json:"id"`
		Name        string `json:"name"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"id":    (i-1)*5 + j,
				"uid":   i,
				"score": j,
			})
			gtest.AssertNil(err)
		}
	}

	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).With(UserDetail{}, UserDetail2{}, UserDetail3{}, UserScores{}).Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 3)
		t.Assert(user.UserDetail.UserDetail3.Uid, 3)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.Uid, 3)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.UserDetail1, nil)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserDetail.UserScores), 5)
		t.Assert(user.UserDetail.UserScores[0].Uid, 3)
		t.Assert(user.UserDetail.UserScores[0].Score, 1)
		t.Assert(user.UserDetail.UserScores[4].Uid, 3)
		t.Assert(user.UserDetail.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).With(UserDetail{}, UserDetail2{}, UserDetail3{}, UserScores{}).Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 4)
		t.Assert(user.UserDetail.UserDetail3.Uid, 4)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.Uid, 4)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.UserDetail1, nil)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserDetail.UserScores), 5)
		t.Assert(user.UserDetail.UserScores[0].Uid, 4)
		t.Assert(user.UserDetail.UserScores[0].Score, 1)
		t.Assert(user.UserDetail.UserScores[4].Uid, 4)
		t.Assert(user.UserDetail.UserScores[4].Score, 5)
	})
}

// Test_Table_Relation_With_MultipleDepends1 tests WithAll on a chain of struct attributes.
func Test_Table_Relation_With_MultipleDepends1(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	defer func() {
		dropTableWithDb(db, "table_a")
		dropTableWithDb(db, "table_b")
		dropTableWithDb(db, "table_c")
	}()
	for _, v := range chDWithMultipleDepends {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	chDOptimizeSchema(db)

	type TableC struct {
		gmeta.Meta `orm:"table_c"`
		Id         int `orm:"id,primary" json:"id"`
		TableBId   int `orm:"table_b_id" json:"table_b_id"`
	}

	type TableB struct {
		gmeta.Meta `orm:"table_b"`
		Id         int     `orm:"id,primary" json:"id"`
		TableAId   int     `orm:"table_a_id" json:"table_a_id"`
		TableC     *TableC `orm:"with:table_b_id=id"  json:"table_c"`
	}

	type TableA struct {
		gmeta.Meta `orm:"table_a"`
		Id         int     `orm:"id,primary" json:"id"`
		TableB     *TableB `orm:"with:table_a_id=id" json:"table_b"`
	}

	// Struct.
	gtest.C(t, func(t *gtest.T) {
		var tableA *TableA
		err := db.Model("table_a").WithAll().Scan(&tableA)
		t.AssertNil(err)
		t.AssertNE(tableA, nil)
		t.Assert(tableA.Id, 1)

		t.AssertNE(tableA.TableB, nil)
		t.AssertNE(tableA.TableB.TableC, nil)
		t.Assert(tableA.TableB.TableAId, 1)
		t.Assert(tableA.TableB.TableC.Id, 100)
		t.Assert(tableA.TableB.TableC.TableBId, 10)
	})

	// Structs
	gtest.C(t, func(t *gtest.T) {
		var tableA []*TableA
		err := db.Model("table_a").WithAll().OrderAsc("id").Scan(&tableA)
		t.AssertNil(err)
		t.Assert(len(tableA), 2)
		t.AssertNE(tableA[0].TableB, nil)
		t.AssertNE(tableA[1].TableB, nil)
		t.AssertNE(tableA[0].TableB.TableC, nil)
		t.AssertNE(tableA[1].TableB.TableC, nil)

		t.Assert(tableA[0].Id, 1)
		t.Assert(tableA[0].TableB.Id, 10)
		t.Assert(tableA[0].TableB.TableC.Id, 100)

		t.Assert(tableA[1].Id, 2)
		t.Assert(tableA[1].TableB.Id, 20)
		t.Assert(tableA[1].TableB.TableC.Id, 300)
	})
}

// Test_Table_Relation_With_MultipleDepends2 tests WithAll on a chain of slice attributes.
func Test_Table_Relation_With_MultipleDepends2(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	defer func() {
		dropTableWithDb(db, "table_a")
		dropTableWithDb(db, "table_b")
		dropTableWithDb(db, "table_c")
	}()
	for _, v := range chDWithMultipleDepends {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	chDOptimizeSchema(db)

	type TableC struct {
		gmeta.Meta `orm:"table_c"`
		Id         int `orm:"id,primary" json:"id"`
		TableBId   int `orm:"table_b_id" json:"table_b_id"`
	}

	type TableB struct {
		gmeta.Meta `orm:"table_b"`
		Id         int       `orm:"id,primary" json:"id"`
		TableAId   int       `orm:"table_a_id" json:"table_a_id"`
		TableC     []*TableC `orm:"with:table_b_id=id"  json:"table_c"`
	}

	type TableA struct {
		gmeta.Meta `orm:"table_a"`
		Id         int       `orm:"id,primary" json:"id"`
		TableB     []*TableB `orm:"with:table_a_id=id" json:"table_b"`
	}

	// Struct.
	gtest.C(t, func(t *gtest.T) {
		var tableA *TableA
		err := db.Model("table_a").WithAll().Scan(&tableA)
		t.AssertNil(err)
		t.AssertNE(tableA, nil)
		t.Assert(tableA.Id, 1)

		t.Assert(len(tableA.TableB), 2)
		t.Assert(tableA.TableB[0].Id, 10)
		t.Assert(tableA.TableB[1].Id, 30)

		t.Assert(len(tableA.TableB[0].TableC), 2)
		t.Assert(len(tableA.TableB[1].TableC), 1)
		t.Assert(tableA.TableB[0].TableC[0].Id, 100)
		t.Assert(tableA.TableB[0].TableC[0].TableBId, 10)
		t.Assert(tableA.TableB[0].TableC[1].Id, 200)
		t.Assert(tableA.TableB[0].TableC[1].TableBId, 10)
		t.Assert(tableA.TableB[1].TableC[0].Id, 400)
		t.Assert(tableA.TableB[1].TableC[0].TableBId, 30)
	})

	// Structs
	gtest.C(t, func(t *gtest.T) {
		var tableA []*TableA
		err := db.Model("table_a").WithAll().OrderAsc("id").Scan(&tableA)
		t.AssertNil(err)
		t.Assert(len(tableA), 2)

		t.Assert(len(tableA[0].TableB), 2)
		t.Assert(tableA[0].TableB[0].Id, 10)
		t.Assert(tableA[0].TableB[1].Id, 30)

		t.Assert(len(tableA[0].TableB[0].TableC), 2)
		t.Assert(len(tableA[0].TableB[1].TableC), 1)
		t.Assert(tableA[0].TableB[0].TableC[0].Id, 100)
		t.Assert(tableA[0].TableB[0].TableC[0].TableBId, 10)
		t.Assert(tableA[0].TableB[0].TableC[1].Id, 200)
		t.Assert(tableA[0].TableB[0].TableC[1].TableBId, 10)
		t.Assert(tableA[0].TableB[1].TableC[0].Id, 400)
		t.Assert(tableA[0].TableB[1].TableC[0].TableBId, 30)

		t.Assert(tableA[1].TableB[0].TableC[0].Id, 300)
		t.Assert(tableA[1].TableB[0].TableC[0].TableBId, 20)

		t.Assert(tableA[1].TableB[1].Id, 40)
		t.Assert(tableA[1].TableB[1].TableAId, 2)
		t.Assert(tableA[1].TableB[1].TableC, nil)
	})
}

// Test_Table_Relation_With_MultipleDepends_Embedded tests WithAll on a chain of embedded struct attributes.
func Test_Table_Relation_With_MultipleDepends_Embedded(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	defer func() {
		dropTableWithDb(db, "table_a")
		dropTableWithDb(db, "table_b")
		dropTableWithDb(db, "table_c")
	}()
	for _, v := range chDWithMultipleDepends {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	chDOptimizeSchema(db)

	type TableC struct {
		gmeta.Meta `orm:"table_c"`
		Id         int `orm:"id,primary" json:"id"`
		TableBId   int `orm:"table_b_id" json:"table_b_id"`
	}

	type TableB struct {
		gmeta.Meta `orm:"table_b"`
		Id         int `orm:"id,primary" json:"id"`
		TableAId   int `orm:"table_a_id" json:"table_a_id"`
		*TableC    `orm:"with:table_b_id=id"  json:"table_c"`
	}

	type TableA struct {
		gmeta.Meta `orm:"table_a"`
		Id         int `orm:"id,primary" json:"id"`
		*TableB    `orm:"with:table_a_id=id" json:"table_b"`
	}

	// Struct.
	gtest.C(t, func(t *gtest.T) {
		var tableA *TableA
		err := db.Model("table_a").WithAll().Scan(&tableA)
		t.AssertNil(err)
		t.AssertNE(tableA, nil)
		t.Assert(tableA.Id, 1)

		t.AssertNE(tableA.TableB, nil)
		t.AssertNE(tableA.TableB.TableC, nil)
		t.Assert(tableA.TableB.TableAId, 1)
		t.Assert(tableA.TableB.TableC.Id, 100)
		t.Assert(tableA.TableB.TableC.TableBId, 10)
	})

	// Structs
	gtest.C(t, func(t *gtest.T) {
		var tableA []*TableA
		err := db.Model("table_a").WithAll().OrderAsc("id").Scan(&tableA)
		t.AssertNil(err)
		t.Assert(len(tableA), 2)
		t.AssertNE(tableA[0].TableB, nil)
		t.AssertNE(tableA[1].TableB, nil)
		t.AssertNE(tableA[0].TableB.TableC, nil)
		t.AssertNE(tableA[1].TableB.TableC, nil)

		t.Assert(tableA[0].Id, 1)
		t.Assert(tableA[0].TableB.Id, 10)
		t.Assert(tableA[0].TableB.TableC.Id, 100)

		t.Assert(tableA[1].Id, 2)
		t.Assert(tableA[1].TableB.Id, 20)
		t.Assert(tableA[1].TableB.TableC.Id, 300)
	})
}

// Test_Table_Relation_WithAll_Embedded_Meta_NameMatchingRule tests WithAll matching the relation fields by name rules.
func Test_Table_Relation_WithAll_Embedded_Meta_NameMatchingRule(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user100"
		tableUserDetail = "user_detail100"
		tableUserScores = "user_scores100"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
id UInt32 NOT NULL,
name String NOT NULL
) ENGINE = MergeTree()
ORDER BY id;
 `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
user_id UInt32 NOT NULL,
address String NOT NULL
) ENGINE = MergeTree()
ORDER BY user_id;
 `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
id UInt32 NOT NULL,
user_id UInt32 NOT NULL,
score UInt32 NOT NULL
) ENGINE = MergeTree()
ORDER BY id;
 `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail100"`
		UserID     int    `json:"user_id"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:user_scores100"`
		ID         int `json:"id"`
		UserID     int `json:"user_id"`
		Score      int `json:"score"`
	}

	// For Test Only
	type UserEmbedded struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user100"`
		UserEmbedded
		UserDetail UserDetail    `orm:"with:user_id=id"`
		UserScores []*UserScores `orm:"with:user_id=id"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"user_id": i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"id":      (i-1)*5 + j,
				"user_id": i,
				"score":   j,
			})
			gtest.AssertNil(err)
		}
	}

	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.ID, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserID, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].UserID, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserID, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
}

// Test_Table_Relation_WithAll_Unscoped tests WithAll with soft deleted relation records and the unscoped tag.
func Test_Table_Relation_WithAll_Unscoped(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user101"
		tableUserDetail = "user_detail101"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
id UInt32 NOT NULL,
name String NOT NULL
) ENGINE = MergeTree()
ORDER BY id;
 `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
user_id UInt32 NOT NULL,
address String NOT NULL,
deleted_at Nullable(DateTime)
) ENGINE = MergeTree()
ORDER BY user_id;
 `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail101"`
		UserID     int         `json:"user_id"`
		Address    string      `json:"address"`
		DeletedAt  *gtime.Time `json:"deleted_at"`
	}

	// For Test Only
	type UserEmbedded struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user101"`
		UserEmbedded
		UserDetail *UserDetail `orm:"with:user_id=id"`
	}
	type UserWithDeletedDetail struct {
		gmeta.Meta `orm:"table:user101"`
		UserEmbedded
		UserDetail *UserDetail `orm:"with:user_id=id, unscoped:true"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"user_id": i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		// Delete detail where i = 3
		if i == 3 {
			_, err = db.Delete(ctx, tableUserDetail, g.Map{
				"user_id": i,
			})
		}
		gtest.AssertNil(err)
	}
	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var user0 User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user0)
		t.AssertNil(err)
		t.Assert(user0.ID, 4)
		t.AssertNE(user0.UserDetail, nil)
		t.AssertNil(user0.UserDetail.DeletedAt)
		t.Assert(user0.UserDetail.UserID, 4)
		t.Assert(user0.UserDetail.Address, `address_4`)

		var user1 User
		err = db.Model(tableUser).WithAll().Where("id", 3).Scan(&user1)
		t.AssertNil(err)
		t.Assert(user1.ID, 3)
		t.AssertNil(user1.UserDetail)

		var user2 UserWithDeletedDetail
		err = db.Model(tableUser).WithAll().Where("id", 3).Scan(&user2)
		t.AssertNil(err)
		t.Assert(user2.ID, 3)
		t.AssertNE(user2.UserDetail, nil)
		t.AssertNE(user2.UserDetail.DeletedAt, nil)
		t.Assert(user2.UserDetail.UserID, 3)
		t.Assert(user2.UserDetail.Address, `address_3`)

		// Unscoped outside test
		var user3 User
		err = db.Model(tableUser).Unscoped().WithAll().Where("id", 3).Scan(&user3)
		t.AssertNil(err)
		t.Assert(user3.ID, 3)
		t.AssertNil(user3.UserDetail)
	})
}

// Test_Table_Relation_WithAll_Order tests WithAll with the order and unscoped tags.
func Test_Table_Relation_WithAll_Order(t *testing.T) {
	db, schema := chDNewSchemaDb()
	defer chDDropDatabase(schema)
	var (
		tableUser       = "user101"
		tableUserDetail = "user_detail101"
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
id UInt32 NOT NULL,
name String NOT NULL
) ENGINE = MergeTree()
ORDER BY id;
 `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
user_id UInt32 NOT NULL,
address String NOT NULL,
deleted_at Nullable(DateTime)
) ENGINE = MergeTree()
ORDER BY user_id;
 `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTableWithDb(db, tableUserDetail)

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail101"`
		UserID     int         `json:"user_id"`
		Address    string      `json:"address"`
		DeletedAt  *gtime.Time `json:"deleted_at"`
	}

	// For Test Only
	type UserEmbedded struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user101"`
		UserEmbedded
		UserDetail *UserDetail `orm:"with:user_id=id"`
	}
	type UserWithDeletedDetail struct {
		gmeta.Meta `orm:"table:user101"`
		UserEmbedded
		UserDetail *UserDetail `orm:"with:user_id=id, order:user_id asc,address desc, unscoped:true"`
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(ctx, tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.AssertNil(err)
		// Detail.
		_, err = db.Insert(ctx, tableUserDetail, g.Map{
			"user_id": i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		// Delete detail where i = 3
		if i == 3 {
			_, err = db.Delete(ctx, tableUserDetail, g.Map{
				"user_id": i,
			})
		}
		gtest.AssertNil(err)
	}
	chDOptimizeSchema(db)
	gtest.C(t, func(t *gtest.T) {
		var user0 User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user0)
		t.AssertNil(err)
		t.Assert(user0.ID, 4)
		t.AssertNE(user0.UserDetail, nil)
		t.AssertNil(user0.UserDetail.DeletedAt)
		t.Assert(user0.UserDetail.UserID, 4)
		t.Assert(user0.UserDetail.Address, `address_4`)

		var user1 User
		err = db.Model(tableUser).WithAll().Where("id", 3).Scan(&user1)
		t.AssertNil(err)
		t.Assert(user1.ID, 3)
		t.AssertNil(user1.UserDetail)

		var user2 UserWithDeletedDetail
		err = db.Model(tableUser).WithAll().Where("id", 3).Scan(&user2)
		t.AssertNil(err)
		t.Assert(user2.ID, 3)
		t.AssertNE(user2.UserDetail, nil)
		t.AssertNE(user2.UserDetail.DeletedAt, nil)
		t.Assert(user2.UserDetail.UserID, 3)
		t.Assert(user2.UserDetail.Address, `address_3`)

		// Unscoped outside test
		var user3 User
		err = db.Model(tableUser).Unscoped().WithAll().Where("id", 3).Scan(&user3)
		t.AssertNil(err)
		t.Assert(user3.ID, 3)
		t.AssertNil(user3.UserDetail)
	})
}
