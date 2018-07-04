// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 对象复用池.
package gpool

import (
    "time"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gtype"
    "errors"
)

// 对象池
type Pool struct {
    list    *glist.List // 可用/闲置的文件指针链表
    idle    int64       // (毫秒)闲置最大时间，超过该时间则被系统回收
    closed  *gtype.Bool // 连接池是否已关闭
    newFunc func()(interface{}, error) // 创建对象的方法定义
}

// 对象池数据项
type poolItem struct {
    expire int64               // (毫秒)过期时间
    value  interface{}         // 对象值
}

// 创建一个对象池，为保证执行效率，过期时间一旦设定之后无法修改
func New(expire int, newFunc...func() (interface{}, error)) *Pool {
    r := &Pool {
        list    : glist.New(),
        idle    : int64(expire),
        closed  : gtype.NewBool(),
    }
    if len(newFunc) > 0 {
        r.newFunc = newFunc[0]
    }
    go r.expireCheckingLoop()
    return r
}

// 放一个临时对象到池中
func (p *Pool) Put(item interface{}) {
    p.list.PushBack(&poolItem{
        expire : gtime.Millisecond() + p.idle,
        value  : item,
    })
}

// 从池中获得一个临时对象
func (p *Pool) Get() (interface{}, error) {
    for !p.closed.Val() {
        if r := p.list.PopFront(); r != nil {
            f := r.(*poolItem)
            if f.expire > gtime.Millisecond() {
                return f.value, nil
            }
        } else {
            break
        }
    }
    if p.newFunc != nil {
        return p.newFunc()
    }
    return nil, errors.New("pool is empty")
}

// 查询当前池中的对象数量
func (p *Pool) Size() int {
    return p.list.Len()
}

// 关闭池
func (p *Pool) Close() {
    p.closed.Set(true)
}

// 超时检测循环
func (p *Pool) expireCheckingLoop() {
    for !p.closed.Val() {
        if r := p.list.PopFront(); r != nil {
            f := r.(*poolItem)
            if f.expire > gtime.Millisecond() {
                p.list.PushFront(f)
                break
            }
        }
        time.Sleep(3 * time.Second)
    }
}