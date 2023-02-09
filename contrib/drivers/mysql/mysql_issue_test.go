// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/guid"
)

// https://github.com/gogf/gf/issues/1934
func Test_Issue1934(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where(" id ", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"], 1)
	})
}

// https://github.com/gogf/gf/issues/1570
func Test_Issue1570(t *testing.T) {
	var (
		tableUser       = "user_" + gtime.TimestampMicroStr()
		tableUserDetail = "user_detail_" + gtime.TimestampMicroStr()
		tableUserScores = "user_scores_" + gtime.TimestampMicroStr()
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  uid int(10) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(45) NOT NULL,
  PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  uid int(10) unsigned NOT NULL AUTO_INCREMENT,
  address varchar(45) NOT NULL,
  PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
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
	gtest.C(t, func(t *gtest.T) {
		var err error
		for i := 1; i <= 5; i++ {
			// User.
			_, err = db.Insert(ctx, tableUser, g.Map{
				"uid":  i,
				"name": fmt.Sprintf(`name_%d`, i),
			})
			t.AssertNil(err)
			// Detail.
			_, err = db.Insert(ctx, tableUserDetail, g.Map{
				"uid":     i,
				"address": fmt.Sprintf(`address_%d`, i),
			})
			t.AssertNil(err)
			// Scores.
			for j := 1; j <= 5; j++ {
				_, err = db.Insert(ctx, tableUserScores, g.Map{
					"uid":   i,
					"score": j,
				})
				t.AssertNil(err)
			}
		}
	})

	// Result ScanList with struct elements and pointer attributes.
	gtest.C(t, func(t *gtest.T) {
		var users []Entity
		// User
		err := db.Model(tableUser).
			Where("uid", g.Slice{3, 4}).
			Fields("uid").
			Order("uid asc").
			ScanList(&users, "User")
		t.AssertNil(err)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, ""})
		t.Assert(users[1].User, &EntityUser{4, ""})
		// Detail
		err = db.Model(tableUserDetail).
			Where("uid", gdb.ListItemValues(users, "User", "Uid")).
			Order("uid asc").
			ScanList(&users, "UserDetail", "User", "uid:Uid")
		t.AssertNil(err)
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		err = db.Model(tableUserScores).
			Where("uid", gdb.ListItemValues(users, "User", "Uid")).
			Order("id asc").
			ScanList(&users, "UserScores", "User", "uid:Uid")
		t.AssertNil(err)
		t.AssertNil(err)
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

// https://github.com/gogf/gf/issues/1401
func Test_Issue1401(t *testing.T) {
	var (
		table1 = "parcels"
		table2 = "parcel_items"
	)
	array := gstr.SplitAndTrim(gtest.DataContent(`issue1401.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table1)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		type NItem struct {
			Id       int `json:"id"`
			ParcelId int `json:"parcel_id"`
		}

		type ParcelItem struct {
			gmeta.Meta `orm:"table:parcel_items"`
			NItem
		}

		type ParcelRsp struct {
			gmeta.Meta `orm:"table:parcels"`
			Id         int           `json:"id"`
			Items      []*ParcelItem `json:"items" orm:"with:parcel_id=Id"`
		}

		parcelDetail := &ParcelRsp{}
		err := db.Model(table1).With(parcelDetail.Items).Where("id", 3).Scan(&parcelDetail)
		t.AssertNil(err)
		t.Assert(parcelDetail.Id, 3)
		t.Assert(len(parcelDetail.Items), 1)
		t.Assert(parcelDetail.Items[0].Id, 2)
		t.Assert(parcelDetail.Items[0].ParcelId, 3)
	})
}

// https://github.com/gogf/gf/issues/1412
func Test_Issue1412(t *testing.T) {
	var (
		table1 = "parcels"
		table2 = "items"
	)
	array := gstr.SplitAndTrim(gtest.DataContent(`issue1412.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table1)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		type Items struct {
			gmeta.Meta `orm:"table:items"`
			Id         int    `json:"id"`
			Name       string `json:"name"`
		}

		type ParcelRsp struct {
			gmeta.Meta `orm:"table:parcels"`
			Id         int   `json:"id"`
			ItemId     int   `json:"item_id"`
			Items      Items `json:"items" orm:"with:Id=ItemId"`
		}

		entity := &ParcelRsp{}
		err := db.Model("parcels").With(Items{}).Where("id=3").Scan(&entity)
		t.AssertNil(err)
		t.Assert(entity.Id, 3)
		t.Assert(entity.ItemId, 0)
		t.Assert(entity.Items.Id, 0)
		t.Assert(entity.Items.Name, "")
	})

	gtest.C(t, func(t *gtest.T) {
		type Items struct {
			gmeta.Meta `orm:"table:items"`
			Id         int    `json:"id"`
			Name       string `json:"name"`
		}

		type ParcelRsp struct {
			gmeta.Meta `orm:"table:parcels"`
			Id         int   `json:"id"`
			ItemId     int   `json:"item_id"`
			Items      Items `json:"items" orm:"with:Id=ItemId"`
		}

		entity := &ParcelRsp{}
		err := db.Model("parcels").With(Items{}).Where("id=30000").Scan(&entity)
		t.AssertNE(err, nil)
		t.Assert(entity.Id, 0)
		t.Assert(entity.ItemId, 0)
		t.Assert(entity.Items.Id, 0)
		t.Assert(entity.Items.Name, "")
	})
}

// https://github.com/gogf/gf/issues/1002
func Test_Issue1002(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	result, err := db.Model(table).Data(g.Map{
		"id":          1,
		"passport":    "port_1",
		"password":    "pass_1",
		"nickname":    "name_2",
		"create_time": "2020-10-27 19:03:33",
	}).Insert()
	gtest.AssertNil(err)
	n, _ := result.RowsAffected()
	gtest.Assert(n, 1)

	// where + string.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Fields("id").Where("create_time>'2020-10-27 19:03:32' and create_time<'2020-10-27 19:03:34'").Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Fields("id").Where("create_time>'2020-10-27 19:03:32' and create_time<'2020-10-27 19:03:34'").Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
	// where + string arguments.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Fields("id").Where("create_time>? and create_time<?", "2020-10-27 19:03:32", "2020-10-27 19:03:34").Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
	// where + gtime.Time arguments.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Fields("id").Where("create_time>? and create_time<?", gtime.New("2020-10-27 19:03:32"), gtime.New("2020-10-27 19:03:34")).Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
	// where + time.Time arguments, UTC.
	gtest.C(t, func(t *gtest.T) {
		t1, _ := time.Parse("2006-01-02 15:04:05", "2020-10-27 11:03:32")
		t2, _ := time.Parse("2006-01-02 15:04:05", "2020-10-27 11:03:34")
		{
			v, err := db.Model(table).Fields("id").Where("create_time>? and create_time<?", t1, t2).Value()
			t.AssertNil(err)
			t.Assert(v.Int(), 1)
		}
	})
	// where + time.Time arguments, +8.
	// gtest.C(t, func(t *gtest.T) {
	//	// Change current timezone to +8 zone.
	//	location, err := time.LoadLocation("Asia/Shanghai")
	//	t.AssertNil(err)
	//	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-10-27 19:03:32", location)
	//	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-10-27 19:03:34", location)
	//	{
	//		v, err := db.Model(table).Fields("id").Where("create_time>? and create_time<?", t1, t2).Value()
	//		t.AssertNil(err)
	//		t.Assert(v.Int(), 1)
	//	}
	//	{
	//		v, err := db.Model(table).Fields("id").Where("create_time>? and create_time<?", t1, t2).FindValue()
	//		t.AssertNil(err)
	//		t.Assert(v.Int(), 1)
	//	}
	//	{
	//		v, err := db.Model(table).Where("create_time>? and create_time<?", t1, t2).FindValue("id")
	//		t.AssertNil(err)
	//		t.Assert(v.Int(), 1)
	//	}
	// })
}

// https://github.com/gogf/gf/issues/1700
func Test_Issue1700(t *testing.T) {
	table := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        id         int(10) unsigned NOT NULL AUTO_INCREMENT,
	        user_id    int(10) unsigned NOT NULL,
	        UserId    int(10) unsigned NOT NULL,
	        PRIMARY KEY (id)
	    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	    `, table,
	)); err != nil {
		gtest.AssertNil(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id     int `orm:"id"`
			Userid int `orm:"user_id"`
			UserId int `orm:"UserId"`
		}
		_, err := db.Model(table).Data(User{
			Id:     1,
			Userid: 2,
			UserId: 3,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).One()
		t.AssertNil(err)
		t.Assert(one, g.Map{
			"id":      1,
			"user_id": 2,
			"UserId":  3,
		})

		for i := 0; i < 1000; i++ {
			var user *User
			err = db.Model(table).Scan(&user)
			t.AssertNil(err)
			t.Assert(user.Id, 1)
			t.Assert(user.Userid, 2)
			t.Assert(user.UserId, 3)
		}
	})
}

// https://github.com/gogf/gf/issues/1701
func Test_Issue1701(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Fields(gdb.Raw("if(id=1,100,null)")).WherePri(1).Value()
		t.AssertNil(err)
		t.Assert(value.String(), 100)
	})
}

// https://github.com/gogf/gf/issues/1733
func Test_Issue1733(t *testing.T) {
	table := "user_" + guid.S()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        id int(8) unsigned zerofill NOT NULL AUTO_INCREMENT,
	        PRIMARY KEY (id)
	    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	    `, table,
	)); err != nil {
		gtest.AssertNil(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for i := 1; i <= 10; i++ {
			_, err := db.Model(table).Data(g.Map{
				"id": i,
			}).Insert()
			t.AssertNil(err)
		}

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 10)
		for i := 0; i < 10; i++ {
			t.Assert(all[i]["id"].Int(), i+1)
		}
	})
}

// https://github.com/gogf/gf/issues/2105
func Test_Issue2105(t *testing.T) {
	table := "issue2105"
	array := gstr.SplitAndTrim(gtest.DataContent(`issue2105.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table)

	type JsonItem struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}
	type Test struct {
		Id   string      `json:"id,omitempty"`
		Json []*JsonItem `json:"json,omitempty"`
	}

	gtest.C(t, func(t *gtest.T) {
		var list []*Test
		err := db.Model(table).Scan(&list)
		t.AssertNil(err)
		t.Assert(len(list), 2)
		t.Assert(len(list[0].Json), 0)
		t.Assert(len(list[1].Json), 3)
	})
}

// https://github.com/gogf/gf/issues/2231
func Test_Issue2231(t *testing.T) {
	var (
		pattern = `(\w+):([\w\-]*):(.*?)@(\w+?)\((.+?)\)/{0,1}([^\?]*)\?{0,1}(.*)`
		link    = `mysql:root:12345678@tcp(127.0.0.1:3306)/a正bc式?loc=Local&parseTime=true`
	)
	gtest.C(t, func(t *gtest.T) {
		match, err := gregex.MatchString(pattern, link)
		t.AssertNil(err)
		t.Assert(match[1], "mysql")
		t.Assert(match[2], "root")
		t.Assert(match[3], "12345678")
		t.Assert(match[4], "tcp")
		t.Assert(match[5], "127.0.0.1:3306")
		t.Assert(match[6], "a正bc式")
		t.Assert(match[7], "loc=Local&parseTime=true")
	})
}

// https://github.com/gogf/gf/issues/2339
func Test_Issue2339(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		model1 := db.Model(table, "u1").Where("id between ? and ?", 1, 9)
		model2 := db.Model("? as u2", model1)
		model3 := db.Model("? as u3", model2)
		all2, err := model2.WhereGT("id", 6).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all2), 3)
		t.Assert(all2[0]["id"], 7)

		all3, err := model3.WhereGT("id", 7).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all3), 2)
		t.Assert(all3[0]["id"], 8)
	})
}

// https://github.com/gogf/gf/issues/2356
func Test_Issue2356(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := "demo_" + guid.S()
		if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        id BIGINT(20) UNSIGNED NOT NULL DEFAULT '0',
	        PRIMARY KEY (id)
	    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	    `, table,
		)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table)

		if _, err := db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (id) VALUES (18446744073709551615);`, table)); err != nil {
			t.AssertNil(err)
		}

		one, err := db.Model(table).One()
		t.AssertNil(err)
		t.AssertEQ(one["id"].Val(), uint64(18446744073709551615))
	})
}
