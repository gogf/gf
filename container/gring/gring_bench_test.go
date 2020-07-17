// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

// go test *.go -bench=".*"

package gring_test

import (
	"testing"

	"github.com/jin502437344/gf/container/gring"
)

var length = 10000
var ringObject = gring.New(length, true)

func BenchmarkRing_Put(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ringObject.Put(i)
			i++
		}
	})
}

func BenchmarkRing_Next(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ringObject.Next()
			i++
		}
	})
}

func BenchmarkRing_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ringObject.Set(i)
			i++
		}
	})
}

func BenchmarkRing_Len(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ringObject.Len()
			i++
		}
	})
}

func BenchmarkRing_Cap(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ringObject.Cap()
			i++
		}
	})
}
