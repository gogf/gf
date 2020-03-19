// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
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
	T *testing.T
}

// Assert checks <value> and <expect> EQUAL.
func (t *T) Assert(value, expect interface{}) {
	Assert(value, expect)
}

// AssertEQ checks <value> and <expect> EQUAL, including their TYPES.
func (t *T) AssertEQ(value, expect interface{}) {
	AssertEQ(value, expect)
}

// AssertNE checks <value> and <expect> NOT EQUAL.
func (t *T) AssertNE(value, expect interface{}) {
	AssertNE(value, expect)
}

// AssertGT checks <value> is GREATER THAN <expect>.
// Notice that, only string, integer and float types can be compared by AssertGT,
// others are invalid.
func (t *T) AssertGT(value, expect interface{}) {
	AssertGT(value, expect)
}

// AssertGE checks <value> is GREATER OR EQUAL THAN <expect>.
// Notice that, only string, integer and float types can be compared by AssertGTE,
// others are invalid.
func (t *T) AssertGE(value, expect interface{}) {
	AssertGE(value, expect)
}

// AssertLT checks <value> is LESS EQUAL THAN <expect>.
// Notice that, only string, integer and float types can be compared by AssertLT,
// others are invalid.
func (t *T) AssertLT(value, expect interface{}) {
	AssertLT(value, expect)
}

// AssertLE checks <value> is LESS OR EQUAL THAN <expect>.
// Notice that, only string, integer and float types can be compared by AssertLTE,
// others are invalid.
func (t *T) AssertLE(value, expect interface{}) {
	AssertLE(value, expect)
}

// AssertIN checks <value> is IN <expect>.
// The <expect> should be a slice,
// but the <value> can be a slice or a basic type variable.
func (t *T) AssertIN(value, expect interface{}) {
	AssertIN(value, expect)
}

// AssertNI checks <value> is NOT IN <expect>.
// The <expect> should be a slice,
// but the <value> can be a slice or a basic type variable.
func (t *T) AssertNI(value, expect interface{}) {
	AssertNI(value, expect)
}

// Error panics with given <message>.
func (t *T) Error(message ...interface{}) {
	Error(message...)
}

// Fatal prints <message> to stderr and exit the process.
func (t *T) Fatal(message ...interface{}) {
	Fatal(message...)
}
