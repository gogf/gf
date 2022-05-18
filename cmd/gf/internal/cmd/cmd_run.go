package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gfsnotify"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
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
gf run main.go -ede "(test)|(vendor)"
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
		ExcludeDirExpr string `name:"excludeDirExpr" short:"ede" brief:"{cRunExcludeDirExprBrief}"`
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

	var reg *regexp.Regexp
	excludeDirExpr := in.ExcludeDirExpr
	if len(excludeDirExpr) > 0 {
		reg, err = regexp.Compile(excludeDirExpr)
		if err != nil {
			mlog.Printf("excludeDirExpr(%s) err: %v", excludeDirExpr, err)
		}
	}

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

	for _, subPath := range fileAllDirs(gfile.RealPath(".")) {
		if reg != nil && reg.MatchString(subPath) {
			mlog.Printf("watcher exclude dir match: %s", subPath)
			continue
		}
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
	mlog.Printf("build: %s", app.File)
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
	mlog.Print(buildCommand)
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
	mlog.Print(runCommand)
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

// fileIsDir checks whether given `path` a directory.
func fileIsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// fileAllDirs returns all sub-folders including itself of given `path` recursively.
func fileAllDirs(path string) (list []string) {
	list = []string{path}
	file, err := os.Open(path)
	if err != nil {
		return list
	}
	defer file.Close()
	names, err := file.Readdirnames(-1)
	if err != nil {
		return list
	}
	for _, name := range names {
		tempPath := fmt.Sprintf("%s%s%s", path, string(filepath.Separator), name)
		if fileIsDir(tempPath) {
			if array := fileAllDirs(tempPath); len(array) > 0 {
				list = append(list, array...)
			}
		}
	}
	return
}
