// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -count=1

package grpool_test

import (
	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/os/grpool"
	"github.com/gogf/gf/g/test/gtest"
	"sync"
	"testing"
	"time"
)


func Test_Basic(t *testing.T) {
	gtest.Case(t, func() {
		wg    := sync.WaitGroup{}
		array := garray.NewArray()
		size  := 100000
		wg.Add(size)
		for i := 0; i < size; i++ {
			grpool.Add(func() {
				array.Append(1)
				wg.Done()
			})
		}
		wg.Wait()
		gtest.Assert(array.Len(), size)
	})

	gtest.Case(t, func() {
		array := garray.NewArray()
		size  := 100000
		pool  := grpool.New(10000)
		for i := 0; i < size; i++ {
			pool.Add(func() {
				array.Append(1)
				time.Sleep(2*time.Second)
			})
		}
		time.Sleep(time.Second)
		gtest.Assert(pool.Size(), 10000)
		gtest.Assert(pool.Jobs(), 90000)
		gtest.Assert(array.Len(), 10000)
		pool.Close()
		time.Sleep(2*time.Second)
		gtest.Assert(pool.Size(), 10000)
		gtest.Assert(pool.Jobs(), 90000)
		gtest.Assert(array.Len(), 10000)
	})
}