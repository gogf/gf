// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gcompress provides kinds of compression algorithms for binary/bytes data.
//
// 数据压缩/解压.
package gcompress

import (
    "bytes"
    "compress/zlib"
    "compress/gzip"
    "io"
)

// 进行zlib压缩
func Zlib(data []byte) []byte {
    if data == nil || len(data) < 13 {
        return data
    }
    var in bytes.Buffer
    w := zlib.NewWriter(&in)
    w.Write(data)
    w.Close()
    return in.Bytes()
}

// 进行zlib解压缩
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
    io.Copy(&out, r)
    return out.Bytes()
}

//做gzip解压缩
func UnGzip(data []byte) []byte {
    var buf bytes.Buffer
	content := bytes.NewReader(data)
	zipdata, err := gzip.NewReader(content)
	if err != nil {
		return nil
	}
	io.Copy(&buf, zipdata)
	zipdata.Close()
	return buf.Bytes()
}

//做gzip压缩
func Gzip(data []byte) []byte {
    var buf bytes.Buffer
	zip := gzip.NewWriter(&buf)
	_, err := zip.Write(data)
	if err != nil {
		return nil
	}
	zip.Close()

	return buf.Bytes()
}