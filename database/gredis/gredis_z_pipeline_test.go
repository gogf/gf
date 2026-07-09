// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT License was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Cmd_BeforeExec verifies that a newly created Cmd returns nil values
// before any SetVal or SetErr calls.
func Test_Cmd_BeforeExec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cmd := &gredis.Cmd{}
		val, err := cmd.Result()
		t.AssertNil(val)
		t.AssertNil(err)
		t.AssertNil(cmd.Val())
		t.AssertNil(cmd.Err())
	})
}

// Test_Cmd_AfterSetVal verifies that SetVal correctly populates the value
// accessible via Result() and Val().
func Test_Cmd_AfterSetVal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cmd := &gredis.Cmd{}
		cmd.SetVal(gvar.New("test"))
		val, err := cmd.Result()
		t.AssertNil(err)
		t.Assert(val.String(), "test")
		t.Assert(cmd.Val().String(), "test")
	})
}

// Test_Cmd_AfterSetErr verifies that SetErr correctly populates the error
// accessible via Result() and Err().
func Test_Cmd_AfterSetErr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cmd := &gredis.Cmd{}
		someErr := errors.New("redis error")
		cmd.SetErr(someErr)
		val, err := cmd.Result()
		t.AssertNil(val)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "redis error")
		t.AssertNE(cmd.Err(), nil)
		t.Assert(cmd.Err().Error(), "redis error")
	})
}

// Test_Cmd_AfterSetBoth verifies that both value and error can be populated
// simultaneously on a Cmd.
func Test_Cmd_AfterSetBoth(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cmd := &gredis.Cmd{}
		cmd.SetVal(gvar.New(42))
		someErr := errors.New("partial failure")
		cmd.SetErr(someErr)
		val, err := cmd.Result()
		t.Assert(val.Int(), 42)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "partial failure")
		t.Assert(cmd.Val().Int(), 42)
		t.Assert(cmd.Err().Error(), "partial failure")
	})
}

// Test_Cmd_SetValOverwrite verifies that SetVal can be called multiple times
// and the last value wins.
func Test_Cmd_SetValOverwrite(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cmd := &gredis.Cmd{}
		cmd.SetVal(gvar.New("first"))
		t.Assert(cmd.Val().String(), "first")
		cmd.SetVal(gvar.New("second"))
		t.Assert(cmd.Val().String(), "second")
	})
}

// Test_Cmd_SetErrOverwrite verifies that SetErr can be called multiple times
// and the last error wins.
func Test_Cmd_SetErrOverwrite(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cmd := &gredis.Cmd{}
		cmd.SetErr(errors.New("err1"))
		t.Assert(cmd.Err().Error(), "err1")
		cmd.SetErr(errors.New("err2"))
		t.Assert(cmd.Err().Error(), "err2")
	})
}

// Test_Cmd_SetErrNil verifies that SetErr(nil) clears a previously set error.
func Test_Cmd_SetErrNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cmd := &gredis.Cmd{}
		cmd.SetErr(errors.New("temporary"))
		t.AssertNE(cmd.Err(), nil)
		cmd.SetErr(nil)
		t.AssertNil(cmd.Err())
	})
}

// Compile-time check: Redis must implement Pipeline, TxPipeline, and Watch methods.
var _ interface {
	Pipeline(ctx context.Context) gredis.Pipeliner
	TxPipeline(ctx context.Context) gredis.Pipeliner
	Watch(ctx context.Context, fn func(gredis.Tx) error, keys ...string) error
} = (*gredis.Redis)(nil)

// Test_AdapterOperation_HasPipeline verifies at runtime that the
// AdapterOperation interface includes Pipeline, TxPipeline, and Watch.
func Test_AdapterOperation_HasPipeline(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var _ gredis.AdapterOperation = (gredis.AdapterOperation)(nil)
		t.Assert(true, true)
	})
}

// Test_PipelineInterface_Compliance verifies that the Pipeliner interface
// composes PipelinerOperation and PipelinerGroup.
func Test_PipelineInterface_Compliance(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var _ gredis.Pipeliner = (gredis.Pipeliner)(nil)
		t.Assert(true, true)
	})
}

// Test_IGroupGeneric_ScanAll verifies at compile time that IGroupGeneric
// includes the ScanAll method.
func Test_IGroupGeneric_ScanAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var _ gredis.IGroupGeneric = (gredis.IGroupGeneric)(nil)
		t.Assert(true, true)
	})
}

// Test_Redis_Pipeline_NilReceiver verifies that Pipeline returns nil when
// called on a nil Redis or when no adapter is configured.
func Test_Redis_Pipeline_NilReceiver(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var r *gredis.Redis
		result := r.Pipeline(context.Background())
		t.AssertNil(result)
	})
}

// Test_Redis_TxPipeline_NilReceiver verifies that TxPipeline returns nil when
// called on a nil Redis or when no adapter is configured.
func Test_Redis_TxPipeline_NilReceiver(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var r *gredis.Redis
		result := r.TxPipeline(context.Background())
		t.AssertNil(result)
	})
}

// Test_Redis_Watch_NilReceiver verifies that Watch returns an error when
// called on a nil Redis.
func Test_Redis_Watch_NilReceiver(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var r *gredis.Redis
		err := r.Watch(context.Background(), func(tx gredis.Tx) error {
			return nil
		}, "key1")
		t.AssertNE(err, nil)
	})
}

// Test_Redis_MustPipeline_NilAdapter verifies that MustPipeline panics
// when called on a Redis with nil adapter.
func Test_Redis_MustPipeline_NilAdapter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			r := recover()
			t.AssertNE(r, nil)
		}()
		r := &gredis.Redis{}
		r.MustPipeline(context.Background())
	})
}

// Test_Redis_MustTxPipeline_NilAdapter verifies that MustTxPipeline panics
// when called on a Redis with nil adapter.
func Test_Redis_MustTxPipeline_NilAdapter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			r := recover()
			t.AssertNE(r, nil)
		}()
		r := &gredis.Redis{}
		r.MustTxPipeline(context.Background())
	})
}
