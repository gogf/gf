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
	cache sync.Map // map[fieldCacheKey]*fieldCacheItem
}

// fieldCacheKey is the composite key for field cache
// Using struct instead of string to properly distinguish types with same name
type fieldCacheKey struct {
	typ              reflect.Type // The actual type (not just string representation)
	bindToAttrName   string
	relationAttrName string
}

// newFieldCacheManager creates field cache manager
func newFieldCacheManager() *fieldCacheManager {
	return &fieldCacheManager{}
}

// fieldCacheInstance global field cache manager instance
var fieldCacheInstance = newFieldCacheManager()

// fieldCacheItem field cache
// Stores deterministic field information that can be safely cached to avoid repeated reflection in loops
// Note: withTag and related batch settings are NOT cached because they contain dynamic semantics
// (e.g., where/order conditions) that may differ across struct definitions with the same type name.
type fieldCacheItem struct {
	// Deterministic field index (can be safely cached)
	bindToAttrIndex       int          // Field index of bound attribute (e.g. UserDetail), -1 for embedded fields
	bindToAttrIndexPath   []int        // Full index path for embedded fields (e.g. []int{1, 2})
	relationAttrIndex     int          // Field index of relation attribute (e.g. User, -1 means none)
	relationAttrIndexPath []int        // Full index path for embedded relation attribute
	isPointerElem         bool         // Whether array element is pointer type
	bindToAttrKind        reflect.Kind // Type of bound attribute

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
	// Build cache key using reflect.Type directly
	cacheKey := fieldCacheKey{
		typ:              arrayItemType,
		bindToAttrName:   bindToAttrName,
		relationAttrName: relationAttrName,
	}

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
			// Try to find embedded field using FieldByName (supports anonymous/embedded struct)
			field, ok := structType.FieldByName(bindToAttrName)
			if !ok {
				return nil, fmt.Errorf(`field "%s" not found in type %s`, bindToAttrName, arrayItemType.String())
			}
			// For embedded fields, field.Index contains the full path
			if len(field.Index) == 1 {
				// Direct field (shouldn't happen as we already checked fieldIndexMap)
				cache.bindToAttrIndex = field.Index[0]
			} else {
				// Embedded field - store the full index path
				cache.bindToAttrIndex = -1 // Mark as embedded field
				cache.bindToAttrIndexPath = field.Index
			}
			cache.bindToAttrKind = field.Type.Kind()
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
			} else {
				// Try to find embedded field
				if field, ok := structType.FieldByName(relationAttrName); ok {
					if len(field.Index) == 1 {
						cache.relationAttrIndex = field.Index[0]
					} else {
						// Embedded field
						cache.relationAttrIndex = -1
						cache.relationAttrIndexPath = field.Index
					}
				}
				// Note: if still not found, keep -1, indicating that arrayElemValue itself should be used
			}
		}
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
