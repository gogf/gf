// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/fileinfo"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
)

// ZipPathWriter compresses `paths` to `writer` using zip compressing algorithm.
// The unnecessary parameter `prefix` indicates the path prefix for zip file.
//
// Note that the parameter `paths` can be either a directory or a file, which
// supports multiple paths join with ','.
func zipPathWriter(paths string, writer io.Writer, prefix ...string) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()
	for _, path := range strings.Split(paths, ",") {
		path = strings.TrimSpace(path)
		if err := doZipPathWriter(path, "", zipWriter, prefix...); err != nil {
			return err
		}
	}
	return nil
}

// doZipPathWriter compresses the file of given `path` and writes the content to `zipWriter`.
// The parameter `exclude` specifies the exclusive file path that is not compressed to `zipWriter`,
// commonly the destination zip file path.
// The unnecessary parameter `prefix` indicates the path prefix for zip file.
func doZipPathWriter(path string, exclude string, zipWriter *zip.Writer, prefix ...string) error {
	var (
		err   error
		files []string
	)
	path, err = gfile.Search(path)
	if err != nil {
		return err
	}
	if gfile.IsDir(path) {
		files, err = gfile.ScanDir(path, "*", true)
		if err != nil {
			return err
		}
	} else {
		files = []string{path}
	}
	headerPrefix := ""
	if len(prefix) > 0 && prefix[0] != "" {
		headerPrefix = prefix[0]
	}
	headerPrefix = strings.TrimRight(headerPrefix, `\/`)
	if len(headerPrefix) > 0 && gfile.IsDir(path) {
		headerPrefix += "/"
	}
	if headerPrefix == "" {
		headerPrefix = gfile.Basename(path)
	}
	headerPrefix = strings.Replace(headerPrefix, `//`, `/`, -1)
	for _, file := range files {
		if exclude == file {
			intlog.Printf(context.TODO(), `exclude file path: %s`, file)
			continue
		}
		if err = zipFile(file, headerPrefix+gfile.Dir(file[len(path):]), zipWriter); err != nil {
			return err
		}
	}
	// Add all directories to zip archive.
	if headerPrefix != "" {
		var name string
		path = headerPrefix
		for {
			name = strings.Replace(gfile.Basename(path), `\`, `/`, -1)
			if err = zipFileVirtual(fileinfo.New(name, 0, os.ModeDir|os.ModePerm, time.Now()), path, zipWriter); err != nil {
				return err
			}
			if path == `/` || !strings.Contains(path, `/`) {
				break
			}
			path = gfile.Dir(path)
		}
	}
	return nil
}

// zipFile compresses the file of given `path` and writes the content to `zw`.
// The parameter `prefix` indicates the path prefix for zip file.
func zipFile(path string, prefix string, zw *zip.Writer) error {
	prefix = strings.Replace(prefix, `//`, `/`, -1)
	file, err := os.Open(path)
	if err != nil {
		err = gerror.Wrapf(err, `os.Open failed for file "%s"`, path)
		return nil
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		err = gerror.Wrapf(err, `read file stat failed for file "%s"`, path)
		return err
	}
	header, err := createFileHeader(info, prefix)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		header.Method = zip.Deflate
	}
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

func zipFileVirtual(info os.FileInfo, path string, zw *zip.Writer) error {
	header, err := createFileHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = path
	if _, err = zw.CreateHeader(header); err != nil {
		err = gerror.Wrapf(err, `create zip header failed for %#v`, header)
		return err
	}
	return nil
}

func createFileHeader(info os.FileInfo, prefix string) (*zip.FileHeader, error) {
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		err = gerror.Wrapf(err, `create file header failed for name "%s"`, info.Name())
		return nil, err
	}
	if len(prefix) > 0 {
		header.Name = prefix + `/` + header.Name
		header.Name = strings.Replace(header.Name, `\`, `/`, -1)
		header.Name, _ = gregex.ReplaceString(`/{2,}`, `/`, header.Name)
	}
	return header, nil
}
