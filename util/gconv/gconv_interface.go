// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// apiString is used for type assert api for String().
type apiString interface {
	String() string
}

// apiError is used for type assert api for Error().
type apiError interface {
	Error() string
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

// apiSet is the interface for custom value assignment.
type apiSet interface {
	Set(value interface{}) (old interface{})
}
