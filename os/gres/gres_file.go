// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"archive/zip"
	"bytes"
	"github.com/gogf/gf/internal/json"
	"io"
	"os"
)

type File struct {
	file     *zip.File
	reader   *bytes.Reader
	resource *Resource
}

// Name returns the name of the file.
func (f *File) Name() string {
	return f.file.Name
}

// Open returns a ReadCloser that provides access to the File's contents.
// Multiple files may be read concurrently.
func (f *File) Open() (io.ReadCloser, error) {
	return f.file.Open()
}

// Content returns the content of the file.
func (f *File) Content() []byte {
	reader, err := f.Open()
	if err != nil {
		return nil
	}
	defer reader.Close()
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, reader); err != nil {
		return nil
	}
	return buffer.Bytes()
}

// FileInfo returns an os.FileInfo for the FileHeader.
func (f *File) FileInfo() os.FileInfo {
	return f.file.FileInfo()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (f *File) MarshalJSON() ([]byte, error) {
	info := f.FileInfo()
	return json.Marshal(map[string]interface{}{
		"name": f.Name(),
		"size": info.Size(),
		"time": info.ModTime(),
		"file": !info.IsDir(),
	})
}
