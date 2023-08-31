package gconv

import (
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gset"
)

// ScanGArray Convert *garray.Array to a Slice of the specified type
func ScanGArray[T any](arr *garray.Array) []T {
	dst := make([]T, 0)
	arr.IteratorAsc(func(k int, v interface{}) bool {
		dst[k] = v.(T)
		return true
	})
	return dst
}

// ScanGSortArray Convert *garray.SortedArray to a Slice of the specified type
func ScanGSortArray[T any](arr *garray.SortedArray) []T {
	dst := make([]T, 0)
	arr.IteratorAsc(func(k int, v interface{}) bool {
		dst[k] = v.(T)
		return true
	})
	return dst
}

// ScanGMap Convert *gmap.Map to a Map of the specified type
func ScanGMap[K comparable, V any](m *gmap.Map) map[K]V {
	dst := make(map[K]V)
	m.Iterator(func(k, v interface{}) bool {
		dst[k.(K)] = v.(V)
		return true
	})
	return dst
}

// ScanGIntAnyMap Convert *gmap.IntAnyMap to a Map of the specified type
func ScanGIntAnyMap[V any](m *gmap.IntAnyMap) map[int]V {
	m2 := make(map[int]V)
	m.Iterator(func(k int, v interface{}) bool {
		m2[k] = v.(V)
		return true
	})
	return m2
}

// ScanGStrAny Convert *gmap.StrAnyMap to a Map of the specified type
func ScanGStrAny[V any](m *gmap.StrAnyMap) map[string]V {
	m2 := make(map[string]V)
	m.Iterator(func(k string, v interface{}) bool {
		m2[k] = v.(V)
		return true
	})
	return m2
}

// ScanGListMap Convert *gmap.ListMap to a Map of the specified type
func ScanGListMap[K comparable, V any](m *gmap.ListMap) map[K]V {
	m2 := make(map[K]V)
	m.Iterator(func(key, value interface{}) bool {
		m2[key.(K)] = value.(V)
		return true
	})
	return m2
}

// ScanGTreeMap Convert *gmap.TreeMap to a Map of the specified type
func ScanGTreeMap[K comparable, V any](m *gmap.TreeMap) map[K]V {
	m2 := make(map[K]V)
	m.Iterator(func(key, value interface{}) bool {
		m2[key.(K)] = value.(V)
		return true
	})
	return m2
}

// ScanGSet Convert *gset.Set to a Slice of the specified type
func ScanGSet[T any](set *gset.Set) []T {
	arr := make([]T, 0)
	set.Iterator(func(v interface{}) bool {
		arr = append(arr, v.(T))
		return true
	})
	return arr
}

// ScanGList Convert *glist.List to a Slice of the specified type
func ScanGList[T any](l *glist.List) []T {
	list := make([]T, 0)
	l.Iterator(func(e *glist.Element) bool {
		list = append(list, e.Value.(T))
		return true
	})
	return list
}
