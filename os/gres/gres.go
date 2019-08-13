// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

const (
	gPACKAGE_TEMPLATE = `package %s

import "github.com/gogf/gf/os/gres"

func init() {
	if err := gres.Add(%s); err != nil {
		panic(err)
	}
}
`
)

var (
	defaultResource = New()
)

func Add(content []byte, prefix ...string) error {
	return defaultResource.Add(content, prefix...)
}

func Load(path string, prefix ...string) error {
	return defaultResource.Load(path, prefix...)
}

func Scan(path string, pattern string, recursive ...bool) []*File {
	return defaultResource.Scan(path, pattern, recursive...)
}

func Dump() {
	defaultResource.Dump()
}
