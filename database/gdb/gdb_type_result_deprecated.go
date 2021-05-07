// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

// Deprecated, use Json instead.
func (r Result) ToJson() string {
	return r.Json()
}

// Deprecated, use Xml instead.
func (r Result) ToXml(rootTag ...string) string {
	return r.Xml(rootTag...)
}

// Deprecated, use List instead.
func (r Result) ToList() List {
	return r.List()
}

// Deprecated, use MapKeyStr instead.
func (r Result) ToStringMap(key string) map[string]Map {
	return r.MapKeyStr(key)
}

// Deprecated, use MapKetInt instead.
func (r Result) ToIntMap(key string) map[int]Map {
	return r.MapKeyInt(key)
}

// Deprecated, use MapKeyUint instead.
func (r Result) ToUintMap(key string) map[uint]Map {
	return r.MapKeyUint(key)
}

// Deprecated, use RecordKeyStr instead.
func (r Result) ToStringRecord(key string) map[string]Record {
	return r.RecordKeyStr(key)
}

// Deprecated, use RecordKetInt instead.
func (r Result) ToIntRecord(key string) map[int]Record {
	return r.RecordKeyInt(key)
}

// Deprecated, use RecordKetUint instead.
func (r Result) ToUintRecord(key string) map[uint]Record {
	return r.RecordKeyUint(key)
}

// Deprecated, use Structs instead.
func (r Result) ToStructs(pointer interface{}) (err error) {
	return r.Structs(pointer)
}
