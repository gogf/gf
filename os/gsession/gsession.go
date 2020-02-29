// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsession implements manager and storage features for sessions.
package gsession

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/grand"
)

var (
	ErrorDisabled = errors.New("this feature is disabled in this storage")
)

// NewSessionId creates and returns a new and unique session id string,
// the length of which is 18 bytes.
func NewSessionId() string {
	return strings.ToUpper(strconv.FormatInt(gtime.TimestampNano(), 36) + grand.S(6))
}
