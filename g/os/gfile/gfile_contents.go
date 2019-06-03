// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
    "io"
    "io/ioutil"
    "os"
)

const (
    // Buffer size for reading file content.
    gREAD_BUFFER = 1024
)

// GetContents returns the file content of <path> as string.
// It returns en empty string if it fails reading.
func GetContents(path string) string {
    return string(GetBinContents(path))
}

// GetBinContents returns the file content of <path> as []byte.
// It returns nil if it fails reading.
func GetBinContents(path string) []byte {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil
    }
    return data
}

// putContents puts binary content to file of <path>.
func putContents(path string, data []byte, flag int, perm int) error {
    // It supports creating file of <path> recursively.
    dir := Dir(path)
    if !Exists(dir) {
        if err := Mkdir(dir); err != nil {
            return err
        }
    }
    // Opening file with given <flag> and <perm>.
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

// Truncate truncates file of <path> to given size by <size>.
func Truncate(path string, size int) error {
    return os.Truncate(path, int64(size))
}

// PutContents puts string <content> to file of <path>.
// It creates file of <path> recursively if it does not exist.
func PutContents(path string, content string) error {
    return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, gDEFAULT_PERM)
}

// PutContentsAppend appends string <content> to file of <path>.
// It creates file of <path> recursively if it does not exist.
func PutContentsAppend(path string, content string) error {
    return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_APPEND, gDEFAULT_PERM)
}

// PutBinContents puts binary <content> to file of <path>.
// It creates file of <path> recursively if it does not exist.
func PutBinContents(path string, content []byte) error {
    return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, gDEFAULT_PERM)
}

// PutBinContentsAppend appends binary <content> to file of <path>.
// It creates file of <path> recursively if it does not exist.
func PutBinContentsAppend(path string, content []byte) error {
    return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_APPEND, gDEFAULT_PERM)
}

// GetNextCharOffset returns the file offset for given <char> starting from <start>.
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

// GetNextCharOffsetByPath returns the file offset for given <char> starting from <start>.
// It opens file of <path> for reading with os.O_RDONLY flag and default perm.
func GetNextCharOffsetByPath(path string, char byte, start int64) int64 {
    if f, err := OpenWithFlagPerm(path, os.O_RDONLY, gDEFAULT_PERM); err == nil {
        defer f.Close()
        return GetNextCharOffset(f, char, start)
    }
    return -1
}

// GetBinContentsTilChar returns the contents of the file as []byte 
// until the next specified byte <char> position.
//
// Note: Returned value contains the character of the last position.
func GetBinContentsTilChar(reader io.ReaderAt, char byte, start int64) ([]byte, int64) {
    if offset := GetNextCharOffset(reader, char, start); offset != -1 {
        return GetBinContentsByTwoOffsets(reader, start, offset + 1), offset
    }
    return nil, -1
}

// GetBinContentsTilCharByPath returns the contents of the file given by <path> as []byte 
// until the next specified byte <char> position.
// It opens file of <path> for reading with os.O_RDONLY flag and default perm.
//
// Note: Returned value contains the character of the last position.
func GetBinContentsTilCharByPath(path string, char byte, start int64) ([]byte, int64) {
    if f, err := OpenWithFlagPerm(path, os.O_RDONLY, gDEFAULT_PERM); err == nil {
        defer f.Close()
        return GetBinContentsTilChar(f, char, start)
    }
    return nil, -1
}

// GetBinContentsByTwoOffsets returns the binary content as []byte from <start> to <end>.
// Note: Returned value does not contain the character of the last position, which means
// it returns content range as [start, end).
func GetBinContentsByTwoOffsets(reader io.ReaderAt, start int64, end int64) []byte {
    buffer := make([]byte, end - start)
    if _, err := reader.ReadAt(buffer, start); err != nil {
        return nil
    }
    return buffer
}

// GetBinContentsByTwoOffsetsByPath returns the binary content as []byte from <start> to <end>.
// Note: Returned value does not contain the character of the last position, which means
// it returns content range as [start, end).
// It opens file of <path> for reading with os.O_RDONLY flag and default perm.
func GetBinContentsByTwoOffsetsByPath(path string, start int64, end int64) []byte {
    if f, err := OpenWithFlagPerm(path, os.O_RDONLY, gDEFAULT_PERM); err == nil {
        defer f.Close()
        return GetBinContentsByTwoOffsets(f, start, end)
    }
    return nil
}