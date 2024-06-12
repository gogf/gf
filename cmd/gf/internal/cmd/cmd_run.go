// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gfsnotify"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

var (
	Run = cRun{}
)

type cRun struct {
	g.Meta `name:"run" usage:"{cRunUsage}" brief:"{cRunBrief}" eg:"{cRunEg}" dc:"{cRunDc}"`
}

type cRunApp struct {
	File       string   // Go run file name.
	Path       string   // Directory storing built binary.
	Options    string   // Extra "go run" options.
	Args       string   // Custom arguments.
	WatchPaths []string // Watch paths for live reload.
}

const (
	cRunUsage = `gf run FILE [OPTION]`
	cRunBrief = `running go codes with hot-compiled-like feature`
	cRunEg    = `
gf run main.go
gf run main.go --args "server -p 8080"
gf run main.go -mod=vendor
gf run main.go -w "manifest/config/*.yaml"
`
	cRunDc = `
The "run" command is used for running go codes with hot-compiled-like feature,
which compiles and runs the go codes asynchronously when codes change.
`
	cRunFileBrief       = `building file path.`
	cRunPathBrief       = `output directory path for built binary file. it's "./" in default`
	cRunExtraBrief      = `the same options as "go run"/"go build" except some options as follows defined`
	cRunArgsBrief       = `custom arguments for your process`
	cRunWatchPathsBrief = `watch additional paths for live reload, separated by ",". i.e. "manifest/config/*.yaml"`
)

var (
	process *gproc.Process
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cRunUsage`:           cRunUsage,
		`cRunBrief`:           cRunBrief,
		`cRunEg`:              cRunEg,
		`cRunDc`:              cRunDc,
		`cRunFileBrief`:       cRunFileBrief,
		`cRunPathBrief`:       cRunPathBrief,
		`cRunExtraBrief`:      cRunExtraBrief,
		`cRunArgsBrief`:       cRunArgsBrief,
		`cRunWatchPathsBrief`: cRunWatchPathsBrief,
	})
}

type (
	cRunInput struct {
		g.Meta     `name:"run" config:"gfcli.run"`
		File       string   `name:"FILE"       arg:"true" brief:"{cRunFileBrief}" v:"required"`
		Path       string   `name:"path"       short:"p"  brief:"{cRunPathBrief}" d:"./"`
		Extra      string   `name:"extra"      short:"e"  brief:"{cRunExtraBrief}"`
		Args       string   `name:"args"       short:"a"  brief:"{cRunArgsBrief}"`
		WatchPaths []string `name:"watchPaths" short:"w"  brief:"{cRunWatchPathsBrief}"`
	}
	cRunOutput struct{}
)

func (c cRun) Index(ctx context.Context, in cRunInput) (out *cRunOutput, err error) {
	// Necessary check.
	if gproc.SearchBinary("go") == "" {
		mlog.Fatalf(`command "go" not found in your environment, please install golang first to proceed this command`)
	}

	if len(in.WatchPaths) == 1 {
		in.WatchPaths = strings.Split(in.WatchPaths[0], ",")
		mlog.Printf("watchPaths: %v", in.WatchPaths)
	}

	app := &cRunApp{
		File:       in.File,
		Path:       filepath.FromSlash(in.Path),
		Options:    in.Extra,
		Args:       in.Args,
		WatchPaths: in.WatchPaths,
	}
	dirty := gtype.NewBool()

	var outputPath = app.genOutputPath()
	callbackFunc := func(event *gfsnotify.Event) {
		if gfile.ExtName(event.Path) != "go" {
			return
		}

		// Variable `dirty` is used for running the changes only one in one second.
		if !dirty.Cas(false, true) {
			return
		}

		// With some delay in case of multiple code changes in very short interval.
		gtimer.SetTimeout(ctx, 1500*gtime.MS, func(ctx context.Context) {
			defer dirty.Set(false)
			mlog.Printf(`watched file changes: %s`, event.String())
			app.Run(ctx, outputPath)
		})
	}

	if len(app.WatchPaths) > 0 {
		for _, path := range app.WatchPaths {
			_, err = gfsnotify.Add(gfile.RealPath(path), callbackFunc)
			if err != nil {
				mlog.Fatal(err)
			}
		}
	} else {
		_, err = gfsnotify.Add(gfile.RealPath("."), callbackFunc)
		if err != nil {
			mlog.Fatal(err)
		}
	}

	go app.Run(ctx, outputPath)

	gproc.AddSigHandlerShutdown(func(sig os.Signal) {
		app.End(ctx, sig, outputPath)
		os.Exit(0)
	})
	gproc.Listen()

	select {}
}

func (app *cRunApp) Run(ctx context.Context, outputPath string) {
	// Rebuild and run the codes.
	mlog.Printf("build: %s", app.File)

	// In case of `pipe: too many open files` error.
	// Build the app.
	buildCommand := fmt.Sprintf(
		`go build -o %s %s %s`,
		outputPath,
		app.Options,
		app.File,
	)
	mlog.Print(buildCommand)
	result, err := gproc.ShellExec(ctx, buildCommand)
	if err != nil {
		mlog.Printf("build error: \n%s%s", result, err.Error())
		return
	}
	// Kill the old process if build successfully.
	if process != nil {
		if err := process.Kill(); err != nil {
			mlog.Debugf("kill process error: %s", err.Error())
		}
	}
	// Run the binary file.
	runCommand := fmt.Sprintf(`%s %s`, outputPath, app.Args)
	mlog.Print(runCommand)
	if runtime.GOOS == "windows" {
		// Special handling for windows platform.
		// DO NOT USE "cmd /c" command.
		process = gproc.NewProcess(outputPath, strings.Fields(app.Args))
	} else {
		process = gproc.NewProcessCmd(runCommand, nil)
	}
	if pid, err := process.Start(ctx); err != nil {
		mlog.Printf("build running error: %s", err.Error())
	} else {
		mlog.Printf("build running pid: %d", pid)
	}
}

func (app *cRunApp) End(ctx context.Context, sig os.Signal, outputPath string) {
	// Delete the binary file.
	// firstly, kill the process.
	if process != nil {
		if err := process.Kill(); err != nil {
			mlog.Debugf("kill process error: %s", err.Error())
		}
	}
	if err := gfile.Remove(outputPath); err != nil {
		mlog.Printf("delete binary file error: %s", err.Error())
	} else {
		mlog.Printf("deleted binary file: %s", outputPath)
	}
}

func (app *cRunApp) genOutputPath() (outputPath string) {
	var renamePath string
	outputPath = gfile.Join(app.Path, gfile.Name(app.File))
	if runtime.GOOS == "windows" {
		outputPath += ".exe"
		if gfile.Exists(outputPath) {
			renamePath = outputPath + "~"
			if err := gfile.Rename(outputPath, renamePath); err != nil {
				mlog.Print(err)
			}
		}
	}
	return filepath.FromSlash(outputPath)
}

func matchWatchPaths(watchPaths []string, eventPath string) bool {
	for _, path := range watchPaths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			mlog.Printf("match watchPath '%s' error: %s", path, err.Error())
			continue
		}
		matched, err := filepath.Match(absPath, eventPath)
		if err != nil {
			mlog.Printf("match watchPath '%s' error: %s", path, err.Error())
			continue
		}
		if matched {
			return true
		}
	}
	return false
}
