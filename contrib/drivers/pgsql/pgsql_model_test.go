// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
	"testing"
)

func Test_Model_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		result, err = db.Model(table).Data(g.Map{
			"id":          "2",
			"uid":         "2",
			"passport":    "t2",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		type User struct {
			Id         int         `gconv:"id"`
			Uid        int         `gconv:"uid"`
			Passport   string      `json:"passport"`
			Password   string      `gconv:"password"`
			Nickname   string      `gconv:"nickname"`
			CreateTime *gtime.Time `json:"create_time"`
		}
		// Model inserting.
		result, err = db.Model(table).Data(User{
			Id:       3,
			Uid:      3,
			Passport: "t3",
			Password: "25d55ad283aa400af464c76d713c07ad",
			Nickname: "name_3",
		}).Insert()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)
		value, err := db.Model(table).Fields("passport").Where("id=3").Value()
		t.AssertNil(err)
		t.Assert(value.String(), "t3")

		result, err = db.Model(table).Data(&User{
			Id:         4,
			Uid:        4,
			Passport:   "t4",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "T4",
			CreateTime: gtime.Now(),
		}).Insert()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)
		value, err = db.Model(table).Fields("passport").Where("id=4").Value()
		t.AssertNil(err)
		t.Assert(value.String(), "t4")

		result, err = db.Model(table).Where("id>?", 1).Delete()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 3)
	})
}

func Test_Model_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", "2").Delete()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_Model_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data("passport", "user_22").Where("passport=?", "user_2").Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data("passport", "user_2").Where("passport='user_22'").Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})

	// Update + Data(string)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data("passport='user_33'").Where("passport='user_3'").Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})

	// Update + Fields(string)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields("passport").Data(g.Map{
			"passport": "user_44",
			"none":     "none",
		}).Where("passport='user_4'").Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_Model_Where(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? and nickname=?", 3, "name_3").One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})

	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Slice{"id", 3}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Slice{"id", 3, "nickname", "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})

	// slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	
	// map like
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"passport like": "user_1%",
		}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0].GMap().Get("id"), 1)
		t.Assert(result[1].GMap().Get("id"), 10)
	})

	// map + slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).Where("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=3", g.Slice{}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=?", g.Slice{3}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 3).Where("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 3).Where("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 30).WhereOr("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 30).WhereOr("nickname", "name_3").Where("id>?", 1).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 30).WhereOr("nickname", "name_3").Where("id>", 1).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})

	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? AND nickname=?", g.Slice{3, "name_3"}...).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? AND nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("passport like ? and nickname like ?", g.Slice{"user_3", "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})

	// map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{"id": 3, "nickname": "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})

	// map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{"id>": 1, "id<": 3}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 2)
	})

	// gmap.Map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})

	// gmap.Map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 2)
	})

	// list map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewListMapFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})

	// list map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewListMapFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 2)
	})

	// tree map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})

	// tree map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 2)
	})

	// struct, automatic mapping and filtering.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Nickname string
		}
		result, err := db.Model(table).Where(User{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)

		result, err = db.Model(table).Where(&User{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// slice single
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id IN(?)", g.Slice{1, 3}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 3)
	})
	// slice + string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("nickname=? AND id IN(?)", "name_3", g.Slice{1, 3}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
	})
	// slice + map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"id":       g.Slice{1, 3},
			"nickname": "name_3",
		}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
	})
	// slice + struct
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Ids      []int  `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := db.Model(table).Where(User{
			Ids:      []int{1, 3},
			Nickname: "name_3",
		}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
	})
}
