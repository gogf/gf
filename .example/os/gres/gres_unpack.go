package main

import (
	"github.com/jin502437344/gf/crypto/gaes"
	"github.com/jin502437344/gf/os/gfile"
	"github.com/jin502437344/gf/os/gres"
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
