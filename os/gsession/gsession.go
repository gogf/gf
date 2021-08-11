// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsession implements manager and storage features for sessions.
package gsession

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/util/guid"
)

var (
	ErrorDisabled = gerror.NewOption(gerror.Option{
		Text: "this feature is disabled in this storage",
		Code: gerror.CodeNotSupported,
	})
)

// NewSessionId creates and returns a new and unique session id string,
// which is in 36 bytes.
func NewSessionId() string {
	return guid.S()
}
