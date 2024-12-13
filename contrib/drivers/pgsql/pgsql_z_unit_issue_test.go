// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
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
			list []map[string]interface{}
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
		Text   interface{}
		Number interface{}
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
