// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package structcache provides struct and field info cache feature to enhance performance for struct converting.
package structcache

import (
	"reflect"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

type AnyConvertFunc func(from any, to reflect.Value) error

// ConvertConfig is the configuration for type converting.
type ConvertConfig struct {
	// map[reflect.Type]*CachedStructInfo
	cachedStructsInfoMap sync.Map

	// customConvertTypeMap is used to store whether field types are registered to custom conversions
	// For example:
	// func (src *TypeA) (dst *TypeB,err error)
	// This map will store `TypeB` for quick judgment during assignment.
	// TODO remove?
	customConvertTypeMap map[reflect.Type]struct{}

	// anyToTypeConvertMap for custom type converting from any to its reflect.Value.
	anyToTypeConvertMap map[reflect.Type]AnyConvertFunc
}

// CommonTypeConverter holds some converting functions of common types for internal usage.
type CommonTypeConverter struct {
	Int64   func(v any) (int64, error)
	Uint64  func(v any) (uint64, error)
	String  func(v any) (string, error)
	Float32 func(v any) (float32, error)
	Float64 func(v any) (float64, error)
	Time    func(v any, format ...string) (time.Time, error)
	GTime   func(v any, format ...string) (*gtime.Time, error)
	Bytes   func(v any) ([]byte, error)
	Bool    func(v any) (bool, error)
}

// NewConvertConfig creates and returns a new ConvertConfig object.
func NewConvertConfig() *ConvertConfig {
	return &ConvertConfig{
		cachedStructsInfoMap: sync.Map{},
		customConvertTypeMap: make(map[reflect.Type]struct{}),
		anyToTypeConvertMap:  make(map[reflect.Type]AnyConvertFunc),
	}
}

// RegisterCustomConvertType registers custom
func (cf *ConvertConfig) RegisterCustomConvertType(fieldType reflect.Type) {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	cf.customConvertTypeMap[fieldType] = struct{}{}
}

// RegisterAnyConvertFunc registers custom type converting function for specified type.
func (cf *ConvertConfig) RegisterAnyConvertFunc(t reflect.Type, convertFunc AnyConvertFunc) {
	cf.anyToTypeConvertMap[t] = convertFunc
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
