// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

// TestBuiltinGTimeConverter tests the builtin converter for gtime.Time types
func TestBuiltinGTimeConverter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Set up test environment with different timezone
		originalLocation := time.Local
		defer func() {
			time.Local = originalLocation
		}()
		
		shanghaiLocation, _ := time.LoadLocation("Asia/Shanghai")
		time.Local = shanghaiLocation
		
		// Test data with various timezones
		utcTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
		gtimeUTC := gtime.NewFromTime(utcTime)
		
		gmtLocation, _ := time.LoadLocation("GMT")
		gmtTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, gmtLocation)
		gtimeGMT := gtime.NewFromTime(gmtTime)
		
		estLocation, _ := time.LoadLocation("America/New_York")
		estTime := time.Date(2025, 9, 16, 7, 32, 42, 878465000, estLocation)
		gtimeEST := gtime.NewFromTime(estTime)
		
		// Test 1: Direct gtime.Time to gtime.Time conversion
		t.Logf("=== Test 1: Direct gtime.Time to gtime.Time conversion ===")
		var result1 gtime.Time
		err := gconv.Struct(gtimeUTC, &result1)
		t.AssertNil(err)
		t.Assert(result1.Location().String(), gtimeUTC.Location().String())
		t.Assert(result1.Equal(gtimeUTC), true)
		t.Logf("Original: %s, Result: %s", gtimeUTC.Time, result1.Time)
		
		// Test 2: *gtime.Time to *gtime.Time conversion
		t.Logf("=== Test 2: *gtime.Time to *gtime.Time conversion ===")
		var result2 *gtime.Time
		err = gconv.Struct(gtimeUTC, &result2)
		t.AssertNil(err)
		t.AssertNE(result2, nil)
		t.Assert(result2.Location().String(), gtimeUTC.Location().String())
		t.Assert(result2.Equal(gtimeUTC), true)
		t.Logf("Original: %s, Result: %s", gtimeUTC.Time, result2.Time)
		
		// Test 3: gtime.Time to *gtime.Time conversion
		t.Logf("=== Test 3: gtime.Time to *gtime.Time conversion ===")
		var result3 *gtime.Time
		err = gconv.Struct(*gtimeUTC, &result3)
		t.AssertNil(err)
		t.AssertNE(result3, nil)
		t.Assert(result3.Location().String(), gtimeUTC.Location().String())
		t.Assert(result3.Equal(gtimeUTC), true)
		t.Logf("Original: %s, Result: %s", gtimeUTC.Time, result3.Time)
		
		// Test 4: Multiple timezone preservation
		testCases := []struct {
			name     string
			input    *gtime.Time
			expected int // expected offset in seconds
		}{
			{"UTC", gtimeUTC, 0},
			{"GMT", gtimeGMT, 0},
			{"EST", gtimeEST, -4 * 3600}, // EST is UTC-4 in September
		}
		
		for _, tc := range testCases {
			t.Logf("=== Test 4.%s: %s timezone preservation ===", tc.name, tc.name)
			var result gtime.Time
			err := gconv.Struct(tc.input, &result)
			t.AssertNil(err)
			
			_, inputOffset := tc.input.Zone()
			_, resultOffset := result.Zone()
			
			t.Assert(resultOffset, inputOffset)
			t.Assert(result.Equal(tc.input), true)
			t.Logf("%s - Original: %s (offset: %d), Result: %s (offset: %d)", 
				tc.name, tc.input.Time, inputOffset, result.Time, resultOffset)
		}
	})
}

// TestBuiltinGTimeConverter_EdgeCases tests edge cases for the builtin gtime converter
func TestBuiltinGTimeConverter_EdgeCases(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test 1: Nil *gtime.Time conversion
		t.Logf("=== Test 1: Nil *gtime.Time conversion ===")
		var nilGtime *gtime.Time = nil
		var result1 gtime.Time
		err := gconv.Struct(nilGtime, &result1)
		t.AssertNil(err)
		t.Assert(result1.IsZero(), true)
		t.Logf("Nil gtime converted to zero gtime: %s", result1.Time)
		
		// Test 2: Zero gtime.Time conversion
		t.Logf("=== Test 2: Zero gtime.Time conversion ===")
		zeroGtime := gtime.Time{}
		var result2 gtime.Time
		err = gconv.Struct(zeroGtime, &result2)
		t.AssertNil(err)
		t.Assert(result2.IsZero(), true)
		t.Logf("Zero gtime preserved: %s", result2.Time)
		
		// Test 3: Conversion with microsecond precision
		t.Logf("=== Test 3: Microsecond precision preservation ===")
		preciseTime := time.Date(2025, 9, 16, 11, 32, 42, 123456789, time.UTC)
		gtimePrecise := gtime.NewFromTime(preciseTime)
		
		var result3 gtime.Time
		err = gconv.Struct(gtimePrecise, &result3)
		t.AssertNil(err)
		t.Assert(result3.Nanosecond(), preciseTime.Nanosecond())
		t.Assert(result3.Equal(gtimePrecise), true)
		t.Logf("Precision preserved - Original: %s, Result: %s", gtimePrecise.Time, result3.Time)
		
		// Test 4: Conversion with different date components
		t.Logf("=== Test 4: Date component preservation ===")
		complexTime := time.Date(2025, 12, 31, 23, 59, 59, 999999999, time.UTC)
		gtimeComplex := gtime.NewFromTime(complexTime)
		
		var result4 gtime.Time
		err = gconv.Struct(gtimeComplex, &result4)
		t.AssertNil(err)
		
		t.Assert(result4.Year(), complexTime.Year())
		t.Assert(result4.Month(), complexTime.Month())
		t.Assert(result4.Day(), complexTime.Day())
		t.Assert(result4.Hour(), complexTime.Hour())
		t.Assert(result4.Minute(), complexTime.Minute())
		t.Assert(result4.Second(), complexTime.Second())
		t.Assert(result4.Nanosecond(), complexTime.Nanosecond())
		t.Logf("Complex time preserved - Original: %s, Result: %s", gtimeComplex.Time, result4.Time)
	})
}

// TestBuiltinGTimeConverter_StructFields tests gtime fields in struct conversion
func TestBuiltinGTimeConverter_StructFields(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Set up timezone environment
		originalLocation := time.Local
		defer func() {
			time.Local = originalLocation
		}()
		
		tokyoLocation, _ := time.LoadLocation("Asia/Tokyo")
		time.Local = tokyoLocation
		
		// Test data
		utcTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
		gtimeUTC := gtime.NewFromTime(utcTime)
		
		// Test struct with gtime.Time field
		type TestStructGTime struct {
			ID        int        `json:"id"`
			CreatedAt gtime.Time `json:"created_at"`
			UpdatedAt *gtime.Time `json:"updated_at"`
		}
		
		// Test 1: Map to struct with gtime fields
		t.Logf("=== Test 1: Map to struct with gtime fields ===")
		mapData := map[string]interface{}{
			"id":         1,
			"created_at": gtimeUTC,
			"updated_at": gtimeUTC,
		}
		
		var result1 TestStructGTime
		err := gconv.Struct(mapData, &result1)
		t.AssertNil(err)
		
		t.Assert(result1.ID, 1)
		t.Assert(result1.CreatedAt.Equal(gtimeUTC), true)
		t.AssertNE(result1.UpdatedAt, nil)
		t.Assert(result1.UpdatedAt.Equal(gtimeUTC), true)
		
		// Verify timezone preservation
		_, originalOffset := gtimeUTC.Zone()
		_, createdOffset := result1.CreatedAt.Zone()
		_, updatedOffset := result1.UpdatedAt.Zone()
		
		t.Assert(createdOffset, originalOffset)
		t.Assert(updatedOffset, originalOffset)
		
		t.Logf("Original: %s (offset: %d)", gtimeUTC.Time, originalOffset)
		t.Logf("CreatedAt: %s (offset: %d)", result1.CreatedAt.Time, createdOffset)
		t.Logf("UpdatedAt: %s (offset: %d)", result1.UpdatedAt.Time, updatedOffset)
		
		// Test 2: Struct to struct conversion
		t.Logf("=== Test 2: Struct to struct conversion ===")
		sourceStruct := TestStructGTime{
			ID:        2,
			CreatedAt: *gtimeUTC,
			UpdatedAt: gtimeUTC,
		}
		
		var result2 TestStructGTime
		err = gconv.Struct(sourceStruct, &result2)
		t.AssertNil(err)
		
		t.Assert(result2.ID, 2)
		t.Assert(result2.CreatedAt.Equal(&sourceStruct.CreatedAt), true)
		t.Assert(result2.UpdatedAt.Equal(sourceStruct.UpdatedAt), true)
		
		t.Logf("Struct to struct conversion successful")
	})
}

// TestBuiltinGTimeConverter_SliceConversion tests slice conversion scenarios
func TestBuiltinGTimeConverter_SliceConversion(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Set up timezone environment
		originalLocation := time.Local
		defer func() {
			time.Local = originalLocation
		}()
		
		berlinLocation, _ := time.LoadLocation("Europe/Berlin")
		time.Local = berlinLocation
		
		// Test data with different timezones
		utcTime1 := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
		utcTime2 := time.Date(2025, 9, 16, 15, 45, 30, 123456000, time.UTC)
		gtimeUTC1 := gtime.NewFromTime(utcTime1)
		gtimeUTC2 := gtime.NewFromTime(utcTime2)
		
		// Test 1: Slice of maps to slice of gtime.Time
		t.Logf("=== Test 1: Slice of maps to slice of gtime.Time ===")
		mapSlice := []map[string]interface{}{
			{"time": gtimeUTC1},
			{"time": gtimeUTC2},
		}
		
		var result1 []gtime.Time
		err := gconv.Structs(mapSlice, &result1)
		t.AssertNil(err)
		t.Assert(len(result1), 2)
		
		// Verify timezone preservation for each element
		for i, result := range result1 {
			expected := []*gtime.Time{gtimeUTC1, gtimeUTC2}[i]
			_, expectedOffset := expected.Zone()
			_, resultOffset := result.Zone()
			
			t.Assert(resultOffset, expectedOffset)
			t.Assert(result.Equal(expected), true)
			t.Logf("Element %d - Expected: %s (offset: %d), Result: %s (offset: %d)",
				i, expected.Time, expectedOffset, result.Time, resultOffset)
		}
		
		// Test 2: Slice of maps to slice of *gtime.Time
		t.Logf("=== Test 2: Slice of maps to slice of *gtime.Time ===")
		var result2 []*gtime.Time
		err = gconv.Structs(mapSlice, &result2)
		t.AssertNil(err)
		t.Assert(len(result2), 2)
		
		for i, result := range result2 {
			t.AssertNE(result, nil)
			expected := []*gtime.Time{gtimeUTC1, gtimeUTC2}[i]
			_, expectedOffset := expected.Zone()
			_, resultOffset := result.Zone()
			
			t.Assert(resultOffset, expectedOffset)
			t.Assert(result.Equal(expected), true)
			t.Logf("Pointer Element %d - Expected: %s (offset: %d), Result: %s (offset: %d)",
				i, expected.Time, expectedOffset, result.Time, resultOffset)
		}
		
		// Test 3: Direct gtime slice conversion
		t.Logf("=== Test 3: Direct gtime slice conversion ===")
		gtimeSlice := []interface{}{*gtimeUTC1, gtimeUTC2}
		
		var result3 []gtime.Time
		err = gconv.Structs(gtimeSlice, &result3)
		t.AssertNil(err)
		t.Assert(len(result3), 2)
		
		for i, result := range result3 {
			expected := []*gtime.Time{gtimeUTC1, gtimeUTC2}[i]
			t.Assert(result.Equal(expected), true)
			t.Logf("Direct Element %d preserved timezone correctly", i)
		}
	})
}