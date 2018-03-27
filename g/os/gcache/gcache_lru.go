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
type _Lru struct {
    list  *glist.List
    data  *gmap.StringInterfaceMap
    queue *gqueue.Queue
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

// 删除指定数据项
func (lru *_Lru) Remove(key string) {
    if v := lru.data.Get(key); v != nil {
        lru.list.Remove(v.(*list.Element))
    }
}

// 添加LRU数据项
func (lru *_Lru) Push(key string) {
    lru.queue.PushBack(key)
}

// 从链表尾删除LRU数据项，并返回对应数据
func (lru *_Lru) Pop() string {
    if v := lru.list.PopBack(); v != nil {
        s := v.(string)
        lru.data.Remove(s)
        return s
    }
    return ""
}

// 从链表头打印LRU链表值
func (lru *_Lru) Print() {
    for _, v := range lru.list.FrontAll() {
        fmt.Printf("%s ", v.(string))
    }
}

// 异步执行协程
func (lru *_Lru) StartAutoLoop() {
    for {
        if v := lru.queue.PopFront(); v != nil {
            s := v.(string)
            // 删除对应链表项
            if v := lru.data.Get(s); v != nil {
                lru.list.Remove(v.(*list.Element))
            }
            // 将数据插入到链表头，并生成新的链表项
            lru.data.Set(s, lru.list.PushFront(s))
        } else {
            break
        }
    }
}