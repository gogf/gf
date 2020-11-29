// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -count=1

package grpool_test

import (
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/grpool"
	"github.com/gogf/gf/test/gtest"
)

func Test_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		wg := sync.WaitGroup{}
		array := garray.NewArray(true)
		size := 100
		wg.Add(size)
		for i := 0; i < size; i++ {
			grpool.Add(func() {
				array.Append(1)
				wg.Done()
			})
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
		wg := sync.WaitGroup{}
		array := garray.NewArray(true)
		size := 100
		pool := grpool.New(10)
		wg.Add(size)
		for i := 0; i < size; i++ {
			pool.Add(func() {
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
			wg    = sync.WaitGroup{}
			array = garray.NewArray(true)
			size  = 100
			pool  = grpool.New(1)
		)
		wg.Add(size)
		for i := 0; i < size; i++ {
			pool.Add(func() {
				defer wg.Done()
				array.Append(1)
			})
		}
		wg.Wait()
		t.Assert(array.Len(), size)
	})
}

func Test_Limit3(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewArray(true)
		size := 1000
		pool := grpool.New(100)
		t.Assert(pool.Cap(), 100)
		for i := 0; i < size; i++ {
			pool.Add(func() {
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
		t.AssertNE(pool.Add(func() {}), nil)
	})
}

func Test_AddWithRecover(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.NewArray(true)
		grpool.AddWithRecover(func() {
			array.Append(1)
			panic(1)
		}, func(err error) {
			array.Append(1)
		})
		grpool.AddWithRecover(func() {
			panic(1)
			array.Append(1)
		})
		time.Sleep(500 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}
