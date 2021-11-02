package main

import (
	"fmt"

	"github.com/gogf/gf/v2/encoding/gbinary"
)

func main() {
	// 传感器状态，0:已下线, 1:开启, 2:关闭， 3:待机
	count := 100
	status := 1

	// 网关编码
	bits := make([]gbinary.Bit, 0)
	for i := 0; i < count; i++ {
		bits = gbinary.EncodeBits(bits, status, 2)
	}
	buffer := gbinary.EncodeBitsToBytes(bits)
	fmt.Println("buffer length:", len(buffer))

	/* 上报过程忽略，这里只展示编码/解码示例 */

	// 平台解码
	alivecount := 0
	sensorbits := gbinary.DecodeBytesToBits(buffer)
	for i := 0; i < len(sensorbits); i += 2 {
		if gbinary.DecodeBits(sensorbits[i:i+2]) == 1 {
			alivecount++
		}
	}
	fmt.Println("alived sensor:", alivecount)
}
