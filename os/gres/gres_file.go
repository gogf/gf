// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"sync"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

// A File provides access to a single file.
// The File interface is the minimum implementation required of the file.
// Directory files should also implement [ReadDirFile].
// A file may implement [io.ReaderAt] or [io.Seeker] as optimizations.
type File interface {
	Name() string
	Path() string
	Content() []byte
	FileInfo() os.FileInfo
	Export(dst string, option ...ExportOption) error

	// For http.File implementation.

	Readdir(count int) ([]os.FileInfo, error)
	io.ReadSeekCloser
}

// File implements the interface fs.File.
type localFile struct {
	name    string      // Name is the file name
	path    string      // Path is the file path
	content []byte      // file content
	file    os.FileInfo // FileInfo is the underlying file info
	fs      FS          // FS is the file system that contains this file
	mu      sync.Mutex  // mu protects concurrent access to the file
}

// Name returns the name of the file
func (f *localFile) Name() string {
	return f.name
}

// Path returns the path of the file
func (f *localFile) Path() string {
	return f.path
}

// FileInfo returns an os.FileInfo describing this file
func (f *localFile) FileInfo() os.FileInfo {
	return f.file
}

// Content returns the file content
func (f *localFile) Content() []byte {
	return f.content
}

// Export exports and saves all its sub files to specified system path `dst` recursively.
func (f *localFile) Export(dst string, option ...ExportOption) error {
	var (
		err          error
		name         string
		path         string
		exportOption ExportOption
		exportFiles  []File
	)
	if f.FileInfo().IsDir() {
		exportFiles = f.fs.ScanDir(f.path, "*", true)
	} else {
		exportFiles = append(exportFiles, f)
	}

	if len(option) > 0 {
		exportOption = option[0]
	}
	for _, exportFile := range exportFiles {
		name = exportFile.Name()
		if exportOption.RemovePrefix != "" {
			name = gstr.TrimLeftStr(name, exportOption.RemovePrefix)
		}
		name = gstr.Trim(name, `\/`)
		if name == "" {
			continue
		}
		path = gfile.Join(dst, name)
		if f.FileInfo().IsDir() {
			err = gfile.Mkdir(path)
		} else {
			err = gfile.PutBytes(path, exportFile.Content())
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Close implements interface of http.File.
func (f *localFile) Close() error {
	return nil
}

// Readdir implements Readdir interface of http.File.
func (f *localFile) Readdir(count int) ([]os.FileInfo, error) {
	files := f.fs.ScanDir(f.Name(), "*", false)
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
func (f *localFile) Read(b []byte) (n int, err error) {
	reader := bytes.NewReader(f.Content())
	if n, err = reader.Read(b); err != nil {
		err = gerror.Wrapf(err, `read content failed`)
	}
	return
}

// Seek implements the io.Seeker interface.
func (f *localFile) Seek(offset int64, whence int) (n int64, err error) {
	reader := bytes.NewReader(f.Content())
	if n, err = reader.Seek(offset, whence); err != nil {
		err = gerror.Wrapf(err, `seek failed for offset %d, whence %d`, offset, whence)
	}
	return
}

func (f *localFile) getReader() (io.ReadSeeker, error) {
	return bytes.NewReader(f.Content()), nil
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (f *localFile) MarshalJSON() ([]byte, error) {
	info := f.FileInfo()
	return gjson.Marshal(map[string]interface{}{
		"name":    f.name,
		"path":    f.path,
		"size":    info.Size(),
		"time":    info.ModTime(),
		"isDir":   info.IsDir(),
		"content": f.Content(),
	})
}

// fileInfo is the internal implementation of os.FileInfo interface.
type fileInfo struct {
	file *localFile
}

// Name returns the base name of the file.
func (fi *fileInfo) Name() string {
	return fi.file.Name()
}

// Size returns the size in bytes of the file.
func (fi *fileInfo) Size() int64 {
	return int64(len(fi.file.Content()))
}

// Mode returns the file mode bits.
func (fi *fileInfo) Mode() fs.FileMode {
	if fi.IsDir() {
		return os.ModeDir | 0755
	}
	return 0644
}

// ModTime returns the modification time.
func (fi *fileInfo) ModTime() time.Time {
	if fi.file.file != nil {
		return fi.file.file.ModTime()
	}
	return time.Now()
}

// IsDir reports whether the file is a directory.
func (fi *fileInfo) IsDir() bool {
	if fi.file.file != nil {
		return fi.file.file.IsDir()
	}
	return false
}

// Sys returns the underlying data source.
func (fi *fileInfo) Sys() interface{} {
	return nil
}
