// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gmd5_test

import (
	"os"
	"testing"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	s = "pibigstar"
	// online generated MD5 value
	result = "d175a1ff66aedde64344785f7f7a3df8"
)

type user struct {
	name     string
	password string
	age      int
}

func TestEncrypt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		encryptString, _ := gmd5.Encrypt(s)
		t.Assert(encryptString, result)

		result := "1427562bb29f88a1161590b76398ab72"
		encrypt, _ := gmd5.Encrypt(123456)
		t.AssertEQ(encrypt, result)
	})

	gtest.C(t, func(t *gtest.T) {
		user := &user{
			name:     "派大星",
			password: "123456",
			age:      23,
		}
		result := "70917ebce8bd2f78c736cda63870fb39"
		encrypt, _ := gmd5.Encrypt(user)
		t.AssertEQ(encrypt, result)
	})
}

func TestEncryptString(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		encryptString, _ := gmd5.EncryptString(s)
		t.Assert(encryptString, result)
	})
}

func TestEncryptFile(t *testing.T) {
	path := "test.text"
	errorPath := "err.txt"
	result := "e6e6e1cd41895beebff16d5452dfce12"
	gtest.C(t, func(t *gtest.T) {
		file, err := os.Create(path)
		defer os.Remove(path)
		defer file.Close()
		t.Assert(err, nil)
		_, _ = file.Write([]byte("Hello Go Frame"))
		encryptFile, _ := gmd5.EncryptFile(path)
		t.AssertEQ(encryptFile, result)
		// when the file is not exist,encrypt will return empty string
		errEncrypt, _ := gmd5.EncryptFile(errorPath)
		t.AssertEQ(errEncrypt, "")
	})

}
