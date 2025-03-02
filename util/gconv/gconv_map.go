// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

type recursiveType string

const (
	recursiveTypeAuto recursiveType = "auto"
	recursiveTypeTrue recursiveType = "true"
)

// MapOption specifies the option for map converting.
type MapOption struct {
	// Deep marks doing Map function recursively, which means if the attribute of given converting value
	// is also a struct/*struct, it automatically calls Map function on this attribute converting it to
	// a map[string]any type variable.
	Deep bool

	// OmitEmpty ignores the attributes that has json `omitempty` tag.
	OmitEmpty bool

	// Tags specifies the converted map key name by struct tag name.
	Tags []string
}

// Map converts any variable `value` to map[string]any. If the parameter `value` is not a
// map/struct/*struct type, then the conversion will fail and returns nil.
//
// If `value` is a struct/*struct object, the second parameter `priorityTagAndFieldName` specifies the most priority
// priorityTagAndFieldName that will be detected, otherwise it detects the priorityTagAndFieldName in order of:
// gconv, json, field name.
func Map(value any, option ...MapOption) map[string]any {
	return defaultConverter.doMapConvert(value, recursiveTypeAuto, false, option...)
}

// MapDeep does Map function recursively, which means if the attribute of `value`
// is also a struct/*struct, calls Map function on this attribute converting it to
// a map[string]any type variable.
// Deprecated: used Map instead.
func MapDeep(value any, tags ...string) map[string]any {
	return defaultConverter.doMapConvert(value, recursiveTypeTrue, false, MapOption{
		Deep: true,
		Tags: tags,
	})
}

// MapStrStr converts `value` to map[string]string.
// Note that there might be data copy for this map type converting.
func MapStrStr(value any, option ...MapOption) map[string]string {
	if r, ok := value.(map[string]string); ok {
		return r
	}
	m := Map(value, option...)
	if len(m) > 0 {
		vMap := make(map[string]string, len(m))
		for k, v := range m {
			vMap[k] = String(v)
		}
		return vMap
	}
	return nil
}

// MapStrStrDeep converts `value` to map[string]string recursively.
// Note that there might be data copy for this map type converting.
// Deprecated: used MapStrStr instead.
func MapStrStrDeep(value any, tags ...string) map[string]string {
	if r, ok := value.(map[string]string); ok {
		return r
	}
	m := MapDeep(value, tags...)
	if len(m) > 0 {
		vMap := make(map[string]string, len(m))
		for k, v := range m {
			vMap[k] = String(v)
		}
		return vMap
	}
	return nil
}
