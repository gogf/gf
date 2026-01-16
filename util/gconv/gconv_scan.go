// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// Scan automatically checks the type of `pointer` and converts `params` to `pointer`.
// It supports various types of parameter conversions, including:
// 1. Basic types (int, string, float, etc.)
// 2. Pointer types
// 3. Slice types
// 4. Map types
// 5. Struct types
//
// The `paramKeyToAttrMap` parameter is used for mapping between attribute names and parameter keys.
// TODO: change `paramKeyToAttrMap` to `ScanOption` to be more scalable; add `DeepCopy` option for `ScanOption`.
func Scan(srcValue any, dstPointer any, paramKeyToAttrMap ...map[string]string) (err error) {
	option := ScanOption{
		ContinueOnError: true,
	}
	if len(paramKeyToAttrMap) > 0 {
		option.ParamKeyToAttrMap = paramKeyToAttrMap[0]
	}
	return defaultConverter.Scan(srcValue, dstPointer, option)
}

// ScanWithOptions automatically checks the type of `dstPointer` and converts `srcValue` to `dstPointer`.
// It is the same as Scan function, but accepts one or more ScanOption values for additional conversion control.
//
// When using ScanWithOptions, the term "omit" means that the assignment from the source to the destination
// is skipped, so the existing value in the destination field is preserved.
//
//   - option.OmitEmpty, when set to true, skips assignment of empty source values (for example: empty strings,
//     zero numeric values, zero time values, empty slices or maps), preserving any existing non-empty values
//     in the destination.
//
//   - option.OmitNil, when set to true, skips assignment of nil source values, preserving the existing values
//     in the destination when the source contains nil.
//
// Example:
//
//	type User struct {
//	    Name  string
//	    Email string
//	}
//
//	dst := &User{Name: "Alice", Email: "alice@example.com"}
//	src := map[string]any{
//	    "Name":  "",
//	    "Email": nil,
//	}
//
//	// With OmitEmpty and OmitNil, empty and nil values in src will not overwrite dst.
//	err := ScanWithOptions(src, dst, ScanOption{OmitEmpty: true, OmitNil: true})
func ScanWithOptions(srcValue any, dstPointer any, option ...ScanOption) (err error) {
	return defaultConverter.Scan(srcValue, dstPointer, option...)
}
