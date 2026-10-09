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

func Test_rewriteSubQueries(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cases := map[string]string{
			// The compound shapes the core builds, with and without the trailing ORDER BY / LIMIT.
			"(SELECT * FROM `a`) UNION (SELECT * FROM `b`)":                            "SELECT * FROM (SELECT * FROM `a`) UNION SELECT * FROM (SELECT * FROM `b`)",
			"(SELECT * FROM `a`) UNION ALL (SELECT * FROM `b`) ORDER BY `id` LIMIT 1":  "SELECT * FROM (SELECT * FROM `a`) UNION ALL SELECT * FROM (SELECT * FROM `b`) ORDER BY `id` LIMIT 1",
			"(SELECT 1) UNION (SELECT 2) UNION (SELECT 3)":                             "SELECT * FROM (SELECT 1) UNION SELECT * FROM (SELECT 2) UNION SELECT * FROM (SELECT 3)",
			"(SELECT * FROM `a` ORDER BY `id` DESC) union (SELECT * FROM `b` LIMIT 5)": "SELECT * FROM (SELECT * FROM `a` ORDER BY `id` DESC) union SELECT * FROM (SELECT * FROM `b` LIMIT 5)",
			"(SELECT 1) EXCEPT (SELECT 2) INTERSECT (SELECT 3)":                        "SELECT * FROM (SELECT 1) EXCEPT SELECT * FROM (SELECT 2) INTERSECT SELECT * FROM (SELECT 3)",
			// Parentheses inside string literals and nested sub-queries do not confuse the matching.
			"(SELECT * FROM `a` WHERE `n`='x)y') UNION (SELECT * FROM `b` WHERE `id` IN (SELECT `id` FROM `c`))": "SELECT * FROM (SELECT * FROM `a` WHERE `n`='x)y') UNION SELECT * FROM (SELECT * FROM `b` WHERE `id` IN (SELECT `id` FROM `c`))",
			// A compound query inside another statement: Count, IN and EXISTS on a union model.
			"SELECT COUNT(1) FROM ((SELECT `id` FROM `a`) UNION (SELECT `id` FROM `b`)) AS T":       "SELECT COUNT(1) FROM (SELECT * FROM (SELECT `id` FROM `a`) UNION SELECT * FROM (SELECT `id` FROM `b`)) AS T",
			"SELECT COUNT(1) FROM ((SELECT `id` FROM `a`) UNION ALL (SELECT `id` FROM `b`)) AS T":   "SELECT COUNT(1) FROM (SELECT * FROM (SELECT `id` FROM `a`) UNION ALL SELECT * FROM (SELECT `id` FROM `b`)) AS T",
			"SELECT * FROM `u` WHERE id IN (((SELECT `id` FROM `a`) UNION (SELECT `id` FROM `b`)))": "SELECT * FROM `u` WHERE id IN (SELECT * FROM (SELECT `id` FROM `a`) UNION SELECT * FROM (SELECT `id` FROM `b`))",
			"SELECT * FROM `u` WHERE EXISTS (((SELECT 1) UNION (SELECT 2)))":                        "SELECT * FROM `u` WHERE EXISTS (SELECT * FROM (SELECT 1) UNION SELECT * FROM (SELECT 2))",
			"SELECT * FROM (SELECT * FROM ((SELECT 1) UNION (SELECT 2)) AS x) AS y":                 "SELECT * FROM (SELECT * FROM (SELECT * FROM (SELECT 1) UNION SELECT * FROM (SELECT 2)) AS x) AS y",
			// A sub-query in doubled parentheses keeps one pair.
			"SELECT * FROM `u` WHERE EXISTS ((SELECT `id` FROM `t` WHERE uid = (u.id)))":      "SELECT * FROM `u` WHERE EXISTS (SELECT `id` FROM `t` WHERE uid = (u.id))",
			"SELECT * FROM `u` WHERE NOT EXISTS ((SELECT 1 FROM `t`)) AND `id`>0":             "SELECT * FROM `u` WHERE NOT EXISTS (SELECT 1 FROM `t`) AND `id`>0",
			"SELECT * FROM `u` WHERE EXISTS ((SELECT 1)) OR EXISTS ((SELECT 2))":              "SELECT * FROM `u` WHERE EXISTS (SELECT 1) OR EXISTS (SELECT 2)",
			"SELECT * FROM `u` WHERE exists ((SELECT 1 FROM `t` WHERE `n`=')('))":             "SELECT * FROM `u` WHERE exists (SELECT 1 FROM `t` WHERE `n`=')(')",
			"SELECT * FROM `u` WHERE `id` IN ((WITH c AS (SELECT 1 AS id) SELECT id FROM c))": "SELECT * FROM `u` WHERE `id` IN (WITH c AS (SELECT 1 AS id) SELECT id FROM c)",
			// Left untouched: no compound query, a list of sub-queries, doubled parentheses around a value.
			"(SELECT 1)":                        "(SELECT 1)",
			"SELECT * FROM `a` WHERE `x`=(1)":   "SELECT * FROM `a` WHERE `x`=(1)",
			"SELECT * FROM `a` WHERE `x`=((1))": "SELECT * FROM `a` WHERE `x`=((1))",
			"INSERT INTO `a`(`x`) VALUES(1)":    "INSERT INTO `a`(`x`) VALUES(1)",
			"SELECT * FROM `u` WHERE EXISTS (SELECT 1 FROM `t` WHERE `x` IN ((SELECT 1),(SELECT 2)))": "SELECT * FROM `u` WHERE EXISTS (SELECT 1 FROM `t` WHERE `x` IN ((SELECT 1),(SELECT 2)))",
			"SELECT 'a UNION (b)' FROM `t` -- (SELECT 1) UNION (SELECT 2)":                            "SELECT 'a UNION (b)' FROM `t` -- (SELECT 1) UNION (SELECT 2)",
			"SELECT 'EXISTS ((SELECT 1))' AS note":                                                    "SELECT 'EXISTS ((SELECT 1))' AS note",
			"UPDATE `t` SET `note`='(SELECT 1) UNION (SELECT 2)' WHERE `id`=1":                        "UPDATE `t` SET `note`='(SELECT 1) UNION (SELECT 2)' WHERE `id`=1",
			"SELECT * FROM `a` WHERE `x`=(1":                                                          "SELECT * FROM `a` WHERE `x`=(1",
		}
		for in, want := range cases {
			t.Assert(rewriteSubQueries(in), want)
		}
	})
}
