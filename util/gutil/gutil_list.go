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
//
// The parameter <list> supports types like:
// []map[string]interface{}
// []map[string]sub-map
// []struct
// []struct:sub-struct
// Note that the sub-map/sub-struct makes sense only if the optional parameter <subKey> is given.
func ListItemValues(list interface{}, key interface{}, subKey ...interface{}) (values []interface{}) {
	var reflectValue reflect.Value
	if v, ok := list.(reflect.Value); ok {
		reflectValue = v
	} else {
		reflectValue = reflect.ValueOf(list)
	}
	reflectKind := reflectValue.Kind()
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		if reflectValue.Len() == 0 {
			return
		}
		values = []interface{}{}
		for i := 0; i < reflectValue.Len(); i++ {
			if value, ok := ItemValue(reflectValue.Index(i), key); ok {
				if len(subKey) > 0 && subKey[0] != nil {
					if subValue, ok := ItemValue(value, subKey[0]); ok {
						value = subValue
					} else {
						continue
					}
				}
				if array, ok := value.([]interface{}); ok {
					values = append(values, array...)
				} else {
					values = append(values, value)
				}
			}
		}
	}
	return
}

// ItemValue retrieves and returns its value of which name/attribute specified by <key>.
// The parameter <item> can be type of map/*map/struct/*struct.
func ItemValue(item interface{}, key interface{}) (value interface{}, found bool) {
	var reflectValue reflect.Value
	if v, ok := item.(reflect.Value); ok {
		reflectValue = v
	} else {
		reflectValue = reflect.ValueOf(item)
	}
	reflectKind := reflectValue.Kind()
	if reflectKind == reflect.Interface {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	var keyValue reflect.Value
	if v, ok := key.(reflect.Value); ok {
		keyValue = v
	} else {
		keyValue = reflect.ValueOf(key)
	}
	switch reflectKind {
	case reflect.Array, reflect.Slice:
		// The <key> must be type of string.
		values := ListItemValues(reflectValue, keyValue.String())
		if values == nil {
			return nil, false
		}
		return values, true

	case reflect.Map:
		v := reflectValue.MapIndex(keyValue)
		if v.IsValid() {
			found = true
			value = v.Interface()
		}

	case reflect.Struct:
		// The <mapKey> must be type of string.
		v := reflectValue.FieldByName(keyValue.String())
		if v.IsValid() {
			found = true
			value = v.Interface()
		}
	}
	return
}

// ListItemValuesUnique retrieves and returns the unique elements of all struct/map with key <key>.
// Note that the parameter <list> should be type of slice which contains elements of map or struct,
// or else it returns an empty slice.
func ListItemValuesUnique(list interface{}, key string, subKey ...interface{}) []interface{} {
	values := ListItemValues(list, key, subKey...)
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
