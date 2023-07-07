// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package conv

// SliceAny is alias of Interfaces.
func SliceAny[T any](v interface{}) []T {
	return Interfaces[T](v)
}

// Interfaces converts `any` to []T.
func Interfaces[T any](val interface{}) []T {
	if val == nil {
		return nil
	}
	var (
		array []T
	)
	switch value := val.(type) {
	case []T:
		array = make([]T, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []interface{}:
		array = make([]T, len(value))
		for k, v := range value {
			array[k], _ = v.(T)
		}
	}
	return array
}
