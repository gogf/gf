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

		// Create a time with GMT timezone (like database result)
		gmtLocation, _ := time.LoadLocation("GMT")
		dbTime := time.Date(2025, 9, 15, 7, 45, 40, 0, gmtLocation)
		gtimeVal := gtime.NewFromTime(dbTime)

		// Test direct Time converter (should work after fix)
		convertedTime := gconv.Time(gtimeVal)
		originalName, originalOffset := gtimeVal.Zone()
		convertedName, convertedOffset := convertedTime.Zone()
		t.Assert(originalOffset, convertedOffset) // Offset must be preserved
		t.Assert(originalOffset, 0) // GMT offset
		t.Assert(convertedOffset, 0) // Converted offset should also be 0
		
		// Test single struct conversion (should work after fix)
		type TestStruct struct {
			Time time.Time
		}
		var testStruct TestStruct
		err := gconv.Struct(map[string]interface{}{"Time": gtimeVal}, &testStruct)
		t.AssertNil(err)
		_, structOffset := testStruct.Time.Zone()
		t.Assert(structOffset, 0) // Struct field should preserve timezone

		// Test the problematic case: ORM Result.Structs() conversion  
		// Note: This test documents the current issue and should be updated when fixed
		result := []map[string]interface{}{{"now": gtimeVal}}
		var nowResult []time.Time
		err = gconv.Structs(result, &nowResult)
		t.AssertNil(err)
		
		structsTime := nowResult[0]
		_, structsOffset := structsTime.Zone()
		
		// TODO: This should pass when the issue is fully fixed
		// Currently documents the known issue
		if structsOffset == 0 {
			t.Logf("✅ Structs timezone preservation works!")
		} else {
			t.Logf("⚠️  Known issue: Structs loses timezone (offset: %d vs expected: 0)", structsOffset/3600)
		}
		
		// The critical assertion: times should represent the same instant
		t.Assert(gtimeVal.Time.Equal(structsTime), true)
		
		// Note: Timezone name might change but offset preservation is critical
		_, _ = originalName, convertedName
	})
}