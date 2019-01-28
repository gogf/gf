// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gregex_test

import (
    "gitee.com/johng/gf/g/util/gregex"
    "testing"
)

var pattern      = `(.+):(\d+)`
var src          = "johng.cn:80"
var srcBytes     = []byte("johng.cn:80")
var replace      = "johng.cn"
var replaceBytes = []byte("johng.cn")

func BenchmarkValidate(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.Validate(pattern)
    }
}

func BenchmarkIsMatch(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.IsMatch(pattern, srcBytes)
    }
}

func BenchmarkIsMatchString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.IsMatchString(pattern, src)
    }
}

func BenchmarkMatch(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.Match(pattern, srcBytes)
    }
}

func BenchmarkMatchString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.MatchString(pattern, src)
    }
}

func BenchmarkMatchAll(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.MatchAll(pattern, srcBytes)
    }
}

func BenchmarkMatchAllString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.MatchAllString(pattern, src)
    }
}

func BenchmarkReplace(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.Replace(pattern, replaceBytes, srcBytes)
    }
}

func BenchmarkReplaceString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.ReplaceString(pattern, replace, src)
    }
}
