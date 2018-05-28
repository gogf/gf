// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gcache

import (
    "testing"
    "strconv"
)


func BenchmarkSet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Set(strconv.Itoa(i), strconv.Itoa(i), 0)
    }
}

func BenchmarkGet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Get(strconv.Itoa(i))
    }
}

func BenchmarkRemove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Remove(strconv.Itoa(i))
    }
}