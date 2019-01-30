// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package gset

import (
	"fmt"
	"gitee.com/johng/gf/g/internal/rwmutex"
)

type StringSet struct {
	mu *rwmutex.RWMutex
	m  map[string]struct{}
}

func NewStringSet(unsafe...bool) *StringSet {
	return &StringSet {
		m  : make(map[string]struct{}),
		mu : rwmutex.New(unsafe...),
	}
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (set *StringSet) Iterator(f func (v string) bool) *StringSet {
    set.mu.RLock()
    defer set.mu.RUnlock()
	for k, _ := range set.m {
		if !f(k) {
			break
		}
	}
	return set
}

// 设置键
func (set *StringSet) Add(item string) *StringSet {
	set.mu.Lock()
	set.m[item] = struct{}{}
	set.mu.Unlock()
	return set
}

// 批量添加设置键
func (set *StringSet) BatchAdd(items []string) *StringSet {
	set.mu.Lock()
	for _, item := range items {
        set.m[item] = struct{}{}
    }
	set.mu.Unlock()
    return set
}

// 键是否存在
func (set *StringSet) Contains(item string) bool {
	set.mu.RLock()
	_, exists := set.m[item]
	set.mu.RUnlock()
	return exists
}

// 删除键值对
func (set *StringSet) Remove(key string) *StringSet {
	set.mu.Lock()
	delete(set.m, key)
	set.mu.Unlock()
    return set
}

// 大小
func (set *StringSet) Size() int {
	set.mu.RLock()
	l := len(set.m)
	set.mu.RUnlock()
	return l
}

// 清空set
func (set *StringSet) Clear() *StringSet {
	set.mu.Lock()
	set.m = make(map[string]struct{})
	set.mu.Unlock()
    return set
}

// 转换为数组
func (set *StringSet) Slice() []string {
	set.mu.RLock()
	ret := make([]string, len(set.m))
	i := 0
	for item := range set.m {
		ret[i] = item
		i++
	}

	set.mu.RUnlock()
	return ret
}

// 转换为字符串
func (set *StringSet) String() string {
	return fmt.Sprint(set.Slice())
}

func (set *StringSet) LockFunc(f func(m map[string]struct{})) *StringSet {
	set.mu.Lock(true)
	defer set.mu.Unlock(true)
	f(set.m)
    return set
}

func (set *StringSet) RLockFunc(f func(m map[string]struct{})) *StringSet {
	set.mu.RLock(true)
	defer set.mu.RUnlock(true)
	f(set.m)
    return set
}

// 判断两个集合是否相等.
func (set *StringSet) Equal(other *StringSet) bool {
	if set == other {
		return true
	}
	set.mu.RLock()
	defer set.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	if len(set.m) != len(other.m) {
		return false
	}
	for key := range set.m {
		if _, ok := other.m[key]; !ok {
			return false
		}
	}
	return true
}

// 判断other集合是否为当前集合的子集.
func (set *StringSet) IsSubset(other *StringSet) bool {
	if set == other {
		return true
	}
	set.mu.RLock()
	defer set.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	if len(set.m) != len(other.m) {
		return false
	}
	for key := range other.m {
		if _, ok := set.m[key]; !ok {
			return false
		}
	}
	return true
}

// 并集, 返回新的集合：属于set或属于other的元素为元素的集合.
func (set *StringSet) Union(other *StringSet) (newSet *StringSet) {
	newSet = NewStringSet(true)
	set.mu.RLock()
	defer set.mu.RUnlock()
	if set != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	for k, v := range set.m {
		newSet.m[k] = v
	}
	if set != other {
		for k, v := range other.m {
			newSet.m[k] = v
		}
	}
	return
}

// 差集, 返回新的集合: 属于set且不属于other的元素为元素的集合.
func (set *StringSet) Diff(other *StringSet) (newSet *StringSet) {
	newSet = NewStringSet(true)
	if set == other {
		return newSet
	}
	set.mu.RLock()
	defer set.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()

	for k, v := range set.m {
		if _, ok := other.m[k]; !ok {
			newSet.m[k] = v
		}
	}
	return
}

// 交集, 返回新的集合: 属于set且属于other的元素为元素的集合.
func (set *StringSet) Inter(other *StringSet) (newSet *StringSet) {
	newSet = NewStringSet(true)
	set.mu.RLock()
	defer set.mu.RUnlock()
	if set != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	for k, v := range set.m {
		if _, ok := other.m[k]; ok {
			newSet.m[k] = v
		}
	}
	return
}

// 补集, 返回新的集合: (前提: set应当为full的子集)属于全集full不属于集合set的元素组成的集合.
func (set *StringSet) Complement(full *StringSet) (newSet *StringSet) {
	newSet = NewStringSet(true)
	set.mu.RLock()
	defer set.mu.RUnlock()
	full.mu.RLock()
	defer full.mu.RUnlock()
	for k, v := range full.m {
		if _, ok := set.m[k]; !ok {
			newSet.m[k] = v
		}
	}
	return
}