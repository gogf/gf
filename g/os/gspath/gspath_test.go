// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gspath

import (
    "testing"
)

var (
    sp = New()
)

func init() {
    sp.Add("/Users/john/Temp")
}

func Benchmark_Search(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sp.Search("1")
    }
}

func Benchmark_Search_None(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sp.Search("1000")
    }
}

func Benchmark_Search_IndexFiles(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sp.Search("1", "index.html")
    }
}
