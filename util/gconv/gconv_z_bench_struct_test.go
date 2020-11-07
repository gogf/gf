// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gconv

import (
	"testing"
)

type structType struct {
	Name  string
	Score int
}

var (
	structMap = map[string]interface{}{
		"name":  "gf",
		"score": 100,
	}
	structPointer = new(structType)
)

func Benchmark_Struct_Basic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Struct(structMap, structPointer)
	}
}
