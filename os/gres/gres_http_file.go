// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"os"
)

type File struct {
	file   *zip.File
	reader *bytes.Reader
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
func (f *File) Content() ([]byte, error) {
	reader, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, reader); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// FileInfo returns an os.FileInfo for the FileHeader.
func (f *File) FileInfo() os.FileInfo {
	return f.file.FileInfo()
}

// Read implements the io.Reader interface.
func (f *File) Read(b []byte) (n int, err error) {
	reader, err := f.getReader()
	if err != nil {
		return 0, err
	}
	return reader.Read(b)
}

// Seek implements the io.Seeker interface.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	reader, err := f.getReader()
	if err != nil {
		return 0, err
	}
	return reader.Seek(offset, whence)
}

func (f *File) getReader() (*bytes.Reader, error) {
	if f.reader == nil {
		content, err := f.Content()
		if err != nil {
			return nil, err
		}
		f.reader = bytes.NewReader(content)
	}
	return f.reader, nil
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (f *File) MarshalJSON() ([]byte, error) {
	info := f.FileInfo()
	return json.Marshal(map[string]interface{}{
		"name": f.Name(),
		"size": info.Size(),
		"time": info.ModTime(),
	})
}
