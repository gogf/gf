// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gsha256_test

import (
	"os"
	"testing"

	"github.com/gogf/gf/v2/crypto/gsha256"
	"github.com/gogf/gf/v2/test/gtest"
)

type user struct {
	name     string
	password string
	age      int
}

func TestEncrypt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		result := "b5568f1b35aeb9eb7528336dea6c211a2cdcec1f333d98141b8adf346717907e"
		s := gsha256.Encrypt("pibigstar")
		t.AssertEQ(s, result)
	})
	gtest.C(t, func(t *gtest.T) {
		user := &user{
			name:     "派大星",
			password: "123456",
			age:      23,
		}
		result := "8e0293ca8e1860ae258a88429d3c14755712059d9562c825557a927718f574f3"
		encrypt := gsha256.Encrypt(user)
		t.AssertEQ(encrypt, result)
	})
}

func TestEncryptFile(t *testing.T) {
	path := "test.text"
	errPath := "err.text"
	gtest.C(t, func(t *gtest.T) {
		result := "8fd86e81f66886d4ef7007c2df565f7f61dce2000d8b67ac7163be547c3115ef"
		file, err := os.Create(path)
		defer os.Remove(path)
		defer file.Close()
		t.AssertNil(err)
		_, _ = file.Write([]byte("Hello Go Frame"))
		encryptFile, err := gsha256.EncryptFile(path)
		t.AssertNil(err)
		t.AssertEQ(encryptFile, result)
		// when the file is not exist,encrypt will return empty string
		errEncrypt, err := gsha256.EncryptFile(errPath)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, "")
	})
}
