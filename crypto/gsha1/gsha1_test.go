// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gsha1_test

import (
	"os"
	"testing"

	"github.com/gogf/gf/crypto/gsha1"
	"github.com/gogf/gf/test/gtest"
)

type user struct {
	name     string
	password string
	age      int
}

func TestEncrypt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		user := &user{
			name:     "派大星",
			password: "123456",
			age:      23,
		}
		result := "97386736e3ee4adee5ca595c78c12129f6032cad"
		encrypt := gsha1.Encrypt(user)
		t.AssertEQ(encrypt, result)
	})
	gtest.C(t, func(t *gtest.T) {
		result := "5b4c1c2a08ca85ddd031ef8627414f4cb2620b41"
		s := gsha1.Encrypt("pibigstar")
		t.AssertEQ(s, result)
	})
}

func TestEncryptFile(t *testing.T) {
	path := "test.text"
	errPath := "err.text"
	gtest.C(t, func(t *gtest.T) {
		result := "8b05d3ba24b8d2374b8f5149d9f3fbada14ea984"
		file, err := os.Create(path)
		defer os.Remove(path)
		defer file.Close()
		t.Assert(err, nil)
		_, _ = file.Write([]byte("Hello Go Frame"))
		encryptFile, _ := gsha1.EncryptFile(path)
		t.AssertEQ(encryptFile, result)
		// when the file is not exist,encrypt will return empty string
		errEncrypt, _ := gsha1.EncryptFile(errPath)
		t.AssertEQ(errEncrypt, "")
	})
}
