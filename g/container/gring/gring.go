// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gring provides a concurrent-safe(alternative) ring(circular lists).
//
// 并发安全环.
package gring

import (
    "container/ring"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/internal/rwmutex"
)

type Ring struct {
    mu    *rwmutex.RWMutex // 互斥锁
    ring  *ring.Ring       // 底层环形数据结构
    len   *gtype.Int       // 数据大小(已使用的大小)
    cap   *gtype.Int       // 总长度(分配的环大小，包括未使用的数据项数量)
    dirty *gtype.Bool      // 标记环是否脏了(需要重新计算大小，当环大小发生改变时做标记)
}

func New(cap int, unsafe...bool) *Ring {
    return &Ring {
        mu    : rwmutex.New(unsafe...),
        ring  : ring.New(cap),
        len   : gtype.NewInt(),
        cap   : gtype.NewInt(cap),
        dirty : gtype.NewBool(),
    }
}

// 返回当前环指向的数据项值
func (r *Ring) Val() interface{} {
    r.mu.RLock()
    v := r.ring.Value
    r.mu.RUnlock()
    return v
}

// 返回当前环已有数据项大小
func (r *Ring) Len() int {
    r.checkAndUpdateLenAndCap()
    return r.len.Val()
}

// 返回当前环总大小(包含未使用长度)
func (r *Ring) Cap() int {
    r.checkAndUpdateLenAndCap()
    return r.cap.Val()
}

// 检测并执行len和cap的更新(两者必须一起更新)
func (r *Ring) checkAndUpdateLenAndCap()  {
    if !r.dirty.Val() {
        return
    }
    totalLen := 0
    emptyLen := 0
    if r.ring != nil {
        r.mu.RLock()
        for p := r.ring.Next(); p != r.ring; p = p.Next() {
            if p.Value == nil {
                emptyLen++
            }
            totalLen++
        }
        r.mu.RUnlock()
    }
    r.cap.Set(totalLen)
    r.len.Set(totalLen - emptyLen)
    r.dirty.Set(false)
}

// 当前位置设置数据项值
func (r *Ring) Set(value interface{}) *Ring {
    r.mu.Lock()
    if r.ring.Value == nil {
        r.len.Add(1)
    }
    r.ring.Value = value
    r.mu.Unlock()
    return r
}

// Set & Next
func (r *Ring) Put(value interface{}) *Ring {
    r.mu.Lock()
    if r.ring.Value == nil {
        r.len.Add(1)
    }
    r.ring.Value = value
    r.ring       = r.ring.Next()
    r.mu.Unlock()
    return r
}

// 环往后(n > 0)或者往前(n < 0)移动n个元素
func (r *Ring) Move(n int) *Ring {
    r.mu.Lock()
    r.ring = r.ring.Move(n)
    r.mu.Unlock()
    return r
}

// 环往前移动1个元素
func (r *Ring) Prev() *Ring {
    r.mu.Lock()
    r.ring = r.ring.Prev()
    r.mu.Unlock()
    return r
}

// 环往后移动1个元素
func (r *Ring) Next() *Ring {
    r.mu.Lock()
    r.ring = r.ring.Next()
    r.mu.Unlock()
    return r
}

// 连接两个环，两个环的大小和位置都有可能会发生改变。
// 1、链接将环r与环s连接，使得r.Next()成为s并返回r.Next()的原始值。r一定不能为空。
// 2、如果r和s指向同一个环，则链接它们会从环中移除r和s之间的元素。
//    删除的元素形成子环，结果是对该子环的引用(如果没有删除元素，结果仍然是r.Next()的原始值，而不是nil)。
// 3、如果r和s指向不同的环，则链接它们会创建一个单独的环，并在r之后插入s的元素。 结果指向插入后s的最后一个元素后面的元素。
func (r *Ring) Link(s *Ring) *Ring {
    r.mu.Lock()
    s.mu.Lock()
    r.ring.Link(s.ring)
    s.mu.Unlock()
    r.mu.Unlock()
    r.dirty.Set(true)
    s.dirty.Set(true)
    return r
}

// 删除环中当前位置往后的n个数据项
func (r *Ring) Unlink(n int) *Ring {
    r.mu.Lock()
    r.ring = r.ring.Unlink(n)
    r.dirty.Set(true)
    r.mu.Unlock()
    return r
}

// 读锁遍历，往后只读遍历，回调函数返回true表示继续遍历，否则退出遍历
func (r *Ring) RLockIteratorNext(f func(value interface{}) bool) {
    r.mu.RLock(true)
    defer r.mu.RUnlock(true)
    if !f(r.ring.Value) {
        return
    }
    for p := r.ring.Next(); p != r.ring; p = p.Next() {
        if !f(p.Value) {
            break
        }
    }
}

// 读锁遍历，往前只读遍历，回调函数返回true表示继续遍历，否则退出遍历
func (r *Ring) RLockIteratorPrev(f func(value interface{}) bool) {
    r.mu.RLock(true)
    defer r.mu.RUnlock(true)
    if !f(r.ring.Value) {
        return
    }
    for p := r.ring.Prev(); p != r.ring; p = p.Prev() {
        if !f(p.Value) {
            break
        }
    }
}

// 写锁遍历，往后写遍历，回调函数返回true表示继续遍历，否则退出遍历
func (r *Ring) LockIteratorNext(f func(item *ring.Ring) bool) {
    r.mu.RLock(true)
    defer r.mu.RUnlock(true)
    if !f(r.ring) {
        return
    }
    for p := r.ring.Next(); p != r.ring; p = p.Next() {
        if !f(p) {
            break
        }
    }
}

// 写锁遍历，往前写遍历，回调函数返回true表示继续遍历，否则退出遍历
func (r *Ring) LockIteratorPrev(f func(item *ring.Ring) bool) {
    r.mu.RLock(true)
    defer r.mu.RUnlock(true)
    if !f(r.ring) {
        return
    }
    for p := r.ring.Prev(); p != r.ring; p = p.Prev() {
        if !f(p) {
            break
        }
    }
}

// 从当前位置，往后只读完整遍历，返回非空数据项值构成的数组
func (r *Ring) SliceNext() []interface{} {
    s := make([]interface{}, 0)
    r.mu.RLock()
    if r.ring.Value != nil {
        s = append(s, r.ring.Value)
    }
    for p := r.ring.Next(); p != r.ring; p = p.Next() {
        if p.Value != nil {
            s = append(s, p.Value)
        }
    }
    r.mu.RUnlock()
    return s
}

// 从当前位置，往前只读完整遍历，返回非空数据项值构成的数组
func (r *Ring) SlicePrev() []interface{} {
    s := make([]interface{}, 0)
    r.mu.RLock()
    if r.ring.Value != nil {
        s = append(s, r.ring.Value)
    }
    for p := r.ring.Prev(); p != r.ring; p = p.Prev() {
        if p.Value != nil {
            s = append(s, p.Value)
        }
    }
    r.mu.RUnlock()
    return s
}