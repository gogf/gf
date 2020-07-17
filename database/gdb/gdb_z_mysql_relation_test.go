// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"testing"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

func Test_Table_Relation(t *testing.T) {
	var (
		tableUser       = "user_" + gtime.TimestampMicroStr()
		tableUserDetail = "user_detail_" + gtime.TimestampMicroStr()
		tableUserScores = "user_scores_" + gtime.TimestampMicroStr()
	)
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
  uid int(10) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(45) NOT NULL,
  PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
  uid int(10) unsigned NOT NULL AUTO_INCREMENT,
  address varchar(45) NOT NULL,
  PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  uid int(10) unsigned NOT NULL,
  score int(10) unsigned NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type EntityUser struct {
		Uid  int    `json:"uid"`
		Name string `json:"name"`
	}
	type EntityUserDetail struct {
		Uid     int    `json:"uid"`
		Address string `json:"address"`
	}
	type EntityUserScores struct {
		Id    int `json:"id"`
		Uid   int `json:"uid"`
		Score int `json:"score"`
	}
	type Entity struct {
		User       *EntityUser
		UserDetail *EntityUserDetail
		UserScores []*EntityUserScores
	}

	// Initialize the data.
	var err error
	for i := 1; i <= 5; i++ {
		// User.
		_, err = db.Insert(tableUser, g.Map{
			"uid":  i,
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
	// MapKeyValue.
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(tableUser).Where("uid", g.Slice{3, 4}).Order("uid asc").All()
		t.Assert(err, nil)
		t.Assert(all.Len(), 2)
		t.Assert(len(all.MapKeyValue("uid")), 2)
		t.Assert(all.MapKeyValue("uid")["3"].Map()["uid"], 3)
		t.Assert(all.MapKeyValue("uid")["4"].Map()["uid"], 4)
		all, err = db.Table(tableUserScores).Where("uid", g.Slice{3, 4}).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(all.Len(), 10)
		t.Assert(len(all.MapKeyValue("uid")), 2)
		t.Assert(len(all.MapKeyValue("uid")["3"].Slice()), 5)
		t.Assert(len(all.MapKeyValue("uid")["4"].Slice()), 5)
		t.Assert(gconv.Map(all.MapKeyValue("uid")["3"].Slice()[0])["uid"], 3)
		t.Assert(gconv.Map(all.MapKeyValue("uid")["3"].Slice()[0])["score"], 1)
		t.Assert(gconv.Map(all.MapKeyValue("uid")["3"].Slice()[4])["uid"], 3)
		t.Assert(gconv.Map(all.MapKeyValue("uid")["3"].Slice()[4])["score"], 5)
	})
	// Result ScanList with struct elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []Entity
		// User
		all, err := db.Table(tableUser).Where("uid", g.Slice{3, 4}).Order("uid asc").All()
		t.Assert(err, nil)
		err = all.ScanList(&users, "User")
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Table(tableUserDetail).Where("uid", gdb.ListItemValues(users, "User", "Uid")).Order("uid asc").All()
		gtest.Assert(err, nil)
		err = all.ScanList(&users, "UserDetail", "User", "uid:Uid")
		t.Assert(err, nil)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Table(tableUserScores).Where("uid", gdb.ListItemValues(users, "User", "Uid")).Order("id asc").All()
		gtest.Assert(err, nil)
		err = all.ScanList(&users, "UserScores", "User", "uid:Uid")
		t.Assert(err, nil)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].Uid, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		all, err := db.Table(tableUser).Where("uid", g.Slice{3, 4}).Order("uid asc").All()
		t.Assert(err, nil)
		err = all.ScanList(&users, "User")
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Table(tableUserDetail).Where("uid", gdb.ListItemValues(users, "User", "Uid")).Order("uid asc").All()
		gtest.Assert(err, nil)
		err = all.ScanList(&users, "UserDetail", "User", "uid:Uid")
		t.Assert(err, nil)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Table(tableUserScores).Where("uid", gdb.ListItemValues(users, "User", "Uid")).Order("id asc").All()
		gtest.Assert(err, nil)
		err = all.ScanList(&users, "UserScores", "User", "uid:Uid")
		t.Assert(err, nil)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].Uid, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with struct elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			Uid  int    `json:"uid"`
			Name string `json:"name"`
		}
		type EntityUserDetail struct {
			Uid     int    `json:"uid"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id    int `json:"id"`
			Uid   int `json:"uid"`
			Score int `json:"score"`
		}
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []Entity
		// User
		all, err := db.Table(tableUser).Where("uid", g.Slice{3, 4}).Order("uid asc").All()
		t.Assert(err, nil)
		err = all.ScanList(&users, "User")
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Table(tableUserDetail).Where("uid", gdb.ListItemValues(users, "User", "Uid")).Order("uid asc").All()
		gtest.Assert(err, nil)
		err = all.ScanList(&users, "UserDetail", "User", "uid:Uid")
		t.Assert(err, nil)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Table(tableUserScores).Where("uid", gdb.ListItemValues(users, "User", "Uid")).Order("id asc").All()
		gtest.Assert(err, nil)
		err = all.ScanList(&users, "UserScores", "User", "uid:Uid")
		t.Assert(err, nil)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].Uid, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with pointer elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			Uid  int    `json:"uid"`
			Name string `json:"name"`
		}
		type EntityUserDetail struct {
			Uid     int    `json:"uid"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id    int `json:"id"`
			Uid   int `json:"uid"`
			Score int `json:"score"`
		}
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []*Entity

		// User
		all, err := db.Table(tableUser).Where("uid", g.Slice{3, 4}).Order("uid asc").All()
		t.Assert(err, nil)
		err = all.ScanList(&users, "User")
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Table(tableUserDetail).Where("uid", gdb.ListItemValues(users, "User", "Uid")).Order("uid asc").All()
		gtest.Assert(err, nil)
		err = all.ScanList(&users, "UserDetail", "User", "uid:Uid")
		t.Assert(err, nil)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Table(tableUserScores).Where("uid", gdb.ListItemValues(users, "User", "Uid")).Order("id asc").All()
		gtest.Assert(err, nil)
		err = all.ScanList(&users, "UserScores", "User", "uid:Uid")
		t.Assert(err, nil)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].Uid, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Model ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		err := db.Table(tableUser).
			Where("uid", g.Slice{3, 4}).
			Order("uid asc").
			ScanList(&users, "User")
		t.Assert(err, nil)
		// Detail
		err = db.Table(tableUserDetail).
			Where("uid", gdb.ListItemValues(users, "User", "Uid")).
			Order("uid asc").
			ScanList(&users, "UserDetail", "User", "uid:Uid")
		gtest.Assert(err, nil)
		// Scores
		err = db.Table(tableUserScores).
			Where("uid", gdb.ListItemValues(users, "User", "Uid")).
			Order("id asc").
			ScanList(&users, "UserScores", "User", "uid:Uid")
		t.Assert(err, nil)

		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})

		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})

		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].Uid, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
}
