// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtest

import (
	"fmt"
	"github.com/gogf/gf/internal/empty"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gogf/gf/debug/gdebug"

	"github.com/gogf/gf/util/gconv"
)

const (
	pathFilterKey = "/test/gtest/gtest"
)

// C creates a unit testing case.
// The parameter `t` is the pointer to testing.T of stdlib (*testing.T).
// The parameter `f` is the closure function for unit testing case.
func C(t *testing.T, f func(t *T)) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n%s", err, gdebug.StackWithFilter(pathFilterKey))
			t.Fail()
		}
	}()
	f(&T{t})
}

// Case creates a unit testing case.
// The parameter `t` is the pointer to testing.T of stdlib (*testing.T).
// The parameter `f` is the closure function for unit testing case.
// Deprecated.
func Case(t *testing.T, f func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n%s", err, gdebug.StackWithFilter(pathFilterKey))
			t.Fail()
		}
	}()
	f()
}

// Assert checks `value` and `expect` EQUAL.
func Assert(value, expect interface{}) {
	rvExpect := reflect.ValueOf(expect)
	if empty.IsNil(value) {
		value = nil
	}
	if rvExpect.Kind() == reflect.Map {
		if err := compareMap(value, expect); err != nil {
			panic(err)
		}
		return
	}
	var (
		strValue  = gconv.String(value)
		strExpect = gconv.String(expect)
	)
	if strValue != strExpect {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v == %v`, strValue, strExpect))
	}
}

// AssertEQ checks `value` and `expect` EQUAL, including their TYPES.
func AssertEQ(value, expect interface{}) {
	// Value assert.
	rvExpect := reflect.ValueOf(expect)
	if empty.IsNil(value) {
		value = nil
	}
	if rvExpect.Kind() == reflect.Map {
		if err := compareMap(value, expect); err != nil {
			panic(err)
		}
		return
	}
	strValue := gconv.String(value)
	strExpect := gconv.String(expect)
	if strValue != strExpect {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v == %v`, strValue, strExpect))
	}
	// Type assert.
	t1 := reflect.TypeOf(value)
	t2 := reflect.TypeOf(expect)
	if t1 != t2 {
		panic(fmt.Sprintf(`[ASSERT] EXPECT TYPE %v[%v] == %v[%v]`, strValue, t1, strExpect, t2))
	}
}

// AssertNE checks `value` and `expect` NOT EQUAL.
func AssertNE(value, expect interface{}) {
	rvExpect := reflect.ValueOf(expect)
	if empty.IsNil(value) {
		value = nil
	}
	if rvExpect.Kind() == reflect.Map {
		if err := compareMap(value, expect); err == nil {
			panic(fmt.Sprintf(`[ASSERT] EXPECT %v != %v`, value, expect))
		}
		return
	}
	var (
		strValue  = gconv.String(value)
		strExpect = gconv.String(expect)
	)
	if strValue == strExpect {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v != %v`, strValue, strExpect))
	}
}

// AssertNQ checks `value` and `expect` NOT EQUAL, including their TYPES.
func AssertNQ(value, expect interface{}) {
	// Type assert.
	t1 := reflect.TypeOf(value)
	t2 := reflect.TypeOf(expect)
	if t1 == t2 {
		panic(
			fmt.Sprintf(
				`[ASSERT] EXPECT TYPE %v[%v] != %v[%v]`,
				gconv.String(value), t1, gconv.String(expect), t2,
			),
		)
	}
	// Value assert.
	AssertNE(value, expect)
}

// AssertGT checks `value` is GREATER THAN `expect`.
// Notice that, only string, integer and float types can be compared by AssertGT,
// others are invalid.
func AssertGT(value, expect interface{}) {
	passed := false
	switch reflect.ValueOf(expect).Kind() {
	case reflect.String:
		passed = gconv.String(value) > gconv.String(expect)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		passed = gconv.Int(value) > gconv.Int(expect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		passed = gconv.Uint(value) > gconv.Uint(expect)

	case reflect.Float32, reflect.Float64:
		passed = gconv.Float64(value) > gconv.Float64(expect)
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v > %v`, value, expect))
	}
}

// AssertGE checks `value` is GREATER OR EQUAL THAN `expect`.
// Notice that, only string, integer and float types can be compared by AssertGTE,
// others are invalid.
func AssertGE(value, expect interface{}) {
	passed := false
	switch reflect.ValueOf(expect).Kind() {
	case reflect.String:
		passed = gconv.String(value) >= gconv.String(expect)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		passed = gconv.Int64(value) >= gconv.Int64(expect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		passed = gconv.Uint64(value) >= gconv.Uint64(expect)

	case reflect.Float32, reflect.Float64:
		passed = gconv.Float64(value) >= gconv.Float64(expect)
	}
	if !passed {
		panic(fmt.Sprintf(
			`[ASSERT] EXPECT %v(%v) >= %v(%v)`,
			value, reflect.ValueOf(value).Kind(),
			expect, reflect.ValueOf(expect).Kind(),
		))
	}
}

// AssertLT checks `value` is LESS EQUAL THAN `expect`.
// Notice that, only string, integer and float types can be compared by AssertLT,
// others are invalid.
func AssertLT(value, expect interface{}) {
	passed := false
	switch reflect.ValueOf(expect).Kind() {
	case reflect.String:
		passed = gconv.String(value) < gconv.String(expect)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		passed = gconv.Int(value) < gconv.Int(expect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		passed = gconv.Uint(value) < gconv.Uint(expect)

	case reflect.Float32, reflect.Float64:
		passed = gconv.Float64(value) < gconv.Float64(expect)
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v < %v`, value, expect))
	}
}

// AssertLE checks `value` is LESS OR EQUAL THAN `expect`.
// Notice that, only string, integer and float types can be compared by AssertLTE,
// others are invalid.
func AssertLE(value, expect interface{}) {
	passed := false
	switch reflect.ValueOf(expect).Kind() {
	case reflect.String:
		passed = gconv.String(value) <= gconv.String(expect)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		passed = gconv.Int(value) <= gconv.Int(expect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		passed = gconv.Uint(value) <= gconv.Uint(expect)

	case reflect.Float32, reflect.Float64:
		passed = gconv.Float64(value) <= gconv.Float64(expect)
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v <= %v`, value, expect))
	}
}

// AssertIN checks `value` is IN `expect`.
// The `expect` should be a slice,
// but the `value` can be a slice or a basic type variable.
// TODO map support.
func AssertIN(value, expect interface{}) {
	var (
		passed     = true
		expectKind = reflect.ValueOf(expect).Kind()
	)
	switch expectKind {
	case reflect.Slice, reflect.Array:
		expectSlice := gconv.Strings(expect)
		for _, v1 := range gconv.Strings(value) {
			result := false
			for _, v2 := range expectSlice {
				if v1 == v2 {
					result = true
					break
				}
			}
			if !result {
				passed = false
				break
			}
		}
	default:
		panic(fmt.Sprintf(`[ASSERT] INVALID EXPECT VALUE TYPE: %v`, expectKind))
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v IN %v`, value, expect))
	}
}

// AssertNI checks `value` is NOT IN `expect`.
// The `expect` should be a slice,
// but the `value` can be a slice or a basic type variable.
// TODO map support.
func AssertNI(value, expect interface{}) {
	var (
		passed     = true
		expectKind = reflect.ValueOf(expect).Kind()
	)
	switch expectKind {
	case reflect.Slice, reflect.Array:
		for _, v1 := range gconv.Strings(value) {
			result := true
			for _, v2 := range gconv.Strings(expect) {
				if v1 == v2 {
					result = false
					break
				}
			}
			if !result {
				passed = false
				break
			}
		}
	default:
		panic(fmt.Sprintf(`[ASSERT] INVALID EXPECT VALUE TYPE: %v`, expectKind))
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v NOT IN %v`, value, expect))
	}
}

// Error panics with given `message`.
func Error(message ...interface{}) {
	panic(fmt.Sprintf("[ERROR] %s", fmt.Sprint(message...)))
}

// Fatal prints `message` to stderr and exit the process.
func Fatal(message ...interface{}) {
	fmt.Fprintf(os.Stderr, "[FATAL] %s\n%s", fmt.Sprint(message...), gdebug.StackWithFilter(pathFilterKey))
	os.Exit(1)
}

// compareMap compares two maps, returns nil if they are equal, or else returns error.
func compareMap(value, expect interface{}) error {
	var (
		rvValue  = reflect.ValueOf(value)
		rvExpect = reflect.ValueOf(expect)
	)
	if empty.IsNil(value) {
		value = nil
	}
	if rvExpect.Kind() == reflect.Map {
		if rvValue.Kind() == reflect.Map {
			if rvExpect.Len() == rvValue.Len() {
				// Turn two interface maps to the same type for comparison.
				// Direct use of rvValue.MapIndex(key).Interface() will panic
				// when the key types are inconsistent.
				mValue := make(map[string]string)
				mExpect := make(map[string]string)
				ksValue := rvValue.MapKeys()
				ksExpect := rvExpect.MapKeys()
				for _, key := range ksValue {
					mValue[gconv.String(key.Interface())] = gconv.String(rvValue.MapIndex(key).Interface())
				}
				for _, key := range ksExpect {
					mExpect[gconv.String(key.Interface())] = gconv.String(rvExpect.MapIndex(key).Interface())
				}
				for k, v := range mExpect {
					if v != mValue[k] {
						return fmt.Errorf(`[ASSERT] EXPECT VALUE map["%v"]:%v == map["%v"]:%v`+
							"\nGIVEN : %v\nEXPECT: %v", k, mValue[k], k, v, mValue, mExpect)
					}
				}
			} else {
				return fmt.Errorf(`[ASSERT] EXPECT MAP LENGTH %d == %d`, rvValue.Len(), rvExpect.Len())
			}
		} else {
			return fmt.Errorf(`[ASSERT] EXPECT VALUE TO BE A MAP, BUT GIVEN "%s"`, rvValue.Kind())
		}
	}
	return nil
}

// AssertNil asserts `value` is nil.
func AssertNil(value interface{}) {
	if empty.IsNil(value) {
		return
	}
	if err, ok := value.(error); ok {
		panic(fmt.Sprintf(`%+v`, err))
	}
	AssertNE(value, nil)
}

// TestDataPath retrieves and returns the testdata path of current package,
// which is used for unit testing cases only.
// The optional parameter `names` specifies the sub-folders/sub-files,
// which will be joined with current system separator and returned with the path.
func TestDataPath(names ...string) string {
	_, path, _ := gdebug.CallerWithFilter(pathFilterKey)
	path = filepath.Dir(path) + string(filepath.Separator) + "testdata"
	for _, name := range names {
		path += string(filepath.Separator) + name
	}
	return path
}

// TestDataContent retrieves and returns the file content for specified testdata path of current package
func TestDataContent(names ...string) string {
	path := TestDataPath(names...)
	if path != "" {
		data, err := ioutil.ReadFile(path)
		if err == nil {
			return string(data)
		}
	}
	return ""
}
