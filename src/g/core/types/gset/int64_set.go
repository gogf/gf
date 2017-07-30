package gset

import (
	"fmt"
	"sync"
)

type Int64Set struct {
	m sync.RWMutex
	M map[int64]struct{}
}

func NewInt64Set() *Int64Set {
	return &Int64Set{M: make(map[int64]struct{})}
}

// 设置键
func (this *Int64Set) Add(item int64) *Int64Set {
	if this.Contains(item) {
		return this
	}
	this.m.Lock()
	this.M[item] = struct{}{}
	this.m.Unlock()
	return this
}

// 批量添加设置键
func (this *Int64Set) BatchAdd(items []int64) *Int64Set {
    count := len(items)
    if count == 0 {
        return this
    }

    todo := make([]int64, 0, count)
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
func (this *Int64Set) Contains(item int64) bool {
	this.m.RLock()
	_, exists := this.M[item]
	this.m.RUnlock()
	return exists
}

// 删除键值对
func (this *Int64Set) Remove(key int64) {
	this.m.Lock()
	delete(this.M, key)
	this.m.Unlock()
}

// 大小
func (this *Int64Set) Size() int {
	this.m.RLock()
	l := len(this.M)
	this.m.RUnlock()
	return l
}

// 清空set
func (this *Int64Set) Clear() {
	this.m.Lock()
	this.M = make(map[int64]struct{})
	this.m.Unlock()
}

// 转换为数组
func (this *Int64Set) Slice() []int64 {
	this.m.RLock()
	ret := make([]int64, len(this.M))
	i := 0
	for item := range this.M {
		ret[i] = item
		i++
	}

	this.m.RUnlock()
	return ret
}

// 转换为字符串
func (this *Int64Set) String() string {
	s := this.Slice()
	return fmt.Sprint(s)
}
