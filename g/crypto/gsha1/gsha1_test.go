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
		encrypt := gsha1.Encrypt(user)
		gtest.AssertNE(encrypt, "")
	})
	gtest.Case(t, func() {
		s := gsha1.Encrypt("pibigstar")
		gtest.AssertNE(s, "")
	})
}

func TestEncryptString(t *testing.T) {
	gtest.Case(t, func() {
		s := gsha1.EncryptString("pibigstar")
		gtest.AssertNE(s, "")
	})
}

func TestEncryptFile(t *testing.T) {
	path := "test.text"
	gtest.Case(t, func() {
		file, err := os.Create(path)
		gtest.Assert(err, nil)
		defer file.Close()
		file.Write([]byte("Hello Go Frame"))
		encryptFile := gsha1.EncryptFile(path)
		gtest.AssertNE(encryptFile, "")
	})
	defer os.Remove(path)
}
