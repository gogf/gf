package update

import (
	"fmt"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/tool/gf/library/mlog"
	"runtime"
)

var (
	cdnUrl  = g.Config("url").GetString("cdn.url")
	homeUrl = g.Config("url").GetString("home.url")
)

func init() {
	if cdnUrl == "" {
		mlog.Fatal("CDN configuration cannot be empty")
	}
	if homeUrl == "" {
		mlog.Fatal("Home configuration cannot be empty")
	}
}

func Run() {
	mlog.Print("checking...")
	md5Url := homeUrl + `/cli/binary/md5`
	latestMd5 := ghttp.GetContent(md5Url, g.Map{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	})
	if latestMd5 == "" {
		mlog.Fatal("get the latest binary md5 failed, may be network issue")
	}
	localMd5, err := gmd5.EncryptFile(gfile.SelfPath())
	if err != nil {
		mlog.Fatal("calculate local binary md5 failed,", err.Error())
	}
	if localMd5 != latestMd5 {
		mlog.Print("downloading...")
		ext := ""
		if runtime.GOOS == "windows" {
			ext = ".exe"
		}
		downloadUrl := fmt.Sprintf(
			`%s/cli/binary/%s_%s/gf%s?%s`,
			cdnUrl,
			runtime.GOOS,
			runtime.GOARCH,
			ext,
			latestMd5,
		)
		mlog.Debugf("HTTP GET %s", downloadUrl)
		res, err := ghttp.Get(downloadUrl)
		if err != nil || res.StatusCode != 200 {
			mlog.Fatalf(
				"downloading failed for %s %s, may be network issue:\n%s",
				runtime.GOOS, runtime.GOARCH, res.ReadAllString(),
			)
		}
		defer res.Close()
		data := res.ReadAll()
		mlog.Print("installing...")
		var (
			binPath    = gfile.SelfPath()
			binDirPath = gfile.SelfDir()
			renamePath = binPath + "~"
		)
		if runtime.GOOS == "windows" {
			// Rename myself for windows.
			if err := gfile.Rename(binPath, renamePath); err != nil {
				mlog.Fatal("rename binary file failed:", err.Error())
			}
			defer gfile.Remove(renamePath)
		} else {
			// Remove the binary for other platforms.
			if gfile.IsWritable(binDirPath) {
				if err := gfile.Remove(binPath); err != nil {
					mlog.Fatal("remove binary failed:", err.Error())
				}
			}
		}
		if err := gfile.PutBytes(binPath, data); err != nil {
			mlog.Fatal("install binary failed:", err.Error())
		}
		if err := gfile.Chmod(binPath, 0777); err != nil {
			mlog.Fatal("chmod binary failed:", err.Error())
		}
		mlog.Print("gf binary is now updated to the latest version")
	} else {
		mlog.Print("it's the latest version, no need updates")
	}
}
