package tpl

import (
	"context"

	_ "github.com/gogf/gf/contrib/drivers/clickhouse/v2"
	_ "github.com/gogf/gf/contrib/drivers/mssql/v2"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/oracle/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

const (
	CGenTplConfig = `gfcli.gen.tpl`
	CGenTplUsage  = `gf gen tpl [OPTION]`
	CGenTplBrief  = `automatically generate template files`
	CGenTplEg     = `
gf gen tpl
gf gen tpl -t default -p ./template
`
	CGenTplAd = `
CONFIGURATION SUPPORT
    Options are also supported by configuration file.
    It's suggested using configuration file instead of command line arguments.
    The configuration node name is "gfcli.gen.tpl" 
`

	CGenTplBriefPath = `output directory path (default: "./template")`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenTplConfig`: CGenTplConfig,
		`CGenTplUsage`:  CGenTplUsage,
		`CGenTplBrief`:  CGenTplBrief,
		`CGenTplEg`:     CGenTplEg,
		`CGenTplAd`:     CGenTplAd,
	})
}

type (
	CGenTpl      struct{}
	CGenTplInput struct {
		g.Meta `name:"tpl" config:"{CGenTplConfig}" usage:"{CGenTplUsage}" brief:"{CGenTplBrief}" eg:"{CGenTplEg}" ad:"{CGenTplAd}"`
		Path   string `name:"path"    short:"p" brief:"{CGenTplBriefPath}" d:"./template"`
		Clear  bool   `name:"clear"   short:"c" brief:"delete old files before generation"`
	}
	CGenTplOutput struct{}
)

func (c CGenTpl) Tpl(ctx context.Context, in CGenTplInput) (out *CGenTplOutput, err error) {
	// Clear old files
	if in.Clear {
		if err := gfile.Remove(in.Path); err != nil {
			return nil, gerror.Wrapf(err, "clear output path failed")
		}
	}

	// Create output directory
	if !gfile.Exists(in.Path) {
		if err := gfile.Mkdir(in.Path); err != nil {
			return nil, gerror.Wrapf(err, "create output directory failed")
		}
	}

	mlog.Print("template files generated successfully!")
	return &CGenTplOutput{}, nil
}
