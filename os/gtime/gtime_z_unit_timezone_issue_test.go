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

		// Simulate ORM Result.Structs() conversion
		result := []map[string]interface{}{{"now": gtimeVal}}
		var nowResult []time.Time
		err := gconv.Structs(result, &nowResult)
		t.AssertNil(err)
		
		convertedTime := nowResult[0]
		
		// The key assertion: timezone offset should be preserved
		originalName, originalOffset := gtimeVal.Zone()
		convertedName, convertedOffset := convertedTime.Zone()
		
		// Offset must be preserved (this is the critical fix)
		t.Assert(originalOffset, convertedOffset)
		
		// Times should represent the same instant
		t.Assert(gtimeVal.Time.Equal(convertedTime), true)
		
		// Both should have 0 offset (GMT/UTC)
		t.Assert(originalOffset, 0)
		t.Assert(convertedOffset, 0)
		
		// Note: Timezone name might change (GMT->UTC) but that's acceptable as long as offset is preserved
		_ = originalName
		_ = convertedName
	})
}