// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// Convert converts the variable `fromValue` to the type `toTypeName`, the type `toTypeName` is specified by string.
//
// The optional parameter `extraParams` is used for additional necessary parameter for this conversion.
// It supports common basic types conversion as its conversion based on type name string.
func Convert(fromValue any, toTypeName string, extraParams ...any) any {
	result, _ := defaultConverter.ConvertWithTypeName(fromValue, toTypeName, ConvertOption{
		ExtraParams:  extraParams,
		SliceOption:  SliceOption{ContinueOnError: true},
		MapOption:    MapOption{ContinueOnError: true},
		StructOption: StructOption{ContinueOnError: true},
	})
	return result
}

// ConvertWithRefer converts the variable `fromValue` to the type referred by value `referValue`.
//
// The optional parameter `extraParams` is used for additional necessary parameter for this conversion.
// It supports common basic types conversion as its conversion based on type name string.
func ConvertWithRefer(fromValue any, referValue any, extraParams ...any) any {
	result, _ := defaultConverter.ConvertWithRefer(fromValue, referValue, ConvertOption{
		ExtraParams:  extraParams,
		SliceOption:  SliceOption{ContinueOnError: true},
		MapOption:    MapOption{ContinueOnError: true},
		StructOption: StructOption{ContinueOnError: true},
	})
	return result
}
