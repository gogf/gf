// Copyright 2017-2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gfile

import (
    "io"
    "io/ioutil"
    "os"
)

const (
    // 方法中涉及到读取的时候的缓冲大小
    gREAD_BUFFER      = 1024
    // 方法中涉及到文件指针池的默认缓存时间(毫秒)
    //gFILE_POOL_EXPIRE = 60000
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
func putContents(path string, data []byte, flag int, perm int) error {
    // 支持目录递归创建
    dir := Dir(path)
    if !Exists(dir) {
        if err := Mkdir(dir); err != nil {
            return err
        }
    }
    // 创建/打开文件
    f, err := OpenWithFlagPerm(path, flag, perm)
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
    return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, gDEFAULT_PERM)
}

// (文本)追加内容到文件末尾
func PutContentsAppend(path string, content string) error {
    return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_APPEND, gDEFAULT_PERM)
}

// (二进制)写入文件内容
func PutBinContents(path string, content []byte) error {
    return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, gDEFAULT_PERM)
}

// (二进制)追加内容到文件末尾
func PutBinContentsAppend(path string, content []byte) error {
    return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_APPEND, gDEFAULT_PERM)
}

// 获得文件内容下一个指定字节的位置
func GetNextCharOffset(reader io.ReaderAt, char byte, start int64) int64 {
    buffer := make([]byte, gREAD_BUFFER)
    offset := start
    for {
        if n, err := reader.ReadAt(buffer, offset); n > 0 {
            for i := 0; i < n; i++ {
                if buffer[i] == char {
                    return int64(i) + offset
                }
            }
            offset += int64(n)
        } else if err != nil {
            break
        }
    }
    return -1
}

// 获得文件内容下一个指定字节的位置
func GetNextCharOffsetByPath(path string, char byte, start int64) int64 {
    if f, err := OpenWithFlagPerm(path, os.O_RDONLY, gDEFAULT_PERM); err == nil {
        defer f.Close()
        return GetNextCharOffset(f, char, start)
    }
    return -1
}

// 获得文件内容直到下一个指定字节的位置(返回值包含该位置字符内容)
func GetBinContentsTilChar(reader io.ReaderAt, char byte, start int64) ([]byte, int64) {
    if offset := GetNextCharOffset(reader, char, start); offset != -1 {
        return GetBinContentsByTwoOffsets(reader, start, offset + 1), offset
    }
    return nil, -1
}

// 获得文件内容直到下一个指定字节的位置(返回值包含该位置字符内容)
func GetBinContentsTilCharByPath(path string, char byte, start int64) ([]byte, int64) {
    if f, err := OpenWithFlagPerm(path, os.O_RDONLY, gDEFAULT_PERM); err == nil {
        defer f.Close()
        return GetBinContentsTilChar(f, char, start)
    }
    return nil, -1
}

// 获得文件内容中两个offset之间的内容 [start, end)
func GetBinContentsByTwoOffsets(reader io.ReaderAt, start int64, end int64) []byte {
    buffer := make([]byte, end - start)
    if _, err := reader.ReadAt(buffer, start); err != nil {
        return nil
    }
    return buffer
}

// 获得文件内容中两个offset之间的内容 [start, end)
func GetBinContentsByTwoOffsetsByPath(path string, start int64, end int64) []byte {
    if f, err := OpenWithFlagPerm(path, os.O_RDONLY, gDEFAULT_PERM); err == nil {
        defer f.Close()
        return GetBinContentsByTwoOffsets(f, start, end)
    }
    return nil
}