<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件管理.
package gfile

import (
    "os"
    "io"
    "io/ioutil"
    "sort"
    "fmt"
    "time"
    "strings"
    "bytes"
    "os/exec"
    "errors"
    "os/user"
    "runtime"
    "path/filepath"
    "gitee.com/johng/gf/g/util/gregx"
    "gitee.com/johng/gf/g/container/gtype"
)

// 封装了常用的文件操作方法，如需更详细的文件控制，请查看官方os包

// 文件分隔符
const (
    Separator = string(filepath.Separator)
)

// 源码的main包所在目录，仅仅会设置一次
var mainPkgPath = gtype.NewInterface()

// 给定文件的绝对路径创建文件
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gfile provides easy-to-use operations for file system.
package gfile

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/text/gregex"
	"github.com/gogf/gf/g/text/gstr"
	"github.com/gogf/gf/g/util/gconv"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

const (
    // Separator for file system.
    Separator     = string(filepath.Separator)
    // Default perm for file opening.
    gDEFAULT_PERM = 0666
)

var (
    // The absolute file path for main package.
    // It can be only checked and set once.
    mainPkgPath = gtype.NewString()
)

// Mkdir creates directories recursively with given <path>.
// The parameter <path> is suggested to be absolute path.
>>>>>>> upstream/master
func Mkdir(path string) error {
    err  := os.MkdirAll(path, os.ModePerm)
    if err != nil {
        return err
    }
    return nil
}

<<<<<<< HEAD
// 给定文件的绝对路径创建文件
func Create(path string) error {
=======
// Create creates file with given <path> recursively.
// The parameter <path> is suggested to be absolute path.
func Create(path string) (*os.File, error) {
>>>>>>> upstream/master
    dir := Dir(path)
    if !Exists(dir) {
        Mkdir(dir)
    }
<<<<<<< HEAD
    f, err  := os.Create(path)
    if err != nil {
        return err
    }
    f.Close()
    return nil
}

// 打开文件
func Open(path string) (*os.File, error) {
    f, err  := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
=======
    return os.Create(path)
}

// Open opens file/directory readonly.
func Open(path string) (*os.File, error) {
    return os.Open(path)
}

// OpenFile opens file/directory with given <flag> and <perm>.
func OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
    return os.OpenFile(path, flag, perm)
}

// OpenWithFlag opens file/directory with default perm and given <flag>.
func OpenWithFlag(path string, flag int) (*os.File, error) {
    f, err  := os.OpenFile(path, flag, gDEFAULT_PERM)
>>>>>>> upstream/master
    if err != nil {
        return nil, err
    }
    return f, nil
}

<<<<<<< HEAD
// 打开文件
func OpenWithFlag(path string, flag int) (*os.File, error) {
    f, err  := os.OpenFile(path, flag, 0666)
=======
// OpenWithFlagPerm opens file/directory with given <flag> and <perm>.
func OpenWithFlagPerm(path string, flag int, perm int) (*os.File, error) {
    f, err  := os.OpenFile(path, flag, os.FileMode(perm))
>>>>>>> upstream/master
    if err != nil {
        return nil, err
    }
    return f, nil
}

<<<<<<< HEAD
// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
    _, err := os.Stat(path)
    if err != nil {
        if os.IsExist(err) {
            return true
        }
        return false
    }
    return true
}

// 判断所给路径是否为文件夹
=======
// Exists checks whether given <path> exist.
func Exists(path string) bool {
    if _, err := os.Stat(path); !os.IsNotExist(err) {
        return true
    }
    return false
}

// IsDir checks whether given <path> a directory.
>>>>>>> upstream/master
func IsDir(path string) bool {
    s, err := os.Stat(path)
    if err != nil {
        return false
    }
    return s.IsDir()
}

<<<<<<< HEAD
// 判断所给路径是否为文件
func IsFile(path string) bool {
    return !IsDir(path)
}

// 获取文件或目录信息
func Info(path string) *os.FileInfo {
    info, err := os.Stat(path)
    if err != nil {
        return nil
    }
    return &info
}

// 修改时间(秒)
func MTime(path string) int64 {
    f, e := os.Stat(path)
    if e != nil {
        return 0
    }
    return f.ModTime().Unix()
}

// 修改时间(毫秒)
func MTimeMillisecond(path string) int64 {
    f, e := os.Stat(path)
    if e != nil {
        return 0
    }
    return int64(f.ModTime().Nanosecond()/1000000)
}

// 文件大小(bytes)
func Size(path string) int64 {
    f, e := os.Stat(path)
    if e != nil {
        return 0
    }
    return f.Size()
}

// 格式化文件大小
func ReadableSize(path string) string {
    return FormatSize(float64(Size(path)))
}

// 格式化文件大小
func FormatSize(raw float64) string {
    var t float64 = 1024
    var d float64 = 1

    if raw < t {
        return fmt.Sprintf("%.2fB", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fK", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fM", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fG", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fT", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fP", raw/d)
    }

    return "TooLarge"
}

// 文件移动/重命名
=======
// Pwd returns absolute path of current working directory.
func Pwd() string {
    path, _ := os.Getwd()
    return path
}

// IsFile checks whether given <path> a file, which means it's not a directory.
func IsFile(path string) bool {
    s, err := os.Stat(path)
    if err != nil {
        return false
    }
    return !s.IsDir()
}

// Alias of Stat.
// See Stat.
func Info(path string) (os.FileInfo, error) {
    return Stat(path)
}

// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.
func Stat(path string) (os.FileInfo, error) {
    return os.Stat(path)
}

// Move renames (moves) <src> to <dst> path.
>>>>>>> upstream/master
func Move(src string, dst string) error {
    return os.Rename(src, dst)
}

<<<<<<< HEAD

// 文件移动/重命名
=======
// Alias of Move.
// See Move.
>>>>>>> upstream/master
func Rename(src string, dst string) error {
    return Move(src, dst)
}

<<<<<<< HEAD
// 文件复制
func Copy(src string, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
=======
// Copy file from <src> to <dst>.
//
// @TODO directory copy support.
func Copy(src string, dst string) error {
    srcFile, err := Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()
    dstFile, err := Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()
>>>>>>> upstream/master
    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        return err
    }
    err = dstFile.Sync()
    if err != nil {
        return err
    }
<<<<<<< HEAD
    srcFile.Close()
    dstFile.Close()
    return nil
}

// 文件/目录删除
=======
    return nil
}

// DirNames returns sub-file names of given directory <path>.
func DirNames(path string) ([]string, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    list, err := f.Readdirnames(-1)
    f.Close()
    if err != nil {
        return nil, err
    }
    return list, nil
}

// Glob returns the names of all files matching pattern or nil
// if there is no matching file. The syntax of patterns is the same
// as in Match. The pattern may describe hierarchical names such as
// /usr/*/bin/ed (assuming the Separator is '/').
//
// Glob ignores file system errors such as I/O errors reading directories.
// The only possible returned error is ErrBadPattern, when pattern
// is malformed.
func Glob(pattern string, onlyNames...bool) ([]string, error) {
    if list, err := filepath.Glob(pattern); err == nil {
        if len(onlyNames) > 0 && onlyNames[0] && len(list) > 0 {
            array := make([]string, len(list))
            for k, v := range list {
                array[k] = Basename(v)
            }
            return array, nil
        }
        return list, nil
    } else {
        return nil, err
    }
}

// Remove deletes all file/directory with <path> parameter.
// If parameter <path> is directory, it deletes it recursively.
>>>>>>> upstream/master
func Remove(path string) error {
    return os.RemoveAll(path)
}

<<<<<<< HEAD
// 文件是否可
func IsReadable(path string) bool {
    result    := true
    file, err := os.OpenFile(path, os.O_RDONLY, 0666)
=======
// IsReadable checks whether given <path> is readable.
func IsReadable(path string) bool {
    result    := true
    file, err := os.OpenFile(path, os.O_RDONLY, gDEFAULT_PERM)
>>>>>>> upstream/master
    if err != nil {
        result = false
    }
    file.Close()
    return result
}

<<<<<<< HEAD
// 文件是否可写
func IsWritable(path string) bool {
    result := true
    if IsDir(path) {
        // 如果是目录，那么创建一个临时文件进行写入测试
        tfile := strings.TrimRight(path, Separator) + Separator + string(time.Now().UnixNano())
        err   := Create(tfile)
        if err != nil || !Exists(tfile){
            result = false
        } else {
            Remove(tfile)
        }
    } else {
        // 如果是文件，那么判断文件是否可打开
        file, err := os.OpenFile(path, os.O_WRONLY, 0666)
=======
// IsWritable checks whether given <path> is writable.
//
// @TODO improve performance; use golang.org/x/sys to cross-plat-form
func IsWritable(path string) bool {
    result := true
    if IsDir(path) {
        // If it's a directory, create a temporary file to test whether it's writable.
        tmpFile := strings.TrimRight(path, Separator) + Separator + gconv.String(time.Now().UnixNano())
        if f, err := Create(tmpFile); err != nil || !Exists(tmpFile){
            result = false
        } else {
            f.Close()
            Remove(tmpFile)
        }
    } else {
        // 如果是文件，那么判断文件是否可打开
        file, err := os.OpenFile(path, os.O_WRONLY, gDEFAULT_PERM)
>>>>>>> upstream/master
        if err != nil {
            result = false
        }
        file.Close()
    }
    return result
}

<<<<<<< HEAD
// 修改文件/目录权限
=======
// See os.Chmod.
>>>>>>> upstream/master
func Chmod(path string, mode os.FileMode) error {
    return os.Chmod(path, mode)
}

<<<<<<< HEAD
// 打开目录，并返回其下一级子目录名称列表，按照文件名称大小写进行排序
func ScanDir(path string) []string {
    f, err := os.Open(path)
    if err != nil {
        return nil
    }

    list, err := f.Readdirnames(-1)
    f.Close()
    if err != nil {
        return nil
    }
    sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
    return list
}

// 将所给定的路径转换为绝对路径
// 并判断文件路径是否存在，如果文件不存在，那么返回空字符串
=======
// ScanDir returns all sub-files with absolute paths of given <path>,
// It scans directory recursively if given parameter <recursive> is true.
func ScanDir(path string, pattern string, recursive ... bool) ([]string, error) {
    list, err := doScanDir(path, pattern, recursive...)
    if err != nil {
        return nil, err
    }
    if len(list) > 0 {
        sort.Strings(list)
    }
    return list, nil
}

// doScanDir is an internal method which scans directory
// and returns the absolute path list of files that are not sorted.
//
// The pattern parameter <pattern> supports multiple file name patterns,
// using the ',' symbol to separate multiple patterns.
//
// It scans directory recursively if given parameter <recursive> is true.
func doScanDir(path string, pattern string, recursive ... bool) ([]string, error) {
    list := ([]string)(nil)
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    names, err := file.Readdirnames(-1)
    if err != nil {
        return nil, err
    }
    for _, name := range names {
        path := fmt.Sprintf("%s%s%s", path, Separator, name)
        if IsDir(path) && len(recursive) > 0 && recursive[0] {
            array, _ := doScanDir(path, pattern, true)
            if len(array) > 0 {
                list = append(list, array...)
            }
        }
        // If it meets pattern, then add it to the result list.
        for _, p := range strings.Split(pattern, ",") {
            if match, err := filepath.Match(strings.TrimSpace(p), name); err == nil && match {
                list = append(list, path)
            }
        }
    }
    return list, nil
}

// RealPath converts the given <path> to its absolute path
// and checks if the file path exists.
// If the file does not exist, return an empty string.
>>>>>>> upstream/master
func RealPath(path string) string {
    p, err := filepath.Abs(path)
    if err != nil {
        return ""
    }
    if !Exists(p) {
        return ""
    }
    return p
}

<<<<<<< HEAD
// (文本)读取文件内容
func GetContents(path string) string {
    return string(GetBinContents(path))
}

// (二进制)读取文件内容
func GetBinContents(path string) []byte {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil
    }
    return data
}

// 写入文件内容
func putContents(path string, data []byte, flag int, perm os.FileMode) error {
    // 支持目录递归创建
    dir := Dir(path)
    if !Exists(dir) {
        Mkdir(dir)
    }
    // 创建/打开文件
    f, err := os.OpenFile(path, flag, perm)
    if err != nil {
        return err
    }
    defer f.Close()
    n, err := f.Write(data)
    if err != nil {
        return err
    } else if n < len(data) {
        return io.ErrShortWrite
    }
    return nil
}

// Truncate
func Truncate(path string, size int) error {
    return os.Truncate(path, int64(size))
}

// (文本)写入文件内容
func PutContents(path string, content string) error {
    return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
}

// (文本)追加内容到文件末尾
func PutContentsAppend(path string, content string) error {
    return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}

// (二进制)写入文件内容
func PutBinContents(path string, content []byte) error {
    return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
}

// (二进制)追加内容到文件末尾
func PutBinContentsAppend(path string, content []byte) error {
    return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}


// 获取当前执行文件的绝对路径
=======
// SelfPath returns absolute file path of current running process(binary).
>>>>>>> upstream/master
func SelfPath() string {
    p, _ := filepath.Abs(os.Args[0])
    return p
}

<<<<<<< HEAD
// 获取当前执行文件的目录绝对路径
=======
// SelfDir returns absolute directory path of current running process(binary).
>>>>>>> upstream/master
func SelfDir() string {
    return filepath.Dir(SelfPath())
}

<<<<<<< HEAD
// 获取指定文件路径的文件名称
=======
// Basename returns the last element of path.
// Trailing path separators are removed before extracting the last element.
// If the path is empty, Base returns ".".
// If the path consists entirely of separators, Base returns a single separator.
>>>>>>> upstream/master
func Basename(path string) string {
    return filepath.Base(path)
}

<<<<<<< HEAD
// 获取指定文件路径的目录地址绝对路径
=======
// Dir returns all but the last element of path, typically the path's directory.
// After dropping the final element, Dir calls Clean on the path and trailing
// slashes are removed.
// If the path is empty, Dir returns ".".
// If the path consists entirely of separators, Dir returns a single separator.
// The returned path does not end in a separator unless it is the root directory.
>>>>>>> upstream/master
func Dir(path string) string {
    return filepath.Dir(path)
}

<<<<<<< HEAD
// 获取指定文件路径的文件扩展名
=======
// Ext returns the file name extension used by path.
// The extension is the suffix beginning at the final dot
// in the final element of path; it is empty if there is
// no dot.
//
// Note: the result contains symbol '.'.
>>>>>>> upstream/master
func Ext(path string) string {
    return filepath.Ext(path)
}

<<<<<<< HEAD
// 获取用户主目录
=======
// Home returns absolute path of current user's home directory.
>>>>>>> upstream/master
func Home() (string, error) {
    u, err := user.Current()
    if nil == err {
        return u.HomeDir, nil
    }
    if "windows" == runtime.GOOS {
        return homeWindows()
    }
    return homeUnix()
}

func homeUnix() (string, error) {
    if home := os.Getenv("HOME"); home != "" {
        return home, nil
    }
    var stdout bytes.Buffer
    cmd := exec.Command("sh", "-c", "eval echo ~$USER")
    cmd.Stdout = &stdout
    if err := cmd.Run(); err != nil {
        return "", err
    }

    result := strings.TrimSpace(stdout.String())
    if result == "" {
        return "", errors.New("blank output when reading home directory")
    }

    return result, nil
}

func homeWindows() (string, error) {
    drive := os.Getenv("HOMEDRIVE")
    path  := os.Getenv("HOMEPATH")
    home  := drive + path
    if drive == "" || path == "" {
        home = os.Getenv("USERPROFILE")
    }
    if home == "" {
        return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
    }

    return home, nil
}

<<<<<<< HEAD
// 获得文件内容下一个指定字节的位置
func GetNextCharOffset(file *os.File, char string, start int64) int64 {
    c := []byte(char)[0]
    b := make([]byte, 1)
    o := start
    for {
        _, err := file.ReadAt(b, o)
        if err != nil {
            return 0
        }
        if b[0] == c {
            return o
        }
        o++
    }
    return 0
}

// 获得文件内容中两个offset之间的内容 [start, end)
func GetBinContentByTwoOffsets(file *os.File, start int64, end int64) []byte {
    buffer := make([]byte, end - start)
    if _, err := file.ReadAt(buffer, start); err != nil {
        return nil
    }
    return buffer
}

// 获取入口函数文件所在目录(main包文件目录)，仅对源码开发环境有效（即仅对生成该可执行文件的系统下有效）
func MainPkgPath() string {
    path := mainPkgPath.Val()
    if path != nil {
        return path.(string)
    }
    f := ""
    for i := 1; i < 10000; i++ {
        if _, file, _, ok := runtime.Caller(i); ok {
            // 不包含go源码路径
            if !gregx.IsMatchString("^" + runtime.GOROOT(), file) {
                f = file
            }
=======
// MainPkgPath returns absolute file path of package main,
// which contains the entrance function main.
//
// It's only available in develop environment.
//
// Note1: Only valid for source development environments,
// IE only valid for systems that generate this executable.
// Note2: When the method is called for the first time, if it is in an asynchronous goroutine,
// the method may not get the main package path.
func MainPkgPath() string {
    path := mainPkgPath.Val()
    if path != "" {
    	if path == "-" {
    		return ""
	    }
        return path
    }
    for i := 1; i < 10000; i++ {
        if _, file, _, ok := runtime.Caller(i); ok {
        	// <file> is separated by '/'
        	if gstr.Contains(file, "/gf/g/") {
        		continue
	        }
        	if Ext(file) != ".go" {
        		continue
	        }
        	// separator of <file> '/' will be converted to Separator.
        	for path = Dir(file); len(path) > 1 && Exists(path) && path[len(path) - 1] != os.PathSeparator; {
        		files, _ := ScanDir(path, "*.go")
        		for _, v := range files {
			        if gregex.IsMatchString(`package\s+main`, GetContents(v)) {
				        mainPkgPath.Set(path)
				        return path
			        }
		        }
        		path = Dir(path)
	        }

>>>>>>> upstream/master
        } else {
            break
        }
    }
<<<<<<< HEAD
    if f != "" {
        p := Dir(f)
        mainPkgPath.Set(p)
        return p
    }
    return ""
}

// 系统临时目录
func TempDir() string {
    return os.TempDir()
}
=======
    // If it fails finding the path, then mark it as "-",
    // which means it will never do this search again.
	mainPkgPath.Set("-")
    return ""
}

// See os.TempDir().
func TempDir() string {
    return os.TempDir()
}
>>>>>>> upstream/master
