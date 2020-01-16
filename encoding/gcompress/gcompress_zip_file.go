// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcompress

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/internal/fileinfo"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
)

// ZipPath compresses <paths> to <dest> using zip compressing algorithm.
// The unnecessary parameter <prefix> indicates the path prefix for zip file.
//
// Note that the parameter <paths> can be either a directory or a file, which
// supports multiple paths join with ','.
func ZipPath(paths, dest string, prefix ...string) error {
	writer, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer writer.Close()
	return ZipPathWriter(paths, writer, prefix...)
}

// ZipPathWriter compresses <paths> to <writer> using zip compressing algorithm.
// The unnecessary parameter <prefix> indicates the path prefix for zip file.
//
// Note that the parameter <paths> can be either a directory or a file, which
// supports multiple paths join with ','.
func ZipPathWriter(paths string, writer io.Writer, prefix ...string) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()
	for _, path := range strings.Split(paths, ",") {
		path = strings.TrimSpace(path)
		if err := doZipPathWriter(path, zipWriter, prefix...); err != nil {
			return err
		}
	}
	return nil
}

func doZipPathWriter(path string, zipWriter *zip.Writer, prefix ...string) error {
	var err error
	var files []string
	realPath, err := gfile.Search(path)
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
		err := zipFile(file, headerPrefix+gfile.Dir(file[len(realPath):]), zipWriter)
		if err != nil {
			return err
		}
	}
	// Add prefix to zip archive.
	path = headerPrefix
	for {
		err := zipFileVirtual(fileinfo.New(gfile.Basename(path), 0, os.ModeDir, time.Now()), path, zipWriter)
		if err != nil {
			return err
		}
		if path == "/" || !strings.Contains(path, "/") {
			break
		}
		path = gfile.Dir(path)
	}
	return nil
}

// UnZipFile decompresses <archive> to <dest> using zip compressing algorithm.
// The parameter <path> specifies the unzipped path of <archive>,
// which can be used to specify part of the archive file to unzip.
//
// Note thate the parameter <dest> should be a directory.
func UnZipFile(archive, dest string, path ...string) error {
	readerCloser, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer readerCloser.Close()
	return unZipFileWithReader(&readerCloser.Reader, dest, path...)
}

// UnZipContent decompresses <data> to <dest> using zip compressing algorithm.
// The parameter <path> specifies the unzipped path of <archive>,
// which can be used to specify part of the archive file to unzip.
//
// Note thate the parameter <dest> should be a directory.
func UnZipContent(data []byte, dest string, path ...string) error {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	return unZipFileWithReader(reader, dest, path...)
}

func unZipFileWithReader(reader *zip.Reader, dest string, path ...string) error {
	prefix := ""
	if len(path) > 0 {
		prefix = gstr.Replace(path[0], `\`, `/`)
	}
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	name := ""
	for _, file := range reader.File {
		name = gstr.Replace(file.Name, `\`, `/`)
		name = gstr.Trim(name, "/")
		if prefix != "" {
			if name[0:len(prefix)] != prefix {
				continue
			}
			name = name[len(prefix):]
		}
		path := filepath.Join(dest, name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}
		dir := filepath.Dir(path)
		if len(dir) > 0 {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					return err
				}
			}
		}
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}

func zipFile(path string, prefix string, zw *zip.Writer) error {
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
