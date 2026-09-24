// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlite

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_wrapUnionOperandsAsSubQueries(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cases := map[string]string{
			// The shapes the core builds, with and without the trailing ORDER BY / LIMIT.
			"(SELECT * FROM `a`) UNION (SELECT * FROM `b`)":                            "SELECT * FROM (SELECT * FROM `a`) UNION SELECT * FROM (SELECT * FROM `b`)",
			"(SELECT * FROM `a`) UNION ALL (SELECT * FROM `b`) ORDER BY `id` LIMIT 1":  "SELECT * FROM (SELECT * FROM `a`) UNION ALL SELECT * FROM (SELECT * FROM `b`) ORDER BY `id` LIMIT 1",
			"(SELECT 1) UNION (SELECT 2) UNION (SELECT 3)":                             "SELECT * FROM (SELECT 1) UNION SELECT * FROM (SELECT 2) UNION SELECT * FROM (SELECT 3)",
			"(SELECT * FROM `a` ORDER BY `id` DESC) union (SELECT * FROM `b` LIMIT 5)": "SELECT * FROM (SELECT * FROM `a` ORDER BY `id` DESC) union SELECT * FROM (SELECT * FROM `b` LIMIT 5)",
			// Parentheses inside string literals and nested sub-queries do not confuse the matching.
			"(SELECT * FROM `a` WHERE `n`='x)y') UNION (SELECT * FROM `b` WHERE `id` IN (SELECT `id` FROM `c`))": "SELECT * FROM (SELECT * FROM `a` WHERE `n`='x)y') UNION SELECT * FROM (SELECT * FROM `b` WHERE `id` IN (SELECT `id` FROM `c`))",
			// Not a compound query: left untouched.
			"(SELECT 1)":                      "(SELECT 1)",
			"SELECT * FROM `a` WHERE `x`=(1)": "SELECT * FROM `a` WHERE `x`=(1)",
			"INSERT INTO `a`(`x`) VALUES(1)":  "INSERT INTO `a`(`x`) VALUES(1)",
		}
		for in, want := range cases {
			t.Assert(wrapUnionOperandsAsSubQueries(in), want)
		}
	})
}

func Test_unwrapDoubledExistsParentheses(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cases := map[string]string{
			"SELECT * FROM `u` WHERE EXISTS ((SELECT `id` FROM `t` WHERE uid = (u.id)))":              "SELECT * FROM `u` WHERE EXISTS (SELECT `id` FROM `t` WHERE uid = (u.id))",
			"SELECT * FROM `u` WHERE NOT EXISTS ((SELECT 1 FROM `t`)) AND `id`>0":                     "SELECT * FROM `u` WHERE NOT EXISTS (SELECT 1 FROM `t`) AND `id`>0",
			"SELECT * FROM `u` WHERE EXISTS ((SELECT 1)) OR EXISTS ((SELECT 2))":                      "SELECT * FROM `u` WHERE EXISTS (SELECT 1) OR EXISTS (SELECT 2)",
			"SELECT * FROM `u` WHERE exists ((SELECT 1 FROM `t` WHERE `n`=')('))":                     "SELECT * FROM `u` WHERE exists (SELECT 1 FROM `t` WHERE `n`=')(')",
			"SELECT * FROM `u` WHERE EXISTS (SELECT 1 FROM `t` WHERE `x` IN ((SELECT 1),(SELECT 2)))": "SELECT * FROM `u` WHERE EXISTS (SELECT 1 FROM `t` WHERE `x` IN ((SELECT 1),(SELECT 2)))",
		}
		for in, want := range cases {
			t.Assert(unwrapDoubledExistsParentheses(in), want)
		}
	})
}
