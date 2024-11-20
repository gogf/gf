// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"bytes"
	"fmt"
	"io/fs"

	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

func PackFsWithOption(dirfs fs.FS, dirPath string, option Option) ([]byte, error) {
	var buffer = bytes.NewBuffer(nil)
	err := zipFsWriter(dirfs, dirPath, buffer, option)
	if err != nil {
		return nil, err
	}
	// Gzip the data bytes to reduce the size.
	return gcompress.Gzip(buffer.Bytes(), 9)
}

// PackToFileWithOption packs the path specified by `srcPaths` to target file `dstPath`.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
func PackFsToFileWithOption(dirfs fs.FS, dirPath string, option Option) error {
	data, err := PackFsWithOption(dirfs, dirPath, option)
	if err != nil {
		return err
	}
	return gfile.PutBytes(dirPath, data)
}

// PackToGoFileWithOption packs the path specified by `srcPaths` to target go file `goFilePath`
// with given package name `pkgName`.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
func PackFsToGoFileWithOption(dirfs fs.FS, dirPath string, goFilePath, pkgName string, option Option) error {
	data, err := PackFsWithOption(dirfs, dirPath, option)
	if err != nil {
		return err
	}
	return gfile.PutContents(
		goFilePath,
		fmt.Sprintf(gstr.TrimLeft(packedGoSourceTemplate), pkgName, gbase64.EncodeToString(data)),
	)
}
