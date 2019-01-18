// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gmlock

import (
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/os/gtimer"
    "time"
)

// 内存锁管理对象
type Locker struct {
    m *gmap.StringInterfaceMap
}

// 创建一把内存锁, 底层使用的是Mutex
func New() *Locker {
    return &Locker{
        m : gmap.NewStringInterfaceMap(),
    }
}

// 内存写锁，如果锁成功返回true，失败则返回false; 过期时间默认为0表示不过期
func (l *Locker) TryLock(key string, expire...time.Duration) bool {
    return l.doLock(key, l.getExpire(expire...), true)
}

// 内存写锁，锁成功返回true，失败时阻塞，当失败时表示有其他写锁存在;过期时间默认为0表示不过期
func (l *Locker) Lock(key string, expire...time.Duration) {
    l.doLock(key, l.getExpire(expire...), false)
}

// 解除基于内存锁的写锁
func (l *Locker) Unlock(key string) {
    if v := l.m.Get(key); v != nil {
        v.(*Mutex).Unlock()
    }
}

// 内存读锁，如果锁成功返回true，失败则返回false; 过期时间单位为秒，默认为0表示不过期
func (l *Locker) TryRLock(key string) bool {
    return l.doRLock(key, true)
}

// 内存写锁，锁成功返回true，失败时阻塞，当失败时表示有写锁存在; 过期时间单位为秒，默认为0表示不过期
func (l *Locker) RLock(key string) {
    l.doRLock(key, false)
}

// 解除基于内存锁的读锁
func (l *Locker) RUnlock(key string) {
    if v := l.m.Get(key); v != nil {
        v.(*Mutex).RUnlock()
    }
}

// 获得过期时间，没有设置时默认为0不过期
func (l *Locker) getExpire(expire...time.Duration) time.Duration {
    e := time.Duration(0)
    if len(expire) > 0 {
        e = expire[0]
    }
    return e
}

// 内存写锁，当try==true时，如果锁成功返回true，失败则返回false；try==false时，成功时立即返回，否则阻塞等待
func (l *Locker) doLock(key string, expire time.Duration, try bool) bool {
    mu := l.getOrNewMutex(key)
    ok := true
    if try {
        ok = mu.TryLock()
    } else {
        mu.Lock()
    }
    if ok && expire > 0 {
        // 异步goroutine计时处理
        wid := mu.wid.Val()
        gtimer.AddOnce(expire, func() {
            if wid == mu.wid.Val() {
                mu.Unlock()
            }
        })
    }
    return ok
}

// 内存读锁，当try==true时，如果锁成功返回true，失败则返回false；try==false时，成功时立即返回，否则阻塞等待
func (l *Locker) doRLock(key string, try bool) bool {
    mu := l.getOrNewMutex(key)
    ok := true
    if try {
        ok = mu.TryRLock()
    } else {
        mu.RLock()
    }
    return ok
}

// 根据指定key查询或者创建新的Mutex
func (l *Locker) getOrNewMutex(key string) (*Mutex) {
    return l.m.GetOrSetFuncLock(key, func() interface{} {
        return NewMutex()
    }).(*Mutex)
}
