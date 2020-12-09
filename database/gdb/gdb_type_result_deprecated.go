// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

// Deprecated.
func (r Result) ToJson() string {
	return r.Json()
}

// Deprecated.
func (r Result) ToXml(rootTag ...string) string {
	return r.Xml(rootTag...)
}

// Deprecated.
func (r Result) ToList() List {
	return r.List()
}

// Deprecated.
func (r Result) ToStringMap(key string) map[string]Map {
	return r.MapKeyStr(key)
}

// Deprecated.
func (r Result) ToIntMap(key string) map[int]Map {
	return r.MapKeyInt(key)
}

// Deprecated.
func (r Result) ToUintMap(key string) map[uint]Map {
	return r.MapKeyUint(key)
}

// Deprecated.
func (r Result) ToStringRecord(key string) map[string]Record {
	return r.RecordKeyStr(key)
}

// Deprecated.
func (r Result) ToIntRecord(key string) map[int]Record {
	return r.RecordKeyInt(key)
}

// Deprecated.
func (r Result) ToUintRecord(key string) map[uint]Record {
	return r.RecordKeyUint(key)
}

// Deprecated.
func (r Result) ToStructs(pointer interface{}) (err error) {
	return r.Structs(pointer)
}
