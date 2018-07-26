// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件指针池
package gfilepool

import (
    "os"
    "sync"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gpool"
    "fmt"
)

// 文件指针池
type Pool struct {
    pool    *gpool.Pool     // 底层对象池
}

// 文件指针池指针
type File struct {
    *os.File                // 底层文件指针
    mu     sync.RWMutex     // 互斥锁
    pool   *Pool            // 所属池
}

// 全局指针池，expire < 0表示不过期，expire = 0表示使用完立即回收，expire > 0表示超时回收
var pools = gmap.NewStringInterfaceMap()

// 获得文件对象，并自动创建指针池
func OpenWithPool(path string, flag int, perm os.FileMode, expire int) (*File, error) {
    key    := fmt.Sprintf("%s&%d&%d&%d", path, flag, expire, perm)
    result := pools.Get(key)
    if result != nil {
        return result.(*Pool).File()
    }
    pool := New(path, flag, perm, expire)
    pools.Set(key, pool)
    return pool.File()
}

// 创建一个文件指针池，expire = 0表示不过期，expire < 0表示使用完立即回收，expire > 0表示超时回收
func New(path string, flag int, perm os.FileMode, expire int) *Pool {
    p     := &Pool {}
    p.pool = gpool.New(expire, func() (interface{}, error) {
        file, err := os.OpenFile(path, flag, perm)
        if err != nil {
            return nil, err
        }
        return &File{
            File : file,
            pool : p,
        }, nil
    })
    p.pool.SetExpireFunc(func(i interface{}) {
        i.(*File).File.Close()
    })
    return p
}

// 获得一个文件打开指针
func (p *Pool) File() (*File, error) {
    if v, err := p.pool.Get(); err != nil {
        return nil, err
    } else {
        return v.(*File), nil
    }
}

// 关闭指针池
func (p *Pool) Close() {
    p.pool.Close()
}

// 获得底层文件指针
func (f *File) Close() {
    f.pool.pool.Put(f)
}
