// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"context"
	"fmt"
	"runtime"

	"github.com/gogf/selfupdate"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

var (
	Up = cUp{}
)

type cUp struct {
	g.Meta `name:"up" brief:"upgrade GoFrame version/tool to latest one in current project" eg:"{cUpEg}" `
}

const (
	gfPackage = `github.com/gogf/gf/`
	cUpEg     = `
gf up
gf up -a
gf up -c
gf up -cf
gf up -a -m=install
gf up -a -m=install -p=github.com/gogf/gf/cmd/gf/v2@latest
`
	cliMethodHttpDownload = "http"
	cliMethodGoInstall    = "install"
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cUpEg`: cUpEg,
	})
}

type cUpInput struct {
	g.Meta               `name:"up" config:"gfcli.up"`
	All                  bool   `name:"all" short:"a" brief:"upgrade both version and cli, auto fix codes" orphan:"true"`
	Cli                  bool   `name:"cli" short:"c" brief:"also upgrade CLI tool" orphan:"true"`
	Fix                  bool   `name:"fix" short:"f" brief:"auto fix codes(it only make sense if cli is to be upgraded)" orphan:"true"`
	CliDownloadingMethod string `name:"cli-download-method" short:"m" brief:"cli upgrade method: http=download binary via HTTP GET, install=upgrade via go install" d:"http"`
	// CliModulePath specifies the module path for CLI installation via go install.
	// This is used when CliDownloadingMethod is set to "install".
	CliModulePath string `name:"cli-module-path" short:"p" brief:"custom cli module path for upgrade CLI tool with go install method" d:"github.com/gogf/gf/cmd/gf/v2@latest"`
}

type cUpOutput struct{}

func (c cUp) Index(ctx context.Context, in cUpInput) (out *cUpOutput, err error) {
	defer func() {
		if err == nil {
			mlog.Print()
			mlog.Print(`👏congratulations! you've upgraded to the latest version of GoFrame! enjoy it!👏`)
			mlog.Print()
		}
	}()

	var doUpgradeVersionOut *doUpgradeVersionOutput
	if in.All {
		in.Cli = true
		in.Fix = true
	}
	if doUpgradeVersionOut, err = c.doUpgradeVersion(ctx, in); err != nil {
		return nil, err
	}

	if in.Cli {
		if err = c.doUpgradeCLI(ctx, in); err != nil {
			return nil, err
		}
	}

	if in.Cli && in.Fix {
		if doUpgradeVersionOut != nil && len(doUpgradeVersionOut.Items) > 0 {
			upgradedPathSet := gset.NewStrSet()
			for _, item := range doUpgradeVersionOut.Items {
				if !upgradedPathSet.AddIfNotExist(item.DirPath) {
					continue
				}
				if err = c.doAutoFixing(ctx, item.DirPath, item.Version); err != nil {
					return nil, err
				}
			}
		}
	}
	return
}

type doUpgradeVersionOutput struct {
	Items []doUpgradeVersionOutputItem
}

type doUpgradeVersionOutputItem struct {
	DirPath string
	Version string
}

func (c cUp) doUpgradeVersion(ctx context.Context, in cUpInput) (out *doUpgradeVersionOutput, err error) {
	mlog.Print(`start upgrading version...`)
	out = &doUpgradeVersionOutput{
		Items: make([]doUpgradeVersionOutputItem, 0),
	}
	type Package struct {
		Name    string
		Version string
	}

	var (
		temp      string
		dirPath   = gfile.Pwd()
		goModPath = gfile.Join(dirPath, "go.mod")
	)
	// It recursively upgrades the go.mod from sub folder to its parent folders.
	for {
		if gfile.Exists(goModPath) {
			var packages []Package
			err = gfile.ReadLines(goModPath, func(line string) error {
				line = gstr.Trim(line)
				line = gstr.TrimLeftStr(line, "require ")
				line = gstr.Trim(line)
				if gstr.HasPrefix(line, gfPackage) {
					array := gstr.SplitAndTrim(line, " ")
					packages = append(packages, Package{
						Name:    array[0],
						Version: array[1],
					})
				}
				return nil
			})
			if err != nil {
				return
			}
			for _, pkg := range packages {
				mlog.Printf(`upgrading "%s" from "%s" to "latest"`, pkg.Name, pkg.Version)
				mlog.Printf(`running command: go get %s@latest`, pkg.Name)
				// go get @latest
				command := fmt.Sprintf(`cd %s && go get %s@latest`, dirPath, pkg.Name)
				if err = gproc.ShellRun(ctx, command); err != nil {
					return
				}
				// go mod tidy
				if err = utils.GoModTidy(ctx, dirPath); err != nil {
					return nil, err
				}
				out.Items = append(out.Items, doUpgradeVersionOutputItem{
					DirPath: dirPath,
					Version: pkg.Version,
				})
			}
			return
		}
		temp = gfile.Dir(dirPath)
		if temp == "" || temp == dirPath {
			return
		}
		dirPath = temp
		goModPath = gfile.Join(dirPath, "go.mod")
	}
}

// doUpgradeCLI downloads the new version binary with process.
func (c cUp) doUpgradeCLI(ctx context.Context, in cUpInput) (err error) {
	mlog.Print(`start upgrading cli...`)
	fmt.Println(` cli upgrade method:`, in.CliDownloadingMethod)
	switch in.CliDownloadingMethod {
	case cliMethodHttpDownload:
		return c.doUpgradeCLIWithHttpDownload(ctx)
	case cliMethodGoInstall:
		return c.doUpgradeCLIWithGoInstall(ctx, in)
	default:
		mlog.Fatalf(`invalid cli upgrade method: "%s", please use "http" or "install"`, in.CliDownloadingMethod)
	}
	return
}

func (c cUp) doUpgradeCLIWithHttpDownload(ctx context.Context) (err error) {
	mlog.Print(`start upgrading cli with http get download...`)
	var (
		downloadUrl = fmt.Sprintf(
			`https://github.com/gogf/gf/releases/latest/download/gf_%s_%s`,
			runtime.GOOS, runtime.GOARCH,
		)
		localSaveFilePath = gfile.SelfPath() + "~"
	)

	if runtime.GOOS == "windows" {
		downloadUrl += ".exe"
	}

	mlog.Printf(`start downloading "%s" to "%s", it may take some time`, downloadUrl, localSaveFilePath)
	err = utils.HTTPDownloadFileWithPercent(downloadUrl, localSaveFilePath)
	if err != nil {
		return err
	}

	defer func() {
		mlog.Printf(`new version cli binary is successfully installed to "%s"`, gfile.SelfPath())
		mlog.Printf(`remove temporary buffer file "%s"`, localSaveFilePath)
		_ = gfile.RemoveFile(localSaveFilePath)
	}()

	// It fails if file not exist or its size is less than 1MB.
	if !gfile.Exists(localSaveFilePath) || gfile.Size(localSaveFilePath) < 1024*1024 {
		mlog.Fatalf(`download "%s" to "%s" failed`, downloadUrl, localSaveFilePath)
	}

	newFile, err := gfile.Open(localSaveFilePath)
	if err != nil {
		return err
	}
	// selfupdate
	err = selfupdate.Apply(newFile, selfupdate.Options{})
	if err != nil {
		return err
	}
	return
}

func (c cUp) doUpgradeCLIWithGoInstall(ctx context.Context, in cUpInput) (err error) {
	mlog.Print(`upgrading cli with go install...`)
	if !genv.Contains("GOPATH") {
		mlog.Fatal(`"GOPATH" environment variable does not exist, please check your go installation`)
	}

	command := fmt.Sprintf(`go install %s`, in.CliModulePath)
	mlog.Printf(`running command: %s`, command)
	err = gproc.ShellRun(ctx, command)
	if err != nil {
		return err
	}

	cliFilePath := gfile.Join(genv.Get("GOPATH").String(), "bin/gf")
	if runtime.GOOS == "windows" {
		cliFilePath += ".exe"
	}

	// It fails if file not exist or its size is less than 1MB.
	if !gfile.Exists(cliFilePath) || gfile.Size(cliFilePath) < 1024*1024 {
		mlog.Fatalf(`go install %s failed, "%s" does not exist or its size is less than 1MB`, in.CliModulePath, cliFilePath)
	}

	newFile, err := gfile.Open(cliFilePath)
	if err != nil {
		return err
	}
	// selfupdate
	err = selfupdate.Apply(newFile, selfupdate.Options{})
	if err != nil {
		return err
	}
	return
}

func (c cUp) doAutoFixing(ctx context.Context, dirPath string, version string) (err error) {
	mlog.Printf(`auto fixing directory path "%s" from version "%s" ...`, dirPath, version)
	command := fmt.Sprintf(`gf fix -p %s`, dirPath)
	_ = gproc.ShellRun(ctx, command)
	return
}
