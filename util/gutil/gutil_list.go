// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"reflect"
)

// ListItemValues retrieves and returns the elements of all item struct/map with key <key>.
// Note that the parameter <list> should be type of slice which contains elements of map or struct,
// or else it returns an empty slice.
func ListItemValues(list interface{}, key interface{}) (values []interface{}) {
	// If the given <list> is the most common used slice type []map[string]interface{},
	// it enhances the performance using type assertion.
	if l, ok := list.([]map[string]interface{}); ok {
		if mapKey, ok := key.(string); ok {
			if len(l) == 0 {
				return
			}
			if _, ok := l[0][mapKey]; !ok {
				return
			}
			values = make([]interface{}, len(l))
			for k, m := range l {
				values[k] = m[mapKey]
			}
			return
		}
	}

	// It uses reflect for common checks and converting.
	var (
		reflectValue = reflect.ValueOf(list)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	values = []interface{}{}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		if reflectValue.Len() == 0 {
			return
		}
		var (
			itemValue     reflect.Value
			givenKeyValue = reflect.ValueOf(key)
		)
		for i := 0; i < reflectValue.Len(); i++ {
			itemValue = reflectValue.Index(i)
			// If the items are type of interface{}.
			if itemValue.Kind() == reflect.Interface {
				itemValue = itemValue.Elem()
			}
			if itemValue.Kind() == reflect.Ptr {
				itemValue = itemValue.Elem()
			}
			switch itemValue.Kind() {
			case reflect.Map:
				v := itemValue.MapIndex(givenKeyValue)
				if v.IsValid() {
					values = append(values, v.Interface())
				}

			case reflect.Struct:
				// The <mapKey> must be type of string.
				v := itemValue.FieldByName(givenKeyValue.String())
				if v.IsValid() {
					values = append(values, v.Interface())
				}
			default:
				return
			}
		}
		return
	default:
		return
	}
}

// ListItemValuesUnique retrieves and returns the unique elements of all struct/map with key <key>.
// Note that the parameter <list> should be type of slice which contains elements of map or struct,
// or else it returns an empty slice.
func ListItemValuesUnique(list interface{}, key string) []interface{} {
	values := ListItemValues(list, key)
	if len(values) > 0 {
		var (
			ok bool
			m  = make(map[interface{}]struct{}, len(values))
		)
		for i := 0; i < len(values); {
			if _, ok = m[values[i]]; ok {
				values = SliceDelete(values, i)
			} else {
				m[values[i]] = struct{}{}
				i++
			}
		}
	}
	return values
}
