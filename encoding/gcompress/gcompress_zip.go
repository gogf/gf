// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcompress

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

// ZipPath compresses `fileOrFolderPaths` to `dstFilePath` using zip compressing algorithm.
//
// The parameter `paths` can be either a directory or a file, which
// supports multiple paths join with ','.
// The unnecessary parameter `prefix` indicates the path prefix for zip file.
func ZipPath(fileOrFolderPaths, dstFilePath string, prefix ...string) error {
	writer, err := os.Create(dstFilePath)
	if err != nil {
		err = gerror.Wrapf(err, `os.Create failed for name "%s"`, dstFilePath)
		return err
	}
	defer writer.Close()
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()
	for _, path := range strings.Split(fileOrFolderPaths, ",") {
		path = strings.TrimSpace(path)
		if err = doZipPathWriter(path, gfile.RealPath(dstFilePath), zipWriter, prefix...); err != nil {
			return err
		}
	}
	return nil
}

// ZipPathWriter compresses `fileOrFolderPaths` to `writer` using zip compressing algorithm.
//
// Note that the parameter `fileOrFolderPaths` can be either a directory or a file, which
// supports multiple paths join with ','.
// The unnecessary parameter `prefix` indicates the path prefix for zip file.
func ZipPathWriter(fileOrFolderPaths string, writer io.Writer, prefix ...string) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()
	for _, path := range strings.Split(fileOrFolderPaths, ",") {
		path = strings.TrimSpace(path)
		if err := doZipPathWriter(path, "", zipWriter, prefix...); err != nil {
			return err
		}
	}
	return nil
}

// ZipPathContent compresses `fileOrFolderPaths` to []byte using zip compressing algorithm.
//
// Note that the parameter `fileOrFolderPaths` can be either a directory or a file, which
// supports multiple paths join with ','.
// The unnecessary parameter `prefix` indicates the path prefix for zip file.
func ZipPathContent(fileOrFolderPaths string, prefix ...string) ([]byte, error) {
	var (
		err    error
		buffer = bytes.NewBuffer(nil)
	)
	if err = ZipPathWriter(fileOrFolderPaths, buffer, prefix...); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// doZipPathWriter compresses given `fileOrFolderPaths` and writes the content to `zipWriter`.
//
// The parameter `fileOrFolderPath` can be either a single file or folder path.
// The parameter `exclude` specifies the exclusive file path that is not compressed to `zipWriter`,
// commonly the destination zip file path.
// The unnecessary parameter `prefix` indicates the path prefix for zip file.
func doZipPathWriter(fileOrFolderPath string, exclude string, zipWriter *zip.Writer, prefix ...string) error {
	var (
		err   error
		files []string
	)
	fileOrFolderPath, err = gfile.Search(fileOrFolderPath)
	if err != nil {
		return err
	}
	if gfile.IsDir(fileOrFolderPath) {
		files, err = gfile.ScanDir(fileOrFolderPath, "*", true)
		if err != nil {
			return err
		}
	} else {
		files = []string{fileOrFolderPath}
	}
	headerPrefix := ""
	if len(prefix) > 0 && prefix[0] != "" {
		headerPrefix = prefix[0]
	}
	headerPrefix = strings.TrimRight(headerPrefix, "\\/")
	if gfile.IsDir(fileOrFolderPath) {
		if len(headerPrefix) > 0 {
			headerPrefix += "/"
		} else {
			headerPrefix = gfile.Basename(fileOrFolderPath)
		}
	}
	headerPrefix = strings.ReplaceAll(headerPrefix, "//", "/")
	for _, file := range files {
		if exclude == file {
			intlog.Printf(context.TODO(), `exclude file path: %s`, file)
			continue
		}
		dir := gfile.Dir(file[len(fileOrFolderPath):])
		if dir == "." {
			dir = ""
		}
		if err = zipFile(file, headerPrefix+dir, zipWriter); err != nil {
			return err
		}
	}
	return nil
}

// UnZipFile decompresses `archive` to `dstFolderPath` using zip compressing algorithm.
//
// The parameter `dstFolderPath` should be a directory.
// The optional parameter `zippedPrefix` specifies the unzipped path of `zippedFilePath`,
// which can be used to specify part of the archive file to unzip.
func UnZipFile(zippedFilePath, dstFolderPath string, zippedPrefix ...string) error {
	readerCloser, err := zip.OpenReader(zippedFilePath)
	if err != nil {
		err = gerror.Wrapf(err, `zip.OpenReader failed for name "%s"`, dstFolderPath)
		return err
	}
	defer readerCloser.Close()
	return unZipFileWithReader(&readerCloser.Reader, dstFolderPath, zippedPrefix...)
}

// UnZipContent decompresses `zippedContent` to `dstFolderPath` using zip compressing algorithm.
//
// The parameter `dstFolderPath` should be a directory.
// The parameter `zippedPrefix` specifies the unzipped path of `zippedContent`,
// which can be used to specify part of the archive file to unzip.
func UnZipContent(zippedContent []byte, dstFolderPath string, zippedPrefix ...string) error {
	reader, err := zip.NewReader(bytes.NewReader(zippedContent), int64(len(zippedContent)))
	if err != nil {
		err = gerror.Wrapf(err, `zip.NewReader failed`)
		return err
	}
	return unZipFileWithReader(reader, dstFolderPath, zippedPrefix...)
}

func unZipFileWithReader(reader *zip.Reader, dstFolderPath string, zippedPrefix ...string) error {
	prefix := ""
	if len(zippedPrefix) > 0 {
		prefix = gstr.Replace(zippedPrefix[0], `\`, `/`)
	}
	if err := os.MkdirAll(dstFolderPath, 0755); err != nil {
		return err
	}
	var (
		name    string
		dstPath string
		dstDir  string
	)
	for _, file := range reader.File {
		name = gstr.Replace(file.Name, `\`, `/`)
		name = gstr.Trim(name, "/")
		if prefix != "" {
			if name[0:len(prefix)] != prefix {
				continue
			}
			name = name[len(prefix):]
		}
		dstPath = filepath.Join(dstFolderPath, name)
		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(dstPath, file.Mode())
			continue
		}
		dstDir = filepath.Dir(dstPath)
		if len(dstDir) > 0 {
			if _, err := os.Stat(dstDir); os.IsNotExist(err) {
				if err = os.MkdirAll(dstDir, 0755); err != nil {
					err = gerror.Wrapf(err, `os.MkdirAll failed for path "%s"`, dstDir)
					return err
				}
			}
		}
		fileReader, err := file.Open()
		if err != nil {
			err = gerror.Wrapf(err, `file.Open failed`)
			return err
		}
		// The fileReader is closed in function doCopyForUnZipFileWithReader.
		if err = doCopyForUnZipFileWithReader(file, fileReader, dstPath); err != nil {
			return err
		}
	}
	return nil
}

func doCopyForUnZipFileWithReader(file *zip.File, fileReader io.ReadCloser, dstPath string) error {
	defer fileReader.Close()
	targetFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		err = gerror.Wrapf(err, `os.OpenFile failed for name "%s"`, dstPath)
		return err
	}
	defer targetFile.Close()

	if _, err = io.Copy(targetFile, fileReader); err != nil {
		err = gerror.Wrapf(err, `io.Copy failed from "%s" to "%s"`, file.Name, dstPath)
		return err
	}
	return nil
}

// zipFile compresses the file of given `filePath` and writes the content to `zw`.
// The parameter `prefix` indicates the path prefix for zip file.
func zipFile(filePath string, prefix string, zw *zip.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		err = gerror.Wrapf(err, `os.Open failed for name "%s"`, filePath)
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		err = gerror.Wrapf(err, `file.Stat failed for name "%s"`, filePath)
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
		err = gerror.Wrapf(err, `zip.Writer.CreateHeader failed for header "%#v"`, header)
		return err
	}
	if !info.IsDir() {
		if _, err = io.Copy(writer, file); err != nil {
			err = gerror.Wrapf(err, `io.Copy failed from "%s" to "%s"`, filePath, header.Name)
			return err
		}
	}
	return nil
}

func createFileHeader(info os.FileInfo, prefix string) (*zip.FileHeader, error) {
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		err = gerror.Wrapf(err, `zip.FileInfoHeader failed for info "%#v"`, info)
		return nil, err
	}

	if len(prefix) > 0 {
		prefix = strings.ReplaceAll(prefix, `\`, `/`)
		prefix = strings.TrimRight(prefix, `/`)
		header.Name = prefix + `/` + header.Name
	}
	return header, nil
}
