// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
)

func withCreateTable(table, columns, autoIncrement string) {
	dropTable(table)
	if _, err := db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (%s)`, table, columns)); err != nil {
		gtest.Error(err)
	}
	if autoIncrement != "" {
		createAutoIncrement(table, autoIncrement, 1)
	}
}

func withCreateTableUser(table string) {
	withCreateTable(table, `
id NUMBER(10) NOT NULL,
name VARCHAR2(45) NOT NULL,
PRIMARY KEY (id)
`, "")
}

func withCreateTableUserDetail(table string) {
	withCreateTable(table, `
user_id NUMBER(10) NOT NULL,
address VARCHAR2(45) NOT NULL,
PRIMARY KEY (user_id)
`, "")
}

func withCreateTableUserScores(table string) {
	withCreateTable(table, `
id NUMBER(10) NOT NULL,
user_id NUMBER(10) NOT NULL,
score NUMBER(10) NOT NULL,
PRIMARY KEY (id)
`, "id")
}

func withCreateMultipleDependsTables() {
	withCreateTable("with_table_a", `
id NUMBER(10) NOT NULL,
alias VARCHAR2(255),
PRIMARY KEY (id)
`, "")
	withCreateTable("with_table_b", `
id NUMBER(10) NOT NULL,
table_a_id NUMBER(10) NOT NULL,
alias VARCHAR2(255),
PRIMARY KEY (id)
`, "")
	withCreateTable("with_table_c", `
id NUMBER(10) NOT NULL,
table_b_id NUMBER(10) NOT NULL,
alias VARCHAR2(255),
PRIMARY KEY (id)
`, "")
	for _, v := range []string{
		`INSERT INTO with_table_a VALUES (1, 'table_a_test1')`,
		`INSERT INTO with_table_a VALUES (2, 'table_a_test2')`,
		`INSERT INTO with_table_b VALUES (10, 1, 'table_b_test1')`,
		`INSERT INTO with_table_b VALUES (20, 2, 'table_b_test2')`,
		`INSERT INTO with_table_b VALUES (30, 1, 'table_b_test3')`,
		`INSERT INTO with_table_b VALUES (40, 2, 'table_b_test4')`,
		`INSERT INTO with_table_c VALUES (100, 10, 'table_c_test1')`,
		`INSERT INTO with_table_c VALUES (200, 10, 'table_c_test2')`,
		`INSERT INTO with_table_c VALUES (300, 20, 'table_c_test3')`,
		`INSERT INTO with_table_c VALUES (400, 30, 'table_c_test4')`,
	} {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
}

func Test_Table_Relation_With_Scan(t *testing.T) {
	var (
		tableUser       = "with_user"
		tableUserDetail = "with_user_detail"
		tableUserScores = "with_user_score"
	)
	withCreateTableUser(tableUser)
	createAutoIncrement(tableUser, "id", 1)
	defer dropTable(tableUser)

	withCreateTableUserDetail(tableUserDetail)
	defer dropTable(tableUserDetail)

	withCreateTableUserScores(tableUserScores)
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_user_detail"`
		UserId     int    `json:"user_id"`
		Address    string `json:"address"`
	}

	type UserScore struct {
		gmeta.Meta `orm:"table:with_user_score"`
		Id         int `json:"id"`
		UserId     int `json:"user_id"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:with_user"`
		Id         int          `json:"id"`
		Name       string       `json:"name"`
		UserDetail *UserDetail  `orm:"with:user_id=id"`
		UserScores []*UserScore `orm:"with:user_id=id"`
	}

	// Initialize the data.
	gtest.C(t, func(t *gtest.T) {
		for i := 1; i <= 5; i++ {
			// User.
			user := User{
				Name: fmt.Sprintf(`name_%d`, i),
			}
			lastInsertId, err := db.Model(user).Data(user).OmitEmpty().InsertAndGetId()
			t.AssertNil(err)
			// Detail.
			userDetail := UserDetail{
				UserId:  int(lastInsertId),
				Address: fmt.Sprintf(`address_%d`, lastInsertId),
			}
			_, err = db.Model(userDetail).Data(userDetail).OmitEmpty().Insert()
			t.AssertNil(err)
			// Scores.
			for j := 1; j <= 5; j++ {
				userScore := UserScore{
					UserId: int(lastInsertId),
					Score:  j,
				}
				_, err = db.Model(userScore).Data(userScore).OmitEmpty().Insert()
				t.AssertNil(err)
			}
		}
	})
	for i := 1; i <= 5; i++ {
		// User.
		user := User{
			Name: fmt.Sprintf(`name_%d`, i),
		}
		lastInsertId, err := db.Model(user).Data(user).OmitEmpty().InsertAndGetId()
		gtest.AssertNil(err)
		// Detail.
		userDetail := UserDetail{
			UserId:  int(lastInsertId),
			Address: fmt.Sprintf(`address_%d`, lastInsertId),
		}
		_, err = db.Model(userDetail).Data(userDetail).Insert()
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			userScore := UserScore{
				UserId: int(lastInsertId),
				Score:  j,
			}
			_, err = db.Model(userScore).Data(userScore).OmitEmpty().Insert()
			gtest.AssertNil(err)
		}
	}
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
		t.Assert(user.UserDetail.UserId, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].UserId, 3)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserId, 3)
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
		t.Assert(user.UserDetail.UserId, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].UserId, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserId, 4)
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
		t.Assert(user.UserDetail.UserId, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].UserId, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserId, 4)
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
		t.Assert(user.UserDetail.UserId, 4)
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
		t.Assert(user.UserScores[0].UserId, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserId, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_With(t *testing.T) {
	var (
		tableUser       = "with_user"
		tableUserDetail = "with_user_detail"
		tableUserScores = "with_user_scores"
	)
	withCreateTableUser(tableUser)
	defer dropTable(tableUser)

	withCreateTableUserDetail(tableUserDetail)
	defer dropTable(tableUserDetail)

	withCreateTableUserScores(tableUserScores)
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_user_detail"`
		UserId     int    `json:"user_id"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:with_user_scores"`
		Id         int `json:"id"`
		UserId     int `json:"user_id"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:with_user"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:user_id=id"`
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
				"user_id": i,
				"score":   j,
			})
			gtest.AssertNil(err)
		}
	}

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
		t.Assert(users[0].UserDetail.UserId, 3)
		t.Assert(users[0].UserDetail.Address, "address_3")
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].UserId, 3)
		t.Assert(users[0].UserScores[4].Score, 5)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.UserId, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].UserId, 4)
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
		t.Assert(users[0].UserDetail.UserId, 3)
		t.Assert(users[0].UserDetail.Address, "address_3")
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].UserId, 3)
		t.Assert(users[0].UserScores[4].Score, 5)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.UserId, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].UserId, 4)
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
		t.Assert(users[0].UserDetail.UserId, 3)
		t.Assert(users[0].UserDetail.Address, "address_3")
		t.Assert(len(users[0].UserScores), 0)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.UserId, 4)
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
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].UserId, 3)
		t.Assert(users[0].UserScores[4].Score, 5)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.Assert(users[1].UserDetail, nil)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].UserId, 4)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_WithAll(t *testing.T) {
	var (
		tableUser       = "with_user"
		tableUserDetail = "with_user_detail"
		tableUserScores = "with_user_scores"
	)
	withCreateTableUser(tableUser)
	defer dropTable(tableUser)

	withCreateTableUserDetail(tableUserDetail)
	defer dropTable(tableUserDetail)

	withCreateTableUserScores(tableUserScores)
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_user_detail"`
		UserId     int    `json:"user_id"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:with_user_scores"`
		Id         int `json:"id"`
		UserId     int `json:"user_id"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:with_user"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:user_id=id"`
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
				"user_id": i,
				"score":   j,
			})
			gtest.AssertNil(err)
		}
	}
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserId, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].UserId, 3)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserId, 3)
		t.Assert(user.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserId, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].UserId, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserId, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_WithAll_List(t *testing.T) {
	var (
		tableUser       = "with_user"
		tableUserDetail = "with_user_detail"
		tableUserScores = "with_user_scores"
	)
	withCreateTableUser(tableUser)
	defer dropTable(tableUser)

	withCreateTableUserDetail(tableUserDetail)
	defer dropTable(tableUserDetail)

	withCreateTableUserScores(tableUserScores)
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_user_detail"`
		UserId     int    `json:"user_id"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:with_user_scores"`
		Id         int `json:"id"`
		UserId     int `json:"user_id"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:with_user"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:user_id=id"`
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
				"user_id": i,
				"score":   j,
			})
			gtest.AssertNil(err)
		}
	}
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		err := db.Model(tableUser).WithAll().Where("id", []int{3, 4}).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, 3)
		t.Assert(users[0].Name, "name_3")
		t.AssertNE(users[0].UserDetail, nil)
		t.Assert(users[0].UserDetail.UserId, 3)
		t.Assert(users[0].UserDetail.Address, "address_3")
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].UserId, 3)
		t.Assert(users[0].UserScores[4].Score, 5)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.UserId, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].UserId, 4)
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
		t.Assert(users[0].UserDetail.UserId, 3)
		t.Assert(users[0].UserDetail.Address, "address_3")
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].UserId, 3)
		t.Assert(users[0].UserScores[4].Score, 5)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.UserId, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].UserId, 4)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_WithAllCondition_List(t *testing.T) {
	var (
		tableUser       = "with_user"
		tableUserDetail = "with_user_detail"
		tableUserScores = "with_user_scores"
	)
	withCreateTableUser(tableUser)
	defer dropTable(tableUser)

	withCreateTableUserDetail(tableUserDetail)
	defer dropTable(tableUserDetail)

	withCreateTableUserScores(tableUserScores)
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_user_detail"`
		UserId     int    `json:"user_id"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:with_user_scores"`
		Id         int `json:"id"`
		UserId     int `json:"user_id"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:with_user"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:user_id=id, where:user_id > 3"`
		UserScores []*UserScores `orm:"with:user_id=id, where:score>1 and score<5, order:score desc"`
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
				"user_id": i,
				"score":   j,
			})
			gtest.AssertNil(err)
		}
	}

	db.SetDebug(true)
	defer db.SetDebug(false)

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
		t.Assert(users[1].UserDetail.UserId, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 3)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 4)
		t.Assert(users[1].UserScores[2].UserId, 4)
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
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 4)
		t.Assert(users[0].UserScores[2].UserId, 3)
		t.Assert(users[0].UserScores[2].Score, 2)

		t.Assert(users[1].Id, 4)
		t.Assert(users[1].Name, "name_4")
		t.AssertNE(users[1].UserDetail, nil)
		t.Assert(users[1].UserDetail.UserId, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 3)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 4)
		t.Assert(users[1].UserScores[2].UserId, 4)
		t.Assert(users[1].UserScores[2].Score, 2)
	})
}

func Test_Table_Relation_WithAll_Embedded_With_SelfMaintained_Attributes(t *testing.T) {
	var (
		tableUser       = "with_user"
		tableUserDetail = "with_user_detail"
		tableUserScores = "with_user_scores"
	)
	withCreateTableUser(tableUser)
	defer dropTable(tableUser)

	withCreateTableUserDetail(tableUserDetail)
	defer dropTable(tableUserDetail)

	withCreateTableUserScores(tableUserScores)
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_user_detail"`
		UserId     int    `json:"user_id"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:with_user_scores"`
		Id         int `json:"id"`
		UserId     int `json:"user_id"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta  `orm:"table:with_user"`
		*UserDetail `orm:"with:user_id=id"`
		Id          int           `json:"id"`
		Name        string        `json:"name"`
		UserScores  []*UserScores `orm:"with:user_id=id"`
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
				"user_id": i,
				"score":   j,
			})
			gtest.AssertNil(err)
		}
	}
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserId, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].UserId, 3)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserId, 3)
		t.Assert(user.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserId, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].UserId, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserId, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_WithAll_Embedded_Without_SelfMaintained_Attributes(t *testing.T) {
	var (
		tableUser       = "with_user"
		tableUserDetail = "with_user_detail"
		tableUserScores = "with_user_scores"
	)
	withCreateTableUser(tableUser)
	defer dropTable(tableUser)

	withCreateTableUserDetail(tableUserDetail)
	defer dropTable(tableUserDetail)

	withCreateTableUserScores(tableUserScores)
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_user_detail"`
		UserId     int    `json:"user_id"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:with_user_scores"`
		Id         int `json:"id"`
		UserId     int `json:"user_id"`
		Score      int `json:"score"`
	}

	// For Test Only
	type UserEmbedded struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	type User struct {
		gmeta.Meta  `orm:"table:with_user"`
		*UserDetail `orm:"with:user_id=id"`
		UserEmbedded
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
				"user_id": i,
				"score":   j,
			})
			gtest.AssertNil(err)
		}
	}
	db.SetDebug(true)
	defer db.SetDebug(false)

	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserId, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].UserId, 3)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserId, 3)
		t.Assert(user.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserId, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].UserId, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserId, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_WithAll_Embedded_WithoutMeta(t *testing.T) {
	var (
		tableUser       = "with_user"
		tableUserDetail = "with_user_detail"
		tableUserScores = "with_user_scores"
	)
	withCreateTableUser(tableUser)
	defer dropTable(tableUser)

	withCreateTableUserDetail(tableUserDetail)
	defer dropTable(tableUserDetail)

	withCreateTableUserScores(tableUserScores)
	defer dropTable(tableUserScores)

	type UserDetailBase struct {
		UserId  int    `json:"user_id"`
		Address string `json:"address"`
	}

	type WithUserDetail struct {
		UserDetailBase
	}

	type WithUserScores struct {
		Id     int `json:"id"`
		UserId int `json:"user_id"`
		Score  int `json:"score"`
	}

	type User struct {
		*WithUserDetail `orm:"with:user_id=id"`
		Id              int               `json:"id"`
		Name            string            `json:"name"`
		UserScores      []*WithUserScores `orm:"with:user_id=id"`
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
				"user_id": i,
				"score":   j,
			})
			gtest.AssertNil(err)
		}
	}
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.WithUserDetail, nil)
		t.Assert(user.WithUserDetail.UserId, 3)
		t.Assert(user.WithUserDetail.Address, `address_3`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].UserId, 3)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserId, 3)
		t.Assert(user.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.WithUserDetail, nil)
		t.Assert(user.WithUserDetail.UserId, 4)
		t.Assert(user.WithUserDetail.Address, `address_4`)
		t.Assert(len(user.UserScores), 5)
		t.Assert(user.UserScores[0].UserId, 4)
		t.Assert(user.UserScores[0].Score, 1)
		t.Assert(user.UserScores[4].UserId, 4)
		t.Assert(user.UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_WithAll_AttributeStructAlsoHasWithTag(t *testing.T) {
	var (
		tableUser       = "with_user"
		tableUserDetail = "with_user_detail"
		tableUserScores = "with_user_scores"
	)
	withCreateTableUser(tableUser)
	defer dropTable(tableUser)

	withCreateTableUserDetail(tableUserDetail)
	defer dropTable(tableUserDetail)

	withCreateTableUserScores(tableUserScores)
	defer dropTable(tableUserScores)

	type UserScores struct {
		gmeta.Meta `orm:"table:with_user_scores"`
		Id         int `json:"id"`
		UserId     int `json:"user_id"`
		Score      int `json:"score"`
	}

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_user_detail"`
		UserId     int           `json:"user_id"`
		Address    string        `json:"address"`
		UserScores []*UserScores `orm:"with:user_id"`
	}

	type User struct {
		gmeta.Meta  `orm:"table:with_user"`
		*UserDetail `orm:"with:user_id=id"`
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
			"user_id": i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"user_id": i,
				"score":   j,
			})
			gtest.AssertNil(err)
		}
	}

	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserId, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserDetail.UserScores), 5)
		t.Assert(user.UserDetail.UserScores[0].UserId, 3)
		t.Assert(user.UserDetail.UserScores[0].Score, 1)
		t.Assert(user.UserDetail.UserScores[4].UserId, 3)
		t.Assert(user.UserDetail.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserId, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserDetail.UserScores), 5)
		t.Assert(user.UserDetail.UserScores[0].UserId, 4)
		t.Assert(user.UserDetail.UserScores[0].Score, 1)
		t.Assert(user.UserDetail.UserScores[4].UserId, 4)
		t.Assert(user.UserDetail.UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_WithAll_AttributeStructAlsoHasWithTag_MoreDeep(t *testing.T) {
	var (
		tableUser       = "with_user"
		tableUserDetail = "with_user_detail"
		tableUserScores = "with_user_scores"
	)
	withCreateTableUser(tableUser)
	defer dropTable(tableUser)

	withCreateTableUserDetail(tableUserDetail)
	defer dropTable(tableUserDetail)

	withCreateTableUserScores(tableUserScores)
	defer dropTable(tableUserScores)

	type UserScores struct {
		gmeta.Meta `orm:"table:with_user_scores"`
		Id         int `json:"id"`
		UserId     int `json:"user_id"`
		Score      int `json:"score"`
	}

	type UserDetail1 struct {
		gmeta.Meta `orm:"table:with_user_detail"`
		UserId     int           `json:"user_id"`
		Address    string        `json:"address"`
		UserScores []*UserScores `orm:"with:user_id"`
	}

	type UserDetail2 struct {
		gmeta.Meta  `orm:"table:with_user_detail"`
		UserId      int           `json:"user_id"`
		Address     string        `json:"address"`
		UserDetail1 *UserDetail1  `orm:"with:user_id"`
		UserScores  []*UserScores `orm:"with:user_id"`
	}

	type UserDetail3 struct {
		gmeta.Meta  `orm:"table:with_user_detail"`
		UserId      int           `json:"user_id"`
		Address     string        `json:"address"`
		UserDetail2 *UserDetail2  `orm:"with:user_id"`
		UserScores  []*UserScores `orm:"with:user_id"`
	}

	type UserDetail struct {
		gmeta.Meta  `orm:"table:with_user_detail"`
		UserId      int           `json:"user_id"`
		Address     string        `json:"address"`
		UserDetail3 *UserDetail3  `orm:"with:user_id"`
		UserScores  []*UserScores `orm:"with:user_id"`
	}

	type User struct {
		gmeta.Meta  `orm:"table:with_user"`
		*UserDetail `orm:"with:user_id=id"`
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
			"user_id": i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"user_id": i,
				"score":   j,
			})
			gtest.AssertNil(err)
		}
	}

	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserId, 3)
		t.Assert(user.UserDetail.UserDetail3.UserId, 3)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.UserId, 3)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.UserDetail1.UserId, 3)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserDetail.UserScores), 5)
		t.Assert(user.UserDetail.UserScores[0].UserId, 3)
		t.Assert(user.UserDetail.UserScores[0].Score, 1)
		t.Assert(user.UserDetail.UserScores[4].UserId, 3)
		t.Assert(user.UserDetail.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserId, 4)
		t.Assert(user.UserDetail.UserDetail3.UserId, 4)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.UserId, 4)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.UserDetail1.UserId, 4)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserDetail.UserScores), 5)
		t.Assert(user.UserDetail.UserScores[0].UserId, 4)
		t.Assert(user.UserDetail.UserScores[0].Score, 1)
		t.Assert(user.UserDetail.UserScores[4].UserId, 4)
		t.Assert(user.UserDetail.UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_With_AttributeStructAlsoHasWithTag_MoreDeep(t *testing.T) {
	var (
		tableUser       = "with_user"
		tableUserDetail = "with_user_detail"
		tableUserScores = "with_user_scores"
	)
	withCreateTableUser(tableUser)
	defer dropTable(tableUser)

	withCreateTableUserDetail(tableUserDetail)
	defer dropTable(tableUserDetail)

	withCreateTableUserScores(tableUserScores)
	defer dropTable(tableUserScores)

	type UserScores struct {
		gmeta.Meta `orm:"table:with_user_scores"`
		Id         int `json:"id"`
		UserId     int `json:"user_id"`
		Score      int `json:"score"`
	}

	type UserDetail1 struct {
		gmeta.Meta `orm:"table:with_user_detail"`
		UserId     int           `json:"user_id"`
		Address    string        `json:"address"`
		UserScores []*UserScores `orm:"with:user_id"`
	}

	type UserDetail2 struct {
		gmeta.Meta  `orm:"table:with_user_detail"`
		UserId      int           `json:"user_id"`
		Address     string        `json:"address"`
		UserDetail1 *UserDetail1  `orm:"with:user_id"`
		UserScores  []*UserScores `orm:"with:user_id"`
	}

	type UserDetail3 struct {
		gmeta.Meta  `orm:"table:with_user_detail"`
		UserId      int           `json:"user_id"`
		Address     string        `json:"address"`
		UserDetail2 *UserDetail2  `orm:"with:user_id"`
		UserScores  []*UserScores `orm:"with:user_id"`
	}

	type UserDetail struct {
		gmeta.Meta  `orm:"table:with_user_detail"`
		UserId      int           `json:"user_id"`
		Address     string        `json:"address"`
		UserDetail3 *UserDetail3  `orm:"with:user_id"`
		UserScores  []*UserScores `orm:"with:user_id"`
	}

	type User struct {
		gmeta.Meta  `orm:"table:with_user"`
		*UserDetail `orm:"with:user_id=id"`
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
			"user_id": i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.AssertNil(err)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(ctx, tableUserScores, g.Map{
				"user_id": i,
				"score":   j,
			})
			gtest.AssertNil(err)
		}
	}

	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).With(UserDetail{}, UserDetail2{}, UserDetail3{}, UserScores{}).Where("id", 3).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 3)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserId, 3)
		t.Assert(user.UserDetail.UserDetail3.UserId, 3)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.UserId, 3)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.UserDetail1, nil)
		t.Assert(user.UserDetail.Address, `address_3`)
		t.Assert(len(user.UserDetail.UserScores), 5)
		t.Assert(user.UserDetail.UserScores[0].UserId, 3)
		t.Assert(user.UserDetail.UserScores[0].Score, 1)
		t.Assert(user.UserDetail.UserScores[4].UserId, 3)
		t.Assert(user.UserDetail.UserScores[4].Score, 5)
	})
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).With(UserDetail{}, UserDetail2{}, UserDetail3{}, UserScores{}).Where("id", 4).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 4)
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.UserId, 4)
		t.Assert(user.UserDetail.UserDetail3.UserId, 4)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.UserId, 4)
		t.Assert(user.UserDetail.UserDetail3.UserDetail2.UserDetail1, nil)
		t.Assert(user.UserDetail.Address, `address_4`)
		t.Assert(len(user.UserDetail.UserScores), 5)
		t.Assert(user.UserDetail.UserScores[0].UserId, 4)
		t.Assert(user.UserDetail.UserScores[0].Score, 1)
		t.Assert(user.UserDetail.UserScores[4].UserId, 4)
		t.Assert(user.UserDetail.UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_With_MultipleDepends1(t *testing.T) {
	defer func() {
		dropTable("with_table_a")
		dropTable("with_table_b")
		dropTable("with_table_c")
	}()
	withCreateMultipleDependsTables()

	type WithTableC struct {
		gmeta.Meta `orm:"with_table_c"`
		Id         int `orm:"id,primary" json:"id"`
		TableBId   int `orm:"table_b_id" json:"table_b_id"`
	}

	type WithTableB struct {
		gmeta.Meta `orm:"with_table_b"`
		Id         int         `orm:"id,primary" json:"id"`
		TableAId   int         `orm:"table_a_id" json:"table_a_id"`
		TableC     *WithTableC `orm:"with:table_b_id=id"  json:"table_c"`
	}

	type WithTableA struct {
		gmeta.Meta `orm:"with_table_a"`
		Id         int         `orm:"id,primary" json:"id"`
		TableB     *WithTableB `orm:"with:table_a_id=id" json:"table_b"`
	}

	db.SetDebug(true)
	defer db.SetDebug(false)

	// Struct.
	gtest.C(t, func(t *gtest.T) {
		var tableA *WithTableA
		err := db.Model("with_table_a").WithAll().Scan(&tableA)
		// g.Dump(tableA)
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
		var tableA []*WithTableA
		err := db.Model("with_table_a").WithAll().OrderAsc("id").Scan(&tableA)
		// g.Dump(tableA)
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

func Test_Table_Relation_With_MultipleDepends2(t *testing.T) {
	defer func() {
		dropTable("with_table_a")
		dropTable("with_table_b")
		dropTable("with_table_c")
	}()
	withCreateMultipleDependsTables()

	type WithTableC struct {
		gmeta.Meta `orm:"with_table_c"`
		Id         int `orm:"id,primary" json:"id"`
		TableBId   int `orm:"table_b_id" json:"table_b_id"`
	}

	type WithTableB struct {
		gmeta.Meta `orm:"with_table_b"`
		Id         int           `orm:"id,primary" json:"id"`
		TableAId   int           `orm:"table_a_id" json:"table_a_id"`
		TableC     []*WithTableC `orm:"with:table_b_id=id"  json:"table_c"`
	}

	type WithTableA struct {
		gmeta.Meta `orm:"with_table_a"`
		Id         int           `orm:"id,primary" json:"id"`
		TableB     []*WithTableB `orm:"with:table_a_id=id" json:"table_b"`
	}

	db.SetDebug(true)
	defer db.SetDebug(false)

	// Struct.
	gtest.C(t, func(t *gtest.T) {
		var tableA *WithTableA
		err := db.Model("with_table_a").WithAll().Scan(&tableA)
		// g.Dump(tableA)
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
		var tableA []*WithTableA
		err := db.Model("with_table_a").WithAll().OrderAsc("id").Scan(&tableA)
		// g.Dump(tableA)
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

func Test_Table_Relation_With_MultipleDepends_Embedded(t *testing.T) {
	defer func() {
		dropTable("with_table_a")
		dropTable("with_table_b")
		dropTable("with_table_c")
	}()
	withCreateMultipleDependsTables()

	type WithTableC struct {
		gmeta.Meta `orm:"with_table_c"`
		Id         int `orm:"id,primary" json:"id"`
		TableBId   int `orm:"table_b_id" json:"table_b_id"`
	}

	type WithTableB struct {
		gmeta.Meta  `orm:"with_table_b"`
		Id          int `orm:"id,primary" json:"id"`
		TableAId    int `orm:"table_a_id" json:"table_a_id"`
		*WithTableC `orm:"with:table_b_id=id"  json:"table_c"`
	}

	type WithTableA struct {
		gmeta.Meta  `orm:"with_table_a"`
		Id          int `orm:"id,primary" json:"id"`
		*WithTableB `orm:"with:table_a_id=id" json:"table_b"`
	}

	db.SetDebug(true)
	defer db.SetDebug(false)

	// Struct.
	gtest.C(t, func(t *gtest.T) {
		var tableA *WithTableA
		err := db.Model("with_table_a").WithAll().Scan(&tableA)
		// g.Dump(tableA)
		t.AssertNil(err)
		t.AssertNE(tableA, nil)
		t.Assert(tableA.Id, 1)

		t.AssertNE(tableA.WithTableB, nil)
		t.AssertNE(tableA.WithTableB.WithTableC, nil)
		t.Assert(tableA.WithTableB.TableAId, 1)
		t.Assert(tableA.WithTableB.WithTableC.Id, 100)
		t.Assert(tableA.WithTableB.WithTableC.TableBId, 10)
	})

	// Structs
	gtest.C(t, func(t *gtest.T) {
		var tableA []*WithTableA
		err := db.Model("with_table_a").WithAll().OrderAsc("id").Scan(&tableA)
		// g.Dump(tableA)
		t.AssertNil(err)
		t.Assert(len(tableA), 2)
		t.AssertNE(tableA[0].WithTableB, nil)
		t.AssertNE(tableA[1].WithTableB, nil)
		t.AssertNE(tableA[0].WithTableB.WithTableC, nil)
		t.AssertNE(tableA[1].WithTableB.WithTableC, nil)

		t.Assert(tableA[0].Id, 1)
		t.Assert(tableA[0].WithTableB.Id, 10)
		t.Assert(tableA[0].WithTableB.WithTableC.Id, 100)

		t.Assert(tableA[1].Id, 2)
		t.Assert(tableA[1].WithTableB.Id, 20)
		t.Assert(tableA[1].WithTableB.WithTableC.Id, 300)
	})
}

func Test_Table_Relation_WithAll_Embedded_Meta_NameMatchingRule(t *testing.T) {
	var (
		tableUser       = "with_user100"
		tableUserDetail = "with_user_detail100"
		tableUserScores = "with_user_scores100"
	)
	withCreateTable(tableUser, `
id NUMBER(10) NOT NULL,
name VARCHAR2(45) NOT NULL,
PRIMARY KEY (id)
`, "")
	defer dropTable(tableUser)

	withCreateTable(tableUserDetail, `
user_id NUMBER(10) NOT NULL,
address VARCHAR2(45) NOT NULL,
PRIMARY KEY (user_id)
`, "")
	defer dropTable(tableUserDetail)

	withCreateTable(tableUserScores, `
id NUMBER(10) NOT NULL,
user_id NUMBER(10) NOT NULL,
score NUMBER(10) NOT NULL,
PRIMARY KEY (id)
`, "id")
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_user_detail100"`
		UserID     int    `json:"user_id"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:with_user_scores100"`
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
		gmeta.Meta `orm:"table:with_user100"`
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
				"user_id": i,
				"score":   j,
			})
			gtest.AssertNil(err)
		}
	}

	// gtest.C(t, func(t *gtest.T) {
	//	var user *User
	//	err := db.Model(tableUser).WithAll().Where("id", 3).Scan(&user)
	//	t.AssertNil(err)
	//	t.Assert(user.ID, 3)
	//	t.AssertNE(user.UserDetail, nil)
	//	t.Assert(user.UserDetail.UserID, 3)
	//	t.Assert(user.UserDetail.Address, `address_3`)
	//	t.Assert(len(user.UserScores), 5)
	//	t.Assert(user.UserScores[0].UserID, 3)
	//	t.Assert(user.UserScores[0].Score, 1)
	//	t.Assert(user.UserScores[4].UserID, 3)
	//	t.Assert(user.UserScores[4].Score, 5)
	// })
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

func Test_Table_Relation_WithAll_Unscoped(t *testing.T) {
	var (
		tableUser       = "with_user101"
		tableUserDetail = "with_user_detail101"
	)
	withCreateTable(tableUser, `
id NUMBER(10) NOT NULL,
name VARCHAR2(45) NOT NULL,
PRIMARY KEY (id)
`, "")
	defer dropTable(tableUser)

	withCreateTable(tableUserDetail, `
user_id NUMBER(10) NOT NULL,
address VARCHAR2(45) NOT NULL,
deleted_at TIMESTAMP DEFAULT NULL,
PRIMARY KEY (user_id)
`, "")
	defer dropTable(tableUserDetail)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_user_detail101"`
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
		gmeta.Meta `orm:"table:with_user101"`
		UserEmbedded
		UserDetail *UserDetail `orm:"with:user_id=id"`
	}
	type UserWithDeletedDetail struct {
		gmeta.Meta `orm:"table:with_user101"`
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

func Test_Table_Relation_WithAll_Order(t *testing.T) {
	var (
		tableUser       = "with_user101"
		tableUserDetail = "with_user_detail101"
	)
	withCreateTable(tableUser, `
id NUMBER(10) NOT NULL,
name VARCHAR2(45) NOT NULL,
PRIMARY KEY (id)
`, "")
	defer dropTable(tableUser)

	withCreateTable(tableUserDetail, `
user_id NUMBER(10) NOT NULL,
address VARCHAR2(45) NOT NULL,
deleted_at TIMESTAMP DEFAULT NULL,
PRIMARY KEY (user_id)
`, "")
	defer dropTable(tableUserDetail)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_user_detail101"`
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
		gmeta.Meta `orm:"table:with_user101"`
		UserEmbedded
		UserDetail *UserDetail `orm:"with:user_id=id"`
	}
	type UserWithDeletedDetail struct {
		gmeta.Meta `orm:"table:with_user101"`
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
