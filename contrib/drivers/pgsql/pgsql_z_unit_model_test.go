// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
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
			Id:         3,
			Uid:        3,
			Passport:   "t3",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "name_3",
			CreateTime: gtime.Now(),
		}).Insert()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)
		value, err := db.Model(table).Fields("passport").Where("id=3").Value() // model value
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

		result, err = db.Model(table).Where("id>?", 1).Delete() // model delete
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 3)
	})
}

func Test_Model_One(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			Nickname   string
			CreateTime string
		}
		data := User{
			Id:         1,
			Passport:   "user_1",
			Password:   "pass_1",
			Nickname:   "name_1",
			CreateTime: "2020-10-10 12:00:01",
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One() // model one
		t.AssertNil(err)
		t.Assert(one["passport"], data.Passport)
		t.Assert(one["create_time"], data.CreateTime)
		t.Assert(one["nickname"], data.Nickname)
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

func Test_Model_Array(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id", g.Slice{1, 2, 3}).All()
		t.AssertNil(err)
		t.Assert(all.Array("id"), g.Slice{1, 2, 3})
		t.Assert(all.Array("nickname"), g.Slice{"name_1", "name_2", "name_3"})
	})
	gtest.C(t, func(t *gtest.T) {
		array, err := db.Model(table).Fields("nickname").Where("id", g.Slice{1, 2, 3}).Array()
		t.AssertNil(err)
		t.Assert(array, g.Slice{"name_1", "name_2", "name_3"})
	})
	gtest.C(t, func(t *gtest.T) {
		array, err := db.Model(table).Array("nickname", "id", g.Slice{1, 2, 3})
		t.AssertNil(err)
		t.Assert(array, g.Slice{"name_1", "name_2", "name_3"})
	})
}

func Test_Model_Scan(t *testing.T) {
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

func Test_Model_Count(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).FieldsEx("id").Where("id>8").Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))
	})
}

func Test_Model_Where(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

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
}

func Test_Model_Save(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var (
			user   User
			count  int
			result sql.Result
			err    error
		)

		result, err = db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "p1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": CreateTime,
		}).OnConflict("id").Save()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		err = db.Model(table).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 1)
		t.Assert(user.Passport, "p1")
		t.Assert(user.Password, "pw1")
		t.Assert(user.NickName, "n1")
		t.Assert(user.CreateTime.String(), CreateTime)

		_, err = db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "p1",
			"password":    "pw2",
			"nickname":    "n2",
			"create_time": CreateTime,
		}).OnConflict("id").Save()
		t.AssertNil(err)

		err = db.Model(table).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Passport, "p1")
		t.Assert(user.Password, "pw2")
		t.Assert(user.NickName, "n2")
		t.Assert(user.CreateTime.String(), CreateTime)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

func Test_Model_Replace(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "t11",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T11",
			"create_time": "2018-10-24 10:00:00",
		}).Replace()
		t.Assert(err, "Replace operation is not supported by pgsql driver")
	})
}

func Test_Model_OnConflict(t *testing.T) {
	var (
		table      = fmt.Sprintf(`%s_%d`, TablePrefix+"test", gtime.TimestampNano())
		uniqueName = fmt.Sprintf(`%s_%d`, TablePrefix+"test_unique", gtime.TimestampNano())
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		   	id bigserial  NOT NULL,
		   	passport varchar(45) NOT NULL,
		   	password varchar(32) NOT NULL,
		   	nickname varchar(45) NOT NULL,
		   	create_time timestamp NOT NULL,
		   	PRIMARY KEY (id),
			CONSTRAINT %s UNIQUE ("passport", "password")
		) ;`, table, uniqueName,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	// string type 1.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict("passport,password").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "n1")
	})

	// string type 2.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict("passport", "password").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "n1")
	})

	// slice.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict(g.Slice{"passport", "password"}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "n1")
	})
}

func Test_Model_OnDuplicate(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// string type 1.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict("id").OnDuplicate("passport,password").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})

	// string type 2.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict("id").OnDuplicate("passport", "password").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})

	// slice.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict("id").OnDuplicate(g.Slice{"passport", "password"}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})

	// map.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict("id").OnDuplicate(g.Map{
			"passport": "nickname",
			"password": "nickname",
		}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["nickname"])
		t.Assert(one["password"], data["nickname"])
		t.Assert(one["nickname"], "name_1")
	})

	// map+raw.
	gtest.C(t, func(t *gtest.T) {
		data := g.MapStrStr{
			"id":          "1",
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict("id").OnDuplicate(g.Map{
			"passport": gdb.Raw("CONCAT(EXCLUDED.passport, '1')"),
			"password": gdb.Raw("CONCAT(EXCLUDED.password, '2')"),
		}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"]+"1")
		t.Assert(one["password"], data["password"]+"2")
		t.Assert(one["nickname"], "name_1")
	})
}

func Test_Model_OnDuplicateEx(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// string type 1.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict("id").OnDuplicateEx("nickname,create_time").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})

	// string type 2.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict("id").OnDuplicateEx("nickname", "create_time").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})

	// slice.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict("id").OnDuplicateEx(g.Slice{"nickname", "create_time"}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})

	// map.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnConflict("id").OnDuplicateEx(g.Map{
			"nickname":    "nickname",
			"create_time": "nickname",
		}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})
}
