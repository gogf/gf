// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structcache

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

// addCustomConvertersOutType
// Save the [outType] of the custom conversion function for
// quick detection of whether a type has implemented custom conversion
func (cf *ConvertConfig) addCustomConvertersOutType(fieldType reflect.Type) {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	cf.customConvertersOutTypes[fieldType] = struct{}{}
}

func (cf *ConvertConfig) checkTypeMaybeCustomConvert(fieldType reflect.Type) bool {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	_, ok := cf.customConvertersOutTypes[fieldType]
	return ok
}

func (cf *ConvertConfig) GetCustomConverters(srcType reflect.Type) (map[converterOutType]converterFunc, bool) {
	registeredOutTypeMap, ok := cf.customConverters[srcType]
	return registeredOutTypeMap, ok
}

func (cf *ConvertConfig) HasNoCustomConverters() bool {
	return len(cf.customConverters) == 0
}

func (cf *ConvertConfig) RegisterConverter(fn interface{}) (err error) {
	var (
		fnReflectType = reflect.TypeOf(fn)
		errType       = reflect.TypeOf((*error)(nil)).Elem()
	)
	if fnReflectType.Kind() != reflect.Func ||
		fnReflectType.NumIn() != 1 || fnReflectType.NumOut() != 2 ||
		!fnReflectType.Out(1).Implements(errType) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"parameter must be type of converter function and defined as pattern `func(T1) (T2, error)`, but defined as `%s`",
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
			"invalid converter function `%s`: invalid input parameter type `%s`, should not be type of pointer",
			fnReflectType.String(), inType.String(),
		)
		return
	}
	if outType.Kind() != reflect.Pointer {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"invalid converter function `%s`: invalid output parameter type `%s` should be type of pointer",
			fnReflectType.String(), outType.String(),
		)
		return
	}

	registeredOutTypeMap, ok := cf.customConverters[inType]
	if !ok {
		registeredOutTypeMap = make(map[converterOutType]converterFunc)
		cf.customConverters[inType] = registeredOutTypeMap
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
	cf.addCustomConvertersOutType(outType)
	return
}
