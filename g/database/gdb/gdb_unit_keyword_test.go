// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/test/gtest"
)

const (
	TABLE_KEYWORD        = "values"
	TABLE_KEYWORD_QUOTED = "`values`"
	KEY_FOR_TEST         = "KeyForTest"
	VAL_FOR_TEST         = "ValForTest"
)

type kv struct {
	Id  int    `gconv:"id"`
	Key string `json:"key"`
	Val string `gconv:"val"`
}

func initKeywordTable() {
	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `values`;\n" +
		"CREATE TABLE `values` (" +
		"	`id` int(11) NOT NULL AUTO_INCREMENT," +
		"	`key` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL," +
		"	`val` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL," +
		"	PRIMARY KEY (`id`)," +
		"	UNIQUE KEY `uix_values_key` (`key`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;\n" +
		"INSERT INTO `values` (`id`, `key`, `val`) VALUES" +
		"(1,	'UserHid',	'1000007')," +
		"(2,	'KeyForTest',	'ValForTest');")); err != nil {
		gtest.Error(err)
	}
	return
}

func Test_Keyword_Insert(t *testing.T) {
	doTest_Keyword_Insert(t, TABLE_KEYWORD)
	doTest_Keyword_Insert(t, TABLE_KEYWORD_QUOTED)
}

func doTest_Keyword_Insert(t *testing.T, table string) {
	initKeywordTable()
	gtest.Case(t, func() {
		result, err := db.Table(table).Filter().Data(g.Map{
			"id":  10,
			"key": "k10",
			"val": "v10",
		}).Insert()
		gtest.Assert(err, nil)
		n, _ := result.LastInsertId()
		gtest.Assert(n, 10)

		result, err = db.Table(table).Filter().Data(g.Map{
			"id":  "11",
			"key": "k11",
			"val": "v11",
		}).Insert()
		gtest.Assert(err, nil)
		n, _ = result.RowsAffected()
		gtest.Assert(n, 1)

		result, err = db.Table(table).Filter().Data(kv{
			Id:  12,
			Key: "k12",
			Val: "v12",
		}).Insert()
		gtest.Assert(err, nil)
		n, _ = result.RowsAffected()
		gtest.Assert(n, 1)
		value, err := db.Table(table).Fields("val").Where("id=12").Value()
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "v12")

		result, err = db.Table(table).Filter().Data(&kv{
			Id:  13,
			Key: "k13",
			Val: "v13",
		}).Insert()
		gtest.Assert(err, nil)
		n, _ = result.RowsAffected()
		gtest.Assert(n, 1)
		value, err = db.Table(table).Fields("val").Where("id=13").Value()
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "v13")

		result, err = db.Table(table).Where("id>?", 10).Delete()
		gtest.Assert(err, nil)
		n, _ = result.RowsAffected()
		gtest.Assert(n, 3)
	})
}

func Test_Keyword_Batch(t *testing.T) {
	doTest_Keyword_Batch(t, TABLE_KEYWORD)
	doTest_Keyword_Batch(t, TABLE_KEYWORD_QUOTED)
}

func doTest_Keyword_Batch(t *testing.T, table string) {
	// bacth insert
	gtest.Case(t, func() {
		initKeywordTable()
		result, err := db.Table(table).Filter().Data(g.List{
			{
				"id":  10,
				"key": "k10",
				"val": "v10",
			},
			{
				"id":  "11",
				"key": "k11",
				"val": "v11",
			},
		}).Batch(1).Insert()
		if err != nil {
			gtest.Error(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 2)
	})

	// batch save
	gtest.Case(t, func() {
		initKeywordTable()
		result, err := db.Table(table).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 2)
		for _, v := range result {
			v["val"].Set(v["val"].String() + v["id"].String())
		}
		r, e := db.Table(table).Data(result).Save()
		gtest.Assert(e, nil)
		n, e := r.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, 4)
	})

	// batch replace
	gtest.Case(t, func() {
		initKeywordTable()
		result, err := db.Table(table).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 2)
		for _, v := range result {
			v["val"].Set(v["val"].String() + v["id"].String())
		}
		r, e := db.Table(table).Data(result).Replace()
		gtest.Assert(e, nil)
		n, e := r.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, 4)
	})
}

func Test_Keyword_Update(t *testing.T) {
	initKeywordTable()
	table := TABLE_KEYWORD
	// UPDATE...LIMIT
	gtest.Case(t, func() {
		result, err := db.Table(table).Data("val", "T100").OrderBy("id desc").Limit(1).Update()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)

		v1, err := db.Table(table).Fields("val").Where("id", 2).Value()
		gtest.Assert(err, nil)
		gtest.Assert(v1.String(), "T100")

		v2, err := db.Table(table).Fields("val").Where("id", 1).Value()
		gtest.Assert(err, nil)
		gtest.AssertNE(v2.String(), "T100")
	})

	gtest.Case(t, func() {
		result, err := db.Table(table).Data("val", "T200").Where("val=?", "T100").Update()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	})

	gtest.Case(t, func() {
		cond := fmt.Sprintf("`key`='%s'", KEY_FOR_TEST)
		result, err := db.Table(table).Data("val", "T100").Where(cond).Update()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	})
}

func Test_Keyword_Clone(t *testing.T) {
	initKeywordTable()
	table := TABLE_KEYWORD

	gtest.Case(t, func() {
		md := db.Table(table).Where("id IN(?)", g.Slice{1, 2})
		count, err := md.Count()
		gtest.Assert(err, nil)

		record, err := md.OrderBy("id DESC").One()
		gtest.Assert(err, nil)

		result, err := md.OrderBy("id ASC").All()
		gtest.Assert(err, nil)

		gtest.Assert(count, 2)
		gtest.Assert(record["id"].Int(), 2)
		gtest.Assert(len(result), 2)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 2)
	})
}

func Test_Keyword_Safe(t *testing.T) {
	initKeywordTable()
	table := TABLE_KEYWORD

	gtest.Case(t, func() {
		md := db.Table(table).Safe(false).Where("id IN(?)", g.Slice{1, 2})
		count, err := md.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 2)

		md.And("id = ?", 1)
		count, err = md.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 1)
	})

	gtest.Case(t, func() {
		md := db.Table(table).Safe(true).Where("id IN(?)", g.Slice{1, 2})
		count, err := md.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 2)

		md.And("id = ?", 1)
		count, err = md.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 2)
	})
}

func Test_Keyword_All(t *testing.T) {
	initKeywordTable()
	table := TABLE_KEYWORD

	gtest.Case(t, func() {
		result, err := db.Table(table).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 2)
	})
	// sql.ErrNoRows
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id<0").All()
		gtest.Assert(result, nil)
		gtest.Assert(err, sql.ErrNoRows)
	})
}

func Test_Keyword_One(t *testing.T) {
	doTest_Keyword_One(t, TABLE_KEYWORD)
	doTest_Keyword_One(t, TABLE_KEYWORD_QUOTED)
}

func doTest_Keyword_One(t *testing.T, table string) {
	initKeywordTable()
	gtest.Case(t, func() {
		_, err := db.Table(table).Where("id", 1).One()
		gtest.Assert(err, nil)
	})

	gtest.Case(t, func() {
		record, err := db.Table(table).Where("key", KEY_FOR_TEST).One()
		gtest.Assert(err, nil)
		gtest.Assert(record["val"].String(), VAL_FOR_TEST)
	})

	gtest.Case(t, func() {
		record, err := db.Table(table).Where("`key`", KEY_FOR_TEST).One()
		gtest.Assert(err, nil)
		gtest.Assert(record["val"].String(), VAL_FOR_TEST)
	})
}

func Test_Keyword_Value(t *testing.T) {
	initKeywordTable()
	table := TABLE_KEYWORD
	gtest.Case(t, func() {
		value, err := db.Table(table).Fields("val").Where("id", 2).Value()
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), VAL_FOR_TEST)
	})
	// sql.ErrNoRows
	gtest.Case(t, func() {
		value, err := db.Table(table).Fields("val").Where("id", 0).Value()
		gtest.Assert(err, sql.ErrNoRows)
		gtest.Assert(value, nil)
	})
}

func Test_Keyword_Count(t *testing.T) {
	initKeywordTable()
	table := TABLE_KEYWORD
	gtest.Case(t, func() {
		count, err := db.Table(table).Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 2)
	})
}

func Test_Keyword_Select(t *testing.T) {
	initKeywordTable()
	table := TABLE_KEYWORD
	gtest.Case(t, func() {
		result, err := db.Table(table).Select()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 2)
	})
}

func Test_Keyword_Struct(t *testing.T) {
	initKeywordTable()
	table := TABLE_KEYWORD
	gtest.Case(t, func() {
		rec := new(kv)
		err := db.Table(table).Where("id=2").Struct(rec)
		gtest.Assert(err, nil)
		gtest.Assert(rec.Key, KEY_FOR_TEST)
		gtest.Assert(rec.Val, VAL_FOR_TEST)
	})
	gtest.Case(t, func() {
		rec := new(kv)
		err := db.Table(table).Where("id=2").Struct(rec)
		gtest.Assert(err, nil)
		gtest.Assert(rec.Key, KEY_FOR_TEST)
		gtest.Assert(rec.Val, VAL_FOR_TEST)
	})
	// Auto creating struct object.
	gtest.Case(t, func() {
		rec := (*kv)(nil)
		err := db.Table(table).Where("id=2").Struct(&rec)
		gtest.Assert(err, nil)
		gtest.Assert(rec.Key, KEY_FOR_TEST)
		gtest.Assert(rec.Val, VAL_FOR_TEST)
	})
	// Just using Scan.
	gtest.Case(t, func() {
		rec := (*kv)(nil)
		err := db.Table(table).Where("id=2").Scan(&rec)
		if err != nil {
			gtest.Error(err)
		}
		gtest.Assert(rec.Key, KEY_FOR_TEST)
		gtest.Assert(rec.Val, VAL_FOR_TEST)
	})
	// sql.ErrNoRows
	gtest.Case(t, func() {
		rec := new(kv)
		err := db.Table(table).Where("id=-1").Struct(rec)
		gtest.Assert(err, sql.ErrNoRows)
	})
}

func Test_Keyword_Structs(t *testing.T) {
	initKeywordTable()
	table := TABLE_KEYWORD

	gtest.Case(t, func() {
		var kvs []kv
		err := db.Table(table).OrderBy("id asc").Structs(&kvs)
		if err != nil {
			gtest.Error(err)
		}
		gtest.Assert(len(kvs), 2)
		gtest.Assert(kvs[0].Id, 1)
		gtest.Assert(kvs[1].Id, 2)
		gtest.Assert(kvs[1].Key, KEY_FOR_TEST)
		gtest.Assert(kvs[1].Val, VAL_FOR_TEST)
	})
	// Auto create struct slice.
	gtest.Case(t, func() {
		var kvs []*kv
		err := db.Table(table).OrderBy("id asc").Structs(&kvs)
		if err != nil {
			gtest.Error(err)
		}
		gtest.Assert(len(kvs), 2)
		gtest.Assert(kvs[0].Id, 1)
		gtest.Assert(kvs[1].Id, 2)
		gtest.Assert(kvs[1].Key, KEY_FOR_TEST)
		gtest.Assert(kvs[1].Val, VAL_FOR_TEST)
	})
	// Just using Scan.
	gtest.Case(t, func() {
		var kvs []*kv
		err := db.Table(table).OrderBy("id asc").Scan(&kvs)
		if err != nil {
			gtest.Error(err)
		}
		gtest.Assert(len(kvs), 2)
		gtest.Assert(kvs[0].Id, 1)
		gtest.Assert(kvs[1].Id, 2)
		gtest.Assert(kvs[1].Key, KEY_FOR_TEST)
		gtest.Assert(kvs[1].Val, VAL_FOR_TEST)
	})
	// sql.ErrNoRows
	gtest.Case(t, func() {
		var kvs []*kv
		err := db.Table(table).Where("id<0").Structs(&kvs)
		gtest.Assert(err, sql.ErrNoRows)
	})
}

func Test_Keyword_Where(t *testing.T) {
	initKeywordTable()
	table := TABLE_KEYWORD

	// string
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id=? and val=?", 2, VAL_FOR_TEST).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 2)
	})
	// slice parameter
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id=? and val=?", g.Slice{2, VAL_FOR_TEST}).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 2)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id=?", g.Slice{2}).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 2)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id", 2).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 2)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id", 2).Where("key", KEY_FOR_TEST).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id", 2).And("key", KEY_FOR_TEST).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id", 30).Or("key", KEY_FOR_TEST).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("(id=30)").Or("key", KEY_FOR_TEST).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})
	// slice
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("`key` like ? and val like ?", g.Slice{KEY_FOR_TEST, VAL_FOR_TEST}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})
	// map
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(g.Map{"id": 2, "key": KEY_FOR_TEST}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})
	// struct
	gtest.Case(t, func() {
		type Key struct {
			Id  int    `json:"id"`
			Key string `gconv:"key"`
		}
		result, err := db.Table(table).Where(Key{2, KEY_FOR_TEST}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)

		result, err = db.Table(table).Where(&Key{2, KEY_FOR_TEST}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})
	// slice single
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id IN(?)", g.Slice{1, 2}).OrderBy("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 2)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 2)
	})
	// slice + string
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("`key`=? AND id IN(?)", KEY_FOR_TEST, g.Slice{1, 2}).OrderBy("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 2)
	})
	// slice + map
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(g.Map{
			"id":  g.Slice{1, 2},
			"key": KEY_FOR_TEST,
		}).OrderBy("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 2)
	})
	// slice + struct
	gtest.Case(t, func() {
		type User struct {
			Ids []int  `json:"id"`
			Key string `gconv:"key"`
		}
		result, err := db.Table(table).Where(User{
			Ids: []int{1, 2},
			Key: KEY_FOR_TEST,
		}).OrderBy("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 2)
	})
}

func Test_Keyword_Delete(t *testing.T) {
	initKeywordTable()
	table := TABLE_KEYWORD

	// DELETE...LIMIT
	gtest.Case(t, func() {
		result, err := db.Table(table).Limit(2).Delete()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 2)
	})

	initKeywordTable()
	gtest.Case(t, func() {
		result, err := db.Table(table).Delete()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 2)
	})
}
