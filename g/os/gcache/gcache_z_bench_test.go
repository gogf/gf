// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*" -benchmem

package gcache_test

import (
    "gitee.com/johng/gf/g/os/gcache"
    "testing"
    "sync"
)

var (
    c    = gcache.New()
    clru = gcache.New(10000)
    mInt = make(map[int]int)
    mMap = make(map[interface{}]interface{})

    muInt = sync.RWMutex{}
    muMap = sync.RWMutex{}
)

func Benchmark_CacheSet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        c.Set(i, i, 0)
    }
}

func Benchmark_CacheGet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        c.Get(i)
    }
}

func Benchmark_CacheRemove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        c.Remove(i)
    }
}

func Benchmark_CacheLruSet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        clru.Set(i, i, 0)
    }
}

func Benchmark_CacheLruGet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        clru.Get(i)
    }
}

func Benchmark_CacheLruRemove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        clru.Remove(i)
    }
}

func Benchmark_InterfaceMapWithLockSet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        muMap.Lock()
        mMap[i] = i
        muMap.Unlock()
    }
}

func Benchmark_InterfaceMapWithLockGet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        muMap.RLock()
        if _, ok := mMap[i]; ok {

        }
        muMap.RUnlock()
    }
}

func Benchmark_InterfaceMapWithLockRemove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        muMap.Lock()
        delete(mMap, i)
        muMap.Unlock()
    }
}

func Benchmark_IntMapWithLockWithLockSet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        muInt.Lock()
        mInt[i] = i
        muInt.Unlock()
    }
}

func Benchmark_IntMapWithLockGet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        muInt.RLock()
        if _, ok := mInt[i]; ok {

        }
        muInt.RUnlock()
    }
}

func Benchmark_IntMapWithLockRemove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        muInt.Lock()
        delete(mInt, i)
        muInt.Unlock()
    }
}
