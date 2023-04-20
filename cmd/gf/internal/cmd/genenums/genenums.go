package genenums

import (
	"context"
	"golang.org/x/tools/go/packages"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
)

type (
	CGenEnums      struct{}
	CGenEnumsInput struct {
		g.Meta `name:"enums" config:"{CGenEnumsConfig}" brief:"{CGenEnumsBrief}" eg:"{CGenEnumsEg}"`
		Src    string `name:"src"  short:"s"  dc:"source folder path to be parsed" d:"."`
		Path   string `name:"path" short:"p"  dc:"output go file path storing enums content" d:"internal/boot/boot_enums.go"`
	}
	CGenEnumsOutput struct{}
)

const (
	CGenEnumsConfig = `gfcli.gen.enums`
	CGenEnumsBrief  = `parse go files in current project and generate enums go file`
	CGenEnumsEg     = `
gf gen enums
gf gen enums -p internal/boot/boot_enums.go
gf gen enums -p internal/boot/boot_enums.go -s .
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
	mlog.Printf(`scanning: %s`, realPath)
	cfg := &packages.Config{
		Dir:   realPath,
		Mode:  pkgLoadMode,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg)
	if err != nil {
		mlog.Fatal(err)
	}
	p := NewEnumsParser()
	p.ParsePackages(pkgs)
	var enumsContent = gstr.ReplaceByMap(consts.TemplateGenEnums, g.MapStrStr{
		"{PackageName}": gfile.Basename(gfile.Dir(in.Path)),
		"{EnumsJson}":   "`" + p.Export() + "`",
	})
	if err = gfile.PutContents(in.Path, enumsContent); err != nil {
		return
	}
	mlog.Printf(`generated: %s`, in.Path)
	mlog.Print("done!")
	return
}
