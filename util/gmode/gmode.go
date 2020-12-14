// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmode provides release mode management for project.
//
// It uses string to mark the mode instead of integer, which is convenient for configuration.
package gmode

import (
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
)

const (
	NOT_SET   = "not-set"
	DEVELOP   = "develop"
	TESTING   = "testing"
	STAGING   = "staging"
	PRODUCT   = "product"
	cmdEnvKey = "gf.gmode"
)

var (
	currentMode = NOT_SET
)

// Set sets the mode for current application.
func Set(mode string) {
	currentMode = mode
}

// SetDevelop sets current mode DEVELOP for current application.
func SetDevelop() {
	Set(DEVELOP)
}

// SetTesting sets current mode TESTING for current application.
func SetTesting() {
	Set(TESTING)
}

// SetStaging sets current mode STAGING for current application.
func SetStaging() {
	Set(STAGING)
}

// SetProduct sets current mode PRODUCT for current application.
func SetProduct() {
	Set(PRODUCT)
}

// Mode returns current application mode set.
func Mode() string {
	// If current mode is not set, do this auto check.
	if currentMode == NOT_SET {
		if v := gcmd.GetWithEnv(cmdEnvKey).String(); v != "" {
			// Mode configured from command argument of environment.
			currentMode = v
		} else {
			// If there are source codes found, it's in develop mode, or else in product mode.
			if gfile.Exists(gdebug.CallerFilePath()) {
				currentMode = DEVELOP
			} else {
				currentMode = PRODUCT
			}
		}
	}
	return currentMode
}

// IsDevelop checks and returns whether current application is running in DEVELOP mode.
func IsDevelop() bool {
	return Mode() == DEVELOP
}

// IsTesting checks and returns whether current application is running in TESTING mode.
func IsTesting() bool {
	return Mode() == TESTING
}

// IsStaging checks and returns whether current application is running in STAGING mode.
func IsStaging() bool {
	return Mode() == STAGING
}

// IsProduct checks and returns whether current application is running in PRODUCT mode.
func IsProduct() bool {
	return Mode() == PRODUCT
}
