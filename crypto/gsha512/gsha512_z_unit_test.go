// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gsha512_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/crypto/gsha512"
	"github.com/gogf/gf/v2/test/gtest"
)

type user struct {
	name     string
	password string
	age      int
}

func TestEncrypt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		result := "c7b81ef31111986759f12df55baf7ea79f9d23557f32656fd271813adc37ab605b793e7c0170180b219a7a66a43a156e04b7563eeab61c4ad04c650b132da269"
		s := gsha512.Encrypt("pibigstar")
		t.AssertEQ(s, result)
	})
	gtest.C(t, func(t *gtest.T) {
		user := &user{
			name:     "派大星",
			password: "123456",
			age:      23,
		}
		result := "fe5e3be3c17e593f89f176833a52b130a6f5d367fd4a65b520cfa6818c4c42f2af133457c75c884554817b36e255130b4164da88c3a1740767153d63a06bdaa5"
		encrypt := gsha512.Encrypt(user)
		t.AssertEQ(encrypt, result)
	})
}

func TestEncryptFile(t *testing.T) {
	path := gtest.DataPath("test.text")
	errPath := gtest.DataPath("err.text")
	gtest.C(t, func(t *gtest.T) {
		result := "2c6df89b4fda8e4c0baa7dc962380c496f1efe6e5c7ffc3bd33175b2e8f8e394716c8ec2e40c70468dd23bbbdc503db480c57b0051705ef5beaa7aec4a9061d5"
		// ensure the testdata directory exists
		dir := filepath.Dir(path)
		err := os.MkdirAll(dir, 0o755)
		t.AssertNil(err)

		file, err := os.Create(path)
		t.AssertNil(err)
		defer func() { _ = os.Remove(path) }()
		defer func() { _ = file.Close() }()
		_, err = file.Write([]byte("Hello Go Frame"))
		t.AssertNil(err)
		encryptFile, err := gsha512.EncryptFile(path)
		t.AssertNil(err)
		t.AssertEQ(encryptFile, result)
		// When the file does not exist, EncryptFile returns an empty string and a non-nil error.
		errEncrypt, err := gsha512.EncryptFile(errPath)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, "")
	})
}
