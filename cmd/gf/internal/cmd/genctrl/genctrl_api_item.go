// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genctrl

import "github.com/gogf/gf/v2/text/gstr"

type apiItem struct {
	Import     string `eg:"demo.com/api/user/v1"`
	Module     string `eg:"user"`
	Version    string `eg:"v1"`
	MethodName string `eg:"GetList"`
}

func (a apiItem) String() string {
	return gstr.Join([]string{
		a.Import, a.Module, a.Version, a.MethodName,
	}, ",")
}
