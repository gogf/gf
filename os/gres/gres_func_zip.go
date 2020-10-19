// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"archive/zip"
	"github.com/gogf/gf/internal/fileinfo"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gfile"
	"io"
	"os"
	"strings"
	"time"
)

// ZipPathWriter compresses <paths> to <writer> using zip compressing algorithm.
// The unnecessary parameter <prefix> indicates the path prefix for zip file.
//
// Note that the parameter <paths> can be either a directory or a file, which
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

// doZipPathWriter compresses the file of given <path> and writes the content to <zipWriter>.
// The parameter <exclude> specifies the exclusive file path that is not compressed to <zipWriter>,
// commonly the destination zip file path.
// The unnecessary parameter <prefix> indicates the path prefix for zip file.
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
	headerPrefix = strings.TrimRight(headerPrefix, "\\/")
	if len(headerPrefix) > 0 && gfile.IsDir(path) {
		headerPrefix += "/"
	}
	if headerPrefix == "" {
		headerPrefix = gfile.Basename(path)
	}
	headerPrefix = strings.Replace(headerPrefix, "//", "/", -1)
	for _, file := range files {
		if exclude == file {
			intlog.Printf(`exclude file path: %s`, file)
			continue
		}
		err := zipFile(file, headerPrefix+gfile.Dir(file[len(path):]), zipWriter)
		if err != nil {
			return err
		}
	}
	// Add all directories to zip archive.
	if headerPrefix != "" {
		var name string
		path = headerPrefix
		for {
			name = gfile.Basename(path)
			err := zipFileVirtual(
				fileinfo.New(name, 0, os.ModeDir|os.ModePerm, time.Now()), path, zipWriter,
			)
			if err != nil {
				return err
			}
			if path == "/" || !strings.Contains(path, "/") {
				break
			}
			path = gfile.Dir(path)
		}
	}
	return nil
}

// zipFile compresses the file of given <path> and writes the content to <zw>.
// The parameter <prefix> indicates the path prefix for zip file.
func zipFile(path string, prefix string, zw *zip.Writer) error {
	prefix = strings.Replace(prefix, "//", "/", -1)
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
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
		return err
	}
	if !info.IsDir() {
		if _, err = io.Copy(writer, file); err != nil {
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
	if _, err := zw.CreateHeader(header); err != nil {
		return err
	}
	return nil
}

func createFileHeader(info os.FileInfo, prefix string) (*zip.FileHeader, error) {
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return nil, err
	}
	if len(prefix) > 0 {
		prefix = strings.Replace(prefix, `\`, `/`, -1)
		prefix = strings.TrimRight(prefix, `/`)
		header.Name = prefix + `/` + header.Name
	}
	return header, nil
}
