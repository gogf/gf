package rwmutex

import "sync"

// RWMutex的封装，支持对并发安全开启/关闭的控制。
// 但是只能初始化时确定并发安全性，不能在运行时动态修改并发安全特性设置。
type RWMutex struct {
    sync.RWMutex
    safe bool
}

func New(safe...bool) *RWMutex {
    mu := new(RWMutex)
    if len(safe) > 0 {
        mu.safe = safe[0]
    } else {
        mu.safe = true
    }
    return mu
}

func (mu *RWMutex) IsSafe() bool {
    return mu.safe
}

func (mu *RWMutex) Lock(force...bool) {
    if mu.safe || (len(force) > 0 && force[0]) {
        mu.RWMutex.Lock()
    }
}

func (mu *RWMutex) Unlock(force...bool) {
    if mu.safe || (len(force) > 0 && force[0]) {
        mu.RWMutex.Unlock()
    }
}

func (mu *RWMutex) RLock(force...bool) {
    if mu.safe || (len(force) > 0 && force[0]) {
        mu.RWMutex.RLock()
    }
}

func (mu *RWMutex) RUnlock(force...bool) {
    if mu.safe || (len(force) > 0 && force[0]) {
        mu.RWMutex.RUnlock()
    }
}