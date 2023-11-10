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
	detail string
}

type cVersionInput struct {
	g.Meta `name:"version"`
}

type cVersionOutput struct{}

func (c cVersion) Index(ctx context.Context, in cVersionInput) (*cVersionOutput, error) {
	c.detail = fmt.Sprintf("%s", gf.VERSION)

	c.appendLine(0, "Welcome to GoFrame!")

	c.appendLine(0, "Env Detail:")
	goVersion, ok := getGoVersion()
	if ok {
		c.appendLine(1, fmt.Sprintf("Go Version: %s", goVersion))
		c.appendLine(1, fmt.Sprintf("GF Version(go.mod): %s", getGoFrameVersion()))
	} else {
		v, err := c.getGFVersionOfCurrentProject()
		if err == nil {
			c.appendLine(1, fmt.Sprintf("GF Version(go.mod): %s", v))
		} else {
			c.appendLine(1, fmt.Sprintf("GF Version(go.mod): %s", err.Error()))
		}
	}

	c.appendLine(0, "CLI Detail:")
	c.appendLine(1, fmt.Sprintf("Installed At: %s", gfile.SelfPath()))
	info := gbuild.Info()
	if info.GoFrame == "" {
		c.appendLine(1, fmt.Sprintf("Builded Go Version: %s", runtime.Version()))
		c.appendLine(1, fmt.Sprintf("Builded GF Version: %s", gf.VERSION))
	} else {
		if info.Git == "" {
			info.Git = "none"
		}
		c.appendLine(1, fmt.Sprintf("Builded Go Version: %s", info.Golang))
		c.appendLine(1, fmt.Sprintf("Builded GF Version: %s", info.GoFrame))
		c.appendLine(1, fmt.Sprintf("Git Commit: %s", info.Git))
		c.appendLine(1, fmt.Sprintf("Builded Time: %s", info.Time))
	}

	c.appendLine(0, "Others Detail:")
	c.appendLine(1, "Docs: https://goframe.org")
	c.appendLine(1, fmt.Sprintf("Now Time: %s", time.Now().Format("2006-01-02 15:04:05")))
	c.print("  ")

	return nil, nil
}

// appendLine description
func (c *cVersion) appendLine(level int, line string) {
	c.detail += "\n" + strings.Repeat("\t", level) + line
}

// print description
func (c *cVersion) print(indent string) {
	c.detail = strings.ReplaceAll(c.detail, "\t", indent)
	mlog.Print(c.detail)
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
			gfVersion += fmt.Sprintf("\n\t\t%s", v)
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
