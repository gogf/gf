// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"bytes"
	"os"

	"github.com/gogf/gf/v2/errors/gerror"
)

// Close implements interface of http.File.
func (f *File) Close() error {
	return nil
}

// Readdir implements Readdir interface of http.File.
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	files := f.resource.ScanDir(f.Name(), "*", false)
	if len(files) > 0 {
		if count <= 0 || count > len(files) {
			count = len(files)
		}
		infos := make([]os.FileInfo, count)
		for k, v := range files {
			infos[k] = v.FileInfo()
		}
		return infos, nil
	}
	return nil, nil
}

// Stat implements Stat interface of http.File.
func (f *File) Stat() (os.FileInfo, error) {
	return f.FileInfo(), nil
}

// Read implements the io.Reader interface.
func (f *File) Read(b []byte) (n int, err error) {
	reader, err := f.getReader()
	if err != nil {
		return 0, err
	}
	if n, err = reader.Read(b); err != nil {
		err = gerror.Wrapf(err, `read content failed`)
	}
	return
}

// Seek implements the io.Seeker interface.
func (f *File) Seek(offset int64, whence int) (n int64, err error) {
	reader, err := f.getReader()
	if err != nil {
		return 0, err
	}
	if n, err = reader.Seek(offset, whence); err != nil {
		err = gerror.Wrapf(err, `seek failed for offset %d, whence %d`, offset, whence)
	}
	return
}

func (f *File) getReader() (*bytes.Reader, error) {
	if f.reader == nil {
		f.reader = bytes.NewReader(f.Content())
	}
	return f.reader, nil
}
