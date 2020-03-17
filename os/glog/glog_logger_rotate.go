// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/encoding/gcompress"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/text/gregex"
)

// rotateFile rotates the current logging file.
func (l *Logger) rotateFile() {
	// Rotation feature is not enabled as rotation file size is zero.
	if l.config.RotateSize == 0 {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	filePath := l.getFilePath()
	// No backups, it then just removes the current logging file.
	if l.config.RotateBackups == 0 {
		if err := gfile.Remove(filePath); err != nil {
			intlog.Print(err)
		}
		intlog.Printf(`%d size exceeds, no backups set, remove original logging file: %s`, l.config.RotateSize, filePath)
		return
	}
	// Else it creates new backup files.
	var (
		dirPath     = gfile.Dir(filePath)
		fileName    = gfile.Name(filePath)
		fileExt     = gfile.Ext(filePath)
		newFilePath = ""
	)
	for {
		// Rename the logging file by adding extra time information to milliseconds, like:
		// access              -> access.20200102190000899
		// access.log          -> access.20200102190000899.log
		// access.20200102.log -> access.20200102.20200102190000899.log
		newFilePath = gfile.Join(
			dirPath,
			fmt.Sprintf(`%s.%s%s`, fileName, gtime.Now().Format("YmdHisu"), fileExt),
		)
		if !gfile.Exists(newFilePath) {
			break
		}
	}
	if err := gfile.Rename(filePath, newFilePath); err != nil {
		panic(err)
	}
}

// rotateChecks timely checks the backups expiration and the compression.
func (l *Logger) rotateChecks() {
	defer func() {
		gtimer.AddOnce(l.config.RotateInterval, l.rotateChecks)
	}()

	// Checks whether file rotation not enabled.
	if l.config.RotateSize == 0 || l.config.RotateBackups == 0 {
		return
	}
	files, _ := gfile.ScanDirFile(l.config.Path, "*.*", true)
	intlog.Printf("logging rotation start checks: %+v", files)
	// Compression.
	needCompressFileArray := garray.NewStrArray()
	if l.config.RotateCompress > 0 {
		for _, file := range files {
			// Eg: access.20200102190000899.gz
			if gfile.ExtName(file) == "gz" {
				continue
			}
			// Eg:
			// access.20200102190000899
			// access.20200102190000899.log
			if gregex.IsMatchString(`.+\.\d{14,}`, file) {
				needCompressFileArray.Append(file)
			}
		}
		if needCompressFileArray.Len() > 0 {
			needCompressFileArray.Iterator(func(_ int, path string) bool {
				err := gcompress.GzipFile(path, path+".gz")
				if err == nil {
					intlog.Printf(`compressed done, remove original logging file: %s`, path)
					if err = gfile.Remove(path); err != nil {
						intlog.Print(err)
					}
				} else {
					intlog.Print(err)
				}
				return true
			})
			// Update the files array.
			files, _ = gfile.ScanDirFile(l.config.Path, "*.*", true)
		}
	}
	// Backups count limit and expiration checks.
	var (
		backupFilesMap          = make(map[string]*garray.SortedArray)
		originalLoggingFilePath = ""
	)
	if l.config.RotateBackups > 0 || l.config.RotateExpire > 0 {
		for _, file := range files {
			originalLoggingFilePath, _ = gregex.ReplaceString(`\.\d{14,}`, "", file)
			if backupFilesMap[originalLoggingFilePath] == nil {
				backupFilesMap[originalLoggingFilePath] = garray.NewSortedArray(func(a, b interface{}) int {
					// Sorted by backup file mtime.
					// The old backup file is put in the head of array.
					file1 := a.(string)
					file2 := b.(string)
					result := gfile.MTimeMillisecond(file1) - gfile.MTimeMillisecond(file2)
					if result <= 0 {
						return -1
					}
					return 1
				})
			}
			if gregex.IsMatchString(`.+\.\d{14,}`, file) {
				backupFilesMap[originalLoggingFilePath].Add(file)
			}
		}
		intlog.Printf(`calculated backup files map: %+v`, backupFilesMap)
		for _, array := range backupFilesMap {
			for i := 0; i < array.Len()-l.config.RotateBackups; i++ {
				path := array.PopLeft().(string)
				intlog.Printf(`remove exceeded backup file: %s`, path)
				if err := gfile.Remove(path); err != nil {
					intlog.Print(err)
				}
			}
		}
		// Expiration checks.
		if l.config.RotateExpire > 0 {
			nowTimestampMilli := gtime.TimestampMilli()
			// As for Golang version < 1.13, there's no method Milliseconds for time.Duration.
			// So we need calculate the milliseconds using its nanoseconds value.
			expireMillisecond := l.config.RotateExpire.Nanoseconds() / 1000000
			for _, array := range backupFilesMap {
				array.Iterator(func(_ int, v interface{}) bool {
					path := v.(string)
					mtime := gfile.MTimeMillisecond(path)
					differ := nowTimestampMilli - mtime
					if differ > expireMillisecond {
						intlog.Printf(
							`%d - %d = %d > %d, remove expired backup file: %s`,
							nowTimestampMilli, mtime, differ,
							expireMillisecond,
							path,
						)
						if err := gfile.Remove(path); err != nil {
							intlog.Print(err)
						}
						return true
					} else {
						return false
					}
				})
			}
		}
	}
}
