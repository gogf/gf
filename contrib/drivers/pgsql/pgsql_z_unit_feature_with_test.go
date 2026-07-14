// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gmeta"
)

func Test_Table_Relation_With_Scan(t *testing.T) {
	var (
		tableUser       = "with_scan_user"
		tableUserDetail = "with_scan_user_detail"
		tableUserScores = "with_scan_user_score"
	)
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user.sql"), tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_detail.sql"), tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_scores.sql"), tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_scan_user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScore struct {
		gmeta.Meta `orm:"table:with_scan_user_score"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:with_scan_user"`
		Id         int          `json:"id"`
		Name       string       `json:"name"`
		UserDetail *UserDetail  `orm:"with:uid=id"`
		UserScores []*UserScore `orm:"with:uid=id"`
	}

	// Initialize the data.
	gtest.C(t, func(t *gtest.T) {
		for i := 1; i <= 5; i++ {
			// User.
			user := User{
				Name: fmt.Sprintf(`name_%d`, i),
			}
			lastInsertId, err := db.Model(tableUser).Data(user).OmitEmpty().InsertAndGetId()
			t.AssertNil(err)
			// Detail.
			userDetail := UserDetail{
				Uid:     int(lastInsertId),
				Address: fmt.Sprintf(`address_%d`, lastInsertId),
			}
			_, err = db.Model(tableUserDetail).Data(userDetail).OmitEmpty().Insert()
			t.AssertNil(err)
			// Scores.
			for j := 1; j <= 5; j++ {
				userScore := UserScore{
					Uid:   int(lastInsertId),
					Score: j,
				}
				_, err = db.Model(tableUserScores).Data(userScore).OmitEmpty().Insert()
				t.AssertNil(err)
			}
		}
	})

	// Scan pointer.
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(tableUser).
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

	// Scan struct.
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).
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

	// With part attribute: UserDetail.
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(tableUser).
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
		err := db.Model(tableUser).
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

func Test_Table_Relation_With(t *testing.T) {
	var (
		tableUser       = "with_rel_user"
		tableUserDetail = "with_rel_user_detail"
		tableUserScores = "with_rel_user_scores"
	)
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user.sql"), tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_detail.sql"), tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_scores.sql"), tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_rel_user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:with_rel_user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:with_rel_user"`
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
				"uid":   i,
				"score": j,
			})
			gtest.AssertNil(err)
		}
	}

	gtest.C(t, func(t *gtest.T) {
		var users []*User
		err := db.Model(tableUser).
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
		err := db.Model(tableUser).
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
		err := db.Model(tableUser).
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

func Test_Table_Relation_WithAll(t *testing.T) {
	var (
		tableUser       = "withall_user"
		tableUserDetail = "withall_user_detail"
		tableUserScores = "withall_user_scores"
	)
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user.sql"), tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_detail.sql"), tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_scores.sql"), tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:withall_user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:withall_user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:withall_user"`
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
				"uid":   i,
				"score": j,
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

func Test_Table_Relation_WithAll_List(t *testing.T) {
	var (
		tableUser       = "withall_list_user"
		tableUserDetail = "withall_list_user_detail"
		tableUserScores = "withall_list_user_scores"
	)
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user.sql"), tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_detail.sql"), tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_scores.sql"), tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:withall_list_user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:withall_list_user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:withall_list_user"`
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
				"uid":   i,
				"score": j,
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

func Test_Table_Relation_WithAllCondition_List(t *testing.T) {
	var (
		tableUser       = "withall_cond_user"
		tableUserDetail = "withall_cond_user_detail"
		tableUserScores = "withall_cond_user_scores"
	)
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user.sql"), tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_detail.sql"), tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_scores.sql"), tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:withall_cond_user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:withall_cond_user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:withall_cond_user"`
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
				"uid":   i,
				"score": j,
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
		t.Assert(users[1].UserDetail.Uid, 4)
		t.Assert(users[1].UserDetail.Address, "address_4")
		t.Assert(len(users[1].UserScores), 3)
		t.Assert(users[1].UserScores[0].Uid, 4)
		t.Assert(users[1].UserScores[0].Score, 4)
		t.Assert(users[1].UserScores[2].Uid, 4)
		t.Assert(users[1].UserScores[2].Score, 2)
	})
}

func Test_Table_Relation_WithAll_Embedded_With_SelfMaintained_Attributes(t *testing.T) {
	var (
		tableUser       = "withall_emsm_user"
		tableUserDetail = "withall_emsm_user_detail"
		tableUserScores = "withall_emsm_user_scores"
	)
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user.sql"), tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_detail.sql"), tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_scores.sql"), tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:withall_emsm_user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:withall_emsm_user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta  `orm:"table:withall_emsm_user"`
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
				"uid":   i,
				"score": j,
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

func Test_Table_Relation_WithAll_Embedded_Without_SelfMaintained_Attributes(t *testing.T) {
	var (
		tableUser       = "withall_emns_user"
		tableUserDetail = "withall_emns_user_detail"
		tableUserScores = "withall_emns_user_scores"
	)
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user.sql"), tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_detail.sql"), tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_scores.sql"), tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:withall_emns_user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:withall_emns_user_scores"`
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
		gmeta.Meta  `orm:"table:withall_emns_user"`
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
				"uid":   i,
				"score": j,
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

func Test_Table_Relation_WithAll_Embedded_WithoutMeta(t *testing.T) {
	var (
		tableUser       = "withall_nometa_user"
		tableUserDetail = "user_detail"
		tableUserScores = "user_scores"
	)
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user.sql"), tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_detail.sql"), tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_scores.sql"), tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

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
				"uid":   i,
				"score": j,
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

func Test_Table_Relation_WithAll_AttributeStructAlsoHasWithTag(t *testing.T) {
	var (
		tableUser       = "withall_nested_user"
		tableUserDetail = "withall_nested_user_detail"
		tableUserScores = "withall_nested_user_scores"
	)
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user.sql"), tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_detail.sql"), tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_scores.sql"), tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type UserScores struct {
		gmeta.Meta `orm:"table:withall_nested_user_scores"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type UserDetail struct {
		gmeta.Meta `orm:"table:withall_nested_user_detail"`
		Uid        int           `json:"uid"`
		Address    string        `json:"address"`
		UserScores []*UserScores `orm:"with:uid"`
	}

	type User struct {
		gmeta.Meta  `orm:"table:withall_nested_user"`
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
				"uid":   i,
				"score": j,
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

func Test_Table_Relation_With_MultipleDepends1(t *testing.T) {
	defer func() {
		dropTable("table_a")
		dropTable("table_b")
		dropTable("table_c")
	}()
	for _, v := range gstr.SplitAndTrim(gfile.GetContents(gtest.DataPath("with_multiple_depends.sql")), ";") {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}

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

	db.SetDebug(true)
	defer db.SetDebug(false)

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

func Test_Table_Relation_With_MultipleDepends2(t *testing.T) {
	defer func() {
		dropTable("table_a")
		dropTable("table_b")
		dropTable("table_c")
	}()
	for _, v := range gstr.SplitAndTrim(gfile.GetContents(gtest.DataPath("with_multiple_depends.sql")), ";") {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}

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

	db.SetDebug(true)
	defer db.SetDebug(false)

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

func Test_Table_Relation_With_MultipleDepends_Embedded(t *testing.T) {
	defer func() {
		dropTable("table_a")
		dropTable("table_b")
		dropTable("table_c")
	}()
	for _, v := range gstr.SplitAndTrim(gfile.GetContents(gtest.DataPath("with_multiple_depends.sql")), ";") {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}

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

	db.SetDebug(true)
	defer db.SetDebug(false)

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

func Test_Table_Relation_WithAll_Embedded_Meta_NameMatchingRule(t *testing.T) {
	var (
		tableUser       = "with_embed_user"
		tableUserDetail = "with_embed_user_detail"
		tableUserScores = "with_embed_user_scores"
	)
	// Drop tables first to ensure clean state
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
id SERIAL PRIMARY KEY,
name varchar(45) NOT NULL
);
 `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
user_id SERIAL PRIMARY KEY,
address varchar(45) NOT NULL
);
 `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
id SERIAL PRIMARY KEY,
user_id integer NOT NULL,
score integer NOT NULL
);
 `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_embed_user_detail"`
		UserID     int    `json:"user_id"`
		Address    string `json:"address"`
	}

	type UserScores struct {
		gmeta.Meta `orm:"table:with_embed_user_scores"`
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
		gmeta.Meta `orm:"table:with_embed_user"`
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
		tableUser       = "with_unscoped_user"
		tableUserDetail = "with_unscoped_user_detail"
	)
	// Drop tables first to ensure clean state
	dropTable(tableUser)
	dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
id SERIAL PRIMARY KEY,
name varchar(45) NOT NULL
);
 `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
user_id SERIAL PRIMARY KEY,
address varchar(45) NOT NULL,
deleted_at timestamp default NULL
);
 `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_unscoped_user_detail"`
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
		gmeta.Meta `orm:"table:with_unscoped_user"`
		UserEmbedded
		UserDetail *UserDetail `orm:"with:user_id=id"`
	}
	type UserWithDeletedDetail struct {
		gmeta.Meta `orm:"table:with_unscoped_user"`
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
		tableUser       = "with_order_user"
		tableUserDetail = "with_order_user_detail"
	)
	// Drop tables first to ensure clean state
	dropTable(tableUser)
	dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
id SERIAL PRIMARY KEY,
name varchar(45) NOT NULL
);
 `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
user_id SERIAL PRIMARY KEY,
address varchar(45) NOT NULL,
deleted_at timestamp default NULL
);
 `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	type UserDetail struct {
		gmeta.Meta `orm:"table:with_order_user_detail"`
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
		gmeta.Meta `orm:"table:with_order_user"`
		UserEmbedded
		UserDetail *UserDetail `orm:"with:user_id=id"`
	}
	type UserWithDeletedDetail struct {
		gmeta.Meta `orm:"table:with_order_user"`
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
