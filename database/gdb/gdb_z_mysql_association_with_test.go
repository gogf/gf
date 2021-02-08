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

func Test_Table_Relation_With(t *testing.T) {
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
		db.SetDebug(true)
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
	})
}
