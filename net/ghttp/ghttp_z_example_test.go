// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleServer_Run() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello world")
	})
	s.SetPort(8999)
	s.Run()
}

// Custom saving file name.
func ExampleUploadFile_Save() {
	s := g.Server()
	s.BindHandler("/upload", func(r *ghttp.Request) {
		file := r.GetUploadFile("TestFile")
		if file == nil {
			r.Response.Write("empty file")
			return
		}
		file.Filename = "MyCustomFileName.txt"
		fileName, err := file.Save(gfile.Temp())
		if err != nil {
			r.Response.Write(err)
			return
		}
		r.Response.Write(fileName)
	})
	s.SetPort(8999)
	s.Run()
}

func ExampleRegisterParseRule() {
	ghttp.RegisterParseRule("slug", func(ctx context.Context, in ghttp.ParseFuncInput) (any, error) {
		value, ok := in.Value.(string)
		if !ok {
			return in.Value, nil
		}
		value = strings.TrimSpace(strings.ToLower(value))
		value = strings.ReplaceAll(value, " ", "-")
		return value, nil
	})
	defer ghttp.DeleteParseRule("slug")

	type CreateReq struct {
		Title string `json:"title" parse:"slug"`
	}

	_ = CreateReq{}
}

func ExampleRequest_Parse() {
	type CreateReq struct {
		Title string   `json:"title" parse:"trim-space" v:"required"`
		Tags  []string `json:"tags" parse:"foreach|trim-space|lower"`
	}

	s := g.Server()
	s.BindHandler("/menu", func(ctx context.Context, req *CreateReq) (res any, err error) {
		return req, nil
	})
	s.SetPort(8999)
	s.Run()
}
