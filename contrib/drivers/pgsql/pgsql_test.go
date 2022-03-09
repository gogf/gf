// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"testing"

	"github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Driver_DoFilter(t *testing.T) {
	var (
		ctx    = gctx.New()
		driver = pgsql.Driver{}
	)
	gtest.C(t, func(t *gtest.T) {
		var data = g.Map{
			`select * from user where (role)::jsonb ?| 'admin'`: `select * from user where (role)::jsonb ?| 'admin'`,
			`select * from user where (role)::jsonb ?| '?'`:     `select * from user where (role)::jsonb ?| '$2'`,
			`select * from user where (role)::jsonb &? '?'`:     `select * from user where (role)::jsonb &? '$2'`,
			`select * from user where (role)::jsonb ? '?'`:      `select * from user where (role)::jsonb ? '$2'`,
			`select * from user where '?'`:                      `select * from user where '$1'`,
		}
		for k, v := range data {
			newSql, _, err := driver.DoFilter(ctx, nil, k, nil)
			t.AssertNil(err)
			t.Assert(newSql, v)
		}
	})
}
