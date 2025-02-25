// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Ctx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db, err := gdb.Instance()
		t.AssertNil(err)

		err1 := db.PingMaster()
		err2 := db.PingSlave()
		t.Assert(err1, nil)
		t.Assert(err2, nil)

		newDb := db.Ctx(context.Background())
		t.AssertNE(newDb, nil)
	})
}

func Test_Ctx_Query(t *testing.T) {
	db.GetLogger().(*glog.Logger).SetCtxKeys("SpanId", "TraceId")
	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)
		ctx := context.WithValue(context.Background(), "TraceId", "12345678")
		ctx = context.WithValue(ctx, "SpanId", "0.1")
		db.Query(ctx, "select 1")
	})
	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)
		db.Query(ctx, "select 2")
	})
}

func Test_Ctx_Model(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	db.GetLogger().(*glog.Logger).SetCtxKeys("SpanId", "TraceId")
	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)
		ctx := context.WithValue(context.Background(), "TraceId", "12345678")
		ctx = context.WithValue(ctx, "SpanId", "0.1")
		db.Model(table).Ctx(ctx).All()
	})
	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)
		db.Model(table).All()
	})
}

type implScanner struct {
	v int
}

func (t *implScanner) Scan(src any) error {
	t.v = 10
	return nil
}

type testImplSqlScanner struct {
	Id       implScanner
	Passport *implScanner
	Password string
}

func Test_CustomConvertFunc(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		var impl testImplSqlScanner
		err := db.Model(table).Scan(&impl)
		t.AssertNil(err)
		t.Assert(impl.Id, implScanner{v: 10})
		t.Assert(impl.Password, "pass_1")
		t.Assert(impl.Passport, &implScanner{v: 10})
	})
}
