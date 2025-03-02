// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv/internal/structcache"
)

type (
	converterInType  = reflect.Type
	converterOutType = reflect.Type
	converterFunc    = reflect.Value
)

// Converter is the manager for type converting.
type Converter interface {
	RegisterTypeConverterFunc(fn any) (err error)
	ConverterForInt
	ConverterForUint
	String(any any) (string, error)
	Bool(any any) (bool, error)
	Bytes(any any) ([]byte, error)
	Float32(any any) (float32, error)
	Float64(any any) (float64, error)

	MapToMap(params any, pointer any, mapping ...map[string]string) (err error)
	MapToMaps(params any, pointer any, paramKeyToAttrMap ...map[string]string) (err error)
	Rune(any any) (rune, error)
	Runes(any any) ([]rune, error)
	Scan(srcValue any, dstPointer any, paramKeyToAttrMap ...map[string]string) (err error)
	Time(any interface{}, format ...string) (time.Time, error)
	Duration(any interface{}) (time.Duration, error)
	GTime(any interface{}, format ...string) (*gtime.Time, error)
}

type ConverterForInt interface {
	Int(any any) (int, error)
	Int8(any any) (int8, error)
	Int16(any any) (int16, error)
	Int32(any any) (int32, error)
	Int64(any any) (int64, error)
}

type ConverterForUint interface {
	Uint(any any) (uint, error)
	Uint8(any any) (uint8, error)
	Uin16(any any) (uint16, error)
	Uint32(any any) (uint32, error)
	Uint64(any any) (uint64, error)
}

// impConverter implements the interface Converter.
type impConverter struct {
	internalConverter    *structcache.Converter
	typeConverterFuncMap map[converterInType]map[converterOutType]converterFunc
}

var (
	intType   = reflect.TypeOf(0)
	int8Type  = reflect.TypeOf(int8(0))
	int16Type = reflect.TypeOf(int16(0))
	int32Type = reflect.TypeOf(int32(0))
	int64Type = reflect.TypeOf(int64(0))

	uintType   = reflect.TypeOf(uint(0))
	uint8Type  = reflect.TypeOf(uint8(0))
	uint16Type = reflect.TypeOf(uint16(0))
	uint32Type = reflect.TypeOf(uint32(0))
	uint64Type = reflect.TypeOf(uint64(0))

	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))

	stringType = reflect.TypeOf("")
	bytesType  = reflect.TypeOf([]byte{})

	boolType = reflect.TypeOf(false)

	timeType  = reflect.TypeOf((*time.Time)(nil)).Elem()
	gtimeType = reflect.TypeOf((*gtime.Time)(nil)).Elem()
)

// NewConverter creates and returns management object for type converting.
func NewConverter() *impConverter {
	cf := &impConverter{
		internalConverter:    structcache.NewConverter(),
		typeConverterFuncMap: make(map[converterInType]map[converterOutType]converterFunc),
	}
	cf.registerBuiltInConverter()
	return cf
}

// RegisterTypeConverterFunc registers custom converter.
// It must be registered before you use this custom converting feature.
// It is suggested to do it in boot procedure of the process.
//
// Note:
//  1. The parameter `fn` must be defined as pattern `func(T1) (T2, error)`.
//     It will convert type `T1` to type `T2`.
//  2. The `T1` should not be type of pointer, but the `T2` should be type of pointer.
func (c *impConverter) RegisterTypeConverterFunc(fn any) (err error) {
	var (
		fnReflectType = reflect.TypeOf(fn)
		errType       = reflect.TypeOf((*error)(nil)).Elem()
	)
	if fnReflectType.Kind() != reflect.Func ||
		fnReflectType.NumIn() != 1 || fnReflectType.NumOut() != 2 ||
		!fnReflectType.Out(1).Implements(errType) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"parameter must be type of converter function and defined as pattern `func(T1) (T2, error)`, "+
				"but defined as `%s`",
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

	registeredOutTypeMap, ok := c.typeConverterFuncMap[inType]
	if !ok {
		registeredOutTypeMap = make(map[converterOutType]converterFunc)
		c.typeConverterFuncMap[inType] = registeredOutTypeMap
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
	c.internalConverter.RegisterTypeConvertFunc(outType)
	return
}

func (c *impConverter) registerBuiltInConverter() {
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForInt64, intType, int8Type, int16Type, int32Type, int64Type,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForUint64, uintType, uint8Type, uint16Type, uint32Type, uint64Type,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForString, stringType,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForFloat64, float32Type, float64Type,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForBool, boolType,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForBytes, bytesType,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForTime, timeType,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForGTime, gtimeType,
	)
}

func (c *impConverter) registerAnyConvertFuncForTypes(convertFunc AnyConvertFunc, types ...reflect.Type) {
	for _, t := range types {
		c.internalConverter.RegisterAnyConvertFunc(t, convertFunc)
	}
}
