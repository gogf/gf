// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import "github.com/gogf/gf/v2/os/gtime"

// iString is used for type assert api for String().
type iString interface {
	String() string
}

// iBool is used for type assert api for Bool().
type iBool interface {
	Bool() bool
}

// iInt64 is used for type assert api for Int64().
type iInt64 interface {
	Int64() int64
}

// iUint64 is used for type assert api for Uint64().
type iUint64 interface {
	Uint64() uint64
}

// iFloat32 is used for type assert api for Float32().
type iFloat32 interface {
	Float32() float32
}

// iFloat64 is used for type assert api for Float64().
type iFloat64 interface {
	Float64() float64
}

// iError is used for type assert api for Error().
type iError interface {
	Error() string
}

// iBytes is used for type assert api for Bytes().
type iBytes interface {
	Bytes() []byte
}

// iInterface is used for type assert api for Interface().
type iInterface interface {
	Interface() interface{}
}

// iInterfaces is used for type assert api for Interfaces().
type iInterfaces interface {
	Interfaces() []interface{}
}

// iFloats is used for type assert api for Floats().
type iFloats interface {
	Floats() []float64
}

// iInts is used for type assert api for Ints().
type iInts interface {
	Ints() []int
}

// iStrings is used for type assert api for Strings().
type iStrings interface {
	Strings() []string
}

// iUints is used for type assert api for Uints().
type iUints interface {
	Uints() []uint
}

// iMapStrAny is the interface support for converting struct parameter to map.
type iMapStrAny interface {
	MapStrAny() map[string]interface{}
}

// iUnmarshalValue is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalValue.
type iUnmarshalValue interface {
	UnmarshalValue(interface{}) error
}

// iUnmarshalText is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalText.
type iUnmarshalText interface {
	UnmarshalText(text []byte) error
}

// iUnmarshalText is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalJSON.
type iUnmarshalJSON interface {
	UnmarshalJSON(b []byte) error
}

// iSet is the interface for custom value assignment.
type iSet interface {
	Set(value interface{}) (old interface{})
}

// iGTime is the interface for gtime.Time converting.
type iGTime interface {
	GTime(format ...string) *gtime.Time
}
