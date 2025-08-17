// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"

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

	resp, err := http.Get(url)
	if err != nil {
		return gerror.Wrapf(err, `download "%s" to "%s" failed`, url, localSaveFilePath)
	}
	defer resp.Body.Close()

	bar := progressbar.NewOptions(int(resp.ContentLength), progressbar.OptionShowBytes(true), progressbar.OptionShowCount())
	writer := io.MultiWriter(out, bar)
	_, err = io.Copy(writer, resp.Body)

	elapsed := time.Since(start)
	if elapsed > time.Minute {
		mlog.Printf(`download completed in %.0fm`, float64(elapsed)/float64(time.Minute))
	} else {
		mlog.Printf(`download completed in %.0fs`, elapsed.Seconds())
	}

	return nil
}
