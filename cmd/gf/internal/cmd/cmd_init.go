// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/allyes"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
)

var (
	// Init .
	Init = cInit{}
)

type cInit struct {
	g.Meta `name:"init" brief:"{cInitBrief}" eg:"{cInitEg}"`
}

const (
	cInitRepoPrefix  = `github.com/gogf/`
	cInitMonoRepo    = `template-mono`
	cInitMonoRepoApp = `template-mono-app`
	cInitSingleRepo  = `template-single`
	cInitBrief       = `create and initialize an empty GoFrame project`
	cInitEg          = `
gf init my-project
gf init my-mono-repo -m
gf init my-mono-repo -a
`
	cInitNameBrief = `
name for the project. It will create a folder with NAME in current directory.
The NAME will also be the module name for the project.
`
	// cInitGitDir the git directory
	cInitGitDir = ".git"
	// cInitGitignore the gitignore file
	cInitGitignore = ".gitignore"
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cInitBrief`:     cInitBrief,
		`cInitEg`:        cInitEg,
		`cInitNameBrief`: cInitNameBrief,
	})
}

type cInitInput struct {
	g.Meta  `name:"init"`
	Name    string `name:"NAME" arg:"true" v:"required" brief:"{cInitNameBrief}"`
	Mono    bool   `name:"mono" short:"m" brief:"initialize a mono-repo instead a single-repo" orphan:"true"`
	MonoApp bool   `name:"monoApp" short:"a" brief:"initialize a mono-repo-app instead a single-repo" orphan:"true"`
	Update  bool   `name:"update" short:"u" brief:"update to the latest goframe version" orphan:"true"`
	Module  string `name:"module" short:"g" brief:"custom go module"`
}

type cInitOutput struct{}

func (c cInit) Index(ctx context.Context, in cInitInput) (out *cInitOutput, err error) {
	var overwrote = false
	if !gfile.IsEmpty(in.Name) && !allyes.Check() {
		s := gcmd.Scanf(`the folder "%s" is not empty, files might be overwrote, continue? [y/n]: `, in.Name)
		if strings.EqualFold(s, "n") {
			return
		}
		overwrote = true
	}
	mlog.Print("initializing...")

	// Create project folder and files.
	var (
		templateRepoName string
		gitignoreFile    = in.Name + "/" + cInitGitignore
	)

	if in.Mono {
		templateRepoName = cInitMonoRepo
	} else if in.MonoApp {
		templateRepoName = cInitMonoRepoApp
	} else {
		templateRepoName = cInitSingleRepo
	}

	err = gres.Export(templateRepoName, in.Name, gres.ExportOption{
		RemovePrefix: templateRepoName,
	})
	if err != nil {
		return
	}

	// build ignoreFiles from the .gitignore file
	ignoreFiles := make([]string, 0, 10)
	ignoreFiles = append(ignoreFiles, cInitGitDir)
	// in.MonoApp is a mono-repo-app, it should ignore the .gitignore file
	if overwrote && !in.MonoApp {
		err = gfile.ReadLines(gitignoreFile, func(line string) error {
			// Add only hidden files or directories
			// If other directories are added, it may cause the entire directory to be ignored
			// such as 'main' in the .gitignore file, but the path is ' D:\main\my-project '
			if line != "" && strings.HasPrefix(line, ".") {
				ignoreFiles = append(ignoreFiles, line)
			}
			return nil
		})

		// if not found the .gitignore file will skip os.ErrNotExist error
		if err != nil && !os.IsNotExist(err) {
			return
		}
	}

	// Get template name and module name.
	if in.Module == "" {
		in.Module = gfile.Basename(gfile.RealPath(in.Name))
	}
	if in.MonoApp {
		pwd := gfile.Pwd() + string(os.PathSeparator) + in.Name
		in.Module = utils.GetImportPath(pwd)
	}

	// Replace template name to project name.
	err = gfile.ReplaceDirFunc(func(path, content string) string {
		for _, ignoreFile := range ignoreFiles {
			if strings.Contains(path, ignoreFile) {
				return content
			}
		}
		return gstr.Replace(gfile.GetContents(path), cInitRepoPrefix+templateRepoName, in.Module)
	}, in.Name, "*", true)
	if err != nil {
		return
	}

	// Update the GoFrame version.
	if in.Update {
		mlog.Print("update goframe...")
		// go get -u github.com/gogf/gf/v2@latest
		updateCommand := `go get -u github.com/gogf/gf/v2@latest`
		if in.Name != "." {
			updateCommand = fmt.Sprintf(`cd %s && %s`, in.Name, updateCommand)
		}
		if err = gproc.ShellRun(ctx, updateCommand); err != nil {
			mlog.Fatal(err)
		}
		// go mod tidy
		gomModTidyCommand := `go mod tidy`
		if in.Name != "." {
			gomModTidyCommand = fmt.Sprintf(`cd %s && %s`, in.Name, gomModTidyCommand)
		}
		if err = gproc.ShellRun(ctx, gomModTidyCommand); err != nil {
			mlog.Fatal(err)
		}
	}

	mlog.Print("initialization done! ")
	if !in.Mono {
		enjoyCommand := `gf run main.go`
		if in.Name != "." {
			enjoyCommand = fmt.Sprintf(`cd %s && %s`, in.Name, enjoyCommand)
		}
		mlog.Printf(`you can now run "%s" to start your journey, enjoy!`, enjoyCommand)
	}
	return
}
