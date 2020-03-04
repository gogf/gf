// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"net/http"
	"net/url"

	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

// CORSOptions is the options for CORS feature.
// See https://www.w3.org/TR/cors/ .
type CORSOptions struct {
	AllowDomain      []string // Used for allowing requests from custom domains
	AllowOrigin      string   // Access-Control-Allow-Origin
	AllowCredentials string   // Access-Control-Allow-Credentials
	ExposeHeaders    string   // Access-Control-Expose-Headers
	MaxAge           int      // Access-Control-Max-Age
	AllowMethods     string   // Access-Control-Allow-Methods
	AllowHeaders     string   // Access-Control-Allow-Headers
}

var (
	// defaultAllowHeaders is the default allowed headers for CORS.
	// It's defined as map for better header key searching performance.
	defaultAllowHeaders = map[string]struct{}{
		"Origin":           {},
		"Accept":           {},
		"Cookie":           {},
		"Authorization":    {},
		"X-Auth-Token":     {},
		"X-Requested-With": {},
		"Content-Type":     {},
	}
)

// DefaultCORSOptions returns the default CORS options,
// which allows any cross-domain request.
func (r *Response) DefaultCORSOptions() CORSOptions {
	options := CORSOptions{
		AllowOrigin:      "*",
		AllowMethods:     HTTP_METHODS,
		AllowCredentials: "true",
		MaxAge:           3628800,
	}
	// Allow all client's custom headers in default.
	if headers := r.Request.Header.Get("Access-Control-Request-Headers"); headers != "" {
		array := gstr.SplitAndTrim(headers, ",")
		for _, header := range array {
			if _, ok := defaultAllowHeaders[header]; !ok {
				options.AllowHeaders += header + ","
			}
		}
		for header, _ := range defaultAllowHeaders {
			if len(options.AllowHeaders) > 0 {
				options.AllowHeaders += ","
			}
			options.AllowHeaders += header
		}
	}
	// Allow all anywhere origin in default.
	if origin := r.Request.Header.Get("Origin"); origin != "" {
		options.AllowOrigin = origin
	} else if referer := r.Request.Referer(); referer != "" {
		if p := gstr.PosR(referer, "/", 6); p != -1 {
			options.AllowOrigin = referer[:p]
		} else {
			options.AllowOrigin = referer
		}
	}
	return options
}

// CORS sets custom CORS options.
// See https://www.w3.org/TR/cors/ .
func (r *Response) CORS(options CORSOptions) {
	if r.CORSAllowedOrigin(options) {
		r.Header().Set("Access-Control-Allow-Origin", options.AllowOrigin)
	}
	if options.AllowCredentials != "" {
		r.Header().Set("Access-Control-Allow-Credentials", options.AllowCredentials)
	}
	if options.ExposeHeaders != "" {
		r.Header().Set("Access-Control-Expose-Headers", options.ExposeHeaders)
	}
	if options.MaxAge != 0 {
		r.Header().Set("Access-Control-Max-Age", gconv.String(options.MaxAge))
	}
	if options.AllowMethods != "" {
		r.Header().Set("Access-Control-Allow-Methods", options.AllowMethods)
	}
	if options.AllowHeaders != "" {
		r.Header().Set("Access-Control-Allow-Headers", options.AllowHeaders)
	}
	// No continue service handling if it's OPTIONS request.
	if gstr.Equal(r.Request.Method, "OPTIONS") {
		if r.Status == 0 {
			r.Status = http.StatusOK
		}
		r.Request.ExitAll()
	}
}

// CORSAllowed checks whether the current request origin is allowed cross-domain.
func (r *Response) CORSAllowedOrigin(options CORSOptions) bool {
	if options.AllowDomain == nil {
		return true
	}
	origin := r.Request.Header.Get("Origin")
	if origin == "" {
		return true
	}
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	for _, v := range options.AllowDomain {
		if gstr.IsSubDomain(parsed.Host, v) {
			return true
		}
	}
	return false
}

// CORSDefault sets CORS with default CORS options,
// which allows any cross-domain request.
func (r *Response) CORSDefault() {
	r.CORS(r.DefaultCORSOptions())
}
