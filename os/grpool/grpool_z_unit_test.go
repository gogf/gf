// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpool_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err   error
			wg    = sync.WaitGroup{}
			array = garray.NewArray(true)
			size  = 100
		)
		wg.Add(size)
		for i := 0; i < size; i++ {
			err = grpool.Add(ctx, func(ctx context.Context) {
				array.Append(1)
				wg.Done()
			})
			t.AssertNil(err)
		}
		wg.Wait()

		time.Sleep(100 * time.Millisecond)

		t.Assert(array.Len(), size)
		t.Assert(grpool.Jobs(), 0)
		t.Assert(grpool.Size(), 0)
	})
}

func Test_Limit1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			wg    = sync.WaitGroup{}
			array = garray.NewArray(true)
			size  = 100
			pool  = grpool.New(10)
		)
		wg.Add(size)
		for i := 0; i < size; i++ {
			pool.Add(ctx, func(ctx context.Context) {
				array.Append(1)
				wg.Done()
			})
		}
		wg.Wait()
		t.Assert(array.Len(), size)
	})
}

func Test_Limit2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err   error
			wg    = sync.WaitGroup{}
			array = garray.NewArray(true)
			size  = 100
			pool  = grpool.New(1)
		)
		wg.Add(size)
		for i := 0; i < size; i++ {
			err = pool.Add(ctx, func(ctx context.Context) {
				defer wg.Done()
				array.Append(1)
			})
			t.AssertNil(err)
		}
		wg.Wait()
		t.Assert(array.Len(), size)
	})
}

func Test_Limit3(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			array = garray.NewArray(true)
			size  = 1000
			pool  = grpool.New(100)
		)
		t.Assert(pool.Cap(), 100)
		for i := 0; i < size; i++ {
			pool.Add(ctx, func(ctx context.Context) {
				array.Append(1)
				time.Sleep(2 * time.Second)
			})
		}
		time.Sleep(time.Second)
		t.Assert(pool.Size(), 100)
		t.Assert(pool.Jobs(), 900)
		t.Assert(array.Len(), 100)
		pool.Close()
		time.Sleep(2 * time.Second)
		t.Assert(pool.Size(), 0)
		t.Assert(pool.Jobs(), 900)
		t.Assert(array.Len(), 100)
		t.Assert(pool.IsClosed(), true)
		t.AssertNE(pool.Add(ctx, func(ctx context.Context) {}), nil)
	})
}

func Test_Limit4(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			array = garray.NewArray(true)
			size  = 1000
			limit = 100
			pool  = grpool.NewWithOption(grpool.PoolOption{
				Limit: 100,
				LimitChanger: func(ctx context.Context, val *atomic.Int64) (changed bool) {
					return val.Swap(int64(limit)) != int64(limit)
				},
			})
		)
		t.Assert(pool.Cap(), 100)
		for i := 0; i < size; i++ {
			pool.Add(ctx, func(ctx context.Context) {
				array.Append(1)
				time.Sleep(2 * time.Second)
			})
		}
		time.Sleep(time.Second)
		t.Assert(pool.Size(), 100)
		t.Assert(pool.Jobs(), 900)
		t.Assert(array.Len(), 100)
		limit = 50
		time.Sleep(time.Second * 2)
		t.Assert(pool.Size(), 50)
		t.Assert(pool.Jobs(), 850)
		t.Assert(array.Len(), 150)
		limit = 100
		time.Sleep(time.Second * 2)
		t.Assert(pool.Size(), 100)
		t.Assert(pool.Jobs(), 750)
		t.Assert(array.Len(), 250)
		pool.Close()
		time.Sleep(2 * time.Second)
		t.Assert(pool.Size(), 0)
		t.Assert(pool.Jobs(), 750)
		t.Assert(array.Len(), 250)
		t.Assert(pool.IsClosed(), true)
		t.AssertNE(pool.Add(ctx, func(ctx context.Context) {}), nil)
	})
}

func Test_ParseAndResume(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			array = garray.NewArray(true)
			size  = 1000
			pool  = grpool.New(100)
		)
		t.Assert(pool.Cap(), 100)
		for i := 0; i < size; i++ {
			pool.Add(ctx, func(ctx context.Context) {
				array.Append(1)
				time.Sleep(2 * time.Second)
			})
		}
		time.Sleep(time.Second)
		t.Assert(pool.Size(), 100)
		t.Assert(pool.Jobs(), 900)
		t.Assert(array.Len(), 100)
		pool.Parse()
		time.Sleep(time.Second * 2)
		t.Assert(pool.Size(), 0)
		t.Assert(pool.Jobs(), 900)
		t.Assert(array.Len(), 100)
		t.Assert(pool.IsParsed(), true)
		pool.Resume()
		time.Sleep(time.Second * 2)
		t.Assert(pool.Size(), 100)
		t.Assert(pool.Jobs(), 800)
		t.Assert(array.Len(), 200)
		t.Assert(pool.IsParsed(), false)
		pool.Close()
		time.Sleep(2 * time.Second)
		t.Assert(pool.Size(), 0)
		t.Assert(pool.Jobs(), 800)
		t.Assert(array.Len(), 200)
		t.Assert(pool.IsClosed(), true)
		t.AssertNE(pool.Add(ctx, func(ctx context.Context) {}), nil)
	})
}

func Test_AddWithRecover(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err   error
			array = garray.NewArray(true)
		)
		err = grpool.AddWithRecover(ctx, func(ctx context.Context) {
			array.Append(1)
			panic(1)
		}, func(ctx context.Context, err error) {
			array.Append(1)
		})
		t.AssertNil(err)
		err = grpool.AddWithRecover(ctx, func(ctx context.Context) {
			panic(1)
			array.Append(1)
		}, nil)
		t.AssertNil(err)

		time.Sleep(500 * time.Millisecond)

		t.Assert(array.Len(), 2)
	})
}
