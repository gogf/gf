// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gring

import (
	"testing"
)

var length = 10000
var r1 = New(length)

func BenchmarkRing_Put(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r1.Put(i)
	}
}

func BenchmarkRing_Next(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r1.Next()
	}
}

func BenchmarkRing_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r1.Set(i)
	}
}

func BenchmarkRing_Len(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r1.Len()
	}
}

func BenchmarkRing_Cap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r1.Cap()
	}
}
