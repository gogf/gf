// Copyright 2017-2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gfile

import (
    "gitee.com/johng/gf/g/os/gfpool"
    "io"
    "io/ioutil"
    "os"
)

// (文本)读取文件内容
func GetContents(path string) string {
    return string(GetBinContents(path))
}

// (二进制)读取文件内容
func GetBinContents(path string) []byte {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil
    }
    return data
}

// 写入文件内容
func putContents(path string, data []byte, flag int, perm os.FileMode) error {
    // 支持目录递归创建
    dir := Dir(path)
    if !Exists(dir) {
        if err := Mkdir(dir); err != nil {
            return err
        }
    }
    // 创建/打开文件，使用文件指针池，默认60秒
    f, err := gfpool.OpenFile(path, flag, perm, 60000)
    if err != nil {
        return err
    }
    defer f.Close()
    if n, err := f.Write(data); err != nil {
        return err
    } else if n < len(data) {
        return io.ErrShortWrite
    }
    return nil
}

// Truncate
func Truncate(path string, size int) error {
    return os.Truncate(path, int64(size))
}

// (文本)写入文件内容
func PutContents(path string, content string) error {
    return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
}

// (文本)追加内容到文件末尾
func PutContentsAppend(path string, content string) error {
    return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}

// (二进制)写入文件内容
func PutBinContents(path string, content []byte) error {
    return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
}

// (二进制)追加内容到文件末尾
func PutBinContentsAppend(path string, content []byte) error {
    return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}

// 获得文件内容下一个指定字节的位置
func GetNextCharOffset(file *os.File, char string, start int64) int64 {
    c := []byte(char)[0]
    b := make([]byte, 1)
    o := start
    for {
        _, err := file.ReadAt(b, o)
        if err != nil {
            return 0
        }
        if b[0] == c {
            return o
        }
        o++
    }
    return 0
}

// 获得文件内容中两个offset之间的内容 [start, end)
func GetBinContentByTwoOffsets(file *os.File, start int64, end int64) []byte {
    buffer := make([]byte, end - start)
    if _, err := file.ReadAt(buffer, start); err != nil {
        return nil
    }
    return buffer
}