// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gring_test

import (
	"testing"

	"github.com/gogf/gf/container/gring"
)

var length = 10000
var ring = gring.New(length, true)

func BenchmarkRing_Put(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ring.Put(i)
			i++
		}
	})
}

func BenchmarkRing_Next(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ring.Next()
			i++
		}
	})
}

func BenchmarkRing_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ring.Set(i)
			i++
		}
	})
}

func BenchmarkRing_Len(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ring.Len()
			i++
		}
	})
}

func BenchmarkRing_Cap(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ring.Cap()
			i++
		}
	})
}
