package cmd

import (
	"context"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	Update = cUpdate{}
)

type cUpdate struct {
	g.Meta `name:"update" brief:"update GoFrame CLI version to latest" eg:"{cUpdateEg}" `
}

const (
	cUpdateEg = `
gf update
`
	updateClientUAPrefix = "gf-cli"
	gfGithubReleaseUrl   = "https://github.com/gogf/gf/releases"
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cUpdateEg`: cUpdateEg,
	})
}

type cUpdateInput struct {
	g.Meta `name:"update"  config:"gfcli.update"`
	Proxy  string `name:"proxy" short:"p" brief:"set proxy to access update repo" usage:"use proxy server string,like http://USER:PASSWORD@IP:PORT or socks5://USER:PASSWORD@IP:PORT"`
	Magic  bool   `name:"magic" short:"m" brief:"use ghproxy magic(for CN users) to get update content" orphan:"true" `
	Force  bool   `name:"force" short:"f"  brief:"force update local cli to latest" orphan:"true"`
}

type cUpdateOutput struct{}

func (c cUpdate) Index(ctx context.Context, in cUpdateInput) (out *cUpdateOutput, err error) {
	client := gclient.New()
	client.SetAgent(fmt.Sprintf("%s/%s", updateClientUAPrefix, gf.VERSION))
	client.SetTimeout(10 * time.Second)
	mlog.Print("Checking latest cli version")
	if in.Proxy != "" {
		mlog.Printf("Using proxy %s", in.Proxy)
	}
	version, err := c.getLatestReleaseVersion(ctx, in.Proxy, client)
	if err != nil {
		return nil, err
	}
	if !in.Force && version == gf.VERSION {
		mlog.Printf("Current cli %s is latest version,no need to update.", version)
		return nil, nil
	}
	doUpdate := "n"
	if !in.Force {
		doUpdate = gcmd.Scanf("Latest remote version:%s,local version:%s,update to latest version?[y/n]", version, gf.VERSION)
	} else {
		mlog.Printf("Force update cli.")
	}
	if strings.EqualFold(doUpdate, "n") {
		mlog.Printf("Do nothing with update.")
		return
	}
	mlog.Printf("Latest remote version:%s,local version:%s", version, gf.VERSION)
	var binaryName string
	if runtime.GOOS == "windows" {
		binaryName = fmt.Sprintf(`%s/%s/download/gf_%s_%s.exe`, runtime.GOOS, runtime.GOARCH)
	} else {
		binaryName = fmt.Sprintf(`gf_%s_%s`, runtime.GOOS, runtime.GOARCH)
	}
	downloadUrl := fmt.Sprintf(`%s/download/%s/%s`, gfGithubReleaseUrl, version, binaryName)
	if in.Magic {
		mlog.Print("Using ghproxy 🪄✨🌟")
		downloadUrl = fmt.Sprintf("https://ghproxy.com/%s", downloadUrl)
	}
	mlog.Printf("Start download from:%s", downloadUrl)

	localSaveFilePath := fmt.Sprintf("%s%s-%s", gfile.Temp(), version, binaryName)
	if err := c.doCliUpdate(ctx, in.Proxy, downloadUrl, localSaveFilePath, client); err != nil {
		return nil, err
	}
	if err := gfile.Chmod(localSaveFilePath, os.FileMode(0777)); err != nil {
		return nil, err
	}
	mlog.Printf("start new installer:%s", localSaveFilePath)
	return nil, gproc.ShellRun(ctx, fmt.Sprintf("%s install", localSaveFilePath))
}

func (c cUpdate) doCliUpdate(ctx context.Context, proxy, srcUrl, localSaveFilePath string, client *gclient.Client) error {
	if proxy != "" {
		mlog.Debugf("use proxy %s download s% to %s", proxy, srcUrl, localSaveFilePath)
	}
	var (
		source     io.Reader
		sourceSize int64
	)
	outputFile, err := os.Create(localSaveFilePath)
	if err != nil {
		return gerror.Wrapf(err, `create localSavedFile error`)
	}
	defer outputFile.Close()
	resp, err := client.Proxy(proxy).Get(ctx, srcUrl)
	if err != nil {
		return gerror.Wrapf(err, "get remote file error")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return gerror.Wrapf(err, "Server return non-200 status: %v", resp.Status)
	}
	i, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	sourceSize = int64(i)
	if sourceSize < 100 {
		return gerror.New("Content-Length get error")
	}
	source = resp.Body
	// start new bar
	bar := pb.Full.Start64(sourceSize)
	// create proxy reader
	barReader := bar.NewProxyReader(source)
	io.Copy(outputFile, barReader)
	// finish bar
	bar.Finish()
	return nil
}

func (c cUpdate) getLatestReleaseVersion(ctx context.Context, proxy string, client *gclient.Client) (string, error) {
	resp, err := client.Proxy(proxy).RedirectLimit(0).Get(ctx, gfGithubReleaseUrl+"/latest")
	if err != nil {
		return "", err
	}
	defer resp.Close()
	redirectLocation := ""
	if resp.StatusCode == 302 {
		redirectLocation = resp.Response.Header.Get("location")
	}
	splitLocationLocation := gstr.Split(redirectLocation, "/")
	if redirectLocation == "" || len(splitLocationLocation) == 0 {
		return "", fmt.Errorf("get latest version error")
	}
	return splitLocationLocation[len(splitLocationLocation)-1], nil
}
