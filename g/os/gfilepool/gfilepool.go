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
    "strconv"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/glist"
    "sync/atomic"
    "sync"
)

// 文件指针池
type Pool struct {
    path    string          // 文件绝对路径
    flag    int             // 文件打开标识
    list    *glist.SafeList // 可用/闲置的文件指针链表
    idlemax int             // 闲置最大时间，超过该时间则被系统回收(秒)
    closed  bool            // 连接池是否已关闭
}

// 文件指针池指针
type File struct {
    sync.RWMutex
    pool   *Pool     // 所属池
    file   *os.File  // 指针对象
    expire int64     // 过期时间(秒)
}

// 全局指针池，expire < 0表示不过期，expire = 0表示使用完立即回收，expire > 0表示超时回收
var pools *gmap.StringInterfaceMap = gmap.NewStringInterfaceMap()

// 获得文件对象，并自动创建指针池
func OpenWithPool(path string, flag int, expire int) (*File, error) {
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
        list    : glist.NewSafeList(),
        idlemax : expire,
    }
    // 独立的线程执行过期清理工作
    if expire != -1 {
        go func(p *Pool) {
            // 遍历可用指针列表，判断是否过期
            for !p.closed {
                r := p.list.Front()
                if r != nil && r.Value != nil {
                    f := r.Value.(*File)
                    // 必须小于，中间有1秒的缓存时间，防止同时获取和判断过期时冲突
                    if f.getExpire() < gtime.Second() {
                        f.destroy()
                        p.list.Remove(r)
                        continue
                    }
                }
                time.Sleep(3 * time.Second)
            }
        }(r)
    }
    return r
}

// 获得一个文件打开指针
func (p *Pool) File() (*File, error) {
    if p.list.Len() > 0 {
        // 遍历可用指针列表，返回一个未过期的指针
        for {
            r := p.list.PopBack()
            if r != nil {
                f := r.(*File)
                // 必须大于
                if f.getExpire() > gtime.Second() {
                    return f, nil
                } else if f.file != nil {
                    f.destroy()
                }
            } else {
                break;
            }
        }
    }
    file, err := os.OpenFile(p.path, p.flag, 0666)
    if err != nil {
        return nil, err
    }
    return &File {
        pool : p,
        file : file,
    }, nil
}

// 关闭指针池
func (p *Pool) Close() {
    p.closed = true
}

// 获得底层文件指针
func (f *File) File() *os.File {
    return f.file
}

// 关闭指针链接(软关闭)，放回池中重复使用
func (f *File) Close() {
    f.setExpire(gtime.Second() + int64(f.pool.idlemax))
    f.pool.list.PushFront(f)
}

// 销毁指针
func (f *File) destroy() {
    f.Lock()
    defer f.Unlock()
    if f.file != nil {
        f.file.Close()
        f.file =  nil
    }
}

// 获取指针过期时间
func (f *File) setExpire(expire int64) {
    atomic.StoreInt64(&f.expire, expire)
}

// 获取指针过期时间
func (f *File) getExpire() int64 {
    return atomic.LoadInt64(&f.expire)
}