// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import "github.com/gogf/gf/g/container/gmap"

var (
    // Customized configuration content.
    configs = gmap.NewStringStringMap()
)

// SetContent sets customized configuration content for specified <file>.
// The <file> is unnecessary param, default is DEFAULT_CONFIG_FILE.
func SetContent(content string, file...string) {
    name := DEFAULT_CONFIG_FILE
    if len(file) > 0 {
        name = file[0]
    }
    configs.Set(name, content)
}

// GetContent returns customized configuration content for specified <file>.
// The <file> is unnecessary param, default is DEFAULT_CONFIG_FILE.
func GetContent(file...string) string {
    name := DEFAULT_CONFIG_FILE
    if len(file) > 0 {
        name = file[0]
    }
    return configs.Get(name)
}

// ClearContent removes all global configuration contents.
func ClearContent() {
    configs.Clear()
}