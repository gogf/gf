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

// TestGTimeTimezonePreservation_ComprehensiveScenarios tests various timezone preservation scenarios
func TestGTimeTimezonePreservation_ComprehensiveScenarios(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Set up test environment with local timezone different from UTC
		originalLocation := time.Local
		defer func() {
			time.Local = originalLocation
		}()

		// Use a timezone that's different from UTC to catch timezone loss issues
		sydneyLocation, _ := time.LoadLocation("Australia/Sydney")
		time.Local = sydneyLocation

		// Test scenarios with different timezones
		testTimezones := []struct {
			name     string
			location *time.Location
		}{
			{"UTC", time.UTC},
			{"GMT", mustLoadLocationComprehensive("GMT")},
			{"EST", mustLoadLocationComprehensive("America/New_York")},
			{"PST", mustLoadLocationComprehensive("America/Los_Angeles")},
			{"JST", mustLoadLocationComprehensive("Asia/Tokyo")},
			{"CET", mustLoadLocationComprehensive("Europe/Paris")},
			{"IST", mustLoadLocationComprehensive("Asia/Kolkata")},
		}

		baseTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)

		for _, tz := range testTimezones {
			t.Logf("=== Testing timezone: %s ===", tz.name)

			// Create time in specific timezone
			testTime := baseTime.In(tz.location)
			gtimeVal := gtime.NewFromTime(testTime)

			originalName, originalOffset := gtimeVal.Zone()
			t.Logf("Original %s time: %s (zone: %s, offset: %d hours)",
				tz.name, gtimeVal.Time, originalName, originalOffset/3600)

			// Test 1: Direct conversion
			convertedTime := gconv.Time(gtimeVal)
			_, convertedOffset := convertedTime.Zone()
			t.Assert(convertedOffset, originalOffset)
			t.Assert(gtimeVal.Time.Equal(convertedTime), true)

			// Test 2: GTime conversion
			reconvertedGTime := gconv.GTime(gtimeVal)
			t.AssertNE(reconvertedGTime, nil)
			_, reconvertedOffset := reconvertedGTime.Zone()
			t.Assert(reconvertedOffset, originalOffset)
			t.Assert(gtimeVal.Equal(reconvertedGTime), true)

			// Test 3: Struct conversion
			type TimeStruct struct {
				Time time.Time `json:"time"`
			}

			var timeStruct TimeStruct
			err := gconv.Struct(map[string]interface{}{"Time": gtimeVal}, &timeStruct)
			t.AssertNil(err)
			_, structOffset := timeStruct.Time.Zone()
			t.Assert(structOffset, originalOffset)
			t.Assert(gtimeVal.Time.Equal(timeStruct.Time), true)

			// Test 4: Structs (slice) conversion
			result := []map[string]interface{}{{"time": gtimeVal}}
			var timeSlice []time.Time
			err = gconv.Structs(result, &timeSlice)
			t.AssertNil(err)
			t.Assert(len(timeSlice), 1)
			_, sliceOffset := timeSlice[0].Zone()
			t.Assert(sliceOffset, originalOffset)
			t.Assert(gtimeVal.Time.Equal(timeSlice[0]), true)

			t.Logf("%s timezone preservation: ✅ PASSED", tz.name)
		}
	})
}

// TestGTimeTimezonePreservation_DatabaseSimulation simulates database timestamp scenarios
func TestGTimeTimezonePreservation_DatabaseSimulation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Simulate application running in Asia/Shanghai
		originalLocation := time.Local
		defer func() {
			time.Local = originalLocation
		}()

		shanghaiLocation, _ := time.LoadLocation("Asia/Shanghai")
		time.Local = shanghaiLocation

		// Simulate different database storage scenarios
		testCases := []struct {
			name        string
			description string
			dbTime      time.Time
			expectedTz  string
		}{
			{
				name:        "UTC_Storage",
				description: "Database stores timestamp in UTC",
				dbTime:      time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC),
				expectedTz:  "UTC",
			},
			{
				name:        "GMT_Storage",
				description: "Database stores timestamp in GMT",
				dbTime:      time.Date(2025, 9, 16, 11, 32, 42, 878465000, mustLoadLocationComprehensive("GMT")),
				expectedTz:  "GMT",
			},
			{
				name:        "Server_Timezone",
				description: "Database timestamp in server timezone",
				dbTime:      time.Date(2025, 9, 16, 19, 32, 42, 878465000, shanghaiLocation),
				expectedTz:  "Asia/Shanghai",
			},
		}

		for _, tc := range testCases {
			t.Logf("=== %s: %s ===", tc.name, tc.description)

			// Create gtime from database time (simulating ORM behavior)
			gtimeFromDB := gtime.NewFromTime(tc.dbTime)
			originalName, originalOffset := gtimeFromDB.Zone()

			t.Logf("Database time: %s (zone: %s, offset: %d)",
				gtimeFromDB.Time, originalName, originalOffset/3600)

			// Simulate ORM query result conversion - the critical path that was failing
			dbResult := []map[string]interface{}{
				{"created_at": gtimeFromDB},
				{"updated_at": gtimeFromDB},
			}

			// Convert to time.Time slice (common ORM usage pattern)
			var timestamps []time.Time
			err := gconv.Structs(dbResult, &timestamps)
			t.AssertNil(err)
			t.Assert(len(timestamps), 2)

			// Verify timezone preservation for both timestamps
			for i, ts := range timestamps {
				_, resultOffset := ts.Zone()
				t.Assert(resultOffset, originalOffset)
				t.Assert(gtimeFromDB.Time.Equal(ts), true)

				t.Logf("Element %d: %s (offset: %d) - ✅ PRESERVED",
					i, ts, resultOffset/3600)
			}

			// Also test struct field conversion
			type DatabaseRecord struct {
				ID        int       `json:"id"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
			}

			recordData := map[string]interface{}{
				"id":         1,
				"created_at": gtimeFromDB,
				"updated_at": gtimeFromDB,
			}

			var record DatabaseRecord
			err = gconv.Struct(recordData, &record)
			t.AssertNil(err)

			_, createdOffset := record.CreatedAt.Zone()
			_, updatedOffset := record.UpdatedAt.Zone()

			t.Assert(createdOffset, originalOffset)
			t.Assert(updatedOffset, originalOffset)
			t.Assert(gtimeFromDB.Time.Equal(record.CreatedAt), true)
			t.Assert(gtimeFromDB.Time.Equal(record.UpdatedAt), true)

			t.Logf("Struct fields: CreatedAt=%s (offset: %d), UpdatedAt=%s (offset: %d) - ✅ PRESERVED",
				record.CreatedAt, createdOffset/3600, record.UpdatedAt, updatedOffset/3600)
		}
	})
}

// TestGTimeTimezonePreservation_PrecisionAndEdgeCases tests precision and edge cases
func TestGTimeTimezonePreservation_PrecisionAndEdgeCases(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Set up test environment
		originalLocation := time.Local
		defer func() {
			time.Local = originalLocation
		}()

		// Use London timezone (has DST transitions)
		londonLocation, _ := time.LoadLocation("Europe/London")
		time.Local = londonLocation

		// Test precision preservation
		t.Logf("=== Precision Preservation Tests ===")

		precisionTests := []struct {
			name  string
			nanos int
		}{
			{"Microseconds", 123456000},
			{"Nanoseconds", 123456789},
			{"Milliseconds", 123000000},
			{"Zero_Nanos", 0},
			{"Max_Nanos", 999999999},
		}

		for _, pt := range precisionTests {
			t.Logf("--- Testing %s precision ---", pt.name)

			testTime := time.Date(2025, 9, 16, 11, 32, 42, pt.nanos, time.UTC)
			gtimeVal := gtime.NewFromTime(testTime)

			// Test through Structs conversion (the problematic path)
			result := []map[string]interface{}{{"time": gtimeVal}}
			var timeSlice []time.Time
			err := gconv.Structs(result, &timeSlice)
			t.AssertNil(err)
			t.Assert(len(timeSlice), 1)

			convertedTime := timeSlice[0]
			t.Assert(convertedTime.Nanosecond(), pt.nanos)
			t.Assert(convertedTime.Equal(testTime), true)

			t.Logf("%s: Original=%d ns, Converted=%d ns - ✅ PRESERVED",
				pt.name, pt.nanos, convertedTime.Nanosecond())
		}

		// Test edge cases
		t.Logf("=== Edge Cases Tests ===")

		// Test 1: Leap year
		leapTime := time.Date(2024, 2, 29, 11, 32, 42, 0, time.UTC)
		gtimeLeap := gtime.NewFromTime(leapTime)

		var leapResult []time.Time
		err := gconv.Structs([]map[string]interface{}{{"time": gtimeLeap}}, &leapResult)
		t.AssertNil(err)
		t.Assert(leapResult[0].Equal(leapTime), true)
		t.Logf("Leap year: %s - ✅ PRESERVED", leapResult[0])

		// Test 2: Year boundaries
		yearBoundary := time.Date(1999, 12, 31, 23, 59, 59, 999999999, time.UTC)
		gtimeYear := gtime.NewFromTime(yearBoundary)

		var yearResult []time.Time
		err = gconv.Structs([]map[string]interface{}{{"time": gtimeYear}}, &yearResult)
		t.AssertNil(err)
		t.Assert(yearResult[0].Equal(yearBoundary), true)
		t.Logf("Year boundary: %s - ✅ PRESERVED", yearResult[0])

		// Test 3: Unix epoch
		epochTime := time.Unix(0, 0).UTC()
		gtimeEpoch := gtime.NewFromTime(epochTime)

		var epochResult []time.Time
		err = gconv.Structs([]map[string]interface{}{{"time": gtimeEpoch}}, &epochResult)
		t.AssertNil(err)
		t.Assert(epochResult[0].Equal(epochTime), true)
		t.Logf("Unix epoch: %s - ✅ PRESERVED", epochResult[0])

		// Test 4: Future date
		futureTime := time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)
		gtimeFuture := gtime.NewFromTime(futureTime)

		var futureResult []time.Time
		err = gconv.Structs([]map[string]interface{}{{"time": gtimeFuture}}, &futureResult)
		t.AssertNil(err)
		t.Assert(futureResult[0].Equal(futureTime), true)
		t.Logf("Future date: %s - ✅ PRESERVED", futureResult[0])
	})
}

// TestGTimeTimezonePreservation_PerformanceRegression tests performance regression
func TestGTimeTimezonePreservation_PerformanceRegression(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create test data
		utcTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
		gtimeVal := gtime.NewFromTime(utcTime)

		// Performance test: Ensure timezone preservation doesn't significantly impact performance
		iterations := 1000

		// Test 1: Direct conversion performance
		start := time.Now()
		for i := 0; i < iterations; i++ {
			_ = gconv.Time(gtimeVal)
		}
		directDuration := time.Since(start)

		// Test 2: Struct conversion performance
		start = time.Now()
		for i := 0; i < iterations; i++ {
			var result time.Time
			_ = gconv.Struct(gtimeVal, &result)
		}
		structDuration := time.Since(start)

		// Test 3: Structs (slice) conversion performance
		mapData := []map[string]interface{}{{"time": gtimeVal}}
		start = time.Now()
		for i := 0; i < iterations; i++ {
			var result []time.Time
			_ = gconv.Structs(mapData, &result)
		}
		sliceDuration := time.Since(start)

		// Performance should be reasonable (not exact assertions, just reasonable bounds)
		t.Logf("Performance Results for %d iterations:", iterations)
		t.Logf("Direct conversion: %v (avg: %v/op)", directDuration, directDuration/time.Duration(iterations))
		t.Logf("Struct conversion: %v (avg: %v/op)", structDuration, structDuration/time.Duration(iterations))
		t.Logf("Slice conversion: %v (avg: %v/op)", sliceDuration, sliceDuration/time.Duration(iterations))

		// Ensure performance is reasonable (under 1ms per operation)
		avgDirect := directDuration / time.Duration(iterations)
		avgStruct := structDuration / time.Duration(iterations)
		avgSlice := sliceDuration / time.Duration(iterations)

		t.Assert(avgDirect < time.Millisecond, true)
		t.Assert(avgStruct < time.Millisecond, true)
		t.Assert(avgSlice < time.Millisecond, true)

		t.Logf("All performance tests passed ✅")
	})
}

// Helper function to load location for comprehensive tests
func mustLoadLocationComprehensive(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
