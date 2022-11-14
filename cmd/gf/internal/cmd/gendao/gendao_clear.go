package gendao

import (
	"context"

	"github.com/gogf/gf/v2/os/gfile"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

func doClear(ctx context.Context, dirPath string) {
	files, err := gfile.ScanDirFile(dirPath, "*.go", true)
	if err != nil {
		mlog.Fatal(err)
	}
	for _, file := range files {
		if utils.IsFileDoNotEdit(file) {
			if err = gfile.Remove(file); err != nil {
				mlog.Print(err)
			}
		}
	}
}
