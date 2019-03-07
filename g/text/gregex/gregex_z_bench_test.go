// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gregex_test

import (
    "github.com/gogf/gf/g/text/gregex"
    "testing"
)

var pattern = `(.+):(\d+)`
var src     = "johng.cn:80"
var replace = "johng.cn"

func BenchmarkValidate(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.Validate(pattern)
    }
}

func BenchmarkIsMatch(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.IsMatch(pattern, []byte(src))
    }
}

func BenchmarkIsMatchString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.IsMatchString(pattern, src)
    }
}

func BenchmarkMatchString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.MatchString(pattern, src)
    }
}

func BenchmarkMatchAllString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.MatchAllString(pattern, src)
    }
}

func BenchmarkReplace(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.Replace(pattern, []byte(replace), []byte(src))
    }
}

func BenchmarkReplaceString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gregex.ReplaceString(pattern, replace, src)
    }
}
