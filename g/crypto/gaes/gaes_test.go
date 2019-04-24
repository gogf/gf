// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gaes_test

import (
	"testing"

	"github.com/gogf/gf/g/crypto/gaes"
	"github.com/gogf/gf/g/test/gtest"
)

var (
	content = []byte("pibigstar")
	// iv 长度必须等于blockSize，只能为16
	iv     = []byte("Hello My GoFrame")
	key_16 = []byte("1234567891234567")
	key_24 = []byte("123456789123456789123456")
	key_32 = []byte("12345678912345678912345678912345")
	keys   = []byte("12345678912345678912345678912346")
)

func TestEncrypt(t *testing.T) {
	gtest.Case(t, func() {
		_, err := gaes.Encrypt(content, key_16)
		gtest.Assert(err, nil)
		_, err = gaes.Encrypt(content, key_24)
		gtest.Assert(err, nil)
		_, err = gaes.Encrypt(content, key_32)
		gtest.Assert(err, nil)
		_, err = gaes.Encrypt(content, key_16, iv)
		gtest.Assert(err, nil)
	})
}

func TestDecrypt(t *testing.T) {
	gtest.Case(t, func() {
		encrypt, err := gaes.Encrypt(content, key_16)
		decrypt, err := gaes.Decrypt(encrypt, key_16)
		gtest.Assert(err, nil)
		gtest.Assert(string(decrypt), string(content))

		encrypt, err = gaes.Encrypt(content, key_24)
		decrypt, err = gaes.Decrypt(encrypt, key_24)
		gtest.Assert(err, nil)
		gtest.Assert(string(decrypt), string(content))

		encrypt, err = gaes.Encrypt(content, key_32)
		decrypt, err = gaes.Decrypt(encrypt, key_32)
		gtest.Assert(err, nil)
		gtest.Assert(string(decrypt), string(content))

		encrypt, err = gaes.Encrypt(content, key_32, iv)
		decrypt, err = gaes.Decrypt(encrypt, key_32, iv)
		gtest.Assert(err, nil)
		gtest.Assert(string(decrypt), string(content))

		encrypt, err = gaes.Encrypt(content, key_32, iv)
		decrypt, err = gaes.Decrypt(encrypt, keys, iv)
		gtest.Assert(err, "invalid padding")
	})
}
