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
)

type (
	converterInType  = reflect.Type
	converterOutType = reflect.Type
	converterFunc    = reflect.Value
)

// customConverters for internal converter storing.
var customConverters = make(map[converterInType]map[converterOutType]converterFunc)

// RegisterConverter to register custom converter.
// It must be registered before you use this custom converting feature.
// It is suggested to do it in boot.
//
// Note:
//  1. The parameter `fn` must be defined as pattern `func(T1) (T2, error)`.
//     It will convert type `T1` to type `T2`.
//  2. The `T1` should not be type of pointer, but the `T2` should be type of pointer.
func RegisterConverter(fn interface{}) (err error) {
	var (
		fnReflectType = reflect.TypeOf(fn)
		errType       = reflect.TypeOf((*error)(nil)).Elem()
	)
	if fnReflectType.Kind() != reflect.Func ||
		fnReflectType.NumIn() != 1 || fnReflectType.NumOut() != 2 ||
		!fnReflectType.Out(1).Implements(errType) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"parameter must be type of function and defined as pattern `func(T1) (T2, error)`, but defined as `%s`",
			fnReflectType.String(),
		)
		return
	}

	// The Key and Value of the converter map should not be pointer.
	var (
		inType  = fnReflectType.In(0)
		outType = fnReflectType.Out(0)
	)
	if inType.Kind() == reflect.Pointer {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"invalid input parameter type `%s`: should not be type of pointer",
			inType.String(),
		)
		return
	}
	if outType.Kind() != reflect.Pointer {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"invalid output parameter type `%s`: should be type of pointer",
			outType.String(),
		)
		return
	}

	registeredOutTypeMap, ok := customConverters[inType]
	if !ok {
		registeredOutTypeMap = make(map[converterOutType]converterFunc)
		customConverters[inType] = registeredOutTypeMap
	}
	if _, ok = registeredOutTypeMap[outType]; ok {
		err = gerror.NewCodef(
			gcode.CodeInvalidOperation,
			"the converter parameter type `%s` to type `%s` has already been registered",
			inType.String(), outType.String(),
		)
		return
	}
	registeredOutTypeMap[outType] = reflect.ValueOf(fn)
	return
}

// callCustomConverter call the custom converter. It will try some possible type.
func callCustomConverter(srcReflectValue reflect.Value, dstReflectValue reflect.Value) (converted bool, err error) {
	if len(customConverters) == 0 {
		return false, nil
	}
	var (
		ok      bool
		srcType = srcReflectValue.Type()
	)
	for srcType.Kind() == reflect.Pointer {
		srcType = srcType.Elem()
	}
	var (
		registeredOutTypeMap    map[converterOutType]converterFunc
		registeredConverterFunc converterFunc
	)
	// firstly, it searches the map by input parameter type.
	registeredOutTypeMap, ok = customConverters[srcType]
	if !ok {
		return false, nil
	}
	var dstType = dstReflectValue.Type()
	if dstType.Kind() == reflect.Pointer && dstReflectValue.Elem().Kind() == reflect.Pointer {
		dstType = dstReflectValue.Elem().Type()
	} else if dstType.Kind() != reflect.Pointer && dstReflectValue.CanAddr() {
		dstType = dstReflectValue.Addr().Type()
	}
	// secondly, it searches the input parameter type map
	// and finds the result converter function by the output parameter type.
	registeredConverterFunc, ok = registeredOutTypeMap[dstType]
	if !ok {
		return false, nil
	}
	// Converter function calling.
	for srcReflectValue.Type() != srcType {
		srcReflectValue = srcReflectValue.Elem()
	}
	result := registeredConverterFunc.Call([]reflect.Value{srcReflectValue})
	if !result[1].IsNil() {
		return false, result[1].Interface().(error)
	}
	// The `result[0]` is a pointer.
	if result[0].IsNil() {
		return false, nil
	}
	var resultValue = result[0]
	for {
		if resultValue.Type() == dstReflectValue.Type() && dstReflectValue.CanSet() {
			dstReflectValue.Set(resultValue)
			converted = true
		} else if dstReflectValue.Kind() == reflect.Pointer {
			if resultValue.Type() == dstReflectValue.Elem().Type() && dstReflectValue.Elem().CanSet() {
				dstReflectValue.Elem().Set(resultValue)
				converted = true
			}
		}
		if converted {
			break
		}
		if resultValue.Kind() == reflect.Pointer {
			resultValue = resultValue.Elem()
		} else {
			break
		}
	}

	return converted, nil
}
