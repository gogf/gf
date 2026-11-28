// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/test/gtest"
)

// https://github.com/gogf/gf/issues/4145
func Test_Issue4145(t *testing.T) {
	// The function of GetOrSetFuncLock can read and fill other keys of the same cache.
	gtest.C(t, func(t *gtest.T) {
		var (
			cache = gcache.New()
			done  = make(chan struct{})
			v     *gvar.Var
			err   error
		)
		go func() {
			defer close(done)
			v, err = cache.GetOrSetFuncLock(ctx, "outer", func(ctx context.Context) (any, error) {
				if _, err := cache.Get(ctx, "missing"); err != nil {
					return nil, err
				}
				inner, err := cache.GetOrSetFuncLock(ctx, "inner", func(ctx context.Context) (any, error) {
					return 1, nil
				}, 0)
				if err != nil {
					return nil, err
				}
				return inner.Int() + 1, nil
			}, 0)
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("nested GetOrSetFuncLock call for another key did not return")
		}
		t.AssertNil(err)
		t.Assert(v, 2)
		inner, _ := cache.Get(ctx, "inner")
		t.Assert(inner, 1)
		outer, _ := cache.Get(ctx, "outer")
		t.Assert(outer, 2)
	})
	// The function still runs only once for concurrent calls on the same key.
	gtest.C(t, func(t *gtest.T) {
		var (
			cache   = gcache.New()
			calls   = gtype.NewInt()
			results = make([]string, 50)
			wg      sync.WaitGroup
		)
		for i := range results {
			wg.Go(func() {
				v, err := cache.GetOrSetFuncLock(ctx, "key", func(ctx context.Context) (any, error) {
					calls.Add(1)
					time.Sleep(10 * time.Millisecond)
					return "value", nil
				}, 0)
				if err != nil {
					results[i] = err.Error()
					return
				}
				results[i] = v.String()
			})
		}
		wg.Wait()
		t.Assert(calls.Val(), 1)
		for _, result := range results {
			t.Assert(result, "value")
		}
	})
	// A value set while the function runs is kept and returned.
	gtest.C(t, func(t *gtest.T) {
		var (
			cache = gcache.New()
			done  = make(chan struct{})
			v     *gvar.Var
			err   error
		)
		go func() {
			defer close(done)
			v, err = cache.GetOrSetFuncLock(ctx, "key", func(ctx context.Context) (any, error) {
				if err := cache.Set(ctx, "key", "set", 0); err != nil {
					return nil, err
				}
				return "func", nil
			}, 0)
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Set inside the function of GetOrSetFuncLock did not return")
		}
		t.AssertNil(err)
		t.Assert(v, "set")
		v, err = cache.Get(ctx, "key")
		t.AssertNil(err)
		t.Assert(v, "set")
	})
}
