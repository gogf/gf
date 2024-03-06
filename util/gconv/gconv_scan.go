// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
)

// Scan automatically checks the type of `pointer` and converts `params` to `pointer`.
// It supports `pointer` in type of `*map/*[]map/*[]*map/*struct/**struct/*[]struct/*[]*struct` for converting.
//
// TODO change `paramKeyToAttrMap` to `ScanOption` to be more scalable; add `DeepCopy` option for `ScanOption`.
func Scan(srcValue interface{}, dstPointer interface{}, paramKeyToAttrMap ...map[string]string) (err error) {
	if srcValue == nil {
		// If `srcValue` is nil, no conversion.
		return nil
	}
	if dstPointer == nil {
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`destination pointer should not be nil`,
		)
	}

	// json converting check.
	ok, err := doConvertWithJsonCheck(srcValue, dstPointer)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	var (
		dstPointerReflectType  reflect.Type
		dstPointerReflectValue reflect.Value
	)
	if v, ok := dstPointer.(reflect.Value); ok {
		dstPointerReflectValue = v
		dstPointerReflectType = v.Type()
	} else {
		dstPointerReflectValue = reflect.ValueOf(dstPointer)
		// do not use dstPointerReflectValue.Type() as dstPointerReflectValue might be zero.
		dstPointerReflectType = reflect.TypeOf(dstPointer)
	}

	// pointer kind validation.
	var dstPointerReflectKind = dstPointerReflectType.Kind()
	if dstPointerReflectKind != reflect.Ptr {
		if dstPointerReflectValue.CanAddr() {
			dstPointerReflectValue = dstPointerReflectValue.Addr()
			dstPointerReflectType = dstPointerReflectValue.Type()
			dstPointerReflectKind = dstPointerReflectType.Kind()
		} else {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`destination pointer should be type of pointer, but got type: %v`,
				dstPointerReflectType,
			)
		}
	}
	// direct assignment checks!
	var srcValueReflectValue reflect.Value
	if v, ok := srcValue.(reflect.Value); ok {
		srcValueReflectValue = v
	} else {
		srcValueReflectValue = reflect.ValueOf(srcValue)
	}
	// if `srcValue` and `dstPointer` are the same type, the do directly assignment.
	// For performance enhancement purpose.
	var dstPointerReflectValueElem = dstPointerReflectValue.Elem()
	// if `srcValue` and `dstPointer` are the same type, the do directly assignment.
	// for performance enhancement purpose.
	if ok = doConvertWithTypeCheck(srcValueReflectValue, dstPointerReflectValueElem); ok {
		return nil
	}

	// do the converting.
	var (
		dstPointerReflectTypeElem     = dstPointerReflectType.Elem()
		dstPointerReflectTypeElemKind = dstPointerReflectTypeElem.Kind()
		keyToAttributeNameMapping     map[string]string
	)
	if len(paramKeyToAttrMap) > 0 {
		keyToAttributeNameMapping = paramKeyToAttrMap[0]
	}
	switch dstPointerReflectTypeElemKind {
	case reflect.Map:
		return doMapToMap(srcValue, dstPointer, paramKeyToAttrMap...)

	case reflect.Array, reflect.Slice:
		var (
			sliceElem     = dstPointerReflectTypeElem.Elem()
			sliceElemKind = sliceElem.Kind()
		)
		for sliceElemKind == reflect.Ptr {
			sliceElem = sliceElem.Elem()
			sliceElemKind = sliceElem.Kind()
		}
		if sliceElemKind == reflect.Map {
			return doMapToMaps(srcValue, dstPointer, paramKeyToAttrMap...)
		}
		return doStructs(srcValue, dstPointer, keyToAttributeNameMapping, "")

	default:
		return doStruct(srcValue, dstPointer, keyToAttributeNameMapping, "")
	}
}

func doConvertWithTypeCheck(srcValueReflectValue, dstPointerReflectValueElem reflect.Value) (ok bool) {
	if !dstPointerReflectValueElem.IsValid() || !srcValueReflectValue.IsValid() {
		return false
	}
	switch {
	// Example:
	// UploadFile    => UploadFile
	// []UploadFile  => []UploadFile
	// *UploadFile   => *UploadFile
	// *[]UploadFile => *[]UploadFile
	// map           => map
	// []map         => []map
	// *[]map        => *[]map
	case dstPointerReflectValueElem.Type() == srcValueReflectValue.Type():
		dstPointerReflectValueElem.Set(srcValueReflectValue)
		return true

	// Example:
	// UploadFile    => *UploadFile
	// []UploadFile  => *[]UploadFile
	// map           => *map
	// []map         => *[]map
	case dstPointerReflectValueElem.Kind() == reflect.Ptr &&
		dstPointerReflectValueElem.Elem().IsValid() &&
		dstPointerReflectValueElem.Elem().Type() == srcValueReflectValue.Type():
		dstPointerReflectValueElem.Elem().Set(srcValueReflectValue)
		return true

	// Example:
	// *UploadFile    => UploadFile
	// *[]UploadFile  => []UploadFile
	// *map           => map
	// *[]map         => []map
	case srcValueReflectValue.Kind() == reflect.Ptr &&
		srcValueReflectValue.Elem().IsValid() &&
		dstPointerReflectValueElem.Type() == srcValueReflectValue.Elem().Type():
		dstPointerReflectValueElem.Set(srcValueReflectValue.Elem())
		return true

	default:
		return false
	}
}

// doConvertWithJsonCheck does json converting check.
// If given `params` is JSON, it then uses json.Unmarshal doing the converting.
func doConvertWithJsonCheck(srcValue interface{}, dstPointer interface{}) (ok bool, err error) {
	switch valueResult := srcValue.(type) {
	case []byte:
		if json.Valid(valueResult) {
			if dstPointerReflectType, ok := dstPointer.(reflect.Value); ok {
				if dstPointerReflectType.Kind() == reflect.Ptr {
					if dstPointerReflectType.IsNil() {
						return false, nil
					}
					return true, json.UnmarshalUseNumber(valueResult, dstPointerReflectType.Interface())
				} else if dstPointerReflectType.CanAddr() {
					return true, json.UnmarshalUseNumber(valueResult, dstPointerReflectType.Addr().Interface())
				}
			} else {
				return true, json.UnmarshalUseNumber(valueResult, dstPointer)
			}
		}

	case string:
		if valueBytes := []byte(valueResult); json.Valid(valueBytes) {
			if dstPointerReflectType, ok := dstPointer.(reflect.Value); ok {
				if dstPointerReflectType.Kind() == reflect.Ptr {
					if dstPointerReflectType.IsNil() {
						return false, nil
					}
					return true, json.UnmarshalUseNumber(valueBytes, dstPointerReflectType.Interface())
				} else if dstPointerReflectType.CanAddr() {
					return true, json.UnmarshalUseNumber(valueBytes, dstPointerReflectType.Addr().Interface())
				}
			} else {
				return true, json.UnmarshalUseNumber(valueBytes, dstPointer)
			}
		}

	default:
		// The `params` might be struct that implements interface function Interface, eg: gvar.Var.
		if v, ok := srcValue.(iInterface); ok {
			return doConvertWithJsonCheck(v.Interface(), dstPointer)
		}
	}
	return false, nil
}
