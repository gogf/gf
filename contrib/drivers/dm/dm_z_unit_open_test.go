// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// dm_z_unit_open_test.go tests DM DSN construction without requiring a DM server.

package dm

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_BuildDataSourceName_TimezoneOffset verifies IANA timezone offsets in DM DSNs.
func Test_BuildDataSourceName_TimezoneOffset(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			config = &gdb.ConfigNode{
				Host:    "localhost",
				Port:    "5236",
				User:    "user",
				Pass:    "pass",
				Name:    "test",
				Charset: "utf8",
			}
			testCases = []struct {
				timezone string
				now      time.Time
				expected string
			}{
				{
					timezone: "Asia/Shanghai",
					now:      time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
					expected: "dm://user:pass@localhost:5236/test?charset=utf8&schema=test&timeZone=480",
				},
				{
					timezone: "America/New_York",
					now:      time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
					expected: "dm://user:pass@localhost:5236/test?charset=utf8&schema=test&timeZone=-300",
				},
				{
					timezone: "America/New_York",
					now:      time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
					expected: "dm://user:pass@localhost:5236/test?charset=utf8&schema=test&timeZone=-240",
				},
				{
					timezone: "Etc/GMT-12",
					now:      time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
					expected: "dm://user:pass@localhost:5236/test?charset=utf8&schema=test&timeZone=720",
				},
				{
					timezone: "Etc/GMT+12",
					now:      time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
					expected: "dm://user:pass@localhost:5236/test?charset=utf8&schema=test&timeZone=-720",
				},
			}
		)

		for _, testCase := range testCases {
			config.Timezone = testCase.timezone
			source, err := buildDataSourceName(config, testCase.now)
			t.AssertNil(err)
			t.Assert(source, testCase.expected)
			t.Assert(config.Timezone, testCase.timezone)
		}
	})
}

// Test_BuildDataSourceName_OptionalParameters verifies DSN options without timezone or port.
func Test_BuildDataSourceName_OptionalParameters(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		source, err := buildDataSourceName(&gdb.ConfigNode{
			Host:    "localhost",
			User:    "user",
			Pass:    "pass",
			Name:    "test",
			Charset: "utf8",
			Extra:   "compress=1",
		}, time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC))

		t.AssertNil(err)
		t.Assert(source, "dm://user:pass@localhost/test?charset=utf8&schema=test&compress=1")
	})
}

// Test_BuildDataSourceName_InvalidTimezone verifies invalid timezone names fail clearly.
func Test_BuildDataSourceName_InvalidTimezone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := buildDataSourceName(&gdb.ConfigNode{
			Name:     "test",
			Timezone: "Invalid/Timezone",
		}, time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC))

		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
	})
}

// Test_BuildDataSourceName_UnsupportedTimezoneOffset verifies unsupported offsets fail clearly.
func Test_BuildDataSourceName_UnsupportedTimezoneOffset(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := buildDataSourceName(&gdb.ConfigNode{
			Name:     "test",
			Timezone: "Pacific/Kiritimati",
		}, time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC))

		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
	})
}
