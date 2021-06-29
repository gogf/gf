package install

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/tool/gf/library/allyes"
	"github.com/gogf/gf/tool/gf/library/mlog"
	"github.com/gogf/gf/util/gconv"
	"runtime"
	"strings"
)

// installFolderPath contains installFolderPath-related data,
// such as path, writable, binaryFilePath, and installed.
type installFolderPath struct {
	path           string
	writable       bool
	binaryFilePath string
	installed      bool
}

// Run does the installation.
func Run() {
	// Ask where to install.
	paths := getInstallPathsData()
	if len(paths) <= 0 {
		mlog.Printf("no path detected, you can manually install gf by copying the binary to path folder.")
		return
	}
	mlog.Printf("I found some installable paths for you(from $PATH): ")
	mlog.Printf("  %2s | %8s | %9s | %s", "Id", "Writable", "Installed", "Path")

	// Print all paths status and determine the default selectedID value.
	var (
		selectedID = -1
		pathSet    = gset.NewStrSet() // Used for repeated items filtering.
	)
	for id, aPath := range paths {
		if !pathSet.AddIfNotExist(aPath.path) {
			continue
		}
		mlog.Printf("  %2d | %8t | %9t | %s", id, aPath.writable, aPath.installed, aPath.path)
		if selectedID == -1 {
			// Use the previously installed path as the most priority choice.
			if aPath.installed {
				selectedID = id
			}
		}
	}
	// If there's no previously installed path, use the first writable path.
	if selectedID == -1 {
		// Order by choosing priority.
		commonPaths := garray.NewStrArrayFrom(g.SliceStr{
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
				if strings.EqualFold(aPath.path, v) {
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
		// Get input and update selectedID.
		input := gcmd.Scanf("please choose one installation destination [default %d]: ", selectedID)
		if input != "" {
			selectedID = gconv.Int(input)
		}
	}

	// Check if out of range.
	if selectedID >= len(paths) || selectedID < 0 {
		mlog.Printf("invalid install destination Id: %d", selectedID)
		return
	}

	// Get selected destination path.
	dstPath := paths[selectedID]

	// Install the new binary.
	err := gfile.CopyFile(gfile.SelfPath(), dstPath.binaryFilePath)
	if err != nil {
		mlog.Printf("install gf binary to '%s' failed: %v", dstPath.path, err)
		mlog.Printf("you can manually install gf by copying the binary to folder: %s", dstPath.path)
	} else {
		mlog.Printf("gf binary is successfully installed to: %s", dstPath.path)
	}

	// Uninstall the old binary.
	for _, aPath := range paths {
		// Do not delete myself.
		if aPath.binaryFilePath != "" &&
			aPath.binaryFilePath != dstPath.binaryFilePath &&
			gfile.SelfPath() != aPath.binaryFilePath {
			gfile.Remove(aPath.binaryFilePath)
		}
	}
}

// IsInstalled returns whether the binary is installed.
func IsInstalled() bool {
	paths := getInstallPathsData()
	for _, aPath := range paths {
		if aPath.installed {
			return true
		}
	}
	return false
}

// GetInstallPathsData returns the installation paths data for the binary.
func getInstallPathsData() []installFolderPath {
	var folderPaths []installFolderPath
	// Pre generate binaryFileName.
	binaryFileName := "gf" + gfile.Ext(gfile.SelfPath())
	switch runtime.GOOS {
	case "darwin":
		folderPaths = checkPathAndAppendToInstallFolderPath(
			folderPaths, "/usr/local/bin", binaryFileName,
		)
	default:
		// Search and find the writable directory path.
		envPath := genv.Get("PATH", genv.Get("Path"))
		if gstr.Contains(envPath, ";") {
			for _, v := range gstr.SplitAndTrim(envPath, ";") {
				folderPaths = checkPathAndAppendToInstallFolderPath(
					folderPaths, v, binaryFileName)
			}
		} else if gstr.Contains(envPath, ":") {
			for _, v := range gstr.SplitAndTrim(envPath, ":") {
				folderPaths = checkPathAndAppendToInstallFolderPath(
					folderPaths, v, binaryFileName)
			}
		} else if envPath != "" {
			folderPaths = checkPathAndAppendToInstallFolderPath(
				folderPaths, envPath, binaryFileName)
		} else {
			folderPaths = checkPathAndAppendToInstallFolderPath(
				folderPaths, "/usr/local/bin", binaryFileName)
		}
	}
	return folderPaths
}

// checkPathAndAppendToInstallFolderPath checks if <path> is writable and already installed.
// It adds the <path> to <folderPaths> if it is writable or already installed, or else it ignores the <path>.
func checkPathAndAppendToInstallFolderPath(folderPaths []installFolderPath, path string, binaryFileName string) []installFolderPath {
	var (
		binaryFilePath = gfile.Join(path, binaryFileName)
		writable       = gfile.IsWritable(path)
		installed      = isInstalled(binaryFilePath)
	)
	if !writable && !installed {
		return folderPaths
	}
	return append(
		folderPaths,
		installFolderPath{
			path:           path,
			writable:       writable,
			binaryFilePath: binaryFilePath,
			installed:      installed,
		})
}

// Check if this gf binary path exists.
func isInstalled(path string) bool {
	return gfile.Exists(path)
}
