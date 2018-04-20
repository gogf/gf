// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件指针池
package gfilepool

import (
    "os"
    "time"
    "sync"
    "strconv"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gtype"
)

// 文件指针池
type Pool struct {
    path    string          // 文件绝对路径
    flag    int             // 文件打开标识
    list    *glist.List     // 可用/闲置的文件指针链表
    idle    int             // 闲置最大时间，超过该时间则被系统回收(秒)
    closed  *gtype.Bool     // 连接池是否已关闭
}

// 文件指针池指针
type PoolItem struct {
    mu     sync.RWMutex
    pool   *Pool            // 所属池
    file   *os.File         // 指针对象
    expire *gtype.Int64     // 过期时间(秒)
}

// 全局指针池，expire < 0表示不过期，expire = 0表示使用完立即回收，expire > 0表示超时回收
var pools = gmap.NewStringInterfaceMap()

// 获得文件对象，并自动创建指针池
func OpenWithPool(path string, flag int, expire int) (*PoolItem, error) {
    key    := path + strconv.Itoa(flag) + strconv.Itoa(expire)
    result := pools.Get(key)
    if result != nil {
        return result.(*Pool).File()
    }
    pool := New(path, flag, expire)
    pools.Set(key, pool)
    return pool.File()
}

// 创建一个文件指针池，expire < 0表示不过期，expire = 0表示使用完立即回收，expire > 0表示超时回收
func New(path string, flag int, expire int) *Pool {
    r := &Pool {
        path    : path,
        flag    : flag,
        list    : glist.New(),
        idle    : expire,
        closed  : gtype.NewBool(),
    }
    // 独立的线程执行过期清理工作
    if expire != -1 {
        go func(p *Pool) {
            // 遍历可用指针列表，判断是否过期
            for !p.closed.Val() {
                if r := p.list.PopFront(); r != nil {
                    f := r.(*PoolItem)
                    // 必须小于，中间有1秒的缓存时间，防止同时获取和判断过期时冲突
                    if f.expire.Val() < gtime.Second() {
                        f.destroy()
                    } else {
                        // 重新推回去
                        p.list.PushFront(f)
                        break
                    }
                }
                time.Sleep(3 * time.Second)
            }
        }(r)
    }
    return r
}

// 获得一个文件打开指针
func (p *Pool) File() (*PoolItem, error) {
    if p.list.Len() > 0 {
        for {
            // 从队列头依次查找，返回一个未过期的指针
            if r := p.list.PopFront(); r != nil {
                f := r.(*PoolItem)
                if f.expire.Val() > gtime.Second() {
                    return f, nil
                } else if f.file != nil {
                    f.destroy()
                }
            } else {
                break
            }
        }
    }
    file, err := os.OpenFile(p.path, p.flag, 0666)
    if err != nil {
        return nil, err
    }
    return &PoolItem {
        pool   : p,
        file   : file,
        expire : gtype.NewInt64(),
    }, nil
}

// 关闭指针池
func (p *Pool) Close() {
    p.closed.Set(true)
}

// 获得底层文件指针
func (f *PoolItem) File() *os.File {
    return f.file
}

// 关闭指针链接(软关闭)，放回池中重复使用
func (f *PoolItem) Close() {
    f.expire.Set(gtime.Second() + int64(f.pool.idle))
    f.pool.list.PushBack(f)
}

// 销毁指针
func (f *PoolItem) destroy() {
    f.file.Close()
}