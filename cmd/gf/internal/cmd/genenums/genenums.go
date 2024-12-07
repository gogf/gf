// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genenums

import (
	"context"

	"golang.org/x/tools/go/packages"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

type (
	CGenEnums      struct{}
	CGenEnumsInput struct {
		g.Meta   `name:"enums" config:"{CGenEnumsConfig}" brief:"{CGenEnumsBrief}" eg:"{CGenEnumsEg}"`
		Src      string   `name:"src"      short:"s"  dc:"source folder path to be parsed" d:"api"`
		Path     string   `name:"path"     short:"p"  dc:"output go file path storing enums content" d:"internal/packed/packed_enums.go"`
		Prefixes []string `name:"prefixes" short:"x"  dc:"only exports packages that starts with specified prefixes"`
	}
	CGenEnumsOutput struct{}
)

const (
	CGenEnumsConfig = `gfcli.gen.enums`
	CGenEnumsBrief  = `parse go files in current project and generate enums go file`
	CGenEnumsEg     = `
gf gen enums
gf gen enums -p internal/packed/packed_enums.go
gf gen enums -p internal/packed/packed_enums.go -s .
gf gen enums -x github.com/gogf
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenEnumsEg`:     CGenEnumsEg,
		`CGenEnumsBrief`:  CGenEnumsBrief,
		`CGenEnumsConfig`: CGenEnumsConfig,
	})
}

func (c CGenEnums) Enums(ctx context.Context, in CGenEnumsInput) (out *CGenEnumsOutput, err error) {
	realPath := gfile.RealPath(in.Src)
	if realPath == "" {
		mlog.Fatalf(`source folder path "%s" does not exist`, in.Src)
	}
	err = gfile.Chdir(realPath)
	if err != nil {
		mlog.Fatal(err)
	}
	mlog.Printf(`scanning for enums: %s`, realPath)
	cfg := &packages.Config{
		Dir:   realPath,
		Mode:  pkgLoadMode,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg)
	if err != nil {
		mlog.Fatal(err)
	}
	p := NewEnumsParser(in.Prefixes)
	p.ParsePackages(pkgs)
	var enumsContent = gstr.ReplaceByMap(consts.TemplateGenEnums, g.MapStrStr{
		"{PackageName}": gfile.Basename(gfile.Dir(in.Path)),
		"{EnumsJson}":   "`" + p.Export() + "`",
	})
	enumsContent = gstr.Trim(enumsContent)
	if err = gfile.PutContents(in.Path, enumsContent); err != nil {
		return
	}
	mlog.Printf(`generated enums go file: %s`, in.Path)
	mlog.Print("done!")
	return
}
