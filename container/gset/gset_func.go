package gset

// ScanGSet Convert *gset.Set to a Slice of the specified type
func ScanGSet[T any](set *Set) []T {
	arr := make([]T, 0)
	set.Iterator(func(v interface{}) bool {
		arr = append(arr, v.(T))
		return true
	})
	return arr
}
