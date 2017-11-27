package main

import (
    "testing"
    "gitee.com/johng/gf/g/os/gcache"
)

var cache *gcache.Cache = gcache.New()

func BenchmarkSet(b *testing.B) {
    b.N = 1000000
    for i := 0; i < 1000000; i ++ {
        cache.Set(string(i), i, 0)
    }
}

func BenchmarkSetWithExpire(b *testing.B) {
    b.N = 1000000
    for i := 0; i < 1000000; i ++ {
        cache.Set(string(i), i, 60)
    }
}

func BenchmarkGet1(b *testing.B) {
    b.N = 1000000
    for i := 0; i < 1000000; i ++ {
        cache.Get(string(i))
    }
}

func BenchmarkGet2(b *testing.B) {
    b.N = 1000000
    for i := 0; i < 1000000; i ++ {
        cache.Get(string(i))
    }
}

func BenchmarkRemove(b *testing.B) {
    b.N = 1000000
    for i := 0; i < 1000000; i ++ {
        cache.Remove(string(i))
    }
}