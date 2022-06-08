// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"compress/gzip"
	"github.com/gogf/gf/v2/internal/intlog"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

type (
	ExcludedExts        map[string]bool
	ExcludedPaths       []string
	ExcludedPathRegexps []*regexp.Regexp
)

type Options struct {
	ExcludedExts        ExcludedExts
	ExcludedPaths       ExcludedPaths
	ExcludedPathRegexps ExcludedPathRegexps
}

type Option func(*Options)

func WithExcludedExts(exts []string) Option {
	return func(o *Options) {
		o.ExcludedExts = NewExcludedExts(exts)
	}
}

func NewExcludedExts(exts []string) ExcludedExts {
	res := make(ExcludedExts)
	for _, ext := range exts {
		res[ext] = true
	}
	return res
}

func (e ExcludedExts) Contains(ext string) bool {
	_, ok := e[ext]
	return ok
}

func WithExcludedPaths(paths []string) Option {
	return func(o *Options) {
		o.ExcludedPaths = NewExcludedPaths(paths)
	}
}

func NewExcludedPaths(paths []string) ExcludedPaths {
	return paths
}

func (e ExcludedPaths) Contains(URI string) bool {
	for _, path := range e {
		if strings.HasPrefix(URI, path) {
			return true
		}
	}
	return false
}

func WithExcludedPathRegexps(regexps []string) Option {
	return func(o *Options) {
		o.ExcludedPathRegexps = NewExcludedPathRegexps(regexps)
	}
}

func NewExcludedPathRegexps(regexps []string) ExcludedPathRegexps {
	res := make([]*regexp.Regexp, 0, len(regexps))
	for _, reg := range regexps {
		res = append(res, regexp.MustCompile(reg))
	}
	return res
}

func (e ExcludedPathRegexps) Contains(URI string) bool {
	for _, reg := range e {
		if reg.MatchString(URI) {
			return true
		}
	}
	return false
}

var (
	DefaultOptions = &Options{
		ExcludedExts: NewExcludedExts([]string{
			".png", ".jpg", "jpeg", ".gif",
		}),
	}
)

const (
	GzipBestCompression    = gzip.BestCompression
	GzipBestSpeed          = gzip.BestSpeed
	GzipDefaultCompression = gzip.DefaultCompression
	GzipNoCompression      = gzip.NoCompression
)

type gzipHandler struct {
	level int
	*Options
}

func newGzipHandler(level int, options ...Option) *gzipHandler {
	handler := &gzipHandler{
		level:   level,
		Options: DefaultOptions,
	}

	for _, o := range options {
		o(handler.Options)
	}

	return handler
}

func (g *gzipHandler) Handle(r *Request) {
	if g.shouldCompress(r) {
		gzipBuffer := bytes.NewBuffer(nil)
		gz, err := gzip.NewWriterLevel(gzipBuffer, g.level)
		if err != nil {
			intlog.Errorf(r.Context(), `%+v`, err)
		} else {
			r.Response.Header().Set("Content-Encoding", "gzip")
			r.Response.Header().Set("Vary", "Accept-Encoding")
			r.Response.Writer = &gzipWriter{r.Response.ResponseWriter, gz, gzipBuffer}
		}
	}

	r.Middleware.Next()
}

func (g *gzipHandler) shouldCompress(r *Request) bool {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") ||
		strings.Contains(r.Header.Get("Connection"), "Upgrade") ||
		strings.Contains(r.Header.Get("Accept"), "text/event-stream") {
		return false
	}

	ext := filepath.Ext(r.URL.Path)
	if g.ExcludedExts.Contains(ext) {
		return false
	}

	if g.ExcludedPaths.Contains(r.URL.Path) {
		return false
	}

	if g.ExcludedPathRegexps.Contains(r.URL.Path) {
		return false
	}

	return true
}

func MiddlewareGzip(level int, options ...Option) HandlerFunc {
	return newGzipHandler(level, options...).Handle
}

type gzipWriter struct {
	*ResponseWriter
	writer     *gzip.Writer
	gzipBuffer *bytes.Buffer
}

func (g *gzipWriter) Flush() {
	g.ResponseWriter.Header().Set("Content-Type", http.DetectContentType(g.buffer.Bytes()))
	g.writer.Write(g.buffer.Bytes())
	g.writer.Close()
	g.buffer.Reset()

	g.gzipBuffer.WriteTo(g.buffer)
	g.ResponseWriter.Flush()
}
