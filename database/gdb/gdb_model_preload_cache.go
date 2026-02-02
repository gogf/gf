// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// Package gdb provides ORM features for popular relationship databases.

package gdb

import (
	"reflect"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/gstructs"
)

// Global cache for struct information using gmap.KVMap
// Use NewKVMapWithChecker to handle typed nil issue for pointer values
var (
	preloadChecker     = func(v *cachedStructInfo) bool { return v == nil }
	preloadStructCache = gmap.NewKVMapWithChecker[reflect.Type, *cachedStructInfo](preloadChecker, true)
)

// cachedStructInfo holds cached struct information.
// Only caches static information (field metadata from gstructs).
// Tag parsing is done every time to maintain flexibility.
type cachedStructInfo struct {
	fields []gstructs.Field // All fields from gstructs (static metadata)
}

// getCachedStructInfo gets or creates cached struct information.
// It uses gmap.KVMap's GetOrSetFuncLock for thread-safe lazy initialization.
func getCachedStructInfo(structType reflect.Type) (*cachedStructInfo, error) {
	// Use GetOrSetFuncLock to ensure thread-safe lazy initialization
	// The function is only executed once per key, even under concurrent access
	cached := preloadStructCache.GetOrSetFuncLock(structType, func() *cachedStructInfo {
		// Use gstructs to get field information (only executed once per type)
		fieldMap, err := gstructs.FieldMap(gstructs.FieldMapInput{
			Pointer:         reflect.New(structType).Interface(),
			RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
		})
		if err != nil {
			return nil // Return nil on error, will not be cached
		}

		// Only cache field information, not parsed tags
		info := &cachedStructInfo{
			fields: make([]gstructs.Field, 0, len(fieldMap)),
		}
		for _, field := range fieldMap {
			info.fields = append(info.fields, field)
		}

		return info
	})

	// Check if the cached value is nil (error case)
	if cached == nil {
		// Re-execute to get the actual error
		fieldMap, err := gstructs.FieldMap(gstructs.FieldMapInput{
			Pointer:         reflect.New(structType).Interface(),
			RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
		})
		if err != nil {
			return nil, err
		}

		// This should not happen, but handle it just in case
		info := &cachedStructInfo{
			fields: make([]gstructs.Field, 0, len(fieldMap)),
		}
		for _, field := range fieldMap {
			info.fields = append(info.fields, field)
		}
		return info, nil
	}

	return cached, nil
}

// ClearPreloadStructCache clears the preload struct cache.
// This is usually not needed unless you're dynamically loading many struct types.
func ClearPreloadStructCache() {
	preloadStructCache.Clear()
}

// GetPreloadCacheStats returns statistics about the preload cache.
// Useful for monitoring and debugging.
func GetPreloadCacheStats() map[string]interface{} {
	totalMemory := 0
	preloadStructCache.Iterator(func(k reflect.Type, v *cachedStructInfo) bool {
		// Rough estimate: each field ~150 bytes
		totalMemory += len(v.fields) * 150
		return true
	})

	return map[string]interface{}{
		"cached_types":    preloadStructCache.Size(),
		"memory_estimate": totalMemory,
	}
}
