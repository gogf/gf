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

// BenchmarkGTimeConverter_ComprehensiveScenarios benchmarks various gtime conversion scenarios
func BenchmarkGTimeConverter_ComprehensiveScenarios(b *testing.B) {
	// Set up test data
	utcTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
	gtimeVal := gtime.NewFromTime(utcTime)
	gtimePtr := gtimeVal
	gtimeValue := *gtimeVal
	
	// Set different local timezone for more realistic testing
	shanghaiLocation, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = shanghaiLocation
	
	// Benchmark 1: Direct type conversions (should be fastest)
	b.Run("DirectGTimeToTime", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = gconv.Time(gtimePtr)
		}
	})
	
	b.Run("DirectGTimeValueToTime", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = gconv.Time(gtimeValue)
		}
	})
	
	b.Run("DirectGTimeToGTime", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = gconv.GTime(gtimePtr)
		}
	})
	
	// Benchmark 2: Builtin converter scenarios
	b.Run("BuiltinGTimeStruct", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result gtime.Time
			_ = gconv.Struct(gtimePtr, &result)
		}
	})
	
	b.Run("BuiltinGTimePtrStruct", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result *gtime.Time
			_ = gconv.Struct(gtimePtr, &result)
		}
	})
	
	b.Run("BuiltinGTimeValueStruct", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result gtime.Time
			_ = gconv.Struct(gtimeValue, &result)
		}
	})
	
	// Benchmark 3: String conversion scenarios
	b.Run("GTimeToString", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = gconv.String(gtimePtr)
		}
	})
	
	b.Run("GTimeValueToString", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = gconv.String(gtimeValue)
		}
	})
	
	b.Run("StringToGTime", func(b *testing.B) {
		timeStr := "2025-09-16T11:32:42.878465Z"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = gconv.GTime(timeStr)
		}
	})
	
	// Benchmark 4: Map conversion scenarios (problematic in original issue)
	b.Run("MapToTime", func(b *testing.B) {
		mapData := map[string]interface{}{"time": gtimePtr}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = gconv.Time(mapData)
		}
	})
	
	b.Run("MapToGTime", func(b *testing.B) {
		mapData := map[string]interface{}{"time": gtimePtr}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = gconv.GTime(mapData)
		}
	})
	
	// Benchmark 5: Struct field conversion scenarios
	b.Run("StructFieldConversion", func(b *testing.B) {
		type TestStruct struct {
			Time time.Time `json:"time"`
		}
		mapData := map[string]interface{}{"Time": gtimePtr}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result TestStruct
			_ = gconv.Struct(mapData, &result)
		}
	})
	
	b.Run("StructGTimeFieldConversion", func(b *testing.B) {
		type TestStruct struct {
			Time gtime.Time `json:"time"`
		}
		mapData := map[string]interface{}{"Time": gtimePtr}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result TestStruct
			_ = gconv.Struct(mapData, &result)
		}
	})
	
	// Benchmark 6: Slice conversion scenarios (the main issue scenario)
	b.Run("SliceConversionToTime", func(b *testing.B) {
		sliceData := []map[string]interface{}{{"time": gtimePtr}}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result []time.Time
			_ = gconv.Structs(sliceData, &result)
		}
	})
	
	b.Run("SliceConversionToGTime", func(b *testing.B) {
		sliceData := []map[string]interface{}{{"time": gtimePtr}}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result []gtime.Time
			_ = gconv.Structs(sliceData, &result)
		}
	})
	
	b.Run("SliceConversionToGTimePtr", func(b *testing.B) {
		sliceData := []map[string]interface{}{{"time": gtimePtr}}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result []*gtime.Time
			_ = gconv.Structs(sliceData, &result)
		}
	})
}

// BenchmarkGTimeConverter_TimezoneImpact benchmarks timezone impact on performance
func BenchmarkGTimeConverter_TimezoneImpact(b *testing.B) {
	// Test performance with different timezones
	timezones := []struct {
		name string
		loc  *time.Location
	}{
		{"UTC", time.UTC},
		{"Shanghai", mustLoadLocation("Asia/Shanghai")},
		{"NewYork", mustLoadLocation("America/New_York")},
		{"London", mustLoadLocation("Europe/London")},
		{"Tokyo", mustLoadLocation("Asia/Tokyo")},
	}
	
	baseTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
	
	for _, tz := range timezones {
		testTime := baseTime.In(tz.loc)
		gtimeVal := gtime.NewFromTime(testTime)
		
		b.Run("DirectConversion_"+tz.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = gconv.Time(gtimeVal)
			}
		})
		
		b.Run("StringConversion_"+tz.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = gconv.String(gtimeVal)
			}
		})
		
		b.Run("StructsConversion_"+tz.name, func(b *testing.B) {
			sliceData := []map[string]interface{}{{"time": gtimeVal}}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var result []time.Time
				_ = gconv.Structs(sliceData, &result)
			}
		})
	}
}

// BenchmarkGTimeConverter_PrecisionImpact benchmarks precision impact on performance
func BenchmarkGTimeConverter_PrecisionImpact(b *testing.B) {
	// Test performance with different precision levels
	precisions := []struct {
		name  string
		nanos int
	}{
		{"Seconds", 0},
		{"Milliseconds", 123000000},
		{"Microseconds", 123456000},
		{"Nanoseconds", 123456789},
	}
	
	baseTime := time.Date(2025, 9, 16, 11, 32, 42, 0, time.UTC)
	
	for _, p := range precisions {
		testTime := baseTime.Add(time.Duration(p.nanos))
		gtimeVal := gtime.NewFromTime(testTime)
		
		b.Run("Conversion_"+p.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = gconv.Time(gtimeVal)
			}
		})
		
		b.Run("StringRoundTrip_"+p.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				str := gconv.String(gtimeVal)
				_ = gconv.GTime(str)
			}
		})
		
		b.Run("StructsConversion_"+p.name, func(b *testing.B) {
			sliceData := []map[string]interface{}{{"time": gtimeVal}}
			b.ResetTimer()  
			for i := 0; i < b.N; i++ {
				var result []time.Time
				_ = gconv.Structs(sliceData, &result)
			}
		})
	}
}

// BenchmarkGTimeConverter_MemoryAllocation benchmarks memory allocation patterns
func BenchmarkGTimeConverter_MemoryAllocation(b *testing.B) {
	utcTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
	gtimeVal := gtime.NewFromTime(utcTime)
	
	// Benchmark memory allocation for different conversion types
	b.Run("DirectConversion_Allocs", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = gconv.Time(gtimeVal)
		}
	})
	
	b.Run("BuiltinConverter_Allocs", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result gtime.Time
			_ = gconv.Struct(gtimeVal, &result)
		}
	})
	
	b.Run("StringConversion_Allocs", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = gconv.String(gtimeVal)
		}
	})
	
	b.Run("SliceConversion_Allocs", func(b *testing.B) {
		sliceData := []map[string]interface{}{{"time": gtimeVal}}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result []time.Time
			_ = gconv.Structs(sliceData, &result)
		}
	})
}

// BenchmarkGTimeConverter_ComparisonWithStandard compares performance with standard library
func BenchmarkGTimeConverter_ComparisonWithStandard(b *testing.B) {
	utcTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
	gtimeVal := gtime.NewFromTime(utcTime)
	timeStr := "2025-09-16T11:32:42.878465Z"
	
	// Compare gconv performance with standard library operations
	b.Run("GConv_TimeConversion", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = gconv.Time(gtimeVal)
		}
	})
	
	b.Run("Standard_TimeParsing", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = time.Parse(time.RFC3339, timeStr)
		}
	})
	
	b.Run("GConv_StringConversion", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = gconv.String(gtimeVal)
		}
	})
	
	b.Run("Standard_TimeFormatting", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = utcTime.Format(time.RFC3339)
		}
	})
	
	b.Run("GConv_StructConversion", func(b *testing.B) {
		type TimeStruct struct {
			Time time.Time `json:"time"`
		}
		mapData := map[string]interface{}{"Time": gtimeVal}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result TimeStruct
			_ = gconv.Struct(mapData, &result)
		}
	})
}

// Helper function
func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}