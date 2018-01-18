package main

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

func main() {
    // 传感器状态，0:已下线, 1:开启, 2:关闭， 3:待机
    count  := 100
    status := 1

    // 编码
    bits := make([]gbinary.Bit, 0)
    for i := 0; i < count; i++ {
        bits = gbinary.EncodeBits(bits, uint(status), 2)
    }
    buffer := gbinary.EncodeBitsToBytes(bits)
    fmt.Println("buffer length:", len(buffer))

    // 解码
    alivecount := 0
    sensorbits := gbinary.DecodeBytesToBits(buffer)
    for i := 0; i < len(sensorbits); i += 2 {
        if gbinary.DecodeBits(sensorbits[i:i+2]) == 1 {
            alivecount++
        }
    }
    fmt.Println("alived sensor:", alivecount)
}
