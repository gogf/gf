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
	"github.com/gogf/gf/container/gtree"
	"github.com/gogf/gf/encoding/gcompress"
	"github.com/gogf/gf/internal/utilbytes"
	"github.com/gogf/gf/os/gfile"
	"strings"
)

type Resource struct {
	Name string
}

var (
	resTree = gtree.NewBTree(10, func(v1, v2 interface{}) int {
		return strings.Compare(v1.(string), v2.(string))
	})
)

func Add(content []byte) error {
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return err
	}
	for _, file := range reader.File {
		resTree.Set(file.Name, file)
	}
	return nil
}

func Dump() {
	resTree.Iterator(func(key, value interface{}) bool {
		fmt.Printf("%7s %s\n", gfile.FormatSize(value.(*zip.File).FileInfo().Size()), key)
		return true
	})
}

func Export(srcPath, goFilePath, pkgName string, keyPrefix ...string) error {
	buffer := bytes.NewBuffer(nil)
	err := gcompress.ZipPathWriter(srcPath, buffer, keyPrefix...)
	if err != nil {
		return err
	}
	return gfile.PutContents(
		goFilePath,
		fmt.Sprintf(
			`package %s

import "github.com/gogf/gf/os/gres"

func init() {
	if err := gres.Add(%s); err != nil {
		panic(err)
	}
}
`, pkgName, utilbytes.Export(buffer.Bytes())),
	)
}
