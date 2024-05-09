// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gconv

import (
	"reflect"
	"testing"
)

type structType struct {
	Name  string
	Score int
	Age   int
	ID    int
}

type structType8 struct {
	Name        string  `json:"name"   `
	CategoryId  string  `json:"category-Id" `
	Price       float64 `json:"price"    `
	Code        string  `json:"code"       `
	Image       string  `json:"image"   `
	Description string  `json:"description" `
	Status      int     `json:"status"   `
	IdType      int     `json:"id-type"`
	Score       int
	Age         int
	ID          int
}

var (
	structMap = map[string]interface{}{
		"name":  "gf",
		"score": 100,
		"Age":   98,
		"ID":    199,
	}

	structMapFields8 = map[string]interface{}{
		"name":  "gf",
		"score": 100,
		"Age":   98,
		"ID":    199,

		"category-Id": "1",
		"price":       198.09,
		"code":        "1",
		"image":       "https://goframe.org",
		"description": "This is the data for testing eight fields",
		"status":      1,
		"id-type":     2,
	}

	structObj = structType{
		Name:  "john",
		Score: 60,
		Age:   98,
		ID:    199,
	}
	structPointer = &structType{
		Name:  "john",
		Score: 60,
	}
	structPointer8   = &structType8{}
	structPointerNil *structType

	// struct slice
	structSliceNil []structType
	structSlice    = []structType{
		{Name: "john", Score: 60},
		{Name: "smith", Score: 100},
	}
	// struct pointer slice
	structPointerSliceNil []*structType
	structPointerSlice    = []*structType{
		{Name: "john", Score: 60},
		{Name: "smith", Score: 100},
	}
)

func Benchmark_Struct_Basic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Struct(structMap, structPointer)
	}
}

func Benchmark_doStruct_Fields8_Basic_MapToStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doStruct(structMapFields8, structPointer8, map[string]string{}, "")
	}
}

// *struct -> **struct
func Benchmark_Reflect_PPStruct_PStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v1 := reflect.ValueOf(&structPointerNil)
		v2 := reflect.ValueOf(structPointer)
		//if v1.Kind() == reflect.Ptr {
		//	if elem := v1.Elem(); elem.Type() == v2.Type() {
		//		elem.Set(v2)
		//	}
		//}
		v1.Elem().Set(v2)
	}
}

func Benchmark_Struct_PPStruct_PStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Struct(structPointer, &structPointerNil)
	}
}

// struct -> *struct
func Benchmark_Reflect_PStruct_Struct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v1 := reflect.ValueOf(structPointer)
		v2 := reflect.ValueOf(structObj)
		//if v1.Kind() == reflect.Ptr {
		//	if elem := v1.Elem(); elem.Type() == v2.Type() {
		//		elem.Set(v2)
		//	}
		//}
		v1.Elem().Set(v2)
	}
}

func Benchmark_Struct_PStruct_Struct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Struct(structObj, structPointer)
	}
}

// []struct -> *[]struct
func Benchmark_Reflect_PStructs_Structs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v1 := reflect.ValueOf(&structSliceNil)
		v2 := reflect.ValueOf(structSlice)
		//if v1.Kind() == reflect.Ptr {
		//	if elem := v1.Elem(); elem.Type() == v2.Type() {
		//		elem.Set(v2)
		//	}
		//}
		v1.Elem().Set(v2)
	}
}

func Benchmark_Structs_PStructs_Structs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Structs(structSlice, &structSliceNil)
	}
}

// []*struct -> *[]*struct
func Benchmark_Reflect_PPStructs_PStructs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v1 := reflect.ValueOf(&structPointerSliceNil)
		v2 := reflect.ValueOf(structPointerSlice)
		//if v1.Kind() == reflect.Ptr {
		//	if elem := v1.Elem(); elem.Type() == v2.Type() {
		//		elem.Set(v2)
		//	}
		//}
		v1.Elem().Set(v2)
	}
}

func Benchmark_Structs_PPStructs_PStructs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Structs(structPointerSlice, &structPointerSliceNil)
	}
}
