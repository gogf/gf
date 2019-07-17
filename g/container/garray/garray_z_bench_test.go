// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package garray_test

import (
	"testing"

	"github.com/gogf/gf/g/container/garray"
)

var (
	sortedIntArray = garray.NewSortedIntArray()
)

func BenchmarkSortedIntArray_Add(b *testing.B) {
	b.N = 1000
	for i := 0; i < b.N; i++ {
		sortedIntArray.Add(i)
	}
}

func BenchmarkSortedIntArray_Search(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sortedIntArray.Search(i)
	}
}

func BenchmarkSortedIntArray_PopLeft(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sortedIntArray.PopLeft()
	}
}

func BenchmarkSortedIntArray_PopRight(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sortedIntArray.PopLeft()
	}
}
