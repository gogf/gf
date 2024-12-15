// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package fs_std

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gres/internal/defines"
)

// HttpFileImp implements the interface fs.File.
type HttpFileImp struct {
	fs         defines.FS    // FS is the file system that contains this file
	fsFile     fs.File       // File is the underlying file object
	readSeeker io.ReadSeeker // ReadCloser is the underlying file object
}

var _ http.File = (*HttpFileImp)(nil)

func NewHttpFile(fs defines.FS, fsFile fs.File) (*HttpFileImp, error) {
	content, err := io.ReadAll(fsFile)
	if err != nil {
		return nil, gerror.WrapCodef(gcode.CodeOperationFailed, err, `read zip file content failed`)
	}
	return &HttpFileImp{
		readSeeker: bytes.NewReader(content),
		fsFile:     fsFile,
		fs:         fs,
	}, nil
}

// Stat returns the FileInfo structure describing file.
func (f *HttpFileImp) Stat() (os.FileInfo, error) {
	return f.fsFile.Stat()
}

// Close implements interface of http.File.
func (f *HttpFileImp) Close() error {
	return nil
}

// Readdir implements Readdir interface of http.File.
func (f *HttpFileImp) Readdir(count int) ([]os.FileInfo, error) {
	info, err := f.fsFile.Stat()
	if err != nil {
		return nil, gerror.WrapCodef(gcode.CodeOperationFailed, err, `get file info failed`)
	}
	files := f.fs.ScanDir(info.Name(), "*", false)
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

// Read implements the io.Reader interface.
func (f *HttpFileImp) Read(b []byte) (n int, err error) {
	if n, err = f.readSeeker.Read(b); err != nil {
		err = gerror.WrapCodef(gcode.CodeOperationFailed, err, `read content failed`)
	}
	return
}

// Seek implements the io.Seeker interface.
func (f *HttpFileImp) Seek(offset int64, whence int) (n int64, err error) {
	if n, err = f.readSeeker.Seek(offset, whence); err != nil {
		err = gerror.Wrapf(err, `seek failed for offset %d, whence %d`, offset, whence)
	}
	return
}
