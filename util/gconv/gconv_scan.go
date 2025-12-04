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
