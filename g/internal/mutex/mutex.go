package mutex

import "sync"

// Mutex的封装，支持对并发安全开启/关闭的控制。
type Mutex struct {
    sync.Mutex
    safe bool
}

func New(unsafe...bool) *Mutex {
    mu := new(Mutex)
    if len(unsafe) > 0 {
        mu.safe = !unsafe[0]
    } else {
        mu.safe = true
    }
    return mu
}

func (mu *Mutex) IsSafe() bool {
    return mu.safe
}

func (mu *Mutex) Lock(force...bool) {
    if mu.safe || (len(force) > 0 && force[0]) {
        mu.Mutex.Lock()
    }
}

func (mu *Mutex) Unlock(force...bool) {
    if mu.safe || (len(force) > 0 && force[0]) {
        mu.Mutex.Unlock()
    }
}
