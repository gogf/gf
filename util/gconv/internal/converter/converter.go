// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package converter provides converting utilities for any types of variables.
package converter

import (
	"reflect"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv/internal/structcache"
)

// AnyConvertFunc is the type for any type converting function.
type AnyConvertFunc = structcache.AnyConvertFunc

// RecursiveType is the type for converting recursively.
type RecursiveType string

const (
	RecursiveTypeAuto RecursiveType = "auto"
	RecursiveTypeTrue RecursiveType = "true"
)

type (
	converterInType  = reflect.Type
	converterOutType = reflect.Type
	converterFunc    = reflect.Value
)

// Converter implements the interface Converter.
type Converter struct {
	internalConverter    *structcache.Converter
	typeConverterFuncMap map[converterInType]map[converterOutType]converterFunc
}

var (
	// Empty strings.
	emptyStringMap = map[string]struct{}{
		"":      {},
		"0":     {},
		"no":    {},
		"off":   {},
		"false": {},
	}
)

// NewConverter creates and returns management object for type converting.
func NewConverter() *Converter {
	cf := &Converter{
		internalConverter:    structcache.NewConverter(),
		typeConverterFuncMap: make(map[converterInType]map[converterOutType]converterFunc),
	}
	cf.registerBuiltInAnyConvertFunc()
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
func (c *Converter) RegisterTypeConverterFunc(f any) (err error) {
	var (
		fReflectType = reflect.TypeOf(f)
		errType      = reflect.TypeOf((*error)(nil)).Elem()
	)
	if fReflectType.Kind() != reflect.Func ||
		fReflectType.NumIn() != 1 || fReflectType.NumOut() != 2 ||
		!fReflectType.Out(1).Implements(errType) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"parameter must be type of converter function and defined as pattern `func(T1) (T2, error)`, "+
				"but defined as `%s`",
			fReflectType.String(),
		)
		return
	}

	// The Key and Value of the converter map should not be pointer.
	var (
		inType  = fReflectType.In(0)
		outType = fReflectType.Out(0)
	)
	if inType.Kind() == reflect.Pointer {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"invalid converter function `%s`: invalid input parameter type `%s`, should not be type of pointer",
			fReflectType.String(), inType.String(),
		)
		return
	}
	if outType.Kind() != reflect.Pointer {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"invalid converter function `%s`: invalid output parameter type `%s` should be type of pointer",
			fReflectType.String(), outType.String(),
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
	registeredOutTypeMap[outType] = reflect.ValueOf(f)
	c.internalConverter.MarkTypeConvertFunc(outType)
	return
}

// RegisterAnyConverterFunc registers custom type converting function for specified types.
func (c *Converter) RegisterAnyConverterFunc(convertFunc AnyConvertFunc, types ...reflect.Type) {
	for _, t := range types {
		c.internalConverter.RegisterAnyConvertFunc(t, convertFunc)
	}
}

func (c *Converter) registerBuiltInAnyConvertFunc() {
	var (
		intType     = reflect.TypeOf(0)
		int8Type    = reflect.TypeOf(int8(0))
		int16Type   = reflect.TypeOf(int16(0))
		int32Type   = reflect.TypeOf(int32(0))
		int64Type   = reflect.TypeOf(int64(0))
		uintType    = reflect.TypeOf(uint(0))
		uint8Type   = reflect.TypeOf(uint8(0))
		uint16Type  = reflect.TypeOf(uint16(0))
		uint32Type  = reflect.TypeOf(uint32(0))
		uint64Type  = reflect.TypeOf(uint64(0))
		float32Type = reflect.TypeOf(float32(0))
		float64Type = reflect.TypeOf(float64(0))
		stringType  = reflect.TypeOf("")
		bytesType   = reflect.TypeOf([]byte{})
		boolType    = reflect.TypeOf(false)
		timeType    = reflect.TypeOf((*time.Time)(nil)).Elem()
		gtimeType   = reflect.TypeOf((*gtime.Time)(nil)).Elem()
	)
	c.RegisterAnyConverterFunc(
		c.builtInAnyConvertFuncForInt64, intType, int8Type, int16Type, int32Type, int64Type,
	)
	c.RegisterAnyConverterFunc(
		c.builtInAnyConvertFuncForUint64, uintType, uint8Type, uint16Type, uint32Type, uint64Type,
	)
	c.RegisterAnyConverterFunc(
		c.builtInAnyConvertFuncForString, stringType,
	)
	c.RegisterAnyConverterFunc(
		c.builtInAnyConvertFuncForFloat64, float32Type, float64Type,
	)
	c.RegisterAnyConverterFunc(
		c.builtInAnyConvertFuncForBool, boolType,
	)
	c.RegisterAnyConverterFunc(
		c.builtInAnyConvertFuncForBytes, bytesType,
	)
	c.RegisterAnyConverterFunc(
		c.builtInAnyConvertFuncForTime, timeType,
	)
	c.RegisterAnyConverterFunc(
		c.builtInAnyConvertFuncForGTime, gtimeType,
	)
}
