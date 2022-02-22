package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/frame/g"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

var (
	Update     = cUpdate{}
	releaseUrl = "https://api.github.com/repos/gogf/gf/releases"
	binaryUrl  = "https://github.com/gogf/gf/releases/download/%s/%s"
	projectUrl = "git@github.com:gogf/gf.git"
	codeUrl    = "https://github.com/gogf/gf/archive/refs/tags/%s"
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

type cmdInfo struct {
	dir     string
	command []string
}

func getRelease(tag string) (release gitRelease, err error) {
	resp, err := http.Get(releaseUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyContent, _ := ioutil.ReadAll(resp.Body)
	var releases []gitRelease
	if err = json.Unmarshal(bodyContent, &releases); err != nil {
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

func cmdDo(runCmd []cmdInfo) error {
	for _, data := range runCmd {
		cmd := exec.Command(data.command[0], data.command[1:]...)
		cmd.Dir = data.dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (c cUpdate) Index(ctx context.Context, in cUpdateInput) (*cUpdateOutput, error) {
	c.tempDir = workDir()
	release, err := getRelease(in.Tag)
	if err != nil {
		return nil, err
	}
	c.release = release
	defer os.RemoveAll(c.tempDir)
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
	runCmd := []cmdInfo{
		{c.tempDir, []string{"wget", url}},
		{c.tempDir, []string{"chmod", "+x", filename}},
		{c.tempDir, []string{"./" + filename, "install"}},
	}
	if err = cmdDo(runCmd); err != nil {
		mlog.Fatalf("do command err: %+v", err)
	}
	return nil
}

func (c cUpdate) updateByClone() (err error) {
	runCmd := []cmdInfo{
		{c.tempDir, []string{"git", "clone", projectUrl}},
		{c.tempDir + "/gf", []string{"git", "fetch"}},
		{c.tempDir + "/gf", []string{"git", "checkout", c.release.TargetCommitish}},
		{c.tempDir + "/gf/cmd/gf", []string{"go", "install"}},
	}
	if err = cmdDo(runCmd); err != nil {
		mlog.Fatalf("do command err: %+v", err)
	}
	return nil
}

func (c cUpdate) updateByCode() (err error) {
	fileName := fmt.Sprintf("%s.tar.gz", c.release.TagName)
	url := fmt.Sprintf(codeUrl, fileName)
	installPath := fmt.Sprintf("%s/gf-%s/cmd/gf", c.tempDir, strings.TrimLeft(c.release.TagName, "v"))
	runCmd := []cmdInfo{
		{c.tempDir, []string{"wget", url}},
		{c.tempDir, []string{"tar", "zxf", fileName}},
		{installPath, []string{"go", "install"}},
	}
	if err = cmdDo(runCmd); err != nil {
		mlog.Fatalf("do command err: %+v", err)
	}
	return nil
}
