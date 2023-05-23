// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

var (
	ctx = context.TODO()
)

func Test_Try(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := `gutil Try test`
		t.Assert(gutil.Try(ctx, func(ctx context.Context) {
			panic(s)
		}), s)
	})
	gtest.C(t, func(t *gtest.T) {
		s := `gutil Try test`
		t.Assert(gutil.Try(ctx, func(ctx context.Context) {
			panic(gerror.New(s))
		}), s)
	})
}

func Test_TryCatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gutil.TryCatch(ctx, func(ctx context.Context) {
			panic("gutil TryCatch test")
		})
	})

	gtest.C(t, func(t *gtest.T) {
		gutil.TryCatch(ctx, func(ctx context.Context) {
			panic("gutil TryCatch test")

		}, func(ctx context.Context, err error) {
			t.Assert(err, "gutil TryCatch test")
		})
	})

	gtest.C(t, func(t *gtest.T) {
		gutil.TryCatch(ctx, func(ctx context.Context) {
			panic(gerror.New("gutil TryCatch test"))

		}, func(ctx context.Context, err error) {
			t.Assert(err, "gutil TryCatch test")
		})
	})
}

func Test_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.IsEmpty(1), false)
	})
}

func Test_Throw(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			t.Assert(recover(), "gutil Throw test")
		}()

		gutil.Throw("gutil Throw test")
	})
}

func Test_Keys(t *testing.T) {
	// not support int
	gtest.C(t, func(t *gtest.T) {
		var val int = 1
		keys := gutil.Keys(reflect.ValueOf(val))
		t.AssertEQ(len(keys), 0)
	})
	// map
	gtest.C(t, func(t *gtest.T) {
		keys := gutil.Keys(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)

		strKeys := gutil.Keys(map[string]interface{}{
			"key1": 1,
			"key2": 2,
		})
		t.AssertIN("key1", strKeys)
		t.AssertIN("key2", strKeys)
	})
	// *map
	gtest.C(t, func(t *gtest.T) {
		keys := gutil.Keys(&map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)
	})
	// *struct
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		keys := gutil.Keys(new(T))
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
	// *struct nil
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := gutil.Keys(pointer)
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
	// **struct nil
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := gutil.Keys(&pointer)
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
}

func Test_Values(t *testing.T) {
	// not support int
	gtest.C(t, func(t *gtest.T) {
		var val int = 1
		keys := gutil.Values(reflect.ValueOf(val))
		t.AssertEQ(len(keys), 0)
	})
	// map
	gtest.C(t, func(t *gtest.T) {
		values := gutil.Values(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN(10, values)
		t.AssertIN(20, values)

		values = gutil.Values(map[string]interface{}{
			"key1": 10,
			"key2": 20,
		})
		t.AssertIN(10, values)
		t.AssertIN(20, values)
	})
	// *map
	gtest.C(t, func(t *gtest.T) {
		keys := gutil.Values(&map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN(10, keys)
		t.AssertIN(20, keys)
	})
	// struct
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		keys := gutil.Values(T{
			A: "1",
			B: 2,
		})
		t.Assert(keys, g.Slice{"1", 2})
	})
}

func TestListToMapByKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		listMap := []map[string]interface{}{
			{"key1": 1, "key2": 2},
			{"key3": 3, "key4": 4},
		}
		t.Assert(gutil.ListToMapByKey(listMap, "key1"), "{\"1\":{\"key1\":1,\"key2\":2}}")
	})
}

func Test_GetOrDefaultStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.GetOrDefaultStr("a", "b"), "b")
		t.Assert(gutil.GetOrDefaultStr("a", "b", "c"), "b")
		t.Assert(gutil.GetOrDefaultStr("a"), "a")
	})
}

func Test_GetOrDefaultAny(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.GetOrDefaultAny("a", "b"), "b")
		t.Assert(gutil.GetOrDefaultAny("a", "b", "c"), "b")
		t.Assert(gutil.GetOrDefaultAny("a"), "a")
	})
}
