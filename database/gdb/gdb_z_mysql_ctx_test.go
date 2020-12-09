// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
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
		t.Assert(err, nil)

		err1 := db.PingMaster()
		err2 := db.PingSlave()
		t.Assert(err1, nil)
		t.Assert(err2, nil)

		newDb := db.Ctx(context.Background())
		t.AssertNE(newDb, nil)
	})
}

func Test_Ctx_Query(t *testing.T) {
	db.GetLogger().SetCtxKeys("SpanId", "TraceId")
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
	db.GetLogger().SetCtxKeys("SpanId", "TraceId")
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
