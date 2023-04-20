// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

// HTTPDownloadFileWithPercent downloads target url file to local path with percent process printing.
func HTTPDownloadFileWithPercent(url string, localSaveFilePath string) error {
	start := time.Now()
	out, err := os.Create(localSaveFilePath)
	if err != nil {
		return gerror.Wrapf(err, `download "%s" to "%s" failed`, url, localSaveFilePath)
	}
	defer out.Close()

	headResp, err := http.Head(url)
	if err != nil {
		return gerror.Wrapf(err, `download "%s" to "%s" failed`, url, localSaveFilePath)
	}
	defer headResp.Body.Close()

	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))
	if err != nil {
		return gerror.Wrap(err, "retrieve Content-Length failed")
	}
	doneCh := make(chan int64)

	go doPrintDownloadPercent(doneCh, localSaveFilePath, int64(size))

	resp, err := http.Get(url)
	if err != nil {
		return gerror.Wrapf(err, `download "%s" to "%s" failed`, url, localSaveFilePath)
	}
	defer resp.Body.Close()

	wroteBytesCount, err := io.Copy(out, resp.Body)
	if err != nil {
		return gerror.Wrapf(err, `download "%s" to "%s" failed`, url, localSaveFilePath)
	}

	doneCh <- wroteBytesCount
	elapsed := time.Since(start)
	if elapsed > time.Minute {
		mlog.Printf(`download completed in %.0fm`, float64(elapsed)/float64(time.Minute))
	} else {
		mlog.Printf(`download completed in %.0fs`, elapsed.Seconds())
	}

	return nil
}

func doPrintDownloadPercent(doneCh chan int64, localSaveFilePath string, total int64) {
	var (
		stop           = false
		lastPercentFmt string
	)
	for {
		select {
		case <-doneCh:
			stop = true

		default:
			file, err := os.Open(localSaveFilePath)
			if err != nil {
				mlog.Fatal(err)
			}
			fi, err := file.Stat()
			if err != nil {
				mlog.Fatal(err)
			}
			size := fi.Size()
			if size == 0 {
				size = 1
			}
			var (
				percent    = float64(size) / float64(total) * 100
				percentFmt = fmt.Sprintf(`%.0f`, percent) + "%"
			)
			if lastPercentFmt != percentFmt {
				lastPercentFmt = percentFmt
				mlog.Print(percentFmt)
			}
		}

		if stop {
			break
		}
		time.Sleep(time.Second)
	}
}
