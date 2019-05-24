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
    "os"
    "io"
    "github.com/gogf/gf/g/util/gconv"
)

// Encrypt encrypts any type of variable using MD5 algorithms.
// It uses gconv package to convert <v> to its bytes type.
func Encrypt(v interface{}) string {
    h := md5.New()
    h.Write([]byte(gconv.Bytes(v)))
    return fmt.Sprintf("%x", h.Sum(nil))
}


// Deprecated.
func EncryptString(v string) string {
	h := md5.New()
	h.Write([]byte(v))
	return fmt.Sprintf("%x", h.Sum(nil))
}


// EncryptFile encrypts file content of <path> using MD5 algorithms.
func EncryptFile(path string) string {
    f, e := os.Open(path)
    if e != nil {
        return ""
    }
    defer f.Close()
    h   := md5.New()
    _, e = io.Copy(h, f)
    if e != nil {
        return ""
    }
    return fmt.Sprintf("%x", h.Sum(nil))
}
