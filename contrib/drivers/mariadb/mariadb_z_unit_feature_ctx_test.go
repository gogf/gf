// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mariadb_test

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

// Test_Ctx_Timeout tests context timeout behavior
func Test_Ctx_Timeout(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Create a context with very short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		// Wait for timeout
		time.Sleep(1 * time.Millisecond)

		// Query should fail due to context timeout
		_, err := db.Model(table).Ctx(ctx).All()
		t.AssertNE(err, nil)
	})
}

// Test_Ctx_Cancel tests context cancellation
func Test_Ctx_Cancel(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithCancel(context.Background())
		// Cancel immediately
		cancel()

		// Query should fail due to cancelled context
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
			// Context should propagate to transaction operations
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

		db.Model(table).Ctx(ctx).Where("id", 1).One()
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

		// Nested query operations should all have context
		result, err := db.Model(table).Ctx(ctx).Where("id>", 0).All()
		t.AssertNil(err)

		if len(result) > 0 {
			// Another query using same context
			_, err = db.Model(table).Ctx(ctx).Where("id", result[0]["id"]).One()
			t.AssertNil(err)
		}
	})
}
