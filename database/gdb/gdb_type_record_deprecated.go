// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

// Deprecated, use Json instead.
func (r Record) ToJson() string {
	return r.Json()
}

// Deprecated, use Xml instead.
func (r Record) ToXml(rootTag ...string) string {
	return r.Xml(rootTag...)
}

// Deprecated, use Map instead.
func (r Record) ToMap() Map {
	return r.Map()
}

// Deprecated, use Struct instead.
func (r Record) ToStruct(pointer interface{}) error {
	return r.Struct(pointer)
}
