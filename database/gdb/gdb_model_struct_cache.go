// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"reflect"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/gstructs"
)

// Global cache for struct information using gmap.KVMap.
// This cache is shared by both legacy With mode and Preload mode.
// Uses NewKVMapWithChecker to handle typed nil issue for pointer values.
var (
	structInfoChecker = func(v *modelStructCacheItem) bool { return v == nil }
	structInfoCache   = gmap.NewKVMapWithChecker[reflect.Type, *modelStructCacheItem](structInfoChecker, true)
)

// modelStructCacheItem holds cached struct information.
// Only caches static information (field metadata from gstructs).
// Tag parsing is done dynamically to maintain flexibility.
//
// IMPORTANT: The Field.Value in cached fields is a zero-value instance.
// Only use Field.Field (StructField), Field.Type(), Field.Name(), Field.Tag() etc.
// DO NOT use Field.Value.Interface() to get actual values from real instances.
type modelStructCacheItem struct {
	fields []gstructs.Field // All fields from gstructs (static metadata)
}

// buildStructCacheItem creates a modelStructCacheItem from a struct type.
// It extracts field information using gstructs and builds the fields slice.
func buildStructCacheItem(structType reflect.Type) (*modelStructCacheItem, error) {
	// Use gstructs to get field information
	fieldMap, err := gstructs.FieldMap(gstructs.FieldMapInput{
		Pointer:         reflect.New(structType).Interface(),
		RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
	})
	if err != nil {
		return nil, err
	}

	// Build cache item with fields slice
	info := &modelStructCacheItem{
		fields: make([]gstructs.Field, 0, len(fieldMap)),
	}
	for _, field := range fieldMap {
		info.fields = append(info.fields, field)
	}

	return info, nil
}

// getCachedStructInfo gets or creates cached struct information.
// It uses gmap.KVMap's GetOrSetFuncLock for thread-safe lazy initialization.
// This function is used by both legacy With mode and Preload mode.
func getCachedStructInfo(structType reflect.Type) (*modelStructCacheItem, error) {
	// Use GetOrSetFuncLock to ensure thread-safe lazy initialization
	// The function is only executed once per key, even under concurrent access
	cached := structInfoCache.GetOrSetFuncLock(structType, func() *modelStructCacheItem {
		info, err := buildStructCacheItem(structType)
		if err != nil {
			return nil // Return nil on error, will not be cached
		}
		return info
	})

	// Check if the cached value is nil (error case)
	if cached == nil {
		// Re-execute to get the actual error
		return buildStructCacheItem(structType)
	}

	return cached, nil
}

// ClearModelStructCache clears the model struct cache.
// This is usually not needed unless you're dynamically loading many struct types.
func ClearModelStructCache() {
	structInfoCache.Clear()
}
