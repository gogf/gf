// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv/internal/structcache"
)

// CommonTypeConverter holds some converting functions of common types for internal usage.
type CommonTypeConverter = structcache.CommonTypeConverter

type (
	converterInType  = reflect.Type
	converterOutType = reflect.Type
	converterFunc    = reflect.Value
)

// ConvertConfig is the configuration for type converting.
type ConvertConfig struct {
	internalConvertConfig *structcache.ConvertConfig
	// customConverters for internal converter storing.
	customConverters map[converterInType]map[converterOutType]converterFunc
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

// NewConvertConfig creates and returns configuration management object for type converting.
func NewConvertConfig() *ConvertConfig {
	cf := &ConvertConfig{
		internalConvertConfig: structcache.NewConvertConfig(),
		customConverters:      make(map[converterInType]map[converterOutType]converterFunc),
	}
	cf.registerBuiltInConverter()
	return cf
}

func (cf *ConvertConfig) registerBuiltInConverter() {
	cf.registerAnyConvertFuncForTypes(
		builtInAnyConvertFuncForInt64, intType, int8Type, int16Type, int32Type, int64Type,
	)
	cf.registerAnyConvertFuncForTypes(
		builtInAnyConvertFuncForUint64, uintType, uint8Type, uint16Type, uint32Type, uint64Type,
	)
	cf.registerAnyConvertFuncForTypes(
		builtInAnyConvertFuncForString, stringType,
	)
	cf.registerAnyConvertFuncForTypes(
		builtInAnyConvertFuncForFloat64, float32Type, float64Type,
	)
	cf.registerAnyConvertFuncForTypes(
		builtInAnyConvertFuncForBool, boolType,
	)
	cf.registerAnyConvertFuncForTypes(
		builtInAnyConvertFuncForBytes, bytesType,
	)
	cf.registerAnyConvertFuncForTypes(
		builtInAnyConvertFuncForTime, timeType,
	)
	cf.registerAnyConvertFuncForTypes(
		builtInAnyConvertFuncForGTime, gtimeType,
	)
}

func (cf *ConvertConfig) registerAnyConvertFuncForTypes(convertFunc AnyConvertFunc, types ...reflect.Type) {
	for _, t := range types {
		cf.internalConvertConfig.RegisterAnyConvertFunc(t, convertFunc)
	}
}
