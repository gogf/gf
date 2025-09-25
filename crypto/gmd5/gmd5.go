// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
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

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

// Encrypt encrypts any type of variable using MD5 algorithms.
// It uses gconv package to convert `v` to its bytes type.
func Encrypt(data any) (encrypt string, err error) {
	return EncryptBytes(gconv.Bytes(data))
}

// MustEncrypt encrypts any type of variable using MD5 algorithms.
// It uses gconv package to convert `v` to its bytes type.
// It panics if any error occurs.
func MustEncrypt(data any) string {
	result, err := Encrypt(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptBytes encrypts `data` using MD5 algorithms.
func EncryptBytes(data []byte) (encrypt string, err error) {
	h := md5.New()
	if _, err = h.Write(data); err != nil {
		err = gerror.Wrap(err, `hash.Write failed`)
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// MustEncryptBytes encrypts `data` using MD5 algorithms.
// It panics if any error occurs.
func MustEncryptBytes(data []byte) string {
	result, err := EncryptBytes(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptString encrypts string `data` using MD5 algorithms.
func EncryptString(data string) (encrypt string, err error) {
	return EncryptBytes([]byte(data))
}

// MustEncryptString encrypts string `data` using MD5 algorithms.
// It panics if any error occurs.
func MustEncryptString(data string) string {
	result, err := EncryptString(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptFile encrypts file content of `path` using MD5 algorithms.
func EncryptFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	if err != nil {
		err = gerror.Wrapf(err, `os.Open failed for name "%s"`, path)
		return "", err
	}
	defer f.Close()
	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		err = gerror.Wrap(err, `io.Copy failed`)
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// MustEncryptFile encrypts file content of `path` using MD5 algorithms.
// It panics if any error occurs.
func MustEncryptFile(path string) string {
	result, err := EncryptFile(path)
	if err != nil {
		panic(err)
	}
	return result
}
