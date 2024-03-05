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
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)

		t.Assert(fields["id"].Key, "pri")
		t.Assert(fields["password"].Key, "uni")

		var list []map[string]interface{}
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
			one, err := db.Model(table).WherePri(i).One()
			t.AssertNil(err)
			t.Assert(one["id"], list[i-1]["id"])
			t.Assert(one["passport"], list[i-1]["passport"])
			t.Assert(one["password"], list[i-1]["password"])
			t.Assert(one["nickname"], list[i-1]["nickname"])
		}
	})
}
