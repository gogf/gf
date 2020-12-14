// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson

import (
	"fmt"
	"time"

	"github.com/gogf/gf/util/gutil"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

// Value returns the json value.
func (j *Json) Value() interface{} {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return *(j.p)
}

// IsNil checks whether the value pointed by <j> is nil.
func (j *Json) IsNil() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.p == nil || *(j.p) == nil
}

// Get retrieves and returns value by specified <pattern>.
// It returns all values of current Json object if <pattern> is given empty or string ".".
// It returns nil if no value found by <pattern>.
//
// We can also access slice item by its index number in <pattern> like:
// "list.10", "array.0.name", "array.0.1.id".
//
// It returns a default value specified by <def> if value for <pattern> is not found.
func (j *Json) Get(pattern string, def ...interface{}) interface{} {
	j.mu.RLock()
	defer j.mu.RUnlock()

	// It returns nil if pattern is empty.
	if pattern == "" {
		return nil
	}

	// It returns all if pattern is ".".
	if pattern == "." {
		return *j.p
	}

	var result *interface{}
	if j.vc {
		result = j.getPointerByPattern(pattern)
	} else {
		result = j.getPointerByPatternWithoutViolenceCheck(pattern)
	}
	if result != nil {
		return *result
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

// GetVar returns a gvar.Var with value by given <pattern>.
func (j *Json) GetVar(pattern string, def ...interface{}) *gvar.Var {
	return gvar.New(j.Get(pattern, def...))
}

// GetVars returns []*gvar.Var with value by given <pattern>.
func (j *Json) GetVars(pattern string, def ...interface{}) []*gvar.Var {
	return gvar.New(j.Get(pattern, def...)).Vars()
}

// GetMap retrieves and returns the value by specified <pattern> as map[string]interface{}.
func (j *Json) GetMap(pattern string, def ...interface{}) map[string]interface{} {
	result := j.Get(pattern, def...)
	if result != nil {
		return gconv.Map(result)
	}
	return nil
}

// GetMapStrStr retrieves and returns the value by specified <pattern> as map[string]string.
func (j *Json) GetMapStrStr(pattern string, def ...interface{}) map[string]string {
	result := j.Get(pattern, def...)
	if result != nil {
		return gconv.MapStrStr(result)
	}
	return nil
}

// GetMaps retrieves and returns the value by specified <pattern> as []map[string]interface{}.
func (j *Json) GetMaps(pattern string, def ...interface{}) []map[string]interface{} {
	result := j.Get(pattern, def...)
	if result != nil {
		return gconv.Maps(result)
	}
	return nil
}

// GetJson gets the value by specified <pattern>,
// and converts it to a un-concurrent-safe Json object.
func (j *Json) GetJson(pattern string, def ...interface{}) *Json {
	return New(j.Get(pattern, def...))
}

// GetJsons gets the value by specified <pattern>,
// and converts it to a slice of un-concurrent-safe Json object.
func (j *Json) GetJsons(pattern string, def ...interface{}) []*Json {
	array := j.GetArray(pattern, def...)
	if len(array) > 0 {
		jsonSlice := make([]*Json, len(array))
		for i := 0; i < len(array); i++ {
			jsonSlice[i] = New(array[i])
		}
		return jsonSlice
	}
	return nil
}

// GetJsonMap gets the value by specified <pattern>,
// and converts it to a map of un-concurrent-safe Json object.
func (j *Json) GetJsonMap(pattern string, def ...interface{}) map[string]*Json {
	m := j.GetMap(pattern, def...)
	if len(m) > 0 {
		jsonMap := make(map[string]*Json, len(m))
		for k, v := range m {
			jsonMap[k] = New(v)
		}
		return jsonMap
	}
	return nil
}

// GetArray retrieves the value by specified <pattern>,
// and converts it to a slice of []interface{}.
func (j *Json) GetArray(pattern string, def ...interface{}) []interface{} {
	return gconv.Interfaces(j.Get(pattern, def...))
}

// GetString retrieves the value by specified <pattern> and converts it to string.
func (j *Json) GetString(pattern string, def ...interface{}) string {
	return gconv.String(j.Get(pattern, def...))
}

// GetBytes retrieves the value by specified <pattern> and converts it to []byte.
func (j *Json) GetBytes(pattern string, def ...interface{}) []byte {
	return gconv.Bytes(j.Get(pattern, def...))
}

// GetBool retrieves the value by specified <pattern>,
// converts and returns it as bool.
// It returns false when value is: "", 0, false, off, nil;
// or returns true instead.
func (j *Json) GetBool(pattern string, def ...interface{}) bool {
	return gconv.Bool(j.Get(pattern, def...))
}

// GetInt retrieves the value by specified <pattern> and converts it to int.
func (j *Json) GetInt(pattern string, def ...interface{}) int {
	return gconv.Int(j.Get(pattern, def...))
}

// GetInt8 retrieves the value by specified <pattern> and converts it to int8.
func (j *Json) GetInt8(pattern string, def ...interface{}) int8 {
	return gconv.Int8(j.Get(pattern, def...))
}

// GetInt16 retrieves the value by specified <pattern> and converts it to int16.
func (j *Json) GetInt16(pattern string, def ...interface{}) int16 {
	return gconv.Int16(j.Get(pattern, def...))
}

// GetInt32 retrieves the value by specified <pattern> and converts it to int32.
func (j *Json) GetInt32(pattern string, def ...interface{}) int32 {
	return gconv.Int32(j.Get(pattern, def...))
}

// GetInt64 retrieves the value by specified <pattern> and converts it to int64.
func (j *Json) GetInt64(pattern string, def ...interface{}) int64 {
	return gconv.Int64(j.Get(pattern, def...))
}

// GetUint retrieves the value by specified <pattern> and converts it to uint.
func (j *Json) GetUint(pattern string, def ...interface{}) uint {
	return gconv.Uint(j.Get(pattern, def...))
}

// GetUint8 retrieves the value by specified <pattern> and converts it to uint8.
func (j *Json) GetUint8(pattern string, def ...interface{}) uint8 {
	return gconv.Uint8(j.Get(pattern, def...))
}

// GetUint16 retrieves the value by specified <pattern> and converts it to uint16.
func (j *Json) GetUint16(pattern string, def ...interface{}) uint16 {
	return gconv.Uint16(j.Get(pattern, def...))
}

// GetUint32 retrieves the value by specified <pattern> and converts it to uint32.
func (j *Json) GetUint32(pattern string, def ...interface{}) uint32 {
	return gconv.Uint32(j.Get(pattern, def...))
}

// GetUint64 retrieves the value by specified <pattern> and converts it to uint64.
func (j *Json) GetUint64(pattern string, def ...interface{}) uint64 {
	return gconv.Uint64(j.Get(pattern, def...))
}

// GetFloat32 retrieves the value by specified <pattern> and converts it to float32.
func (j *Json) GetFloat32(pattern string, def ...interface{}) float32 {
	return gconv.Float32(j.Get(pattern, def...))
}

// GetFloat64 retrieves the value by specified <pattern> and converts it to float64.
func (j *Json) GetFloat64(pattern string, def ...interface{}) float64 {
	return gconv.Float64(j.Get(pattern, def...))
}

// GetFloats retrieves the value by specified <pattern> and converts it to []float64.
func (j *Json) GetFloats(pattern string, def ...interface{}) []float64 {
	return gconv.Floats(j.Get(pattern, def...))
}

// GetInts retrieves the value by specified <pattern> and converts it to []int.
func (j *Json) GetInts(pattern string, def ...interface{}) []int {
	return gconv.Ints(j.Get(pattern, def...))
}

// GetStrings retrieves the value by specified <pattern> and converts it to []string.
func (j *Json) GetStrings(pattern string, def ...interface{}) []string {
	return gconv.Strings(j.Get(pattern, def...))
}

// GetInterfaces is alias of GetArray.
// See GetArray.
func (j *Json) GetInterfaces(pattern string, def ...interface{}) []interface{} {
	return gconv.Interfaces(j.Get(pattern, def...))
}

// GetTime retrieves the value by specified <pattern> and converts it to time.Time.
func (j *Json) GetTime(pattern string, format ...string) time.Time {
	return gconv.Time(j.Get(pattern), format...)
}

// GetDuration retrieves the value by specified <pattern> and converts it to time.Duration.
func (j *Json) GetDuration(pattern string, def ...interface{}) time.Duration {
	return gconv.Duration(j.Get(pattern, def...))
}

// GetGTime retrieves the value by specified <pattern> and converts it to *gtime.Time.
func (j *Json) GetGTime(pattern string, format ...string) *gtime.Time {
	return gconv.GTime(j.Get(pattern), format...)
}

// Set sets value with specified <pattern>.
// It supports hierarchical data access by char separator, which is '.' in default.
func (j *Json) Set(pattern string, value interface{}) error {
	return j.setValue(pattern, value, false)
}

// Remove deletes value with specified <pattern>.
// It supports hierarchical data access by char separator, which is '.' in default.
func (j *Json) Remove(pattern string) error {
	return j.setValue(pattern, nil, true)
}

// Contains checks whether the value by specified <pattern> exist.
func (j *Json) Contains(pattern string) bool {
	return j.Get(pattern) != nil
}

// Len returns the length/size of the value by specified <pattern>.
// The target value by <pattern> should be type of slice or map.
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

// Append appends value to the value by specified <pattern>.
// The target value by <pattern> should be type of slice.
func (j *Json) Append(pattern string, value interface{}) error {
	p := j.getPointerByPattern(pattern)
	if p == nil {
		return j.Set(fmt.Sprintf("%s.0", pattern), value)
	}
	switch (*p).(type) {
	case []interface{}:
		return j.Set(fmt.Sprintf("%s.%d", pattern, len((*p).([]interface{}))), value)
	}
	return fmt.Errorf("invalid variable type of %s", pattern)
}

// GetStruct retrieves the value by specified <pattern> and converts it to specified object
// <pointer>. The <pointer> should be the pointer to an object.
func (j *Json) GetStruct(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return gconv.Struct(j.Get(pattern), pointer, mapping...)
}

// GetStructDeep does GetStruct recursively.
// Deprecated, use GetStruct instead.
func (j *Json) GetStructDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return gconv.StructDeep(j.Get(pattern), pointer, mapping...)
}

// GetStructs converts any slice to given struct slice.
func (j *Json) GetStructs(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return gconv.Structs(j.Get(pattern), pointer, mapping...)
}

// GetStructsDeep converts any slice to given struct slice recursively.
// Deprecated, use GetStructs instead.
func (j *Json) GetStructsDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return gconv.StructsDeep(j.Get(pattern), pointer, mapping...)
}

// GetScan automatically calls Struct or Structs function according to the type of parameter
// <pointer> to implement the converting..
func (j *Json) GetScan(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return gconv.Scan(j.Get(pattern), pointer, mapping...)
}

// GetScanDeep automatically calls StructDeep or StructsDeep function according to the type of
// parameter <pointer> to implement the converting..
func (j *Json) GetScanDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return gconv.ScanDeep(j.Get(pattern), pointer, mapping...)
}

// GetMapToMap retrieves the value by specified <pattern> and converts it to specified map variable.
// See gconv.MapToMap.
func (j *Json) GetMapToMap(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return gconv.MapToMap(j.Get(pattern), pointer, mapping...)
}

// GetMapToMapDeep retrieves the value by specified <pattern> and converts it to specified map
// variable recursively.
// See gconv.MapToMapDeep.
func (j *Json) GetMapToMapDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return gconv.MapToMapDeep(j.Get(pattern), pointer, mapping...)
}

// GetMapToMaps retrieves the value by specified <pattern> and converts it to specified map slice
// variable.
// See gconv.MapToMaps.
func (j *Json) GetMapToMaps(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return gconv.MapToMaps(j.Get(pattern), pointer, mapping...)
}

// GetMapToMapsDeep retrieves the value by specified <pattern> and converts it to specified map slice
// variable recursively.
// See gconv.MapToMapsDeep.
func (j *Json) GetMapToMapsDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	return gconv.MapToMapsDeep(j.Get(pattern), pointer, mapping...)
}

// Map converts current Json object to map[string]interface{}.
// It returns nil if fails.
func (j *Json) Map() map[string]interface{} {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return gconv.Map(*(j.p))
}

// Array converts current Json object to []interface{}.
// It returns nil if fails.
func (j *Json) Array() []interface{} {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return gconv.Interfaces(*(j.p))
}

// Struct converts current Json object to specified object.
// The <pointer> should be a pointer type of *struct.
func (j *Json) Struct(pointer interface{}, mapping ...map[string]string) error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return gconv.Struct(*(j.p), pointer, mapping...)
}

// Structs converts current Json object to specified object slice.
// The <pointer> should be a pointer type of []struct/*struct.
func (j *Json) Structs(pointer interface{}, mapping ...map[string]string) error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return gconv.Structs(*(j.p), pointer, mapping...)
}

// Scan automatically calls Struct or Structs function according to the type of parameter
// <pointer> to implement the converting..
func (j *Json) Scan(pointer interface{}, mapping ...map[string]string) error {
	return gconv.Scan(*(j.p), pointer, mapping...)
}

// MapToMap converts current Json object to specified map variable.
// The parameter of <pointer> should be type of *map.
func (j *Json) MapToMap(pointer interface{}, mapping ...map[string]string) error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return gconv.MapToMap(*(j.p), pointer, mapping...)
}

// MapToMaps converts current Json object to specified map variable slice.
// The parameter of <pointer> should be type of []map/*map.
func (j *Json) MapToMaps(pointer interface{}, mapping ...map[string]string) error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return gconv.MapToMaps(*(j.p), pointer, mapping...)
}

// Dump prints current Json object with more manually readable.
func (j *Json) Dump() {
	j.mu.RLock()
	defer j.mu.RUnlock()
	gutil.Dump(*j.p)
}

// Export returns <j> as a string with more manually readable.
func (j *Json) Export() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return gutil.Export(*j.p)
}
