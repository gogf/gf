package main

import (
    "sync"
)

func main() {
    mu := sync.RWMutex{}
    mu.RLocker()
    mu.Lock()
}