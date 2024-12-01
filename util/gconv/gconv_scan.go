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
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// Scan automatically checks the type of `pointer` and converts `params` to `pointer`.
//
// TODO change `paramKeyToAttrMap` to `ScanOption` to be more scalable; add `DeepCopy` option for `ScanOption`.
func Scan(srcValue any, dstPointer any, paramKeyToAttrMap ...map[string]string) (err error) {
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
	var (
		dstPointerReflectValueElem     = dstPointerReflectValue.Elem()
		dstPointerReflectValueElemKind = dstPointerReflectValueElem.Kind()
	)
	// Handle multiple level pointers
	if dstPointerReflectValueElemKind == reflect.Ptr {
		// Create new value for pointer dereference
		nextLevelPtr := reflect.New(dstPointerReflectValueElem.Type().Elem())
		// Recursively scan into the dereferenced pointer
		if err = Scan(srcValueReflectValue, nextLevelPtr, paramKeyToAttrMap...); err == nil {
			dstPointerReflectValueElem.Set(nextLevelPtr)
		}
		return
	}

	// if `srcValue` and `dstPointer` are the same type, the do directly assignment.
	// for performance enhancement purpose.
	if ok := doConvertWithTypeCheck(srcValueReflectValue, dstPointerReflectValueElem); ok {
		return nil
	}

	switch dstPointerReflectValueElemKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dstPointerReflectValueElem.SetInt(Int64(srcValue))
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dstPointerReflectValueElem.SetUint(Uint64(srcValue))
		return nil
	case reflect.Float32, reflect.Float64:
		dstPointerReflectValueElem.SetFloat(Float64(srcValue))
		return nil
	case reflect.String:
		dstPointerReflectValueElem.SetString(String(srcValue))
		return nil
	case reflect.Bool:
		dstPointerReflectValueElem.SetBool(Bool(srcValue))
		return nil
	case reflect.Slice:
		// The slice element is struct.
		var (
			dstElemType = dstPointerReflectValueElem.Type().Elem()
			dstElemKind = dstElemType.Kind()
		)
		// The slice element might be type of pointer.
		if dstElemKind == reflect.Ptr {
			dstElemType = dstElemType.Elem()
			dstElemKind = dstElemType.Kind()
		}
		if dstElemKind == reflect.Struct || dstElemKind == reflect.Map {
			return doScanForComplicatedTypes(srcValue, dstPointer, dstPointerReflectType, paramKeyToAttrMap...)
		}
		// Handle slice type conversions
		var srcValueReflectValueKind = srcValueReflectValue.Kind()
		if srcValueReflectValueKind == reflect.Slice || srcValueReflectValueKind == reflect.Array {
			var (
				srcLen   = srcValueReflectValue.Len()
				newSlice = reflect.MakeSlice(dstPointerReflectValueElem.Type(), srcLen, srcLen)
			)
			for i := 0; i < srcLen; i++ {
				srcElem := srcValueReflectValue.Index(i).Interface()
				switch dstElemType.Kind() {
				case reflect.String:
					newSlice.Index(i).SetString(String(srcElem))
				case reflect.Int:
					newSlice.Index(i).SetInt(Int64(srcElem))
				case reflect.Int64:
					newSlice.Index(i).SetInt(Int64(srcElem))
				case reflect.Float64:
					newSlice.Index(i).SetFloat(Float64(srcElem))
				case reflect.Bool:
					newSlice.Index(i).SetBool(Bool(srcElem))
				default:
					return Scan(
						srcElem, newSlice.Index(i).Addr().Interface(), paramKeyToAttrMap...,
					)
				}
			}
			dstPointerReflectValueElem.Set(newSlice)
			return nil
		}
		return doScanForComplicatedTypes(srcValue, dstPointer, dstPointerReflectType, paramKeyToAttrMap...)

	default:
		return doScanForComplicatedTypes(srcValue, dstPointer, dstPointerReflectType, paramKeyToAttrMap...)
	}
}

func doScanForComplicatedTypes(
	srcValue, dstPointer any,
	dstPointerReflectType reflect.Type,
	paramKeyToAttrMap ...map[string]string,
) error {
	// json converting check.
	ok, err := doConvertWithJsonCheck(srcValue, dstPointer)
	if err != nil {
		return err
	}
	if ok {
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

// doConvertWithTypeCheck supports `pointer` in type of `*map/*[]map/*[]*map/*struct/**struct/*[]struct/*[]*struct`
// for converting.
func doConvertWithTypeCheck(srcValueReflectValue, dstPointerReflectValueElem reflect.Value) (ok bool) {
	if !dstPointerReflectValueElem.IsValid() || !srcValueReflectValue.IsValid() {
		return false
	}
	switch {
	// Examples:
	// UploadFile       => UploadFile
	// []UploadFile     => []UploadFile
	// *UploadFile      => *UploadFile
	// *[]UploadFile    => *[]UploadFile
	// map[int][int]    => map[int][int]
	// []map[int][int]  => []map[int][int]
	// *[]map[int][int] => *[]map[int][int]
	case dstPointerReflectValueElem.Type() == srcValueReflectValue.Type():
		dstPointerReflectValueElem.Set(srcValueReflectValue)
		return true

	// Examples:
	// UploadFile      => *UploadFile
	// []UploadFile    => *[]UploadFile
	// map[int][int]   => *map[int][int]
	// []map[int][int] => *[]map[int][int]
	case dstPointerReflectValueElem.Kind() == reflect.Ptr &&
		dstPointerReflectValueElem.Elem().IsValid() &&
		dstPointerReflectValueElem.Elem().Type() == srcValueReflectValue.Type():
		dstPointerReflectValueElem.Elem().Set(srcValueReflectValue)
		return true

	// Examples:
	// *UploadFile      => UploadFile
	// *[]UploadFile    => []UploadFile
	// *map[int][int]   => map[int][int]
	// *[]map[int][int] => []map[int][int]
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
func doConvertWithJsonCheck(srcValue any, dstPointer any) (ok bool, err error) {
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
		if v, ok := srcValue.(localinterface.IInterface); ok {
			return doConvertWithJsonCheck(v.Interface(), dstPointer)
		}
	}
	return false, nil
}
