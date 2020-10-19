// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gregex_test

import (
	"regexp"
	"testing"

	"github.com/gogf/gf/text/gregex"
)

var pattern = `(\w+).+\-\-\s*(.+)`
var src = `GF is best! -- John`

func Benchmark_GF_IsMatchString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gregex.IsMatchString(pattern, src)
	}
}

func Benchmark_GF_MatchString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gregex.MatchString(pattern, src)
	}
}

func Benchmark_Compile(b *testing.B) {
	var wcdRegexp = regexp.MustCompile(pattern)
	for i := 0; i < b.N; i++ {
		wcdRegexp.MatchString(src)
	}
}

func Benchmark_Compile_Actual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wcdRegexp := regexp.MustCompile(pattern)
		wcdRegexp.MatchString(src)
	}
}
