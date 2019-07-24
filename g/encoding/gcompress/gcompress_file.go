// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcompress

import (
	"archive/zip"
	"bytes"
	"github.com/gogf/gf/g/os/gfile"
	"io"
	"os"
	"path/filepath"
)

// Zip compresses <path> to <dest> using zip compressing algorithm.
func ZipPath(path, dest string, prefix ...string) error {
	d, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	files, err := gfile.ScanDir(path, "*.*", true)
	if err != nil {
		return err
	}
	pathRealPath := gfile.RealPath(path)
	destRealPath := gfile.RealPath(dest)
	headerPrefix := ""
	if len(prefix) > 0 {
		headerPrefix = prefix[0]
	}
	for _, file := range files {
		if destRealPath == file {
			continue
		}
		err := zipFile(file, headerPrefix+gfile.Dir(file[len(pathRealPath):]), w)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnZipFile decompresses <archive> to <dest> using zip compressing algorithm.
func UnZipFile(archive, dest string) error {
	readerCloser, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer readerCloser.Close()
	return unZipFileWithReader(&readerCloser.Reader, dest)
}

// UnZipContent decompresses <data> to <dest> using zip compressing algorithm.
func UnZipContent(data []byte, dest string) error {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	return unZipFileWithReader(reader, dest)
}

func unZipFileWithReader(reader *zip.Reader, dest string) error {
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)
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
		header.Name = prefix + "/" + header.Name
	} else {
		header.Name = header.Name
	}

	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	if _, err = io.Copy(writer, file); err != nil {
		return err
	}

	return nil
}
