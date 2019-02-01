// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gset provides kinds of concurrent-safe(alternative) sets.
//
// 并发安全集合.
package gset

import (
    "fmt"
    "gitee.com/johng/gf/g/internal/rwmutex"
)

type Set struct {
    mu *rwmutex.RWMutex
    m  map[interface{}]struct{}
}

func New(unsafe...bool) *Set {
    return NewSet(unsafe...)
}

func NewSet(unsafe...bool) *Set {
    return &Set{
        m  : make(map[interface{}]struct{}),
        mu : rwmutex.New(unsafe...),
    }
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (set *Set) Iterator(f func (v interface{}) bool) *Set {
    set.mu.RLock()
    defer set.mu.RUnlock()
    for k, _ := range set.m {
        if !f(k) {
            break
        }
    }
    return set
}

// 添加
func (set *Set) Add(item interface{}) *Set {
    set.mu.Lock()
    set.m[item] = struct{}{}
    set.mu.Unlock()
    return set
}

// 批量添加
func (set *Set) BatchAdd(items []interface{}) *Set {
    set.mu.Lock()
    for _, item := range items {
        set.m[item] = struct{}{}
    }
    set.mu.Unlock()
    return set
}

// 键是否存在
func (set *Set) Contains(item interface{}) bool {
    set.mu.RLock()
    _, exists := set.m[item]
    set.mu.RUnlock()
    return exists
}

// 删除键值对
func (set *Set) Remove(key interface{}) *Set {
    set.mu.Lock()
    delete(set.m, key)
    set.mu.Unlock()
    return set
}

// 大小
func (set *Set) Size() int {
    set.mu.RLock()
    l := len(set.m)
    set.mu.RUnlock()
    return l
}

// 清空set
func (set *Set) Clear() *Set {
    set.mu.Lock()
    set.m = make(map[interface{}]struct{})
    set.mu.Unlock()
    return set
}

// 转换为数组
func (set *Set) Slice() []interface{} {
    set.mu.RLock()
    i   := 0
    ret := make([]interface{}, len(set.m))
    for item := range set.m {
        ret[i] = item
        i++
    }
    set.mu.RUnlock()
    return ret
}

// 转换为字符串
func (set *Set) String() string {
    return fmt.Sprint(set.Slice())
}

// 写锁操作
func (set *Set) LockFunc(f func(m map[interface{}]struct{})) *Set {
    set.mu.Lock(true)
    defer set.mu.Unlock(true)
    f(set.m)
    return set
}

// 读锁操作
func (set *Set) RLockFunc(f func(m map[interface{}]struct{})) *Set {
    set.mu.RLock(true)
    defer set.mu.RUnlock(true)
    f(set.m)
    return set
}

// 判断两个集合是否相等.
func (set *Set) Equal(other *Set) bool {
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

// 判断当前集合是否为other集合的子集.
func (set *Set) IsSubsetOf(other *Set) bool {
    if set == other {
        return true
    }
    set.mu.RLock()
    defer set.mu.RUnlock()
    other.mu.RLock()
    defer other.mu.RUnlock()
    for key := range set.m {
        if _, ok := other.m[key]; !ok {
            return false
        }
    }
    return true
}

// 并集, 返回新的集合：属于set或属于others的元素为元素的集合.
func (set *Set) Union(others ... *Set) (newSet *Set) {
    newSet = NewSet(true)
    set.mu.RLock()
    defer set.mu.RUnlock()
    for _, other := range others {
        if set != other {
            other.mu.RLock()
        }
        for k, v := range set.m {
            newSet.m[k] = v
        }
        if set != other {
            for k, v := range other.m {
                newSet.m[k] = v
            }
        }
        if set != other {
            other.mu.RUnlock()
        }
    }

    return
}

// 差集, 返回新的集合: 属于set且不属于others的元素为元素的集合.
func (set *Set) Diff(others...*Set) (newSet *Set) {
    newSet = NewSet(true)
    set.mu.RLock()
    defer set.mu.RUnlock()
    for _, other := range others {
        if set == other {
            continue
        }
        other.mu.RLock()
        for k, v := range set.m {
            if _, ok := other.m[k]; !ok {
                newSet.m[k] = v
            }
        }
        other.mu.RUnlock()
    }
    return
}

// 交集, 返回新的集合: 属于set且属于others的元素为元素的集合.
func (set *Set) Intersect(others...*Set) (newSet *Set) {
    newSet = NewSet(true)
    set.mu.RLock()
    defer set.mu.RUnlock()
    for _, other := range others {
        if set != other {
            other.mu.RLock()
        }
        for k, v := range set.m {
            if _, ok := other.m[k]; ok {
                newSet.m[k] = v
            }
        }
        if set != other {
            other.mu.RUnlock()
        }
    }
    return
}

// 补集, 返回新的集合: (前提: set应当为full的子集)属于全集full不属于集合set的元素组成的集合.
// 如果给定的full集合不是set的全集时，返回full与set的差集.
func (set *Set) Complement(full *Set) (newSet *Set) {
    newSet = NewSet(true)
    set.mu.RLock()
    defer set.mu.RUnlock()
    if set != full {
        full.mu.RLock()
        defer full.mu.RUnlock()
    }
    for k, v := range full.m {
        if _, ok := set.m[k]; !ok {
            newSet.m[k] = v
        }
    }
    return
}