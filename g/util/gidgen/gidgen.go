// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// ID生成器.
// 内部采用了通道+缓冲池来实现高效的ID递增生成，
// 非常适合高并发下使用
package gidgen

import "math"

// ID生成器管理对象
type Gen struct {
    ch chan uint
}

// 创建一个ID生成器，并给定ID池大小
func New (bufsize int) *Gen {
    g := &Gen {
        ch : make(chan uint, bufsize),
    }
    go g.startLoop()
    return g
}

// 内部循环，当最大值使用完之后重新从1开始获取
func (g *Gen) startLoop() {
    for {
        // 当ch达到缓冲池大小，会阻塞，只要有线程取出值，再立即填充
        for i := uint(1); i < uint(math.MaxUint64); i++ {
            g.ch <- i
        }
    }
}

// 从池中获取一个ID返回(uint)
func (g *Gen) Uint() uint {
    return <- g.ch
}

// 从池中获取一个ID返回(int)
func (g *Gen) Int() int {
    i := int(<- g.ch & 0x7FFFFFFFFFFFFFFF)
    // 可能是int与uint之间的临界点
    if i == 0 {
        i = int(<- g.ch & 0x7FFFFFFFFFFFFFFF)
    }
    return i
}

