// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcompress provides kinds of compression algorithms for binary/bytes data.
package gcompress

import (
    "bytes"
    "compress/zlib"
    "compress/gzip"
    "io"
)

// Zlib compresses <data> with zlib algorithm.
func Zlib(data []byte) []byte {
    if data == nil || len(data) < 13 {
        return data
    }
    var in bytes.Buffer
    w   := zlib.NewWriter(&in)
    _, _ = w.Write(data)
    _    = w.Close()
    return in.Bytes()
}

// UnZlib decompresses <data> with zlib algorithm.
func UnZlib(data []byte) []byte {
    if data == nil || len(data) < 13 {
        return data
    }
    b := bytes.NewReader(data)
    var out bytes.Buffer
    r, err := zlib.NewReader(b)
    if err != nil {
        return nil
    }
    _, _ = io.Copy(&out, r)
    return out.Bytes()
}

// Gzip compresses <data> with gzip algorithm.
func Gzip(data []byte) []byte {
	var buf bytes.Buffer
	zip    := gzip.NewWriter(&buf)
	_, err := zip.Write(data)
	if err != nil {
		return nil
	}
	_ = zip.Close()
	return buf.Bytes()
}

// UnGzip decompresses <data> with gzip algorithm.
func UnGzip(data []byte) []byte {
    var buf bytes.Buffer
	content      := bytes.NewReader(data)
	zipData, err := gzip.NewReader(content)
	if err != nil {
		return nil
	}
	_, _ = io.Copy(&buf, zipData)
	_    = zipData.Close()
	return buf.Bytes()
}

