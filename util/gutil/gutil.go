// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gutil provides utility functions.
package gutil

import (
	"reflect"

	"github.com/gogf/gf/v2/util/gconv"
)

const (
	dumpIndent = `    `
)

// Keys retrieves and returns the keys from given map or struct.
func Keys(mapOrStruct interface{}) (keysOrAttrs []string) {
	keysOrAttrs = make([]string, 0)
	if m, ok := mapOrStruct.(map[string]interface{}); ok {
		for k := range m {
			keysOrAttrs = append(keysOrAttrs, k)
		}
		return
	}
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
	)
	if v, ok := mapOrStruct.(reflect.Value); ok {
		reflectValue = v
	} else {
		reflectValue = reflect.ValueOf(mapOrStruct)
	}
	reflectKind = reflectValue.Kind()
	for reflectKind == reflect.Ptr {
		if !reflectValue.IsValid() || reflectValue.IsNil() {
			reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
			reflectKind = reflectValue.Kind()
		} else {
			reflectValue = reflectValue.Elem()
			reflectKind = reflectValue.Kind()
		}
	}
	switch reflectKind {
	case reflect.Map:
		for _, k := range reflectValue.MapKeys() {
			keysOrAttrs = append(keysOrAttrs, gconv.String(k.Interface()))
		}
	case reflect.Struct:
		var (
			fieldType   reflect.StructField
			reflectType = reflectValue.Type()
		)
		for i := 0; i < reflectValue.NumField(); i++ {
			fieldType = reflectType.Field(i)
			if fieldType.Anonymous {
				keysOrAttrs = append(keysOrAttrs, Keys(reflectValue.Field(i))...)
			} else {
				keysOrAttrs = append(keysOrAttrs, fieldType.Name)
			}
		}
	}
	return
}

// Values retrieves and returns the values from given map or struct.
func Values(mapOrStruct interface{}) (values []interface{}) {
	values = make([]interface{}, 0)
	if m, ok := mapOrStruct.(map[string]interface{}); ok {
		for _, v := range m {
			values = append(values, v)
		}
		return
	}
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
	)
	if v, ok := mapOrStruct.(reflect.Value); ok {
		reflectValue = v
	} else {
		reflectValue = reflect.ValueOf(mapOrStruct)
	}
	reflectKind = reflectValue.Kind()
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Map:
		for _, k := range reflectValue.MapKeys() {
			values = append(values, reflectValue.MapIndex(k).Interface())
		}
	case reflect.Struct:
		var (
			fieldType   reflect.StructField
			reflectType = reflectValue.Type()
		)
		for i := 0; i < reflectValue.NumField(); i++ {
			fieldType = reflectType.Field(i)
			if fieldType.Anonymous {
				values = append(values, Values(reflectValue.Field(i))...)
			} else {
				values = append(values, reflectValue.Field(i).Interface())
			}
		}
	}
	return
}
