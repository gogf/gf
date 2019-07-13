// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmd5 provides useful API for MD5 encryption algorithms.
package gmd5

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"

	"github.com/gogf/gf/g/errors/gerror"
	"github.com/gogf/gf/g/util/gconv"
)

// Encrypt encrypts any type of variable using MD5 algorithms.
// It uses gconv package to convert <v> to its bytes type.
func Encrypt(v interface{}) (encrypt string, err error) {
	h := md5.New()
	if _, err = h.Write([]byte(gconv.Bytes(v))); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// EncryptString is alias of Encrypt.
// Deprecated.
func EncryptString(v string) (encrypt string, err error) {
	return Encrypt(v)
}

// EncryptFile encrypts file content of <path> using MD5 algorithms.
func EncryptFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		err = gerror.Wrap(f.Close(), "file closing error")
	}()
	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
