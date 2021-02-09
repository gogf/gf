// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gmeta"
	"testing"
)

func Test_Table_Relation_With_Scan(t *testing.T) {
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
id int(10) unsigned NOT NULL AUTO_INCREMENT,
name varchar(45) NOT NULL,
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
uid int(10) unsigned NOT NULL AUTO_INCREMENT,
address varchar(45) NOT NULL,
PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
id int(10) unsigned NOT NULL AUTO_INCREMENT,
uid int(10) unsigned NOT NULL,
score int(10) unsigned NOT NULL,
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

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
		_, err = db.Insert(tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.Assert(err, nil)
		// Detail.
		_, err = db.Insert(tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.Assert(err, nil)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(tableUserScores, g.Map{
				"uid":   i,
				"score": j,
			})
			gtest.Assert(err, nil)
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
			With(UserScores{}).
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

func Test_Table_Relation_With_ScanList(t *testing.T) {
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
id int(10) unsigned NOT NULL AUTO_INCREMENT,
name varchar(45) NOT NULL,
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
uid int(10) unsigned NOT NULL AUTO_INCREMENT,
address varchar(45) NOT NULL,
PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
id int(10) unsigned NOT NULL AUTO_INCREMENT,
uid int(10) unsigned NOT NULL,
score int(10) unsigned NOT NULL,
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

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
		_, err = db.Insert(tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.Assert(err, nil)
		// Detail.
		_, err = db.Insert(tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.Assert(err, nil)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(tableUserScores, g.Map{
				"uid":   i,
				"score": j,
			})
			gtest.Assert(err, nil)
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

func Test_Table_Relation_WithAll_Scan(t *testing.T) {
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
id int(10) unsigned NOT NULL AUTO_INCREMENT,
name varchar(45) NOT NULL,
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
uid int(10) unsigned NOT NULL AUTO_INCREMENT,
address varchar(45) NOT NULL,
PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
id int(10) unsigned NOT NULL AUTO_INCREMENT,
uid int(10) unsigned NOT NULL,
score int(10) unsigned NOT NULL,
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

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
		_, err = db.Insert(tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.Assert(err, nil)
		// Detail.
		_, err = db.Insert(tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.Assert(err, nil)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(tableUserScores, g.Map{
				"uid":   i,
				"score": j,
			})
			gtest.Assert(err, nil)
		}
	}
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

func Test_Table_Relation_WithAll_ScanList(t *testing.T) {
	var (
		tableUser       = "user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
id int(10) unsigned NOT NULL AUTO_INCREMENT,
name varchar(45) NOT NULL,
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
uid int(10) unsigned NOT NULL AUTO_INCREMENT,
address varchar(45) NOT NULL,
PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
id int(10) unsigned NOT NULL AUTO_INCREMENT,
uid int(10) unsigned NOT NULL,
score int(10) unsigned NOT NULL,
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
 `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

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
		_, err = db.Insert(tableUser, g.Map{
			"id":   i,
			"name": fmt.Sprintf(`name_%d`, i),
		})
		gtest.Assert(err, nil)
		// Detail.
		_, err = db.Insert(tableUserDetail, g.Map{
			"uid":     i,
			"address": fmt.Sprintf(`address_%d`, i),
		})
		gtest.Assert(err, nil)
		// Scores.
		for j := 1; j <= 5; j++ {
			_, err = db.Insert(tableUserScores, g.Map{
				"uid":   i,
				"score": j,
			})
			gtest.Assert(err, nil)
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
