// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcache

import (
    "fmt"
    "container/list"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gqueue"
    "gitee.com/johng/gf/g/container/gmap"

)

// LRU算法实现对象，底层双向链表使用了标准库的list.List
type memCacheLru struct {
    data  *gmap.Map     // 记录键名与链表中的位置项指针
    list  *glist.List   // 键名历史记录链表
    queue *gqueue.Queue // 事件队列
}

// 创建LRU管理对象
func newMemCacheLru() *memCacheLru {
    lru := &memCacheLru {
        list  : glist.New(),
        data  : gmap.New(),
        queue : gqueue.New(),
    }
    go lru.StartAutoLoop()
    return lru
}

// 关闭LRU对象
func (lru *memCacheLru) Close() {
    lru.queue.Close()
}

// 删除指定数据项
func (lru *memCacheLru) Remove(key interface{}) {
    if v := lru.data.Get(key); v != nil {
        lru.data.Remove(key)
        lru.list.Remove(v.(*list.Element))
    }
}

// 当前LRU数据大小
func (lru *memCacheLru) Size() int {
    return lru.data.Size()
}

// 添加LRU数据项
func (lru *memCacheLru) Push(key interface{}) {
    lru.queue.Push(key)
}

// 从链表尾删除LRU数据项，并返回对应数据
func (lru *memCacheLru) Pop() interface{} {
    if v := lru.list.PopBack(); v != nil {
        lru.data.Remove(v)
        return v
    }
    return nil
}

// 从链表头打印LRU链表值
func (lru *memCacheLru) Print() {
    for _, v := range lru.list.FrontAll() {
        fmt.Printf("%v ", v)
    }
}

// 异步执行协程，将queue中的数据同步到list中
func (lru *memCacheLru) StartAutoLoop() {
    for {
        if v := lru.queue.Pop(); v != nil {
            // 删除对应链表项
            if v := lru.data.Get(v); v != nil {
                lru.list.Remove(v.(*list.Element))
            }
            // 将数据插入到链表头，并记录对应的链表项到哈希表中，便于检索
            lru.data.Set(v, lru.list.PushFront(v))
        } else {
            break
        }
    }
}