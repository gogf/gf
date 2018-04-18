// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package test

import (
    "testing"
)

var m1 = make(map[int]int)
var m2 = make(map[interface{}]interface{})

func BenchmarkMapIntInt_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        m1[i] = i
    }
}

func BenchmarkMapIntInt_Search(b *testing.B) {
    for i := 0; i < b.N; i++ {
        if _, ok := m1[i]; ok {

        }
    }
}

func BenchmarkMapIntInt_Remove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        delete(m1, i)
    }
}


func BenchmarkMapInterface_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        m2[i] = i
    }
}

func BenchmarkMapInterface_Search(b *testing.B) {
    for i := 0; i < b.N; i++ {
        if _, ok := m2[i]; ok {

        }
    }
}

func BenchmarkMapInterface_Remove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        delete(m2, i)
    }
}