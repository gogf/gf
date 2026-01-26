// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gogf/gf/v2/text/gstr"
)

// ==================== 字段元数据缓存管理器 ====================
// 设计原则：
// 1. 仅缓存确定性信息（字段索引、类型判断）
// 2. 保留动态查找能力（嵌入字段、大小写不敏感）
// 3. 使用 sync.Map 保证高性能并发读取

// fieldCacheManager 字段缓存管理器
type fieldCacheManager struct {
	cache sync.Map // map[string]*fieldCache
}

// newFieldCacheManager 创建字段缓存管理器
func newFieldCacheManager() *fieldCacheManager {
	return &fieldCacheManager{}
}

// fieldCacheMgr 全局字段缓存管理器实例
var fieldCacheMgr = newFieldCacheManager()

// fieldCache 字段缓存
// 存储可以安全缓存的确定性字段信息，避免在循环内重复反射
type fieldCache struct {
	// 确定性字段索引（可安全缓存）
	bindToAttrIndex   int          // 绑定属性的字段索引（如 UserDetail）
	relationAttrIndex int          // 关系属性的字段索引（如 User，-1表示无）
	isPointerElem     bool         // 数组元素是否为指针类型
	bindToAttrKind    reflect.Kind // 绑定属性的类型

	// 字段名映射（支持大小写不敏感查找）
	fieldNameMap  map[string]string // lowercase -> OriginalName
	fieldIndexMap map[string]int    // FieldName -> Index
}

// getOrBuild 获取或构建缓存（线程安全）
func (m *fieldCacheManager) getOrBuild(
	arrayItemType reflect.Type,
	bindToAttrName string,
	relationAttrName string,
) (*fieldCache, error) {
	// 构建缓存键
	cacheKey := m.buildCacheKey(arrayItemType, bindToAttrName, relationAttrName)

	// 快速路径：缓存命中
	if cached, ok := m.cache.Load(cacheKey); ok {
		return cached.(*fieldCache), nil
	}

	// 慢速路径：构建缓存
	cache, err := m.buildCache(arrayItemType, bindToAttrName, relationAttrName)
	if err != nil {
		return nil, err
	}

	// 存储到缓存（如果并发构建，只有一个会被保存）
	actual, _ := m.cache.LoadOrStore(cacheKey, cache)
	return actual.(*fieldCache), nil
}

// buildCacheKey 构建缓存键
func (m *fieldCacheManager) buildCacheKey(
	typ reflect.Type,
	bindToAttrName string,
	relationAttrName string,
) string {
	// 使用类型的唯一标识 + 字段名组合
	return fmt.Sprintf("%s|%s|%s", typ.String(), bindToAttrName, relationAttrName)
}

// buildCache 构建字段访问缓存
func (m *fieldCacheManager) buildCache(
	arrayItemType reflect.Type,
	bindToAttrName string,
	relationAttrName string,
) (*fieldCache, error) {
	cache := &fieldCache{
		relationAttrIndex: -1, // 默认值
		fieldNameMap:      make(map[string]string),
		fieldIndexMap:     make(map[string]int),
	}

	// 获取实际的结构体类型
	structType := arrayItemType
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
		cache.isPointerElem = true
	}

	if structType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("arrayItemType must be struct or pointer to struct, got: %s", arrayItemType.Kind())
	}

	// 遍历所有字段，构建字段映射
	numField := structType.NumField()
	for i := 0; i < numField; i++ {
		field := structType.Field(i)
		fieldName := field.Name

		cache.fieldIndexMap[fieldName] = i
		cache.fieldNameMap[gstr.ToLower(fieldName)] = fieldName
	}

	// 查找 bindToAttrName 字段索引
	if idx, ok := cache.fieldIndexMap[bindToAttrName]; ok {
		cache.bindToAttrIndex = idx
		field := structType.Field(idx)
		cache.bindToAttrKind = field.Type.Kind()
	} else if originalName, ok := cache.fieldNameMap[gstr.ToLower(bindToAttrName)]; ok {
		// 大小写不敏感查找
		cache.bindToAttrIndex = cache.fieldIndexMap[originalName]
		field := structType.Field(cache.bindToAttrIndex)
		cache.bindToAttrKind = field.Type.Kind()
	} else {
		return nil, fmt.Errorf(`field "%s" not found in type %s`, bindToAttrName, arrayItemType.String())
	}

	// 查找 relationAttrName 字段索引（可选）
	if relationAttrName != "" {
		if idx, ok := cache.fieldIndexMap[relationAttrName]; ok {
			cache.relationAttrIndex = idx
		} else if originalName, ok := cache.fieldNameMap[gstr.ToLower(relationAttrName)]; ok {
			cache.relationAttrIndex = cache.fieldIndexMap[originalName]
		}
		// 注意：如果找不到，保持 -1，表示需要使用 arrayElemValue 本身
	}

	return cache, nil
}

// clear 清空所有缓存（测试或热更新时使用）
func (m *fieldCacheManager) clear() {
	m.cache.Range(func(key, value any) bool {
		m.cache.Delete(key)
		return true
	})
}

// stats 获取缓存统计信息
func (m *fieldCacheManager) stats() (count int) {
	m.cache.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

// ClearFieldCache 清空字段缓存（供外部调用）
// 用于测试或应用热更新场景
func ClearFieldCache() {
	fieldCacheMgr.clear()
}

// GetFieldCacheStats 获取字段缓存统计信息（供监控使用）
func GetFieldCacheStats() int {
	return fieldCacheMgr.stats()
}
