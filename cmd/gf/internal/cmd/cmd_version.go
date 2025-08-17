// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"bytes"
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

const (
	defaultIndent = "{{indent}}"
)

type cVersion struct {
	g.Meta `name:"version" brief:"show version information of current binary"`
}

type cVersionInput struct {
	g.Meta `name:"version"`
}

type cVersionOutput struct{}

func (c cVersion) Index(ctx context.Context, in cVersionInput) (*cVersionOutput, error) {
	detailBuffer := &detailBuffer{}
	detailBuffer.WriteString(fmt.Sprintf("%s", gf.VERSION))

	detailBuffer.appendLine(0, "Welcome to GoFrame!")

	detailBuffer.appendLine(0, "Env Detail:")
	goVersion, ok := getGoVersion()
	if ok {
		detailBuffer.appendLine(1, fmt.Sprintf("Go Version: %s", goVersion))
		detailBuffer.appendLine(1, fmt.Sprintf("GF Version(go.mod): %s", getGoFrameVersion(2)))
	} else {
		v, err := c.getGFVersionOfCurrentProject()
		if err == nil {
			detailBuffer.appendLine(1, fmt.Sprintf("GF Version(go.mod): %s", v))
		} else {
			detailBuffer.appendLine(1, fmt.Sprintf("GF Version(go.mod): %s", err.Error()))
		}
	}

	detailBuffer.appendLine(0, "CLI Detail:")
	detailBuffer.appendLine(1, fmt.Sprintf("Installed At: %s", gfile.SelfPath()))
	info := gbuild.Info()
	if info.GoFrame == "" {
		detailBuffer.appendLine(1, fmt.Sprintf("Built Go Version: %s", runtime.Version()))
		detailBuffer.appendLine(1, fmt.Sprintf("Built GF Version: %s", gf.VERSION))
	} else {
		if info.Git == "" {
			info.Git = "none"
		}
		detailBuffer.appendLine(1, fmt.Sprintf("Built Go Version: %s", info.Golang))
		detailBuffer.appendLine(1, fmt.Sprintf("Built GF Version: %s", info.GoFrame))
		detailBuffer.appendLine(1, fmt.Sprintf("Git Commit: %s", info.Git))
		detailBuffer.appendLine(1, fmt.Sprintf("Built Time: %s", info.Time))
	}

	detailBuffer.appendLine(0, "Others Detail:")
	detailBuffer.appendLine(1, "Docs: https://goframe.org")
	detailBuffer.appendLine(1, fmt.Sprintf("Now : %s", time.Now().Format(time.RFC3339)))

	mlog.Print(detailBuffer.replaceAllIndent("  "))
	return nil, nil
}

// detailBuffer is a buffer for detail information.
type detailBuffer struct {
	bytes.Buffer
}

// appendLine appends a line to the buffer with given indent level.
func (d *detailBuffer) appendLine(indentLevel int, line string) {
	d.WriteString(fmt.Sprintf("\n%s%s", strings.Repeat(defaultIndent, indentLevel), line))
}

// replaceAllIndent replaces the tab with given indent string and prints the buffer content.
func (d *detailBuffer) replaceAllIndent(indentStr string) string {
	return strings.ReplaceAll(d.String(), defaultIndent, indentStr)
}

// getGoFrameVersion returns the goframe version of current project using.
func getGoFrameVersion(indentLevel int) (gfVersion string) {
	pkgInfo, err := gproc.ShellExec(context.Background(), `go list -f "{{if (not .Main)}}{{.Path}}@{{.Version}}{{end}}" -m all`)
	if err != nil {
		return "cannot find go.mod"
	}
	pkgList := gstr.Split(pkgInfo, "\n")
	for _, v := range pkgList {
		if strings.HasPrefix(v, "github.com/gogf/gf") {
			gfVersion += fmt.Sprintf("\n%s%s", strings.Repeat(defaultIndent, indentLevel), v)
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
