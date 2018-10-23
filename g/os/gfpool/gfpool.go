// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件指针池
package gfpool

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
    os.File                // 底层文件指针
    mu     sync.RWMutex     // 互斥锁
    pool   *Pool            // 所属池
    flag   int              // 打开标志
    perm   os.FileMode      // 打开权限
    path   string           // 绝对路径
}

// 全局指针池，expire < 0表示不过期，expire = 0表示使用完立即回收，expire > 0表示超时回收
var pools = gmap.NewStringInterfaceMap()

// 获得文件对象，并自动创建指针池(过期时间单位：毫秒)
func Open(path string, flag int, perm os.FileMode, expire...int) (*File, error) {
    fpExpire := 0
    if len(expire) > 0 {
        fpExpire = expire[0]
    }
    key    := fmt.Sprintf("%s&%d&%d&%d", path, flag, expire, perm)
    result := pools.Get(key)
    if result != nil {
        return result.(*Pool).File()
    }
    pool := New(path, flag, perm, fpExpire)
    pools.Set(key, pool)
    return pool.File()
}

func OpenFile(path string, flag int, perm os.FileMode, expire...int) (*File, error) {
    return Open(path, flag, perm, expire...)
}

// 创建一个文件指针池，expire = 0表示不过期，expire < 0表示使用完立即回收，expire > 0表示超时回收，默认值为0不过期
// 过期时间单位：毫秒
func New(path string, flag int, perm os.FileMode, expire...int) *Pool {
    fpExpire := 0
    if len(expire) > 0 {
        fpExpire = expire[0]
    }
    p     := &Pool {}
    p.pool = gpool.New(fpExpire, func() (interface{}, error) {
        file, err := os.OpenFile(path, flag, perm)
        if err != nil {
            return nil, err
        }
        return &File{
            File : *file,
            pool : p,
            flag : flag,
            perm : perm,
            path : path,
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
        f := v.(*File)
        if f.flag & os.O_CREATE > 0 {
            if _, err := os.Stat(f.path); os.IsNotExist(err) {
                if file, err := os.OpenFile(f.path, f.flag, f.perm); err != nil {
                    return nil, err
                } else {
                    f.File = *file
                }
            }
        }
        if f.flag & os.O_TRUNC > 0 {
            if stat, err := f.Stat(); err == nil {
                if stat.Size() > 0 {
                    if err := f.Truncate(0); err != nil {
                        return nil, err
                    }
                }
            }
        }
        if f.flag & os.O_APPEND > 0 {
            if _, err := f.Seek(0, 2); err != nil {
                return nil, err
            }
        } else {
            f.Seek(0, 0)
        }
        return f, nil
    }
}

// 关闭指针池(返回error是标准库io.ReadWriteCloser接口实现)
func (p *Pool) Close() error {
    p.pool.Close()
    return nil
}

// 获得底层文件指针(返回error是标准库io.ReadWriteCloser接口实现)
func (f *File) Close() error {
    f.pool.pool.Put(f)
    return nil
}
