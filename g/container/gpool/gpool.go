// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gpool provides a object-reusable concurrent-safe pool.
//
// 对象复用池.
package gpool

import (
    "errors"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/gtimer"
    "time"
)

// 对象池
type Pool struct {
    list       *glist.List                // 可用/闲置的文件指针链表
    closed     *gtype.Bool                // 连接池是否已关闭
    Expire     int64                      // (毫秒)闲置最大时间，超过该时间则被系统回收
    NewFunc    func()(interface{}, error) // 创建对象的方法定义
    ExpireFunc func(interface{})          // 对象的过期销毁方法(当池对象销毁需要执行额外的销毁操作时，需要定义该方法)
                                          // 例如: net.Conn, os.File等对象都需要执行额外关闭操作
}

// 对象池数据项
type poolItem struct {
    expire int64               // (毫秒)过期时间
    value  interface{}         // 对象值
}

// 对象创建方法类型
type NewFunc    func() (interface{}, error)

// 对象过期方法类型
type ExpireFunc func(interface{})

// 创建一个对象池，为保证执行效率，过期时间一旦设定之后无法修改
// expire = 0表示不过期，expire < 0表示使用完立即回收，expire > 0表示超时回收
// 注意过期时间单位为**毫秒**
func New(expire int, newFunc NewFunc, expireFunc...ExpireFunc) *Pool {
    r := &Pool {
        list    : glist.New(),
        closed  : gtype.NewBool(),
        Expire  : int64(expire),
        NewFunc : newFunc,
    }
    if len(expireFunc) > 0 {
        r.ExpireFunc = expireFunc[0]
    }
    gtimer.AddSingleton(time.Second, r.checkExpire)
    return r
}

// 放一个临时对象到池中
func (p *Pool) Put(value interface{}) {
    item := &poolItem {
        value : value,
    }
    if p.Expire == 0 {
        item.expire = 0
    } else {
        item.expire = gtime.Millisecond() + p.Expire
    }
    p.list.PushBack(item)
}

// 清空对象池
func (p *Pool) Clear() {
    p.list.RemoveAll()
}

// 从池中获得一个临时对象
func (p *Pool) Get() (interface{}, error) {
    for !p.closed.Val() {
        if r := p.list.PopFront(); r != nil {
            f := r.(*poolItem)
            if f.expire == 0 || f.expire > gtime.Millisecond() {
                return f.value, nil
            }
        } else {
            break
        }
    }
    if p.NewFunc != nil {
        return p.NewFunc()
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
func (p *Pool) checkExpire() {
    if p.closed.Val() {
        gtimer.Exit()
    }
    for {
        if r := p.list.PopFront(); r != nil {
            item := r.(*poolItem)
            if item.expire == 0 || item.expire > gtime.Millisecond() {
                p.list.PushFront(item)
                break
            }
            if p.ExpireFunc != nil {
                p.ExpireFunc(item.value)
            }
        } else {
            break
        }
    }
}