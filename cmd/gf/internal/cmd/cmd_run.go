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
	"time"

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

var Run = cRun{}

type cRun struct {
	g.Meta `name:"run" usage:"{cRunUsage}" brief:"{cRunBrief}" eg:"{cRunEg}" dc:"{cRunDc}"`
}

type cRunApp struct {
	File        string   // Go run file name.
	Path        string   // Directory storing built binary.
	Options     string   // Extra "go run" options.
	Args        string   // Custom arguments.
	WatchPaths  []string // Watch paths for live reload.
	IgnorePaths []string // Ignore paths for file watching.
}

const (
	cRunUsage = `gf run FILE [OPTION]`
	cRunBrief = `running go codes with hot-compiled-like feature`
	cRunEg    = `
gf run main.go
gf run main.go --args "server -p 8080"
gf run main.go -mod=vendor
gf run main.go -w internal/,api/
gf run main.go -i vendor/*,*.pb.go,node_modules/*
gf run main.go -w app/,manifest/ -i .git/*,.github/*,dist/*
gf run main.go -p ./bin -w . -i "test/*,tmp/*,*.log"
gf run main.go -w service/,model/ -i "frontend/*,web/*,build/*"
`
	cRunDc = `
The "run" command is used for running go codes with hot-compiled-like feature,
which compiles and runs the go codes asynchronously when codes change.
`
	cRunFileBrief        = `building file path.`
	cRunPathBrief        = `output directory path for built binary file. it's "./" in default`
	cRunExtraBrief       = `the same options as "go run"/"go build" except some options as follows defined`
	cRunArgsBrief        = `custom arguments for your process`
	cRunWatchPathsBrief  = `watch additional paths for live reload, separated by ",". i.e. "internal/,api/"`
	cRunIgnorePathsBrief = `ignore paths for file watching, separated by ",". i.e. "vendor/*,*.pb.go,node_modules/*,.git/*"`
)

var process *gproc.Process

func init() {
	gtag.Sets(g.MapStrStr{
		`cRunUsage`:            cRunUsage,
		`cRunBrief`:            cRunBrief,
		`cRunEg`:               cRunEg,
		`cRunDc`:               cRunDc,
		`cRunFileBrief`:        cRunFileBrief,
		`cRunPathBrief`:        cRunPathBrief,
		`cRunExtraBrief`:       cRunExtraBrief,
		`cRunArgsBrief`:        cRunArgsBrief,
		`cRunWatchPathsBrief`:  cRunWatchPathsBrief,
		`cRunIgnorePathsBrief`: cRunIgnorePathsBrief,
	})
}

type (
	cRunInput struct {
		g.Meta      `name:"run" config:"gfcli.run"`
		File        string   `name:"FILE"        arg:"true" brief:"{cRunFileBrief}" v:"required"`
		Path        string   `name:"path"        short:"p"  brief:"{cRunPathBrief}" d:"./"`
		Extra       string   `name:"extra"       short:"e"  brief:"{cRunExtraBrief}"`
		Args        string   `name:"args"        short:"a"  brief:"{cRunArgsBrief}"`
		WatchPaths  []string `name:"watchPaths"  short:"w"  brief:"{cRunWatchPathsBrief}"`
		IgnorePaths []string `name:"ignorePaths" short:"i"  brief:"{cRunIgnorePathsBrief}"`
	}
	cRunOutput struct{}
)

func (c cRun) Index(ctx context.Context, in cRunInput) (out *cRunOutput, err error) {
	if !gfile.Exists(in.File) {
		mlog.Fatalf(`given file "%s" not found`, in.File)
	}
	if !gfile.IsFile(in.File) {
		mlog.Fatalf(`given "%s" is not a file`, in.File)
	}
	// Necessary check.
	if gproc.SearchBinary("go") == "" {
		mlog.Fatalf(`command "go" not found in your environment, please install golang first to proceed this command`)
	}

	if len(in.WatchPaths) == 1 {
		in.WatchPaths = strings.Split(in.WatchPaths[0], ",")
		mlog.Printf("watchPaths: %v", in.WatchPaths)
	}

	if len(in.IgnorePaths) == 1 {
		in.IgnorePaths = strings.Split(in.IgnorePaths[0], ",")
		mlog.Printf("ignorePaths: %v", in.IgnorePaths)
	}

	app := &cRunApp{
		File:        in.File,
		Path:        filepath.FromSlash(in.Path),
		Options:     in.Extra,
		Args:        in.Args,
		WatchPaths:  in.WatchPaths,
		IgnorePaths: in.IgnorePaths,
	}
	dirty := gtype.NewBool()

	outputPath := app.genOutputPath()
	callbackFunc := func(event *gfsnotify.Event) {
		if !event.IsWrite() && !event.IsCreate() && !event.IsRemove() && !event.IsRename() {
			return
		}

		// Check if the file extension is 'go'.
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

	// Get all paths to watch after filtering.
	watchPaths := app.getWatchPaths()
	for _, path := range watchPaths {
		_, err = gfsnotify.Add(gfile.RealPath(path), callbackFunc)
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
		if sig != nil && runtime.GOOS != "windows" {
			if err := process.Signal(sig); err != nil {
				mlog.Debugf("send signal to process error: %s", err.Error())
				if err := process.Kill(); err != nil {
					mlog.Debugf("kill process error: %s", err.Error())
				}
			} else {
				waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
				defer cancel()
				done := make(chan error, 1)
				go func() {
					select {
					case <-waitCtx.Done():
						done <- waitCtx.Err()
					case done <- process.Wait():
					}
				}()
				err := <-done
				if err != nil {
					mlog.Debugf("process wait error: %s", err.Error())
					if err := process.Kill(); err != nil {
						mlog.Debugf("kill process error: %s", err.Error())
					}
				} else {
					mlog.Debug("process exited gracefully")
				}
			}
		} else {
			if err := process.Kill(); err != nil {
				mlog.Debugf("kill process error: %s", err.Error())
			}
		}
	}
	if err := gfile.RemoveFile(outputPath); err != nil {
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

// getWatchPaths returns all paths to watch after filtering with ignore patterns.
func (app *cRunApp) getWatchPaths() []string {
	// Collect all root paths to scan.
	roots := []string{"."}
	if len(app.WatchPaths) > 0 {
		roots = app.WatchPaths
	}

	// Scan all paths and filter with ignore patterns.
	var watchPaths []string
	seen := make(map[string]bool)

	for _, root := range roots {
		absRoot := gfile.RealPath(root)
		if absRoot == "" {
			mlog.Printf("watch path '%s' not found, skipping", root)
			continue
		}

		// Check if the root itself should be ignored.
		if app.shouldIgnorePath(absRoot) {
			mlog.Printf("ignoring path: %s", absRoot)
			continue
		}

		// Scan directory recursively with custom filtering to avoid scanning ignored directories.
		files, err := app.scanDirWithFilter(absRoot, "*", true)
		if err != nil {
			mlog.Printf("scan directory '%s' error: %s", absRoot, err.Error())
			continue
		}

		// Filter files with ignore patterns.
		for _, file := range files {
			if seen[file] {
				continue
			}
			seen[file] = true

			if !app.shouldIgnorePath(file) {
				watchPaths = append(watchPaths, file)
			} else {
				mlog.Printf("ignoring path: %s", file)
			}
		}
	}

	if len(watchPaths) == 0 {
		mlog.Printf("no paths to watch after filtering, watching current directory")
		return []string{"."}
	}

	mlog.Printf("watching %d paths", len(watchPaths))
	// for _, v := range watchPaths {
	// 	mlog.Printf("path: %s", v)
	// }
	return watchPaths
}

// scanDirWithFilter scans directory recursively but skips ignored directories to improve performance.
func (app *cRunApp) scanDirWithFilter(path string, pattern string, recursive bool) ([]string, error) {
	if !recursive {
		return gfile.ScanDir(path, pattern, false)
	}

	var result []string
	files, err := gfile.ScanDir(path, pattern, false)
	if err != nil {
		return nil, err
	}
	result = append(result, files...)

	// Get subdirectories
	subDirs, err := gfile.ScanDir(path, "*", false)
	if err != nil {
		return nil, err
	}

	for _, subDir := range subDirs {
		if !gfile.IsDir(subDir) {
			continue
		}

		// Check if this directory should be ignored
		if app.shouldIgnorePath(subDir) {
			mlog.Printf("skipping ignored directory: %s", subDir)
			continue
		}

		// Recursively scan this directory
		subFiles, err := app.scanDirWithFilter(subDir, pattern, true)
		if err != nil {
			return nil, err
		}
		result = append(result, subFiles...)
	}

	return result, nil
}

// shouldIgnorePath checks if the given file path should be ignored based on ignore patterns.
func (app *cRunApp) shouldIgnorePath(path string) bool {
	if len(app.IgnorePaths) == 0 {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		mlog.Printf("convert path to absolute error: %s", err.Error())
		return false
	}

	// Get the file name for pattern matching.
	fileName := filepath.Base(path)

	for _, pattern := range app.IgnorePaths {
		// Normalize the pattern to handle "./" or "../" prefixes.
		normalizedPattern := filepath.Clean(pattern)

		// Match against file name.
		if matched, _ := filepath.Match(normalizedPattern, fileName); matched {
			return true
		}

		// Match against full path.
		if matched, _ := filepath.Match(normalizedPattern, absPath); matched {
			return true
		}

		// Match against relative path from the current working directory.
		if cwd, err := os.Getwd(); err == nil {
			if rel, err := filepath.Rel(cwd, absPath); err == nil {
				// Match against relative path.
				if matched, _ := filepath.Match(normalizedPattern, rel); matched {
					return true
				}

				// Check if the path starts with a directory pattern (e.g., "vendor/", "node_modules/")
				if strings.Contains(normalizedPattern, string(filepath.Separator)) {
					if strings.HasPrefix(rel, normalizedPattern) {
						return true
					}
				}

				// Check if any part of the path matches the pattern.
				pathParts := strings.Split(rel, string(filepath.Separator))
				for _, part := range pathParts {
					if matched, _ := filepath.Match(normalizedPattern, part); matched {
						return true
					}
				}
			}
		}
	}

	return false
}
