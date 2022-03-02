// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson

import (
	"fmt"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gutil"
)

// Interface returns the json value.
func (j *Json) Interface() interface{} {
	if j == nil {
		return nil
	}
	j.mu.RLock()
	defer j.mu.RUnlock()
	return *(j.p)
}

// Var returns the json value as *gvar.Var.
func (j *Json) Var() *gvar.Var {
	return gvar.New(j.Interface())
}

// IsNil checks whether the value pointed by `j` is nil.
func (j *Json) IsNil() bool {
	if j == nil {
		return true
	}
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.p == nil || *(j.p) == nil
}

// Get retrieves and returns value by specified `pattern`.
// It returns all values of current Json object if `pattern` is given empty or string ".".
// It returns nil if no value found by `pattern`.
//
// We can also access slice item by its index number in `pattern` like:
// "list.10", "array.0.name", "array.0.1.id".
//
// It returns a default value specified by `def` if value for `pattern` is not found.
func (j *Json) Get(pattern string, def ...interface{}) *gvar.Var {
	if j == nil {
		return nil
	}
	j.mu.RLock()
	defer j.mu.RUnlock()

	// It returns nil if pattern is empty.
	if pattern == "" {
		return nil
	}

	result := j.getPointerByPattern(pattern)
	if result != nil {
		return gvar.New(*result)
	}
	if len(def) > 0 {
		return gvar.New(def[0])
	}
	return nil
}

// GetJson gets the value by specified `pattern`,
// and converts it to a un-concurrent-safe Json object.
func (j *Json) GetJson(pattern string, def ...interface{}) *Json {
	return New(j.Get(pattern, def...).Val())
}

// GetJsons gets the value by specified `pattern`,
// and converts it to a slice of un-concurrent-safe Json object.
func (j *Json) GetJsons(pattern string, def ...interface{}) []*Json {
	array := j.Get(pattern, def...).Array()
	if len(array) > 0 {
		jsonSlice := make([]*Json, len(array))
		for i := 0; i < len(array); i++ {
			jsonSlice[i] = New(array[i])
		}
		return jsonSlice
	}
	return nil
}

// GetJsonMap gets the value by specified `pattern`,
// and converts it to a map of un-concurrent-safe Json object.
func (j *Json) GetJsonMap(pattern string, def ...interface{}) map[string]*Json {
	m := j.Get(pattern, def...).Map()
	if len(m) > 0 {
		jsonMap := make(map[string]*Json, len(m))
		for k, v := range m {
			jsonMap[k] = New(v)
		}
		return jsonMap
	}
	return nil
}

// Set sets value with specified `pattern`.
// It supports hierarchical data access by char separator, which is '.' in default.
func (j *Json) Set(pattern string, value interface{}) error {
	return j.setValue(pattern, value, false)
}

// MustSet performs as Set, but it panics if any error occurs.
func (j *Json) MustSet(pattern string, value interface{}) {
	if err := j.Set(pattern, value); err != nil {
		panic(err)
	}
}

// Remove deletes value with specified `pattern`.
// It supports hierarchical data access by char separator, which is '.' in default.
func (j *Json) Remove(pattern string) error {
	return j.setValue(pattern, nil, true)
}

// MustRemove performs as Remove, but it panics if any error occurs.
func (j *Json) MustRemove(pattern string) {
	if err := j.Remove(pattern); err != nil {
		panic(err)
	}
}

// Contains checks whether the value by specified `pattern` exist.
func (j *Json) Contains(pattern string) bool {
	return j.Get(pattern) != nil
}

// Len returns the length/size of the value by specified `pattern`.
// The target value by `pattern` should be type of slice or map.
// It returns -1 if the target value is not found, or its type is invalid.
func (j *Json) Len(pattern string) int {
	p := j.getPointerByPattern(pattern)
	if p != nil {
		switch (*p).(type) {
		case map[string]interface{}:
			return len((*p).(map[string]interface{}))
		case []interface{}:
			return len((*p).([]interface{}))
		default:
			return -1
		}
	}
	return -1
}

// Append appends value to the value by specified `pattern`.
// The target value by `pattern` should be type of slice.
func (j *Json) Append(pattern string, value interface{}) error {
	p := j.getPointerByPattern(pattern)
	if p == nil || *p == nil {
		if pattern == "." {
			return j.Set("0", value)
		}
		return j.Set(fmt.Sprintf("%s.0", pattern), value)
	}
	switch (*p).(type) {
	case []interface{}:
		if pattern == "." {
			return j.Set(fmt.Sprintf("%d", len((*p).([]interface{}))), value)
		}
		return j.Set(fmt.Sprintf("%s.%d", pattern, len((*p).([]interface{}))), value)
	}
	return gerror.NewCodef(gcode.CodeInvalidParameter, "invalid variable type of %s", pattern)
}

// MustAppend performs as Append, but it panics if any error occurs.
func (j *Json) MustAppend(pattern string, value interface{}) {
	if err := j.Append(pattern, value); err != nil {
		panic(err)
	}
}

// Map converts current Json object to map[string]interface{}.
// It returns nil if fails.
func (j *Json) Map() map[string]interface{} {
	return j.Var().Map()
}

// Array converts current Json object to []interface{}.
// It returns nil if fails.
func (j *Json) Array() []interface{} {
	return j.Var().Array()
}

// Scan automatically calls Struct or Structs function according to the type of parameter
// `pointer` to implement the converting.
func (j *Json) Scan(pointer interface{}, mapping ...map[string]string) error {
	return j.Var().Scan(pointer, mapping...)
}

// Dump prints current Json object with more manually readable.
func (j *Json) Dump() {
	if j == nil {
		return
	}
	j.mu.RLock()
	defer j.mu.RUnlock()
	gutil.Dump(*j.p)
}
