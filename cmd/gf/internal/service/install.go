package service

import (
	"context"
	"runtime"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/allyes"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	Install = serviceInstall{}
)

type serviceInstall struct{}

type serviceInstallAvailablePath struct {
	dirPath   string
	filePath  string
	writable  bool
	installed bool
}

func (s serviceInstall) Run(ctx context.Context) (err error) {
	// Ask where to install.
	paths := s.getAvailablePaths()
	if len(paths) <= 0 {
		mlog.Printf("no path detected, you can manually install gf by copying the binary to path folder.")
		return
	}
	mlog.Printf("I found some installable paths for you(from $PATH): ")
	mlog.Printf("  %2s | %8s | %9s | %s", "Id", "Writable", "Installed", "Path")

	// Print all paths status and determine the default selectedID value.
	var (
		selectedID = -1
		newPaths   []serviceInstallAvailablePath
		pathSet    = gset.NewStrSet() // Used for repeated items filtering.
	)
	for _, path := range paths {
		if !pathSet.AddIfNotExist(path.dirPath) {
			continue
		}
		newPaths = append(newPaths, path)
	}
	paths = newPaths
	for id, path := range paths {
		mlog.Printf("  %2d | %8t | %9t | %s", id, path.writable, path.installed, path.dirPath)
		if selectedID == -1 {
			// Use the previously installed path as the most priority choice.
			if path.installed {
				selectedID = id
			}
		}
	}
	// If there's no previously installed path, use the first writable path.
	if selectedID == -1 {
		// Order by choosing priority.
		commonPaths := garray.NewStrArrayFrom(g.SliceStr{
			s.getGoPathBin(),
			`/usr/local/bin`,
			`/usr/bin`,
			`/usr/sbin`,
			`C:\Windows`,
			`C:\Windows\system32`,
			`C:\Go\bin`,
			`C:\Program Files`,
			`C:\Program Files (x86)`,
		})
		// Check the common installation directories.
		commonPaths.Iterator(func(k int, v string) bool {
			for id, aPath := range paths {
				if strings.EqualFold(aPath.dirPath, v) {
					selectedID = id
					return false
				}
			}
			return true
		})
		if selectedID == -1 {
			selectedID = 0
		}
	}

	if allyes.Check() {
		// Use the default selectedID.
		mlog.Printf("please choose one installation destination [default %d]: %d", selectedID, selectedID)
	} else {
		for {
			// Get input and update selectedID.
			var (
				inputID int
				input   = gcmd.Scanf("please choose one installation destination [default %d]: ", selectedID)
			)
			if input != "" {
				inputID = gconv.Int(input)
			}
			// Check if out of range.
			if inputID >= len(paths) || inputID < 0 {
				mlog.Printf("invalid install destination Id: %d", inputID)
				continue
			}
			selectedID = inputID
			break
		}
	}

	// Get selected destination path.
	dstPath := paths[selectedID]

	// Install the new binary.
	err = gfile.CopyFile(gfile.SelfPath(), dstPath.filePath)
	if err != nil {
		mlog.Printf("install gf binary to '%s' failed: %v", dstPath.dirPath, err)
		mlog.Printf("you can manually install gf by copying the binary to folder: %s", dstPath.dirPath)
	} else {
		mlog.Printf("gf binary is successfully installed to: %s", dstPath.dirPath)
	}

	// Uninstall the old binary.
	for _, path := range paths {
		// Do not delete myself.
		if path.filePath != "" && path.filePath != dstPath.filePath && gfile.SelfPath() != path.filePath {
			_ = gfile.Remove(path.filePath)
		}
	}
	return
}

// IsInstalled checks and returns whether the binary is installed.
func (s serviceInstall) IsInstalled() bool {
	paths := s.getAvailablePaths()
	for _, aPath := range paths {
		if aPath.installed {
			return true
		}
	}
	return false
}

// getGoPathBinFilePath retrieves ad returns the GOPATH/bin path for binary.
func (s serviceInstall) getGoPathBin() string {
	if goPath := genv.Get(`GOPATH`).String(); goPath != "" {
		return gfile.Join(goPath, "bin")
	}
	return ""
}

// getAvailablePaths returns the installation paths data for the binary.
func (s serviceInstall) getAvailablePaths() []serviceInstallAvailablePath {
	var (
		folderPaths    []serviceInstallAvailablePath
		binaryFileName = "gf" + gfile.Ext(gfile.SelfPath())
	)
	// $GOPATH/bin
	if goPathBin := s.getGoPathBin(); goPathBin != "" {
		folderPaths = s.checkAndAppendToAvailablePath(
			folderPaths, goPathBin, binaryFileName,
		)
	}
	switch runtime.GOOS {
	case "darwin":
		darwinInstallationCheckPaths := []string{"/usr/local/bin"}
		for _, v := range darwinInstallationCheckPaths {
			folderPaths = s.checkAndAppendToAvailablePath(
				folderPaths, v, binaryFileName,
			)
		}
		fallthrough

	default:
		// Search and find the writable directory path.
		envPath := genv.Get("PATH", genv.Get("Path").String()).String()
		if gstr.Contains(envPath, ";") {
			for _, v := range gstr.SplitAndTrim(envPath, ";") {
				folderPaths = s.checkAndAppendToAvailablePath(
					folderPaths, v, binaryFileName,
				)
			}
		} else if gstr.Contains(envPath, ":") {
			for _, v := range gstr.SplitAndTrim(envPath, ":") {
				folderPaths = s.checkAndAppendToAvailablePath(
					folderPaths, v, binaryFileName,
				)
			}
		} else if envPath != "" {
			folderPaths = s.checkAndAppendToAvailablePath(
				folderPaths, envPath, binaryFileName,
			)
		} else {
			folderPaths = s.checkAndAppendToAvailablePath(
				folderPaths, "/usr/local/bin", binaryFileName,
			)
		}
	}
	return folderPaths
}

// checkAndAppendToAvailablePath checks if `path` is writable and already installed.
// It adds the `path` to `folderPaths` if it is writable or already installed, or else it ignores the `path`.
func (s serviceInstall) checkAndAppendToAvailablePath(folderPaths []serviceInstallAvailablePath, dirPath string, binaryFileName string) []serviceInstallAvailablePath {
	var (
		filePath  = gfile.Join(dirPath, binaryFileName)
		writable  = gfile.IsWritable(dirPath)
		installed = gfile.Exists(filePath)
	)
	if !writable && !installed {
		return folderPaths
	}
	return append(
		folderPaths,
		serviceInstallAvailablePath{
			dirPath:   dirPath,
			writable:  writable,
			filePath:  filePath,
			installed: installed,
		})
}
