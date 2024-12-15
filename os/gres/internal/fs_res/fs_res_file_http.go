// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package fs_res

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gres/internal/defines"
)

// HttpFileImp implements the interface fs.File.
type HttpFileImp struct {
	fs         defines.FS    // FS is the file system that contains this file
	zipFile    *zip.File     // File is the underlying file object
	readSeeker io.ReadSeeker // ReadCloser is the underlying file object
}

var _ http.File = (*HttpFileImp)(nil)

func NewHttpFile(fs defines.FS, zipFile *zip.File) (*HttpFileImp, error) {
	readCloser, err := zipFile.Open()
	if err != nil {
		return nil, gerror.WrapCodef(gcode.CodeOperationFailed, err, `open zip file failed`)
	}
	content, err := io.ReadAll(readCloser)
	if err != nil {
		return nil, gerror.WrapCodef(gcode.CodeOperationFailed, err, `read zip file content failed`)
	}
	return &HttpFileImp{
		readSeeker: bytes.NewReader(content),
		zipFile:    zipFile,
		fs:         fs,
	}, nil
}

// Stat returns the FileInfo structure describing file.
func (f *HttpFileImp) Stat() (os.FileInfo, error) {
	return f.zipFile.FileInfo(), nil
}

// Close implements interface of http.File.
func (f *HttpFileImp) Close() error {
	return nil
}

// Readdir implements Readdir interface of http.File.
func (f *HttpFileImp) Readdir(count int) ([]os.FileInfo, error) {
	files := f.fs.ScanDir(f.zipFile.Name, "*", false)
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
