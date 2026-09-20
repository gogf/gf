// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

// Benchmarks for the struct validation path, which is the mostly used path for
// the standard HTTP handler signature validation.
package gvalid_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/util/gvalid"
)

// benchCheckUser is a small struct for the common validation benchmarking.
type benchCheckUser struct {
	Name    string `json:"name"    v:"required|length:1,30#name is required"`
	Age     int    `json:"age"     v:"between:1,100#age should be between 1 and 100"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Tags    string `json:"tags"`
	Remark  string `json:"remark"`
}

var (
	benchCheckCtx       = context.Background()
	benchCheckUserValue = benchCheckUser{
		Name:    "john",
		Age:     18,
		Email:   "john@example.com",
		Address: "street 1",
		Tags:    "a,b",
		Remark:  "hello",
	}
	benchCheckUserMap = map[string]any{
		"name":    "john",
		"age":     18,
		"email":   "john@example.com",
		"address": "street 1",
		"tags":    "a,b",
		"remark":  "hello",
	}
)

// benchBigReq simulates a real request struct with many fields and rules.
type benchBigReq struct {
	Id        int    `json:"id"        v:"required|min:1#id is required"`
	Name      string `json:"name"      v:"required|length:1,30#name is required"`
	Email     string `json:"email"     v:"email"`
	Phone     string `json:"phone"`
	Age       int    `json:"age"       v:"between:1,100"`
	Gender    string `json:"gender"    v:"in:male,female"`
	City      string `json:"city"`
	Address   string `json:"address"`
	ZipCode   string `json:"zipCode"`
	Country   string `json:"country"`
	Remark    string `json:"remark"`
	Tags      string `json:"tags"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	Score     int    `json:"score"     v:"min:0|max:1000"`
	Level     int    `json:"level"`
	Status    int    `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Extra1    string `json:"extra1"`
	Extra2    string `json:"extra2"`
	Extra3    string `json:"extra3"`
	Extra4    string `json:"extra4"`
	Extra5    string `json:"extra5"`
}

var benchBigReqValue = benchBigReq{
	Id: 1, Name: "john", Email: "john@example.com", Phone: "123456",
	Age: 30, Gender: "male", City: "sh", Address: "street 1", ZipCode: "200000",
	Country: "cn", Remark: "hello", Tags: "a,b", Avatar: "a.png", Nickname: "j",
	Score: 100, Level: 2, Status: 1, CreatedAt: "2026-01-01", UpdatedAt: "2026-01-02",
	Extra1: "x", Extra2: "y", Extra3: "z", Extra4: "u", Extra5: "v",
}

// Test_ZZBenchSanity ensures the benchmarks below really trigger the validations.
func Test_ZZBenchSanity(t *testing.T) {
	// The valid data passes.
	if err := gvalid.New().Bail().Data(&benchCheckUserValue).Run(benchCheckCtx); err != nil {
		t.Fatalf("expect no error, but got: %v", err)
	}
	if err := gvalid.New().Bail().Data(&benchBigReqValue).Run(benchCheckCtx); err != nil {
		t.Fatalf("expect no error, but got: %v", err)
	}
	// The invalid data fails, which means the validation is really running.
	var user = benchCheckUser{Name: "", Age: 200}
	if err := gvalid.New().Bail().Data(&user).Run(benchCheckCtx); err == nil {
		t.Fatal("expect validation error, but got nil")
	}
	var req = benchBigReq{Age: 200}
	if err := gvalid.New().Bail().Data(&req).Run(benchCheckCtx); err == nil {
		t.Fatal("expect validation error, but got nil")
	}
}

// Benchmark_ParseTagValue measures the parsing of one validation tag value.
func Benchmark_ParseTagValue(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gvalid.ParseTagValue(`required|length:1,30#name is required`)
	}
}

// Benchmark_CheckStruct_Data measures the struct validation with the struct data.
func Benchmark_CheckStruct_Data(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := gvalid.New().Bail().Data(&benchCheckUserValue).Run(benchCheckCtx); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark_CheckStruct_Assoc measures the struct validation with the associated
// data map, which is the way the standard HTTP handler uses.
func Benchmark_CheckStruct_Assoc(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := gvalid.New().Bail().Data(&benchCheckUserValue).Assoc(benchCheckUserMap).Run(benchCheckCtx); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark_CheckStruct_BigReq measures the struct validation with a real request
// struct, which has many fields and validation rules.
func Benchmark_CheckStruct_BigReq(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := gvalid.New().Bail().Data(&benchBigReqValue).Run(benchCheckCtx); err != nil {
			b.Fatal(err)
		}
	}
}
