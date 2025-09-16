// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// BenchmarkTime_TimezonePreservation benchmarks the timezone preservation optimization
func BenchmarkTime_TimezonePreservation(b *testing.B) {
	// Create test data
	gmtLocation, _ := time.LoadLocation("GMT")
	dbTime := time.Date(2025, 9, 15, 7, 45, 40, 0, gmtLocation)
	gtimeVal := gtime.NewFromTime(dbTime)

	b.ResetTimer()

	b.Run("DirectGTimeConversion", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = gconv.Time(gtimeVal)
		}
	})

	b.Run("MapToTimeConversion", func(b *testing.B) {
		mapData := map[string]interface{}{"now": gtimeVal}
		for i := 0; i < b.N; i++ {
			_ = gconv.Time(mapData)
		}
	})

	b.Run("StructsConversion", func(b *testing.B) {
		result := []map[string]interface{}{{"now": gtimeVal}}
		for i := 0; i < b.N; i++ {
			var nowResult []time.Time
			_ = gconv.Structs(result, &nowResult)
		}
	})
}

// BenchmarkGTime_Optimization benchmarks the GTime function optimizations
func BenchmarkGTime_Optimization(b *testing.B) {
	// Create test data
	gmtLocation, _ := time.LoadLocation("GMT")
	dbTime := time.Date(2025, 9, 15, 7, 45, 40, 0, gmtLocation)
	gtimeVal := gtime.NewFromTime(dbTime)

	b.ResetTimer()

	b.Run("DirectGTimeToGTime", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = gconv.GTime(gtimeVal)
		}
	})

	b.Run("TimeToGTime", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = gconv.GTime(dbTime)
		}
	})

	b.Run("StringToGTime", func(b *testing.B) {
		timeStr := "2025-09-15T07:45:40Z"
		for i := 0; i < b.N; i++ {
			_ = gconv.GTime(timeStr)
		}
	})
}
