// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
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

// Type wraps reflect.Type for additional features.
type Type struct {
	reflect.Type
}

// Field contains information of a struct field .
type Field struct {
	Value    reflect.Value       // The underlying value of the field.
	Field    reflect.StructField // The underlying field of the field.
	TagValue string              // Retrieved tag value. There might be more than one tags in the field, but only one can be retrieved according to calling function rules.
}
