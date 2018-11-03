// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件指针池
package gfpool

import (
    "fmt"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gpool"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/third/github.com/fsnotify/fsnotify"
    "os"
    "sync"
)

// 文件指针池
type Pool struct {
    id         *gtype.Int        // 指针池ID，用以识别指针池是否重建
    pool       *gpool.Pool       // 底层对象池
    inited     *gtype.Bool       // 是否初始化(在执行第一次File方法后初始化)
    watcher    *fsnotify.Watcher // 文件监控对象
    closeChan  chan struct{}     // 关闭事件
    expire     int               // 过期时间
}

// 文件指针池指针
type File struct {
    os.File                // 底层文件指针
    mu     sync.RWMutex     // 互斥锁
    pool   *Pool            // 所属池
    poolid int              // 所属池ID，如果池ID不同表示池已经重建，那么该文件指针也应当销毁，不能重新丢到原有的池中
    flag   int              // 打开标志
    perm   os.FileMode      // 打开权限
    path   string           // 绝对路径
}

// 全局指针池，expire < 0表示不过期，expire = 0表示使用完立即回收，expire > 0表示超时回收
var pools = gmap.NewStringInterfaceMap()

// 获得文件对象，并自动创建指针池(过期时间单位：毫秒)
func Open(path string, flag int, perm os.FileMode, expire...int) (file *File, err error) {
    fpExpire := 0
    if len(expire) > 0 {
        fpExpire = expire[0]
    }
    pool := pools.GetOrSetFuncLock(fmt.Sprintf("%s&%d&%d&%d", path, flag, expire, perm), func() interface{} {
        if p, e := New(path, flag, perm, fpExpire); e == nil {
            return p
        } else {
            err = e
        }
        return nil
    }).(*Pool)
    if pool == nil {
        return nil, err
    }
    return pool.File()
}

func OpenFile(path string, flag int, perm os.FileMode, expire...int) (file *File, err error) {
    return Open(path, flag, perm, expire...)
}

// 创建一个文件指针池，expire = 0表示不过期，expire < 0表示使用完立即回收，expire > 0表示超时回收，默认值为0不过期
// 过期时间单位：毫秒
func New(path string, flag int, perm os.FileMode, expire...int) (*Pool, error) {
    fpExpire := 0
    if len(expire) > 0 {
        fpExpire = expire[0]
    }
    p := &Pool {
        id        : gtype.NewInt(),
        expire    : fpExpire,
        inited    : gtype.NewBool(),
        closeChan : make(chan struct{}),
    }
    p.pool = newFilePool(p, path, flag, perm, fpExpire)
    if watcher, err := fsnotify.NewWatcher(); err == nil {
        p.watcher = watcher
    } else {
        return nil, err
    }
    return p, nil
}

// 创建文件指针池
func newFilePool(p *Pool, path string, flag int, perm os.FileMode, expire int) *gpool.Pool {
    pool := gpool.New(expire, func() (interface{}, error) {
        file, err := os.OpenFile(path, flag, perm)
        if err != nil {
            return nil, err
        }
        return &File{
            File   : *file,
            pool   : p,
            poolid : p.id.Val(),
            flag   : flag,
            perm   : perm,
            path   : path,
        }, nil
    })
    pool.SetExpireFunc(func(i interface{}) {
        i.(*File).File.Close()
    })
    return pool
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

        if !p.inited.Set(true) {
            if err := p.watcher.Add(f.path); err != nil {
                p.inited.Set(false)
            }
            go func() {
                for {
                    select {
                        // 关闭事件
                        case <- p.closeChan:
                            return

                        // 监听事件
                        case ev := <- p.watcher.Events:
                            // 如果文件被删除或者重命名，立即重建指针池
                            if ev.Op & fsnotify.Remove == fsnotify.Remove || ev.Op & fsnotify.Rename == fsnotify.Rename {
                                // 原有的指针都不要了
                                p.id.Add(1)
                                // Clear相当于重建指针池
                                p.pool.Clear()
                                // 为保证原子操作，但又不想加锁，
                                // 这里再执行一次原子Add，将在两次Add中间可能分配出去的文件指针丢弃掉
                                p.id.Add(1)
                            }
                    }
                }
            }()
        }
        return f, nil
    }
}

// 关闭指针池(返回error是标准库io.ReadWriteCloser接口实现)
func (p *Pool) Close() error {
    close(p.closeChan)
    p.pool.Close()
    return nil
}

// 获得底层文件指针(返回error是标准库io.ReadWriteCloser接口实现)
func (f *File) Close() error {
    if f.poolid == f.pool.id.Val() {
        f.pool.pool.Put(f)
    }
    return nil
}
