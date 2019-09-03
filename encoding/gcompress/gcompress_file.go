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

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
)

// ZipPath compresses <path> to <dest> using zip compressing algorithm.
// The unnecessary parameter <prefix> indicates the path prefix for zip file.
func ZipPath(path, dest string, prefix ...string) error {
	writer, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer writer.Close()
	return ZipPathWriter(path, writer, prefix...)
}

// ZipPathWriter compresses <path> to <writer> using zip compressing algorithm.
// The unnecessary parameter <prefix> indicates the path prefix for zip file.
func ZipPathWriter(path string, writer io.Writer, prefix ...string) error {
	realPath, err := gfile.Search(path)
	if err != nil {
		return err
	}
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()
	files, err := gfile.ScanDir(path, "*", true)
	if err != nil {
		return err
	}
	headerPrefix := ""
	if len(prefix) > 0 {
		headerPrefix = prefix[0]
	}
	headerPrefix = strings.Trim(headerPrefix, "\\/")
	// If path is a directory, add it to the zip prefix.
	if gfile.IsDir(realPath) {
		headerPrefix = headerPrefix + "/" + gfile.Basename(realPath)
	}
	for _, file := range files {
		err := zipFile(file, headerPrefix+gfile.Dir(file[len(realPath):]), zipWriter)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnZipFile decompresses <archive> to <dest> using zip compressing algorithm.
// The parameter <path> specifies the unzipped path of <archive>,
// which can be used to specify part of the archive file to unzip.
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
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	if len(prefix) > 0 {
		prefix = strings.Replace(prefix, `\`, `/`, -1)
		prefix = strings.TrimRight(prefix, `/`)
		header.Name = prefix + `/` + header.Name
	} else {
		header.Name = header.Name
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
