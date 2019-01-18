// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gmlock implements a thread-safe memory locker.
// 
// 内存锁.
package gmlock

import "time"

var (
    locker = New()
)

// 内存写锁，如果锁成功返回true，失败则返回false;过期时间单位为秒，默认为0表示不过期
func TryLock(key string, expire...time.Duration) bool {
    return locker.TryLock(key, expire...)
}

// 内存写锁，锁成功返回true，失败时阻塞，当失败时表示有其他写锁存在;过期时间单位为秒，默认为0表示不过期
func Lock(key string, expire...time.Duration) {
    locker.Lock(key, expire...)
}

// 解除基于内存锁的写锁
func Unlock(key string) {
    locker.Unlock(key)
}

// 内存读锁，如果锁成功返回true，失败则返回false; 过期时间单位为秒，默认为0表示不过期
func TryRLock(key string) bool {
    return locker.TryRLock(key)
}

// 内存写锁，锁成功返回true，失败时阻塞，当失败时表示有写锁存在; 过期时间单位为秒，默认为0表示不过期
func RLock(key string) {
    locker.RLock(key)
}

// 解除基于内存锁的读锁
func RUnlock(key string) {
    locker.RUnlock(key)
}