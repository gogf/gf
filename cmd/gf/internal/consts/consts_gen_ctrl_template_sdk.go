// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package consts

const TemplateGenCtrlSdkPkgNew = `
// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package {PkgName}

import (
	"fmt"

	"github.com/gogf/gf/contrib/sdk/httpclient/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
)

type implementer struct {
	config httpclient.Config
}

func New(config httpclient.Config) IClient {
	if !gstr.HasPrefix(config.URL, "http") {
		config.URL = fmt.Sprintf("http://%s", config.URL)
	}
	if config.Logger == nil {
		config.Logger = g.Log()
	}
	return &implementer{
		config: config,
	}
}

`

const TemplateGenCtrlSdkIClient = `
// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package {PkgName}

import (
)

type IClient interface {
}
`

const TemplateGenCtrlSdkImplementer = `
// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package {PkgName}

import (
	"context"

	"github.com/gogf/gf/contrib/sdk/httpclient/v2"
	"github.com/gogf/gf/v2/text/gstr"

{ImportPaths}
)

type implementer{ImplementerName} struct {
	*httpclient.Client
}

`

const TemplateGenCtrlSdkImplementerNew = `
func (i *implementer) {ImplementerName}() {Module}.I{ImplementerName} {
	var (
		client = httpclient.New(i.config)
		prefix = gstr.TrimRight(i.config.URL, "/") + "{VersionPrefix}"
	)
	client.Client = client.Prefix(prefix)
	return &implementer{ImplementerName}{client}
}
`

const TemplateGenCtrlSdkImplementerFunc = `
func (i *implementer{ImplementerName}) {MethodName}(ctx context.Context, req *{Version}.{MethodName}Req) (res *{Version}.{MethodName}Res, err error) {
	err = i.Request(ctx, req, &res)
	return
}
`
