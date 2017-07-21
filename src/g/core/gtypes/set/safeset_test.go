package set

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestSafeSet(t *testing.T) {
	ss := NewSafeSet()
	// Add & Contains & Size
	ss.Add("a1")
	if !(ss.Contains("a1") && ss.Size() == 1) {
		t.Error("error, Add & Contains & Size")
	}
	// Remove
	ss.Clear()
	ss.Add("a1")
	ss.Add("a2")
	ss.Remove("a1")
	ss.Remove("a3")
	if !(!ss.Contains("a1") && ss.Contains("a2")) {
		t.Error("error, Remove")
	}
	ss.Remove("a2")
	if !(!ss.Contains("a2") && ss.Size() == 0) {
		t.Error("error, Remove")
	}
	// Clear
	ss.Clear()
	ss.Add("a1")
	ss.Clear()
	if !(ss.Size() == 0) {
		t.Error("error, Clear")
	}
	// ToSlice
	ss.Clear()
	ss.Add("a1")
	ss.Add("a2")
	ss.Add("a1")
	ar := ss.ToSlice()
	if !(len(ar) == 2 && ArrayContains(ar, "a1") && ArrayContains(ar, "a2")) {
		t.Error("error, ToSlice")
	}
}

func BenchmarkSafeSetAdd(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	ss := NewSafeSet()
	chars := []string{}
	for i := 0; i < b.N; i++ {
		chars = append(chars, fmt.Sprintf("%d", i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ss.Add(chars[i])
	}
}

func BenchmarkSafeSetConcurrent(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	ss := NewSafeSet()
	chars := []string{}
	for i := 0; i < b.N; i++ {
		chars = append(chars, fmt.Sprintf("%d", i))
	}
	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func(base int) {
			for i := 0; i < each; i++ {
				ss.Add(chars[base+i])
			}
			wg.Done()
		}(i * each)
	}
	wg.Wait()
}

func BenchmarkSafeSetRemove(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	ss := NewSafeSet()
	for i := 0; i < b.N; i++ {
		ss.Add(fmt.Sprintf("%d", i))
	}
	chars := []string{}
	for i := 0; i < b.N; i++ {
		chars = append(chars, fmt.Sprintf("%d", i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ss.Remove(chars[i])
	}
}

func BenchmarkSafeSetRemoveConcurrent(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	ss := NewSafeSet()
	for i := 0; i < b.N; i++ {
		ss.Add(fmt.Sprintf("%d", i))
	}

	chars := []string{}
	for i := 0; i < b.N; i++ {
		chars = append(chars, fmt.Sprintf("%d", i))
	}

	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func(base int) {
			for i := 0; i < each; i++ {
				ss.Remove(chars[i])
			}
			wg.Done()
		}(i * each)
	}
	wg.Wait()
}

func BenchmarkSafeSetContains(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	ss := NewSafeSet()
	for i := 0; i < b.N; i++ {
		ss.Add(fmt.Sprintf("%d", i))
	}
	chars := []string{}
	for i := 0; i < b.N; i++ {
		chars = append(chars, fmt.Sprintf("%d", i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ss.Contains(chars[i])
	}
}

func BenchmarkSafeSetContainsConcurrent(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	ss := NewSafeSet()
	for i := 0; i < b.N; i++ {
		ss.Add(fmt.Sprintf("%d", i))
	}

	chars := []string{}
	for i := 0; i < b.N; i++ {
		chars = append(chars, fmt.Sprintf("%d", i))
	}

	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func(base int) {
			for i := 0; i < each; i++ {
				ss.Contains(chars[i])
			}
			wg.Done()
		}(i * each)
	}
	wg.Wait()
}

func BenchmarkSafeSetSize(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	ss := NewSafeSet()
	for i := 0; i < b.N; i++ {
		ss.Add(fmt.Sprintf("%d", i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ss.Size()
	}
}

func ArrayContains(arr []string, item string) bool {
	for _, key := range arr {
		if key == item {
			return true
		}
	}
	return false
}
