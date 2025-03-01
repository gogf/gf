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
	// 接口类型
	typ reflect.Type
	fn  convertFn
}

type ConvertConfig struct {
	name              string
	parseConvertFuncs map[reflect.Type]convertFn
	// 用于存储接口类型的转换函数
	interfaceConvertFuncs []interfaceTypeConvert
	// map[reflect.Type]*CachedStructInfo
	cachedStructsInfoMap sync.Map

	// 下面的两个map是用于自定义转换的
	customConverters map[converterInType]map[converterOutType]converterFunc
	// customConvertTypeMap is used to store whether field types are registered to custom conversions
	// For example:
	// func (src *TypeA) (dst *TypeB,err error)
	// This map will store `TypeB` for quick judgment during assignment.
	customConvertersOutTypes map[reflect.Type]struct{}
}

func NewConvertConfig(name string) *ConvertConfig {
	return &ConvertConfig{
		name:                     name,
		parseConvertFuncs:        make(map[reflect.Type]convertFn),
		customConverters:         make(map[converterInType]map[converterOutType]converterFunc),
		customConvertersOutTypes: make(map[reflect.Type]struct{}),
	}
}

func (cf *ConvertConfig) RegisterTypeConvertFunc(typ reflect.Type, f convertFn) {
	if typ == nil || f == nil {
		panic("parameter cannot be empty")
	}
	if typ.Kind() == reflect.Interface {
		// 由于接口类型不能被实例化，在每次获取convertFn时，
		// 都需要遍历所有的接口类型convertFn，以此判断类型是否实现了接口
		cf.interfaceConvertFuncs = append(cf.interfaceConvertFuncs,
			interfaceTypeConvert{
				typ: typ,
				fn:  f,
			},
		)
		return
	}
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	cf.parseConvertFuncs[typ] = f
}

// 检查typ是否实现了接口
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
	v, ok := cf.cachedStructsInfoMap.Load(structType)
	if ok {
		return v.(*CachedStructInfo), ok
	}
	return nil, false
}

func (cf *ConvertConfig) setCachedConvertStructInfo(structType reflect.Type, cachedStructInfo *CachedStructInfo) {
	cf.cachedStructsInfoMap.Store(structType, cachedStructInfo)
}

// 获取typ的convertFn
func (cf *ConvertConfig) getTypeConvertFunc(typ reflect.Type) (fn convertFn) {
	ptr := 0
	// 如果是指针类型，获取指针指向的类型
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		ptr++
	}
	fn = cf.parseConvertFuncs[typ]
	// 如果没有找到，检查是否是接口类型
	if fn == nil {
		// 需要变更为指针类型
		typ = reflect.PointerTo(typ)
		// 检查是否实现了接口
		fn = cf.checkTypeImplInterface(typ)
	}
	// 如果找到了，需要根据指针级数，获取对应的convertFn
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

// 如果是指针类型，那么需要再外层包裹一层函数
// 假设从string转换到*int
// 那么根据上面getTypeConvertFunc函数的逻辑，
// 我们可以得到一个string=>int的convertFn
// 由于我们需要的是string=>*int的convertFn
// 所以需要再外层包裹一层函数，提前将*int解引用变为int
// 然后再去调用string=>int的convertFn
func getPtrConvertFunc(fn convertFn) convertFn {
	if fn == nil {
		panic("the conversion function cannot be empty")
	}
	return func(from any, to reflect.Value) error {
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}
		// TODO:
		// from = nil
		// to = nil ??
		return fn(from, to.Elem())
	}
}
