// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gstructs provides functions for struct information retrieving.
package gstructs

import (
	"reflect"

	"github.com/gogf/gf/v2/errors/gerror"
)

// Type wraps reflect.Type for additional features.
type Type struct {
	reflect.Type
}

// Field contains information of a struct field .
type Field struct {
	Value reflect.Value       // The underlying value of the field.
	Field reflect.StructField // The underlying field of the field.

	// Retrieved tag name. It depends TagValue.
	TagName string

	// Retrieved tag value.
	// There might be more than one tags in the field,
	// but only one can be retrieved according to calling function rules.
	TagValue string
}

// FieldsInput is the input parameter struct type for function Fields.
type FieldsInput struct {
	// Pointer should be type of struct/*struct.
	// TODO this attribute name is not suitable, which would make confuse.
	Pointer interface{}

	// RecursiveOption specifies the way retrieving the fields recursively if the attribute
	// is an embedded struct. It is RecursiveOptionNone in default.
	RecursiveOption RecursiveOption
}

// FieldMapInput is the input parameter struct type for function FieldMap.
type FieldMapInput struct {
	// Pointer should be type of struct/*struct.
	// TODO this attribute name is not suitable, which would make confuse.
	Pointer interface{}

	// PriorityTagArray specifies the priority tag array for retrieving from high to low.
	// If it's given `nil`, it returns map[name]Field, of which the `name` is attribute name.
	PriorityTagArray []string

	// RecursiveOption specifies the way retrieving the fields recursively if the attribute
	// is an embedded struct. It is RecursiveOptionNone in default.
	RecursiveOption RecursiveOption
}

type RecursiveOption int

const (
	RecursiveOptionNone          RecursiveOption = iota // No recursively retrieving fields as map if the field is an embedded struct.
	RecursiveOptionEmbedded                             // Recursively retrieving fields as map if the field is an embedded struct.
	RecursiveOptionEmbeddedNoTag                        // Recursively retrieving fields as map if the field is an embedded struct and the field has no tag.
)

// Fields retrieves and returns the fields of `pointer` as slice.
func Fields(in FieldsInput) ([]Field, error) {
	var (
		ok                   bool
		fieldFilterMap       = make(map[string]struct{})
		retrievedFields      = make([]Field, 0)
		currentLevelFieldMap = make(map[string]Field)
		rangeFields, err     = getFieldValues(in.Pointer)
	)
	if err != nil {
		return nil, err
	}

	for index := 0; index < len(rangeFields); index++ {
		field := rangeFields[index]
		currentLevelFieldMap[field.Name()] = field
	}

	for index := 0; index < len(rangeFields); index++ {
		field := rangeFields[index]
		if _, ok = fieldFilterMap[field.Name()]; ok {
			continue
		}
		if field.IsEmbedded() {
			if in.RecursiveOption != RecursiveOptionNone {
				switch in.RecursiveOption {
				case RecursiveOptionEmbeddedNoTag:
					if field.TagStr() != "" {
						break
					}
					fallthrough

				case RecursiveOptionEmbedded:
					structFields, err := Fields(FieldsInput{
						Pointer:         field.Value,
						RecursiveOption: in.RecursiveOption,
					})
					if err != nil {
						return nil, err
					}
					// The current level fields can overwrite the sub-struct fields with the same name.
					for i := 0; i < len(structFields); i++ {
						var (
							structField = structFields[i]
							fieldName   = structField.Name()
						)
						if _, ok = fieldFilterMap[fieldName]; ok {
							continue
						}
						fieldFilterMap[fieldName] = struct{}{}
						if v, ok := currentLevelFieldMap[fieldName]; !ok {
							retrievedFields = append(retrievedFields, structField)
						} else {
							retrievedFields = append(retrievedFields, v)
						}
					}
					continue
				}
			}
			continue
		}
		fieldFilterMap[field.Name()] = struct{}{}
		retrievedFields = append(retrievedFields, field)
	}
	return retrievedFields, nil
}

// FieldMap retrieves and returns struct field as map[name/tag]Field from `pointer`.
//
// The parameter `pointer` should be type of struct/*struct.
//
// The parameter `priority` specifies the priority tag array for retrieving from high to low.
// If it's given `nil`, it returns map[name]Field, of which the `name` is attribute name.
//
// The parameter `recursive` specifies whether retrieving the fields recursively if the attribute
// is an embedded struct.
//
// Note that it only retrieves the exported attributes with first letter upper-case from struct.
func FieldMap(in FieldMapInput) (map[string]Field, error) {
	fields, err := getFieldValues(in.Pointer)
	if err != nil {
		return nil, err
	}
	var (
		tagValue string
		mapField = make(map[string]Field)
	)
	for _, field := range fields {
		// Only retrieve exported attributes.
		if !field.IsExported() {
			continue
		}
		tagValue = ""
		for _, p := range in.PriorityTagArray {
			tagValue = field.Tag(p)
			if tagValue != "" && tagValue != "-" {
				break
			}
		}
		tempField := field
		tempField.TagValue = tagValue
		if tagValue != "" {
			mapField[tagValue] = tempField
		} else {
			if in.RecursiveOption != RecursiveOptionNone && field.IsEmbedded() {
				switch in.RecursiveOption {
				case RecursiveOptionEmbeddedNoTag:
					if field.TagStr() != "" {
						mapField[field.Name()] = tempField
						break
					}
					fallthrough

				case RecursiveOptionEmbedded:
					m, err := FieldMap(FieldMapInput{
						Pointer:          field.Value,
						PriorityTagArray: in.PriorityTagArray,
						RecursiveOption:  in.RecursiveOption,
					})
					if err != nil {
						return nil, err
					}
					for k, v := range m {
						if _, ok := mapField[k]; !ok {
							tempV := v
							mapField[k] = tempV
						}
					}
				}
			} else {
				mapField[field.Name()] = tempField
			}
		}
	}
	return mapField, nil
}

// StructType retrieves and returns the struct Type of specified struct/*struct.
// The parameter `object` should be either type of struct/*struct/[]struct/[]*struct.
func StructType(object interface{}) (*Type, error) {
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
		reflectType  reflect.Type
	)
	if rv, ok := object.(reflect.Value); ok {
		reflectValue = rv
	} else {
		reflectValue = reflect.ValueOf(object)
	}
	reflectKind = reflectValue.Kind()
	for {
		switch reflectKind {
		case reflect.Ptr:
			if !reflectValue.IsValid() || reflectValue.IsNil() {
				// If pointer is type of *struct and nil, then automatically create a temporary struct.
				reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
				reflectKind = reflectValue.Kind()
			} else {
				reflectValue = reflectValue.Elem()
				reflectKind = reflectValue.Kind()
			}

		case reflect.Array, reflect.Slice:
			reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
			reflectKind = reflectValue.Kind()

		default:
			goto exitLoop
		}
	}

exitLoop:
	if reflectKind != reflect.Struct {
		return nil, gerror.Newf(
			`invalid object kind "%s", kind of "struct" is required`,
			reflectKind,
		)
	}
	reflectType = reflectValue.Type()
	return &Type{
		Type: reflectType,
	}, nil
}
