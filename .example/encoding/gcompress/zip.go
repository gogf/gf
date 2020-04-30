package main

import (
	"archive/zip"
	"fmt"
	"github.com/gogf/gf/encoding/gcompress"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// srcFile could be a single file or a directory
func Zip(srcFile string, destZip string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
		// header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

func main() {
	src := `/Users/john/Workspace/Go/GOPATH/src/github.com/gogf/gf/test`
	dst := `/Users/john/Workspace/Go/GOPATH/src/github.com/gogf/gf/test.zip`
	//src := `/Users/john/Workspace/Go/GOPATH/src/github.com/gogf/gf/README.MD`
	//dst := `/Users/john/Workspace/Go/GOPATH/src/github.com/gogf/gf/README.MD.zip`
	fmt.Println(gcompress.ZipPath(src, dst))
	//fmt.Println(Zip(src, dst))
}
