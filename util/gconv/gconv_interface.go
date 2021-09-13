// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import "github.com/gogf/gf/os/gtime"

// apiString is used for type assert api for String().
type apiString interface {
	String() string
}

// apiBool is used for type assert api for Bool().
type apiBool interface {
	Bool() bool
}

// apiInt64 is used for type assert api for Int64().
type apiInt64 interface {
	Int64() int64
}

// apiUint64 is used for type assert api for Uint64().
type apiUint64 interface {
	Uint64() uint64
}

// apiFloat32 is used for type assert api for Float32().
type apiFloat32 interface {
	Float32() float32
}

// apiFloat64 is used for type assert api for Float64().
type apiFloat64 interface {
	Float64() float64
}

// apiError is used for type assert api for Error().
type apiError interface {
	Error() string
}

// apiBytes is used for type assert api for Bytes().
type apiBytes interface {
	Bytes() []byte
}

// apiInterfaces is used for type assert api for Interfaces().
type apiInterfaces interface {
	Interfaces() []interface{}
}

// apiFloats is used for type assert api for Floats().
type apiFloats interface {
	Floats() []float64
}

// apiInts is used for type assert api for Ints().
type apiInts interface {
	Ints() []int
}

// apiStrings is used for type assert api for Strings().
type apiStrings interface {
	Strings() []string
}

// apiUints is used for type assert api for Uints().
type apiUints interface {
	Uints() []uint
}

// apiMapStrAny is the interface support for converting struct parameter to map.
type apiMapStrAny interface {
	MapStrAny() map[string]interface{}
}

// apiUnmarshalValue is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface apiUnmarshalValue.
type apiUnmarshalValue interface {
	UnmarshalValue(interface{}) error
}

// apiUnmarshalText is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface apiUnmarshalText.
type apiUnmarshalText interface {
	UnmarshalText(text []byte) error
}

// apiUnmarshalText is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface apiUnmarshalJSON.
type apiUnmarshalJSON interface {
	UnmarshalJSON(b []byte) error
}

// apiSet is the interface for custom value assignment.
type apiSet interface {
	Set(value interface{}) (old interface{})
}

// apiGTime is the interface for gtime.Time converting.
type apiGTime interface {
	GTime(format ...string) *gtime.Time
}
