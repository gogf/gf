// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package fs_res

import (
	"archive/zip"
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres/internal/defines"
	"github.com/gogf/gf/v2/text/gstr"
)

// FileImp implements the interface fs.File.
type FileImp struct {
	file *zip.File  // File is the underlying file object
	fs   defines.FS // FS is the file system that contains this file
}

var _ defines.File = (*FileImp)(nil)

func (f *FileImp) Name() string {
	return f.file.Name
}

// FileInfo returns an os.FileInfo describing this file
func (f *FileImp) FileInfo() os.FileInfo {
	return f.file.FileInfo()
}

// Stat returns the FileInfo structure describing file.
func (f *FileImp) Stat() (os.FileInfo, error) {
	return f.FileInfo(), nil
}

func (f *FileImp) Open() (io.ReadCloser, error) {
	return f.file.Open()
}

func (f *FileImp) HttpFile() (http.File, error) {
	return NewHttpFile(f.fs, f.file)
}

// Content returns the file content
func (f *FileImp) Content() []byte {
	readCloser, err := f.file.Open()
	if err != nil {
		intlog.Error(context.Background(), err)
		return nil
	}
	defer readCloser.Close()
	content, err := io.ReadAll(readCloser)
	if err != nil {
		intlog.Error(context.Background(), err)
		return nil
	}
	return content
}

// Export exports and saves all its sub files to specified system path `dst` recursively.
func (f *FileImp) Export(dst string, option ...defines.ExportOption) error {
	var (
		err          error
		name         string
		path         string
		exportOption defines.ExportOption
		exportFiles  []defines.File
	)
	if f.FileInfo().IsDir() {
		exportFiles = f.fs.ScanDir(f.Name(), "*", true)
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

type jsonFileInfo struct {
	Name  string
	Size  int64
	Time  time.Time
	IsDir bool
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (f *FileImp) MarshalJSON() ([]byte, error) {
	info := f.FileInfo()
	return gjson.Marshal(jsonFileInfo{
		Name:  f.Name(),
		Size:  info.Size(),
		Time:  info.ModTime(),
		IsDir: info.IsDir(),
	})
}
