// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
	"io"
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"strings"
)

type GetExtractorFn func(filed interface{}) FileExtractor
type FileExtractor interface {
	GetFile() io.ReadCloser
	GetName() string
	GetExt() string
	Size() int64
}

// UploadFile wraps the multipart uploading file with more and convenient features.
type UploadFile struct {
	FileExtractor
}

type base64Extractor struct {
	Content string
}

func NewBase64Extractor(field interface{}) FileExtractor {
	return base64Extractor{Content: gconv.String(field)}
}

func (b base64Extractor) Size() int64 {
	return int64(len(b.Content))
}

func (b base64Extractor) GetFile() io.ReadCloser {
	i := gstr.PosI(b.Content, "base64,")
	s := gstr.SubStr(b.Content, i+len("base64,"))
	bs, _ := gbase64.DecodeString(s)
	return ioutil.NopCloser(bytes.NewReader(bs))
}

func (b base64Extractor) GetName() string {
	return strings.ToLower(strconv.FormatInt(gtime.TimestampNano(), 36) + grand.S(6))
}

func (b base64Extractor) GetExt() string {
	ext, err := GetFileExt(b.GetFile())
	if err != nil {
		return ""
	}
	return ext
}
func GetFileExt(r io.Reader) (string, error) {
	reader, err := mimetype.DetectReader(r)
	if err != nil {
		return "", nil
	}
	return reader.Extension(), nil
}

func NewMultipartExtractor(field interface{}) FileExtractor {
	file := field.(*multipart.FileHeader)
	return multipartExtractor{file}
}

type multipartExtractor struct {
	*multipart.FileHeader
}

func (m multipartExtractor) Size() int64 {
	return m.FileHeader.Size
}

func (m multipartExtractor) GetFile() io.ReadCloser {
	file, err := m.FileHeader.Open()
	if err != nil {
		return ioutil.NopCloser(bytes.NewReader(nil))
	}
	return file
}

func (m multipartExtractor) GetName() string {
	return gfile.Basename(m.Filename)
}

func (m multipartExtractor) GetExt() string {
	return gfile.Ext(m.Filename)
}

// UploadFiles is array type for *UploadFile.
type UploadFiles []*UploadFile

// Save saves the single uploading file to directory path and returns the saved file name.
//
// The parameter <dirPath> should be a directory path or it returns error.
//
// The parameter <randomlyRename> specifies whether randomly renames the file name, which
// make sense if the <path> is a directory.
//
// Note that it will overwrite the target file if there's already a same name file exist.
func (f *UploadFile) Save(dirPath string, randomlyRename ...bool) (filename string, err error) {
	if f == nil {
		return
	}
	if !gfile.Exists(dirPath) {
		if err = gfile.Mkdir(dirPath); err != nil {
			return
		}
	} else if !gfile.IsDir(dirPath) {
		return "", errors.New(`parameter "dirPath" should be a directory path`)
	}

	file := f.GetFile()
	defer file.Close()

	name := gfile.Basename(f.GetName())
	if len(randomlyRename) > 0 && randomlyRename[0] {
		name = strings.ToLower(strconv.FormatInt(gtime.TimestampNano(), 36) + grand.S(6))
		name = name + f.GetExt()
	}
	filePath := gfile.Join(dirPath, name)
	newFile, err := gfile.Create(filePath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()
	intlog.Printf(`save upload file: %s`, filePath)
	if _, err := io.Copy(newFile, file); err != nil {
		return "", err
	}
	return gfile.Basename(filePath), nil
}

// Save saves all uploading files to specified directory path and returns the saved file names.
//
// The parameter <dirPath> should be a directory path or it returns error.
//
// The parameter <randomlyRename> specifies whether randomly renames all the file names.
func (fs UploadFiles) Save(dirPath string, randomlyRename ...bool) (filenames []string, err error) {
	if len(fs) == 0 {
		return nil, nil
	}
	for _, f := range fs {
		if filename, err := f.Save(dirPath, randomlyRename...); err != nil {
			return filenames, err
		} else {
			filenames = append(filenames, filename)
		}
	}
	return
}

// GetUploadFile retrieves and returns the uploading file with specified form name.
// This function is used for retrieving single uploading file object, which is
// uploaded using multipart form content type.
//
// Note that the <name> is the file field name of the multipart form from client.
func (r *Request) GetUploadFile(name string, fn ...GetExtractorFn) *UploadFile {
	uploadFiles := r.GetUploadFiles(name)
	if len(uploadFiles) > 0 {
		return uploadFiles[0]
	}
	return nil
}

// GetUploadFiles retrieves and returns multiple uploading files with specified form name.
// This function is used for retrieving multiple uploading file objects, which are
// uploaded using multipart form content type.
//
// Note that the <name> is the file field name of the multipart form from client.
func (r *Request) GetUploadFiles(name string, fn ...GetExtractorFn) UploadFiles {
	var fun GetExtractorFn = NewMultipartExtractor
	if len(fn) > 0 {
		fun = fn[0]
	}
	multipartFiles := r.GetMultipartFiles(name)
	if len(multipartFiles) > 0 {
		uploadFiles := make(UploadFiles, len(multipartFiles))
		for k, v := range multipartFiles {
			uploadFiles[k] = &UploadFile{
				fun(v),
			}
		}
		return uploadFiles
	}
	return nil
}
