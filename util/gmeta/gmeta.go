// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmeta provides embedded meta data feature for struct.
package gmeta

import (
	"maps"
	"reflect"
	"sync"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gstructs"
)

// Meta is used as an embedded attribute for struct to enabled metadata feature.
type Meta struct{}

// metaAttributeName is the attribute name of metadata in struct.
const metaAttributeName = "Meta"

// metaType holds the reflection. Type of Meta, used for efficient type comparison.
var metaType = reflect.TypeOf(Meta{})

// metaCache caches parsed metadata by struct type, mapping reflect.Type to map[string]string.
var metaCache sync.Map

// Data retrieves and returns all metadata from `object`.
// The returned map is a copy, so it is safe for the caller to modify it.
func Data(object any) map[string]string {
	data, ok := cachedData(object)
	if !ok {
		return nil
	}
	return maps.Clone(data)
}

// Get retrieves and returns specified metadata by `key` from `object`.
func Get(object any, key string) *gvar.Var {
	data, _ := cachedData(object)
	v, ok := data[key]
	if !ok {
		return nil
	}
	return gvar.New(v)
}

// cachedData returns the metadata of `object` from cache, parsing and caching it on first access
// of its struct type. The returned map is shared and must not be modified.
func cachedData(object any) (map[string]string, bool) {
	reflectType, err := gstructs.StructType(object)
	if err != nil {
		return nil, false
	}
	if data, ok := metaCache.Load(reflectType.Type); ok {
		return data.(map[string]string), true
	}
	data := map[string]string{}
	if field, ok := reflectType.FieldByName(metaAttributeName); ok {
		if field.Type == metaType {
			data = gstructs.ParseTag(string(field.Tag))
		}
	}
	metaCache.Store(reflectType.Type, data)
	return data, true
}
