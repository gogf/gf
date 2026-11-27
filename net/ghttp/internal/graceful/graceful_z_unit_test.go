// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file covers raw listener matching by listening address.

package graceful

import (
	"io"
	"net"
	"strconv"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_New_RawListener(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		t.AssertNil(err)
		defer func() { t.AssertNil(listener.Close()) }()

		var (
			port   = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
			config = ServerConfig{Listeners: []net.Listener{listener}}
		)
		t.Assert(New(":"+port, 0, io.Discard, config).rawListener == listener, true)
		t.Assert(New("[::]:"+port, 0, io.Discard, config).rawListener == listener, true)
		t.AssertNil(New("bad::"+port, 0, io.Discard, config).rawListener)
	})
}
