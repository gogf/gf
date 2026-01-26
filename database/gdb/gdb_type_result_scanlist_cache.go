// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// Field metadata cache manager
// Design principle:
// 1. Cache deterministic information only (field index, type judgment)
// 2. Retain dynamic lookup capability (embedded fields, case insensitive)
// 3. Use sync.Map to ensure high performance concurrent access

// fieldCacheManager field cache manager
type fieldCacheManager struct {
	cache sync.Map // map[string]*fieldCacheItem
}

// newFieldCacheManager creates field cache manager
func newFieldCacheManager() *fieldCacheManager {
	return &fieldCacheManager{}
}

// fieldCacheInstance global field cache manager instance
var fieldCacheInstance = newFieldCacheManager()

// fieldCacheItem field cache
// Stores deterministic field information that can be safely cached to avoid repeated reflection in loops
type fieldCacheItem struct {
	// Deterministic field index (can be safely cached)
	bindToAttrIndex   int          // Field index of bound attribute (e.g. UserDetail)
	relationAttrIndex int          // Field index of relation attribute (e.g. User, -1 means none)
	isPointerElem     bool         // Whether array element is pointer type
	bindToAttrKind    reflect.Kind // Type of bound attribute

	// Field name mapping (supports case-insensitive lookup)
	fieldNameMap  map[string]string // lowercase -> OriginalName
	fieldIndexMap map[string]int    // FieldName -> Index
}

// getOrSet gets or sets cache (thread-safe)
func (m *fieldCacheManager) getOrSet(
	arrayItemType reflect.Type,
	bindToAttrName string,
	relationAttrName string,
) (*fieldCacheItem, error) {
	// Build cache key
	cacheKey := m.buildCacheKey(arrayItemType, bindToAttrName, relationAttrName)

	// Fast path: cache hit
	if cached, ok := m.cache.Load(cacheKey); ok {
		return cached.(*fieldCacheItem), nil
	}

	// Slow path: build cache
	cache, err := m.buildCache(arrayItemType, bindToAttrName, relationAttrName)
	if err != nil {
		return nil, err
	}

	// Store to cache (if built concurrently, only one will be saved)
	actual, _ := m.cache.LoadOrStore(cacheKey, cache)
	return actual.(*fieldCacheItem), nil
}

// buildCacheKey builds the cache key
func (m *fieldCacheManager) buildCacheKey(
	typ reflect.Type,
	bindToAttrName string,
	relationAttrName string,
) string {
	// Estimate capacity: type name + two field names + 2 separators
	var builder strings.Builder
	typeName := typ.String()
	builder.Grow(len(typeName) + len(bindToAttrName) + len(relationAttrName) + 2)

	builder.WriteString(typeName)
	builder.WriteByte('|')
	builder.WriteString(bindToAttrName)
	builder.WriteByte('|')
	builder.WriteString(relationAttrName)

	return builder.String()
}

// buildCache builds field access cache
func (m *fieldCacheManager) buildCache(
	arrayItemType reflect.Type,
	bindToAttrName string,
	relationAttrName string,
) (*fieldCacheItem, error) {
	// Get the actual struct type
	structType := arrayItemType
	isPointerElem := false
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
		isPointerElem = true
	}

	if structType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("arrayItemType must be struct or pointer to struct, got: %s", arrayItemType.Kind())
	}

	numField := structType.NumField()
	cache := &fieldCacheItem{
		relationAttrIndex: -1,
		isPointerElem:     isPointerElem,
		fieldNameMap:      make(map[string]string, numField), // Pre-allocate capacity
		fieldIndexMap:     make(map[string]int, numField),    // Pre-allocate capacity
	}

	// Iterate all fields, build field mapping
	for i := 0; i < numField; i++ {
		field := structType.Field(i)
		fieldName := field.Name

		cache.fieldIndexMap[fieldName] = i
		cache.fieldNameMap[strings.ToLower(fieldName)] = fieldName
	}

	// Find bindToAttrName field index
	if idx, ok := cache.fieldIndexMap[bindToAttrName]; ok {
		cache.bindToAttrIndex = idx
		field := structType.Field(idx)
		cache.bindToAttrKind = field.Type.Kind()
	} else {
		// Case-insensitive lookup
		lowerBindName := strings.ToLower(bindToAttrName)
		if originalName, ok := cache.fieldNameMap[lowerBindName]; ok {
			cache.bindToAttrIndex = cache.fieldIndexMap[originalName]
			field := structType.Field(cache.bindToAttrIndex)
			cache.bindToAttrKind = field.Type.Kind()
		} else {
			return nil, fmt.Errorf(`field "%s" not found in type %s`, bindToAttrName, arrayItemType.String())
		}
	}

	// Find relationAttrName field index (optional)
	if relationAttrName != "" {
		if idx, ok := cache.fieldIndexMap[relationAttrName]; ok {
			cache.relationAttrIndex = idx
		} else {
			// Case-insensitive lookup
			lowerRelName := strings.ToLower(relationAttrName)
			if originalName, ok := cache.fieldNameMap[lowerRelName]; ok {
				cache.relationAttrIndex = cache.fieldIndexMap[originalName]
			}
		}
		// Note: if not found, keep -1, indicating that arrayElemValue itself should be used
	}

	return cache, nil
}

// clear clears all cache (used for testing or hot updates)
func (m *fieldCacheManager) clear() {
	m.cache.Clear()
}

// ClearFieldCache clears field cache (for external calls)
// Used for testing or application hot update scenarios
func ClearFieldCache() {
	fieldCacheInstance.clear()
}
