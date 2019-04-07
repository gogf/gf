// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime_test

import (
	"testing"

	"github.com/gogf/gf/g/os/gtime"
)

func Benchmark_Second(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gtime.Second()
	}
}

func Benchmark_Millisecond(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gtime.Millisecond()
	}
}

func Benchmark_Microsecond(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gtime.Microsecond()
	}
}

func Benchmark_Nanosecond(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gtime.Nanosecond()
	}
}

func Benchmark_StrToTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gtime.StrToTime("2018-02-09T20:46:17.897Z")
	}
}

func Benchmark_ParseTimeFromContent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gtime.ParseTimeFromContent("2018-02-09T20:46:17.897Z")
	}
}

func Benchmark_NewFromTimeStamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gtime.NewFromTimeStamp(1542674930)
	}
}

func Benchmark_Date(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gtime.Date()
	}
}

func Benchmark_Datetime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gtime.Datetime()
	}
}

func Benchmark_SetTimeZone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gtime.SetTimeZone("Asia/Shanghai")
	}
}
