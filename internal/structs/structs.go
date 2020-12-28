// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package structs provides functions for struct conversion.
//
// Inspired and improved from: https://github.com/fatih/structs
package structs

import (
	"reflect"
)

// Field contains information of a struct field .
type Field struct {
	value reflect.Value
	field reflect.StructField
	// Retrieved tag value. There might be more than one tags in the field,
	// but only one can be retrieved according to calling function rules.
	TagValue string
}

// Tag returns the value associated with key in the tag string. If there is no
// such key in the tag, Tag returns the empty string.
func (f *Field) Tag(key string) string {
	return f.field.Tag.Get(key)
}

// Value returns the underlying value of the field. It panics if the field
// is not exported.
func (f *Field) Value() interface{} {
	return f.value.Interface()
}

// IsEmbedded returns true if the given field is an anonymous field (embedded)
func (f *Field) IsEmbedded() bool {
	return f.field.Anonymous
}

// IsExported returns true if the given field is exported.
func (f *Field) IsExported() bool {
	return f.field.PkgPath == ""
}

// Name returns the name of the given field
func (f *Field) Name() string {
	return f.field.Name
}
