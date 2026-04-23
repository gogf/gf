// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
)

// RuleUnique implements `unique` rule:
// Array/slice elements or map values should be unique.
// For slices/arrays of struct, it supports checking uniqueness by field name using `unique:FieldName`.
//
// Format: unique
// Format: unique:FieldName
type RuleUnique struct{}

func init() {
	Register(RuleUnique{})
}

func (r RuleUnique) Name() string {
	return "unique"
}

func (r RuleUnique) Message() string {
	return "The {field} field must contain unique values"
}

func (r RuleUnique) Run(in RunInput) error {
	reflectValue := indirectValue(reflect.ValueOf(in.Value.Val()))
	if !reflectValue.IsValid() {
		return nil
	}
	switch reflectValue.Kind() {
	case reflect.Slice, reflect.Array:
		return checkUniqueSliceOrArray(reflectValue, in.RulePattern, in.Message)
	case reflect.Map:
		return checkUniqueMapValues(reflectValue, in.Message)
	}
	return gerror.NewCodef(
		gcode.CodeInvalidParameter,
		`validation rule "%s" only supports slice, array or map values`,
		"unique",
	)
}

func checkUniqueSliceOrArray(reflectValue reflect.Value, fieldName string, message string) error {
	seen := make(map[string]struct{}, reflectValue.Len())
	for i := 0; i < reflectValue.Len(); i++ {
		elem := reflectValue.Index(i)
		if fieldName != "" {
			var err error
			elem, err = getUniqueStructFieldValue(elem, fieldName)
			if err != nil {
				return err
			}
		}
		key, err := makeUniqueValueKey(elem)
		if err != nil {
			return err
		}
		if _, ok := seen[key]; ok {
			return errors.New(message)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func checkUniqueMapValues(reflectValue reflect.Value, message string) error {
	seen := make(map[string]struct{}, reflectValue.Len())
	for _, mapKey := range reflectValue.MapKeys() {
		key, err := makeUniqueValueKey(reflectValue.MapIndex(mapKey))
		if err != nil {
			return err
		}
		if _, ok := seen[key]; ok {
			return errors.New(message)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func getUniqueStructFieldValue(value reflect.Value, fieldName string) (reflect.Value, error) {
	value = indirectValue(value)
	if !value.IsValid() {
		return reflect.Value{}, nil
	}
	if value.Kind() != reflect.Struct {
		return reflect.Value{}, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`validation rule "unique" with field name only supports slice or array of struct values`,
		)
	}
	fieldValue := value.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return reflect.Value{}, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid field name "%s" for validation rule "unique"`,
			fieldName,
		)
	}
	return fieldValue, nil
}

func makeUniqueValueKey(value reflect.Value) (string, error) {
	value = indirectValue(value)
	if !value.IsValid() {
		return "<nil>", nil
	}
	if value.CanInterface() && value.Type().Comparable() {
		return fmt.Sprintf("%T:%v", value.Interface(), value.Interface()), nil
	}
	if value.CanInterface() {
		content, err := json.Marshal(value.Interface())
		if err == nil {
			return fmt.Sprintf("%s:%s", value.Type().String(), string(content)), nil
		}
	}
	return "", gerror.NewCodef(
		gcode.CodeInvalidParameter,
		`value type "%s" is not supported by validation rule "unique"`,
		value.Type().String(),
	)
}

func indirectValue(value reflect.Value) reflect.Value {
	for value.IsValid() && (value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}
	return value
}
