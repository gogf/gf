// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtest

import (
	"testing"
)

// T is the testing unit case management object.
type T struct {
	*testing.T
}

// Assert checks `value` and `expect` EQUAL.
func (t *T) Assert(value, expect any) {
	Assert(value, expect)
}

// AssertEQ checks `value` and `expect` EQUAL, including their TYPES.
func (t *T) AssertEQ(value, expect any) {
	AssertEQ(value, expect)
}

// AssertNE checks `value` and `expect` NOT EQUAL.
func (t *T) AssertNE(value, expect any) {
	AssertNE(value, expect)
}

// AssertNQ checks `value` and `expect` NOT EQUAL, including their TYPES.
func (t *T) AssertNQ(value, expect any) {
	AssertNQ(value, expect)
}

// AssertGT checks `value` is GREATER THAN `expect`.
// Notice that, only string, integer and float types can be compared by AssertGT,
// others are invalid.
func (t *T) AssertGT(value, expect any) {
	AssertGT(value, expect)
}

// AssertGE checks `value` is GREATER OR EQUAL THAN `expect`.
// Notice that, only string, integer and float types can be compared by AssertGTE,
// others are invalid.
func (t *T) AssertGE(value, expect any) {
	AssertGE(value, expect)
}

// AssertLT checks `value` is LESS EQUAL THAN `expect`.
// Notice that, only string, integer and float types can be compared by AssertLT,
// others are invalid.
func (t *T) AssertLT(value, expect any) {
	AssertLT(value, expect)
}

// AssertLE checks `value` is LESS OR EQUAL THAN `expect`.
// Notice that, only string, integer and float types can be compared by AssertLTE,
// others are invalid.
func (t *T) AssertLE(value, expect any) {
	AssertLE(value, expect)
}

// AssertIN checks `value` is IN `expect`.
// The `expect` should be a slice,
// but the `value` can be a slice or a basic type variable.
func (t *T) AssertIN(value, expect any) {
	AssertIN(value, expect)
}

// AssertNI checks `value` is NOT IN `expect`.
// The `expect` should be a slice,
// but the `value` can be a slice or a basic type variable.
func (t *T) AssertNI(value, expect any) {
	AssertNI(value, expect)
}

// AssertNil asserts `value` is nil.
func (t *T) AssertNil(value any) {
	AssertNil(value)
}

// Error panics with given `message`.
func (t *T) Error(message ...any) {
	Error(message...)
}

// Fatal prints `message` to stderr and exit the process.
func (t *T) Fatal(message ...any) {
	Fatal(message...)
}
