package cmd

import (
	"context"
	"fmt"
	"runtime"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gfsnotify"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
)

var (
	Run = cRun{}
)

type cRun struct {
	g.Meta `name:"run" usage:"{cRunUsage}" brief:"{cRunBrief}" eg:"{cRunEg}" dc:"{cRunDc}"`
}

type cRunApp struct {
	File    string // Go run file name.
	Path    string // Directory storing built binary.
	Options string // Extra "go run" options.
	Args    string // Custom arguments.
}

const (
	cRunUsage = `gf run FILE [OPTION]`
	cRunBrief = `running go codes with hot-compiled-like feature`
	cRunEg    = `
gf run main.go
gf run main.go --args "server -p 8080"
gf run main.go -mod=vendor
gf run main.go -d "(test)|(vendor)"
`
	cRunDc = `
The "run" command is used for running go codes with hot-compiled-like feature,
which compiles and runs the go codes asynchronously when codes change.
`
	cRunFileBrief           = `building file path.`
	cRunPathBrief           = `output directory path for built binary file. it's "manifest/output" in default`
	cRunExtraBrief          = `the same options as "go run"/"go build" except some options as follows defined`
	cRunExcludeDirExprBrief = `exclude directory expression, which is used for excluding some directories from watching.`
)

var (
	process *gproc.Process
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cRunUsage`:               cRunUsage,
		`cRunBrief`:               cRunBrief,
		`cRunEg`:                  cRunEg,
		`cRunDc`:                  cRunDc,
		`cRunFileBrief`:           cRunFileBrief,
		`cRunPathBrief`:           cRunPathBrief,
		`cRunExtraBrief`:          cRunExtraBrief,
		`cRunExcludeDirExprBrief`: cRunExcludeDirExprBrief,
	})
}

type (
	cRunInput struct {
		g.Meta         `name:"run"`
		File           string `name:"FILE"  arg:"true" brief:"{cRunFileBrief}" v:"required"`
		Path           string `name:"path"  short:"p"  brief:"{cRunPathBrief}" d:"./"`
		Extra          string `name:"extra" short:"e"  brief:"{cRunExtraBrief}"`
		ExcludeDirExpr string `name:"excludeDirExpr" short:"d" brief:"{cRunExcludeDirExprBrief}"`
	}
	cRunOutput struct{}
)

func (c cRun) Index(ctx context.Context, in cRunInput) (out *cRunOutput, err error) {
	// Necessary check.
	if gproc.SearchBinary("go") == "" {
		mlog.Fatalf(`command "go" not found in your environment, please install golang first to proceed this command`)
	}

	app := &cRunApp{
		File:    in.File,
		Path:    in.Path,
		Options: in.Extra,
	}
	dirty := gtype.NewBool()

	currentPath := gfile.RealPath(".")

	// exclude dir
	hasReg := false
	excludeDirExpr := in.ExcludeDirExpr
	excludeDirCount := 0
	if len(excludeDirExpr) > 0 {
		err := gregex.Validate(excludeDirExpr)
		if err != nil {
			mlog.Printf("exclude directory expression[%s] err: %v", excludeDirExpr, err)
		} else {
			mlog.Printf("exclude directory expression[%s]", excludeDirExpr)
			hasReg = true
		}
	}

	listDir, err := gfile.ScanDirFunc(currentPath, "*", true, func(path string) string {
		if gfile.IsDir(path) {
			if hasReg && gregex.IsMatchString(excludeDirExpr, path) {
				excludeDirCount++
				mlog.Debugf("exclude directory path: %s", path)
				return ""
			}
			return path
		}
		return ""
	})

	if hasReg {
		mlog.Printf("exclude directory count: %d \n", excludeDirCount)
	}

	listDir = append(listDir, currentPath)

	callback := func(event *gfsnotify.Event) {
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
			mlog.Printf(`go file changes: %s`, event.String())
			app.Run()
		})
	}

	for _, subPath := range listDir {
		mlog.Debugf("watch directory: %v", subPath)
		_, err = gfsnotify.Add(subPath, callback, false)
		if err != nil {
			mlog.Fatal(err)
		}
	}

	go app.Run()
	select {}
}

func (app *cRunApp) Run() {
	// Rebuild and run the codes.
	renamePath := ""
	mlog.Printf("build file: %s", app.File)
	outputPath := gfile.Join(app.Path, gfile.Name(app.File))
	if runtime.GOOS == "windows" {
		outputPath += ".exe"
		if gfile.Exists(outputPath) {
			renamePath = outputPath + "~"
			if err := gfile.Rename(outputPath, renamePath); err != nil {
				mlog.Print(err)
			}
		}
	}
	// In case of `pipe: too many open files` error.
	// Build the app.
	buildCommand := fmt.Sprintf(
		`go build -o %s %s %s`,
		outputPath,
		app.Options,
		app.File,
	)
	mlog.Printf("buildCommand: %s", buildCommand)
	result, err := gproc.ShellExec(buildCommand)
	if err != nil {
		mlog.Printf("build error: \n%s%s", result, err.Error())
		return
	}
	// Kill the old process if build successfully.
	if process != nil {
		if err := process.Kill(); err != nil {
			mlog.Debugf("kill process error: %s", err.Error())
			//return
		}
	}
	// Run the binary file.
	runCommand := fmt.Sprintf(`%s %s`, outputPath, app.Args)
	mlog.Printf("runCommand: %s", runCommand)
	if runtime.GOOS == "windows" {
		// Special handling for windows platform.
		// DO NOT USE "cmd /c" command.
		process = gproc.NewProcess(outputPath, gstr.SplitAndTrim(" ", app.Args))
	} else {
		process = gproc.NewProcessCmd(outputPath, gstr.SplitAndTrim(" ", app.Args))
	}
	if pid, err := process.Start(); err != nil {
		mlog.Printf("build running error: %s", err.Error())
	} else {
		mlog.Printf("build running pid: %d", pid)
	}
}
