// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests shutdown of prepared listeners before net/http starts serving.

package graceful

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// TestPreparedTLSCertificateRetry leaves caller certificates untouched until loading succeeds.
func TestPreparedTLSCertificateRetry(t *testing.T) {
	for _, embedded := range []bool{false, true} {
		name := "file"
		if embedded {
			name = "embedded"
		}
		t.Run(name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					path      = filepath.Join("..", "..", "testdata", "https", "files")
					certFile  = filepath.Join(path, "server.crt")
					keyFile   = filepath.Join(path, "server.key")
					tlsConfig = &tls.Config{}
					server    = New("127.0.0.1:0", 0, io.Discard, ServerConfig{Logger: glog.NewWithWriter(io.Discard)})
				)
				defer server.Close(context.Background())
				if embedded {
					prefix := "tls-retry-" + guid.S()
					data, err := gres.PackWithOption(path, gres.Option{Prefix: prefix})
					t.AssertNil(err)
					t.AssertNil(gres.Add(string(data)))
					certFile = prefix + "/server.crt"
					keyFile = prefix + "/server.key"
					t.Assert(gres.Contains(certFile), true)
				}
				t.AssertNE(server.CreateListenerTLS(certFile, keyFile+".missing", tlsConfig), nil)
				t.Assert(len(tlsConfig.Certificates), 0)
				t.AssertNil(server.CreateListenerTLS(certFile, keyFile, tlsConfig))
				t.Assert(len(tlsConfig.Certificates), 1)
				t.AssertNE(tlsConfig.Certificates[0].PrivateKey, nil)
			})
		})
	}
}

// preparedCloseFailureListener reports an unexpected failure after releasing its socket.
type preparedCloseFailureListener struct {
	net.Listener
}

// Close releases the socket and returns a failure that shutdown must log.
func (l *preparedCloseFailureListener) Close() error {
	if err := l.Listener.Close(); err != nil {
		return err
	}
	return gerror.New("prepared listener cleanup failed")
}

// TestPreparedListenerCloseFailure verifies unexpected cleanup errors remain visible.
func TestPreparedListenerCloseFailure(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var output bytes.Buffer
		server := New("127.0.0.1:0", 0, io.Discard, ServerConfig{Logger: glog.NewWithWriter(&output)})
		t.AssertNil(server.CreateListener())
		server.setRawListener(&preparedCloseFailureListener{Listener: server.getRawListener()})
		server.Close(context.Background())
		t.Assert(strings.Contains(output.String(), "prepared listener cleanup failed"), true)
	})
}

// TestPreparedListenerShutdown releases a listener even when Serve has not run yet.
func TestPreparedListenerShutdown(t *testing.T) {
	for _, forceful := range []bool{false, true} {
		name := "graceful"
		if forceful {
			name = "forceful"
		}
		t.Run(name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				logger := glog.NewWithWriter(io.Discard)
				server := New("127.0.0.1:0", 0, io.Discard, ServerConfig{Logger: logger, GracefulShutdownTimeout: 1})
				t.AssertNil(server.CreateListener())
				address := server.GetListenedAddress()
				defer server.Close(context.Background())
				if forceful {
					server.Close(context.Background())
				} else {
					server.Shutdown(context.Background())
				}
				listener, err := net.Listen("tcp", address)
				t.AssertNil(err)
				t.AssertNil(listener.Close())
			})
		})
	}
}
