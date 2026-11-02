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
	"strings"
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

// metaCache caches metadata by struct type, mapping reflect.Type to *metaCacheItem.
var metaCache sync.Map

// metaCacheItem is the cached metadata of a struct type.
type metaCacheItem struct {
	data map[string]string // data is the parsed metadata, used when tag is empty.
	tag  string            // tag is the raw tag that may contain gtag placeholders and is parsed on every read.
}

// Data retrieves and returns all metadata from `object`.
// The returned map is a copy, so it is safe for the caller to modify it.
func Data(object any) map[string]string {
	item := cachedItem(object)
	if item == nil {
		return nil
	}
	if item.tag != "" {
		return gstructs.ParseTag(item.tag)
	}
	return maps.Clone(item.data)
}

// Get retrieves and returns specified metadata by `key` from `object`.
func Get(object any, key string) *gvar.Var {
	var data map[string]string
	if item := cachedItem(object); item != nil {
		data = item.data
		if item.tag != "" {
			data = gstructs.ParseTag(item.tag)
		}
	}
	v, ok := data[key]
	if !ok {
		return nil
	}
	return gvar.New(v)
}

// cachedItem returns the cached metadata of `object`, creating it on first access of its struct type.
// It returns nil if `object` is not a struct.
func cachedItem(object any) *metaCacheItem {
	reflectType, err := gstructs.StructType(object)
	if err != nil {
		return nil
	}
	if item, ok := metaCache.Load(reflectType.Type); ok {
		return item.(*metaCacheItem)
	}
	item := &metaCacheItem{data: map[string]string{}}
	if field, ok := reflectType.FieldByName(metaAttributeName); ok && field.Type == metaType {
		tag := string(field.Tag)
		// Tags with gtag placeholders (possibly escaped) keep being parsed on read,
		// so later gtag.Set/SetOver calls still take effect.
		if strings.ContainsAny(tag, `{\`) {
			item.tag = tag
		} else {
			item.data = gstructs.ParseTag(tag)
		}
	}
	metaCache.Store(reflectType.Type, item)
	return item
}
