package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gbuild"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

var (
	Version = cVersion{}
)

type cVersion struct {
	g.Meta `name:"version" brief:"show version information of current binary"`
}

type cVersionInput struct {
	g.Meta `name:"version"`
}

type cVersionOutput struct{}

func (c cVersion) Index(ctx context.Context, in cVersionInput) (*cVersionOutput, error) {
	info := gbuild.Info()
	if info.Git == "" {
		info.Git = "none"
	}
	mlog.Printf(`GoFrame CLI Tool %s, https://goframe.org`, gf.VERSION)
	gfVersion, err := c.getGFVersionOfCurrentProject()
	if err != nil {
		gfVersion = err.Error()
	} else {
		gfVersion = gfVersion + " in current go.mod"
	}
	mlog.Printf(`GoFrame Version: %s`, gfVersion)
	mlog.Printf(`CLI Installed At: %s`, gfile.SelfPath())
	if info.GoFrame == "" {
		mlog.Print(`Current is a custom installed version, no installation information.`)
		return nil, nil
	}

	mlog.Print(gstr.Trim(fmt.Sprintf(`
CLI Built Detail:
  Go Version:  %s
  GF Version:  %s
  Git Commit:  %s
  Build Time:  %s
`, info.Golang, info.GoFrame, info.Git, info.Time)))
	return nil, nil
}

// getGFVersionOfCurrentProject checks and returns the GoFrame version current project using.
func (c cVersion) getGFVersionOfCurrentProject() (string, error) {
	goModPath := gfile.Join(gfile.Pwd(), "go.mod")
	if gfile.Exists(goModPath) {
		lines := gstr.SplitAndTrim(gfile.GetContents(goModPath), "\n")
		for _, line := range lines {
			line = gstr.Trim(line)
			// Version 1.
			match, err := gregex.MatchString(`^github\.com/gogf/gf\s+(.+)$`, line)
			if err != nil {
				return "", err
			}
			if len(match) <= 1 {
				// Version > 1.
				match, err = gregex.MatchString(`^github\.com/gogf/gf/v\d\s+(.+)$`, line)
				if err != nil {
					return "", err
				}
			}
			if len(match) > 1 {
				return gstr.Trim(match[1]), nil
			}
		}

		return "", gerror.New("cannot find goframe requirement in go.mod")
	} else {
		return "", gerror.New("cannot find go.mod")
	}
}
