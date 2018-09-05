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

func (mu *RWMutex) Lock() {
    if mu.safe {
        mu.RWMutex.Lock()
    }
}

func (mu *RWMutex) Unlock() {
    if mu.safe {
        mu.RWMutex.Unlock()
    }
}

func (mu *RWMutex) RLock() {
    if mu.safe {
        mu.RWMutex.RLock()
    }
}

func (mu *RWMutex) RUnlock() {
    if mu.safe {
        mu.RWMutex.RUnlock()
    }
}