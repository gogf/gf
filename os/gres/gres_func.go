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

func Pack(srcPath string, keyPrefix ...string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := gcompress.ZipPathWriter(srcPath, buffer, keyPrefix...)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func PackToFile(srcPath, dstPath string, keyPrefix ...string) error {
	data, err := Pack(srcPath, keyPrefix...)
	if err != nil {
		return err
	}
	return gfile.PutBytes(dstPath, data)
}

func PackToGoFile(srcPath, goFilePath, pkgName string, keyPrefix ...string) error {
	data, err := Pack(srcPath, keyPrefix...)
	if err != nil {
		return err
	}
	return gfile.PutContents(
		goFilePath, fmt.Sprintf(gPACKAGE_TEMPLATE, pkgName, utilbytes.Export(data)),
	)
}

func Unpack(path string) ([]*File, error) {
	realPath, err := gfile.Search(path)
	if err != nil {
		return nil, err
	}
	return UnpackContent(gfile.GetBytes(realPath))
}

func UnpackContent(content []byte) ([]*File, error) {
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return nil, err
	}
	array := make([]*File, len(reader.File))
	for i, file := range reader.File {
		array[i] = &File{zipFile: file}
	}
	return array, nil
}
