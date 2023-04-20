// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"bytes"
	"context"

	"github.com/olekukonko/tablewriter"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

var (
	Env = cEnv{}
)

type cEnv struct {
	g.Meta `name:"env" brief:"show current Golang environment variables"`
}

type cEnvInput struct {
	g.Meta `name:"env"`
}

type cEnvOutput struct{}

func (c cEnv) Index(ctx context.Context, in cEnvInput) (out *cEnvOutput, err error) {
	result, err := gproc.ShellExec(ctx, "go env")
	if err != nil {
		mlog.Fatal(err)
	}
	if result == "" {
		mlog.Fatal(`retrieving Golang environment variables failed, did you install Golang?`)
	}
	var (
		lines  = gstr.Split(result, "\n")
		buffer = bytes.NewBuffer(nil)
	)
	array := make([][]string, 0)
	for _, line := range lines {
		line = gstr.Trim(line)
		if line == "" {
			continue
		}
		if gstr.Pos(line, "set ") == 0 {
			line = line[4:]
		}
		match, _ := gregex.MatchString(`(.+?)=(.*)`, line)
		if len(match) < 3 {
			mlog.Fatalf(`invalid Golang environment variable: "%s"`, line)
		}
		array = append(array, []string{gstr.Trim(match[1]), gstr.Trim(match[2])})
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	tw.AppendBulk(array)
	tw.Render()
	mlog.Print(buffer.String())
	return
}
