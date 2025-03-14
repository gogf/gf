// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package structcache provides struct and field info cache feature to enhance performance for struct converting.
package structcache

import (
	"context"
	"reflect"
	"runtime"
	"sync"

	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

type interfaceTypeConverter struct {
	interfaceType reflect.Type
	convertFunc   AnyConvertFunc
}

// Converter is the configuration for type converting.
type Converter struct {
	// map[reflect.Type]*CachedStructInfo
	cachedStructsInfoMap sync.Map

	// anyToTypeConvertMap for custom type converting from any to its reflect.Value.
	anyToTypeConvertMap map[reflect.Type]AnyConvertFunc

	// interfaceToTypeConvertMap used for converting any interface type
	// the reason why map is not used here, is because interface types cannot be instantiated
	interfaceToTypeConvertMap []interfaceTypeConverter

	// typeConverterFuncMarkMap is used to store whether field types are registered to custom conversions
	typeConverterFuncMarkMap map[reflect.Type]struct{}
}

// AnyConvertFunc is the function type for converting any to specified type.
// Note that the parameter `to` is usually a pointer type.
type AnyConvertFunc func(from any, to reflect.Value) error

// NewConverter creates and returns a new Converter object.
func NewConverter() *Converter {
	return &Converter{
		cachedStructsInfoMap:     sync.Map{},
		typeConverterFuncMarkMap: make(map[reflect.Type]struct{}),
		anyToTypeConvertMap:      make(map[reflect.Type]AnyConvertFunc),
	}
}

// MarkTypeConvertFunc marks converting function registered for custom type.
func (cf *Converter) MarkTypeConvertFunc(fieldType reflect.Type) {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	cf.typeConverterFuncMarkMap[fieldType] = struct{}{}
}

// RegisterAnyConvertFunc registers custom type converting function for specified type.
func (cf *Converter) RegisterAnyConvertFunc(dstType reflect.Type, convertFunc AnyConvertFunc) {
	if dstType == nil || convertFunc == nil {
		return
	}
	for dstType.Kind() == reflect.Ptr {
		dstType = dstType.Elem()
	}
	if dstType.Kind() == reflect.Interface {
		cf.interfaceToTypeConvertMap = append(cf.interfaceToTypeConvertMap, interfaceTypeConverter{
			interfaceType: dstType,
			convertFunc:   convertFunc,
		})
		return
	}
	cf.anyToTypeConvertMap[dstType] = convertFunc
	intlog.Printf(
		context.Background(),
		`RegisterAnyConvertFunc: %s -> %s`,
		dstType.String(), runtime.FuncForPC(reflect.ValueOf(convertFunc).Pointer()).Name(),
	)
}

// GetAnyConvertFuncByType retrieves and returns the converting function for specified type.
func (cf *Converter) GetAnyConvertFuncByType(dstType reflect.Type) AnyConvertFunc {
	if dstType.Kind() == reflect.Ptr {
		dstType = dstType.Elem()
	}
	return cf.anyToTypeConvertMap[dstType]
}

// IsAnyConvertFuncEmpty checks whether there's any converting function registered.
func (cf *Converter) IsAnyConvertFuncEmpty() bool {
	return len(cf.anyToTypeConvertMap) == 0
}

func (cf *Converter) checkTypeImplInterface(t reflect.Type) AnyConvertFunc {
	if t.Kind() != reflect.Ptr {
		t = reflect.PointerTo(t)
	}
	for _, inter := range cf.interfaceToTypeConvertMap {
		if t.Implements(inter.interfaceType) {
			return inter.convertFunc
		}
	}
	return nil
}

var (
	implUnmarshalText  = reflect.TypeOf((*localinterface.IUnmarshalText)(nil)).Elem()
	implUnmarshalJson  = reflect.TypeOf((*localinterface.IUnmarshalJSON)(nil)).Elem()
	implUnmarshalValue = reflect.TypeOf((*localinterface.IUnmarshalValue)(nil)).Elem()
)

func checkTypeIsCommonInterface(field reflect.StructField) bool {
	isCommonInterface := false
	switch field.Type.String() {
	case "time.Time", "*time.Time":
		// default convert.

	case "gtime.Time", "*gtime.Time":
		// default convert.

	default:
		// Implemented three types of interfaces that must be pointer types, otherwise it is meaningless
		if field.Type.Kind() != reflect.Ptr {
			field.Type = reflect.PointerTo(field.Type)
		}
		switch {
		case field.Type.Implements(implUnmarshalText):
			isCommonInterface = true

		case field.Type.Implements(implUnmarshalJson):
			isCommonInterface = true

		case field.Type.Implements(implUnmarshalValue):
			isCommonInterface = true
		}
	}
	return isCommonInterface
}
