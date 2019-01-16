// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gmd5 provides useful API for MD5 encryption/decryption algorithms.
package gmd5

import (
    "crypto/md5"
    "fmt"
    "os"
    "io"
    "gitee.com/johng/gf/g/util/gconv"
)

// 将任意类型的变量进行md5摘要(注意map等非排序变量造成的不同结果)
func Encrypt(v interface{}) string {
    h := md5.New()
    h.Write([]byte(gconv.Bytes(v)))
    return fmt.Sprintf("%x", h.Sum(nil))
}

// 将字符串进行MD5哈希摘要计算
func EncryptString(v string) string {
    h := md5.New()
    h.Write([]byte(v))
    return fmt.Sprintf("%x", h.Sum(nil))
}

// 将文件内容进行MD5哈希摘要计算
func EncryptFile(path string) string {
    f, e := os.Open(path)
    if e != nil {
        return ""
    }
    defer f.Close()
    h := md5.New()
    _, e = io.Copy(h, f)
    if e != nil {
        return ""
    }
    return fmt.Sprintf("%x", h.Sum(nil))
}
