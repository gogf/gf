// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
	"strings"
	"testing"
	"time"
)

func Test_Model_InnerJoin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createInitTable("user1")
		table2 := createInitTable("user2")

		defer dropTable(table1)
		defer dropTable(table2)

		res, err := db.Model(table1).Where("id > ?", 5).Delete()
		if err != nil {
			t.Fatal(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}

		t.Assert(n, 5)

		result, err := db.Model(table1+" u1").InnerJoin(table2+" u2", "u1.id = u2.id").Order("u1.id").All()
		if err != nil {
			t.Fatal(err)
		}

		t.Assert(len(result), 5)

		result, err = db.Model(table1+" u1").InnerJoin(table2+" u2", "u1.id = u2.id").Where("u1.id > ?", 1).Order("u1.id").All()
		if err != nil {
			t.Fatal(err)
		}

		t.Assert(len(result), 4)
	})
}

func Test_Model_LeftJoin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createInitTable("user1")
		table2 := createInitTable("user2")

		defer dropTable(table1)
		defer dropTable(table2)

		res, err := db.Model(table2).Where("id > ?", 3).Delete()
		if err != nil {
			t.Fatal(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		} else {
			t.Assert(n, 7)
		}

		result, err := db.Model(table1+" u1").LeftJoin(table2+" u2", "u1.id = u2.id").All()
		if err != nil {
			t.Fatal(err)
		}

		t.Assert(len(result), 10)

		result, err = db.Model(table1+" u1").LeftJoin(table2+" u2", "u1.id = u2.id").Where("u1.id > ? ", 2).All()
		if err != nil {
			t.Fatal(err)
		}

		t.Assert(len(result), 8)
	})
}

func Test_Model_RightJoin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createInitTable("user1")
		table2 := createInitTable("user2")

		defer dropTable(table1)
		defer dropTable(table2)

		res, err := db.Model(table1).Where("id > ?", 3).Delete()
		if err != nil {
			t.Fatal(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}

		t.Assert(n, 7)

		result, err := db.Model(table1+" u1").RightJoin(table2+" u2", "u1.id = u2.id").All()
		if err != nil {
			t.Fatal(err)
		}
		t.Assert(len(result), 10)

		result, err = db.Model(table1+" u1").RightJoin(table2+" u2", "u1.id = u2.id").Where("u1.id > 2").All()
		if err != nil {
			t.Fatal(err)
		}
		t.Assert(len(result), 1)
	})
}

func TestPage(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	db.SetDebug(true)
	result, err := db.Model(table).Page(1, 2).Order("ID").All()
	gtest.Assert(err, nil)
	fmt.Println("page:1--------", result)
	gtest.Assert(len(result), 2)
	gtest.Assert(result[0]["ID"], 1)
	gtest.Assert(result[1]["ID"], 2)

	result, err = db.Model(table).Page(2, 2).Order("ID").All()
	gtest.Assert(err, nil)
	fmt.Println("page: 2--------", result)
	gtest.Assert(len(result), 2)
	gtest.Assert(result[0]["ID"], 3)
	gtest.Assert(result[1]["ID"], 4)

	result, err = db.Model(table).Page(3, 2).Order("ID").All()
	gtest.Assert(err, nil)
	fmt.Println("page:3 --------", result)
	gtest.Assert(len(result), 2)
	gtest.Assert(result[0]["ID"], 5)

	result, err = db.Model(table).Page(2, 3).All()
	gtest.Assert(err, nil)
	gtest.Assert(len(result), 3)
	gtest.Assert(result[0]["ID"], 4)
	gtest.Assert(result[1]["ID"], 5)
	gtest.Assert(result[2]["ID"], 6)
}

func Test_Model_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	//db.SetDebug(true)
	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"ID":          1,
			"UID":         1,
			"PASSPORT":    "t1",
			"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
			"NICKNAME":    "name_1",
			"CREATE_TIME": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)

		result, err = db.Model(table).Data(g.Map{
			"ID":          "2",
			"UID":         "2",
			"PASSPORT":    "t2",
			"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
			"NICKNAME":    "name_2",
			"CREATE_TIME": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)

		type User struct {
			Id         int         `gconv:"ID"`
			Uid        int         `gconv:"uid"`
			Passport   string      `json:"PASSPORT"`
			Password   string      `gconv:"PASSWORD"`
			Nickname   string      `gconv:"NICKNAME"`
			CreateTime *gtime.Time `json:"CREATE_TIME"`
		}
		// Model inserting.
		result, err = db.Model(table).Data(User{
			Id:         3,
			Uid:        3,
			Passport:   "t3",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "name_3",
			CreateTime: gtime.Now(),
		}).Insert()
		t.AssertNil(err)

		value, err := db.Model(table).Fields("PASSPORT").Where("id=3").Value()
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

		value, err = db.Model(table).Fields("PASSPORT").Where("id=4").Value()
		t.AssertNil(err)
		t.Assert(value.String(), "t4")

		result, err = db.Model(table).Where("id>?", 1).Delete()
		t.AssertNil(err)
		_, _ = result.RowsAffected()

	})
}

func Test_Model_Insert_Time(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"ID":          1,
			"PASSPORT":    "t1",
			"PASSWORD":    "p1",
			"NICKNAME":    "n1",
			"CREATE_TIME": "2020-10-10 20:09:18",
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).One("ID", 1)
		t.AssertNil(err)
		t.Assert(one["PASSPORT"].String(), data["PASSPORT"])
		t.Assert(one["CREATE_TIME"].String(), "2020-10-10 20:09:18")
		t.Assert(one["NICKNAME"].String(), data["NICKNAME"])
	})
}

func Test_Model_Batch(t *testing.T) {
	// batch insert
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		_, err := db.Model(table).Data(g.List{
			{
				"ID":          2,
				"uid":         2,
				"PASSPORT":    "t2",
				"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
				"NICKNAME":    "name_2",
				"CREATE_TIME": gtime.Now().String(),
			},
			{
				"ID":          3,
				"uid":         3,
				"PASSPORT":    "t3",
				"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
				"NICKNAME":    "name_3",
				"CREATE_TIME": gtime.Now().String(),
			},
		}).Batch(1).Insert()
		if err != nil {
			gtest.Error(err)
		}
	})

}

func Test_Model_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data("PASSPORT", "user_22").Where("passport=?", "user_2").Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data("PASSPORT", "user_2").Where("passport='user_22'").Update()
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
		result, err := db.Model(table).Fields("PASSPORT").Data(g.Map{
			"PASSPORT": "user_44",
			"none":     "none",
		}).Where("passport='user_4'").Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_Model_All(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id<0").All()
		t.Assert(result, nil)
		t.AssertNil(err)
	})
}

func Test_Model_One(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		record, err := db.Model(table).Where("ID", 1).One()
		t.AssertNil(err)
		t.Assert(record["NICKNAME"].String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		record, err := db.Model(table).Where("ID", 0).One()
		t.AssertNil(err)
		t.Assert(record, nil)
	})
}

func Test_Model_Value(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Fields("NICKNAME").Where("ID", 1).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Fields("NICKNAME").Where("ID", 0).Value()
		t.AssertNil(err)
		t.Assert(value, nil)
	})
}

func Test_Model_Array(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("ID", g.Slice{1, 2, 3}).All()
		t.AssertNil(err)
		for k, v := range all.Array("ID") {
			t.Assert(v.Int(), k+1)
		}
		t.Assert(all.Array("NICKNAME"), g.Slice{"name_1", "name_2", "name_3"})
	})
	gtest.C(t, func(t *gtest.T) {
		array, err := db.Model(table).Fields("NICKNAME").Where("ID", g.Slice{1, 2, 3}).Array()
		t.AssertNil(err)
		t.Assert(array, g.Slice{"name_1", "name_2", "name_3"})
	})
	gtest.C(t, func(t *gtest.T) {
		array, err := db.Model(table).Array("NICKNAME", "ID", g.Slice{1, 2, 3})
		t.AssertNil(err)
		t.Assert(array, g.Slice{"name_1", "name_2", "name_3"})
	})
}

func Test_Model_Count(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
	// Count with cache, check internal ctx data feature.
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 10; i++ {
			count, err := db.Model(table).Cache(gdb.CacheOption{
				Duration: time.Second * 10,
				Name:     guid.S(),
				Force:    false,
			}).Count()
			t.AssertNil(err)
			t.Assert(count, TableSize)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).FieldsEx("ID").Where("id>8").Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Fields("distinct id").Where("id>8").Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
	// COUNT...LIMIT...
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Page(1, 2).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})

}

func Test_Model_Select(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id         int
		Passport   string
		Password   string
		NickName   string
		CreateTime gtime.Time
	}
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.Model(table).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
	})
}

func Test_Model_Struct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
	})
	// Auto creating struct object.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := (*User)(nil)
		err := db.Model(table).Where("id=1").Scan(&user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
	})
	// Just using Scan.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := (*User)(nil)
		err := db.Model(table).Where("id=1").Scan(&user)
		if err != nil {
			gtest.Error(err)
		}
		t.Assert(user.NickName, "name_1")
	})
	// sql.ErrNoRows
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=-1").Scan(user)
		t.Assert(err, sql.ErrNoRows)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var user *User
		err := db.Model(table).Where("id=-1").Scan(&user)
		t.AssertNil(err)
	})
}

func Test_Model_Scan(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
	})
	// sql.ErrNoRows
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var (
			user  = new(User)
			users = new([]*User)
		)
		err1 := db.Model(table).Where("id < 0").Scan(user)
		err2 := db.Model(table).Where("id < 0").Scan(users)
		t.Assert(err1, sql.ErrNoRows)
		t.Assert(err2, nil)
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
		t.Assert(result["ID"].Int(), 3)
	})

	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Slice{"ID", 3}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Slice{"ID", 3, "NICKNAME", "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})

	// slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	// map like
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"passport like": "user_1%",
		}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0].GMap().Get("ID"), 1)
		t.Assert(result[1].GMap().Get("ID"), 10)
	})
	// map + slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"ID":       g.Slice{1, 2, 3},
			"PASSPORT": g.Slice{"user_2", "user_3"},
		}).Where("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=3", g.Slice{}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=?", g.Slice{3}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("ID", 3).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("ID", 3).Where("NICKNAME", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("ID", 3).Where("NICKNAME", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("ID", 30).WhereOr("NICKNAME", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("ID", 30).WhereOr("NICKNAME", "name_3").Where("id>?", 1).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("ID", 30).WhereOr("NICKNAME", "name_3").Where("id>", 1).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? AND nickname=?", g.Slice{3, "name_3"}...).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? AND nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("passport like ? and nickname like ?", g.Slice{"user_3", "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{"ID": 3, "NICKNAME": "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{"id>": 1, "id<": 3}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 2)
	})

	// gmap.Map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewFrom(g.MapAnyAny{"ID": 3, "NICKNAME": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// gmap.Map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 2)
	})

	// list map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewListMapFrom(g.MapAnyAny{"ID": 3, "NICKNAME": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// list map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewListMapFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 2)
	})

	// tree map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"ID": 3, "NICKNAME": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// tree map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 2)
	})

	// complicated where 1
	gtest.C(t, func(t *gtest.T) {
		// db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id > 0":             nil,
			"ID":                 g.Slice{1, 2, 3},
		}
		result, err := db.Model(table).Where(conditions).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
	})
	// complicated where 2
	gtest.C(t, func(t *gtest.T) {
		// db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id >= ?":            1,
			"id in(?)":           g.Slice{1, 2, 3},
		}
		result, err := db.Model(table).Where(conditions).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
	})
	// struct, automatic mapping and filtering.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Nickname string
		}
		result, err := db.Model(table).Where(User{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)

		result, err = db.Model(table).Where(&User{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// slice single
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id IN(?)", g.Slice{1, 3}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[1]["ID"].Int(), 3)
	})
	// slice + string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("nickname=? AND id IN(?)", "name_3", g.Slice{1, 3}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 3)
	})
	// slice + map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"ID":       g.Slice{1, 3},
			"NICKNAME": "name_3",
		}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 3)
	})
	// slice + struct
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Ids      []int  `json:"ID"`
			Nickname string `gconv:"NICKNAME"`
		}
		result, err := db.Model(table).Where(User{
			Ids:      []int{1, 3},
			Nickname: "name_3",
		}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 3)
	})
}

func Test_Model_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, TableSize)
	})
}

func Test_Model_Having(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id, count(*)").Where("id > 1").Group("ID").Having("id > 8").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})

}

func Test_Model_Distinct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id > 1").Distinct().Count()
		t.AssertNil(err)
		t.Assert(count, 9)
	})
}

//not support
/*
func Test_Model_Min_Max(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table, "t").Fields("min(t.id)").Where("id > 1").Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table, "t").Fields("max(t.id)").Where("id > 1").Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 10)
	})
}
*/
func Test_Model_HasTable(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	//db.SetDebug(true)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetCore().HasTable(strings.ToUpper(table))
		t.Assert(result, true)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetCore().HasTable("table12321")
		t.Assert(result, false)
		t.AssertNil(err)
	})
}

func Test_Model_HasField(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).HasField("ID")
		t.Assert(result, true)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).HasField("id123")
		t.Assert(result, false)
		t.AssertNil(err)
	})
}

func Test_Model_WhereIn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereIn("ID", g.Slice{1, 2, 3, 4}).WhereIn("ID", g.Slice{3, 4, 5}).OrderAsc("ID").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["ID"], 3)
		t.Assert(result[1]["ID"], 4)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereIn("ID", g.Slice{}).OrderAsc("ID").All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).OmitEmptyWhere().WhereIn("ID", g.Slice{}).OrderAsc("ID").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
	})
}

func Test_Model_WhereNotIn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereNotIn("ID", g.Slice{1, 2, 3, 4}).WhereNotIn("ID", g.Slice{3, 4, 5}).OrderAsc("ID").All()
		t.AssertNil(err)
		t.Assert(len(result), 5)
		t.Assert(result[0]["ID"], 6)
		t.Assert(result[1]["ID"], 7)
	})
}

func Test_Model_WhereOrIn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrIn("ID", g.Slice{1, 2, 3, 4}).WhereOrIn("ID", g.Slice{3, 4, 5}).OrderAsc("ID").All()
		t.AssertNil(err)
		t.Assert(len(result), 5)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[4]["ID"], 5)
	})
}

func Test_Model_WhereLike(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereLike("NICKNAME", "name%").OrderAsc("ID").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[TableSize-1]["ID"], TableSize)
	})
}

func Test_Model_WhereNotLike(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereNotLike("NICKNAME", "name%").OrderAsc("ID").All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

func Test_Model_WhereOrLike(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrLike("NICKNAME", "namexxx%").WhereOrLike("NICKNAME", "name%").OrderAsc("ID").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[TableSize-1]["ID"], TableSize)
	})
}

func Test_Model_WhereOrNotLike(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrNotLike("NICKNAME", "namexxx%").WhereOrNotLike("NICKNAME", "name%").OrderAsc("ID").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[TableSize-1]["ID"], TableSize)
	})
}

/* not support the "AS"
func Test_Model_Raw(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.
			Raw(fmt.Sprintf("select * from %s where id in (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("ID", 8).
			WhereIn("ID", g.Slice{1, 2, 3, 4, 5, 6, 7}).
			OrderDesc("ID").
			All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["ID"], 7)
		t.Assert(all[1]["ID"], 5)
	})

	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("select * from %s where id in (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("ID", 8).
			WhereIn("ID", g.Slice{1, 2, 3, 4, 5, 6, 7}).
			OrderDesc("ID").
			Count()
		t.AssertNil(err)
		t.Assert(count, 6)
	})
}

func Test_Model_FieldCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("ID").FieldCount("ID", "total").Group("ID").OrderAsc("ID").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["ID"], 1)
		t.Assert(all[0]["total"], 1)
	})
}

func Test_Model_FieldMax(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("ID").FieldMax("ID", "total").Group("ID").OrderAsc("ID").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["ID"], 1)
		t.Assert(all[0]["total"], 1)
	})
}*/
