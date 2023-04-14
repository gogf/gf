package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/container/gset"
	"runtime"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
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
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cUpEg`: cUpEg,
	})
}

type cUpInput struct {
	g.Meta `name:"up"  config:"gfcli.up"`
	All    bool `name:"all" short:"a" brief:"upgrade both version and cli, auto fix codes" orphan:"true"`
	Cli    bool `name:"cli" short:"c" brief:"also upgrade CLI tool" orphan:"true"`
	Fix    bool `name:"fix" short:"f" brief:"auto fix codes(it only make sense if cli is to be upgraded)" orphan:"true"`
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
		if err = c.doUpgradeCLI(ctx); err != nil {
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
				// go get -u
				command := fmt.Sprintf(`cd %s && go get -u %s@latest`, dirPath, pkg.Name)
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
func (c cUp) doUpgradeCLI(ctx context.Context) (err error) {
	mlog.Print(`start upgrading cli...`)
	var (
		downloadUrl = fmt.Sprintf(
			`https://github.com/gogf/gf/releases/latest/download/gf_%s_%s`,
			runtime.GOOS, runtime.GOARCH,
		)
		localSaveFilePath = gfile.SelfPath() + "~"
	)
	mlog.Printf(`start downloading "%s" to "%s", it may take some time`, downloadUrl, localSaveFilePath)
	err = utils.HTTPDownloadFileWithPercent(downloadUrl, localSaveFilePath)
	if err != nil {
		return err
	}

	defer func() {
		mlog.Printf(`new version cli binary is successfully installed to "%s"`, gfile.SelfPath())
		mlog.Printf(`remove temporary buffer file "%s"`, localSaveFilePath)
		_ = gfile.Remove(localSaveFilePath)
	}()

	// It fails if file not exist or its size is less than 1MB.
	if !gfile.Exists(localSaveFilePath) || gfile.Size(localSaveFilePath) < 1024*1024 {
		mlog.Fatalf(`download "%s" to "%s" failed`, downloadUrl, localSaveFilePath)
	}

	// It replaces self binary with new version cli binary.
	switch runtime.GOOS {
	case "windows":
		if err := gfile.Rename(localSaveFilePath, gfile.SelfPath()); err != nil {
			mlog.Fatalf(`install failed: %s`, err.Error())
		}

	default:
		if err := gfile.PutBytes(gfile.SelfPath(), gfile.GetBytes(localSaveFilePath)); err != nil {
			mlog.Fatalf(`install failed: %s`, err.Error())
		}
	}
	return
}

func (c cUp) doAutoFixing(ctx context.Context, dirPath string, version string) (err error) {
	mlog.Printf(`auto fixing directory path "%s" from version "%s" ...`, dirPath, version)
	command := fmt.Sprintf(`gf fix -p %s`, dirPath)
	_ = gproc.ShellRun(ctx, command)
	return
}
