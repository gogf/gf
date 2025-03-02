// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import "github.com/gogf/gf/v2/internal/json"

// SliceMap converts `value` to []map[string]any.
// Note that it automatically checks and converts json string to []map if `value` is string/[]byte.
func (c *Converter) SliceMap(value any, sliceOption SliceOption, mapOption MapOption) ([]map[string]any, error) {
	if value == nil {
		return nil, nil
	}
	switch r := value.(type) {
	case string:
		list := make([]map[string]any, 0)
		if len(r) > 0 && r[0] == '[' && r[len(r)-1] == ']' {
			if err := json.UnmarshalUseNumber([]byte(r), &list); err != nil {
				return nil, err
			}
			return list, nil
		}
		return nil, nil

	case []byte:
		list := make([]map[string]any, 0)
		if len(r) > 0 && r[0] == '[' && r[len(r)-1] == ']' {
			if err := json.UnmarshalUseNumber(r, &list); err != nil {
				return nil, err
			}
			return list, nil
		}
		return nil, nil

	case []map[string]any:
		return r, nil

	default:
		array, err := c.SliceAny(value, sliceOption)
		if err != nil {
			return nil, err
		}
		if len(array) == 0 {
			return nil, nil
		}
		list := make([]map[string]any, len(array))
		for k, v := range array {
			m, err := c.Map(v, mapOption)
			if err != nil && sliceOption.BreakOnError {
				return nil, err
			}
			list[k] = m
		}
		return list, nil
	}
}
