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
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

// File implements the interface fs.File.
type File struct {
	name     string      // Name is the file name
	path     string      // Path is the file path
	file     os.FileInfo // FileInfo is the underlying file info
	reader   io.Reader   // Reader is the file content reader
	resource []byte      // Resource is the file content in binary format
	fs       FS          // FS is the file system that contains this file
}

// Name returns the name of the file
func (f *File) Name() string {
	return f.name
}

// Path returns the path of the file
func (f *File) Path() string {
	return f.path
}

// Open opens the file for reading
func (f *File) Open() error {
	if f.reader == nil && len(f.resource) > 0 {
		f.reader = bytes.NewReader(f.resource)
	}
	return nil
}

// Close closes the file
func (f *File) Close() error {
	if closer, ok := f.reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Read reads up to len(p) bytes into p
func (f *File) Read(p []byte) (n int, err error) {
	if f.reader == nil {
		if err := f.Open(); err != nil {
			return 0, err
		}
	}
	return f.reader.Read(p)
}

// Seek implements the io.Seeker interface
func (f *File) Seek(offset int64, whence int) (int64, error) {
	if seeker, ok := f.reader.(io.Seeker); ok {
		return seeker.Seek(offset, whence)
	}
	return 0, fs.ErrInvalid
}

// FileInfo returns an os.FileInfo describing this file
func (f *File) FileInfo() os.FileInfo {
	return f.file
}

// Stat returns the FileInfo structure describing file
func (f *File) Stat() (os.FileInfo, error) {
	return f.file, nil
}

// Content returns the file content
func (f *File) Content() []byte {
	if len(f.resource) > 0 {
		return f.resource
	}
	buffer := new(bytes.Buffer)
	if err := f.Open(); err != nil {
		return nil
	}
	defer f.Close()
	if _, err := io.Copy(buffer, f); err != nil {
		return nil
	}
	f.resource = buffer.Bytes()
	return f.resource
}

// Export exports and saves all its sub files to specified system path `dst` recursively.
func (f *File) Export(dst string, option ...ExportOption) error {
	var (
		err          error
		name         string
		path         string
		exportOption ExportOption
		exportFiles  []*File
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
		if exportFile.FileInfo().IsDir() {
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

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (f *File) MarshalJSON() ([]byte, error) {
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
	file *File
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
