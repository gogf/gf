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
	"os"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/os/gspath"
	"github.com/gogf/gf/v2/text/gstr"
)

// SetPath sets the configuration `directory` path for file search.
// The parameter `path` can be absolute or relative `directory` path,
// but absolute `directory` path is strongly recommended.
//
// Note that this parameter is a path to a directory not a file.
func (a *AdapterFile) SetPath(directoryPath string) (err error) {
	var (
		isDir    = false
		realPath = ""
	)
	if file := gres.Get(directoryPath); file != nil {
		realPath = directoryPath
		isDir = file.FileInfo().IsDir()
	} else {
		// Absolute path.
		realPath = gfile.RealPath(directoryPath)
		if realPath == "" {
			// Relative path.
			a.searchPaths.RLockFunc(func(array []string) {
				for _, v := range array {
					if searchedPath, _ := gspath.Search(v, directoryPath); searchedPath != "" {
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
		if a.searchPaths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf(
				`SetPath failed: cannot find directory "%s" in following paths:`,
				directoryPath,
			))
			a.searchPaths.RLockFunc(func(array []string) {
				for k, v := range array {
					buffer.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf(
				`SetPath failed: path "%s" does not exist`,
				directoryPath,
			))
		}
		return gerror.New(buffer.String())
	}
	// Should be a directory.
	if !isDir {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`SetPath failed: path "%s" should be directory type`,
			directoryPath,
		)
	}
	// Repeated path check.
	if a.searchPaths.Search(realPath) != -1 {
		return nil
	}
	a.jsonMap.Clear()
	a.searchPaths.Clear()
	a.searchPaths.Append(realPath)
	intlog.Print(context.TODO(), "SetPath:", realPath)
	return nil
}

// AddPath adds an absolute or relative `directory` path to the search paths.
//
// Note that this parameter is paths to a directories not files.
func (a *AdapterFile) AddPath(directoryPaths ...string) (err error) {
	for _, directoryPath := range directoryPaths {
		if err = a.doAddPath(directoryPath); err != nil {
			return err
		}
	}
	return nil
}

// doAddPath adds an absolute or relative `directory` path to the search paths.
func (a *AdapterFile) doAddPath(directoryPath string) (err error) {
	var (
		isDir    = false
		realPath = ""
	)
	// It firstly checks the resource manager,
	// and then checks the filesystem for the path.
	if file := gres.Get(directoryPath); file != nil {
		realPath = directoryPath
		isDir = file.FileInfo().IsDir()
	} else {
		// Absolute path.
		realPath = gfile.RealPath(directoryPath)
		if realPath == "" {
			// Relative path.
			a.searchPaths.RLockFunc(func(array []string) {
				for _, v := range array {
					if searchedPath, _ := gspath.Search(v, directoryPath); searchedPath != "" {
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
		if a.searchPaths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf(
				`AddPath failed: cannot find directory "%s" in following paths:`,
				directoryPath,
			))
			a.searchPaths.RLockFunc(func(array []string) {
				for k, v := range array {
					buffer.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf(
				`AddPath failed: path "%s" does not exist`,
				directoryPath,
			))
		}
		return gerror.New(buffer.String())
	}
	if !isDir {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`AddPath failed: path "%s" should be directory type`,
			directoryPath,
		)
	}
	// Repeated path check.
	if a.searchPaths.Search(realPath) != -1 {
		return nil
	}
	a.searchPaths.Append(realPath)
	intlog.Print(context.TODO(), "AddPath:", realPath)
	return nil
}

// GetPaths returns the searching directory path array of current configuration manager.
func (a *AdapterFile) GetPaths() []string {
	return a.searchPaths.Slice()
}

// doGetFilePath returns the absolute configuration file path for the given filename by `file`.
// If `file` is not passed, it returns the configuration file path of the default name.
// It returns an empty `path` string and an error if the given `file` does not exist.
func (a *AdapterFile) doGetFilePath(fileName string) (filePath string) {
	var (
		tempPath string
		resFile  *gres.File
		fileInfo os.FileInfo
	)
	// Searching resource manager.
	if !gres.IsEmpty() {
		for _, tryFolder := range resourceTryFolders {
			tempPath = tryFolder + fileName
			if resFile = gres.Get(tempPath); resFile != nil {
				fileInfo, _ = resFile.Stat()
				if fileInfo != nil && !fileInfo.IsDir() {
					filePath = resFile.Name()
					return
				}
			}
		}
		a.searchPaths.RLockFunc(func(array []string) {
			for _, searchPath := range array {
				for _, tryFolder := range resourceTryFolders {
					tempPath = searchPath + tryFolder + fileName
					if resFile = gres.Get(tempPath); resFile != nil {
						fileInfo, _ = resFile.Stat()
						if fileInfo != nil && !fileInfo.IsDir() {
							filePath = resFile.Name()
							return
						}
					}
				}
			}
		})
	}

	a.autoCheckAndAddMainPkgPathToSearchPaths()

	// Searching local file system.
	if filePath == "" {
		// Absolute path.
		if filePath = gfile.RealPath(fileName); filePath != "" && !gfile.IsDir(filePath) {
			return
		}
		a.searchPaths.RLockFunc(func(array []string) {
			for _, searchPath := range array {
				searchPath = gstr.TrimRight(searchPath, `\/`)
				for _, tryFolder := range localSystemTryFolders {
					relativePath := gstr.TrimRight(
						gfile.Join(tryFolder, fileName),
						`\/`,
					)
					if filePath, _ = gspath.Search(searchPath, relativePath); filePath != "" &&
						!gfile.IsDir(filePath) {
						return
					}
				}
			}
		})
	}
	return
}

// GetFilePath returns the absolute configuration file path for the given filename by `file`.
// If `file` is not passed, it returns the configuration file path of the default name.
// It returns an empty `path` string and an error if the given `file` does not exist.
func (a *AdapterFile) GetFilePath(fileName ...string) (filePath string, err error) {
	var (
		fileExtName  string
		tempFileName string
		usedFileName = a.defaultName
	)
	if len(fileName) > 0 {
		usedFileName = fileName[0]
	}
	fileExtName = gfile.ExtName(usedFileName)
	if filePath = a.doGetFilePath(usedFileName); (filePath == "" || gfile.IsDir(filePath)) && !gstr.InArray(supportedFileTypes, fileExtName) {
		// If it's not using default configuration or its configuration file is not available,
		// it searches the possible configuration file according to the name and all supported
		// file types.
		for _, fileType := range supportedFileTypes {
			tempFileName = fmt.Sprintf(`%s.%s`, usedFileName, fileType)
			if filePath = a.doGetFilePath(tempFileName); filePath != "" {
				break
			}
		}
	}
	// If it cannot find the filePath of `file`, it formats and returns a detailed error.
	if filePath == "" {
		var buffer = bytes.NewBuffer(nil)
		if a.searchPaths.Len() > 0 {
			if !gstr.InArray(supportedFileTypes, fileExtName) {
				buffer.WriteString(fmt.Sprintf(
					`possible config files "%s" or "%s" not found in resource manager or following system searching paths:`,
					usedFileName, fmt.Sprintf(`%s.%s`, usedFileName, gstr.Join(supportedFileTypes, "/")),
				))
			} else {
				buffer.WriteString(fmt.Sprintf(
					`specified config file "%s" not found in resource manager or following system searching paths:`,
					usedFileName,
				))
			}
			a.searchPaths.RLockFunc(func(array []string) {
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
			buffer.WriteString(fmt.Sprintf(`cannot find config file "%s" with no filePath configured`, usedFileName))
		}
		err = gerror.NewCode(gcode.CodeNotFound, buffer.String())
	}
	return
}
