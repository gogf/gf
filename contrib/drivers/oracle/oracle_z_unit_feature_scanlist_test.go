// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func scanListTableName(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, gtime.TimestampMicro()%1e9)
}

func Test_Table_Relation_One(t *testing.T) {
	var (
		tableUser       = scanListTableName("user")
		tableUserDetail = scanListTableName("user_detail")
		tableUserScores = scanListTableName("user_scores")
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  name VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  address VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id NUMBER(10) NOT NULL,
  user_id NUMBER(10) NOT NULL,
  score NUMBER(10) NOT NULL,
  course VARCHAR2(45) NOT NULL,
  PRIMARY KEY (id)
)
    `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type EntityUser struct {
		UserId int    `orm:"user_id"`
		Name   string `orm:"name"`
	}

	type EntityUserDetail struct {
		UserId  int    `orm:"user_id"`
		Address string `orm:"address"`
	}

	type EntityUserScores struct {
		Id     int    `orm:"id"`
		UserId int    `orm:"user_id"`
		Score  int    `orm:"score"`
		Course string `orm:"course"`
	}

	type Entity struct {
		User       *EntityUser
		UserDetail *EntityUserDetail
		UserScores []*EntityUserScores
	}

	// Initialize the data.
	var err error
	gtest.C(t, func(t *gtest.T) {
		err = db.Transaction(context.TODO(), func(ctx context.Context, tx gdb.TX) error {
			r, err := tx.Model(tableUser).Save(EntityUser{
				UserId: 1,
				Name:   "john",
			})
			if err != nil {
				return err
			}
			lastInsertId, err := r.LastInsertId()
			if err != nil {
				return err
			}
			t.Assert(lastInsertId, 0)
			uid := int64(1)
			_, err = tx.Model(tableUserDetail).Save(EntityUserDetail{
				UserId:  int(uid),
				Address: "Beijing DongZhiMen #66",
			})
			if err != nil {
				return err
			}
			_, err = tx.Model(tableUserScores).Save(g.Slice{
				EntityUserScores{Id: 1, UserId: int(uid), Score: 100, Course: "math"},
				EntityUserScores{Id: 2, UserId: int(uid), Score: 99, Course: "physics"},
			})
			return err
		})
		t.AssertNil(err)
	})
	// Data check.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(tableUser).All()
		t.AssertNil(err)
		t.Assert(r.Len(), 1)
		t.Assert(r[0]["USER_ID"].Int(), 1)
		t.Assert(r[0]["NAME"].String(), "john")

		r, err = db.Model(tableUserDetail).Where("user_id", r[0]["USER_ID"].Int()).All()
		t.AssertNil(err)
		t.Assert(r.Len(), 1)
		t.Assert(r[0]["USER_ID"].Int(), 1)
		t.Assert(r[0]["ADDRESS"].String(), `Beijing DongZhiMen #66`)

		r, err = db.Model(tableUserScores).Where("user_id", r[0]["USER_ID"].Int()).All()
		t.AssertNil(err)
		t.Assert(r.Len(), 2)
		t.Assert(r[0]["USER_ID"].Int(), 1)
		t.Assert(r[1]["USER_ID"].Int(), 1)
		t.Assert(r[0]["COURSE"].String(), `math`)
		t.Assert(r[1]["COURSE"].String(), `physics`)
	})
	// Entity query.
	gtest.C(t, func(t *gtest.T) {
		var user Entity
		// SELECT * FROM `user` WHERE `name`='john'
		err := db.Model(tableUser).Scan(&user.User, "name", "john")
		t.AssertNil(err)

		// SELECT * FROM `user_detail` WHERE `user_id`=1
		err = db.Model(tableUserDetail).Scan(&user.UserDetail, "user_id", user.User.UserId)
		t.AssertNil(err)

		// SELECT * FROM `user_scores` WHERE `user_id`=1
		err = db.Model(tableUserScores).Scan(&user.UserScores, "user_id", user.User.UserId)
		t.AssertNil(err)
		t.Assert(len(user.UserScores), 2)

		t.Assert(user.User, EntityUser{
			UserId: 1,
			Name:   "john",
		})
		t.Assert(user.UserDetail, EntityUserDetail{
			UserId:  1,
			Address: "Beijing DongZhiMen #66",
		})
		t.Assert(user.UserScores, []EntityUserScores{
			{Id: 1, UserId: 1, Course: "math", Score: 100},
			{Id: 2, UserId: 1, Course: "physics", Score: 99},
		})
	})
}

func Test_Table_Relation_Many(t *testing.T) {
	var (
		tableUser       = scanListTableName("user")
		tableUserDetail = scanListTableName("user_detail")
		tableUserScores = scanListTableName("user_scores")
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  name VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  address VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id NUMBER(10) NOT NULL,
  user_id NUMBER(10) NOT NULL,
  score NUMBER(10) NOT NULL,
  PRIMARY KEY (id)
)
    `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	createAutoIncrement(tableUserScores, "id", 1)
	defer dropTable(tableUserScores)

	type EntityUser struct {
		UserId int    `json:"user_id"`
		Name   string `json:"name"`
	}
	type EntityUserDetail struct {
		UserId  int    `json:"user_id"`
		Address string `json:"address"`
	}
	type EntityUserScores struct {
		Id     int `json:"id"`
		UserId int `json:"user_id"`
		Score  int `json:"score"`
	}
	type Entity struct {
		User       *EntityUser
		UserDetail *EntityUserDetail
		UserScores []*EntityUserScores
	}

	// Initialize the data.
	gtest.C(t, func(t *gtest.T) {
		var err error
		for i := 1; i <= 5; i++ {
			// User.
			_, err = db.Insert(ctx, tableUser, g.Map{
				"user_id": i,
				"name":    fmt.Sprintf(`name_%d`, i),
			})
			t.AssertNil(err)
			// Detail.
			_, err = db.Insert(ctx, tableUserDetail, g.Map{
				"user_id": i,
				"address": fmt.Sprintf(`address_%d`, i),
			})
			t.AssertNil(err)
			// Scores.
			for j := 1; j <= 5; j++ {
				_, err = db.Insert(ctx, tableUserScores, g.Map{
					"user_id": i,
					"score":   j,
				})
				t.AssertNil(err)
			}
		}
	})

	// MapKeyValue.
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		t.Assert(all.Len(), 2)
		t.Assert(len(all.MapKeyValue("USER_ID")), 2)
		t.Assert(all.MapKeyValue("USER_ID")["3"].Map()["USER_ID"], 3)
		t.Assert(all.MapKeyValue("USER_ID")["4"].Map()["USER_ID"], 4)
		all, err = db.Model(tableUserScores).Where("user_id", g.Slice{3, 4}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(all.Len(), 10)
		t.Assert(len(all.MapKeyValue("USER_ID")), 2)
		t.Assert(len(all.MapKeyValue("USER_ID")["3"].Slice()), 5)
		t.Assert(len(all.MapKeyValue("USER_ID")["4"].Slice()), 5)
		t.Assert(gconv.Map(all.MapKeyValue("USER_ID")["3"].Slice()[0])["USER_ID"], 3)
		t.Assert(gconv.Map(all.MapKeyValue("USER_ID")["3"].Slice()[0])["SCORE"], 1)
		t.Assert(gconv.Map(all.MapKeyValue("USER_ID")["3"].Slice()[4])["USER_ID"], 3)
		t.Assert(gconv.Map(all.MapKeyValue("USER_ID")["3"].Slice()[4])["SCORE"], 5)
	})
	// Result ScanList with struct elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id:UserId")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "user_id:UserId")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id:UserId")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "user_id:UserId")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with struct elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			UserId int    `json:"user_id"`
			Name   string `json:"name"`
		}
		type EntityUserDetail struct {
			UserId  int    `json:"user_id"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id     int `json:"id"`
			UserId int `json:"user_id"`
			Score  int `json:"score"`
		}
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id:UserId")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "user_id:UserId")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with pointer elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			UserId int    `json:"user_id"`
			Name   string `json:"name"`
		}
		type EntityUserDetail struct {
			UserId  int    `json:"user_id"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id     int `json:"id"`
			UserId int `json:"user_id"`
			Score  int `json:"score"`
		}
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []*Entity

		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id:UserId")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "user_id:UserId")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Model ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		err := db.Model(tableUser).
			Where("user_id", g.Slice{3, 4}).
			Order("user_id asc").
			ScanList(&users, "User")
		t.AssertNil(err)
		// Detail
		err = db.Model(tableUserDetail).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("user_id asc").
			ScanList(&users, "UserDetail", "User", "user_id:UserId")
		t.AssertNil(err)
		// Scores
		err = db.Model(tableUserScores).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("id asc").
			ScanList(&users, "UserScores", "User", "user_id:UserId")
		t.AssertNil(err)

		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})

		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})

		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_Many_ModelScanList(t *testing.T) {
	var (
		tableUser       = scanListTableName("user")
		tableUserDetail = scanListTableName("user_detail")
		tableUserScores = scanListTableName("user_scores")
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  name VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  address VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id NUMBER(10) NOT NULL,
  user_id NUMBER(10) NOT NULL,
  score NUMBER(10) NOT NULL,
  PRIMARY KEY (id)
)
    `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	createAutoIncrement(tableUserScores, "id", 1)
	defer dropTable(tableUserScores)

	type EntityUser struct {
		UserId int    `json:"user_id"`
		Name   string `json:"name"`
	}
	type EntityUserDetail struct {
		UserId  int    `json:"user_id"`
		Address string `json:"address"`
	}
	type EntityUserScores struct {
		Id     int `json:"id"`
		UserId int `json:"user_id"`
		Score  int `json:"score"`
	}
	type Entity struct {
		User       *EntityUser
		UserDetail *EntityUserDetail
		UserScores []*EntityUserScores
	}

	// Initialize the data.
	gtest.C(t, func(t *gtest.T) {
		var err error
		for i := 1; i <= 5; i++ {
			// User.
			_, err = db.Insert(ctx, tableUser, g.Map{
				"user_id": i,
				"name":    fmt.Sprintf(`name_%d`, i),
			})
			t.AssertNil(err)
			// Detail.
			_, err = db.Insert(ctx, tableUserDetail, g.Map{
				"user_id": i,
				"address": fmt.Sprintf(`address_%d`, i),
			})
			t.AssertNil(err)
			// Scores.
			for j := 1; j <= 5; j++ {
				_, err = db.Insert(ctx, tableUserScores, g.Map{
					"user_id": i,
					"score":   j,
				})
				t.AssertNil(err)
			}
		}
	})

	//db.SetDebug(true)
	// Result ScanList with struct elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []Entity
		// User
		err := db.Model(tableUser).
			Where("user_id", g.Slice{3, 4}).
			Order("user_id asc").
			ScanList(&users, "User")
		t.AssertNil(err)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})

		// Detail
		err = db.Model(tableUserDetail).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("user_id asc").
			ScanList(&users, "UserDetail", "User", "user_id")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})

		// Scores
		err = db.Model(tableUserScores).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("id asc").
			ScanList(&users, "UserScores", "User", "user_id:UserId")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_Many_RelationKeyCaseInsensitive(t *testing.T) {
	var (
		tableUser       = scanListTableName("user")
		tableUserDetail = scanListTableName("user_detail")
		tableUserScores = scanListTableName("user_scores")
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  name VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  address VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id NUMBER(10) NOT NULL,
  user_id NUMBER(10) NOT NULL,
  score NUMBER(10) NOT NULL,
  PRIMARY KEY (id)
)
    `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	createAutoIncrement(tableUserScores, "id", 1)
	defer dropTable(tableUserScores)

	type EntityUser struct {
		UserId int    `json:"user_id"`
		Name   string `json:"name"`
	}
	type EntityUserDetail struct {
		UserId  int    `json:"user_id"`
		Address string `json:"address"`
	}
	type EntityUserScores struct {
		Id     int `json:"id"`
		UserId int `json:"user_id"`
		Score  int `json:"score"`
	}
	type Entity struct {
		User       *EntityUser
		UserDetail *EntityUserDetail
		UserScores []*EntityUserScores
	}

	// Initialize the data.
	gtest.C(t, func(t *gtest.T) {
		var err error
		for i := 1; i <= 5; i++ {
			// User.
			_, err = db.Insert(ctx, tableUser, g.Map{
				"user_id": i,
				"name":    fmt.Sprintf(`name_%d`, i),
			})
			t.AssertNil(err)
			// Detail.
			_, err = db.Insert(ctx, tableUserDetail, g.Map{
				"user_id": i,
				"address": fmt.Sprintf(`address_%d`, i),
			})
			t.AssertNil(err)
			// Scores.
			for j := 1; j <= 5; j++ {
				_, err = db.Insert(ctx, tableUserScores, g.Map{
					"user_id": i,
					"score":   j,
				})
				t.AssertNil(err)
			}
		}
	})

	// MapKeyValue.
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		t.Assert(all.Len(), 2)
		t.Assert(len(all.MapKeyValue("USER_ID")), 2)
		t.Assert(all.MapKeyValue("USER_ID")["3"].Map()["USER_ID"], 3)
		t.Assert(all.MapKeyValue("USER_ID")["4"].Map()["USER_ID"], 4)
		all, err = db.Model(tableUserScores).Where("user_id", g.Slice{3, 4}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(all.Len(), 10)
		t.Assert(len(all.MapKeyValue("USER_ID")), 2)
		t.Assert(len(all.MapKeyValue("USER_ID")["3"].Slice()), 5)
		t.Assert(len(all.MapKeyValue("USER_ID")["4"].Slice()), 5)
		t.Assert(gconv.Map(all.MapKeyValue("USER_ID")["3"].Slice()[0])["USER_ID"], 3)
		t.Assert(gconv.Map(all.MapKeyValue("USER_ID")["3"].Slice()[0])["SCORE"], 1)
		t.Assert(gconv.Map(all.MapKeyValue("USER_ID")["3"].Slice()[4])["USER_ID"], 3)
		t.Assert(gconv.Map(all.MapKeyValue("USER_ID")["3"].Slice()[4])["SCORE"], 5)
	})
	// Result ScanList with struct elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id:user_id")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "user_id:user_id")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "UserId:USER_ID")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "UserId:USER_ID")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with struct elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			UserId int    `json:"user_id"`
			Name   string `json:"name"`
		}
		type EntityUserDetail struct {
			UserId  int    `json:"user_id"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id     int `json:"id"`
			UserId int `json:"user_id"`
			Score  int `json:"score"`
		}
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id:UserID")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "UserID:UserId")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with pointer elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			UserId int    `json:"user_id"`
			Name   string `json:"name"`
		}
		type EntityUserDetail struct {
			UserId  int    `json:"user_id"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id     int `json:"id"`
			UserId int `json:"user_id"`
			Score  int `json:"score"`
		}
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []*Entity

		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id:UserId")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "User_Id:UserId")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Model ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		err := db.Model(tableUser).
			Where("user_id", g.Slice{3, 4}).
			Order("user_id asc").
			ScanList(&users, "User")
		t.AssertNil(err)
		// Detail
		err = db.Model(tableUserDetail).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("user_id asc").
			ScanList(&users, "UserDetail", "User", "user_id:UserId")
		t.AssertNil(err)
		// Scores
		err = db.Model(tableUserScores).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("id asc").
			ScanList(&users, "UserScores", "User", "user_id:UserId")
		t.AssertNil(err)

		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})

		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})

		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_Many_TheSameRelationNames(t *testing.T) {
	var (
		tableUser       = scanListTableName("user")
		tableUserDetail = scanListTableName("user_detail")
		tableUserScores = scanListTableName("user_scores")
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  name VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  address VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id NUMBER(10) NOT NULL,
  user_id NUMBER(10) NOT NULL,
  score NUMBER(10) NOT NULL,
  PRIMARY KEY (id)
)
    `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	createAutoIncrement(tableUserScores, "id", 1)
	defer dropTable(tableUserScores)

	type EntityUser struct {
		UserId int    `json:"user_id"`
		Name   string `json:"name"`
	}
	type EntityUserDetail struct {
		UserId  int    `json:"user_id"`
		Address string `json:"address"`
	}
	type EntityUserScores struct {
		Id     int `json:"id"`
		UserId int `json:"user_id"`
		Score  int `json:"score"`
	}
	type Entity struct {
		User       *EntityUser
		UserDetail *EntityUserDetail
		UserScores []*EntityUserScores
	}

	// Initialize the data.
	gtest.C(t, func(t *gtest.T) {
		var err error
		for i := 1; i <= 5; i++ {
			// User.
			_, err = db.Insert(ctx, tableUser, g.Map{
				"user_id": i,
				"name":    fmt.Sprintf(`name_%d`, i),
			})
			t.AssertNil(err)
			// Detail.
			_, err = db.Insert(ctx, tableUserDetail, g.Map{
				"user_id": i,
				"address": fmt.Sprintf(`address_%d`, i),
			})
			t.AssertNil(err)
			// Scores.
			for j := 1; j <= 5; j++ {
				_, err = db.Insert(ctx, tableUserScores, g.Map{
					"user_id": i,
					"score":   j,
				})
				t.AssertNil(err)
			}
		}
	})

	// Result ScanList with struct elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "user_id")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "UserId")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "USER_ID")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with struct elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			UserId int    `json:"user_id"`
			Name   string `json:"name"`
		}
		type EntityUserDetail struct {
			UserId  int    `json:"user_id"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id     int `json:"id"`
			UserId int `json:"user_id"`
			Score  int `json:"score"`
		}
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "UserID")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "UserId")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with pointer elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			UserId int    `json:"user_id"`
			Name   string `json:"name"`
		}
		type EntityUserDetail struct {
			UserId  int    `json:"user_id"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id     int `json:"id"`
			UserId int `json:"user_id"`
			Score  int `json:"score"`
		}
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []*Entity

		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "USER_ID")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Model ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		err := db.Model(tableUser).
			Where("user_id", g.Slice{3, 4}).
			Order("user_id asc").
			ScanList(&users, "User")
		t.AssertNil(err)
		// Detail
		err = db.Model(tableUserDetail).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("user_id asc").
			ScanList(&users, "UserDetail", "User", "user_id")
		t.AssertNil(err)
		// Scores
		err = db.Model(tableUserScores).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("id asc").
			ScanList(&users, "UserScores", "User", "user_id")
		t.AssertNil(err)

		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})

		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})

		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
}

func Test_Table_Relation_EmptyData(t *testing.T) {
	var (
		tableUser       = scanListTableName("user")
		tableUserDetail = scanListTableName("user_detail")
		tableUserScores = scanListTableName("user_scores")
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  name VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  address VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id NUMBER(10) NOT NULL,
  user_id NUMBER(10) NOT NULL,
  score NUMBER(10) NOT NULL,
  PRIMARY KEY (id)
)
    `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type EntityUser struct {
		UserId int    `json:"user_id"`
		Name   string `json:"name"`
	}
	type EntityUserDetail struct {
		UserId  int    `json:"user_id"`
		Address string `json:"address"`
	}
	type EntityUserScores struct {
		Id     int `json:"id"`
		UserId int `json:"user_id"`
		Score  int `json:"score"`
	}
	type Entity struct {
		User       *EntityUser
		UserDetail *EntityUserDetail
		UserScores []*EntityUserScores
	}

	// Result ScanList with struct elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 0)
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id:user_id")
		t.AssertNil(err)

		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "user_id:user_id")
		t.AssertNil(err)
	})
	// Result ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 0)

		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "UserId:USER_ID")
		t.AssertNil(err)

		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "UserId:USER_ID")
		t.AssertNil(err)
	})

	// Result ScanList with struct elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			UserId int    `json:"user_id"`
			Name   string `json:"name"`
		}
		type EntityUserDetail struct {
			UserId  int    `json:"user_id"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id     int `json:"id"`
			UserId int `json:"user_id"`
			Score  int `json:"score"`
		}
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)

		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id:UserID")
		t.AssertNil(err)

		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "UserID:UserId")
		t.AssertNil(err)
	})

	// Result ScanList with pointer elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			UserId int    `json:"user_id"`
			Name   string `json:"name"`
		}
		type EntityUserDetail struct {
			UserId  int    `json:"user_id"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id     int `json:"id"`
			UserId int `json:"user_id"`
			Score  int `json:"score"`
		}
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []*Entity

		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 0)
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id:UserId")
		t.AssertNil(err)

		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "User_Id:UserId")
		t.AssertNil(err)
	})

	// Model ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		err := db.Model(tableUser).
			Where("user_id", g.Slice{3, 4}).
			Order("user_id asc").
			ScanList(&users, "User")
		t.AssertNil(err)
		// Detail
		err = db.Model(tableUserDetail).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("user_id asc").
			ScanList(&users, "UserDetail", "User", "user_id:UserId")
		t.AssertNil(err)
		// Scores
		err = db.Model(tableUserScores).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("id asc").
			ScanList(&users, "UserScores", "User", "user_id:UserId")
		t.AssertNil(err)

		t.Assert(len(users), 0)
	})
}

func Test_Table_Relation_NoneEqualDataSize(t *testing.T) {
	var (
		tableUser       = scanListTableName("user")
		tableUserDetail = scanListTableName("user_detail")
		tableUserScores = scanListTableName("user_scores")
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  name VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  address VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id NUMBER(10) NOT NULL,
  user_id NUMBER(10) NOT NULL,
  score NUMBER(10) NOT NULL,
  PRIMARY KEY (id)
)
    `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type EntityUser struct {
		UserId int    `json:"user_id"`
		Name   string `json:"name"`
	}
	type EntityUserDetail struct {
		UserId  int    `json:"user_id"`
		Address string `json:"address"`
	}
	type EntityUserScores struct {
		Id     int `json:"id"`
		UserId int `json:"user_id"`
		Score  int `json:"score"`
	}
	type Entity struct {
		User       *EntityUser
		UserDetail *EntityUserDetail
		UserScores []*EntityUserScores
	}

	// Initialize the data.
	gtest.C(t, func(t *gtest.T) {
		var err error
		for i := 1; i <= 5; i++ {
			// User.
			_, err = db.Insert(ctx, tableUser, g.Map{
				"user_id": i,
				"name":    fmt.Sprintf(`name_%d`, i),
			})
			t.AssertNil(err)
			// Detail.
			// _, err = db.Insert(ctx, tableUserDetail, g.Map{
			//	"user_id":     i,
			//	"address": fmt.Sprintf(`address_%d`, i),
			// })
			// t.AssertNil(err)
			// Scores.
			// for j := 1; j <= 5; j++ {
			//	_, err = db.Insert(ctx, tableUserScores, g.Map{
			//		"user_id":   i,
			//		"score": j,
			//	})
			//	t.AssertNil(err)
			// }
		}
	})

	// Result ScanList with struct elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, nil)
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "user_id")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 0)
	})

	// Result ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "UserId")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, nil)
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "USER_ID")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 0)
	})

	// Result ScanList with struct elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			UserId int    `json:"user_id"`
			Name   string `json:"name"`
		}
		type EntityUserDetail struct {
			UserId  int    `json:"user_id"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id     int `json:"id"`
			UserId int `json:"user_id"`
			Score  int `json:"score"`
		}
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []Entity
		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "UserID")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, EntityUserDetail{})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "UserId")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 0)
	})

	// Result ScanList with pointer elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			UserId int    `json:"user_id"`
			Name   string `json:"name"`
		}
		type EntityUserDetail struct {
			UserId  int    `json:"user_id"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id     int `json:"id"`
			UserId int `json:"user_id"`
			Score  int `json:"score"`
		}
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []*Entity

		// User
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "User")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})
		// Detail
		all, err = db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "User", "user_id")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, EntityUserDetail{})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "User", "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "User", "USER_ID")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 0)
	})

	// Model ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		err := db.Model(tableUser).
			Where("user_id", g.Slice{3, 4}).
			Order("user_id asc").
			ScanList(&users, "User")
		t.AssertNil(err)
		// Detail
		err = db.Model(tableUserDetail).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("user_id asc").
			ScanList(&users, "UserDetail", "User", "user_id")
		t.AssertNil(err)
		// Scores
		err = db.Model(tableUserScores).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("id asc").
			ScanList(&users, "UserScores", "User", "user_id")
		t.AssertNil(err)

		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, "name_3"})
		t.Assert(users[1].User, &EntityUser{4, "name_4"})

		t.Assert(users[0].UserDetail, nil)

		t.Assert(len(users[0].UserScores), 0)
	})
}

func Test_Table_Relation_EmbeddedStruct1(t *testing.T) {
	var (
		tableUser       = scanListTableName("user")
		tableUserDetail = scanListTableName("user_detail")
		tableUserScores = scanListTableName("user_scores")
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  name VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  address VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id NUMBER(10) NOT NULL,
  user_id NUMBER(10) NOT NULL,
  score NUMBER(10) NOT NULL,
  PRIMARY KEY (id)
)
    `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	createAutoIncrement(tableUserScores, "id", 1)
	defer dropTable(tableUserScores)

	type EntityUser struct {
		UserId int    `json:"user_id"`
		Name   string `json:"name"`
	}
	type EntityUserDetail struct {
		*EntityUser
		UserId  int    `json:"user_id"`
		Address string `json:"address"`
	}
	type EntityUserScores struct {
		*EntityUser
		*EntityUserDetail
		Id     int `json:"id"`
		UserId int `json:"user_id"`
		Score  int `json:"score"`
	}

	// Initialize the data.
	gtest.C(t, func(t *gtest.T) {
		var err error
		for i := 1; i <= 5; i++ {
			// User.
			_, err = db.Insert(ctx, tableUser, g.Map{
				"user_id": i,
				"name":    fmt.Sprintf(`name_%d`, i),
			})
			t.AssertNil(err)
			// Detail.
			_, err = db.Insert(ctx, tableUserDetail, g.Map{
				"user_id": i,
				"address": fmt.Sprintf(`address_%d`, i),
			})
			t.AssertNil(err)
			// Scores.
			for j := 1; j <= 5; j++ {
				_, err = db.Insert(ctx, tableUserScores, g.Map{
					"user_id": i,
					"score":   j,
				})
				t.AssertNil(err)
			}
		}
	})

	gtest.C(t, func(t *gtest.T) {
		var (
			err    error
			scores []*EntityUserScores
		)
		// SELECT * FROM `user_scores`
		err = db.Model(tableUserScores).Scan(&scores)
		t.AssertNil(err)

		// SELECT * FROM `user_scores` WHERE `user_id` IN(1,2,3,4,5)
		err = db.Model(tableUser).
			Where("user_id", gdb.ListItemValuesUnique(&scores, "UserId")).
			ScanList(&scores, "EntityUser", "user_id:UserId")
		t.AssertNil(err)

		// SELECT * FROM `user_detail` WHERE `user_id` IN(1,2,3,4,5)
		err = db.Model(tableUserDetail).
			Where("user_id", gdb.ListItemValuesUnique(&scores, "UserId")).
			ScanList(&scores, "EntityUserDetail", "user_id:UserId")
		t.AssertNil(err)

		// Assertions.
		t.Assert(len(scores), 25)
		t.Assert(scores[0].Id, 1)
		t.Assert(scores[0].UserId, 1)
		t.Assert(scores[0].Name, "name_1")
		t.Assert(scores[0].Address, "address_1")
		t.Assert(scores[24].Id, 25)
		t.Assert(scores[24].UserId, 5)
		t.Assert(scores[24].Name, "name_5")
		t.Assert(scores[24].Address, "address_5")
	})
}

func Test_Table_Relation_EmbeddedStruct2(t *testing.T) {
	var (
		tableUser       = scanListTableName("user")
		tableUserDetail = scanListTableName("user_detail")
		tableUserScores = scanListTableName("user_scores")
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  name VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  address VARCHAR2(45) NOT NULL,
  PRIMARY KEY (user_id)
)
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id NUMBER(10) NOT NULL,
  user_id NUMBER(10) NOT NULL,
  score NUMBER(10) NOT NULL,
  PRIMARY KEY (id)
)
    `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	createAutoIncrement(tableUserScores, "id", 1)
	defer dropTable(tableUserScores)

	type EntityUser struct {
		UserId int    `json:"user_id"`
		Name   string `json:"name"`
	}
	type EntityUserDetail struct {
		UserId  int    `json:"user_id"`
		Address string `json:"address"`
	}
	type EntityUserScores struct {
		Id     int `json:"id"`
		UserId int `json:"user_id"`
		Score  int `json:"score"`
	}
	type Entity struct {
		*EntityUser
		UserDetail *EntityUserDetail
		UserScores []*EntityUserScores
	}

	// Initialize the data.
	gtest.C(t, func(t *gtest.T) {
		var err error
		for i := 1; i <= 5; i++ {
			// User.
			_, err = db.Insert(ctx, tableUser, g.Map{
				"user_id": i,
				"name":    fmt.Sprintf(`name_%d`, i),
			})
			t.AssertNil(err)
			// Detail.
			_, err = db.Insert(ctx, tableUserDetail, g.Map{
				"user_id": i,
				"address": fmt.Sprintf(`address_%d`, i),
			})
			t.AssertNil(err)
			// Scores.
			for j := 1; j <= 5; j++ {
				_, err = db.Insert(ctx, tableUserScores, g.Map{
					"user_id": i,
					"score":   j,
				})
				t.AssertNil(err)
			}
		}
	})

	// MapKeyValue.
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").All()
		t.AssertNil(err)
		t.Assert(all.Len(), 2)
		t.Assert(len(all.MapKeyValue("USER_ID")), 2)
		t.Assert(all.MapKeyValue("USER_ID")["3"].Map()["USER_ID"], 3)
		t.Assert(all.MapKeyValue("USER_ID")["4"].Map()["USER_ID"], 4)
		all, err = db.Model(tableUserScores).Where("user_id", g.Slice{3, 4}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(all.Len(), 10)
		t.Assert(len(all.MapKeyValue("USER_ID")), 2)
		t.Assert(len(all.MapKeyValue("USER_ID")["3"].Slice()), 5)
		t.Assert(len(all.MapKeyValue("USER_ID")["4"].Slice()), 5)
		t.Assert(gconv.Map(all.MapKeyValue("USER_ID")["3"].Slice()[0])["USER_ID"], 3)
		t.Assert(gconv.Map(all.MapKeyValue("USER_ID")["3"].Slice()[0])["SCORE"], 1)
		t.Assert(gconv.Map(all.MapKeyValue("USER_ID")["3"].Slice()[4])["USER_ID"], 3)
		t.Assert(gconv.Map(all.MapKeyValue("USER_ID")["3"].Slice()[4])["SCORE"], 5)
	})

	// Result ScanList with struct elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []Entity
		// User
		err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].EntityUser, &EntityUser{3, "name_3"})
		t.Assert(users[1].EntityUser, &EntityUser{4, "name_4"})
		// Detail
		all, err := db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "user_id")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "user_id")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].EntityUser, &EntityUser{3, "name_3"})
		t.Assert(users[1].EntityUser, &EntityUser{4, "name_4"})
		// Detail
		all, err := db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "user_id")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "user_id")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with struct elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			UserId int    `json:"user_id"`
			Name   string `json:"name"`
		}
		type EntityUserDetail struct {
			UserId  int    `json:"user_id"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id     int `json:"id"`
			UserId int `json:"user_id"`
			Score  int `json:"score"`
		}
		type Entity struct {
			EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []Entity
		// User
		err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].EntityUser, &EntityUser{3, "name_3"})
		t.Assert(users[1].EntityUser, &EntityUser{4, "name_4"})
		// Detail
		all, err := db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "user_id")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "user_id")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Result ScanList with pointer elements and struct attributes.
	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			UserId int    `json:"user_id"`
			Name   string `json:"name"`
		}
		type EntityUserDetail struct {
			UserId  int    `json:"user_id"`
			Address string `json:"address"`
		}
		type EntityUserScores struct {
			Id     int `json:"id"`
			UserId int `json:"user_id"`
			Score  int `json:"score"`
		}
		type Entity struct {
			EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}
		var users []*Entity

		// User
		err := db.Model(tableUser).Where("user_id", g.Slice{3, 4}).Order("user_id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].EntityUser, &EntityUser{3, "name_3"})
		t.Assert(users[1].EntityUser, &EntityUser{4, "name_4"})
		// Detail
		all, err := db.Model(tableUserDetail).Where("user_id", gdb.ListItemValues(users, "UserId")).Order("user_id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserDetail", "user_id")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		all, err = db.Model(tableUserScores).Where("user_id", gdb.ListItemValues(users, "UserId")).Order("id asc").All()
		t.AssertNil(err)
		err = all.ScanList(&users, "UserScores", "user_id")
		t.AssertNil(err)
		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})

	// Model ScanList with pointer elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []*Entity
		// User
		err := db.Model(tableUser).
			Where("user_id", g.Slice{3, 4}).
			Order("user_id asc").
			Scan(&users)
		t.AssertNil(err)
		// Detail
		err = db.Model(tableUserDetail).
			Where("user_id", gdb.ListItemValues(users, "UserId")).
			Order("user_id asc").
			ScanList(&users, "UserDetail", "user_id:UserId")
		t.AssertNil(err)
		// Scores
		err = db.Model(tableUserScores).
			Where("user_id", gdb.ListItemValues(users, "UserId")).
			Order("id asc").
			ScanList(&users, "UserScores", "user_id:UserId")
		t.AssertNil(err)

		t.Assert(len(users), 2)
		t.Assert(users[0].EntityUser, &EntityUser{3, "name_3"})
		t.Assert(users[1].EntityUser, &EntityUser{4, "name_4"})

		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})

		t.Assert(len(users[0].UserScores), 5)
		t.Assert(len(users[1].UserScores), 5)
		t.Assert(users[0].UserScores[0].UserId, 3)
		t.Assert(users[0].UserScores[0].Score, 1)
		t.Assert(users[0].UserScores[4].Score, 5)
		t.Assert(users[1].UserScores[0].UserId, 4)
		t.Assert(users[1].UserScores[0].Score, 1)
		t.Assert(users[1].UserScores[4].Score, 5)
	})
}
