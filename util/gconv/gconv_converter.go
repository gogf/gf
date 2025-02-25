// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
)

type (
	converterInType  = reflect.Type
	converterOutType = reflect.Type
	converterFunc    = reflect.Value
)

// RegisterConverter to register custom converter.
// It must be registered before you use this custom converting feature.
// It is suggested to do it in boot procedure of the process.
//
// Note:
//  1. The parameter `fn` must be defined as pattern `func(T1) (T2, error)`.
//     It will convert type `T1` to type `T2`.
//  2. The `T1` should not be type of pointer, but the `T2` should be type of pointer.
func RegisterConverter(fn interface{}) (err error) {
	err = defaultConfig.RegisterConverter(fn)
	return
}

func getRegisteredConverterFuncAndSrcType(
	srcReflectValue, dstReflectValueForRefer reflect.Value,
) (f converterFunc, srcType reflect.Type, ok bool) {
	if defaultConfig.HasNoCustomConverters() {
		return reflect.Value{}, nil, false
	}
	srcType = srcReflectValue.Type()
	for srcType.Kind() == reflect.Pointer {
		srcType = srcType.Elem()
	}
	var registeredOutTypeMap map[converterOutType]converterFunc
	// firstly, it searches the map by input parameter type.
	registeredOutTypeMap, ok = defaultConfig.GetCustomConverters(srcType)
	if !ok {
		return reflect.Value{}, nil, false
	}
	var dstType = dstReflectValueForRefer.Type()
	if dstType.Kind() == reflect.Pointer {
		// Might be **struct, which is support as designed.
		if dstType.Elem().Kind() == reflect.Pointer {
			dstType = dstType.Elem()
		}
	} else if dstReflectValueForRefer.IsValid() && dstReflectValueForRefer.CanAddr() {
		dstType = dstReflectValueForRefer.Addr().Type()
	} else {
		dstType = reflect.PointerTo(dstType)
	}
	// secondly, it searches the input parameter type map
	// and finds the result converter function by the output parameter type.
	f, ok = registeredOutTypeMap[dstType]
	if !ok {
		return reflect.Value{}, nil, false
	}
	return
}

func callCustomConverterWithRefer(
	srcReflectValue, referReflectValue reflect.Value,
) (dstReflectValue reflect.Value, converted bool, err error) {
	registeredConverterFunc, srcType, ok := getRegisteredConverterFuncAndSrcType(srcReflectValue, referReflectValue)
	if !ok {
		return reflect.Value{}, false, nil
	}
	dstReflectValue = reflect.New(referReflectValue.Type()).Elem()
	converted, err = doCallCustomConverter(srcReflectValue, dstReflectValue, registeredConverterFunc, srcType)
	return
}

// callCustomConverter call the custom converter. It will try some possible type.
func callCustomConverter(srcReflectValue, dstReflectValue reflect.Value) (converted bool, err error) {
	registeredConverterFunc, srcType, ok := getRegisteredConverterFuncAndSrcType(srcReflectValue, dstReflectValue)
	if !ok {
		return false, nil
	}
	return doCallCustomConverter(srcReflectValue, dstReflectValue, registeredConverterFunc, srcType)
}

func doCallCustomConverter(
	srcReflectValue reflect.Value,
	dstReflectValue reflect.Value,
	registeredConverterFunc converterFunc,
	srcType reflect.Type,
) (converted bool, err error) {
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
