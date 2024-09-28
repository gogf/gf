// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/service"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

var (
	GF = cGF{}
)

type cGF struct {
	g.Meta `name:"gf" ad:"{cGFAd}"`
}

const (
	cGFAd = `
ADDITIONAL
    Use "gf COMMAND -h" for details about a command.
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cGFAd`: cGFAd,
	})
}

type cGFInput struct {
	g.Meta  `name:"gf"`
	Yes     bool `short:"y" name:"yes"     brief:"all yes for all command without prompt ask"   orphan:"true"`
	Version bool `short:"v" name:"version" brief:"show version information of current binary"   orphan:"true"`
	Debug   bool `short:"d" name:"debug"   brief:"show internal detailed debugging information" orphan:"true"`
}

type cGFOutput struct{}

func (c cGF) Index(ctx context.Context, in cGFInput) (out *cGFOutput, err error) {
	// Version.
	if in.Version {
		_, err = Version.Index(ctx, cVersionInput{})
		return
	}

	answer := "n"
	// No argument or option, do installation checks.
	if data, isInstalled := service.Install.IsInstalled(); !isInstalled {
		mlog.Print("hi, it seems it's the first time you installing gf cli.")
		answer = gcmd.Scanf("do you want to install gf(%s) binary to your system? [y/n]: ", gf.VERSION)
	} else if !data.IsSelf {
		mlog.Print("hi, you have installed gf cli.")
		answer = gcmd.Scanf("do you want to install gf(%s) binary to your system? [y/n]: ", gf.VERSION)
	}
	if strings.EqualFold(answer, "y") {
		if err = service.Install.Run(ctx); err != nil {
			return
		}
		gcmd.Scan("press `Enter` to exit...")
		return
	}

	// Print help content.
	gcmd.CommandFromCtx(ctx).Print()
	return
}
