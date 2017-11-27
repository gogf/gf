package main

import (
    "testing"
    "gitee.com/johng/gf/g/container/gmap"
)




func BenchmarkSet1(b *testing.B) {
    a  := gmap.NewStringStringMap()
    m1 := make(map[string]string)
    m2 := make(map[string]string)
    for i := 0; i < 1000000; i ++ {
        m1[string(i)] = string(i)
    }
    for i := 0; i < 1000000; i ++ {
        m2[string(i)] = string(i) + "_2"
    }
    a.BatchSet(m1)
    a.BatchSet(m2)
}

func BenchmarkSet2 (b *testing.B) {
    a  := gmap.NewStringStringMap()
    m1 := make(map[string]string)
    m2 := make(map[string]string)
    for i := 0; i < 1000000; i ++ {
        m1[string(i)] = string(i)
    }
    for i := 0; i < 1000000; i ++ {
        m2[string(i)] = string(i) + "_2"
    }
    a.BatchSet2(m1)
    a.BatchSet2(m2)
}