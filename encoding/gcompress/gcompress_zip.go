// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcompress

import (
	"archive/zip"
	"bytes"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()
	for _, path := range strings.Split(paths, ",") {
		path = strings.TrimSpace(path)
		if err := doZipPathWriter(path, gfile.RealPath(dest), zipWriter, prefix...); err != nil {
			return err
		}
	}
	return nil
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
	var err error
	var files []string
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
	if gfile.IsDir(path) {
		if len(headerPrefix) > 0 {
			headerPrefix += "/"
		} else {
			headerPrefix = gfile.Basename(path)
		}

	}
	headerPrefix = strings.Replace(headerPrefix, "//", "/", -1)
	for _, file := range files {
		if exclude == file {
			intlog.Printf(`exclude file path: %s`, file)
			continue
		}
		dir := gfile.Dir(file[len(path):])
		if dir == "." {
			dir = ""
		}
		err := zipFile(file, headerPrefix+dir, zipWriter)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnZipFile decompresses <archive> to <dest> using zip compressing algorithm.
// The optional parameter <path> specifies the unzipped path of <archive>,
// which can be used to specify part of the archive file to unzip.
//
// Note that the parameter <dest> should be a directory.
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
// Note that the parameter <dest> should be a directory.
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

// zipFile compresses the file of given <path> and writes the content to <zw>.
// The parameter <prefix> indicates the path prefix for zip file.
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

	if info.IsDir() {
		header.Name += "/"
	} else {
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
