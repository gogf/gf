// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"archive/zip"
	"bytes"
	"fmt"

	"github.com/gogf/gf/encoding/gcompress"
	"github.com/gogf/gf/internal/utilbytes"
	"github.com/gogf/gf/os/gfile"
)

const (
	gPACKAGE_TEMPLATE = `package %s

import "github.com/gogf/gf/os/gres"

func init() {
	if err := gres.Add(%s); err != nil {
		panic(err)
	}
}
`
)

// Pack packs the path specified by <srcPaths> into bytes.
// The unnecessary parameter <keyPrefix> indicates the prefix for each file
// packed into the result bytes.
//
// Note that parameter <srcPaths> supports multiple paths join with ','.
func Pack(srcPaths string, keyPrefix ...string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	headerPrefix := ""
	if len(keyPrefix) > 0 && keyPrefix[0] != "" {
		headerPrefix = keyPrefix[0]
	}
	err := gcompress.ZipPathWriter(srcPaths, buffer, headerPrefix)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// PackToFile packs the path specified by <srcPaths> to target file <dstPath>.
// The unnecessary parameter <keyPrefix> indicates the prefix for each file
// packed into the result bytes.
//
// Note that parameter <srcPaths> supports multiple paths join with ','.
func PackToFile(srcPaths, dstPath string, keyPrefix ...string) error {
	data, err := Pack(srcPaths, keyPrefix...)
	if err != nil {
		return err
	}
	return gfile.PutBytes(dstPath, data)
}

// PackToGoFile packs the path specified by <srcPaths> to target go file <goFilePath>
// with given package name <pkgName>.
//
// The unnecessary parameter <keyPrefix> indicates the prefix for each file
// packed into the result bytes.
//
// Note that parameter <srcPaths> supports multiple paths join with ','.
func PackToGoFile(srcPath, goFilePath, pkgName string, keyPrefix ...string) error {
	data, err := Pack(srcPath, keyPrefix...)
	if err != nil {
		return err
	}
	return gfile.PutContents(
		goFilePath, fmt.Sprintf(gPACKAGE_TEMPLATE, pkgName, utilbytes.Export(data)),
	)
}

// Unpack unpacks the content specified by <path> to []*File.
func Unpack(path string) ([]*File, error) {
	realPath, err := gfile.Search(path)
	if err != nil {
		return nil, err
	}
	return UnpackContent(gfile.GetBytes(realPath))
}

// UnpackContent unpacks the content to []*File.
func UnpackContent(content []byte) ([]*File, error) {
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return nil, err
	}
	array := make([]*File, len(reader.File))
	for i, file := range reader.File {
		array[i] = &File{file: file}
	}
	return array, nil
}
