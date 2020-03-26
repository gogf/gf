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
	"time"
)

// rotateFileBySize rotates the current logging file according to the
// configured rotation size.
func (l *Logger) rotateFileBySize(now time.Time) {
	if l.config.RotateSize <= 0 {
		return
	}
	l.rmu.Lock()
	defer l.rmu.Unlock()
	if err := l.doRotateFile(l.getFilePath(now)); err != nil {
		panic(err)
	}
}

// doRotateFile rotates the given logging file.
func (l *Logger) doRotateFile(filePath string) error {
	// No backups, it then just removes the current logging file.
	if l.config.RotateBackLimit == 0 {
		if err := gfile.Remove(filePath); err != nil {
			return err
		}
		intlog.Printf(`%d size exceeds, no backups set, remove original logging file: %s`, l.config.RotateSize, filePath)
		return nil
	}
	// Else it creates new backup files.
	var (
		dirPath     = gfile.Dir(filePath)
		fileName    = gfile.Name(filePath)
		fileExt     = gfile.Ext(filePath)
		newFilePath = ""
	)
	// Rename the logging file by adding extra time information to milliseconds, like:
	// access.log          -> access.20200102190000899.log
	// access.20200102.log -> access.20200102.20200102190000899.log
	newFilePath = gfile.Join(
		dirPath,
		fmt.Sprintf(`%s.%s%s`, fileName, gtime.Now().Format("YmdHisu"), fileExt),
	)
	if err := gfile.Rename(filePath, newFilePath); err != nil {
		return err
	}
	return nil
}

// rotateChecksTimely timely checks the backups expiration and the compression.
func (l *Logger) rotateChecksTimely() {
	defer gtimer.AddOnce(l.config.RotateInterval, l.rotateChecksTimely)
	// Checks whether file rotation not enabled.
	if l.config.RotateSize <= 0 && l.config.RotateExpire == 0 {
		return
	}
	var (
		now      = time.Now()
		pattern  = "*.log, *.gz"
		files, _ = gfile.ScanDirFile(l.config.Path, pattern, true)
	)
	intlog.Printf("logging rotation start checks: %+v", files)
	// =============================================================
	// Rotation expire file checks.
	// =============================================================
	if l.config.RotateExpire > 0 {
		var (
			mtime         time.Time
			subDuration   time.Duration
			expireRotated bool
		)
		for _, file := range files {
			if gfile.ExtName(file) == "gz" {
				continue
			}
			mtime = gfile.MTime(file)
			subDuration = now.Sub(mtime)
			if subDuration > l.config.RotateExpire {
				expireRotated = true
				intlog.Printf(
					`%v - %v = %v > %v, rotation expire logging file: %s`,
					now, mtime, subDuration, l.config.RotateExpire, file,
				)
				if err := l.doRotateFile(file); err != nil {
					intlog.Error(err)
				}
			}
		}
		if expireRotated {
			// Update the files array.
			files, _ = gfile.ScanDirFile(l.config.Path, pattern, true)
		}
	}

	// =============================================================
	// Rotated file compression.
	// =============================================================
	needCompressFileArray := garray.NewStrArray()
	if l.config.RotateBackCompress > 0 {
		for _, file := range files {
			// Eg: access.20200102190000899.log.gz
			if gfile.ExtName(file) == "gz" {
				continue
			}
			// Eg:
			// access.20200102190000899
			// access.20200102190000899.log
			if gregex.IsMatchString(`.+\.\d{14,}`, gfile.Name(file)) {
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
			files, _ = gfile.ScanDirFile(l.config.Path, pattern, true)
		}
	}

	// =============================================================
	// Backups count limit and expiration checks.
	// =============================================================
	var (
		backupFilesMap          = make(map[string]*garray.SortedArray)
		originalLoggingFilePath = ""
	)
	if l.config.RotateBackLimit > 0 || l.config.RotateBackExpire > 0 {
		for _, file := range files {
			originalLoggingFilePath, _ = gregex.ReplaceString(`\.\d{14,}`, "", file)
			if backupFilesMap[originalLoggingFilePath] == nil {
				backupFilesMap[originalLoggingFilePath] = garray.NewSortedArray(func(a, b interface{}) int {
					// Sorted by rotated/backup file mtime.
					// The old rotated/backup file is put in the head of array.
					file1 := a.(string)
					file2 := b.(string)
					result := gfile.MTimestampMilli(file1) - gfile.MTimestampMilli(file2)
					if result <= 0 {
						return -1
					}
					return 1
				})
			}
			// Check if this file a rotated/backup file.
			if gregex.IsMatchString(`.+\.\d{14,}`, gfile.Name(file)) {
				backupFilesMap[originalLoggingFilePath].Add(file)
			}
		}
		intlog.Printf(`calculated backup files map: %+v`, backupFilesMap)
		for _, array := range backupFilesMap {
			diff := array.Len() - l.config.RotateBackLimit
			for i := 0; i < diff; i++ {
				path := array.PopLeft().(string)
				intlog.Printf(`remove exceeded backup file: %s`, path)
				if err := gfile.Remove(path); err != nil {
					intlog.Print(err)
				}
			}
		}
		// Backup expiration checks.
		if l.config.RotateBackExpire > 0 {
			var (
				mtime       time.Time
				subDuration time.Duration
			)
			for _, array := range backupFilesMap {
				array.Iterator(func(_ int, v interface{}) bool {
					path := v.(string)
					mtime = gfile.MTime(path)
					subDuration = now.Sub(mtime)
					if subDuration > l.config.RotateBackExpire {
						intlog.Printf(
							`%v - %v = %v > %v, remove expired backup file: %s`,
							now, mtime, subDuration, l.config.RotateBackExpire, path,
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
