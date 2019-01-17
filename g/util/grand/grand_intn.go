// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package grand

import (
    "crypto/rand"
    "encoding/binary"
)

const (
    gBUFFER_SIZE = 10000 // 缓冲区uint32数量大小
)

var (
    bufferChan = make(chan uint32, gBUFFER_SIZE)
)

// 使用缓冲区实现快速的随机数生成
func init() {
    step   := 0
    buffer := make([]byte, 1024)
    go func() {
        for {
            if n, err := rand.Read(buffer); err != nil {
                panic(err)
            } else {
                // 使用缓冲区数据进行一次完整的随机数生成
                for i := 0; i < n - 4; {
                    bufferChan <- binary.LittleEndian.Uint32(buffer[i : i + 4])
                    i ++
                }
                // 充分利用缓冲区数据，随机索引递增
                step = int(buffer[0])%10
                for i := 0; i < n - 4; {
                    bufferChan <- binary.BigEndian.Uint32(buffer[i : i + 4])
                    i += step
                }
            }
        }
    }()
}

// 自定义的 rand.Intn ，绝对随机
func intn (max int) int {
    n := int(<- bufferChan)%max
    if n < 0 {
        return -n
    }
    return n
}
