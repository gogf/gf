// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfcmd

import (
	"context"
	"runtime"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd"
	_ "github.com/gogf/gf/cmd/gf/v2/internal/packed"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/allyes"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

const (
	cliFolderName = `hack`
)

// Command manages the CLI command of `gf`.
// This struct can be globally accessible and extended with custom struct.
type Command struct {
	*gcmd.Command
}

// Run starts running the command according the command line arguments and options.
func (c *Command) Run(ctx context.Context) {
	defer func() {
		if exception := recover(); exception != nil {
			if err, ok := exception.(error); ok {
				mlog.Print(err.Error())
			} else {
				panic(gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception))
			}
		}
	}()

	// CLI configuration, using the `hack/config.yaml` in priority.
	if path, _ := gfile.Search(cliFolderName); path != "" {
		if adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
			if err := adapter.SetPath(path); err != nil {
				mlog.Fatal(err)
			}
		}
	}

	// zsh alias "git fetch" conflicts checks.
	handleZshAlias()

	// -y option checks.
	allyes.Init()

	// just run.
	if err := c.RunWithError(ctx); err != nil {
		// Exit with error message and exit code 1.
		// It is very important to exit the command process with code 1.
		mlog.Fatalf(`%+v`, err)
	}
}

// GetCommand retrieves and returns the root command of CLI `gf`.
func GetCommand(ctx context.Context) (*Command, error) {
	root, err := gcmd.NewFromObject(cmd.GF)
	if err != nil {
		panic(err)
	}
	err = root.AddObject(
		cmd.Up,
		cmd.Env,
		cmd.Fix,
		cmd.Run,
		cmd.Gen,
		cmd.Tpl,
		cmd.Init,
		cmd.Pack,
		cmd.Build,
		cmd.Docker,
		cmd.Install,
		cmd.Version,
		cmd.Doc,
	)
	if err != nil {
		return nil, err
	}
	command := &Command{
		root,
	}
	return command, nil
}

// zsh alias "git fetch" conflicts checks.
func handleZshAlias() {
	if runtime.GOOS == "windows" {
		return
	}
	if home, err := gfile.Home(); err == nil {
		zshPath := gfile.Join(home, ".zshrc")
		if gfile.Exists(zshPath) {
			var (
				aliasCommand = `alias gf=gf`
				content      = gfile.GetContents(zshPath)
			)
			if !gstr.Contains(content, aliasCommand) {
				_ = gfile.PutContentsAppend(zshPath, "\n"+aliasCommand+"\n")
			}
		}
	}
}
