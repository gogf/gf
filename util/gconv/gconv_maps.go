// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv/internal/converter"
)

// SliceMap is alias of Maps.
func SliceMap(any any, option ...MapOption) []map[string]any {
	return Maps(any, option...)
}

// SliceMapDeep is alias of MapsDeep.
// Deprecated: used SliceMap instead.
func SliceMapDeep(any any) []map[string]any {
	return MapsDeep(any)
}

// Maps converts `value` to []map[string]any.
// Note that it automatically checks and converts json string to []map if `value` is string/[]byte.
func Maps(value any, option ...MapOption) []map[string]any {
	mapOption := MapOption{
		ContinueOnError: true,
	}
	if len(option) > 0 {
		mapOption = option[0]
	}
	result, _ := defaultConverter.SliceMap(value, SliceMapOption{
		MapOption: mapOption,
		SliceOption: converter.SliceOption{
			ContinueOnError: true,
		},
	})
	return result
}

// MapsDeep converts `value` to []map[string]any recursively.
//
// TODO completely implement the recursive converting for all types.
// Deprecated: used Maps instead.
func MapsDeep(value any, tags ...string) []map[string]any {
	if value == nil {
		return nil
	}
	switch r := value.(type) {
	case string:
		list := make([]map[string]any, 0)
		if len(r) > 0 && r[0] == '[' && r[len(r)-1] == ']' {
			if err := json.UnmarshalUseNumber([]byte(r), &list); err != nil {
				return nil
			}
			return list
		} else {
			return nil
		}

	case []byte:
		list := make([]map[string]any, 0)
		if len(r) > 0 && r[0] == '[' && r[len(r)-1] == ']' {
			if err := json.UnmarshalUseNumber(r, &list); err != nil {
				return nil
			}
			return list
		} else {
			return nil
		}

	case []map[string]any:
		list := make([]map[string]any, len(r))
		for k, v := range r {
			list[k] = MapDeep(v, tags...)
		}
		return list

	default:
		array := Interfaces(value)
		if len(array) == 0 {
			return nil
		}
		list := make([]map[string]any, len(array))
		for k, v := range array {
			list[k] = MapDeep(v, tags...)
		}
		return list
	}
}
