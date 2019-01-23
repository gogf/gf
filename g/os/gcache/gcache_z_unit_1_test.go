// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*" -benchmem

package gcache_test

import (
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)

func TestCache_Set(t *testing.T) {
    gtest.Case(t, func() {
        cache := gcache.New()
        cache.Set(1, 11, 0)
        gtest.Assert(cache.Get(1), 11)
    })
}

func TestCache_Set_Expire(t *testing.T) {
    gtest.Case(t, func() {
        cache := gcache.New()
        cache.Set(2, 22, 100)
        gtest.Assert(cache.Get(2), 22)
        time.Sleep(200*time.Millisecond)
        gtest.Assert(cache.Get(2), nil)
        time.Sleep(3*time.Second)
        gtest.Assert(cache.Size(), 0)
    })
}

func TestCache_Keys_Values(t *testing.T) {
    gtest.Case(t, func() {
        cache := gcache.New()
        for i := 0; i < 10; i++ {
            cache.Set(i, i*10, 0)
        }
        gtest.Assert(len(cache.Keys()), 10)
        gtest.Assert(len(cache.Values()), 10)
        gtest.AssertIN(0, cache.Keys())
        gtest.AssertIN(90, cache.Values())
    })
}

func TestCache_LRU(t *testing.T) {
    gtest.Case(t, func() {
        cache := gcache.New(2)
        for i := 0; i < 10; i++ {
            cache.Set(i, i, 0)
        }

        gtest.Assert(cache.Size(), 10)
        gtest.Assert(cache.Get(6), 6)
        time.Sleep(3*time.Second)
        gtest.Assert(cache.Size(), 2)
        gtest.Assert(cache.Get(6), 6)
        gtest.Assert(cache.Get(1), nil)
    })
}