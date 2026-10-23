// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"time"

	"github.com/gogf/gf/v2/os/gsession"
)

// SetSessionMaxAge sets the SessionMaxAge for server.
func (s *Server) SetSessionMaxAge(ttl time.Duration) {
	s.config.SessionMaxAge = ttl
}

// SetSessionIdName sets the SessionIdName for server.
func (s *Server) SetSessionIdName(name string) {
	s.config.SessionIdName = name
}

// SetSessionIdFunc sets the custom session id creating function for server.
// It is called lazily: only when a request without a session id from cookie or header
// actually uses the session, and again when the session id is regenerated.
//
// The function can only use request information (header, cookie, client ip, url, context).
// It must not:
//   - call any method of r.Session, which causes infinite recursion in session initialization;
//   - create the id from session data, as the session is not created yet at that time;
//   - return the same id on each call, as the id must be new for every generation;
//   - return an empty id, which is treated as an error and results in a 500 response.
//
// It must be set before Start or Run is called on this server, and cannot be changed after that.
// When using file storage, the returned id is used as the session file name, so it must not
// contain path separators, ".." or other invalid filename characters; only letters, digits,
// "_" and "-" are recommended.
//
// Passing nil restores the built-in generator.
func (s *Server) SetSessionIdFunc(f func(r *Request) string) {
	s.sessionIdFunc = f
}

// SetSessionStorage sets the SessionStorage for server.
func (s *Server) SetSessionStorage(storage gsession.Storage) {
	s.config.SessionStorage = storage
}

// SetSessionCookieOutput sets the SetSessionCookieOutput for server.
func (s *Server) SetSessionCookieOutput(enabled bool) {
	s.config.SessionCookieOutput = enabled
}

// SetSessionCookieMaxAge sets the SessionCookieMaxAge for server.
func (s *Server) SetSessionCookieMaxAge(maxAge time.Duration) {
	s.config.SessionCookieMaxAge = maxAge
}

// GetSessionMaxAge returns the SessionMaxAge of server.
func (s *Server) GetSessionMaxAge() time.Duration {
	return s.config.SessionMaxAge
}

// GetSessionIdName returns the SessionIdName of server.
func (s *Server) GetSessionIdName() string {
	return s.config.SessionIdName
}

// GetSessionCookieMaxAge returns the SessionCookieMaxAge of server.
func (s *Server) GetSessionCookieMaxAge() time.Duration {
	return s.config.SessionCookieMaxAge
}
