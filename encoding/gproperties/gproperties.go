// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gproperties provides accessing and converting for .properties content.
package gproperties

import (
	"bytes"
	"sort"
	"strings"

	"github.com/magiconair/properties"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
)

// Decode converts properties format to map.
func Decode(data []byte) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	pr, err := properties.Load(data, properties.UTF8)
	if err != nil || pr == nil {
		err = gerror.Wrapf(err, `Lib magiconair load Properties data failed.`)
		return nil, err
	}
	for _, key := range pr.Keys() {
		value, _ := pr.Get(key)
		setNestedValue(res, key, value)
	}
	return res, nil
}

// Encode converts map to properties format.
func Encode(data map[string]interface{}) (res []byte, err error) {
	pr := properties.NewProperties()
	flattened := flattenMap(data, ".", "")

	keys := make([]string, 0, len(flattened))
	for key := range flattened {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		_, _, err := pr.Set(key, gconv.String(flattened[key]))
		if err != nil {
			err = gerror.Wrapf(err, `Sets the property key to the corresponding value failed.`)
			return nil, err
		}
	}

	var buf bytes.Buffer
	_, err = pr.Write(&buf, properties.UTF8)
	if err != nil {
		err = gerror.Wrapf(err, `Properties Write buf failed.`)
		return nil, err
	}
	return buf.Bytes(), nil
}

// ToJson convert .properties format to JSON.
func ToJson(data []byte) (res []byte, err error) {
	prMap, err := Decode(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(prMap)
}

// setNestedValue sets a value in a nested map based on a dot-separated key.
func setNestedValue(m map[string]interface{}, key string, value interface{}) {
	path := strings.Split(key, ".")
	lastKey := strings.ToLower(path[len(path)-1])
	deepestMap := deepSearch(m, path[0:len(path)-1])
	deepestMap[lastKey] = value
}

// deepSearch scans deep maps, following the key indexes listed in the sequence "path".
func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		if m[k] == nil {
			m[k] = make(map[string]interface{})
		}
		m = m[k].(map[string]interface{})
	}
	return m
}

// flattenMap recursively flattens the given map into a new map.
func flattenMap(m map[string]interface{}, delimiter, prefix string) map[string]interface{} {
	shadow := make(map[string]interface{})
	for k, val := range m {
		fullKey := prefix + k
		switch v := val.(type) {
		case map[string]interface{}:
			for subKey, subVal := range flattenMap(v, delimiter, fullKey+delimiter) {
				shadow[subKey] = subVal
			}
		case map[interface{}]interface{}:
			for subKey, subVal := range flattenMap(gconv.Map(v), delimiter, fullKey+delimiter) {
				shadow[subKey] = subVal
			}
		default:
			shadow[strings.ToLower(fullKey)] = v
		}
	}
	return shadow
}
