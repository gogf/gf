// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
)

var (
	// DefaultReadBuffer is the buffer size for reading file content.
	DefaultReadBuffer = 1024
)

// GetContents returns the file content of `path` as string.
// It returns en empty string if it fails reading.
func GetContents(path string) string {
	return string(GetBytes(path))
}

// GetBytes returns the file content of `path` as []byte.
// It returns nil if it fails reading.
func GetBytes(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

// putContents puts binary content to file of `path`.
func putContents(path string, data []byte, flag int, perm os.FileMode) error {
	// It supports creating file of `path` recursively.
	dir := Dir(path)
	if !Exists(dir) {
		if err := Mkdir(dir); err != nil {
			return err
		}
	}
	// Opening file with given `flag` and `perm`.
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

// Truncate truncates file of `path` to given size by `size`.
func Truncate(path string, size int) error {
	return os.Truncate(path, int64(size))
}

// PutContents puts string `content` to file of `path`.
// It creates file of `path` recursively if it does not exist.
func PutContents(path string, content string) error {
	return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, DefaultPermOpen)
}

// PutContentsAppend appends string `content` to file of `path`.
// It creates file of `path` recursively if it does not exist.
func PutContentsAppend(path string, content string) error {
	return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_APPEND, DefaultPermOpen)
}

// PutBytes puts binary `content` to file of `path`.
// It creates file of `path` recursively if it does not exist.
func PutBytes(path string, content []byte) error {
	return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, DefaultPermOpen)
}

// PutBytesAppend appends binary `content` to file of `path`.
// It creates file of `path` recursively if it does not exist.
func PutBytesAppend(path string, content []byte) error {
	return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_APPEND, DefaultPermOpen)
}

// GetNextCharOffset returns the file offset for given `char` starting from `start`.
func GetNextCharOffset(reader io.ReaderAt, char byte, start int64) int64 {
	buffer := make([]byte, DefaultReadBuffer)
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

// GetNextCharOffsetByPath returns the file offset for given `char` starting from `start`.
// It opens file of `path` for reading with os.O_RDONLY flag and default perm.
func GetNextCharOffsetByPath(path string, char byte, start int64) int64 {
	if f, err := OpenWithFlagPerm(path, os.O_RDONLY, DefaultPermOpen); err == nil {
		defer f.Close()
		return GetNextCharOffset(f, char, start)
	}
	return -1
}

// GetBytesTilChar returns the contents of the file as []byte
// until the next specified byte `char` position.
//
// Note: Returned value contains the character of the last position.
func GetBytesTilChar(reader io.ReaderAt, char byte, start int64) ([]byte, int64) {
	if offset := GetNextCharOffset(reader, char, start); offset != -1 {
		return GetBytesByTwoOffsets(reader, start, offset+1), offset
	}
	return nil, -1
}

// GetBytesTilCharByPath returns the contents of the file given by `path` as []byte
// until the next specified byte `char` position.
// It opens file of `path` for reading with os.O_RDONLY flag and default perm.
//
// Note: Returned value contains the character of the last position.
func GetBytesTilCharByPath(path string, char byte, start int64) ([]byte, int64) {
	if f, err := OpenWithFlagPerm(path, os.O_RDONLY, DefaultPermOpen); err == nil {
		defer f.Close()
		return GetBytesTilChar(f, char, start)
	}
	return nil, -1
}

// GetBytesByTwoOffsets returns the binary content as []byte from `start` to `end`.
// Note: Returned value does not contain the character of the last position, which means
// it returns content range as [start, end).
func GetBytesByTwoOffsets(reader io.ReaderAt, start int64, end int64) []byte {
	buffer := make([]byte, end-start)
	if _, err := reader.ReadAt(buffer, start); err != nil {
		return nil
	}
	return buffer
}

// GetBytesByTwoOffsetsByPath returns the binary content as []byte from `start` to `end`.
// Note: Returned value does not contain the character of the last position, which means
// it returns content range as [start, end).
// It opens file of `path` for reading with os.O_RDONLY flag and default perm.
func GetBytesByTwoOffsetsByPath(path string, start int64, end int64) []byte {
	if f, err := OpenWithFlagPerm(path, os.O_RDONLY, DefaultPermOpen); err == nil {
		defer f.Close()
		return GetBytesByTwoOffsets(f, start, end)
	}
	return nil
}

// ReadLines reads file content line by line, which is passed to the callback function `callback` as string.
// It matches each line of text, separated by chars '\r' or '\n', stripped any trailing end-of-line marker.
//
// Note that the parameter passed to callback function might be an empty value, and the last non-empty line
// will be passed to callback function `callback` even if it has no newline marker.
func ReadLines(file string, callback func(text string) error) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err = callback(scanner.Text()); err != nil {
			return err
		}
	}
	return nil
}

// ReadLinesBytes reads file content line by line, which is passed to the callback function `callback` as []byte.
// It matches each line of text, separated by chars '\r' or '\n', stripped any trailing end-of-line marker.
//
// Note that the parameter passed to callback function might be an empty value, and the last non-empty line
// will be passed to callback function `callback` even if it has no newline marker.
func ReadLinesBytes(file string, callback func(bytes []byte) error) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err = callback(scanner.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
