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

	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/os/grpool"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_Basic(t *testing.T) {
	gtest.Case(t, func() {
		wg := sync.WaitGroup{}
		array := garray.NewArray()
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
		gtest.Assert(array.Len(), size)
		gtest.Assert(grpool.Jobs(), 0)
		gtest.Assert(grpool.Size(), 0)
	})
}

func Test_Limit1(t *testing.T) {
	gtest.Case(t, func() {
		wg := sync.WaitGroup{}
		array := garray.NewArray()
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
		gtest.Assert(array.Len(), size)
	})
}

func Test_Limit2(t *testing.T) {
	gtest.Case(t, func() {
		wg := sync.WaitGroup{}
		array := garray.NewArray()
		size := 100
		pool := grpool.New(1)
		wg.Add(size)
		for i := 0; i < size; i++ {
			pool.Add(func() {
				array.Append(1)
				wg.Done()
			})
		}
		wg.Wait()
		gtest.Assert(array.Len(), size)
	})
}

func Test_Limit3(t *testing.T) {
	gtest.Case(t, func() {
		array := garray.NewArray()
		size := 1000
		pool := grpool.New(100)
		gtest.Assert(pool.Cap(), 100)
		for i := 0; i < size; i++ {
			pool.Add(func() {
				array.Append(1)
				time.Sleep(2 * time.Second)
			})
		}
		time.Sleep(time.Second)
		gtest.Assert(pool.Size(), 100)
		gtest.Assert(pool.Jobs(), 900)
		gtest.Assert(array.Len(), 100)
		pool.Close()
		time.Sleep(2 * time.Second)
		gtest.Assert(pool.Size(), 0)
		gtest.Assert(pool.Jobs(), 900)
		gtest.Assert(array.Len(), 100)
		gtest.Assert(pool.IsClosed(), true)
		gtest.AssertNE(pool.Add(func() {}), nil)

	})
}
