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

type ConvertConfig struct {
	name              string
	parseConvertFuncs map[reflect.Type]convertFn
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
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	cf.parseConvertFuncs[typ] = f

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
		return nil
	}
	for i := 0; i < ptr; i++ {
		fn = getPtrConvertFunc(fn)
	}
	return fn
}

func getPtrConvertFunc(
	convertFunc convertFn,
) convertFn {
	if convertFunc == nil {
		panic("The conversion function cannot be empty")
	}
	return func(from any, to reflect.Value) error {
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}
		// from = nil
		// to = nil ??
		return convertFunc(from, to.Elem())
	}
}
