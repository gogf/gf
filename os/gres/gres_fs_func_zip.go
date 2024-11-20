// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/fileinfo"
	"github.com/gogf/gf/v2/os/gfile"
)

func zipFsWriter(dirfs fs.FS, dirPath string, writer io.Writer, option ...Option) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()
	if err := doZipFsWriter(dirfs, dirPath, zipWriter, option...); err != nil {
		return err
	}
	return nil
}

func doZipFsWriter(dirfs fs.FS, dirPath string, zipWriter *zip.Writer, option ...Option) error {
	var (
		err         error
		files       []string
		usedOption  Option
		listfsfiles func(dirfs fs.FS) ([]string, error)
	)
	listfsfiles = func(dirfs fs.FS) ([]string, error) {

		files, err := fs.Glob(dirfs, "*")
		if err != nil {
			return nil, err
		}
		result := make([]string, 0, len(files))
		for _, f := range files {
			result = append(result, f)
			finfo, err := fs.Stat(dirfs, f)
			if err != nil {
				return nil, err
			}
			if finfo.IsDir() {
				subfs, _ := fs.Sub(dirfs, f)
				if subresult, err := listfsfiles(subfs); err != nil {
					return nil, err
				} else {
					for _, rf := range subresult {
						result = append(result, fmt.Sprintf("%s/%s", f, rf))
					}
				}
			}
		}
		return result, nil
	}
	if len(option) > 0 {
		usedOption = option[0]
	}
	if files, err = listfsfiles(dirfs); err != nil {
		return err
	}

	headerPrefix := usedOption.Prefix
	if !(headerPrefix == "/") {
		headerPrefix = strings.TrimRight(headerPrefix, `\/`)
	}

	if headerPrefix != "" {
		headerPrefix += "/"
	}

	if headerPrefix == "" {
		if usedOption.KeepPath {
			// It keeps the path from file system to zip info in resource manager.
			// Usually for relative path, it makes little sense for absolute path.
			headerPrefix = dirPath
		} else {
			headerPrefix = filepath.Base(dirPath)
		}
	}

	headerPrefix = strings.ReplaceAll(headerPrefix, `//`, `/`)
	for _, file := range files {
		// It here calculates the file name prefix, especially packing the directory.
		// Eg:
		// path: dir1
		// file: dir1/dir2/file
		// file[len(absolutePath):] => /dir2/file
		// gfile.Dir(subFilePath)   => /dir2
		var subFilePath string = filepath.Clean(filepath.Dir(file))
		if subFilePath == "." {
			subFilePath = ""
		}
		if err = zipFsFile(dirfs, file, headerPrefix+"/"+subFilePath, zipWriter); err != nil {
			return err
		}
	}
	// Add all directories to zip archive.
	if headerPrefix != "" {
		var (
			name    string
			tmpPath = headerPrefix
		)
		for {
			name = strings.ReplaceAll(gfile.Basename(tmpPath), "\\", `/`)
			err = zipFileVirtual(fileinfo.New(name, 0, os.ModeDir|os.ModePerm, time.Now()), tmpPath, zipWriter)
			if err != nil {
				return err
			}
			if tmpPath == `/` || !strings.Contains(tmpPath, `/`) {
				break
			}
			tmpPath = gfile.Dir(tmpPath)
		}
	}
	return nil
}

func zipFsFile(fspath fs.FS, path string, prefix string, zw *zip.Writer) error {
	prefix = strings.ReplaceAll(prefix, `//`, `/`)
	file, err := fspath.Open(path)
	if err != nil {
		return gerror.Wrapf(err, `fs.ReadFile failed for path "%s"`, path)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		err = gerror.Wrapf(err, `read file stat failed for path "%s"`, path)
		return err
	}

	header, err := createFileHeader(info, prefix)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		// Default compression level.
		header.Method = zip.Deflate
	}
	// Zip header containing the info of a zip file.
	writer, err := zw.CreateHeader(header)
	if err != nil {
		err = gerror.Wrapf(err, `create zip header failed for %#v`, header)
		return err
	}
	if !info.IsDir() {
		if _, err = io.Copy(writer, file); err != nil {
			err = gerror.Wrapf(err, `io.Copy failed for file "%s"`, path)
			return err
		}
	}
	return nil
}
