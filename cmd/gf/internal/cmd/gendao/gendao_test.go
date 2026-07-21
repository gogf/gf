// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

// Test_Issue4659_ShardingDigitSuffix: non-numeric tables must not be treated as shards.
func Test_Issue4659_ShardingDigitSuffix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Pattern a_? with digit-only suffix.
		base, ok := matchShardingPattern("a_?", "a_1")
		t.Assert(ok, true)
		t.Assert(base, "a")

		base, ok = matchShardingPattern("a_?", "a_2")
		t.Assert(ok, true)
		t.Assert(base, "a")

		// Letter / multi-segment tables are not shards of "a".
		_, ok = matchShardingPattern("a_?", "a_b")
		t.Assert(ok, false)
		_, ok = matchShardingPattern("a_?", "a_c")
		t.Assert(ok, false)
		_, ok = matchShardingPattern("a_?", "a_b_1")
		t.Assert(ok, false)

		// Longer pattern still works for a_b_1.
		base, ok = matchShardingPattern("a_b_?", "a_b_1")
		t.Assert(ok, true)
		t.Assert(base, "a_b")

		// users_001 style multi-digit suffixes.
		base, ok = matchShardingPattern("users_?", "users_0001")
		t.Assert(ok, true)
		t.Assert(base, "users")
	})
}

// Test containsWildcard function.
func Test_containsWildcard(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(containsWildcard("trade_*"), true)
		t.Assert(containsWildcard("user_?"), true)
		t.Assert(containsWildcard("*"), true)
		t.Assert(containsWildcard("?"), true)
		t.Assert(containsWildcard("trade_order"), false)
		t.Assert(containsWildcard(""), false)
	})
}

// Test patternToRegex function.
func Test_patternToRegex(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// * should become .*
		t.Assert(patternToRegex("trade_*"), "trade_.*")
		// ? should become .
		t.Assert(patternToRegex("user_???"), "user_...")
		// Mixed
		t.Assert(patternToRegex("*_order_?"), ".*_order_.")
		// No wildcards - should escape special regex chars
		t.Assert(patternToRegex("trade_order"), "trade_order")
		// Just *
		t.Assert(patternToRegex("*"), ".*")
	})
}

// Test filterTablesByPatterns with * wildcard.
func Test_filterTablesByPatterns_Star(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		allTables := []string{"trade_order", "trade_item", "user_info", "user_log", "config"}

		// Single pattern with *
		result := filterTablesByPatterns(allTables, []string{"trade_*"})
		t.Assert(len(result), 2)
		t.AssertIN("trade_order", result)
		t.AssertIN("trade_item", result)

		// Multiple patterns with *
		result = filterTablesByPatterns(allTables, []string{"trade_*", "user_*"})
		t.Assert(len(result), 4)
		t.AssertIN("trade_order", result)
		t.AssertIN("trade_item", result)
		t.AssertIN("user_info", result)
		t.AssertIN("user_log", result)
	})
}

// Test filterTablesByPatterns with ? wildcard.
func Test_filterTablesByPatterns_Question(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		allTables := []string{"trade_order", "trade_item", "user_info", "user_log", "config"}

		// ? matches single character: user_log (3 chars) but not user_info (4 chars)
		result := filterTablesByPatterns(allTables, []string{"user_???"})
		t.Assert(len(result), 1)
		t.AssertIN("user_log", result)
		t.AssertNI("user_info", result)

		// user_???? should match user_info (4 chars)
		result = filterTablesByPatterns(allTables, []string{"user_????"})
		t.Assert(len(result), 1)
		t.AssertIN("user_info", result)
		t.AssertNI("user_log", result)
	})
}

// Test filterTablesByPatterns with mixed patterns and exact names.
func Test_filterTablesByPatterns_Mixed(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		allTables := []string{"trade_order", "trade_item", "user_info", "user_log", "config"}

		// Pattern + exact name
		result := filterTablesByPatterns(allTables, []string{"trade_*", "config"})
		t.Assert(len(result), 3)
		t.AssertIN("trade_order", result)
		t.AssertIN("trade_item", result)
		t.AssertIN("config", result)
		t.AssertNI("user_info", result)
		t.AssertNI("user_log", result)
	})
}

// Test filterTablesByPatterns with exact names only (backward compatibility).
func Test_filterTablesByPatterns_ExactNames(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		allTables := []string{"trade_order", "trade_item", "user_info", "user_log", "config"}

		// Exact names only
		result := filterTablesByPatterns(allTables, []string{"trade_order", "config"})
		t.Assert(len(result), 2)
		t.AssertIN("trade_order", result)
		t.AssertIN("config", result)
		t.AssertNI("trade_item", result)
	})
}

// Test filterTablesByPatterns - no duplicates when table matches multiple patterns.
func Test_filterTablesByPatterns_NoDuplicates(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		allTables := []string{"trade_order", "trade_item", "user_info"}

		// trade_order matches both patterns, should only appear once
		result := filterTablesByPatterns(allTables, []string{"trade_*", "trade_order"})
		t.Assert(len(result), 2) // trade_order, trade_item

		// Count occurrences of trade_order
		count := 0
		for _, v := range result {
			if v == "trade_order" {
				count++
			}
		}
		t.Assert(count, 1) // No duplicates
	})
}

// Test filterTablesByPatterns - pattern matches nothing.
func Test_filterTablesByPatterns_NoMatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		allTables := []string{"trade_order", "trade_item", "user_info"}

		// Pattern that matches nothing
		result := filterTablesByPatterns(allTables, []string{"nonexistent_*"})
		t.Assert(len(result), 0)
	})
}

// Test filterTablesByPatterns - empty input.
func Test_filterTablesByPatterns_Empty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		allTables := []string{"trade_order", "trade_item"}

		// Empty patterns
		result := filterTablesByPatterns(allTables, []string{})
		t.Assert(len(result), 0)

		// Empty tables
		result = filterTablesByPatterns([]string{}, []string{"trade_*"})
		t.Assert(len(result), 0)
	})
}

// Test filterTablesByPatterns - "*" matches all tables.
func Test_filterTablesByPatterns_MatchAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		allTables := []string{"trade_order", "trade_item", "user_info", "user_log", "config"}

		// "*" should match all tables
		result := filterTablesByPatterns(allTables, []string{"*"})
		t.Assert(len(result), 5)
	})
}

// Test filterTablesByPatterns - non-existent exact table name should be skipped.
func Test_filterTablesByPatterns_NonExistent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		allTables := []string{"trade_order", "trade_item", "user_info"}

		// Mix of existing and non-existing tables
		result := filterTablesByPatterns(allTables, []string{"trade_order", "nonexistent", "user_info"})
		t.Assert(len(result), 2)
		t.AssertIN("trade_order", result)
		t.AssertIN("user_info", result)
		t.AssertNI("nonexistent", result)
	})
}

func Test_formatFileName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(formatFileName("sys_i18n_message", ""), "sys_i_18_n_message")
		t.Assert(formatFileName("sys_i18n_message", string(gstr.SnakeFirstUpper)), "sys_i18n_message")
		t.Assert(formatFileName("SYS_I18N_MESSAGE", string(gstr.SnakeFirstUpper)), "sys_i18n_message")
		t.Assert(formatFileName("user_test", string(gstr.SnakeFirstUpper)), "user_test_table")
	})
}
