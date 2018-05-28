// Copyright 2017-2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gregx

import (
    "testing"
)

var pattern = `(.+):(\d+)`
var src     = "johng.cn:80"
var replace = "johng.cn"

func BenchmarkValidate(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Validate(pattern)
    }
}

func BenchmarkIsMatch(b *testing.B) {
    for i := 0; i < b.N; i++ {
        IsMatch(pattern, []byte(src))
    }
}

func BenchmarkIsMatchString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        IsMatchString(pattern, src)
    }
}

func BenchmarkMatchString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        MatchString(pattern, src)
    }
}

func BenchmarkMatchAllString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        MatchAllString(pattern, src)
    }
}

func BenchmarkReplace(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Replace(pattern, []byte(replace), []byte(src))
    }
}

func BenchmarkReplaceString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ReplaceString(pattern, replace, src)
    }
}
