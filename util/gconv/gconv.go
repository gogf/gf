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
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
	"github.com/gogf/gf/v2/util/gconv/internal/structcache"
)

var (
	// Empty strings.
	emptyStringMap = map[string]struct{}{
		"":      {},
		"0":     {},
		"no":    {},
		"off":   {},
		"false": {},
	}
)

// IUnmarshalValue is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface IUnmarshalValue.
type IUnmarshalValue = localinterface.IUnmarshalValue

func init() {
	// register common converters for internal usage.
	structcache.RegisterCommonConverter(structcache.CommonConverter{
		Int64:   Int64,
		Uint64:  Uint64,
		String:  String,
		Float32: Float32,
		Float64: Float64,
		Time:    Time,
		GTime:   GTime,
		Bytes:   Bytes,
		Bool:    Bool,
	})
}
