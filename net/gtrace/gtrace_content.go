// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gtrace provides convenience wrapping functionality for tracing feature using OpenTelemetry.
package gtrace

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/gogf/gf/v2/text/gstr"
	"io"
)

// SafeContent safe trace content
func SafeContent(data []byte) (string, error) {

	var err error
	if IsGzipped(data) {
		if data, err = unCompressGzip(data); err != nil {
			return string(data), err
		}
	}

	content := string(data)
	if gstr.LenRune(content) > MaxContentLogSize() {

		content = gstr.StrLimitRune(content, MaxContentLogSize(), "...")
	}

	return content, nil
}

func unCompressGzip(data []byte) ([]byte, error) {

	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf(`read gzip response err:%+v`, err)
	}
	defer reader.Close()

	uncompressed, err := io.ReadAll(reader)
	if err != nil {
		return data, fmt.Errorf(`get uncompress value err:%+v`, err)
	}

	return uncompressed, nil
}

func IsGzipped(input []byte) bool {

	if len(input) < 2 {
		return false
	}

	return (int64(input[0])&0xff | (int64(int64(input[1])<<8) & 0xff00)) == 0x8b1f
}
