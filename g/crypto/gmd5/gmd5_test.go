// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gmd5_test

import (
	"os"
	"testing"

	"github.com/gogf/gf/g/crypto/gmd5"
	"github.com/gogf/gf/g/test/gtest"
)

var (
	s = "pibigstar"
	// 根据在线工具生成的md5值
	result = "d175a1ff66aedde64344785f7f7a3df8"
)

func TestEncrypt(t *testing.T) {
	gtest.Case(t, func() {
		encryptString := gmd5.Encrypt(s)
		gtest.Assert(encryptString, result)
	})
}

func TestEncryptString(t *testing.T) {
	gtest.Case(t, func() {
		encryptString := gmd5.EncryptString(s)
		gtest.Assert(encryptString, result)
	})
}

func TestEncryptFile(t *testing.T) {
	path := "test.text"
	gtest.Case(t, func() {
		file, err := os.Create(path)
		gtest.Assert(err, nil)
		defer file.Close()
		file.Write([]byte("Hello Go Frame"))
		encryptFile := gmd5.EncryptFile(path)
		gtest.AssertNE(encryptFile, "")
	})
	os.Remove(path)
}
