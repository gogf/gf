// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsha512 provides useful API for SHA512 encryption algorithms.
package gsha512

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"os"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

// Encrypt encrypts any type of variable using SHA512 algorithms.
// It uses package gconv to convert `v` to its bytes type.
func Encrypt(v any) string {
	bs := sha512.Sum512(gconv.Bytes(v))
	return hex.EncodeToString(bs[:])
}

// EncryptFile encrypts file content of `path` using SHA512 algorithms.
func EncryptFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	if err != nil {
		err = gerror.Wrapf(err, `os.Open failed for name "%s"`, path)
		return "", err
	}
	defer f.Close()
	h := sha512.New()
	_, err = io.Copy(h, f)
	if err != nil {
		err = gerror.Wrap(err, `io.Copy failed`)
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// MustEncryptFile encrypts file content of `path` using the SHA512 algorithm.
// It panics if any error occurs.
func MustEncryptFile(path string) string {
	result, err := EncryptFile(path)
	if err != nil {
		panic(err)
	}
	return result
}
