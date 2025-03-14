// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gconv implements powerful and convenient converting functionality for any types of variables.
//
// This package should keep much fewer dependencies with other packages.
package gconv

import (
	"reflect"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv/internal/converter"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
	"github.com/gogf/gf/v2/util/gconv/internal/structcache"
)

// Converter is the manager for type converting.
type Converter interface {
	ConverterForConvert
	ConverterForRegister
	ConverterForInt
	ConverterForUint
	ConverterForTime
	ConverterForFloat
	ConverterForMap
	ConverterForSlice
	ConverterForStruct
	ConverterForBasic
}

// ConverterForBasic is the basic converting interface.
type ConverterForBasic interface {
	Scan(srcValue, dstPointer any, option ...ScanOption) (err error)
	String(any any) (string, error)
	Bool(any any) (bool, error)
	Rune(any any) (rune, error)
}

// ConverterForTime is the converting interface for time.
type ConverterForTime interface {
	Time(v any, format ...string) (time.Time, error)
	Duration(v any) (time.Duration, error)
	GTime(v any, format ...string) (*gtime.Time, error)
}

// ConverterForInt is the converting interface for integer.
type ConverterForInt interface {
	Int(v any) (int, error)
	Int8(v any) (int8, error)
	Int16(v any) (int16, error)
	Int32(v any) (int32, error)
	Int64(v any) (int64, error)
}

// ConverterForUint is the converting interface for unsigned integer.
type ConverterForUint interface {
	Uint(v any) (uint, error)
	Uint8(v any) (uint8, error)
	Uint16(v any) (uint16, error)
	Uint32(v any) (uint32, error)
	Uint64(v any) (uint64, error)
}

// ConverterForFloat is the converting interface for float.
type ConverterForFloat interface {
	Float32(v any) (float32, error)
	Float64(v any) (float64, error)
}

// ConverterForMap is the converting interface for map.
type ConverterForMap interface {
	Map(v any, option ...MapOption) (map[string]any, error)
	MapStrStr(v any, option ...MapOption) (map[string]string, error)
}

// ConverterForSlice is the converting interface for slice.
type ConverterForSlice interface {
	Bytes(v any) ([]byte, error)
	Runes(v any) ([]rune, error)
	SliceAny(v any, option ...SliceOption) ([]any, error)
	SliceFloat32(v any, option ...SliceOption) ([]float32, error)
	SliceFloat64(v any, option ...SliceOption) ([]float64, error)
	SliceInt(v any, option ...SliceOption) ([]int, error)
	SliceInt32(v any, option ...SliceOption) ([]int32, error)
	SliceInt64(v any, option ...SliceOption) ([]int64, error)
	SliceUint(v any, option ...SliceOption) ([]uint, error)
	SliceUint32(v any, option ...SliceOption) ([]uint32, error)
	SliceUint64(v any, option ...SliceOption) ([]uint64, error)
	SliceStr(v any, option ...SliceOption) ([]string, error)
	SliceMap(v any, option ...SliceMapOption) ([]map[string]any, error)
}

// ConverterForStruct is the converting interface for struct.
type ConverterForStruct interface {
	Struct(params, pointer any, option ...StructOption) (err error)
	Structs(params, pointer any, option ...StructsOption) (err error)
}

// ConverterForConvert is the converting interface for custom converting.
type ConverterForConvert interface {
	ConvertWithRefer(fromValue, referValue any, option ...ConvertOption) (any, error)
	ConvertWithTypeName(fromValue any, toTypeName string, option ...ConvertOption) (any, error)
}

// ConverterForRegister is the converting interface for custom converter registration.
type ConverterForRegister interface {
	RegisterTypeConverterFunc(f any) error
	RegisterAnyConverterFunc(f AnyConvertFunc, types ...reflect.Type)
}

type (
	// AnyConvertFunc is the function type for converting any to specified type.
	AnyConvertFunc = structcache.AnyConvertFunc

	// MapOption specifies the option for map converting.
	MapOption = converter.MapOption

	// SliceOption is the option for Slice type converting.
	SliceOption = converter.SliceOption

	// SliceMapOption is the option for SliceMap function.
	SliceMapOption = converter.SliceMapOption

	// ScanOption is the option for the Scan function.
	ScanOption = converter.ScanOption

	// StructOption is the option for Struct converting.
	StructOption = converter.StructOption

	// StructsOption is the option for Structs function.
	StructsOption = converter.StructsOption

	// ConvertOption is the option for converting.
	ConvertOption = converter.ConvertOption
)

// IUnmarshalValue is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface IUnmarshalValue.
type IUnmarshalValue = localinterface.IUnmarshalValue

var (
	// defaultConverter is the default management object converting.
	defaultConverter = converter.NewConverter()
)

// NewConverter creates and returns management object for type converting.
func NewConverter() Converter {
	return converter.NewConverter()
}

// RegisterConverter registers custom converter.
// Deprecated: use RegisterTypeConverterFunc instead for clear
func RegisterConverter(fn any) (err error) {
	return RegisterTypeConverterFunc(fn)
}

// RegisterTypeConverterFunc registers custom converter.
func RegisterTypeConverterFunc(fn any) (err error) {
	return defaultConverter.RegisterTypeConverterFunc(fn)
}

// RegisterAnyConverterFunc registers custom type converting function for specified type.
func RegisterAnyConverterFunc(f AnyConvertFunc, types ...reflect.Type) {
	defaultConverter.RegisterAnyConverterFunc(f, types...)
}
