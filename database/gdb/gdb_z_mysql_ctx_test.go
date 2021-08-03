// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/test/gtest"
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
	db.GetLogger().(gdb.LoggerImp).SetCtxKeys("SpanId", "TraceId")
	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)
		ctx := context.WithValue(context.Background(), "TraceId", "12345678")
		ctx = context.WithValue(ctx, "SpanId", "0.1")
		db.Ctx(ctx).Query("select 1")
	})
	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)
		db.Query("select 2")
	})
}

func Test_Ctx_Model(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	db.GetLogger().(gdb.LoggerImp).SetCtxKeys("SpanId", "TraceId")
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

func Test_Ctx_Strict(t *testing.T) {
	table := createInitTableWithDb(dbCtxStrict)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := dbCtxStrict.Query("select 1")
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		r, err := dbCtxStrict.Model(table).All()
		t.AssertNE(err, nil)
		t.Assert(len(r), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		r, err := dbCtxStrict.Model(table).Ctx(context.TODO()).All()
		t.AssertNil(err)
		t.Assert(len(r), TableSize)
	})
}
