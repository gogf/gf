// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genservice

type logicItem struct {
	Receiver    string              `eg:"sUser"`
	MethodName  string              `eg:"GetList"`
	InputParam  []map[string]string `eg:"ctx context.Context, cond *SearchInput"`
	OutputParam []map[string]string `eg:"list []*User, err error"`
	Comment     string              `eg:"Get user list"`
}
