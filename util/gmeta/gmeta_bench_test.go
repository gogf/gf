// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmeta_test

import (
	"github.com/gogf/gf/v2/util/gmeta"
	"testing"
)

type A struct {
	gmeta.Meta `tag:"123" orm:"456"`
	Id         int
	Name       string
}

var (
	a1 A
	a2 *A
)

func Benchmark_Data_Struct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gmeta.Data(a1)
	}
}

func Benchmark_Data_Pointer1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gmeta.Data(a2)
	}
}

func Benchmark_Data_Pointer2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gmeta.Data(&a2)
	}
}

func Benchmark_Data_Get_Struct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gmeta.Get(a1, "tag")
	}
}

func Benchmark_Data_Get_Pointer1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gmeta.Get(a2, "tag")
	}
}

func Benchmark_Data_Get_Pointer2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gmeta.Get(&a2, "tag")
	}
}
