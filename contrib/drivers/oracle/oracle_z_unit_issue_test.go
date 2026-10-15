// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/guid"
)

// Test_Issue4842 tests that the INTERVAL types, reported by the driver as IntervalDS_DTY
// and IntervalYM_DTY, are not detected as integers and read back as 0.
// See https://github.com/gogf/gf/issues/4842
func Test_Issue4842(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, `SELECT TO_DSINTERVAL('2 12:23:34.456') AS v FROM dual`)
		t.AssertNil(err)
		t.Assert(one["V"].String(), `+02 12:23:34.456000`)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, `SELECT TO_YMINTERVAL('1-10') AS v FROM dual`)
		t.AssertNil(err)
		t.Assert(one["V"].String(), `+01-10`)
	})
	gtest.C(t, func(t *gtest.T) {
		// Genuine number types must keep being detected as numbers.
		one, err := db.GetOne(ctx, `SELECT CAST(42 AS NUMBER(10)) AS v FROM dual`)
		t.AssertNil(err)
		t.Assert(one["V"].Int(), 42)
	})
}

// https://github.com/gogf/gf/issues/1380
func Test_Issue1380(t *testing.T) {
	type GiftImage struct {
		Uid    string `json:"uid"`
		Url    string `json:"url"`
		Status string `json:"status"`
		Name   string `json:"name"`
	}

	type GiftComment struct {
		Name     string `json:"name"`
		Field    string `json:"field"`
		Required bool   `json:"required"`
	}

	type Prop struct {
		Name   string   `json:"name"`
		Values []string `json:"values"`
	}

	type Sku struct {
		GiftId      int64  `json:"gift_id"`
		Name        string `json:"name"`
		ScorePrice  int    `json:"score_price"`
		MarketPrice int    `json:"market_price"`
		CostPrice   int    `json:"cost_price"`
		Stock       int    `json:"stock"`
	}

	type Covers struct {
		List []GiftImage `json:"list"`
	}

	type GiftEntity struct {
		Id                   int64         `json:"id"`
		StoreId              int64         `json:"store_id"`
		GiftType             int           `json:"gift_type"`
		GiftName             string        `json:"gift_name"`
		Description          string        `json:"description"`
		Covers               Covers        `json:"covers"`
		Cover                string        `json:"cover"`
		GiftCategoryId       []int64       `json:"gift_category_id"`
		HasProps             bool          `json:"has_props"`
		OutSn                string        `json:"out_sn"`
		IsLimitSell          bool          `json:"is_limit_sell"`
		LimitSellType        int           `json:"limit_sell_type"`
		LimitSellCycle       string        `json:"limit_sell_cycle"`
		LimitSellCycleCount  int           `json:"limit_sell_cycle_count"`
		LimitSellCustom      bool          `json:"limit_sell_custom"`
		LimitCustomerTags    []int64       `json:"limit_customer_tags"`
		ScorePrice           int           `json:"score_price"`
		MarketPrice          float64       `json:"market_price"`
		CostPrice            int           `json:"cost_price"`
		Stock                int           `json:"stock"`
		Props                []Prop        `json:"props"`
		Skus                 []Sku         `json:"skus"`
		ExpressType          []string      `json:"express_type"`
		Comments             []GiftComment `json:"comments"`
		Content              string        `json:"content"`
		AtLeastRechargeCount int           `json:"at_least_recharge_count"`
		Status               int           `json:"status"`
	}

	table := "jfy_gift"
	dropTable(table)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `1380.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table)
	if _, err := db.Model(table).Data(g.Map{
		"id":                      17,
		"gift_name":               "GIFT",
		"at_least_recharge_count": 1,
		"comments":                `[{"name": "身份证", "field": "idcard", "required": false}, {"name": "留言2", "field": "text", "required": false}]`,
		"content":                 `<p>礼品详情</p>`,
		"cost_price":              0.00,
		"cover":                   "",
		"covers":                  `{"list": [{"uid": "vc-upload-1629292486099-3", "url": "https://cdn.taobao.com/sULsYiwaOPjsKGoBXwKtuewPzACpBDfQ.jpg", "name": "O1CN01OH6PIP1Oc5ot06U17_!!922361725.jpg", "status": "done"}, {"uid": "vc-upload-1629292486099-4", "url": "https://cdn.taobao.com/lqLHDcrFTgNvlWyXfLYZwmsrODzIBtFH.jpg", "name": "O1CN018hBckI1Oc5ouc8ppl_!!922361725.jpg", "status": "done"}, {"uid": "vc-upload-1629292486099-5", "url": "https://cdn.taobao.com/pvqyutXckICmHhbPBQtrVLHuMlXuGxUg.jpg", "name": "O1CN0185Ubp91Oc5osQTTcc_!!922361725.jpg", "status": "done"}]}`,
		"description":             "支持个性定制的父亲节老师长辈的专属礼物",
		"express_type":            `["快递包邮", "同城配送"]`,
		"gift_type":               1,
		"has_props":               0,
		"is_limit_sell":           0,
		"limit_customer_tags":     `[]`,
		"limit_sell_custom":       0,
		"limit_sell_cycle":        "day",
		"limit_sell_cycle_count":  0,
		"limit_sell_type":         1,
		"market_price":            0.00,
		"out_sn":                  "259402",
		"props":                   `[{"name": "颜色", "values": ["红色", "蓝色"]}]`,
		"skus":                    `[{"name": "red", "stock": 10, "gift_id": 1, "cost_price": 80, "score_price": 188, "market_price": 388}, {"name": "blue", "stock": 100, "gift_id": 2, "cost_price": 81, "score_price": 200, "market_price": 288}]`,
		"score_price":             10.00,
		"stock":                   0,
		"create_at":               gtime.New("2021-08-18 21:26:13"),
		"store_id":                100004,
		"status":                  99,
		"view_count":              0,
		"sell_count":              0,
	}).Insert(); err != nil {
		gtest.Error(err)
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			entity = new(GiftEntity)
			err    = db.Model(table).Where("id", 17).Scan(entity)
		)
		t.AssertNil(err)
		t.Assert(len(entity.Skus), 2)

		t.Assert(entity.Skus[0].Name, "red")
		t.Assert(entity.Skus[0].Stock, 10)
		t.Assert(entity.Skus[0].GiftId, 1)
		t.Assert(entity.Skus[0].CostPrice, 80)
		t.Assert(entity.Skus[0].ScorePrice, 188)
		t.Assert(entity.Skus[0].MarketPrice, 388)

		t.Assert(entity.Skus[1].Name, "blue")
		t.Assert(entity.Skus[1].Stock, 100)
		t.Assert(entity.Skus[1].GiftId, 2)
		t.Assert(entity.Skus[1].CostPrice, 81)
		t.Assert(entity.Skus[1].ScorePrice, 200)
		t.Assert(entity.Skus[1].MarketPrice, 288)

		t.Assert(entity.Id, 17)
		t.Assert(entity.StoreId, 100004)
		t.Assert(entity.GiftType, 1)
		t.Assert(entity.GiftName, "GIFT")
		t.Assert(entity.Description, "支持个性定制的父亲节老师长辈的专属礼物")
		t.Assert(len(entity.Covers.List), 3)
		t.Assert(entity.OutSn, "259402")
		t.Assert(entity.LimitCustomerTags, "[]")
		t.Assert(entity.ScorePrice, 10)
		t.Assert(len(entity.Props), 1)
		t.Assert(len(entity.Comments), 2)
		t.Assert(entity.Status, 99)
		t.Assert(entity.Content, `<p>礼品详情</p>`)
	})
}

// https://github.com/gogf/gf/issues/1934
func Test_Issue1934(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where(" id ", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)
	})
}

// https://github.com/gogf/gf/issues/1570
func Test_Issue1570(t *testing.T) {
	var (
		tableUser       = fmt.Sprintf("user_%d", gtime.TimestampMicro()%1e9)
		tableUserDetail = fmt.Sprintf("user_detail_%d", gtime.TimestampMicro()%1e9)
		tableUserScores = fmt.Sprintf("user_scores_%d", gtime.TimestampMicro()%1e9)
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  user_id NUMBER(10) NOT NULL,
  name    VARCHAR2(45) NOT NULL,
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
  id      NUMBER(10) NOT NULL,
  user_id NUMBER(10) NOT NULL,
  score   NUMBER(10) NOT NULL,
  PRIMARY KEY (id)
)
    `, tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)
	createAutoIncrement(tableUserScores, "ID", 1)

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
		err := db.Model(tableUser).
			Where("user_id", g.Slice{3, 4}).
			Fields("user_id").
			Order("user_id asc").
			ScanList(&users, "User")
		t.AssertNil(err)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, ""})
		t.Assert(users[1].User, &EntityUser{4, ""})
		// Detail
		err = db.Model(tableUserDetail).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("user_id asc").
			ScanList(&users, "UserDetail", "User", "user_id:UserId")
		t.AssertNil(err)
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		err = db.Model(tableUserScores).
			Where("user_id", gdb.ListItemValues(users, "User", "UserId")).
			Order("id asc").
			ScanList(&users, "UserScores", "User", "user_id:UserId")
		t.AssertNil(err)
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

// https://github.com/gogf/gf/issues/1401
func Test_Issue1401(t *testing.T) {
	var (
		table1 = "parcels"
		table2 = "parcel_items"
	)
	dropTable(table1)
	dropTable(table2)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `1401.sql`), ";")
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
	dropTable(table1)
	dropTable(table2)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `1412.sql`), ";")
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
	table := fmt.Sprintf("t_issue1002_%d", gtime.TimestampMicro()%1e9)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE %s (
        id          NUMBER(10)   NOT NULL,
        passport    VARCHAR2(45) NOT NULL,
        password    VARCHAR2(45) NOT NULL,
        nickname    VARCHAR2(45) NOT NULL,
        create_time TIMESTAMP,
        PRIMARY KEY (id)
    )
    `, table,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	result, err := db.Model(table).Data(g.Map{
		"id":          1,
		"passport":    "port_1",
		"password":    "pass_1",
		"nickname":    "name_2",
		"create_time": gtime.New("2020-10-27 19:03:33"),
	}).Insert()
	gtest.AssertNil(err)
	n, _ := result.RowsAffected()
	gtest.Assert(n, 1)

	// where + string.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Fields("id").Where("create_time>TIMESTAMP '2020-10-27 19:03:32' and create_time<TIMESTAMP '2020-10-27 19:03:34'").Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Fields("id").Where("create_time>TIMESTAMP '2020-10-27 19:03:32' and create_time<TIMESTAMP '2020-10-27 19:03:34'").Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
	// where + string arguments.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Fields("id").Where("create_time>TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS') and create_time<TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS')", "2020-10-27 19:03:32", "2020-10-27 19:03:34").Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
	// where + gtime.Time arguments.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Fields("id").Where("create_time>? and create_time<?", gtime.New("2020-10-27 19:03:32"), gtime.New("2020-10-27 19:03:34")).Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
	// where + time.Time arguments, local.
	gtest.C(t, func(t *gtest.T) {
		t1, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-10-27 19:03:32", time.Local)
		t2, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-10-27 19:03:34", time.Local)
		{
			v, err := db.Model(table).Fields("id").Where("create_time>? and create_time<?", t1, t2).Value()
			t.AssertNil(err)
			t.Assert(v.Int(), 1)
		}
	})
}

// https://github.com/gogf/gf/issues/1700
func Test_Issue1700(t *testing.T) {
	table := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        id         NUMBER(10) NOT NULL,
	        user_id    NUMBER(10) NOT NULL,
	        UserId     NUMBER(10) NOT NULL,
	        PRIMARY KEY (id)
	    )
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
		for i := 1; i <= 20; i++ {
			_, err := db.Model(table).Data(User{
				Id:     i,
				Userid: 2,
				UserId: 3,
			}).Insert()
			t.AssertNil(err)

			one, err := db.Model(table).WherePri(i).One()
			t.AssertNil(err)
			t.Assert(one, g.Map{
				"ID":      i,
				"USER_ID": 2,
				"USERID":  3,
			})
		}

		for i := 0; i < 100; i++ {
			var user *User
			err := db.Model(table).WherePri(1).Scan(&user)
			t.AssertNil(err)
			t.Assert(user.Id, 1)
		}
		// Note: Userid and UserId are not asserted after Scan. Oracle folds the column UserId to
		// USERID, so neither USER_ID nor USERID equals an orm tag and gconv falls back to fuzzy
		// matching, which ignores case and underscores: both columns match both fields, the pair is
		// assigned in map iteration order and the choice is cached in LastFuzzyKey, see
		// https://github.com/gogf/gf/blob/80bf57a149ed164eb9f7ac0f5217bcc1cd59d036/util/gconv/internal/converter/converter_struct.go#L311-L318.
	})
}

// https://github.com/gogf/gf/issues/1701
func Test_Issue1701(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Fields(gdb.Raw("DECODE(id,1,100,NULL)")).WherePri(1).Value()
		t.AssertNil(err)
		t.Assert(value.String(), 100)
	})
}

// https://github.com/gogf/gf/issues/2105
func Test_Issue2105(t *testing.T) {
	table := "issue2105"
	dropTable(table)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `2105.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table)
	if _, err := db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s VALUES (?, ?)`, table),
		"2", `[{"Name": "任务类型", "Value": "高价值"}, {"Name": "优先级", "Value": "高"}, {"Name": "是否亮点功能", "Value": "是"}]`,
	); err != nil {
		gtest.Error(err)
	}

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
		link    = `oracle:system:oracle@tcp(127.0.0.1:1521)/a正bc式?PREFETCH_ROWS=25&TIMEOUT=60`
	)
	gtest.C(t, func(t *gtest.T) {
		match, err := gregex.MatchString(pattern, link)
		t.AssertNil(err)
		t.Assert(match[1], "oracle")
		t.Assert(match[2], "system")
		t.Assert(match[3], "oracle")
		t.Assert(match[4], "tcp")
		t.Assert(match[5], "127.0.0.1:1521")
		t.Assert(match[6], "a正bc式")
		t.Assert(match[7], "PREFETCH_ROWS=25&TIMEOUT=60")
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
		t.Assert(all2[0]["ID"], 7)

		all3, err := model3.WhereGT("id", 7).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all3), 2)
		t.Assert(all3[0]["ID"], 8)
	})
}

// https://github.com/gogf/gf/issues/2356
func Test_Issue2356(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("demo_%d", gtime.TimestampMicro())
		if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        id NUMBER(20) DEFAULT 0 NOT NULL,
	        PRIMARY KEY (id)
	    )
	    `, table,
		)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table)

		if _, err := db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (id) VALUES (18446744073709551615)`, table)); err != nil {
			t.AssertNil(err)
		}

		one, err := db.Model(table).One()
		t.AssertNil(err)
		t.Assert(one["ID"].String(), "18446744073709551615")
		t.AssertEQ(one["ID"].Uint64(), uint64(18446744073709551615))
	})
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("demo_%d", gtime.TimestampMicro())
		if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        id NUMBER(20) DEFAULT 0 NOT NULL,
	        PRIMARY KEY (id)
	    )
	    `, table,
		)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table)

		if _, err := db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (id) VALUES (18446744073709551614)`, table)); err != nil {
			t.AssertNil(err)
		}

		one, err := db.Model(table).One()
		t.AssertNil(err)
		t.Assert(one["ID"].String(), "18446744073709551614")
		t.AssertEQ(one["ID"].Uint64(), uint64(18446744073709551614))
	})
}

// https://github.com/gogf/gf/issues/2427
func Test_Issue2427(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := fmt.Sprintf("demo_%d", gtime.TimestampMicro())
		if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
    id        NUMBER(10)   NOT NULL,
    passport  VARCHAR2(45) NOT NULL,
    password  VARCHAR2(45) NOT NULL,
    nickname  VARCHAR2(45) NOT NULL,
    create_at TIMESTAMP(6) DEFAULT NULL,
    update_at TIMESTAMP(6) DEFAULT NULL,
    PRIMARY KEY (id)
)
	    `, table,
		)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table)

		_, err1 := db.Model(table).Delete()
		t.Assert(err1, `there should be WHERE condition statement for DELETE operation`)

		_, err2 := db.Model(table).Where(g.Map{}).Delete()
		t.Assert(err2, `there should be WHERE condition statement for DELETE operation`)

		_, err3 := db.Model(table).Where(1).Delete()
		t.AssertNE(err3, nil)
		t.AssertIN("ORA-00920", err3.Error())

		_, err4 := db.Model(table).Where("1=1").Delete()
		t.AssertNil(err4)
	})
}

// https://github.com/gogf/gf/issues/2561
func Test_Issue2561(t *testing.T) {
	table := fmt.Sprintf("t_issue2561_%d", gtime.TimestampMicro()%1e9)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE %s (
        id          NUMBER(10) NOT NULL,
        passport    VARCHAR2(45),
        password    VARCHAR2(32),
        nickname    VARCHAR2(45),
        create_time VARCHAR2(45),
        PRIMARY KEY (id)
    )
    `, table,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := g.Slice{
			User{
				Id:       1,
				Passport: "user_1",
			},
			User{
				Id:       2,
				Password: "pass_2",
			},
			User{
				Id:       3,
				Password: "pass_3",
			},
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		n, _ := result.RowsAffected()
		t.Assert(n, 3)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `1`)
		t.Assert(one[`PASSPORT`], `user_1`)
		t.Assert(one[`PASSWORD`], ``)
		t.Assert(one[`NICKNAME`], ``)
		t.Assert(one[`CREATE_TIME`], ``)

		one, err = db.Model(table).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `2`)
		t.Assert(one[`PASSPORT`], ``)
		t.Assert(one[`PASSWORD`], `pass_2`)
		t.Assert(one[`NICKNAME`], ``)
		t.Assert(one[`CREATE_TIME`], ``)

		one, err = db.Model(table).WherePri(3).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `3`)
		t.Assert(one[`PASSPORT`], ``)
		t.Assert(one[`PASSWORD`], `pass_3`)
		t.Assert(one[`NICKNAME`], ``)
		t.Assert(one[`CREATE_TIME`], ``)
	})
}

// https://github.com/gogf/gf/issues/2439
func Test_Issue2439(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			suffix = fmt.Sprintf("_%d", gtime.TimestampMicro()%1e9)
			a      = "a" + suffix
			b      = "b" + suffix
			c      = "c" + suffix
		)
		defer dropTable(a)
		defer dropTable(b)
		defer dropTable(c)
		array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `2439.sql`), ";")
		for _, v := range array {
			v, err := gregex.ReplaceString(`(TABLE|INTO) ([abc]) `, "${1} ${2}"+suffix+" ", v)
			t.AssertNil(err)
			_, err = db.Exec(ctx, v)
			t.AssertNil(err)
		}

		orm := db.Model(a)
		orm = orm.InnerJoin(
			c, fmt.Sprintf("%s.id=%s.id", a, c),
		)
		orm = orm.InnerJoinOnField(b, "id")
		whereFormat := fmt.Sprintf(
			"(%s.%s LIKE ?) ",
			b, "name",
		)
		orm = orm.WhereOrf(
			whereFormat,
			"%a%",
		)
		r, err := orm.All()
		t.AssertNil(err)
		t.Assert(len(r), 1)
		t.Assert(r[0]["ID"], 2)
		t.Assert(r[0]["NAME"], "a")
	})
}

// https://github.com/gogf/gf/issues/2782
func Test_Issue2787(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table)

		condWhere, _ := m.Builder().
			Where("id", "").
			Where(m.Builder().
				Where("nickname", "foo").
				WhereOr("password", "abc123")).
			Where("passport", "pp").
			Build()
		t.Assert(condWhere, `(id=?) AND (((nickname=?) OR (password=?))) AND (passport=?)`)

		condWhere, _ = m.OmitEmpty().Builder().
			Where("id", "").
			Where(m.Builder().
				Where("nickname", "foo").
				WhereOr("password", "abc123")).
			Where("passport", "pp").
			Build()
		t.Assert(condWhere, `((nickname=?) OR (password=?)) AND (passport=?)`)

		condWhere, _ = m.OmitEmpty().Builder().
			Where(m.Builder().
				Where("nickname", "foo").
				WhereOr("password", "abc123")).
			Where("id", "").
			Where("passport", "pp").
			Build()
		t.Assert(condWhere, `((nickname=?) OR (password=?)) AND (passport=?)`)
	})
}

// https://github.com/gogf/gf/issues/2907
func Test_Issue2907(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		var (
			orm = db.Model(table)
			err error
		)

		orm = orm.WherePrefixNotIn(
			table,
			"id",
			[]int{
				1,
				2,
			},
		)
		all, err := orm.OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize-2)
		t.Assert(all[0]["ID"], 3)
	})
}

// https://github.com/gogf/gf/issues/3086
func Test_Issue3086(t *testing.T) {
	table := "issue3086_user"
	dropTable(table)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `3086.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := g.Slice{
			User{
				Id:       1,
				Passport: "user_1",
			},
			User{
				Id:       1,
				Passport: "user_2",
			},
		}
		_, err := db.Model(table).Data(data).Batch(10).Insert()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := g.Slice{
			User{
				Id:       3,
				Passport: "user_1",
			},
			User{
				Id:       4,
				Passport: "user_2",
			},
		}
		result, err := db.Model(table).Data(data).Batch(10).Insert()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 2)
	})
}

// https://github.com/gogf/gf/issues/3204
func Test_Issue3204(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		makeNullable(t, table, "password", "nickname")
	})

	// where
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any `orm:"id,omitempty"`
			Passport   any `orm:"passport,omitempty"`
			Password   any `orm:"password,omitempty"`
			Nickname   any `orm:"nickname,omitempty"`
			CreateTime any `orm:"create_time,omitempty"`
		}
		where := User{
			Id:       2,
			Passport: "",
		}
		all, err := db.Model(table).Where(where).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 2)
	})
	// data
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any `orm:"id,omitempty"`
			Passport   any `orm:"passport,omitempty"`
			Password   any `orm:"password,omitempty"`
			Nickname   any `orm:"nickname,omitempty"`
			CreateTime any `orm:"create_time,omitempty"`
		}
		var (
			err      error
			sqlArray []string
			data     = User{
				Id:       20,
				Passport: "passport_20",
				Password: "",
			}
		)
		sqlArray, err = gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err = db.Ctx(ctx).Model(table).Data(data).InsertAndGetId()
			return err
		})
		t.AssertNil(err)
		t.AssertIN(sqlArray[len(sqlArray)-1], g.SliceStr{
			fmt.Sprintf("INSERT INTO %s(ID,PASSPORT) VALUES(20,'passport_20')", table),
			fmt.Sprintf("INSERT INTO %s(PASSPORT,ID) VALUES('passport_20',20)", table),
		})
		// Note: the id returned by InsertAndGetId is not asserted. When the data carries the primary
		// key the driver adds no RETURNING clause and the id is 0; that block of DoInsert is rewritten
		// by https://github.com/gogf/gf/pull/4880.
	})
	// update data
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any `orm:"id,omitempty"`
			Passport   any `orm:"passport,omitempty"`
			Password   any `orm:"password,omitempty"`
			Nickname   any `orm:"nickname,omitempty"`
			CreateTime any `orm:"create_time,omitempty"`
		}
		var (
			err      error
			sqlArray []string
			data     = User{
				Passport: "passport_1",
				Password: "",
				Nickname: "",
			}
		)
		sqlArray, err = gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err = db.Ctx(ctx).Model(table).Data(data).WherePri(1).Update()
			return err
		})
		t.AssertNil(err)
		t.Assert(
			gstr.Contains(sqlArray[len(sqlArray)-1], "SET PASSPORT='passport_1' WHERE ID=1"),
			true,
		)
	})
}

// https://github.com/gogf/gf/issues/3218
func Test_Issue3218(t *testing.T) {
	table := "issue3218_sys_config"
	dropTable(table)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `3218.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table)
	if _, err := db.Exec(ctx, fmt.Sprintf(
		`INSERT INTO %s VALUES (?, ?, ?, TIMESTAMP '2023-12-19 14:08:25', TIMESTAMP '2023-12-19 14:08:25')`, table),
		49, "site", `{"banned_ip":"22","filings":"2222","fixed_page":"","site_name":"22","version":"22"}`,
	); err != nil {
		gtest.Error(err)
	}
	gtest.C(t, func(t *gtest.T) {
		type SysConfigInfo struct {
			Name  string            `json:"name"`
			Value map[string]string `json:"value"`
		}
		var configData *SysConfigInfo
		err := db.Model(table).Scan(&configData)
		t.AssertNil(err)
		t.Assert(configData, &SysConfigInfo{
			Name: "site",
			Value: map[string]string{
				"fixed_page": "",
				"site_name":  "22",
				"version":    "22",
				"banned_ip":  "22",
				"filings":    "2222",
			},
		})
	})
}

// https://github.com/gogf/gf/issues/2552
func Test_Issue2552_ClearTableFieldsAll(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	showTableKey := `USER_TAB_COLUMNS`
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err := db.Model(table).Ctx(ctx).Insert(g.Map{
				"passport":    guid.S(),
				"password":    guid.S(),
				"nickname":    guid.S(),
				"create_time": gtime.Now().String(),
			})
			return err
		})
		t.AssertNil(err)
		t.Assert(gstr.Contains(gstr.Join(sqlArray, "|"), showTableKey), true)

		ctx = context.Background()
		sqlArray, err = gdb.CatchSQL(ctx, func(ctx context.Context) error {
			one, err := db.Model(table).Ctx(ctx).One()
			t.Assert(len(one), 6)
			return err
		})
		t.AssertNil(err)
		t.Assert(gstr.Contains(gstr.Join(sqlArray, "|"), showTableKey), false)

		_, err = db.Exec(ctx, fmt.Sprintf("alter table %s drop column nickname", table))
		t.AssertNil(err)

		err = db.GetCore().ClearTableFieldsAll(ctx)
		t.AssertNil(err)

		ctx = context.Background()
		sqlArray, err = gdb.CatchSQL(ctx, func(ctx context.Context) error {
			one, err := db.Model(table).Ctx(ctx).One()
			t.Assert(len(one), 5)
			return err
		})
		t.AssertNil(err)
		t.Assert(gstr.Contains(gstr.Join(sqlArray, "|"), showTableKey), true)
	})
}

// https://github.com/gogf/gf/issues/2552
func Test_Issue2552_ClearTableFields(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	showTableKey := `USER_TAB_COLUMNS`
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err := db.Model(table).Ctx(ctx).Insert(g.Map{
				"passport":    guid.S(),
				"password":    guid.S(),
				"nickname":    guid.S(),
				"create_time": gtime.Now().String(),
			})
			return err
		})
		t.AssertNil(err)
		t.Assert(gstr.Contains(gstr.Join(sqlArray, "|"), showTableKey), true)

		ctx = context.Background()
		sqlArray, err = gdb.CatchSQL(ctx, func(ctx context.Context) error {
			one, err := db.Model(table).Ctx(ctx).One()
			t.Assert(len(one), 6)
			return err
		})
		t.AssertNil(err)
		t.Assert(gstr.Contains(gstr.Join(sqlArray, "|"), showTableKey), false)

		_, err = db.Exec(ctx, fmt.Sprintf("alter table %s drop column nickname", table))
		t.AssertNil(err)

		err = db.GetCore().ClearTableFields(ctx, table)
		t.AssertNil(err)

		ctx = context.Background()
		sqlArray, err = gdb.CatchSQL(ctx, func(ctx context.Context) error {
			one, err := db.Model(table).Ctx(ctx).One()
			t.Assert(len(one), 5)
			return err
		})
		t.AssertNil(err)
		t.Assert(gstr.Contains(gstr.Join(sqlArray, "|"), showTableKey), true)
	})
}

// https://github.com/gogf/gf/issues/2643
func Test_Issue2643(t *testing.T) {
	table := "issue2643"
	dropTable(table)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `2643.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			expectKey1 = "SELECT s.name,replace(concat(lpad(s.id, 6, '0'),s.name),',','') code FROM issue2643"
			expectKey2 = "SELECT CASE WHEN dept='物资部' THEN '物资部' ELSE '其他' END dept,sum(s.value) FROM issue2643"
			expectKey3 = "GROUP BY CASE WHEN dept='物资部' THEN '物资部' ELSE '其他' END"
			err1, err2 error
		)
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err1 = db.Ctx(ctx).Model(table).As("s").Fields(
				"s.name",
				"replace(concat(lpad(s.id, 6, '0'),s.name),',','') code",
			).All()
			_, err2 = db.Ctx(ctx).Model(table).As("s").Fields(
				"CASE WHEN dept='物资部' THEN '物资部' ELSE '其他' END dept",
				"sum(s.value)",
			).Group("CASE WHEN dept='物资部' THEN '物资部' ELSE '其他' END").All()
			return nil
		})
		t.AssertNil(err)
		sqlContent := gstr.Join(sqlArray, "\n")
		t.Assert(gstr.Contains(sqlContent, expectKey1), true)
		t.Assert(gstr.Contains(sqlContent, expectKey2), true)
		t.Assert(gstr.Contains(sqlContent, expectKey3), true)
		t.AssertNil(err1)
		t.AssertNil(err2)
	})
}

// https://github.com/gogf/gf/issues/3238
func Test_Issue3238(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 20; i++ {
			var (
				ids       []gdb.Value
				nicknames []gdb.Value
				errs      []error
			)
			_, err := db.Ctx(ctx).Model(table).Hook(gdb.HookHandler{
				Select: func(ctx context.Context, in *gdb.HookSelectInput) (result gdb.Result, err error) {
					result, err = in.Next(ctx)
					if err != nil {
						return
					}
					var wg sync.WaitGroup
					ids = make([]gdb.Value, len(result))
					nicknames = make([]gdb.Value, len(result))
					errs = make([]error, len(result))
					for index := range result {
						wg.Add(1)
						go func(index int) {
							defer wg.Done()
							if ids[index], errs[index] = db.Ctx(ctx).Model(table).WherePri(1).Value(`id`); errs[index] != nil {
								return
							}
							nicknames[index], errs[index] = db.Ctx(ctx).Model(table).WherePri(1).Value(`nickname`)
						}(index)
					}
					wg.Wait()
					return
				},
			},
			).All()
			t.AssertNil(err)
			t.Assert(len(errs), TableSize)
			for index := range errs {
				t.AssertNil(errs[index])
				t.Assert(ids[index].Int(), 1)
				t.Assert(nicknames[index].String(), "name_1")
			}
		}
	})
}

// https://github.com/gogf/gf/issues/3649
func Test_Issue3649(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		sql, err := gdb.CatchSQL(context.Background(), func(ctx context.Context) (err error) {
			user := db.Model(table).Ctx(ctx)
			_, err = user.Where("create_time = ?", gdb.Raw("TO_CHAR(SYSDATE,'YYYY-MM-DD HH24:MI:SS')")).WhereLT("create_time", gdb.Raw("TO_CHAR(SYSDATE,'YYYY-MM-DD HH24:MI:SS')")).Count()
			return
		})
		t.AssertNil(err)
		sqlStr := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE (create_time = TO_CHAR(SYSDATE,'YYYY-MM-DD HH24:MI:SS')) AND (create_time < TO_CHAR(SYSDATE,'YYYY-MM-DD HH24:MI:SS'))", table)
		t.Assert(sql[0], sqlStr)
	})
}

// https://github.com/gogf/gf/issues/3754
func Test_Issue3754(t *testing.T) {
	table := "issue3754"
	dropTable(table)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `3754.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fieldsEx := []string{"delete_at", "create_at", "update_at"}
		// Insert.
		dataInsert := g.Map{
			"id":   1,
			"name": "name_1",
		}
		r, err := db.Model(table).Data(dataInsert).FieldsEx(fieldsEx).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		oneInsert, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInsert["ID"].Int(), 1)
		t.Assert(oneInsert["NAME"].String(), "name_1")
		t.Assert(oneInsert["DELETE_AT"].String(), "")
		t.Assert(oneInsert["CREATE_AT"].String(), "")
		t.Assert(oneInsert["UPDATE_AT"].String(), "")

		// Update.
		dataUpdate := g.Map{
			"name": "name_1000",
		}
		r, err = db.Model(table).Data(dataUpdate).FieldsEx(fieldsEx).WherePri(1).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["ID"].Int(), 1)
		t.Assert(oneUpdate["NAME"].String(), "name_1000")
		t.Assert(oneUpdate["DELETE_AT"].String(), "")
		t.Assert(oneUpdate["CREATE_AT"].String(), "")
		t.Assert(oneUpdate["UPDATE_AT"].String(), "")

		// FieldsEx does not affect Delete operation.
		r, err = db.Model(table).FieldsEx(fieldsEx).WherePri(1).Delete()
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
		oneDeleteUnscoped, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneDeleteUnscoped["ID"].Int(), 1)
		t.Assert(oneDeleteUnscoped["NAME"].String(), "name_1000")
		t.AssertNE(oneDeleteUnscoped["DELETE_AT"].String(), "")
		t.Assert(oneDeleteUnscoped["CREATE_AT"].String(), "")
		t.Assert(oneDeleteUnscoped["UPDATE_AT"].String(), "")
	})
}

// https://github.com/gogf/gf/issues/3626
func Test_Issue3626(t *testing.T) {
	table := "issue3626"
	dropTable(table)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `3626.sql`), ";")
	defer dropTable(table)
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}

	// Insert.
	gtest.C(t, func(t *gtest.T) {
		dataInsert := g.Map{
			"id":   1,
			"name": "name_1",
		}
		r, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		oneInsert, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInsert["ID"].Int(), 1)
		t.Assert(oneInsert["NAME"].String(), "name_1")
	})

	var (
		cacheKey  = guid.S()
		cacheFunc = func(duration time.Duration) gdb.HookHandler {
			return gdb.HookHandler{
				Select: func(ctx context.Context, in *gdb.HookSelectInput) (result gdb.Result, err error) {
					get, err := db.GetCache().Get(ctx, cacheKey)
					if err == nil && !get.IsEmpty() {
						err = get.Scan(&result)
						if err == nil {
							return result, nil
						}
					}
					result, err = in.Next(ctx)
					if err != nil {
						return nil, err
					}
					if result == nil || result.Len() < 1 {
						result = make(gdb.Result, 0)
					}
					_ = db.GetCache().Set(ctx, cacheKey, result, duration)
					return
				},
			}
		}
	)
	gtest.C(t, func(t *gtest.T) {
		defer db.GetCache().Clear(ctx)
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
		count, err = db.Model(table).Hook(cacheFunc(time.Hour)).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
		count, err = db.Model(table).Hook(cacheFunc(time.Hour)).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// https://github.com/gogf/gf/issues/3932
func Test_Issue3932(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Order("id", "desc").One()
		t.AssertNil(err)
		t.Assert(one["ID"], 10)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Order("id desc").One()
		t.AssertNil(err)
		t.Assert(one["ID"], 10)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Order("id desc, nickname asc").One()
		t.AssertNil(err)
		t.Assert(one["ID"], 10)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Order("id desc", "nickname asc").One()
		t.AssertNil(err)
		t.Assert(one["ID"], 10)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Order("id desc").Order("nickname asc").One()
		t.AssertNil(err)
		t.Assert(one["ID"], 10)
	})
}

// https://github.com/gogf/gf/issues/3968
func Test_Issue3968(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var hook = gdb.HookHandler{
			Select: func(ctx context.Context, in *gdb.HookSelectInput) (result gdb.Result, err error) {
				result, err = in.Next(ctx)
				if err != nil {
					return nil, err
				}
				if result != nil {
					for i := range result {
						result[i]["location"] = gvar.New("ny")
					}
				}
				return
			},
		}
		var (
			count  int
			result gdb.Result
		)
		err := db.Model(table).Hook(hook).ScanAndCount(&result, &count, false)
		t.AssertNil(err)
		t.Assert(count, 10)
		t.Assert(len(result), 10)
	})
}

// https://github.com/gogf/gf/issues/3915
func Test_Issue3915(t *testing.T) {
	table := "issue3915"
	dropTable(table)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `3915.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("a < b").All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 1)

		all, err = db.Model(table).Where(gdb.Raw("a < b")).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 1)

		all, err = db.Model(table).WhereLT("a", gdb.Raw("b")).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 1)
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("a > b").All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 2)

		all, err = db.Model(table).Where(gdb.Raw("a > b")).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 2)

		all, err = db.Model(table).WhereGT("a", gdb.Raw("b")).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 2)
	})
}

type RoleBase struct {
	gmeta.Meta  `orm:"table:sys_role"`
	Name        string      `json:"name"           description:"角色名称"     `
	Code        string      `json:"code"           description:"角色 code"    `
	Description string      `json:"description"    description:"描述信息"     `
	Weight      int         `json:"weight"         description:"排序"         `
	StatusId    int         `json:"statusId"       description:"发布状态"     `
	CreatedAt   *gtime.Time `json:"createdAt"      description:""             `
	UpdatedAt   *gtime.Time `json:"updatedAt"      description:""             `
}

type Role struct {
	gmeta.Meta `orm:"table:sys_role"`
	RoleBase
	Id     uint    `json:"id"          description:""`
	Status *Status `json:"status"       description:"发布状态"     orm:"with:id=status_id"        `
}

type StatusBase struct {
	gmeta.Meta `orm:"table:sys_status"`
	En         string `json:"en"        description:"英文名称"    `
	Cn         string `json:"cn"        description:"中文名称"    `
	Weight     int    `json:"weight"    description:"排序权重"    `
}

type Status struct {
	gmeta.Meta `orm:"table:sys_status"`
	StatusBase
	Id uint `json:"id"          description:""`
}

// https://github.com/gogf/gf/issues/2119
func Test_Issue2119(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tables := []string{
			"sys_role",
			"sys_status",
		}

		dropTable(tables[0])
		dropTable(tables[1])
		defer dropTable(tables[0])
		defer dropTable(tables[1])
		_ = tables
		array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `2119.sql`), ";")
		for _, v := range array {
			_, err := db.Exec(ctx, v)
			t.AssertNil(err)
		}
		roles := make([]*Role, 0)
		err := db.Ctx(context.Background()).Model(&Role{}).WithAll().Scan(&roles)
		t.AssertNil(err)
		expectStatus := []*Status{
			{
				StatusBase: StatusBase{
					En:     "undecided",
					Cn:     "未决定",
					Weight: 800,
				},
				Id: 2,
			},
			{
				StatusBase: StatusBase{
					En:     "on line",
					Cn:     "上线",
					Weight: 900,
				},
				Id: 1,
			},
			{
				StatusBase: StatusBase{
					En:     "on line",
					Cn:     "上线",
					Weight: 900,
				},
				Id: 1,
			},
			{
				StatusBase: StatusBase{
					En:     "on line",
					Cn:     "上线",
					Weight: 900,
				},
				Id: 1,
			},
			{
				StatusBase: StatusBase{
					En:     "on line",
					Cn:     "上线",
					Weight: 900,
				},
				Id: 1,
			},
		}

		for i := 0; i < len(roles); i++ {
			t.Assert(roles[i].Status, expectStatus[i])
		}
	})
}

// https://github.com/gogf/gf/issues/4034
func Test_Issue4034(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := "issue4034"
		dropTable(table)
		array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `4034.sql`), ";")
		for _, v := range array {
			_, err := db.Exec(ctx, v)
			t.AssertNil(err)
		}
		defer dropTable(table)

		err := issue4034SaveDeviceAndToken(ctx, table)
		t.AssertNil(err)
	})
}

func issue4034SaveDeviceAndToken(ctx context.Context, table string) error {
	return db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := issue4034SaveAppDevice(ctx, table, tx); err != nil {
			return err
		}
		return nil
	})
}

func issue4034SaveAppDevice(ctx context.Context, table string, tx gdb.TX) error {
	_, err := db.Model(table).Safe().Ctx(ctx).TX(tx).Data(g.Map{
		"id":       1,
		"passport": "111",
		"password": "222",
		"nickname": "333",
	}).Save()
	return err
}

// https://github.com/gogf/gf/issues/4086
func Test_Issue4086(t *testing.T) {
	table := "issue4086"
	dropTable(table)
	defer dropTable(table)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `4086.sql`), ";")
	for _, v := range array {
		_, err := db.Exec(ctx, v)
		gtest.AssertNil(err)
	}

	gtest.C(t, func(t *gtest.T) {
		type ProxyParam struct {
			ProxyId      int64   `json:"proxyId" orm:"proxy_id"`
			RecommendIds []int64 `json:"recommendIds" orm:"recommend_ids"`
			Photos       []int64 `json:"photos" orm:"photos"`
		}

		var proxyParamList []*ProxyParam
		err := db.Model(table).Ctx(ctx).Scan(&proxyParamList)
		t.AssertNil(err)
		t.Assert(len(proxyParamList), 2)
		t.Assert(proxyParamList, []*ProxyParam{
			{
				ProxyId:      1,
				RecommendIds: []int64{584, 585},
				Photos:       nil,
			},
			{
				ProxyId:      2,
				RecommendIds: []int64{},
				Photos:       nil,
			},
		})
	})

	gtest.C(t, func(t *gtest.T) {
		type ProxyParam struct {
			ProxyId      int64     `json:"proxyId" orm:"proxy_id"`
			RecommendIds []int64   `json:"recommendIds" orm:"recommend_ids"`
			Photos       []float32 `json:"photos" orm:"photos"`
		}

		var proxyParamList []*ProxyParam
		err := db.Model(table).Ctx(ctx).Scan(&proxyParamList)
		t.AssertNil(err)
		t.Assert(len(proxyParamList), 2)
		t.Assert(proxyParamList, []*ProxyParam{
			{
				ProxyId:      1,
				RecommendIds: []int64{584, 585},
				Photos:       nil,
			},
			{
				ProxyId:      2,
				RecommendIds: []int64{},
				Photos:       nil,
			},
		})
	})

	gtest.C(t, func(t *gtest.T) {
		type ProxyParam struct {
			ProxyId      int64    `json:"proxyId" orm:"proxy_id"`
			RecommendIds []int64  `json:"recommendIds" orm:"recommend_ids"`
			Photos       []string `json:"photos" orm:"photos"`
		}

		var proxyParamList []*ProxyParam
		err := db.Model(table).Ctx(ctx).Scan(&proxyParamList)
		t.AssertNil(err)
		t.Assert(len(proxyParamList), 2)
		t.Assert(proxyParamList, []*ProxyParam{
			{
				ProxyId:      1,
				RecommendIds: []int64{584, 585},
				Photos:       nil,
			},
			{
				ProxyId:      2,
				RecommendIds: []int64{},
				Photos:       nil,
			},
		})
	})

	gtest.C(t, func(t *gtest.T) {
		type ProxyParam struct {
			ProxyId      int64   `json:"proxyId" orm:"proxy_id"`
			RecommendIds []int64 `json:"recommendIds" orm:"recommend_ids"`
			Photos       []any   `json:"photos" orm:"photos"`
		}

		var proxyParamList []*ProxyParam
		err := db.Model(table).Ctx(ctx).Scan(&proxyParamList)
		t.AssertNil(err)
		t.Assert(len(proxyParamList), 2)
		t.Assert(proxyParamList, []*ProxyParam{
			{
				ProxyId:      1,
				RecommendIds: []int64{584, 585},
				Photos:       nil,
			},
			{
				ProxyId:      2,
				RecommendIds: []int64{},
				Photos:       nil,
			},
		})
	})
	gtest.C(t, func(t *gtest.T) {
		type ProxyParam struct {
			ProxyId      int64   `json:"proxyId" orm:"proxy_id"`
			RecommendIds []int64 `json:"recommendIds" orm:"recommend_ids"`
			Photos       string  `json:"photos" orm:"photos"`
		}

		var proxyParamList []*ProxyParam
		err := db.Model(table).Ctx(ctx).Scan(&proxyParamList)
		t.AssertNil(err)
		t.Assert(len(proxyParamList), 2)
		t.Assert(proxyParamList, []*ProxyParam{
			{
				ProxyId:      1,
				RecommendIds: []int64{584, 585},
				Photos:       "null",
			},
			{
				ProxyId:      2,
				RecommendIds: []int64{},
				Photos:       "",
			},
		})
	})
	gtest.C(t, func(t *gtest.T) {
		type ProxyParam struct {
			ProxyId      int64           `json:"proxyId" orm:"proxy_id"`
			RecommendIds string          `json:"recommendIds" orm:"recommend_ids"`
			Photos       json.RawMessage `json:"photos" orm:"photos"`
		}

		var proxyParamList []*ProxyParam
		err := db.Model(table).Ctx(ctx).Scan(&proxyParamList)
		t.AssertNil(err)
		t.Assert(len(proxyParamList), 2)
		t.Assert(proxyParamList, []*ProxyParam{
			{
				ProxyId:      1,
				RecommendIds: "[584, 585]",
				Photos:       json.RawMessage("null"),
			},
			{
				ProxyId:      2,
				RecommendIds: "[]",
				Photos:       json.RawMessage("null"),
			},
		})
	})
}

// https://github.com/gogf/gf/issues/4500
// Raw() Count ignores Where condition
func Test_Issue4500(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// Test 1: Raw SQL with WHERE + external Where condition + Count
	// This tests that formatCondition correctly uses AND when Raw SQL already has WHERE
	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id IN (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("id", 8).
			Count()
		t.AssertNil(err)
		// Raw SQL: id IN (1,5,7,8,9,10) = 6 records
		// Where: id < 8 filters to {1,5,7} = 3 records
		t.Assert(count, 3)
	})

	// Test 2: Raw SQL without WHERE + external Where condition + Count
	// This tests that formatCondition correctly adds WHERE
	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s", table)).
			WhereLT("id", 5).
			Count()
		t.AssertNil(err)
		// Raw SQL: all 10 records
		// Where: id < 5 = {1,2,3,4} = 4 records
		t.Assert(count, 4)
	})

	// Test 3: Raw + Where + ScanAndCount
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Passport string
		}
		var users []User
		var total int
		err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id IN (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("id", 8).
			ScanAndCount(&users, &total, false)
		t.AssertNil(err)
		// Both scan result and count should respect Where condition
		t.Assert(len(users), 3)
		t.Assert(total, 3)
	})

	// Test 4: Raw + multiple Where conditions + Count
	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id > ?", table), 0).
			WhereLT("id", 5).
			WhereGTE("id", 2).
			Count()
		t.AssertNil(err)
		// Raw: id > 0 (all 10 records)
		// Where: id < 5 AND id >= 2 = {2, 3, 4} = 3 records
		t.Assert(count, 3)
	})

	// Test 5: Raw SQL with no external Where + Count (baseline test)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id IN (?)", table), g.Slice{1, 2, 3}).
			Count()
		t.AssertNil(err)
		// Should count 3 records
		t.Assert(count, 3)
	})

	// Test 6: Verify All() still works correctly with Raw + Where
	gtest.C(t, func(t *gtest.T) {
		all, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id IN (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("id", 8).
			All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
	})
}

// https://github.com/gogf/gf/issues/4697
func Test_Issue4697(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Fields("") should be treated as Fields() and select all fields
		result, err := db.Model(table).Fields("").Limit(1).All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(len(result[0]), 6)
	})

	gtest.C(t, func(t *gtest.T) {
		// Fields("", "id") should ignore empty string and only select "id"
		result, err := db.Model(table).Fields("", "id").Limit(1).All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(len(result[0]), 1)
		t.AssertNE(result[0]["ID"], nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Fields("id", "", "nickname") should ignore empty string
		result, err := db.Model(table).Fields("id", "", "nickname").Limit(1).All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(len(result[0]), 2)
		t.AssertNE(result[0]["ID"], nil)
		t.AssertNE(result[0]["NICKNAME"], nil)
	})
}

// https://github.com/gogf/gf/issues/4698
func Test_Issue4698(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// Test 1: AllAndCount with multiple fields should generate valid COUNT SQL
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id, nickname").AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(len(result), TableSize)
		t.AssertNE(result[0]["ID"], nil)
		t.AssertNE(result[0]["NICKNAME"], nil)
		t.Assert(result[0]["PASSPORT"], nil)
	})

	// Test 2: AllAndCount(false) with multiple fields
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id, nickname").AllAndCount(false)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(len(result), TableSize)
	})

	// Test 3: ScanAndCount with multiple fields
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Nickname string
		}
		var users []User
		var total int
		err := db.Model(table).Fields("id, nickname").ScanAndCount(&users, &total, true)
		t.AssertNil(err)
		t.Assert(total, TableSize)
		t.Assert(len(users), TableSize)
		t.AssertGT(users[0].Id, 0)
		t.AssertNE(users[0].Nickname, "")
	})

	// Test 4: AllAndCount with single field and useFieldForCount=true
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id").AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(len(result), TableSize)
		t.Assert(len(result[0]), 1)
	})

	// Test 5: AllAndCount with Where condition
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id, nickname").Where("id<?", 5).AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, 4)
		t.Assert(len(result), 4)
	})

	// Test 6: Distinct + AllAndCount(false) should use COUNT(1), not COUNT(DISTINCT 1)
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("nickname").Distinct().AllAndCount(false)
		t.AssertNil(err)
		// COUNT(1) should return total rows, not distinct count
		t.Assert(count, TableSize)
		t.Assert(len(result), TableSize)
	})

	// Test 7: Distinct + AllAndCount(true) with single field should use COUNT(DISTINCT nickname)
	gtest.C(t, func(t *gtest.T) {
		_, count, err := db.Model(table).Fields("nickname").Distinct().AllAndCount(true)
		t.AssertNil(err)
		// COUNT(DISTINCT nickname) should return distinct count
		t.Assert(count, TableSize)
	})

	// Test 8: Distinct + multiple fields + AllAndCount(true) should fallback to COUNT(1)
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id, nickname").Distinct().AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(len(result), TableSize)
	})
}
