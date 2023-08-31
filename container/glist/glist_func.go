package glist

// ScanGList Convert *glist.List to a Slice of the specified type
func ScanGList[T any](l *List) []T {
	list := make([]T, 0)
	l.Iterator(func(e *Element) bool {
		list = append(list, e.Value.(T))
		return true
	})
	return list
}
