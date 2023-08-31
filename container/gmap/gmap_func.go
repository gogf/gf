package gmap

// ScanGMap Convert *gmap.Map to a Map of the specified type
func ScanGMap[K comparable, V any](m *Map) map[K]V {
	dst := make(map[K]V)
	m.Iterator(func(k, v interface{}) bool {
		dst[k.(K)] = v.(V)
		return true
	})
	return dst
}

// ScanGIntAnyMap Convert *gmap.IntAnyMap to a Map of the specified type
func ScanGIntAnyMap[V any](m *IntAnyMap) map[int]V {
	m2 := make(map[int]V)
	m.Iterator(func(k int, v interface{}) bool {
		m2[k] = v.(V)
		return true
	})
	return m2
}

// ScanGStrAny Convert *gmap.StrAnyMap to a Map of the specified type
func ScanGStrAny[V any](m *StrAnyMap) map[string]V {
	m2 := make(map[string]V)
	m.Iterator(func(k string, v interface{}) bool {
		m2[k] = v.(V)
		return true
	})
	return m2
}

// ScanGListMap Convert *gmap.ListMap to a Map of the specified type
func ScanGListMap[K comparable, V any](m *ListMap) map[K]V {
	m2 := make(map[K]V)
	m.Iterator(func(key, value interface{}) bool {
		m2[key.(K)] = value.(V)
		return true
	})
	return m2
}

// ScanGTreeMap Convert *gmap.TreeMap to a Map of the specified type
func ScanGTreeMap[K comparable, V any](m *TreeMap) map[K]V {
	m2 := make(map[K]V)
	m.Iterator(func(key, value interface{}) bool {
		m2[key.(K)] = value.(V)
		return true
	})
	return m2
}
