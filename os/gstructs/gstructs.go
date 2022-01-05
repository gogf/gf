// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gstructs provides functions for struct information retrieving.
package gstructs

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

// FieldsInput is the input parameter struct type for function Fields.
type FieldsInput struct {
	// Pointer should be type of struct/*struct.
	Pointer interface{}

	// RecursiveOption specifies the way retrieving the fields recursively if the attribute
	// is an embedded struct. It is RecursiveOptionNone in default.
	RecursiveOption int
}

// FieldMapInput is the input parameter struct type for function FieldMap.
type FieldMapInput struct {
	// Pointer should be type of struct/*struct.
	Pointer interface{}

	// PriorityTagArray specifies the priority tag array for retrieving from high to low.
	// If it's given `nil`, it returns map[name]Field, of which the `name` is attribute name.
	PriorityTagArray []string

	// RecursiveOption specifies the way retrieving the fields recursively if the attribute
	// is an embedded struct. It is RecursiveOptionNone in default.
	RecursiveOption int
}

const (
	RecursiveOptionNone          = 0 // No recursively retrieving fields as map if the field is an embedded struct.
	RecursiveOptionEmbedded      = 1 // Recursively retrieving fields as map if the field is an embedded struct.
	RecursiveOptionEmbeddedNoTag = 2 // Recursively retrieving fields as map if the field is an embedded struct and the field has no tag.
)
