// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/os/gspath"
	"github.com/gogf/gf/v2/text/gstr"
)

// SetPath sets the configuration directory path for file search.
// The parameter `path` can be absolute or relative path,
// but absolute path is strongly recommended.
func (c *AdapterFile) SetPath(path string) (err error) {
	var (
		isDir    = false
		realPath = ""
	)
	if file := gres.Get(path); file != nil {
		realPath = path
		isDir = file.FileInfo().IsDir()
	} else {
		// Absolute path.
		realPath = gfile.RealPath(path)
		if realPath == "" {
			// Relative path.
			c.searchPaths.RLockFunc(func(array []string) {
				for _, v := range array {
					if searchedPath, _ := gspath.Search(v, path); searchedPath != "" {
						realPath = searchedPath
						break
					}
				}
			})
		}
		if realPath != "" {
			isDir = gfile.IsDir(realPath)
		}
	}
	// Path not exist.
	if realPath == "" {
		buffer := bytes.NewBuffer(nil)
		if c.searchPaths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf(`SetPath failed: cannot find directory "%s" in following paths:`, path))
			c.searchPaths.RLockFunc(func(array []string) {
				for k, v := range array {
					buffer.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf(`SetPath failed: path "%s" does not exist`, path))
		}
		return gerror.New(buffer.String())
	}
	// Should be a directory.
	if !isDir {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`SetPath failed: path "%s" should be directory type`,
			path,
		)
	}
	// Repeated path check.
	if c.searchPaths.Search(realPath) != -1 {
		return nil
	}
	c.jsonMap.Clear()
	c.searchPaths.Clear()
	c.searchPaths.Append(realPath)
	intlog.Print(context.TODO(), "SetPath:", realPath)
	return nil
}

// AddPath adds an absolute or relative path to the search paths.
func (c *AdapterFile) AddPath(path string) (err error) {
	var (
		isDir    = false
		realPath = ""
	)
	// It firstly checks the resource manager,
	// and then checks the filesystem for the path.
	if file := gres.Get(path); file != nil {
		realPath = path
		isDir = file.FileInfo().IsDir()
	} else {
		// Absolute path.
		realPath = gfile.RealPath(path)
		if realPath == "" {
			// Relative path.
			c.searchPaths.RLockFunc(func(array []string) {
				for _, v := range array {
					if searchedPath, _ := gspath.Search(v, path); searchedPath != "" {
						realPath = searchedPath
						break
					}
				}
			})
		}
		if realPath != "" {
			isDir = gfile.IsDir(realPath)
		}
	}
	if realPath == "" {
		buffer := bytes.NewBuffer(nil)
		if c.searchPaths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf(`AddPath failed: cannot find directory "%s" in following paths:`, path))
			c.searchPaths.RLockFunc(func(array []string) {
				for k, v := range array {
					buffer.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf(`AddPath failed: path "%s" does not exist`, path))
		}
		return gerror.New(buffer.String())
	}
	if !isDir {
		return gerror.NewCodef(gcode.CodeInvalidParameter, `AddPath failed: path "%s" should be directory type`, path)
	}
	// Repeated path check.
	if c.searchPaths.Search(realPath) != -1 {
		return nil
	}
	c.searchPaths.Append(realPath)
	intlog.Print(context.TODO(), "AddPath:", realPath)
	return nil
}

// GetFilePath returns the absolute configuration file path for the given filename by `file`.
// If `file` is not passed, it returns the configuration file path of the default name.
// It returns an empty `path` string and an error if the given `file` does not exist.
func (c *AdapterFile) GetFilePath(fileName ...string) (path string, err error) {
	var (
		tempPath     string
		usedFileName = c.defaultName
	)
	if len(fileName) > 0 {
		usedFileName = fileName[0]
	}
	// Searching resource manager.
	if !gres.IsEmpty() {
		for _, tryFolder := range resourceTryFolders {
			tempPath = tryFolder + usedFileName
			if file := gres.Get(tempPath); file != nil {
				path = file.Name()
				return
			}
		}
		c.searchPaths.RLockFunc(func(array []string) {
			for _, searchPath := range array {
				for _, tryFolder := range resourceTryFolders {
					tempPath = searchPath + tryFolder + usedFileName
					if file := gres.Get(tempPath); file != nil {
						path = file.Name()
						return
					}
				}
			}
		})
	}

	c.autoCheckAndAddMainPkgPathToSearchPaths()

	// Searching local file system.
	if path == "" {
		// Absolute path.
		if path = gfile.RealPath(usedFileName); path != "" {
			return
		}
		c.searchPaths.RLockFunc(func(array []string) {
			for _, searchPath := range array {
				searchPath = gstr.TrimRight(searchPath, `\/`)
				for _, tryFolder := range localSystemTryFolders {
					relativePath := gstr.TrimRight(
						gfile.Join(tryFolder, usedFileName),
						`\/`,
					)
					if path, _ = gspath.Search(searchPath, relativePath); path != "" {
						return
					}
				}
			}
		})
	}

	// If it cannot find the path of `file`, it formats and returns a detailed error.
	if path == "" {
		var buffer = bytes.NewBuffer(nil)
		if c.searchPaths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf(
				`config file "%s" not found in resource manager or the following system searching paths:`,
				usedFileName,
			))
			c.searchPaths.RLockFunc(func(array []string) {
				index := 1
				for _, searchPath := range array {
					searchPath = gstr.TrimRight(searchPath, `\/`)
					for _, tryFolder := range localSystemTryFolders {
						buffer.WriteString(fmt.Sprintf(
							"\n%d. %s",
							index, gfile.Join(searchPath, tryFolder),
						))
						index++
					}
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf(`cannot find config file "%s" with no path configured`, usedFileName))
		}
		err = gerror.NewCode(gcode.CodeNotFound, buffer.String())
	}
	return
}
