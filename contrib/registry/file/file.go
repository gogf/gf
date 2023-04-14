// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package file implements service Registry and Discovery using file.
package file

import (
	"time"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gfile"
)

var (
	_ gsvc.Registry = &Registry{}
)

const (
	updateAtKey           = "UpdateAt"
	serviceTTL            = 20 * time.Second
	serviceUpdateInterval = 10 * time.Second
	defaultSeparator      = "#"
)

// Registry implements interface Registry using file.
// This implement is usually for testing only.
type Registry struct {
	path string // Local storing folder path for Services.
}

// New creates and returns a gsvc.Registry implements using file.
func New(path string) gsvc.Registry {
	if !gfile.Exists(path) {
		_ = gfile.Mkdir(path)
	}
	return &Registry{
		path: path,
	}
}
