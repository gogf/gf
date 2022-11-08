package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
)

var (
	Up = cUp{}
)

type cUp struct {
	g.Meta `name:"up" brief:"upgrade GoFrame version/tool to latest one in current project" eg:"{cUpEg}" `
}

const (
	gfPackagePrefix = `github.com/gogf/gf/`
	cUpEg           = `
gf up
gf up -a
gf up -c
gf up -f -c
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cUpEg`: cUpEg,
	})
}

type cUpInput struct {
	g.Meta `name:"up"  config:"gfcli.up"`
	All    bool `name:"all" short:"a" brief:"upgrade both version and cli, auto fix codes" orphan:"true"`
	Fix    bool `name:"fix" short:"f" brief:"auto fix codes" orphan:"true"`
	Cli    bool `name:"cli" short:"c" brief:"also upgrade CLI tool (not supported yet)" orphan:"true"`
}

type cUpOutput struct{}

func (c cUp) Index(ctx context.Context, in cUpInput) (out *cUpOutput, err error) {
	defer mlog.Print(`done!`)

	if in.All {
		in.Cli = true
		in.Fix = true
	}
	if err = c.doUpgradeVersion(ctx); err != nil {
		return nil, err
	}
	if in.Fix {
		if err = c.doAutoFixing(ctx); err != nil {
			return nil, err
		}
	}
	//if in.Cli {
	//	if err = c.doUpgradeCLI(ctx); err != nil {
	//		return nil, err
	//	}
	//}
	return
}

func (c cUp) doUpgradeVersion(ctx context.Context) (err error) {
	mlog.Print(`start upgrading version...`)

	var (
		dir  = gfile.Pwd()
		temp string
		path = gfile.Join(dir, "go.mod")
	)
	for {
		if gfile.Exists(path) {
			var packages []string
			err = gfile.ReadLines(path, func(line string) error {
				line = gstr.Trim(line)
				if gstr.HasPrefix(line, gfPackagePrefix) {
					pkg := gstr.Explode(" ", line)[0]
					packages = append(packages, pkg)
				}
				return nil
			})
			if err != nil {
				return
			}
			for _, pkg := range packages {
				mlog.Printf(`upgrading %s`, pkg)
				command := fmt.Sprintf(`go get -u %s@latest`, pkg)
				if err = gproc.ShellRun(ctx, command); err != nil {
					return
				}
			}
			return
		}
		temp = gfile.Dir(dir)
		if temp == "" || temp == dir {
			return
		}
		dir = temp
		path = gfile.Join(dir, "go.mod")
	}
}

func (c cUp) doUpgradeCLI(ctx context.Context) (err error) {
	mlog.Print(`start upgrading cli...`)

	return
}

func (c cUp) doAutoFixing(ctx context.Context) (err error) {
	mlog.Print(`start auto fixing...`)
	err = cFix{}.doFix(ctx)
	return
}
