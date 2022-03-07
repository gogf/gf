// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"net/http"
	"time"
)

// SetCookieMaxAge sets the CookieMaxAge for server.
func (s *Server) SetCookieMaxAge(ttl time.Duration) {
	s.config.CookieMaxAge = ttl
}

// SetCookiePath sets the CookiePath for server.
func (s *Server) SetCookiePath(path string) {
	s.config.CookiePath = path
}

// SetCookieDomain sets the CookieDomain for server.
func (s *Server) SetCookieDomain(domain string) {
	s.config.CookieDomain = domain
}

// GetCookieMaxAge returns the CookieMaxAge of server.
func (s *Server) GetCookieMaxAge() time.Duration {
	return s.config.CookieMaxAge
}

// GetCookiePath returns the CookiePath of server.
func (s *Server) GetCookiePath() string {
	return s.config.CookiePath
}

// GetCookieDomain returns CookieDomain of server.
func (s *Server) GetCookieDomain() string {
	return s.config.CookieDomain
}

// GetCookieSameSite return CookieSameSite of server.
func (s *Server) GetCookieSameSite() http.SameSite {
	switch s.config.CookieSameSite {
	case "lax":
		return http.SameSiteLaxMode
	case "none":
		return http.SameSiteNoneMode
	case "strict":
		return http.SameSiteStrictMode
	default:
		return http.SameSiteDefaultMode
	}
}

func (s *Server) GetCookieSecure() bool {
	return s.config.CookieSecure
}

func (s *Server) GetCookieHttpOnly() bool {
	return s.config.CookieHttpOnly
}
