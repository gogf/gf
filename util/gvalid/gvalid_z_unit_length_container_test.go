// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gvalid"
)

// runLengthRule validates a single `value` against a single length related `rule`
// and reports whether it passed.
func runLengthRule(rule string, value any) bool {
	return gvalid.New().
		Rules(map[string]string{"f": rule}).
		Data(map[string]any{"f": value}).
		Run(context.Background()) == nil
}

// Test_LengthRules_Container asserts that the length related rules measure a slice,
// an array and a map by the number of their elements.
//
// Before this behavior existed those rules measured the number of unicode runes of
// the JSON text a container is serialized into, so the verdict depended on how wide
// the elements happened to print: `[]int{1}` passed `size:3` because `[1]` is three
// characters long, while `[]int{1, 2, 3}` failed it because `[1,2,3]` is seven.
func Test_LengthRules_Container(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type item struct{ Id int }
		twoStrings := []string{"a", "b"}

		// The element count is what decides, whatever the elements print like.
		t.Assert(runLengthRule("size:3", []int{1, 2, 3}), true)
		t.Assert(runLengthRule("size:3", []int{1}), false)
		t.Assert(runLengthRule("size:3", []string{"a", "b", "c"}), true)
		t.Assert(runLengthRule("size:3", []uint64{100000, 100001, 100002}), true)
		t.Assert(runLengthRule("size:3", []*item{{1}, {2}, {3}}), true)

		// Wide elements no longer make a short slice look long.
		t.Assert(runLengthRule("max-length:3", []uint64{100000, 100001}), true)
		t.Assert(runLengthRule("max-length:3", []uint64{1, 2, 3, 4}), false)
		t.Assert(runLengthRule("min-length:2", []string{"a"}), false)
		t.Assert(runLengthRule("min-length:2", []string{"a", "b"}), true)
		t.Assert(runLengthRule("length:1,5", []int{1, 2, 3}), true)
		t.Assert(runLengthRule("length:1,5", []int{1, 2, 3, 4, 5, 6}), false)

		// Arrays and maps are measured the same way.
		t.Assert(runLengthRule("size:3", [3]int{1, 2, 3}), true)
		t.Assert(runLengthRule("size:3", map[string]int{"a": 1, "b": 2, "c": 3}), true)
		t.Assert(runLengthRule("size:3", map[string]int{"a": 1}), false)

		// A pointer to a container is dereferenced, as `required` does.
		t.Assert(runLengthRule("size:2", &twoStrings), true)
		t.Assert(runLengthRule("size:3", &twoStrings), false)

		// A rune slice needs no special case: its element count is its rune count.
		t.Assert(runLengthRule("size:3", []rune("你好吗")), true)
	})
}

// Test_LengthRules_Container_AgreesWithRequired asserts that `required` and the length
// rules apply the same notion of length to one and the same value.
//
// `required` reports a container as empty by reflect.Value.Len, see isRequiredEmpty, so
// an empty slice must be rejected by `min-length:1` for the two rules of a tag such as
// `required|min-length:1` to agree.
func Test_LengthRules_Container_AgreesWithRequired(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		empty := []string{}
		t.Assert(runLengthRule("required", empty), false)
		t.Assert(runLengthRule("min-length:1", empty), false)
		t.Assert(runLengthRule("size:0", empty), true)

		one := []string{"a"}
		t.Assert(runLengthRule("required", one), true)
		t.Assert(runLengthRule("min-length:1", one), true)
	})
}

// Test_LengthRules_String_Unchanged pins the behavior that must not move.
//
// Without these the change could not be told apart from an implementation that measures
// every value by reflect.Value.Len, which would silently redefine the length of a string
// as its byte count and that of a byte slice as its number of bytes.
func Test_LengthRules_String_Unchanged(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		text := "abc"

		// A string is still measured in unicode runes, not in bytes.
		t.Assert(runLengthRule("size:3", "abc"), true)
		t.Assert(runLengthRule("size:3", "你好吗"), true)
		t.Assert(runLengthRule("size:9", "你好吗"), false)
		t.Assert(runLengthRule("length:1,5", "abcdefg"), false)
		t.Assert(runLengthRule("size:3", &text), true)

		// A byte slice carries text, so it stays on the string side and is measured
		// in the unicode runes of that text rather than in its number of bytes.
		t.Assert(runLengthRule("size:3", []byte("abc")), true)
		t.Assert(runLengthRule("size:2", []byte("你好")), true)
		t.Assert(runLengthRule("size:6", []byte("你好")), false)
		t.Assert(runLengthRule("length:1,5", []byte("abcdefg")), false)

		// A number keeps being measured by the width of its decimal form.
		t.Assert(runLengthRule("size:3", 123), true)
		t.Assert(runLengthRule("size:1", 123), false)
	})
}
