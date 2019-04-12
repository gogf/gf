// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gsha1_test

import (
	"os"
	"testing"

	"github.com/gogf/gf/g/crypto/gsha1"
	"github.com/gogf/gf/g/test/gtest"
)

type user struct {
	name     string
	password string
	age      int
}

func TestEncrypt(t *testing.T) {
	gtest.Case(t, func() {
		user := &user{
			name:     "派大星",
			password: "123456",
			age:      23,
		}
		result := "97386736e3ee4adee5ca595c78c12129f6032cad"
		encrypt := gsha1.Encrypt(user)
		gtest.AssertEQ(encrypt, result)
	})
	gtest.Case(t, func() {
		result := "5b4c1c2a08ca85ddd031ef8627414f4cb2620b41"
		s := gsha1.Encrypt("pibigstar")
		gtest.AssertEQ(s, result)
	})
}

func TestEncryptString(t *testing.T) {
	gtest.Case(t, func() {
		result := "5b4c1c2a08ca85ddd031ef8627414f4cb2620b41"
		s := gsha1.EncryptString("pibigstar")
		gtest.AssertEQ(s, result)
	})
}

func TestEncryptFile(t *testing.T) {
	path := "test.text"
	errPath := "err.text"
	gtest.Case(t, func() {
		result := "8b05d3ba24b8d2374b8f5149d9f3fbada14ea984"
		file, err := os.Create(path)
		defer os.Remove(path)
		defer file.Close()
		gtest.Assert(err, nil)
		file.Write([]byte("Hello Go Frame"))
		encryptFile := gsha1.EncryptFile(path)
		gtest.AssertEQ(encryptFile, result)
		// when the file is not exist,encrypt will return empty string
		errEncrypt := gsha1.EncryptFile(errPath)
		gtest.AssertEQ(errEncrypt,"")
	})
}
