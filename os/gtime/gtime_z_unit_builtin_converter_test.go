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
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

// TestBuiltinGTimeConverter_Issue4429 tests the specific builtin converter fix for issue #4429
func TestBuiltinGTimeConverter_Issue4429(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Set up test environment to match issue scenario
		originalLocation := time.Local
		defer func() {
			time.Local = originalLocation
		}()
		
		// Simulate the issue environment: local timezone is Asia/Shanghai (+8)
		shanghaiLocation, _ := time.LoadLocation("Asia/Shanghai")
		time.Local = shanghaiLocation
		
		// Test data that matches the exact issue scenario
		// Database returns UTC time with microseconds
		utcTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
		gtimeVal := gtime.NewFromTime(utcTime)
		
		originalName, originalOffset := gtimeVal.Zone()
		t.Logf("Original gtimeVal: %s (zone: %s, offset: %d)", 
			gtimeVal.Time, originalName, originalOffset/3600)
		t.Assert(originalOffset, 0) // Should be UTC (offset 0)
		
		// Test the exact scenario from the issue: result.Structs(&nowResult)
		// This simulates the ORM query result conversion
		result := []map[string]interface{}{{"now": gtimeVal}}
		var nowResult []time.Time
		err := gconv.Structs(result, &nowResult)
		t.AssertNil(err)
		t.Assert(len(nowResult), 1)
		
		structsTime := nowResult[0]
		structsName, structsOffset := structsTime.Zone()
		
		t.Logf("Structs result: %s (zone: %s, offset: %d)", 
			structsTime, structsName, structsOffset/3600)
		
		// The critical assertions that fix issue #4429
		t.Assert(structsOffset, 0)
		t.Assert(gtimeVal.Time.Equal(structsTime), true)
		t.Assert(structsTime.Nanosecond(), utcTime.Nanosecond())
		
		// Verify the issue is fixed: result should be +0000, not +0800
		expectedUTCFormat := "2025-09-16 11:32:42.878465 +0000 UTC"
		actualFormat := structsTime.String()
		t.Assert(actualFormat, expectedUTCFormat)
		
		t.Logf("✅ Issue #4429 FIXED: Original +0000 preserved (not converted to +0800)")
	})
}

// TestBuiltinGTimeConverter_DirectAssignment tests direct assignment optimization
func TestBuiltinGTimeConverter_DirectAssignment(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test the enhanced builtin converter's direct assignment feature
		originalLocation := time.Local
		defer func() {
			time.Local = originalLocation
		}()
		
		// Set different local timezone to test independence
		parisLocation, _ := time.LoadLocation("Europe/Paris")
		time.Local = parisLocation
		
		// Test Case 1: gtime.Time to gtime.Time (value to value)
		t.Logf("=== Test Case 1: gtime.Time to gtime.Time ===")
		
		utcTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
		sourceGTime := *gtime.NewFromTime(utcTime)
		
		var targetGTime gtime.Time
		err := gconv.Struct(sourceGTime, &targetGTime)
		t.AssertNil(err)
		
		// Verify direct assignment preserved everything
		t.Assert(targetGTime.Equal(&sourceGTime), true)
		t.Assert(targetGTime.Location().String(), sourceGTime.Location().String())
		t.Assert(targetGTime.Nanosecond(), sourceGTime.Nanosecond())
		
		_, sourceOffset := sourceGTime.Zone()
		_, targetOffset := targetGTime.Zone()
		t.Assert(targetOffset, sourceOffset)
		
		t.Logf("Source: %s, Target: %s - ✅ DIRECT ASSIGNMENT", sourceGTime.Time, targetGTime.Time)
		
		// Test Case 2: *gtime.Time to *gtime.Time (pointer to pointer)
		t.Logf("=== Test Case 2: *gtime.Time to *gtime.Time ===")
		
		sourcePtr := gtime.NewFromTime(utcTime)
		var targetPtr *gtime.Time
		
		err = gconv.Struct(sourcePtr, &targetPtr)
		t.AssertNil(err)
		t.AssertNE(targetPtr, nil)
		
		// Verify pointer assignment
		t.Assert(targetPtr.Equal(sourcePtr), true)
		t.Assert(targetPtr.Location().String(), sourcePtr.Location().String())
		
		t.Logf("Source Ptr: %s, Target Ptr: %s - ✅ DIRECT ASSIGNMENT", sourcePtr.Time, targetPtr.Time)
		
		// Test Case 3: gtime.Time to *gtime.Time (value to pointer)
		t.Logf("=== Test Case 3: gtime.Time to *gtime.Time ===")
		
		var targetFromValue *gtime.Time
		err = gconv.Struct(sourceGTime, &targetFromValue)
		t.AssertNil(err)
		t.AssertNE(targetFromValue, nil)
		
		t.Assert(targetFromValue.Equal(&sourceGTime), true)
		t.Assert(targetFromValue.Location().String(), sourceGTime.Location().String())
		
		t.Logf("Source Value: %s, Target Ptr: %s - ✅ DIRECT ASSIGNMENT", sourceGTime.Time, targetFromValue.Time)
		
		// Test Case 4: *gtime.Time to gtime.Time (pointer to value)
		t.Logf("=== Test Case 4: *gtime.Time to gtime.Time ===")
		
		var targetFromPtr gtime.Time
		err = gconv.Struct(sourcePtr, &targetFromPtr)
		t.AssertNil(err)
		
		t.Assert(targetFromPtr.Equal(sourcePtr), true)
		t.Assert(targetFromPtr.Location().String(), sourcePtr.Location().String())
		
		t.Logf("Source Ptr: %s, Target Value: %s - ✅ DIRECT ASSIGNMENT", sourcePtr.Time, targetFromPtr.Time)
	})
}

// TestBuiltinGTimeConverter_FallbackPaths tests fallback conversion paths
func TestBuiltinGTimeConverter_FallbackPaths(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test scenarios where builtin converter falls back to general conversion
		
		// Test 1: String to gtime.Time (should use general conversion)
		t.Logf("=== Test 1: String to gtime.Time fallback ===")
		
		timeStr := "2025-09-16T11:32:42Z"
		var gtimeFromStr gtime.Time
		err := gconv.Struct(timeStr, &gtimeFromStr)
		t.AssertNil(err)
		
		// Should still preserve timezone from RFC3339 format
		_, offset := gtimeFromStr.Zone()
		t.Assert(offset, 0) // UTC offset from Z suffix
		t.Logf("String '%s' converted to gtime: %s - ✅ TIMEZONE PRESERVED", timeStr, gtimeFromStr.Time)
		
		// Test 2: Integer timestamp to gtime.Time
		t.Logf("=== Test 2: Integer timestamp to gtime.Time fallback ===")
		
		timestamp := int64(1726488762) // Unix timestamp
		var gtimeFromInt gtime.Time
		err = gconv.Struct(timestamp, &gtimeFromInt)
		t.AssertNil(err)
		
		expectedTime := time.Unix(timestamp, 0).UTC()
		t.Assert(gtimeFromInt.Unix(), expectedTime.Unix())
		t.Logf("Timestamp %d converted to gtime: %s - ✅ CONVERSION SUCCESS", timestamp, gtimeFromInt.Time)
		
		// Test 3: time.Time to gtime.Time (should use general conversion)
		t.Logf("=== Test 3: time.Time to gtime.Time fallback ===")
		
		goTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
		var gtimeFromGoTime gtime.Time
		err = gconv.Struct(goTime, &gtimeFromGoTime)
		t.AssertNil(err)
		
		t.Assert(gtimeFromGoTime.Time.Equal(goTime), true)
		_, gtimeOffset := gtimeFromGoTime.Zone()
		_, goTimeOffset := goTime.Zone()
		t.Assert(gtimeOffset, goTimeOffset)
		t.Logf("time.Time %s converted to gtime: %s - ✅ TIMEZONE PRESERVED", goTime, gtimeFromGoTime.Time)
	})
}

// TestBuiltinGTimeConverter_NilAndZeroHandling tests nil and zero value handling
func TestBuiltinGTimeConverter_NilAndZeroHandling(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test 1: Nil *gtime.Time to gtime.Time
		t.Logf("=== Test 1: Nil *gtime.Time to gtime.Time ===")
		
		var nilGTime *gtime.Time = nil
		var resultGTime gtime.Time
		err := gconv.Struct(nilGTime, &resultGTime)
		t.AssertNil(err)
		t.Assert(resultGTime.IsZero(), true)
		t.Logf("Nil gtime converted to zero gtime: %s", resultGTime.Time)
		
		// Test 2: Nil *gtime.Time to *gtime.Time
		t.Logf("=== Test 2: Nil *gtime.Time to *gtime.Time ===")
		
		var resultPtr *gtime.Time
		err = gconv.Struct(nilGTime, &resultPtr)
		t.AssertNil(err)
		t.AssertNE(resultPtr, nil) // Should create new gtime.Time, not remain nil
		t.Assert(resultPtr.IsZero(), true)
		t.Logf("Nil gtime converted to zero gtime pointer: %s", resultPtr.Time)
		
		// Test 3: Zero gtime.Time to gtime.Time
		t.Logf("=== Test 3: Zero gtime.Time to gtime.Time ===")
		
		zeroGTime := gtime.Time{}
		var resultZero gtime.Time
		err = gconv.Struct(zeroGTime, &resultZero)
		t.AssertNil(err)
		t.Assert(resultZero.IsZero(), true)
		t.Assert(resultZero.Equal(&zeroGTime), true)
		t.Logf("Zero gtime preserved: %s", resultZero.Time)
		
		// Test 4: Zero gtime.Time in struct
		t.Logf("=== Test 4: Zero gtime.Time in struct ===")
		
		type TestStruct struct {
			ZeroTime gtime.Time  `json:"zero_time"`
			NilTime  *gtime.Time `json:"nil_time"`
		}
		
		inputData := map[string]interface{}{
			"zero_time": gtime.Time{},
			"nil_time":  (*gtime.Time)(nil),
		}
		
		var resultStruct TestStruct
		err = gconv.Struct(inputData, &resultStruct)
		t.AssertNil(err)
		
		t.Assert(resultStruct.ZeroTime.IsZero(), true)
		t.AssertNE(resultStruct.NilTime, nil)
		t.Assert(resultStruct.NilTime.IsZero(), true)
		
		t.Logf("Struct with zero/nil times: ZeroTime=%s, NilTime=%s", 
			resultStruct.ZeroTime.Time, resultStruct.NilTime.Time)
	})
}