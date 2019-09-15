package main

import (
	"github.com/gogf/gf/crypto/gaes"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gres"
)

var (
	CryptoKey = []byte("x76cgqt36i9c863bzmotuf8626dxiwu0")
)

func main() {
	binContent := gfile.GetBytes("data.bin")
	binContent, err := gaes.Decrypt(binContent, CryptoKey)
	if err != nil {
		panic(err)
	}
	if err := gres.Add(binContent); err != nil {
		panic(err)
	}
	gres.Dump()
}
