// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package localinterface defines some interfaces for converting usage.
package localinterface

import "github.com/gogf/gf/v2/os/gtime"

// IVal is used for type assert api for String().
type IVal interface {
	Val() interface{}
}

// IString is used for type assert api for String().
type IString interface {
	String() string
}

// IBool is used for type assert api for Bool().
type IBool interface {
	Bool() bool
}

// IInt64 is used for type assert api for Int64().
type IInt64 interface {
	Int64() int64
}

// IUint64 is used for type assert api for Uint64().
type IUint64 interface {
	Uint64() uint64
}

// IFloat32 is used for type assert api for Float32().
type IFloat32 interface {
	Float32() float32
}

// IFloat64 is used for type assert api for Float64().
type IFloat64 interface {
	Float64() float64
}

// IError is used for type assert api for Error().
type IError interface {
	Error() string
}

// IBytes is used for type assert api for Bytes().
type IBytes interface {
	Bytes() []byte
}

// IInterface is used for type assert api for Interface().
type IInterface interface {
	Interface() interface{}
}

// IInterfaces is used for type assert api for Interfaces().
type IInterfaces interface {
	Interfaces() []interface{}
}

// IFloats is used for type assert api for Floats().
type IFloats interface {
	Floats() []float64
}

// IInts is used for type assert api for Ints().
type IInts interface {
	Ints() []int
}

// IStrings is used for type assert api for Strings().
type IStrings interface {
	Strings() []string
}

// IUints is used for type assert api for Uints().
type IUints interface {
	Uints() []uint
}

// IMapStrAny is the interface support for converting struct parameter to map.
type IMapStrAny interface {
	MapStrAny() map[string]interface{}
}

// IUnmarshalValue is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalValue.
type IUnmarshalValue interface {
	UnmarshalValue(interface{}) error
}

// IUnmarshalText is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalText.
type IUnmarshalText interface {
	UnmarshalText(text []byte) error
}

// IUnmarshalJSON is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalJSON.
type IUnmarshalJSON interface {
	UnmarshalJSON(b []byte) error
}

// ISet is the interface for custom value assignment.
type ISet interface {
	Set(value interface{}) (old interface{})
}

// IGTime is the interface for gtime.Time converting.
type IGTime interface {
	GTime(format ...string) *gtime.Time
}
