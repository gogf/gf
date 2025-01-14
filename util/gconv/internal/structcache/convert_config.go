// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structcache

import (
	"reflect"
	"sync"
)

type convertFn = func(from any, to reflect.Value) error

type interfaceTypeConvert struct {
	typ reflect.Type
	fn  convertFn
}

type ConvertConfig struct {
	name                  string
	parseConvertFuncs     map[reflect.Type]convertFn
	interfaceConvertFuncs []interfaceTypeConvert
	// map[reflect.Type]*CachedStructInfo
	cachedStructsInfoMap sync.Map
}

func NewConvertConfig(name string) *ConvertConfig {
	return &ConvertConfig{
		name:              name,
		parseConvertFuncs: make(map[reflect.Type]convertFn),
	}
}

func (cf *ConvertConfig) RegisterTypeConvertFunc(typ reflect.Type, f convertFn) {
	if typ == nil || f == nil {
		panic("Parameter cannot be empty")
	}
	if typ.Kind() == reflect.Interface && typ.NumMethod() > 0 {
		panic("Please register using the [RegisterInterfaceTypeConvertFunc] function")
	}
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	cf.parseConvertFuncs[typ] = f
}

// RegisterInterfaceTypeConvertFunc
// typ.Kind == reflect.Interface
// [typ] must be an interface type and cannot be an empty interface
func (cf *ConvertConfig) RegisterInterfaceTypeConvertFunc(typ reflect.Type, f convertFn) {
	if typ == nil || f == nil {
		panic("Parameter cannot be empty")
	}
	if typ.Kind() != reflect.Interface {
		panic("Please register using the [RegisterTypeConvertFunc] function")
	}
	if typ.NumMethod() == 0 {
		panic("Please register using the [RegisterTypeConvertFunc] function")
	}
	cf.interfaceConvertFuncs = append(cf.interfaceConvertFuncs,
		interfaceTypeConvert{
			typ: typ,
			fn:  f,
		},
	)
}

func (cf *ConvertConfig) checkTypeImplInterface(typ reflect.Type) convertFn {
	for _, inter := range cf.interfaceConvertFuncs {
		if typ.Implements(inter.typ) {
			return inter.fn
		}
	}
	return nil
}

// RegisterDefaultConvertFuncs
// Register some commonly used type conversion functions,
// Example:
// 1.int,int8,int16,int32,int64
// 2.uint,uint8,uint16,uint32,uint64
// 3.float32, float64
// 4.bool
// 5.string,[]byte
// 6.time.Time,gtime.Time
func (cf *ConvertConfig) RegisterDefaultConvertFuncs() {
	registerDefaultConvertFuncs(cf)
}

func (cf *ConvertConfig) getCachedConvertStructInfo(structType reflect.Type) (*CachedStructInfo, bool) {
	// Temporarily enabled as an experimental feature
	v, ok := cf.cachedStructsInfoMap.Load(structType)
	if ok {
		return v.(*CachedStructInfo), ok
	}
	return nil, false
}

func (cf *ConvertConfig) storeCachedStructInfo(structType reflect.Type, cachedStructInfo *CachedStructInfo) {
	// Temporarily enabled as an experimental feature
	cf.cachedStructsInfoMap.Store(structType, cachedStructInfo)
}

func (cf *ConvertConfig) getTypeConvertFunc(typ reflect.Type) (fn convertFn) {
	ptr := 0
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		ptr++
	}
	fn = cf.parseConvertFuncs[typ]
	if fn == nil {
		// TODO is value type  typ.Addr
		typ = reflect.PointerTo(typ)
		fn = cf.checkTypeImplInterface(typ)
	}
	if fn != nil {
		fn = ptrConvertFunc(ptr, fn)
	}
	return fn
}

func ptrConvertFunc(ptr int, fn convertFn) convertFn {
	for i := 0; i < ptr; i++ {
		fn = getPtrConvertFunc(fn)
	}
	return fn
}

func getPtrConvertFunc(fn convertFn) convertFn {
	if fn == nil {
		panic("The conversion function cannot be empty")
	}
	return func(from any, to reflect.Value) error {
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}
		// from = nil
		// to = nil ??
		return fn(from, to.Elem())
	}
}
