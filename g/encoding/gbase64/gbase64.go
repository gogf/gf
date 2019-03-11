// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gbase64 provides useful API for BASE64 encoding/decoding algorithms.
package gbase64

import (
	"encoding/base64"
)

// base64 encode
func Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// base64 decode
func Decode(str string) (string, error) {
	s, e := base64.StdEncoding.DecodeString(str)
	return string(s), e
}

//对文件进行base64位编码,可以通过配置文件节点"base64FileExt" 来控制，如:"jpg,jpeg,gif,png,ico,pdf"
func EncodeFile(filename string) (string, error) {

	var (
		files []byte //
		err   error
		types string //默认允许的文件类型的扩展名，多个用逗号分隔，可以通过配置文件来获取
	)

	files, err = ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	ext := gstr.Trim(gfile.Ext(filename), ".")
	if ext == "" {
		return "", errors.New("The file has no extension")
	}

	ext = gstr.ToLower(ext)
	types = "" //默认文件类型为图片
	allowstr := g.Config().GetString("base64FileExt")
	if allowstr == "" {
		allowstr = "jpg,jpeg,gif,png,ico,pdf" //默认允许的类型
	}

	if gstr.Contains(allowstr, ext) {
		types = "image"
	}

	if ext == "pdf" && gstr.Contains(allowstr, "pdf") {
		types = "application"
	}

	//如果不在允许的上传类型之中，则返回
	if types == "" {
		return "", errors.New("Only 'images' or 'pdf' files are allowed to be encoded")
	}

	imageBase64 := gbase64.Encode(string(files))
	rr := "data:" + types + "/" + ext + ";base64," + imageBase64
	return rr, nil

}
