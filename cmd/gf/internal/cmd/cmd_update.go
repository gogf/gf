package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

const (
	releaseUrl = "https://api.github.com/repos/gogf/gf/releases"
	binaryUrl  = "https://github.com/gogf/gf/releases/download/%s/%s"
	projectUrl = "git@github.com:gogf/gf.git"
	codeUrl    = "https://github.com/gogf/gf/archive/refs/tags/%s"
)

var (
	Update = cUpdate{}
)

type cUpdate struct {
	g.Meta  `name:"update" brief:"update gf-cli to last stable tag"`
	tempDir string
	release gitRelease
}

type cUpdateInput struct {
	g.Meta `name:"update"`
	Tag    string `name:"tag" short:"t" brief:"select a tag, default is latest"`
	Way    string `name:"way" short:"w" brief:"download way, You can choose wget、clone or code, default is wget binary file"`
}
type cUpdateOutput struct{}

func getGitFileName() string {
	return fmt.Sprintf("gf_%s_%s", runtime.GOOS, runtime.GOARCH)
}

type gitRelease struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
}

func getRelease(tag string) (release gitRelease, err error) {
	var releases []gitRelease
	if err = g.Client().GetVar(context.Background(), releaseUrl).Scan(&releases); err != nil {
		return
	}
	if len(releases) == 0 {
		mlog.Fatal("There is no tags！")
		return
	}
	release = releases[0]
	if len(tag) > 0 {
		for _, r := range releases {
			if r.TagName == tag {
				release = r
			}
		}
	}
	return
}

func workDir() (workPath string) {
	home, err := os.UserHomeDir()
	if err != nil {
		mlog.Fatalf("create the tmp file err: %+v", err)
	}
	workPath = path.Join(home, ".gf")
	if _, err = os.Stat(workPath); os.IsNotExist(err) {
		if err = os.MkdirAll(workPath, 0o700); err != nil {
			mlog.Fatalf("create the workPath err: %+v", err)
		}
	}
	return
}

func (c cUpdate) Index(ctx context.Context, in cUpdateInput) (*cUpdateOutput, error) {
	c.tempDir = workDir()
	defer os.RemoveAll(c.tempDir)
	release, err := getRelease(in.Tag)
	if err != nil {
		return nil, err
	}
	c.release = release
	if in.Way == "clone" {
		err = c.updateByClone()
	} else if in.Way == "code" {
		err = c.updateByCode()
	} else {
		err = c.updateByWget()
	}
	return nil, err
}

func (c cUpdate) updateByWget() (err error) {
	filename := getGitFileName()
	url := fmt.Sprintf(binaryUrl, c.release.TagName, filename)
	if err = cmdDo(c.tempDir, "wget", url); err != nil {
		mlog.Fatalf(`wget project err: %+v`, err)
	}
	if err = gfile.Chmod(path.Join(c.tempDir, filename), os.ModePerm); err != nil {
		mlog.Fatalf("change permission err: %+v", err)
	}
	if err = cmdDo(c.tempDir, "./"+filename, "install"); err != nil {
		mlog.Fatalf(`install project err: %+v`, err)
	}
	return nil
}

func (c cUpdate) updateByClone() (err error) {
	if err = cmdDo(c.tempDir, "git", "clone", projectUrl); err != nil {
		mlog.Fatalf(`clone project err: %+v`, err)
	}
	if err = cmdDo(path.Join(c.tempDir, "gf"), "git", "fetch"); err != nil {
		mlog.Fatalf(`fetch project err: %+v`, err)
	}
	if err = cmdDo(path.Join(c.tempDir, "gf"), "git", "checkout", c.release.TargetCommitish); err != nil {
		mlog.Fatalf(`checkout project err: %+v`, err)
	}
	if err = cmdDo(path.Join(c.tempDir, "gf/cmd/gf"), "go", "install"); err != nil {
		mlog.Fatalf(`install err: %+v`, err)
	}
	return nil
}

func (c cUpdate) updateByCode() (err error) {
	fileName := fmt.Sprintf("%s.zip", c.release.TagName)
	url := fmt.Sprintf(codeUrl, fileName)
	installPath := fmt.Sprintf("%s/gf-%s/cmd/gf", c.tempDir, strings.TrimLeft(c.release.TagName, "v"))
	if err = cmdDo(c.tempDir, "wget", url); err != nil {
		mlog.Fatalf("download file err: %+v", err)
	}
	if err = gcompress.UnZipFile(path.Join(c.tempDir, fileName), path.Join(c.tempDir)); err != nil {
		mlog.Fatalf("UnGzipFile err: %+v", err)
	}
	if err = cmdDo(installPath, "go", "install"); err != nil {
		mlog.Fatalf("install err: %+v", err)
	}
	return nil
}

func cmdDo(dir, name string, params ...string) error {
	if !strings.HasPrefix(name, ".") {
		bin := gproc.SearchBinary(name)
		if len(bin) == 0 {
			mlog.Fatalf(fmt.Sprintf(`command "%s" not found in your environment, please install %s first to proceed this command`, name, name))
		}
	}
	cmd := exec.Command(name, params...)
	if len(dir) > 0 {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
