// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsha1 provides useful API for SHA1 encryption algorithms.
package gsha1

import (
    "crypto/sha1"
    "encoding/hex"
    "os"
    "io"
    "github.com/gogf/gf/g/util/gconv"
)

// Encrypt encrypts any type of variable using SHA1 algorithms.
// It uses gconv package to convert <v> to its bytes type.
func Encrypt(v interface{}) string {
    r := sha1.Sum(gconv.Bytes(v))
    return hex.EncodeToString(r[:])
}

// Deprecated.
func EncryptString(s string) string {
	r := sha1.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

// EncryptFile encrypts file content of <path> using SHA1 algorithms.
func EncryptFile(path string) string {
    f, e := os.Open(path)
    if e != nil {
        return ""
    }
    defer f.Close()
    h := sha1.New()
    _, e = io.Copy(h, f)
    if e != nil {
        return ""
    }
    return hex.EncodeToString(h.Sum(nil))
}