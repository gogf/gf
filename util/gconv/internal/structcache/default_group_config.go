// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structcache

import (
	"reflect"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// CommonConverter holds some converting functions of common types for internal usage.
type CommonConverter struct {
	Int64   func(any interface{}) int64
	Uint64  func(any interface{}) uint64
	String  func(any interface{}) string
	Float32 func(any interface{}) float32
	Float64 func(any interface{}) float64
	Time    func(any interface{}, format ...string) time.Time
	GTime   func(any interface{}, format ...string) *gtime.Time
	Bytes   func(any interface{}) []byte
	Bool    func(any interface{}) bool
}

var (
	// localCommonConverter holds some converting functions of common types for internal usage.
	localCommonConverter CommonConverter
)

var defaultConfig = NewConvertConfig("gconv.convert.default")

func GetDefaultConfig() *ConvertConfig {
	return defaultConfig
}

// RegisterCommonConverter registers the CommonConverter for local usage.
func RegisterCommonConverter(commonConverter CommonConverter) {
	localCommonConverter = commonConverter
}

func init() {
	registerDefaultConvertFuncs(defaultConfig)
}

func registerDefaultConvertFuncs(cfg *ConvertConfig) {
	registerManyTypesConvertFn(cfg, intConvertFunc, intType, int8Type, int16Type, int32Type, int64Type)
	registerManyTypesConvertFn(cfg, uintConvertFunc, uintType, uint8Type, uint16Type, uint32Type, uint64Type)
	registerManyTypesConvertFn(cfg, floatConvertFunc, float32Type, float64Type)
	registerManyTypesConvertFn(cfg, stringConvertFunc, stringType)
	registerManyTypesConvertFn(cfg, boolConvertFunc, boolType)
	registerManyTypesConvertFn(cfg, bytesConvertFunc, bytesType)
	registerManyTypesConvertFn(cfg, timeConvertFunc, timeType)
	registerManyTypesConvertFn(cfg, gtimeConvertFunc, gtimeType)
}

func registerManyTypesConvertFn(cfg *ConvertConfig, fn convertFn, typs ...reflect.Type) {
	for _, typ := range typs {
		cfg.RegisterTypeConvertFunc(typ, fn)
	}
}

func intConvertFunc(from any, to reflect.Value) error {
	to.SetInt(localCommonConverter.Int64(from))
	return nil
}

func uintConvertFunc(from any, to reflect.Value) error {
	to.SetUint(localCommonConverter.Uint64(from))
	return nil
}

func floatConvertFunc(from any, to reflect.Value) error {
	to.SetFloat(localCommonConverter.Float64(from))
	return nil
}

func stringConvertFunc(from any, to reflect.Value) error {
	to.SetString(localCommonConverter.String(from))
	return nil
}

func boolConvertFunc(from any, to reflect.Value) error {
	to.SetBool(localCommonConverter.Bool(from))
	return nil
}

func bytesConvertFunc(from any, to reflect.Value) error {
	to.SetBytes(localCommonConverter.Bytes(from))
	return nil
}

func timeConvertFunc(from any, to reflect.Value) error {
	*to.Addr().Interface().(*time.Time) = localCommonConverter.Time(from)
	return nil
}

func gtimeConvertFunc(from any, to reflect.Value) error {
	v := localCommonConverter.GTime(from)
	if v == nil {
		v = gtime.New()
	}
	*to.Addr().Interface().(*gtime.Time) = *v
	return nil
}
