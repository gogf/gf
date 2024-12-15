// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres/internal/fs_res"
)

// Pack packs the path specified by `srcPaths` into bytes.
// The unnecessary parameter `keyPrefix` indicates the prefix for each file
// packed into the result bytes.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
//
// Deprecated: use PackWithOption instead.
func Pack(srcPaths string, keyPrefix ...string) ([]byte, error) {
	option := PackOption{}
	if len(keyPrefix) > 0 && keyPrefix[0] != "" {
		option.Prefix = keyPrefix[0]
	}
	return PackWithOption(srcPaths, option)
}

// PackToFile packs the path specified by `srcPaths` to target file `dstPath`.
// The unnecessary parameter `keyPrefix` indicates the prefix for each file
// packed into the result bytes.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
//
// Deprecated: use PackToFileWithOption instead.
func PackToFile(srcPaths, dstPath string, keyPrefix ...string) error {
	data, err := Pack(srcPaths, keyPrefix...)
	if err != nil {
		return err
	}
	return gfile.PutBytes(dstPath, data)
}

// PackToGoFile packs the path specified by `srcPaths` to target go file `goFilePath`
// with given package name `pkgName`.
//
// The unnecessary parameter `keyPrefix` indicates the prefix for each file
// packed into the result bytes.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
//
// Deprecated: use PackToGoFileWithOption instead.
func PackToGoFile(srcPath, goFilePath, pkgName string, keyPrefix ...string) error {
	option := PackOption{}
	if len(keyPrefix) > 0 && keyPrefix[0] != "" {
		option.Prefix = keyPrefix[0]
	}
	return PackToGoFileWithOption(srcPath, goFilePath, pkgName, option)
}

// PackWithOption packs the path specified by `srcPaths` into bytes.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
func PackWithOption(srcPaths string, option PackOption) ([]byte, error) {
	return fs_res.PackWithOption(srcPaths, option)
}

// PackToFileWithOption packs the path specified by `srcPaths` to target file `dstPath`.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
func PackToFileWithOption(srcPaths, dstPath string, option PackOption) error {
	return fs_res.PackToFileWithOption(srcPaths, dstPath, option)
}

// PackToGoFileWithOption packs the path specified by `srcPaths` to target go file `goFilePath`
// with given package name `pkgName`.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
func PackToGoFileWithOption(srcPath, goFilePath, pkgName string, option PackOption) error {
	return fs_res.PackToGoFileWithOption(srcPath, goFilePath, pkgName, option)
}

// Unpack unpacks the content specified by `path` to []*File.
func Unpack(path string) ([]File, error) {
	return fs_res.Unpack(path)
}

// UnpackContent unpacks the content to []File.
func UnpackContent(content string) ([]File, error) {
	return fs_res.UnpackContent(content)
}
