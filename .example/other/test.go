package main

import (
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
)

// bytesToHexString converts binary content to hex string content.
func bytesToHexStr(b []byte) string {
	dst := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(dst, b)
	return gconv.UnsafeBytesToStr(dst)
}

func main() {
	b := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		b[i] = byte(grand.N(0, 255))
	}

	fmt.Println(bytesToHexStr(b))
	fmt.Println(len(b))
	fmt.Println(len(bytesToHexStr(b)))
	fmt.Println(gbase64.EncodeToString(b))
	fmt.Println(len(gbase64.EncodeToString(b)))
}
