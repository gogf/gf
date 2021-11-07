// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils

import "fmt"

// ListToMapByKey converts `list` to a map[string]interface{} of which key is specified by `key`.
// Note that the item value may be type of slice.
func ListToMapByKey(list []map[string]interface{}, key string) map[string]interface{} {
	var (
		s              = ""
		m              = make(map[string]interface{})
		tempMap        = make(map[string][]interface{})
		hasMultiValues bool
	)
	for _, item := range list {
		if k, ok := item[key]; ok {
			s = fmt.Sprintf(`%v`, k)
			tempMap[s] = append(tempMap[s], item)
			if len(tempMap[s]) > 1 {
				hasMultiValues = true
			}
		}
	}
	for k, v := range tempMap {
		if hasMultiValues {
			m[k] = v
		} else {
			m[k] = v[0]
		}
	}
	return m
}
