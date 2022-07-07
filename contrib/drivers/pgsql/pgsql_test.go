// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
)

func init() {
	fmt.Println("init pgsql db")
	myDB, err := gdb.New(gdb.ConfigNode{
		Type: "pgsql",
		// Link: "pgsql://postgres:123456@127.0.0.1:5432/user?sslmode=disable",
		User:  "postgres",
		Pass:  "123456",
		Host:  "127.0.0.1",
		Port:  "5432",
		Name:  "postgres",
		Debug: true,
	})
	if err != nil {
		glog.Fatal(context.TODO(), err.Error())
		return
	}
	list, err := myDB.Tables(context.TODO())
	if err != nil {
		glog.Fatal(context.TODO(), err.Error())
		return
	}
	fmt.Println(list)
	// fmt.Println(myDB.Model("user").InsertAndGetId(g.Map{"name": "hailaz"}))

	res, err := myDB.Model("user").Insert(g.List{
		{"name": "john_1"},
		{"name": "john_2"},
		{"name": "john_3"},
	})
	fmt.Println(err)
	fmt.Print("LastInsertId: ")
	fmt.Println(res.LastInsertId())
	fmt.Print("RowsAffected: ")
	fmt.Println(res.RowsAffected())

}

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
