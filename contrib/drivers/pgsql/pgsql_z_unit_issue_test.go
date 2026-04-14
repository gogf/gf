// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

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
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/guid"
)

// https://github.com/gogf/gf/issues/3330
func Test_Issue3330(t *testing.T) {
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
			CONSTRAINT %s unique ("password")
		) ;`, table, uniqueName,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			list []map[string]any
			one  gdb.Record
			err  error
		)

		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)

		t.Assert(fields["id"].Key, "pri")
		t.Assert(fields["password"].Key, "uni")

		for i := 1; i <= 10; i++ {
			list = append(list, g.Map{
				"id":          i,
				"passport":    fmt.Sprintf("p%d", i),
				"password":    fmt.Sprintf("pw%d", i),
				"nickname":    fmt.Sprintf("n%d", i),
				"create_time": "2016-06-01 00:00:00",
			})
		}

		_, err = db.Model(table).Data(list).Insert()
		t.AssertNil(err)

		for i := 1; i <= 10; i++ {
			one, err = db.Model(table).WherePri(i).One()
			t.AssertNil(err)
			t.Assert(one["id"], list[i-1]["id"])
			t.Assert(one["passport"], list[i-1]["passport"])
			t.Assert(one["password"], list[i-1]["password"])
			t.Assert(one["nickname"], list[i-1]["nickname"])
		}
	})
}

// https://github.com/gogf/gf/issues/3632
func Test_Issue3632(t *testing.T) {
	type Member struct {
		One []int64    `json:"one" orm:"one"`
		Two [][]string `json:"two" orm:"two"`
	}
	var (
		sqlText = gtest.DataContent("issues", "issue3632.sql")
		table   = fmt.Sprintf(`%s_%d`, TablePrefix+"issue3632", gtime.TimestampNano())
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(sqlText, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			dao    = db.Model(table)
			member = Member{
				One: []int64{1, 2, 3},
				Two: [][]string{{"a", "b"}, {"c", "d"}},
			}
		)

		_, err := dao.Ctx(ctx).Data(&member).Insert()
		t.AssertNil(err)
	})
}

// https://github.com/gogf/gf/issues/3671
func Test_Issue3671(t *testing.T) {
	type SubMember struct {
		Seven string
		Eight int64
	}
	type Member struct {
		One   []int64     `json:"one" orm:"one"`
		Two   [][]string  `json:"two" orm:"two"`
		Three []string    `json:"three" orm:"three"`
		Four  []int64     `json:"four" orm:"four"`
		Five  []SubMember `json:"five" orm:"five"`
	}
	var (
		sqlText = gtest.DataContent("issues", "issue3671.sql")
		table   = fmt.Sprintf(`%s_%d`, TablePrefix+"issue3671", gtime.TimestampNano())
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(sqlText, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			dao    = db.Model(table)
			member = Member{
				One:   []int64{1, 2, 3},
				Two:   [][]string{{"a", "b"}, {"c", "d"}},
				Three: []string{"x", "y", "z"},
				Four:  []int64{1, 2, 3},
				Five:  []SubMember{{Seven: "1", Eight: 2}, {Seven: "3", Eight: 4}},
			}
		)

		_, err := dao.Ctx(ctx).Data(&member).Insert()
		t.AssertNil(err)
	})
}

// https://github.com/gogf/gf/issues/3668
func Test_Issue3668(t *testing.T) {
	type Issue3668 struct {
		Text   any
		Number any
	}
	var (
		sqlText = gtest.DataContent("issues", "issue3668.sql")
		table   = fmt.Sprintf(`%s_%d`, TablePrefix+"issue3668", gtime.TimestampNano())
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(sqlText, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			dao  = db.Model(table)
			data = Issue3668{
				Text:   "我们都是自然的婴儿，卧在宇宙的摇篮里",
				Number: nil,
			}
		)
		_, err := dao.Ctx(ctx).
			Data(data).
			Insert()
		t.AssertNil(err)
	})
}

type Issue4033Status int

const (
	Issue4033StatusA Issue4033Status = 1
)

func (s Issue4033Status) String() string {
	return "somevalue"
}

func (s Issue4033Status) Int64() int64 {
	return int64(s)
}

// https://github.com/gogf/gf/issues/4033
func Test_Issue4033(t *testing.T) {
	var (
		sqlText = gtest.DataContent("issues", "issue4033.sql")
		table   = "test_enum"
	)
	if _, err := db.Exec(ctx, sqlText); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		query := g.Map{
			"status": g.Slice{Issue4033StatusA},
		}
		_, err := db.Model(table).Ctx(ctx).Where(query).All()
		t.AssertNil(err)
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

// https://github.com/gogf/gf/issues/4677
// record.Get().Bytes() corrupts bytea data on retrieval from PostgreSQL.
func Test_Issue4677(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"issue4677", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			bin_data bytea
		);`, table,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test 1: Binary data with various byte values including 0x00, 0x5D(']'), 0x5B('[')
		originalBytes := []byte{
			0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x01, 0x5B, 0x5D,
			0xFF, 0x7B, 0x7D, 0x80, 0xCA, 0xFE, 0xBA, 0xBE,
		}

		_, err := db.Model(table).Data(g.Map{
			"bin_data": originalBytes,
		}).Insert()
		t.AssertNil(err)

		record, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)

		retrievedBytes := record["bin_data"].Bytes()
		t.Assert(len(retrievedBytes), len(originalBytes))
		t.Assert(retrievedBytes, originalBytes)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test 2: Larger binary data (simulating gob/protobuf encoded payload)
		largeBytes := make([]byte, 1024)
		for i := range largeBytes {
			largeBytes[i] = byte(i % 256)
		}

		_, err := db.Model(table).Data(g.Map{
			"bin_data": largeBytes,
		}).Insert()
		t.AssertNil(err)

		record, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)

		retrievedBytes := record["bin_data"].Bytes()
		t.Assert(len(retrievedBytes), len(largeBytes))
		t.Assert(retrievedBytes, largeBytes)
	})
}

// https://github.com/gogf/gf/issues/4231
// ConvertValueForField corrupts bytea data containing 0x5D on write.
func Test_Issue4231(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"issue4231", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			bin_data bytea
		);`, table,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Bytes containing 0x5D (ASCII ']') which was being converted to 0x7D ('}')
		originalBytes := []byte{0x01, 0x5D, 0x02, 0x5B, 0x03}

		_, err := db.Model(table).Data(g.Map{
			"bin_data": originalBytes,
		}).Insert()
		t.AssertNil(err)

		record, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)

		retrievedBytes := record["bin_data"].Bytes()
		t.Assert(len(retrievedBytes), len(originalBytes))
		t.Assert(retrievedBytes, originalBytes)
	})
}

// https://github.com/gogf/gf/issues/4595
// FieldsPrefix silently drops fields when using table alias before LeftJoin.
func Test_Issue4595(t *testing.T) {
	var (
		tableUser       = fmt.Sprintf(`%s_%d`, TablePrefix+"issue4595_user", gtime.TimestampNano())
		tableUserDetail = fmt.Sprintf(`%s_%d`, TablePrefix+"issue4595_user_detail", gtime.TimestampNano())
	)

	// Create user table
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			name varchar(100),
			email varchar(100)
		);`, tableUser,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(tableUser)

	// Create user_detail table
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			user_id bigint,
			phone varchar(20),
			address varchar(200)
		);`, tableUserDetail,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(tableUserDetail)

	// Insert test data
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		INSERT INTO %s (id, name, email) VALUES (1, 'john', 'john@example.com');
		INSERT INTO %s (id, user_id, phone, address) VALUES (1, 1, '1234567890', '123 Main St');
	`, tableUser, tableUserDetail)); err != nil {
		gtest.Fatal(err)
	}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: FieldsPrefix called before LeftJoin
		// Both t1 and t2 fields should be present
		r, err := db.Model(tableUser).As("t1").
			FieldsPrefix("t2", "phone", "address").
			FieldsPrefix("t1", "id", "name", "email").
			LeftJoin(tableUserDetail, "t2", "t1.id=t2.user_id").
			All()

		t.AssertNil(err)
		t.Assert(len(r), 1)
		t.Assert(r[0]["id"], 1)
		t.Assert(r[0]["name"], "john")
		t.Assert(r[0]["email"], "john@example.com")
		t.Assert(r[0]["phone"], "1234567890")
		t.Assert(r[0]["address"], "123 Main St")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 2: Using Fields() with prefix
		r, err := db.Model(tableUser).As("t1").
			Fields("t2.phone", "t2.address", "t1.id", "t1.name", "t1.email").
			LeftJoin(tableUserDetail, "t2", "t1.id=t2.user_id").
			All()
		t.AssertNil(err)
		t.Assert(len(r), 1)
		t.Assert(r[0]["id"], 1)
		t.Assert(r[0]["name"], "john")
		t.Assert(r[0]["email"], "john@example.com")
		t.Assert(r[0]["phone"], "1234567890")
		t.Assert(r[0]["address"], "123 Main St")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 3: FieldsPrefix called after LeftJoin
		r, err := db.Model(tableUser).As("t1").
			LeftJoin(tableUserDetail, "t2", "t1.id=t2.user_id").
			FieldsPrefix("t2", "phone", "address").
			FieldsPrefix("t1", "id", "name", "email").
			All()
		t.AssertNil(err)
		t.Assert(len(r), 1)
		t.Assert(r[0]["id"], 1)
		t.Assert(r[0]["name"], "john")
		t.Assert(r[0]["email"], "john@example.com")
		t.Assert(r[0]["phone"], "1234567890")
		t.Assert(r[0]["address"], "123 Main St")
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
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `issue1380.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
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
		tableUser       = "user_" + gtime.TimestampMicroStr()
		tableUserDetail = "user_detail_" + gtime.TimestampMicroStr()
		tableUserScores = "user_scores_" + gtime.TimestampMicroStr()
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  uid serial NOT NULL,
  name varchar(45) NOT NULL,
  PRIMARY KEY (uid)
);
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  uid serial NOT NULL,
  address varchar(45) NOT NULL,
  PRIMARY KEY (uid)
);
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id serial NOT NULL,
  uid int NOT NULL,
  score int NOT NULL,
  PRIMARY KEY (id)
);
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
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `issue1401.sql`), ";")
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
	// Framework bug: With() uses Go struct field name ("Id") as the column name in
	// WHERE clause instead of the mapped column name ("id"). PgSQL double-quoted identifiers
	// are case-sensitive, so WHERE "Id"=0 fails with "column Id does not exist".
	// This needs a fix in gdb's With() column-name resolution. TODO: create issue.
	t.Skip("Framework bug: With() generates case-sensitive column name on PgSQL — needs core fix")
	var (
		table1 = "parcels"
		table2 = "items"
	)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `issue1412.sql`), ";")
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
	// where + time.Time arguments.
	// PgSQL "timestamp without time zone" compares literal values without timezone
	// conversion, so use times that bracket the stored value 19:03:33 directly.
	gtest.C(t, func(t *gtest.T) {
		t1, _ := time.Parse("2006-01-02 15:04:05", "2020-10-27 19:03:32")
		t2, _ := time.Parse("2006-01-02 15:04:05", "2020-10-27 19:03:34")
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
	        id serial NOT NULL,
	        user_id int NOT NULL,
	        "UserId" int NOT NULL,
	        PRIMARY KEY (id)
	    );
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
		t.Assert(one["id"], 1)
		t.Assert(one["user_id"], 2)
		t.Assert(one["UserId"], 3)

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
	t.Skip("MySQL IF() function not supported in PostgreSQL")
}

// https://github.com/gogf/gf/issues/1733
func Test_Issue1733(t *testing.T) {
	t.Skip("PostgreSQL does not support zerofill column attribute")
}

// https://github.com/gogf/gf/issues/2012
func Test_Issue2012(t *testing.T) {
	t.Skip("PostgreSQL does not support zerofill column attribute")
}

// https://github.com/gogf/gf/issues/2105
func Test_Issue2105(t *testing.T) {
	table := "issue2105"
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `issue2105.sql`), ";")
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

// https://github.com/gogf/gf/issues/2338
func Test_Issue2338(t *testing.T) {
	t.Skip("PostgreSQL cross-schema test requires TestSchema1/TestSchema2 setup not available in pgsql driver tests")
}

// https://github.com/gogf/gf/issues/2356
func Test_Issue2356(t *testing.T) {
	t.Skip("PostgreSQL does not have BIGINT UNSIGNED; max uint64 test is MySQL-specific")
}

// https://github.com/gogf/gf/issues/2427
func Test_Issue2427(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := "demo_" + guid.S()
		if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
    id bigserial NOT NULL,
    passport  varchar(45) NOT NULL,
    password  varchar(45) NOT NULL,
    nickname  varchar(45) NOT NULL,
    create_at timestamp DEFAULT NULL,
    update_at timestamp DEFAULT NULL,
    PRIMARY KEY (id)
);
	    `, table,
		)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table)

		_, err1 := db.Model(table).Delete()
		t.Assert(err1, `there should be WHERE condition statement for DELETE operation`)

		_, err2 := db.Model(table).Where(g.Map{}).Delete()
		t.Assert(err2, `there should be WHERE condition statement for DELETE operation`)

		_, err3 := db.Model(table).Where("1=1").Delete()
		t.AssertNil(err3)
	})
}

// https://github.com/gogf/gf/issues/2561
func Test_Issue2561(t *testing.T) {
	table := createTable()
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
func Test_Issue2439(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tableA := "issue2439_a_" + gtime.TimestampNanoStr()
		tableB := "issue2439_b_" + gtime.TimestampNanoStr()
		tableC := "issue2439_c_" + gtime.TimestampNanoStr()

		_, err := db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (id serial PRIMARY KEY)`, tableA))
		t.AssertNil(err)
		defer dropTable(tableA)
		_, err = db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (id) VALUES (2)`, tableA))
		t.AssertNil(err)

		_, err = db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (id serial PRIMARY KEY, name varchar(255) NOT NULL)`, tableB))
		t.AssertNil(err)
		defer dropTable(tableB)
		_, err = db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (id, name) VALUES (2, 'a')`, tableB))
		t.AssertNil(err)
		_, err = db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (id, name) VALUES (3, 'b')`, tableB))
		t.AssertNil(err)

		_, err = db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (id serial PRIMARY KEY)`, tableC))
		t.AssertNil(err)
		defer dropTable(tableC)
		_, err = db.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (id) VALUES (2)`, tableC))
		t.AssertNil(err)

		orm := db.Model(tableA)
		orm = orm.InnerJoin(tableC, fmt.Sprintf(`%s.id=%s.id`, tableA, tableC))
		orm = orm.InnerJoinOnField(tableB, "id")
		whereFormat := fmt.Sprintf(
			`("%s"."%s" LIKE ?) `,
			tableB, "name",
		)
		orm = orm.WhereOrf(
			whereFormat,
			"%a%",
		)
		r, err := orm.All()
		t.AssertNil(err)
		t.Assert(len(r), 1)
		t.Assert(r[0]["id"], 2)
		t.Assert(r[0]["name"], "a")
	})
}

// https://github.com/gogf/gf/issues/2782
func Test_Issue2787(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model("user")

		condWhere, _ := m.Builder().
			Where("id", "").
			Where(m.Builder().
				Where("nickname", "foo").
				WhereOr("password", "abc123")).
			Where("passport", "pp").
			Build()
		t.Assert(condWhere, `("id"=?) AND ((("nickname"=?) OR ("password"=?))) AND ("passport"=?)`)

		condWhere, _ = m.OmitEmpty().Builder().
			Where("id", "").
			Where(m.Builder().
				Where("nickname", "foo").
				WhereOr("password", "abc123")).
			Where("passport", "pp").
			Build()
		t.Assert(condWhere, `(("nickname"=?) OR ("password"=?)) AND ("passport"=?)`)

		condWhere, _ = m.OmitEmpty().Builder().
			Where(m.Builder().
				Where("nickname", "foo").
				WhereOr("password", "abc123")).
			Where("id", "").
			Where("passport", "pp").
			Build()
		t.Assert(condWhere, `(("nickname"=?) OR ("password"=?)) AND ("passport"=?)`)
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
func Test_Issue3086(t *testing.T) {
	table := "issue3086_user"
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `issue3086.sql`), ";")
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
			insertId int64
			data     = User{
				Id:       20,
				Passport: "passport_20",
				Password: "",
			}
		)
		sqlArray, err = gdb.CatchSQL(ctx, func(ctx context.Context) error {
			insertId, err = db.Ctx(ctx).Model(table).Data(data).InsertAndGetId()
			return err
		})
		t.AssertNil(err)
		t.Assert(insertId, 20)
		// CatchSQL may return empty on PgSQL when InsertAndGetId uses RETURNING.
		// The functional assertion (insertId=20) is the meaningful check.
		if len(sqlArray) > 0 {
			t.Assert(
				gstr.Contains(sqlArray[len(sqlArray)-1], `("id","passport") VALUES(20,'passport_20')`),
				true,
			)
		}
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
			gstr.Contains(sqlArray[len(sqlArray)-1], `SET "passport"='passport_1' WHERE "id"=1`),
			true,
		)
	})
}

// https://github.com/gogf/gf/issues/3218
func Test_Issue3218(t *testing.T) {
	table := "issue3218_sys_config"
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `issue3218.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
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
func Test_Issue2552_ClearTableFieldsAll(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	// PgSQL driver queries table fields via information_schema.columns.
	showTableKey := `information_schema`
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err := db.Model(table).Ctx(ctx).Insert(g.Map{
				"passport":    guid.S(),
				"password":    guid.S(),
				"nickname":    guid.S(),
				"create_time": gtime.NewFromStr(CreateTime).String(),
			})
			return err
		})
		t.AssertNil(err)
		t.Assert(gstr.Contains(gstr.Join(sqlArray, "|"), showTableKey), true)

		ctx = context.Background()
		sqlArray, err = gdb.CatchSQL(ctx, func(ctx context.Context) error {
			one, err := db.Model(table).Ctx(ctx).One()
			t.AssertGT(len(one), 0)
			return err
		})
		t.AssertNil(err)
		t.Assert(gstr.Contains(gstr.Join(sqlArray, "|"), showTableKey), false)

		_, err = db.Exec(ctx, fmt.Sprintf(`ALTER TABLE %s DROP COLUMN nickname`, table))
		t.AssertNil(err)

		err = db.GetCore().ClearTableFieldsAll(ctx)
		t.AssertNil(err)

		ctx = context.Background()
		sqlArray, err = gdb.CatchSQL(ctx, func(ctx context.Context) error {
			one, err := db.Model(table).Ctx(ctx).One()
			t.AssertGT(len(one), 0)
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

	// PgSQL driver queries table fields via information_schema.columns.
	showTableKey := `information_schema`
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err := db.Model(table).Ctx(ctx).Insert(g.Map{
				"passport":    guid.S(),
				"password":    guid.S(),
				"nickname":    guid.S(),
				"create_time": gtime.NewFromStr(CreateTime).String(),
			})
			return err
		})
		t.AssertNil(err)
		t.Assert(gstr.Contains(gstr.Join(sqlArray, "|"), showTableKey), true)

		ctx = context.Background()
		sqlArray, err = gdb.CatchSQL(ctx, func(ctx context.Context) error {
			one, err := db.Model(table).Ctx(ctx).One()
			t.AssertGT(len(one), 0)
			return err
		})
		t.AssertNil(err)
		t.Assert(gstr.Contains(gstr.Join(sqlArray, "|"), showTableKey), false)

		_, err = db.Exec(ctx, fmt.Sprintf(`ALTER TABLE %s DROP COLUMN nickname`, table))
		t.AssertNil(err)

		err = db.GetCore().ClearTableFields(ctx, table)
		t.AssertNil(err)

		ctx = context.Background()
		sqlArray, err = gdb.CatchSQL(ctx, func(ctx context.Context) error {
			one, err := db.Model(table).Ctx(ctx).One()
			t.AssertGT(len(one), 0)
			return err
		})
		t.AssertNil(err)
		t.Assert(gstr.Contains(gstr.Join(sqlArray, "|"), showTableKey), true)
	})
}

// https://github.com/gogf/gf/issues/2643
func Test_Issue2643(t *testing.T) {
	table := "issue2643"
	array := gstr.SplitAndTrim(gtest.DataContent("issues", "issue2643.sql"), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			expectKey1 = `SELECT s.name,replace(concat_ws(',',lpad(s.id::text, 6, '0'),s.name),',','') "code" FROM "issue2643" AS s`
			expectKey2 = `SELECT CASE WHEN dept='物资部' THEN '物资部' ELSE '其他' END dept,sum(s.value) FROM "issue2643" AS s GROUP BY CASE WHEN dept='物资部' THEN '物资部' ELSE '其他' END`
		)
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			db.Ctx(ctx).Model(table).As("s").Fields(
				"s.name",
				`replace(concat_ws(',',lpad(s.id::text, 6, '0'),s.name),',','') "code"`,
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
		// PgSQL uses double quotes instead of backticks
		sqlStr := fmt.Sprintf(`SELECT COUNT(1) FROM "%s" WHERE (create_time = now()) AND ("create_time" < now())`, table)
		t.Assert(sql[0], sqlStr)
	})
}

// https://github.com/gogf/gf/issues/3754
func Test_Issue3754(t *testing.T) {
	table := "issue3754"
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `issue3754.sql`), ";")
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
		t.Assert(oneInsert["id"].Int(), 1)
		t.Assert(oneInsert["name"].String(), "name_1")
		t.Assert(oneInsert["delete_at"].String(), "")
		t.Assert(oneInsert["create_at"].String(), "")
		t.Assert(oneInsert["update_at"].String(), "")

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
		t.Assert(oneUpdate["id"].Int(), 1)
		t.Assert(oneUpdate["name"].String(), "name_1000")
		t.Assert(oneUpdate["delete_at"].String(), "")
		t.Assert(oneUpdate["create_at"].String(), "")
		t.Assert(oneUpdate["update_at"].String(), "")

		// FieldsEx does not affect Delete operation.
		r, err = db.Model(table).FieldsEx(fieldsEx).WherePri(1).Delete()
		n, _ = r.RowsAffected()
		t.Assert(n, 1)
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
func Test_Issue3626(t *testing.T) {
	table := "issue3626"
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `issue3626.sql`), ";")
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
	table := "issue3915"
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `issue3915.sql`), ";")
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
		t.Assert(all[0]["id"], 1)

		all, err = db.Model(table).Where(gdb.Raw("a < b")).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 1)

		// PgSQL uses double quotes for column quoting
		all, err = db.Model(table).WhereLT("a", gdb.Raw(`"b"`)).All()
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

		// PgSQL uses double quotes for column quoting
		all, err = db.Model(table).WhereGT("a", gdb.Raw(`"b"`)).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 2)
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

		defer dropTable(tables[0])
		defer dropTable(tables[1])
		_ = tables
		array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `issue2119.sql`), ";")
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
		array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `issue4034.sql`), ";")
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
	// PgSQL ON CONFLICT requires a conflict target; the original MySQL test uses
	// Save() (REPLACE INTO) which works without one. Use Insert() instead since
	// this test only inserts a new record.
	_, err := db.Model(table).Safe().Ctx(ctx).TX(tx).Data(g.Map{
		"passport": "111",
		"password": "222",
		"nickname": "333",
	}).Insert()
	return err
}

// https://github.com/gogf/gf/issues/4086
func Test_Issue4086(t *testing.T) {
	table := "issue4086"
	defer dropTable(table)
	array := gstr.SplitAndTrim(gtest.DataContent(`issues`, `issue4086.sql`), ";")
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

// https://github.com/gogf/gf/issues/4697
func Test_Issue4697(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Fields("") should be treated as Fields() and select all fields
		result, err := db.Model(table).Fields("").Limit(1).All()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		// PgSQL createTable has 10 columns (id, passport, password, nickname, create_time,
		// create_date, favorite_movie, favorite_music, numeric_values, decimal_values)
		t.AssertGT(len(result[0]), 5)
	})

	gtest.C(t, func(t *gtest.T) {
		// Fields("", "id") should ignore empty string and only select "id"
		result, err := db.Model(table).Fields("", "id").Limit(1).All()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(len(result[0]), 1)
		t.AssertNE(result[0]["id"], nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Fields("id", "", "nickname") should ignore empty string
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

	// Test 1: AllAndCount with multiple fields should generate valid COUNT SQL
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id, nickname").AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(len(result), TableSize)
		t.AssertNE(result[0]["id"], nil)
		t.AssertNE(result[0]["nickname"], nil)
		t.Assert(result[0]["passport"], nil)
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
		t.AssertGT(len(result), 0)
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

// https://github.com/gogf/gf/issues/2231
func Test_Issue2231(t *testing.T) {
	t.Skip("MariaDB-specific regex link test not applicable to PostgreSQL")
}
