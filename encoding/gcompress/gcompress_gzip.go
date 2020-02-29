// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcompress

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Gzip compresses <data> with gzip algorithm.
// The optional parameter <level> specifies the compression level from
// 1 to 9 which means from none to the best compression.
//
// Note that it returns error if given <level> is invalid.
func Gzip(data []byte, level ...int) ([]byte, error) {
	var writer *gzip.Writer
	var buf bytes.Buffer
	var err error
	if len(level) > 0 {
		writer, err = gzip.NewWriterLevel(&buf, level[0])
		if err != nil {
			return nil, err
		}
	} else {
		writer = gzip.NewWriter(&buf)
	}
	if _, err = writer.Write(data); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnGzip decompresses <data> with gzip algorithm.
func UnGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(&buf, reader); err != nil {
		return nil, err
	}
	if err = reader.Close(); err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}
