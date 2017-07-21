package list

import (
	"runtime"
	"sync"
	"testing"
)

type DataItem struct {
	Name string
	Age  int
}

func NewDataItem(name string, age int) *DataItem {
	return &DataItem{Name: name, Age: age}
}

func TestSafeList(t *testing.T) {
	sl := NewSafeList()
	// init test
	length := sl.Len()
	item := sl.PopBack()
	items := sl.PopBackBy(10)
	if !(length == 0 && item == nil && len(items) == 0) {
		t.Error("error, init test")
	}
	// PushFront && Front && Len
	sl.PushFront(NewDataItem("n0", 0))
	sl.PushFront(NewDataItem("n1", 1))
	sl.PushFront(NewDataItem("n2", 2))
	item = sl.Front()
	if !(sl.Len() == 3 && item.(*DataItem).Name == "n2") {
		t.Error("error, PushFront && Front && Len")
	}
	// RemoveAll
	sl.RemoveAll()
	if !(sl.Len() == 0) {
		t.Error("error, RemoveAll")
	}
	// Remove
	sl.RemoveAll()
	e1 := sl.PushFront(NewDataItem("n0", 0))
	sl.PushFront(NewDataItem("n1", 1))
	item = sl.Remove(e1)
	if !(sl.Len() == 1 && item != nil && item.(*DataItem).Name == "n0") {
		t.Error("error, Remove")
	}
	// PopBack
	sl.RemoveAll()
	sl.PushFront(NewDataItem("n0", 0))
	sl.PushFront(NewDataItem("n1", 1))
	sl.PushFront(NewDataItem("n2", 2))
	item = sl.PopBack()
	if !(item.(*DataItem).Name == "n0") {
		t.Error("error, PopBack")
	}
	// PopBackBy
	sl.RemoveAll()
	sl.PushFront(NewDataItem("n0", 0))
	sl.PushFront(NewDataItem("n1", 1))
	sl.PushFront(NewDataItem("n2", 2))
	items = sl.PopBackBy(10)
	if !(sl.Len() == 0 && len(items) == 3 && items[0].(*DataItem).Name == "n0") {
		t.Error("error, PopBackBy")
	}
	// PopBackAll
	sl.RemoveAll()
	sl.PushFront(NewDataItem("n0", 0))
	sl.PushFront(NewDataItem("n1", 1))
	sl.PushFront(NewDataItem("n2", 2))
	items = sl.PopBackAll()
	if !(sl.Len() == 0 && len(items) == 3 && items[0].(*DataItem).Name == "n0") {
		t.Error("error, PopBackAll")
	}
	// FrontAll
	sl.RemoveAll()
	sl.PushFront(NewDataItem("n0", 0))
	sl.PushFront(NewDataItem("n1", 0))
	items = sl.FrontAll()
	if !(sl.Len() == 2 && len(items) == 2 && items[0].(*DataItem).Name == "n1" && items[1].(*DataItem).Name == "n0") {
		t.Error("error, FrontAll")
	}
	// BackAll
	sl.RemoveAll()
	sl.PushFront(NewDataItem("n0", 0))
	sl.PushFront(NewDataItem("n1", 0))
	items = sl.BackAll()
	if !(sl.Len() == 2 && len(items) == 2 && items[0].(*DataItem).Name == "n0" && items[1].(*DataItem).Name == "n1") {
		t.Error("error, BackAll")
	}
	// Other
	sl.RemoveAll()
	sl.PushFront(true)
	item = sl.Front()
	if !(item.(bool) == true) {
		t.Error("error, bool")
	}
	sl.PushFront(1024)
	item = sl.Front()
	if !(item.(int) == 1024) {
		t.Error("error, int")
	}
	sl.PushFront("hello world")
	item = sl.Front()
	if !(item.(string) == "hello world") {
		t.Error("error, string")
	}
	sl.PushFront(*NewDataItem("n0", 0))
	item = sl.Front()
	if !(item.(DataItem).Name == "n0") {
		t.Error("error, DataItem struct")
	}
}

func TestSafeListLimited(t *testing.T) {
	sl := NewSafeListLimited(2)
	// init test
	length := sl.Len()
	item := sl.PopBack()
	items := sl.PopBackBy(10)
	if !(length == 0 && item == nil && len(items) == 0) {
		t.Error("error, init test")
	}
	// PushFront && Front && Len
	sl.PushFront(NewDataItem("n0", 0))
	b1 := sl.PushFront(NewDataItem("n1", 1))
	b2 := sl.PushFront(NewDataItem("n2", 2)) //limited
	item = sl.Front()
	if !(b1 && !b2 && sl.Len() == 2 && item.(*DataItem).Name == "n1") {
		t.Error("error, PushFront && Front && Len")
	}
	// RemoveAll
	sl.RemoveAll()
	if !(sl.Len() == 0) {
		t.Error("error, RemoveAll")
	}
	// PopBack
	sl.RemoveAll()
	sl.PushFront(NewDataItem("n0", 0))
	sl.PushFront(NewDataItem("n1", 1))
	item = sl.PopBack()
	if !(item.(*DataItem).Name == "n0") {
		t.Error("error, PopBack")
	}
	// PopBackBy
	sl.RemoveAll()
	sl.PushFront(NewDataItem("n0", 0))
	sl.PushFront(NewDataItem("n1", 1))
	items = sl.PopBackBy(10)
	if !(sl.Len() == 0 && len(items) == 2 && items[1].(*DataItem).Name == "n1") {
		t.Error("error, PopBackBy")
	}
	// Other
	sl.RemoveAll()
	sl.PushFront(true)
	item = sl.Front()
	if !(item.(bool) == true) {
		t.Error("error, bool")
	}
	sl.PushFront(1024)
	item = sl.Front()
	if !(item.(int) == 1024) {
		t.Error("error, int")
	}

	sl.RemoveAll()
	sl.PushFront("hello world")
	item = sl.Front()
	if !(item.(string) == "hello world") {
		t.Error("error, string")
	}
	sl.PushFront(*NewDataItem("n0", 0))
	item = sl.Front()
	if !(item.(DataItem).Name == "n0") {
		t.Error("error, DataItem struct")
	}
}

func BenchmarkSafeListPushFront(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	sl := NewSafeList()
	item := NewDataItem("n0", 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sl.PushFront(item)
	}
}
func BenchmarkSafeListPushFrontConcurrent(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	sl := NewSafeList()
	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	item := NewDataItem("n0", 0)
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < each; j++ {
				sl.PushFront(item)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
func BenchmarkSafeListPopBack(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	sl := NewSafeList()
	item := NewDataItem("n0", 0)
	for i := 0; i < b.N; i++ {
		sl.PushFront(item)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sl.PopBack()
	}
}
func BenchmarkSafeListPopBackConcurrent(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	sl := NewSafeList()
	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	item := NewDataItem("n0", 0)
	for i := 0; i < b.N; i++ {
		sl.PushFront(item)
	}
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < each; j++ {
				sl.PopBack()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkSafeListLimitedPushFront(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	sl := NewSafeListLimited(b.N)
	item := NewDataItem("n0", 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sl.PushFront(item)
	}
}
func BenchmarkSafeListLimitedPushFrontConcurrent(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	sl := NewSafeListLimited(b.N)
	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	item := NewDataItem("n0", 0)
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < each; j++ {
				sl.PushFront(item)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
func BenchmarkSafeListLimitedPopBack(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	sl := NewSafeListLimited(b.N)
	item := NewDataItem("n0", 0)
	for i := 0; i < b.N; i++ {
		sl.PushFront(item)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sl.PopBack()
	}
}
func BenchmarkSafeListLimitedPopBackConcurrent(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	sl := NewSafeListLimited(b.N)
	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	item := NewDataItem("n0", 0)
	for i := 0; i < b.N; i++ {
		sl.PushFront(item)
	}
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < each; j++ {
				sl.PopBack()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
