// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gtag"
)

var (
	Gen = cGen{}
)

type cGen struct {
	g.Meta `name:"gen" brief:"{cGenBrief}" dc:"{cGenDc}"`
	cGenDao
	cGenEnums
	cGenCtrl
	cGenPb
	cGenPbEntity
	cGenService
}

const (
	cGenBrief = `automatically generate go files for dao/do/entity/pb/pbentity`
	cGenDc    = `
The "gen" command is designed for multiple generating purposes. 
It's currently supporting generating go files for ORM models, protobuf and protobuf entity files.
Please use "gf gen dao -h" for specified type help.
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cGenBrief`: cGenBrief,
		`cGenDc`:    cGenDc,
	})
}
