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

func waitUntil(timeout time.Duration, condition func() bool) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return condition()
}

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
		var limit atomic.Int64
		limit.Store(100)
		var (
			array = garray.NewArray(true)
			size  = 1000
			pool  = grpool.NewWithOption(grpool.PoolOption{
				Limit: 100,
				LimitChanger: func(ctx context.Context, val *atomic.Int64) (changed bool) {
					v := limit.Load()
					return val.Swap(v) != v
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
		t.Assert(waitUntil(2*time.Second, func() bool {
			return pool.Size() == 100 && pool.Jobs() <= 900 && array.Len() > 0
		}), true)
		t.Assert(pool.Jobs()+array.Len(), size)

		limit.Store(50)
		t.Assert(waitUntil(4*time.Second, func() bool {
			return pool.Size() <= 50 && pool.Size() >= 0
		}), true)
		t.Assert(pool.Size(), 50)
		t.Assert(pool.Jobs()+array.Len(), size)

		jobsBeforeIncrease := pool.Jobs()
		arrayBeforeIncrease := array.Len()
		limit.Store(100)
		t.Assert(waitUntil(4*time.Second, func() bool {
			return pool.Size() > 50
		}), true)
		t.AssertLE(pool.Size(), 100)
		t.AssertGT(pool.Jobs(), 0)
		t.AssertLT(pool.Jobs(), jobsBeforeIncrease)
		t.AssertGT(array.Len(), arrayBeforeIncrease)
		t.Assert(pool.Jobs()+array.Len(), size)

		pool.Close()
		time.Sleep(2 * time.Second)
		t.Assert(pool.Size(), 0)
		t.Assert(pool.Jobs()+array.Len(), size)
		t.Assert(pool.IsClosed(), true)
		t.AssertNE(pool.Add(ctx, func(ctx context.Context) {}), nil)
	})
}

func Test_PauseAndResume(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err                 error
			started             atomic.Int64
			completed           atomic.Int64
			firstBatchReleased  = make(chan struct{})
			secondBatchReleased = make(chan struct{})
			size                = 20
			pool                = grpool.New(10)
		)
		waitForCondition := func(condition func() bool) {
			t.Helper()
			deadline := time.Now().Add(5 * time.Second)
			for time.Now().Before(deadline) {
				if condition() {
					return
				}
				time.Sleep(10 * time.Millisecond)
			}
			t.Fatal("timed out waiting for pool state")
		}
		t.Assert(pool.Cap(), 10)
		for i := 0; i < size; i++ {
			err = pool.Add(ctx, func(ctx context.Context) {
				current := started.Add(1)
				if current <= 10 {
					<-firstBatchReleased
				} else {
					<-secondBatchReleased
				}
				completed.Add(1)
			})
			t.AssertNil(err)
		}
		waitForCondition(func() bool {
			return pool.Size() == 10 && pool.Jobs() == 10 && started.Load() == 10
		})
		t.Assert(pool.Pause(), true)
		close(firstBatchReleased)
		waitForCondition(func() bool {
			return pool.Size() == 0 && pool.Jobs() == 10 && started.Load() == 10 && completed.Load() == 10
		})
		t.Assert(pool.IsPaused(), true)
		t.Assert(pool.Resume(), true)
		waitForCondition(func() bool {
			return pool.Size() == 10 && pool.Jobs() == 0 && started.Load() == 20 && completed.Load() == 10
		})
		t.Assert(pool.IsPaused(), false)
		close(secondBatchReleased)
		waitForCondition(func() bool {
			return pool.Size() == 0 && pool.Jobs() == 0 && completed.Load() == 20
		})
		pool.Close()
		t.Assert(pool.IsClosed(), true)
		t.Assert(pool.Size(), 0)
		t.Assert(pool.Jobs(), 0)
		t.Assert(completed.Load(), int64(20))
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
