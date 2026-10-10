// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

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

const chCCreateTime = "2018-10-24 10:00:00"

// chCCreateNullableTable creates a table with the columns of the default table of the MySQL tests,
// which are all nullable and include create_date.
func chCCreateNullableTable() string {
	name := fmt.Sprintf(`user_%d`, gtime.TimestampNano())
	dropTable(name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id          UInt32,
			passport    Nullable(String),
			password    Nullable(String),
			nickname    Nullable(String),
			create_time Nullable(DateTime64(6)),
			create_date Nullable(Date)
		) ENGINE = MergeTree()
		ORDER BY id
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

// chCCreateNullableInitTable creates a table with chCCreateNullableTable and inserts TableSize
// records like the default table of the MySQL tests.
func chCCreateNullableInitTable() string {
	name := chCCreateNullableTable()
	list := make(g.List, 0, TableSize)
	for i := 1; i <= TableSize; i++ {
		list = append(list, g.Map{
			"id":          i,
			"passport":    fmt.Sprintf(`user_%d`, i),
			"password":    fmt.Sprintf(`pass_%d`, i),
			"nickname":    fmt.Sprintf(`name_%d`, i),
			"create_time": gtime.NewFromStr(chCCreateTime).String(),
		})
	}
	if _, err := db.Insert(ctx, name, list); err != nil {
		gtest.Fatal(err)
	}
	return name
}

// chCNewSchemaDb creates a database with a unique name and returns a database object using it, so
// that tables whose names are fixed in `orm` tags do not collide with other tests.
func chCNewSchemaDb() (gdb.DB, string) {
	schema := fmt.Sprintf(`chc_schema_%d`, gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", schema)); err != nil {
		gtest.Fatal(err)
	}
	schemaDb, err := gdb.New(gdb.ConfigNode{
		Host:  "127.0.0.1",
		Port:  "9000",
		User:  "default",
		Name:  schema,
		Type:  "clickhouse",
		Extra: "mutations_sync=1",
	})
	if err != nil {
		gtest.Fatal(err)
	}
	return schemaDb, schema
}

// chCDropSchema drops the database created by chCNewSchemaDb.
func chCDropSchema(schema string) {
	if _, err := db.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", schema)); err != nil {
		gtest.Error(err)
	}
}

// chCExecSqls executes the statements in order and fails on the first error.
func chCExecSqls(execDb gdb.DB, sqls ...string) {
	for _, sql := range sqls {
		if _, err := execDb.Exec(ctx, sql); err != nil {
			gtest.Fatal(err)
		}
	}
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

	table := "jfy_gift_" + gtime.TimestampNanoStr()
	chCExecSqls(db, fmt.Sprintf(`
		CREATE TABLE %s (
			id                      UInt32,
			gift_name               String,
			at_least_recharge_count UInt32 DEFAULT 1,
			comments                String,
			content                 String,
			cost_price              Nullable(Decimal(10, 2)),
			cover                   String,
			covers                  String,
			description             String DEFAULT '',
			express_type            String,
			gift_type               Int32,
			has_props               UInt8 DEFAULT 0,
			is_limit_sell           UInt8 DEFAULT 0,
			limit_customer_tags     String,
			limit_sell_custom       UInt8 DEFAULT 0,
			limit_sell_cycle        String DEFAULT '',
			limit_sell_cycle_count  Int32,
			limit_sell_type         Int8,
			market_price            Decimal(10, 2),
			out_sn                  String,
			props                   String,
			skus                    String,
			score_price             Decimal(10, 2),
			stock                   Int32,
			create_at               DateTime,
			store_id                Int32,
			status                  Nullable(UInt32) DEFAULT 1,
			view_count              Int32 DEFAULT 0,
			sell_count              Nullable(Int32) DEFAULT 0
		) ENGINE = MergeTree()
		ORDER BY id`, table,
	), fmt.Sprintf(`
		INSERT INTO %s VALUES (17, 'GIFT', 1, '[{"name": "身份证", "field": "idcard", "required": false}, {"name": "留言2", "field": "text", "required": false}]', '<p>礼品详情</p>', 0.00, '', '{"list": [{"uid": "vc-upload-1629292486099-3", "url": "https://cdn.taobao.com/sULsYiwaOPjsKGoBXwKtuewPzACpBDfQ.jpg", "name": "O1CN01OH6PIP1Oc5ot06U17_!!922361725.jpg", "status": "done"}, {"uid": "vc-upload-1629292486099-4", "url": "https://cdn.taobao.com/lqLHDcrFTgNvlWyXfLYZwmsrODzIBtFH.jpg", "name": "O1CN018hBckI1Oc5ouc8ppl_!!922361725.jpg", "status": "done"}, {"uid": "vc-upload-1629292486099-5", "url": "https://cdn.taobao.com/pvqyutXckICmHhbPBQtrVLHuMlXuGxUg.jpg", "name": "O1CN0185Ubp91Oc5osQTTcc_!!922361725.jpg", "status": "done"}]}', '支持个性定制的父亲节老师长辈的专属礼物', '["快递包邮", "同城配送"]', 1, 0, 0, '[]', 0, 'day', 0, 1, 0.00, '259402', '[{"name": "颜色", "values": ["红色", "蓝色"]}]', '[{"name": "red", "stock": 10, "gift_id": 1, "cost_price": 80, "score_price": 188, "market_price": 388}, {"name": "blue", "stock": 100, "gift_id": 2, "cost_price": 81, "score_price": 200, "market_price": 288}]', 10.00, 0, '2021-08-18 21:26:13', 100004, 99, 0, 0)`,
		table,
	))
	defer dropTable(table)

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
		t.Assert(one["id"], 1)
	})
}

// https://github.com/gogf/gf/issues/1570
func Test_Issue1570(t *testing.T) {
	var (
		tableUser       = "user_" + gtime.TimestampNanoStr()
		tableUserDetail = "user_detail_" + gtime.TimestampNanoStr()
		tableUserScores = "user_scores_" + gtime.TimestampNanoStr()
	)
	chCExecSqls(db,
		fmt.Sprintf(`CREATE TABLE %s (uid UInt32, name String) ENGINE = MergeTree() ORDER BY uid`, tableUser),
		fmt.Sprintf(`CREATE TABLE %s (uid UInt32, address String) ENGINE = MergeTree() ORDER BY uid`, tableUserDetail),
		fmt.Sprintf(`CREATE TABLE %s (id UInt32, uid UInt32, score UInt32) ENGINE = MergeTree() ORDER BY id`, tableUserScores),
	)
	defer dropTable(tableUser)
	defer dropTable(tableUserDetail)
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
		var (
			err     error
			scoreId = 0
		)
		for i := 1; i <= 5; i++ {
			_, err = db.Insert(ctx, tableUser, g.Map{
				"uid":  i,
				"name": fmt.Sprintf(`name_%d`, i),
			})
			t.AssertNil(err)
			_, err = db.Insert(ctx, tableUserDetail, g.Map{
				"uid":     i,
				"address": fmt.Sprintf(`address_%d`, i),
			})
			t.AssertNil(err)
			for j := 1; j <= 5; j++ {
				scoreId++
				_, err = db.Insert(ctx, tableUserScores, g.Map{
					"id":    scoreId,
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
		t.Assert(len(users), 2)
		t.Assert(users[0].User, &EntityUser{3, ""})
		t.Assert(users[1].User, &EntityUser{4, ""})
		// Detail
		err = db.Model(tableUserDetail).
			Where("uid", gdb.ListItemValues(users, "User", "Uid")).
			Order("uid asc").
			ScanList(&users, "UserDetail", "User", "uid:Uid")
		t.AssertNil(err)
		t.Assert(users[0].UserDetail, &EntityUserDetail{3, "address_3"})
		t.Assert(users[1].UserDetail, &EntityUserDetail{4, "address_4"})
		// Scores
		err = db.Model(tableUserScores).
			Where("uid", gdb.ListItemValues(users, "User", "Uid")).
			Order("id asc").
			ScanList(&users, "UserScores", "User", "uid:Uid")
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
	schemaDb, schema := chCNewSchemaDb()
	defer chCDropSchema(schema)
	chCExecSqls(schemaDb,
		`CREATE TABLE parcel_items (id Int32, parcel_id Nullable(Int32), name Nullable(String)) ENGINE = MergeTree() ORDER BY id`,
		`INSERT INTO parcel_items VALUES (1, 1, '新品'), (2, 3, '新品2')`,
		`CREATE TABLE parcels (id Int32) ENGINE = MergeTree() ORDER BY id`,
		`INSERT INTO parcels VALUES (1), (2), (3)`,
	)

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
		err := schemaDb.Model("parcels").With(parcelDetail.Items).Where("id", 3).Scan(&parcelDetail)
		t.AssertNil(err)
		t.Assert(parcelDetail.Id, 3)
		t.Assert(len(parcelDetail.Items), 1)
		t.Assert(parcelDetail.Items[0].Id, 2)
		t.Assert(parcelDetail.Items[0].ParcelId, 3)
	})
}

// https://github.com/gogf/gf/issues/1412
// Note: the column names of ClickHouse are case-sensitive, so the with tag names the column `id`
// where the MySQL test writes `Id`.
func Test_Issue1412(t *testing.T) {
	schemaDb, schema := chCNewSchemaDb()
	defer chCDropSchema(schema)
	chCExecSqls(schemaDb,
		`CREATE TABLE items (id Int32, name Nullable(String)) ENGINE = MergeTree() ORDER BY id`,
		`INSERT INTO items VALUES (1, '金秋产品1'), (2, '金秋产品2')`,
		`CREATE TABLE parcels (id Int32, item_id Nullable(Int32)) ENGINE = MergeTree() ORDER BY id`,
		`INSERT INTO parcels VALUES (1, 1), (2, 2), (3, 0)`,
	)

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
			Items      Items `json:"items" orm:"with:id=ItemId"`
		}

		entity := &ParcelRsp{}
		err := schemaDb.Model("parcels").With(Items{}).Where("id=3").Scan(&entity)
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
			Items      Items `json:"items" orm:"with:id=ItemId"`
		}

		entity := &ParcelRsp{}
		err := schemaDb.Model("parcels").With(Items{}).Where("id=30000").Scan(&entity)
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

	_, err := db.Model(table).Data(g.Map{
		"id":          1,
		"passport":    "port_1",
		"password":    "pass_1",
		"nickname":    "name_2",
		"create_time": "2020-10-27 19:03:33",
	}).Insert()
	gtest.AssertNil(err)

	// where + string.
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
		v, err := db.Model(table).Fields("id").Where("create_time>? and create_time<?", t1, t2).Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
}

// https://github.com/gogf/gf/issues/1700
func Test_Issue1700(t *testing.T) {
	table := "user_" + gtime.Now().TimestampNanoStr()
	chCExecSqls(db, fmt.Sprintf(`
		CREATE TABLE %s (
			id      UInt32,
			user_id UInt32,
			UserId  UInt32
		) ENGINE = MergeTree()
		ORDER BY id`, table,
	))
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
// Note: ClickHouse has no ZEROFILL, so a plain UInt32 column holds the ids.
func Test_Issue1733(t *testing.T) {
	table := "user_" + guid.S()
	chCExecSqls(db, fmt.Sprintf(`CREATE TABLE %s (id UInt32) ENGINE = MergeTree() ORDER BY id`, table))
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

// https://github.com/gogf/gf/issues/2012
// Note: ClickHouse 24.11 has no TIME type, so a DateTime column holds the value. A time of day
// parsed by gtime carries the date 0000-01-01, outside of the range of DateTime, so the inserts
// are checked but not the stored time.
func Test_Issue2012(t *testing.T) {
	table := "time_only_" + guid.S()
	chCExecSqls(db, fmt.Sprintf(`
		CREATE TABLE %s (
			id        UInt32,
			time_only Nullable(DateTime)
		) ENGINE = MergeTree()
		ORDER BY id`, table,
	))
	defer dropTable(table)

	type TimeOnly struct {
		Id       int         `json:"id"`
		TimeOnly *gtime.Time `json:"timeOnly"`
	}

	gtest.C(t, func(t *gtest.T) {
		timeOnly := gtime.New("15:04:05")
		m := db.Model(table)

		_, err := m.Insert(TimeOnly{
			TimeOnly: gtime.New(timeOnly),
		})
		t.AssertNil(err)

		_, err = m.Insert(g.Map{
			"time_only": timeOnly,
		})
		t.AssertNil(err)

		_, err = m.Insert("time_only", timeOnly)
		t.AssertNil(err)

		count, err := db.Model(table).WhereNotNull("time_only").Count()
		t.AssertNil(err)
		t.Assert(count, 3)
	})
}

// https://github.com/gogf/gf/issues/2105
// Note: ClickHouse returns the records in no particular order without ORDER BY, so they are ordered by id.
func Test_Issue2105(t *testing.T) {
	table := "issue2105_" + gtime.TimestampNanoStr()
	chCExecSqls(db,
		fmt.Sprintf(`CREATE TABLE %s (id String, json Nullable(String)) ENGINE = MergeTree() ORDER BY id`, table),
		fmt.Sprintf(`INSERT INTO %s VALUES ('1', NULL)`, table),
		fmt.Sprintf(`INSERT INTO %s VALUES ('2', '[{"Name": "任务类型", "Value": "高价值"}, {"Name": "优先级", "Value": "高"}, {"Name": "是否亮点功能", "Value": "是"}]')`, table),
	)
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
		err := db.Model(table).Order("id").Scan(&list)
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
		chCExecSqls(db, fmt.Sprintf(`CREATE TABLE %s (id UInt64 DEFAULT 0) ENGINE = MergeTree() ORDER BY id`, table))
		defer dropTable(table)

		_, err := db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (id) VALUES (18446744073709551615)`, table))
		t.AssertNil(err)

		one, err := db.Model(table).One()
		t.AssertNil(err)
		t.AssertEQ(one["id"].Val(), uint64(18446744073709551615))
	})
}

// https://github.com/gogf/gf/issues/2338
// Note: ClickHouse names a selected column after its table when the joined tables share its name,
// so the nickname of the joined table is read as `b.nickname`.
func Test_Issue2338(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		schemaDb1, schema1 := chCNewSchemaDb()
		defer chCDropSchema(schema1)
		schemaDb2, schema2 := chCNewSchemaDb()
		defer chCDropSchema(schema2)
		table1 := "demo_" + guid.S()
		table2 := "demo_" + guid.S()
		for i, schemaDb := range []gdb.DB{schemaDb1, schemaDb2} {
			chCExecSqls(schemaDb, fmt.Sprintf(`
				CREATE TABLE %s (
					id        UInt32 COMMENT 'User ID',
					nickname  Nullable(String) COMMENT 'User Nickname',
					create_at Nullable(DateTime64(6)) COMMENT 'Created Time',
					update_at Nullable(DateTime64(6)) COMMENT 'Updated Time'
				) ENGINE = MergeTree()
				ORDER BY id`, []string{table1, table2}[i],
			))
		}

		var err error
		_, err = schemaDb1.Model(table1).Insert(g.Map{
			"id":       1,
			"nickname": "name_1",
		})
		t.AssertNil(err)

		_, err = schemaDb2.Model(table2).Insert(g.Map{
			"id":       1,
			"nickname": "name_2",
		})
		t.AssertNil(err)

		tableName1 := fmt.Sprintf(`%s.%s`, schema1, table1)
		tableName2 := fmt.Sprintf(`%s.%s`, schema2, table2)
		all, err := db.Model(tableName1).As(`a`).
			LeftJoin(tableName2+" b", `a.id=b.id`).
			Fields(`a.id`, `b.nickname`).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["b.nickname"], "name_2")
	})

	gtest.C(t, func(t *gtest.T) {
		schemaDb1, schema1 := chCNewSchemaDb()
		defer chCDropSchema(schema1)
		schemaDb2, schema2 := chCNewSchemaDb()
		defer chCDropSchema(schema2)
		table1 := "demo_" + guid.S()
		table2 := "demo_" + guid.S()
		for i, schemaDb := range []gdb.DB{schemaDb1, schemaDb2} {
			chCExecSqls(schemaDb, fmt.Sprintf(`
				CREATE TABLE %s (
					id         UInt32 COMMENT 'User ID',
					nickname   Nullable(String) COMMENT 'User Nickname',
					create_at  Nullable(DateTime64(6)) COMMENT 'Created Time',
					update_at  Nullable(DateTime64(6)) COMMENT 'Updated Time',
					deleted_at Nullable(DateTime64(6)) COMMENT 'Deleted Time'
				) ENGINE = MergeTree()
				ORDER BY id`, []string{table1, table2}[i],
			))
		}

		var err error
		_, err = schemaDb1.Model(table1).Insert(g.Map{
			"id":       1,
			"nickname": "name_1",
		})
		t.AssertNil(err)

		_, err = schemaDb2.Model(table2).Insert(g.Map{
			"id":       1,
			"nickname": "name_2",
		})
		t.AssertNil(err)

		tableName1 := fmt.Sprintf(`%s.%s`, schema1, table1)
		tableName2 := fmt.Sprintf(`%s.%s`, schema2, table2)
		all, err := db.Model(tableName1).As(`a`).
			LeftJoin(tableName2+" b", `a.id=b.id`).
			Fields(`a.id`, `b.nickname`).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["b.nickname"], "name_2")
	})
}

// https://github.com/gogf/gf/issues/2427
func Test_Issue2427(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := "demo_" + guid.S()
		chCExecSqls(db, fmt.Sprintf(`
			CREATE TABLE %s (
				id        UInt32 COMMENT 'User ID',
				passport  String COMMENT 'User Passport',
				password  String COMMENT 'User Password',
				nickname  String COMMENT 'User Nickname',
				create_at Nullable(DateTime64(6)) COMMENT 'Created Time',
				update_at Nullable(DateTime64(6)) COMMENT 'Updated Time'
			) ENGINE = MergeTree()
			ORDER BY id`, table,
		))
		defer dropTable(table)

		_, err1 := db.Model(table).Delete()
		t.Assert(err1, `there should be WHERE condition statement for DELETE operation`)

		_, err2 := db.Model(table).Where(g.Map{}).Delete()
		t.Assert(err2, `there should be WHERE condition statement for DELETE operation`)

		_, err3 := db.Model(table).Where(1).Delete()
		t.AssertNil(err3)
	})
}

// https://github.com/gogf/gf/issues/2561
// Note: ClickHouse reports no affected rows, so the inserted records are counted instead.
func Test_Issue2561(t *testing.T) {
	table := chCCreateNullableTable()
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
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		n, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, 3)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], ``)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)

		one, err = db.Model(table).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `2`)
		t.Assert(one[`passport`], ``)
		t.Assert(one[`password`], `pass_2`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)

		one, err = db.Model(table).WherePri(3).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `3`)
		t.Assert(one[`passport`], ``)
		t.Assert(one[`password`], `pass_3`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)
	})
}

// https://github.com/gogf/gf/issues/2439
// Note: ClickHouse names the columns of `SELECT *` after their tables when the joined tables share
// their names, so the id of the main table is read as `a.id`.
func Test_Issue2439(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		schemaDb, schema := chCNewSchemaDb()
		defer chCDropSchema(schema)
		chCExecSqls(schemaDb,
			`CREATE TABLE a (id Int32) ENGINE = MergeTree() ORDER BY id`,
			`INSERT INTO a (id) VALUES (2)`,
			`CREATE TABLE b (id Int32, name String) ENGINE = MergeTree() ORDER BY id`,
			`INSERT INTO b (id, name) VALUES (2, 'a')`,
			`INSERT INTO b (id, name) VALUES (3, 'b')`,
			`CREATE TABLE c (id Int32) ENGINE = MergeTree() ORDER BY id`,
			`INSERT INTO c (id) VALUES (2)`,
		)

		orm := schemaDb.Model("a")
		orm = orm.InnerJoin(
			"c", "a.id=c.id",
		)
		orm = orm.InnerJoinOnField("b", "id")
		whereFormat := fmt.Sprintf(
			"(`%s`.`%s` LIKE ?) ",
			"b", "name",
		)
		orm = orm.WhereOrf(
			whereFormat,
			"%a%",
		)
		r, err := orm.All()
		t.AssertNil(err)
		t.Assert(len(r), 1)
		t.Assert(r[0]["a.id"], 2)
		t.Assert(r[0]["name"], "a")
	})
}

// https://github.com/gogf/gf/issues/2782
// Note: the driver quotes no identifiers, so the conditions hold bare column names.
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
		t.Assert(condWhere, "(id=?) AND (((nickname=?) OR (password=?))) AND (passport=?)")

		condWhere, _ = m.OmitEmpty().Builder().
			Where("id", "").
			Where(m.Builder().
				Where("nickname", "foo").
				WhereOr("password", "abc123")).
			Where("passport", "pp").
			Build()
		t.Assert(condWhere, "((nickname=?) OR (password=?)) AND (passport=?)")

		condWhere, _ = m.OmitEmpty().Builder().
			Where(m.Builder().
				Where("nickname", "foo").
				WhereOr("password", "abc123")).
			Where("id", "").
			Where("passport", "pp").
			Build()
		t.Assert(condWhere, "((nickname=?) OR (password=?)) AND (passport=?)")
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
		t.Assert(all[0]["id"], 3)
	})
}

// https://github.com/gogf/gf/issues/3086
// Note: the primary key of a MergeTree table is not a unique constraint, so the batch holding a
// duplicate id is inserted instead of failing, and the inserted records are counted as ClickHouse
// reports no affected rows.
func Test_Issue3086(t *testing.T) {
	table := "issue3086_user_" + gtime.TimestampNanoStr()
	chCExecSqls(db, fmt.Sprintf(`
		CREATE TABLE %s (
			id        UInt32 COMMENT 'User ID',
			passport  String COMMENT 'User Passport',
			password  Nullable(String) COMMENT 'User Password',
			nickname  Nullable(String) COMMENT 'User Nickname',
			create_at Nullable(DateTime) COMMENT 'Created Time',
			update_at Nullable(DateTime) COMMENT 'Updated Time'
		) ENGINE = MergeTree()
		ORDER BY id`, table,
	))
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
		t.AssertNil(err)

		n, err := db.Model(table).Where("id", 1).Count()
		t.AssertNil(err)
		t.Assert(n, 2)
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
		_, err := db.Model(table).Data(data).Batch(10).Insert()
		t.AssertNil(err)

		n, err := db.Model(table).WhereIn("id", g.Slice{3, 4}).Count()
		t.AssertNil(err)
		t.Assert(n, 2)
	})
}

// https://github.com/gogf/gf/issues/3204
// Note: ClickHouse reports no last insert id, so InsertAndGetId must fail as unsupported, as
// db.InsertAndGetId does, and the omitted fields are checked on a plain Insert.
func Test_Issue3204(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

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
		t.Assert(all[0]["id"], 2)
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
		// Note: the driver sends the inserted records in a batch of the underlying driver, which
		// CatchSQL does not see, as the core catches the statements in DoQuery and DoExec only;
		// the omitted fields are checked on the inserted record instead.
		sqlArray, err = gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err = db.Ctx(ctx).Model(table).Data(data).Insert()
			return err
		})
		t.AssertNil(err)
		t.Assert(len(sqlArray), 0)
		one, err := db.Model(table).Where("id", 20).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "passport_20")
		t.Assert(one["nickname"], "")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any `orm:"id,omitempty"`
			Passport   any `orm:"passport,omitempty"`
			Password   any `orm:"password,omitempty"`
			Nickname   any `orm:"nickname,omitempty"`
			CreateTime any `orm:"create_time,omitempty"`
		}
		insertId, err := db.Model(table).Data(User{
			Id:       21,
			Passport: "passport_21",
		}).InsertAndGetId()
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "LastInsertId is not supported by ClickHouse")
		t.Assert(insertId, 0)
		count, err := db.Model(table).Where("id", 21).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
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
			gstr.Contains(sqlArray[len(sqlArray)-1], "passport='passport_1' WHERE id=1"),
			true,
		)
	})
}

// https://github.com/gogf/gf/issues/3218
func Test_Issue3218(t *testing.T) {
	table := "issue3218_sys_config_" + gtime.TimestampNanoStr()
	chCExecSqls(db,
		fmt.Sprintf(`
			CREATE TABLE %s (
				id         Int32,
				name       Nullable(String) COMMENT '配置名称',
				value      Nullable(String) COMMENT '配置值',
				created_at Nullable(DateTime) COMMENT '创建时间',
				updated_at Nullable(DateTime) COMMENT '更新时间'
			) ENGINE = MergeTree()
			ORDER BY id`, table,
		),
		fmt.Sprintf(
			`INSERT INTO %s VALUES (49, 'site', '{"banned_ip":"22","filings":"2222","fixed_page":"","site_name":"22","version":"22"}', '2023-12-19 14:08:25', '2023-12-19 14:08:25')`,
			table,
		),
	)
	defer dropTable(table)
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
// Note: the driver reads the fields of a table from `system`.columns rather than SHOW FULL COLUMNS.
func Test_Issue2552_ClearTableFieldsAll(t *testing.T) {
	table := chCCreateNullableTable()
	defer dropTable(table)

	showTableKey := "`system`.columns"
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err := db.Model(table).Ctx(ctx).Insert(g.Map{
				"passport":    guid.S(),
				"password":    guid.S(),
				"nickname":    guid.S(),
				"create_time": gtime.NewFromStr(chCCreateTime).String(),
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

		_, err = db.Exec(ctx, fmt.Sprintf("alter table %s drop column `nickname`", table))
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
// Note: the driver reads the fields of a table from `system`.columns rather than SHOW FULL COLUMNS.
func Test_Issue2552_ClearTableFields(t *testing.T) {
	table := chCCreateNullableTable()
	defer dropTable(table)

	showTableKey := "`system`.columns"
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err := db.Model(table).Ctx(ctx).Insert(g.Map{
				"passport":    guid.S(),
				"password":    guid.S(),
				"nickname":    guid.S(),
				"create_time": gtime.NewFromStr(chCCreateTime).String(),
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

		_, err = db.Exec(ctx, fmt.Sprintf("alter table %s drop column `nickname`", table))
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
// Note: the driver quotes no identifiers, so the table name is bare in the statements.
func Test_Issue2643(t *testing.T) {
	table := "issue2643_" + gtime.TimestampNanoStr()
	chCExecSqls(db, fmt.Sprintf(`
		CREATE TABLE %s (
			id    Nullable(Int32),
			name  Nullable(String),
			value Nullable(Int32),
			dept  Nullable(String)
		) ENGINE = MergeTree()
		ORDER BY tuple()`, table,
	))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			expectKey1 = fmt.Sprintf("SELECT s.name,replace(concat_ws(',',lpad(s.id, 6, '0'),s.name),',','') `code` FROM %s AS s", table)
			expectKey2 = fmt.Sprintf("SELECT CASE WHEN dept='物资部' THEN '物资部' ELSE '其他' END dept,sum(s.value) FROM %s AS s GROUP BY CASE WHEN dept='物资部' THEN '物资部' ELSE '其他' END", table)
		)
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			db.Ctx(ctx).Model(table).As("s").Fields(
				"s.name",
				"replace(concat_ws(',',lpad(s.id, 6, '0'),s.name),',','') `code`",
			).All()
			db.Ctx(ctx).Model(table).As("s").Fields(
				"CASE WHEN dept='物资部' THEN '物资部' ELSE '其他' END dept",
				"sum(s.value)",
			).Group("CASE WHEN dept='物资部' THEN '物资部' ELSE '其他' END").All()
			return nil
		})
		t.AssertNil(err)
		sqlContent := gstr.Join(sqlArray, "\n")
		t.Assert(gstr.Contains(sqlContent, expectKey1), true)
		t.Assert(gstr.Contains(sqlContent, expectKey2), true)
	})
}

// https://github.com/gogf/gf/issues/3238
func Test_Issue3238(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			_, err := db.Ctx(ctx).Model(table).Hook(gdb.HookHandler{
				Select: func(ctx context.Context, in *gdb.HookSelectInput) (result gdb.Result, err error) {
					result, err = in.Next(ctx)
					if err != nil {
						return
					}
					var wg sync.WaitGroup
					for _, record := range result {
						wg.Add(1)
						go func(record gdb.Record) {
							defer wg.Done()
							id, _ := db.Ctx(ctx).Model(table).WherePri(1).Value(`id`)
							nickname, _ := db.Ctx(ctx).Model(table).WherePri(1).Value(`nickname`)
							t.Assert(id.Int(), 1)
							t.Assert(nickname.String(), "name_1")
						}(record)
					}
					wg.Wait()
					return
				},
			},
			).All()
			t.AssertNil(err)
		}
	})
}

// https://github.com/gogf/gf/issues/3649
// Note: the driver quotes no identifiers, so the statement holds bare names.
func Test_Issue3649(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		sql, err := gdb.CatchSQL(context.Background(), func(ctx context.Context) (err error) {
			user := db.Model(table).Ctx(ctx)
			_, err = user.Where("create_time = ?", gdb.Raw("now()")).WhereLT("create_time", gdb.Raw("now()")).Count()
			return
		})
		t.AssertNil(err)
		sqlStr := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE (create_time = now()) AND (create_time < now())", table)
		t.Assert(sql[0], sqlStr)
	})
}

// https://github.com/gogf/gf/issues/3754
// Note: ClickHouse reports no affected rows, so the records are checked instead.
func Test_Issue3754(t *testing.T) {
	table := "issue3754_" + gtime.TimestampNanoStr()
	chCExecSqls(db, fmt.Sprintf(`
		CREATE TABLE %s (
			id        Int32,
			name      Nullable(String),
			create_at Nullable(DateTime),
			update_at Nullable(DateTime),
			delete_at Nullable(DateTime)
		) ENGINE = MergeTree()
		ORDER BY id`, table,
	))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fieldsEx := []string{"delete_at", "create_at", "update_at"}
		// Insert.
		dataInsert := g.Map{
			"id":   1,
			"name": "name_1",
		}
		_, err := db.Model(table).Data(dataInsert).FieldsEx(fieldsEx).Insert()
		t.AssertNil(err)

		oneInsert, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["name"].String(), "name_1")
		t.Assert(oneInsert["delete_at"].String(), "")
		t.Assert(oneInsert["create_at"].String(), "")
		t.Assert(oneInsert["update_at"].String(), "")

		// Update.
		dataUpdate := g.Map{
			"name": "name_1000",
		}
		_, err = db.Model(table).Data(dataUpdate).FieldsEx(fieldsEx).WherePri(1).Update()
		t.AssertNil(err)

		oneUpdate, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["delete_at"].String(), "")
		t.Assert(oneUpdate["create_at"].String(), "")
		t.Assert(oneUpdate["update_at"].String(), "")

		// FieldsEx does not affect Delete operation.
		_, err = db.Model(table).FieldsEx(fieldsEx).WherePri(1).Delete()
		t.AssertNil(err)
		count, err := db.Model(table).WherePri(1).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
		oneDeleteUnscoped, err := db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneDeleteUnscoped["id"].Int(), 1)
		t.Assert(oneDeleteUnscoped["name"].String(), "name_1000")
		t.AssertNE(oneDeleteUnscoped["delete_at"].String(), "")
		t.Assert(oneDeleteUnscoped["create_at"].String(), "")
		t.Assert(oneDeleteUnscoped["update_at"].String(), "")
	})
}

// https://github.com/gogf/gf/issues/3626
// Note: ClickHouse reports no affected rows, so the inserted record is checked instead.
func Test_Issue3626(t *testing.T) {
	table := "issue3626_" + gtime.TimestampNanoStr()
	chCExecSqls(db, fmt.Sprintf(`
		CREATE TABLE %s (
			id   Int32,
			name Nullable(String)
		) ENGINE = MergeTree()
		ORDER BY id`, table,
	))
	defer dropTable(table)

	// Insert.
	gtest.C(t, func(t *gtest.T) {
		dataInsert := g.Map{
			"id":   1,
			"name": "name_1",
		}
		_, err := db.Model(table).Data(dataInsert).Insert()
		t.AssertNil(err)

		oneInsert, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["name"].String(), "name_1")
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
		t.Assert(one["id"], 10)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Order("id desc").One()
		t.AssertNil(err)
		t.Assert(one["id"], 10)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Order("id desc, nickname asc").One()
		t.AssertNil(err)
		t.Assert(one["id"], 10)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Order("id desc", "nickname asc").One()
		t.AssertNil(err)
		t.Assert(one["id"], 10)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Order("id desc").Order("nickname asc").One()
		t.AssertNil(err)
		t.Assert(one["id"], 10)
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
	table := "issue3915_" + gtime.TimestampNanoStr()
	chCExecSqls(db,
		fmt.Sprintf(`
			CREATE TABLE %s (
				id UInt32 COMMENT 'user id',
				a  Nullable(Float32) COMMENT 'user name',
				b  Nullable(Float32) COMMENT 'user status'
			) ENGINE = MergeTree()
			ORDER BY id`, table,
		),
		fmt.Sprintf("INSERT INTO %s (`id`,`a`,`b`) VALUES (1,1,2)", table),
		fmt.Sprintf("INSERT INTO %s (`id`,`a`,`b`) VALUES (2,5,4)", table),
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("a < b").All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 1)

		all, err = db.Model(table).Where(gdb.Raw("a < b")).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 1)

		all, err = db.Model(table).WhereLT("a", gdb.Raw("`b`")).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 1)
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("a > b").All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 2)

		all, err = db.Model(table).Where(gdb.Raw("a > b")).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 2)

		all, err = db.Model(table).WhereGT("a", gdb.Raw("`b`")).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 2)
	})
}

type ChCRoleBase struct {
	gmeta.Meta  `orm:"table:sys_role"`
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Description string      `json:"description"`
	Weight      int         `json:"weight"`
	StatusId    int         `json:"statusId"`
	CreatedAt   *gtime.Time `json:"createdAt"`
	UpdatedAt   *gtime.Time `json:"updatedAt"`
}

type ChCRole struct {
	gmeta.Meta `orm:"table:sys_role"`
	ChCRoleBase
	Id     uint       `json:"id"`
	Status *ChCStatus `json:"status" orm:"with:id=status_id"`
}

type ChCStatusBase struct {
	gmeta.Meta `orm:"table:sys_status"`
	En         string `json:"en"`
	Cn         string `json:"cn"`
	Weight     int    `json:"weight"`
}

type ChCStatus struct {
	gmeta.Meta `orm:"table:sys_status"`
	ChCStatusBase
	Id uint `json:"id"`
}

// https://github.com/gogf/gf/issues/2119
// Note: ClickHouse returns the records in no particular order without ORDER BY, so they are
// ordered by id, and the tables fixed by the orm tags are created in a database of their own.
func Test_Issue2119(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		schemaDb, schema := chCNewSchemaDb()
		defer chCDropSchema(schema)
		chCExecSqls(schemaDb,
			`CREATE TABLE sys_role (
				id          UInt32,
				name        String DEFAULT '' COMMENT '角色名称||s,r',
				code        String DEFAULT '' COMMENT '角色 code||s,r',
				description String DEFAULT '' COMMENT '描述信息|text',
				weight      UInt32 DEFAULT 0 COMMENT '排序||r|min:0#发布状态不能小于 0',
				status_id   UInt32 DEFAULT 1 COMMENT '发布状态|hasOne|f:status,fk:id',
				created_at  Nullable(DateTime),
				updated_at  Nullable(DateTime)
			) ENGINE = MergeTree()
			ORDER BY id`,
			`INSERT INTO sys_role VALUES (1, '开发人员', 'developer', '123123', 900, 2, '2022-09-03 21:25:03', '2022-09-09 23:35:23')`,
			`INSERT INTO sys_role VALUES (2, '管理员', 'admin', '', 800, 1, '2022-09-03 21:25:03', '2022-09-09 23:00:17')`,
			`INSERT INTO sys_role VALUES (3, '运营', 'operator', '', 700, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03')`,
			`INSERT INTO sys_role VALUES (4, '客服', 'service', '', 600, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03')`,
			`INSERT INTO sys_role VALUES (5, '收银', 'account', '', 500, 1, '2022-09-03 21:25:03', '2022-09-03 21:25:03')`,
			`CREATE TABLE sys_status (
				id     UInt32,
				en     String DEFAULT '' COMMENT '英文名称',
				cn     String DEFAULT '' COMMENT '中文名称',
				weight UInt32 DEFAULT 0 COMMENT '排序权重'
			) ENGINE = MergeTree()
			ORDER BY id`,
			`INSERT INTO sys_status VALUES (1, 'on line', '上线', 900)`,
			`INSERT INTO sys_status VALUES (2, 'undecided', '未决定', 800)`,
			`INSERT INTO sys_status VALUES (3, 'off line', '下线', 700)`,
		)

		roles := make([]*ChCRole, 0)
		err := schemaDb.Ctx(context.Background()).Model(&ChCRole{}).WithAll().Order("id").Scan(&roles)
		t.AssertNil(err)
		expectStatus := []*ChCStatus{
			{
				ChCStatusBase: ChCStatusBase{
					En:     "undecided",
					Cn:     "未决定",
					Weight: 800,
				},
				Id: 2,
			},
			{
				ChCStatusBase: ChCStatusBase{
					En:     "on line",
					Cn:     "上线",
					Weight: 900,
				},
				Id: 1,
			},
			{
				ChCStatusBase: ChCStatusBase{
					En:     "on line",
					Cn:     "上线",
					Weight: 900,
				},
				Id: 1,
			},
			{
				ChCStatusBase: ChCStatusBase{
					En:     "on line",
					Cn:     "上线",
					Weight: 900,
				},
				Id: 1,
			},
			{
				ChCStatusBase: ChCStatusBase{
					En:     "on line",
					Cn:     "上线",
					Weight: 900,
				},
				Id: 1,
			},
		}

		t.Assert(len(roles), len(expectStatus))
		for i := 0; i < len(roles); i++ {
			t.Assert(roles[i].Status, expectStatus[i])
		}
	})
}

// https://github.com/gogf/gf/issues/4034
// Note: ClickHouse has no transactions, so Transaction fails as unsupported before the Save.
func Test_Issue4034(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := "issue4034_" + gtime.TimestampNanoStr()
		chCExecSqls(db, fmt.Sprintf(`
			CREATE TABLE %s (
				id         UInt32,
				passport   Nullable(String),
				password   Nullable(String),
				nickname   Nullable(String),
				created_at DateTime DEFAULT now(),
				updated_at DateTime DEFAULT now()
			) ENGINE = MergeTree()
			ORDER BY id`, table,
		))
		defer dropTable(table)

		err := chCIssue4034SaveDeviceAndToken(ctx, table)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:Transaction")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

func chCIssue4034SaveDeviceAndToken(ctx context.Context, table string) error {
	return db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := chCIssue4034SaveAppDevice(ctx, table, tx); err != nil {
			return err
		}
		return nil
	})
}

func chCIssue4034SaveAppDevice(ctx context.Context, table string, tx gdb.TX) error {
	_, err := db.Model(table).Safe().Ctx(ctx).TX(tx).Data(g.Map{
		"passport": "111",
		"password": "222",
		"nickname": "333",
	}).Save()
	return err
}

// https://github.com/gogf/gf/issues/4086
// Note: ClickHouse returns the records in no particular order without ORDER BY, so they are ordered by proxy_id.
func Test_Issue4086(t *testing.T) {
	table := "issue4086_" + gtime.TimestampNanoStr()
	defer dropTable(table)
	chCExecSqls(db,
		fmt.Sprintf(`
			CREATE TABLE %s (
				proxy_id      Int64,
				recommend_ids Nullable(String),
				photos        Nullable(String)
			) ENGINE = MergeTree()
			ORDER BY proxy_id`, table,
		),
		fmt.Sprintf("INSERT INTO %s (`proxy_id`, `recommend_ids`, `photos`) VALUES (1, '[584, 585]', 'null')", table),
		fmt.Sprintf("INSERT INTO %s (`proxy_id`, `recommend_ids`, `photos`) VALUES (2, '[]', NULL)", table),
	)

	gtest.C(t, func(t *gtest.T) {
		type ProxyParam struct {
			ProxyId      int64   `json:"proxyId" orm:"proxy_id"`
			RecommendIds []int64 `json:"recommendIds" orm:"recommend_ids"`
			Photos       []int64 `json:"photos" orm:"photos"`
		}

		var proxyParamList []*ProxyParam
		err := db.Model(table).Ctx(ctx).Order("proxy_id").Scan(&proxyParamList)
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
		err := db.Model(table).Ctx(ctx).Order("proxy_id").Scan(&proxyParamList)
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
		err := db.Model(table).Ctx(ctx).Order("proxy_id").Scan(&proxyParamList)
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
		err := db.Model(table).Ctx(ctx).Order("proxy_id").Scan(&proxyParamList)
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
		err := db.Model(table).Ctx(ctx).Order("proxy_id").Scan(&proxyParamList)
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
		err := db.Model(table).Ctx(ctx).Order("proxy_id").Scan(&proxyParamList)
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
func Test_Issue4500(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// Raw SQL with WHERE + external Where condition + Count.
	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id IN (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("id", 8).
			Count()
		t.AssertNil(err)
		t.Assert(count, 3)
	})

	// Raw SQL without WHERE + external Where condition + Count.
	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s", table)).
			WhereLT("id", 5).
			Count()
		t.AssertNil(err)
		t.Assert(count, 4)
	})

	// Raw + Where + ScanAndCount.
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
		t.Assert(len(users), 3)
		t.Assert(total, 3)
	})

	// Raw + multiple Where conditions + Count.
	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id > ?", table), 0).
			WhereLT("id", 5).
			WhereGTE("id", 2).
			Count()
		t.AssertNil(err)
		t.Assert(count, 3)
	})

	// Raw SQL with no external Where + Count.
	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id IN (?)", table), g.Slice{1, 2, 3}).
			Count()
		t.AssertNil(err)
		t.Assert(count, 3)
	})

	// All with Raw + Where.
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
	table := chCCreateNullableInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields("").Limit(1).All()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(len(result[0]), 6)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields("", "id").Limit(1).All()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(len(result[0]), 1)
		t.AssertNE(result[0]["id"], nil)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields("id", "", "nickname").Limit(1).All()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(len(result[0]), 2)
		t.AssertNE(result[0]["id"], nil)
		t.AssertNE(result[0]["nickname"], nil)
	})
}

// https://github.com/gogf/gf/issues/4698
func Test_Issue4698(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// AllAndCount with multiple fields should generate valid COUNT SQL.
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id, nickname").AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(len(result), TableSize)
		t.AssertNE(result[0]["id"], nil)
		t.AssertNE(result[0]["nickname"], nil)
		t.Assert(result[0]["passport"], nil)
	})

	// AllAndCount(false) with multiple fields.
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id, nickname").AllAndCount(false)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(len(result), TableSize)
	})

	// ScanAndCount with multiple fields.
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

	// AllAndCount with single field and useFieldForCount=true.
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id").AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(len(result), TableSize)
		t.Assert(len(result[0]), 1)
	})

	// AllAndCount with Where condition.
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id, nickname").Where("id<?", 5).AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, 4)
		t.Assert(len(result), 4)
	})

	// Distinct + AllAndCount(false) should use COUNT(1), not COUNT(DISTINCT 1).
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("nickname").Distinct().AllAndCount(false)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.AssertGT(len(result), 0)
	})

	// Distinct + AllAndCount(true) with single field should use COUNT(DISTINCT nickname).
	gtest.C(t, func(t *gtest.T) {
		_, count, err := db.Model(table).Fields("nickname").Distinct().AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})

	// Distinct + multiple fields + AllAndCount(true) should fallback to COUNT(1).
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id, nickname").Distinct().AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(len(result), TableSize)
	})
}

// Test_ClickHouse_TableFields_OtherSchema tests that the fields of a table do not include the
// columns of a table with the same name in another database.
func Test_ClickHouse_TableFields_OtherSchema(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		schemaDb1, schema1 := chCNewSchemaDb()
		defer chCDropSchema(schema1)
		schemaDb2, schema2 := chCNewSchemaDb()
		defer chCDropSchema(schema2)
		table := "same_name_" + gtime.TimestampNanoStr()
		chCExecSqls(schemaDb1, fmt.Sprintf(`CREATE TABLE %s (id UInt32, name String) ENGINE = MergeTree() ORDER BY id`, table))
		chCExecSqls(schemaDb2, fmt.Sprintf(`CREATE TABLE %s (id UInt32, other String, extra String) ENGINE = MergeTree() ORDER BY id`, table))

		fields, err := schemaDb1.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(len(fields), 2)
		t.AssertNE(fields["name"], nil)
		t.Assert(fields["other"], nil)

		fields, err = db.TableFields(ctx, table, schema2)
		t.AssertNil(err)
		t.Assert(len(fields), 3)
		t.AssertNE(fields["extra"], nil)
		t.Assert(fields["name"], nil)
	})
}
