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

type AnyConvertFunc = structcache.AnyConvertFunc

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

var (
	// defaultConverter is the default management object converting.
	defaultConverter = NewConverter()
)

// RegisterConverter registers custom converter.
// Deprecated: use RegisterTypeConverterFunc instead for clear
func RegisterConverter(fn any) (err error) {
	return RegisterTypeConverterFunc(fn)
}

// RegisterTypeConverterFunc registers custom converter.
func RegisterTypeConverterFunc(fn any) (err error) {
	return defaultConverter.RegisterTypeConverterFunc(fn)
}
