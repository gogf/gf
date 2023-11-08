// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gbuild"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
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
	var (
		version   = gf.VERSION
		gfVersion string
		envDetail string
		cliDetail = getBuildDetail()
		docDetail string
	)

	goVersion, ok := getGoVersion()
	if ok {
		gfVersion = getGoFrameVersion()
		goVersion = "\n  Go Version: " + goVersion
	} else {
		v, err := c.getGFVersionOfCurrentProject()
		if err == nil {
			gfVersion = v
		}
	}

	envDetail = fmt.Sprintf(`
Env Detail:
  CLI Installed At: %s%s
  GoFrame Version(go.mod): %s
`, gfile.SelfPath(), goVersion, gfVersion)

	docDetail = fmt.Sprintf(`
Others Detail:
  Docs: https://goframe.org
  Now Time: %s`, time.Now().Format("2006-01-02 15:04:05"))

	mlog.Print(version +
		envDetail +
		cliDetail +
		docDetail)

	return nil, nil
}

// getGoFrameVersion returns the goframe version of current project using.
func getGoFrameVersion() (gfVersion string) {
	pkgInfo, err := gproc.ShellExec(context.Background(), `go list -f "{{if (not .Main)}}{{.Path}}@{{.Version}}{{end}}" -m all`)
	if err != nil {
		return ""
	}
	pkgList := gstr.Split(pkgInfo, "\n")
	for _, v := range pkgList {
		if strings.HasPrefix(v, "github.com/gogf/gf") {
			gfVersion += fmt.Sprintf("\n    %s", v)
		}
	}
	return
}

// getGoVersion returns the go version
func getGoVersion() (goVersion string, ok bool) {
	goVersion, err := gproc.ShellExec(context.Background(), "go version")
	if err != nil {
		return "", false
	}
	goVersion = gstr.TrimLeftStr(goVersion, "go version ")
	goVersion = gstr.TrimRightStr(goVersion, "\n")
	return goVersion, true
}

// getBuildDetail returns the build information of current binary.
func getBuildDetail() (cliDetail string) {
	cliDetail = fmt.Sprintf(`
GoFrame CLI Build Detail:`)

	// build info
	info := gbuild.Info()
	if info.Git == "" {
		info.Git = "none"
	}
	if info.GoFrame == "" {
		cliDetail += fmt.Sprintf(`
  Go Version: %s
`, runtime.Version())
		return
	}

	cliDetail += fmt.Sprintf(`
  Go Version:  %s
  GF Version:  %s
  Git Commit:  %s
  Build Time:  %s
`, info.Golang, info.GoFrame, info.Git, info.Time)
	return
}

// getGFVersionOfCurrentProject checks and returns the GoFrame version current project using.
func (c cVersion) getGFVersionOfCurrentProject() (string, error) {
	goModPath := gfile.Join(gfile.Pwd(), "go.mod")
	if gfile.Exists(goModPath) {
		lines := gstr.SplitAndTrim(gfile.GetContents(goModPath), "\n")
		for _, line := range lines {
			line = gstr.Trim(line)
			line = gstr.TrimLeftStr(line, "require ")
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
