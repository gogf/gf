// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

// Home returns absolute path of current user's home directory.
// The optional parameter <names> specifies the its sub-folders/sub-files,
// which will be joined with current system separator and returned with the path.
func Home(names ...string) (string, error) {
	path, err := getHomePath()
	if err != nil {
		return "", err
	}
	for _, name := range names {
		path += Separator + name
	}
	return path, nil
}

// getHomePath returns absolute path of current user's home directory.
func getHomePath() (string, error) {
	u, err := user.Current()
	if nil == err {
		return u.HomeDir, nil
	}
	if "windows" == runtime.GOOS {
		return homeWindows()
	}
	return homeUnix()
}

// homeUnix retrieves and returns the home on unix system.
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

// homeWindows retrieves and returns the home on windows system.
func homeWindows() (string, error) {
	var (
		drive = os.Getenv("HOMEDRIVE")
		path  = os.Getenv("HOMEPATH")
		home  = drive + path
	)
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}
