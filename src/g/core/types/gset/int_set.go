package gset

import (
	"fmt"
	"sync"
)

type IntSet struct {
	m sync.RWMutex
	M map[int]struct{}
}

func NewIntSet() *IntSet {
	return &IntSet{M: make(map[int]struct{})}
}

// 设置键
func (this *IntSet) Set(item int) *IntSet {
	if this.Contains(item) {
		return this
	}
	this.m.Lock()
	this.M[item] = struct{}{}
	this.m.Unlock()
	return this
}

// 批量添加设置键
func (this *IntSet) BatchSet(items []int) *IntSet {
    count := len(items)
    if count == 0 {
        return this
    }

    todo := make([]int, 0, count)
    this.m.RLock()
    for i := 0; i < count; i++ {
        _, exists := this.M[items[i]]
        if exists {
            continue
        }

        todo = append(todo, items[i])
    }
    this.m.RUnlock()

    count = len(todo)
    if count == 0 {
        return this
    }

    this.m.Lock()
    for i := 0; i < count; i++ {
        this.M[todo[i]] = struct{}{}
    }
    this.m.Unlock()
    return this
}

// 键是否存在
func (this *IntSet) Contains(item int) bool {
	this.m.RLock()
	_, exists := this.M[item]
	this.m.RUnlock()
	return exists
}

// 大小
func (this *IntSet) Size() int {
	this.m.RLock()
	l := len(this.M)
	this.m.RUnlock()
	return l
}

// 清空set
func (this *IntSet) Clear() {
	this.m.Lock()
	this.M = make(map[int]struct{})
	this.m.Unlock()
}

// 转换为数组
func (this *IntSet) Slice() []int {
	this.m.RLock()
	ret := make([]int, len(this.M))
	i := 0
	for item := range this.M {
		ret[i] = item
		i++
	}

	this.m.RUnlock()
	return ret
}

// 转换为字符串
func (this *IntSet) String() string {
	s := this.Slice()
	return fmt.Sprint(s)
}
