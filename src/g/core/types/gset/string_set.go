package gset

import (
	"fmt"
	"sync"
)

type StringSet struct {
	m sync.RWMutex
	M map[string]struct{}
}

func NewStringSet() *StringSet {
	return &StringSet{M: make(map[string]struct{})}
}

// 设置键
func (this *StringSet) Add(item string) *StringSet {
	if this.Contains(item) {
		return this
	}
	this.m.Lock()
	this.M[item] = struct{}{}
	this.m.Unlock()
	return this
}

// 批量添加设置键
func (this *StringSet) BatchAdd(items []string) *StringSet {
    count := len(items)
    if count == 0 {
        return this
    }

    todo := make([]string, 0, count)
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
func (this *StringSet) Contains(item string) bool {
	this.m.RLock()
	_, exists := this.M[item]
	this.m.RUnlock()
	return exists
}

// 删除键值对
func (this *StringSet) Remove(key string) {
	this.m.Lock()
	delete(this.M, key)
	this.m.Unlock()
}

// 大小
func (this *StringSet) Size() int {
	this.m.RLock()
	l := len(this.M)
	this.m.RUnlock()
	return l
}

// 清空set
func (this *StringSet) Clear() {
	this.m.Lock()
	this.M = make(map[string]struct{})
	this.m.Unlock()
}

// 转换为数组
func (this *StringSet) Slice() []string {
	this.m.RLock()
	ret := make([]string, len(this.M))
	i := 0
	for item := range this.M {
		ret[i] = item
		i++
	}

	this.m.RUnlock()
	return ret
}

// 转换为字符串
func (this *StringSet) String() string {
	s := this.Slice()
	return fmt.Sprint(s)
}
