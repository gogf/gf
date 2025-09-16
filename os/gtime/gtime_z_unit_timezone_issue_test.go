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

// Test for issue #4429: gtime timezone preservation during struct conversion
func TestTime_Issue4429_TimezonePreservation(t1 *testing.T) {
	gtest.C(t1, func(t *gtest.T) {
		// Set local timezone to simulate the issue environment
		originalLocation := time.Local
		defer func() {
			time.Local = originalLocation
		}()

		shanghaiLocation, _ := time.LoadLocation("Asia/Shanghai")
		time.Local = shanghaiLocation

		// Create a time with GMT timezone (like database result with microseconds)
		// This matches the exact scenario from the user's screenshot
		utcTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
		gtimeVal := gtime.NewFromTime(utcTime)

		// Verify the original has the expected timezone
		originalName, originalOffset := gtimeVal.Zone()
		t.Assert(originalOffset, 0) // UTC/GMT offset
		t.Logf("Original: %s (timezone: %s, offset: %d)", gtimeVal.Time, originalName, originalOffset/3600)

		// Test direct Time converter (should work after fix)
		convertedTime := gconv.Time(gtimeVal)
		convertedName, convertedOffset := convertedTime.Zone()
		t.Assert(originalOffset, convertedOffset) // Offset must be preserved
		t.Assert(convertedOffset, 0)              // Converted offset should also be 0

		// Test single struct conversion (should work after fix)
		type TestStruct struct {
			Time time.Time
		}
		var testStruct TestStruct
		err := gconv.Struct(map[string]interface{}{"Time": gtimeVal}, &testStruct)
		t.AssertNil(err)
		_, structOffset := testStruct.Time.Zone()
		t.Assert(structOffset, 0) // Struct field should preserve timezone

		// Test the main problematic case: ORM Result.Structs() conversion
		// This is the exact scenario from the user's screenshot
		result := []map[string]interface{}{{"now": gtimeVal}}
		var nowResult []time.Time
		err = gconv.Structs(result, &nowResult)
		t.AssertNil(err)

		structsTime := nowResult[0]
		structsName, structsOffset := structsTime.Zone()

		// Log the actual results for debugging
		t.Logf("Structs result: %s (timezone: %s, offset: %d)", structsTime, structsName, structsOffset/3600)

		// This should now work with the enhanced fix
		t.Assert(structsOffset, 0)                       // Timezone offset should be preserved (UTC/GMT = 0)
		t.Assert(gtimeVal.Time.Equal(structsTime), true) // Same instant in time

		// Test that precision is preserved
		t.Assert(structsTime.Nanosecond(), utcTime.Nanosecond()) // Microsecond precision should be preserved

		// Test edge cases for robustness

		// Test empty map
		emptyMapResult := []map[string]interface{}{{}}
		var emptyResult []time.Time
		err = gconv.Structs(emptyMapResult, &emptyResult)
		t.AssertNil(err)
		t.Assert(len(emptyResult), 1)
		t.Assert(emptyResult[0].IsZero(), true)

		// Test nil gtime value
		nilResult := []map[string]interface{}{{"time": (*gtime.Time)(nil)}}
		var nilTimeResult []time.Time
		err = gconv.Structs(nilResult, &nilTimeResult)
		t.AssertNil(err)
		t.Assert(len(nilTimeResult), 1)
		t.Assert(nilTimeResult[0].IsZero(), true)

		// Test with different timezone (not just UTC)
		gmtLocation, _ := time.LoadLocation("GMT")
		gmtTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, gmtLocation)
		gtimeGMT := gtime.NewFromTime(gmtTime)
		
		gmtResult := []map[string]interface{}{{"now": gtimeGMT}}
		var gmtNowResult []time.Time
		err = gconv.Structs(gmtResult, &gmtNowResult)
		t.AssertNil(err)
		
		gmtFinalTime := gmtNowResult[0]
		_, gmtFinalOffset := gmtFinalTime.Zone()
		t.Assert(gmtFinalOffset, 0) // GMT should also be preserved as 0 offset
		t.Assert(gtimeGMT.Time.Equal(gmtFinalTime), true)

		// Note: Timezone name might change but offset preservation is critical
		_, _ = originalName, convertedName
	})
}
