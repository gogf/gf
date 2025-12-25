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

type watchPath struct {
	Path      string
	Recursive bool
}

type cRunApp struct {
	File           string   // Go run file name.
	Path           string   // Directory storing built binary.
	Options        string   // Extra "go run" options.
	Args           string   // Custom arguments.
	WatchPaths     []string // Watch paths for live reload.
	IgnorePatterns []string // Custom ignore patterns.
}

const (
	cRunUsage = `gf run FILE [OPTION]`
	cRunBrief = `running go codes with hot-compiled-like feature`
	cRunEg    = `
gf run main.go
gf run main.go --args "server -p 8080"
gf run main.go -mod=vendor
gf run main.go -w internal,api
gf run main.go -i ".git,node_modules"
`
	cRunDc = `
The "run" command is used for running go codes with hot-compiled-like feature,
which compiles and runs the go codes asynchronously when codes change.
`
	cRunFileBrief          = `building file path.`
	cRunPathBrief          = `output directory path for built binary file. it's "./" in default`
	cRunExtraBrief         = `the same options as "go run"/"go build" except some options as follows defined`
	cRunArgsBrief          = `custom arguments for your process`
	cRunWatchPathsBrief    = `watch additional paths for live reload, separated by ",". i.e. "internal,api"`
	cRunIgnorePatternBrief = `custom ignore patterns for watch, separated by ",". i.e. ".git,node_modules". default patterns: node_modules, vendor, .*, _*`
)

var process *gproc.Process

func init() {
	gtag.Sets(g.MapStrStr{
		`cRunUsage`:              cRunUsage,
		`cRunBrief`:              cRunBrief,
		`cRunEg`:                 cRunEg,
		`cRunDc`:                 cRunDc,
		`cRunFileBrief`:          cRunFileBrief,
		`cRunPathBrief`:          cRunPathBrief,
		`cRunExtraBrief`:         cRunExtraBrief,
		`cRunArgsBrief`:          cRunArgsBrief,
		`cRunWatchPathsBrief`:    cRunWatchPathsBrief,
		`cRunIgnorePatternBrief`: cRunIgnorePatternBrief,
	})
}

type (
	cRunInput struct {
		g.Meta         `name:"run" config:"gfcli.run"`
		File           string   `name:"FILE"           arg:"true" brief:"{cRunFileBrief}" v:"required"`
		Path           string   `name:"path"           short:"p"  brief:"{cRunPathBrief}" d:"./"`
		Extra          string   `name:"extra"          short:"e"  brief:"{cRunExtraBrief}"`
		Args           string   `name:"args"           short:"a"  brief:"{cRunArgsBrief}"`
		WatchPaths     []string `name:"watchPaths"     short:"w"  brief:"{cRunWatchPathsBrief}"`
		IgnorePatterns []string `name:"ignorePatterns" short:"i"  brief:"{cRunIgnorePatternBrief}"`
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
		parts := strings.Split(in.WatchPaths[0], ",")
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		in.WatchPaths = parts
		mlog.Printf("watchPaths: %v", in.WatchPaths)
	}

	if len(in.IgnorePatterns) == 1 {
		parts := strings.Split(in.IgnorePatterns[0], ",")
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		in.IgnorePatterns = parts
		mlog.Printf("ignorePatterns: %v", in.IgnorePatterns)
	}

	app := &cRunApp{
		File:           in.File,
		Path:           filepath.FromSlash(in.Path),
		Options:        in.Extra,
		Args:           in.Args,
		WatchPaths:     in.WatchPaths,
		IgnorePatterns: in.IgnorePatterns,
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

	// Get directories to watch (recursive or non-recursive monitoring).
	watchPaths := app.getWatchPaths()
	for _, wp := range watchPaths {
		mlog.Printf("watchPaths: %v", wp)
		option := gfsnotify.WatchOption{NoRecursive: !wp.Recursive}
		_, err = gfsnotify.Add(wp.Path, callbackFunc, option)
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

// getWatchPaths uses DFS to find the minimal set of directories to watch.
// Rule: if a directory and all its descendants have no ignored subdirectories, watch it;
// otherwise, recurse into valid children and watch the current directory non-recursively.
func (app *cRunApp) getWatchPaths() []watchPath {
	roots := []string{"."}
	if len(app.WatchPaths) > 0 {
		roots = app.WatchPaths
	}

	// Use custom ignore patterns if provided, otherwise use default.
	ignorePatterns := defaultIgnorePatterns
	if len(app.IgnorePatterns) > 0 {
		ignorePatterns = app.IgnorePatterns
	}

	var watchPaths []watchPath

	for _, root := range roots {
		absRoot := gfile.RealPath(root)
		if absRoot == "" {
			mlog.Printf("watch path '%s' not found, skipping", root)
			continue
		}
		if isIgnoredDirName(absRoot, ignorePatterns) {
			continue
		}
		app.collectWatchPaths(absRoot, ignorePatterns, &watchPaths)
	}

	if len(watchPaths) == 0 {
		mlog.Printf("no directories to watch, using current directory")
		if absCur := gfile.RealPath("."); absCur != "" {
			return []watchPath{{Path: absCur, Recursive: true}}
		}
		return []watchPath{{Path: ".", Recursive: true}}
	}

	mlog.Printf("watching %d paths", len(watchPaths))
	for _, wp := range watchPaths {
		recursiveStr := "recursive"
		if !wp.Recursive {
			recursiveStr = "non-recursive"
		}
		mlog.Debugf("  - %s (%s)", wp.Path, recursiveStr)
	}
	return watchPaths
}

// collectWatchPaths performs a DFS traversal to collect the minimal set of directories to watch.
func (app *cRunApp) collectWatchPaths(dir string, ignorePatterns []string, watchPaths *[]watchPath) {
	// Check if this directory or any immediate child is ignored.
	hasIgnoredChild := false
	entries, err := gfile.ScanDir(dir, "*", false)
	if err != nil {
		mlog.Printf("scan directory '%s' error: %s", dir, err.Error())
		// If we can't scan the directory, add it to watch list as fallback
		*watchPaths = append(*watchPaths, watchPath{Path: dir, Recursive: true})
		return
	}

	// Check for ignored directories in immediate children
	for _, entry := range entries {
		mlog.Printf("entry: %s", entry)
		if !gfile.IsDir(entry) {
			continue
		}
		if isIgnoredDirName(entry, ignorePatterns) {
			hasIgnoredChild = true
			break
		}
	}

	if !hasIgnoredChild {
		// No ignored descendants, watch this directory (recursive watch covers all).
		*watchPaths = append(*watchPaths, watchPath{Path: dir, Recursive: true})
	} else {
		// Has ignored immediate children, watch current directory non-recursively to catch top-level files,
		// and recurse into valid subdirectories recursively.
		*watchPaths = append(*watchPaths, watchPath{Path: dir, Recursive: false})
		for _, entry := range entries {
			if !gfile.IsDir(entry) {
				continue
			}
			if !isIgnoredDirName(entry, ignorePatterns) {
				app.collectWatchPaths(entry, ignorePatterns, watchPaths)
			}
		}
	}
}

// defaultIgnorePatterns contains glob patterns for directory names that should be ignored when watching.
// These directories typically contain third-party code or non-source files.
// Patterns support glob syntax: * matches any sequence of characters, ? matches single character.
var defaultIgnorePatterns = []string{
	"node_modules",
	"vendor",
	".*", // All hidden directories (covers .git, .svn, .hg, .idea, .vscode, etc.)
	"_*", // Directories starting with underscore
}

// isIgnoredDirName checks if a directory name matches any ignored pattern.
// It accepts either a full path or just the directory name.
func isIgnoredDirName(name string, ignorePatterns []string) bool {
	baseName := gfile.Basename(name)
	for _, pattern := range ignorePatterns {
		if matched, _ := filepath.Match(pattern, baseName); matched {
			return true
		}
	}
	return false
}
