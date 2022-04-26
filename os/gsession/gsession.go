// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsession implements manager and storage features for sessions.
package gsession

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/guid"
)

var (
	// ErrorDisabled is used for marking certain interface function not used.
	ErrorDisabled = gerror.NewOption(gerror.Option{
		Text: "this feature is disabled in this storage",
		Code: gcode.CodeNotSupported,
	})
)

// NewSessionId creates and returns a new and unique session id string,
// which is in 32 bytes.
func NewSessionId() string {
	return guid.S()
}
