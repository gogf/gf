package cmd

import (
	"context"

	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:        "tpl",
		Brief:       "Project scaffolding tool",
		Description: "A CLI tool for generating Go projects from remote templates",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// parser.Command.Print()
            // Just hint user
            println("Please use 'tpl init' to generate project, or 'tpl init -h' for help.")
			return nil
		},
	}
)
