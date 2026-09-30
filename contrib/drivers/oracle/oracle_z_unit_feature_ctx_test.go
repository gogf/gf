// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"context"
	"testing"
	"time"

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
		_, err := db.Query(ctx, "select 1 from dual")
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)
		_, err := db.Query(ctx, "select 2 from dual")
		t.AssertNil(err)
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
		all, err := db.Model(table).Ctx(ctx).All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
	})
	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)
		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
	})
}

// Test_Ctx_Timeout tests context timeout behavior
func Test_Ctx_Timeout(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(1 * time.Millisecond)

		_, err := db.Model(table).Ctx(ctx).All()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		// Note: on Windows the statement runs to its end before ORA-01013 is returned, as go-ora v2.9.0
		// cannot interrupt it there (sendOOB in network/net_windows.go does nothing).
		_, err := db.Exec(ctx, "BEGIN DBMS_LOCK.SLEEP(2); END;")
		t.AssertNE(err, nil)
		t.AssertIN("ORA-01013", err.Error())

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

// Test_Ctx_Cancel tests context cancellation
func Test_Ctx_Cancel(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := db.Model(table).Ctx(ctx).All()
		t.AssertNE(err, nil)
	})
}

// Test_Ctx_Propagation_Transaction tests context propagation in transaction
func Test_Ctx_Propagation_Transaction(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	db.GetLogger().(*glog.Logger).SetCtxKeys("TraceId")

	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)

		ctx := context.WithValue(context.Background(), "TraceId", "tx_trace_123")
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).Ctx(ctx).Where("id", 1).One()
			return err
		})
		t.AssertNil(err)
	})
}

// Test_Ctx_Multiple_Values tests context with multiple values
func Test_Ctx_Multiple_Values(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	db.GetLogger().(*glog.Logger).SetCtxKeys("TraceId", "RequestId", "UserId")

	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)

		ctx := context.WithValue(context.Background(), "TraceId", "trace_001")
		ctx = context.WithValue(ctx, "RequestId", "req_002")
		ctx = context.WithValue(ctx, "UserId", "user_003")

		one, err := db.Model(table).Ctx(ctx).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)
	})
}

// Test_Ctx_Nested_Operations tests context in nested operations
func Test_Ctx_Nested_Operations(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	db.GetLogger().(*glog.Logger).SetCtxKeys("TraceId")

	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)

		ctx := context.WithValue(context.Background(), "TraceId", "nested_trace")

		result, err := db.Model(table).Ctx(ctx).Where("id>", 0).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)

		one, err := db.Model(table).Ctx(ctx).Where("id", result[0]["ID"]).One()
		t.AssertNil(err)
		t.Assert(one["ID"], result[0]["ID"])
	})
}
