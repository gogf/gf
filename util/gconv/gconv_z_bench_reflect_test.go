// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gconv_test

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Id   int
	Name string
}

var ptr = []*testStruct{
	{
		Id:   1,
		Name: "test1",
	},
	{
		Id:   2,
		Name: "test2",
	},
}

func init() {
	for i := 1; i <= 1000; i++ {
		ptr = append(ptr, &testStruct{
			Id:   1,
			Name: "test1",
		})
	}
}

func Benchmark_Reflect_ValueOf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflect.ValueOf(ptr)
	}
}

func Benchmark_Reflect_ValueOf_Kind(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflect.ValueOf(ptr).Kind()
	}
}

func Benchmark_Reflect_ValueOf_Interface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflect.ValueOf(ptr).Interface()
	}
}

func Benchmark_Reflect_ValueOf_Len(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflect.ValueOf(ptr).Len()
	}
}
