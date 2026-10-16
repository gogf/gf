// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_formatDecimal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			text   string
			scale  int64
			expect string
		}{
			{".5", 2, "0.50"},
			{"-.5", 2, "-0.50"},
			{"-.25", 2, "-0.25"},
			{"1500", 2, "1500.00"},
			{"-1500", 2, "-1500.00"},
			{"0", 2, "0.00"},
			{"0", 0, "0"},
			{"12", 0, "12"},
			{"1.25", 0, "1.25"},
			{"69.965", 0, "69.965"},
			{"99999999.99", 2, "99999999.99"},
			{".000000000000000000000000000001", 30, "0.000000000000000000000000000001"},
			{"12345678.12345678901234567890123456789", 30, "12345678.123456789012345678901234567890"},
		}
		for _, c := range cases {
			t.Assert(formatDecimal(c.text, c.scale), c.expect)
		}
	})
}

func Test_parseDecimal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, text := range []string{"12345678.123456789012345678901234567890", "-0.5", "0", " 1500 ", "1e3"} {
			decimal, ok := parseDecimal(text)
			t.Assert(ok, true)
			t.Assert(string(decimal), text)
			_, err := decimal.Value()
			t.AssertNil(err)
		}
		for _, text := range []string{"", " ", "abc", "1.2.3", "Inf", "-Inf", "NaN", "0x10"} {
			_, ok := parseDecimal(text)
			t.Assert(ok, false)
		}
	})
}

func Test_releaseSavePointRegex(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, sql := range []string{
			"RELEASE SAVEPOINT transaction1",
			"release savepoint transaction1",
			` RELEASE  SAVEPOINT "point" `,
		} {
			t.Assert(releaseSavePointRegex.MatchString(sql), true)
		}
		for _, sql := range []string{
			"SAVEPOINT transaction1",
			"ROLLBACK TO SAVEPOINT transaction1",
			"SELECT 'RELEASE SAVEPOINT a' FROM DUAL",
			"RELEASE SAVEPOINT a; DROP TABLE t",
		} {
			t.Assert(releaseSavePointRegex.MatchString(sql), false)
		}
	})
}
