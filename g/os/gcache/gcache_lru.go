// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcache

import (
    "fmt"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gqueue"
    "gitee.com/johng/gf/g/container/gmap"
    "container/list"
)

// LRU算法实现对象，底层双向链表使用了标准库的list.List
type _Lru struct {
    list  *glist.List
    data  *gmap.StringInterfaceMap
    queue *gqueue.Queue
}

// 数据项结构
type _LruItem struct {
    key  string
    size int
}

func newLru() *_Lru {
    lru := &_Lru {
        list  : glist.New(),
        data  : gmap.NewStringInterfaceMap(),
        queue : gqueue.New(),
    }
    go lru.StartAutoLoop()
    return lru
}

// 关闭LRU对象
func (lru *_Lru) Close() {
    lru.queue.Close()
}

// 添加LRU数据项
func (lru *_Lru) Push(key string, size int) {
    lru.queue.PushBack(&_LruItem{key, size})
}

// 从链表尾删除LRU数据项，并返回对应数据
func (lru *_Lru) Pop() *_LruItem {
    if v := lru.list.PopBack(); v != nil {
        item := v.(*_LruItem)
        lru.data.Remove(item.key)
        return item
    }
    return nil
}

// 从链表头打印LRU链表值
func (lru *_Lru) Print() {
    for _, v := range lru.list.FrontAll() {
        fmt.Printf("%s ", v.(*_LruItem).key)
    }
}

// 异步执行协程
func (lru *_Lru) StartAutoLoop() {
    for {
        if v := lru.queue.PopFront(); v != nil {
            item := v.(*_LruItem)
            // 删除对应链表项
            if v := lru.data.Get(item.key); v != nil {
                lru.list.Remove(v.(*list.Element))
            }
            // 将数据插入到链表头，并生成新的链表项
            lru.data.Set(item.key, lru.list.PushFront(item))
        } else {
            break
        }
    }
}