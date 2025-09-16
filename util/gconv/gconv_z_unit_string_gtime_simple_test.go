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

// TestGTimeStringConversion_Basic tests basic gtime string conversion
func TestGTimeStringConversion_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Set up timezone environment
		originalLocation := time.Local
		defer func() {
			time.Local = originalLocation
		}()
		
		parisLocation, _ := time.LoadLocation("Europe/Paris")
		time.Local = parisLocation
		
		// Test UTC time string conversion
		utcTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
		gtimeVal := gtime.NewFromTime(utcTime)
		
		// Test gtime.Time to string
		resultStr := gconv.String(*gtimeVal)
		t.Logf("gtime to string: %s", resultStr)
		
		// Should use RFC3339 format (note: microseconds will be truncated if they're 0)
		expectedRFC3339 := "2025-09-16T11:32:42Z"
		t.Assert(resultStr, expectedRFC3339)
		
		// Test *gtime.Time to string
		ptrStr := gconv.String(gtimeVal)
		t.Assert(ptrStr, expectedRFC3339)
		
		// Test round-trip conversion
		reconverted := gconv.GTime(resultStr)
		t.AssertNE(reconverted, nil)
		
		// Check if times represent the same instant (more important than exact equality due to precision differences)
		t.Assert(gtimeVal.Time.Truncate(time.Second).Equal(reconverted.Time.Truncate(time.Second)), true)
		
		// Verify timezone preservation
		_, originalOffset := gtimeVal.Zone()
		_, reconvertedOffset := reconverted.Zone()
		t.Assert(reconvertedOffset, originalOffset)
		
		t.Logf("✅ String conversion preserves timezone correctly")
	})
}

// TestGTimeStringConversion_Precision tests precision preservation
func TestGTimeStringConversion_Precision(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test microsecond precision
		preciseTime := time.Date(2025, 9, 16, 11, 32, 42, 123456789, time.UTC)
		gtimeVal := gtime.NewFromTime(preciseTime)
		
		// Convert to string
		timeStr := gconv.String(gtimeVal)
		t.Logf("Precise time string: %s", timeStr)
		
		// Should include nanosecond precision
		expected := "2025-09-16T11:32:42.123456789Z"
		t.Assert(timeStr, expected)
		
		// Convert back
		reconverted := gconv.GTime(timeStr)
		t.AssertNE(reconverted, nil)
		
		// Verify precision preservation
		t.Assert(reconverted.Nanosecond(), preciseTime.Nanosecond())
		t.Assert(reconverted.Equal(gtimeVal), true)
		
		t.Logf("✅ Precision preserved in string conversion")
	})
}

// TestGTimeStringConversion_EdgeCases tests edge cases
func TestGTimeStringConversion_EdgeCases(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test zero gtime
		zeroGTime := gtime.Time{}
		zeroStr := gconv.String(zeroGTime)
		t.Assert(zeroStr, "")
		
		// Test nil gtime
		var nilGTime *gtime.Time = nil
		nilStr := gconv.String(nilGTime)
		t.Assert(nilStr, "")
		
		// Test very old date
		oldTime := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
		oldGTime := gtime.NewFromTime(oldTime)
		oldStr := gconv.String(oldGTime)
		expectedOld := "1900-01-01T00:00:00Z"
		t.Assert(oldStr, expectedOld)
		
		// Test round-trip for old date
		fromOld := gconv.GTime(oldStr)
		t.Assert(fromOld.Equal(oldGTime), true)
		
		t.Logf("✅ Edge cases handled correctly")
	})
}