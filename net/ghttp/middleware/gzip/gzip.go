package gzip

import (
	"bytes"
	"compress/gzip"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/net/ghttp"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

var defaultCompressOptions = &CompressOptions{
	ExcludeExt: nil,
	UseGzip:    func(r *ghttp.Request) bool { return true },
}

type gzipHandler struct {
	options CompressOptions
	gzPool  sync.Pool
}

type CompressOptions struct {
	ExcludeExt *gset.Set
	UseGzip    func(r *ghttp.Request) bool
}

func Compress(level int, options *CompressOptions) func(r *ghttp.Request) {
	var gzPool sync.Pool
	gzPool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, level)
		if err != nil {
			panic(err)
		}
		return gz
	}

	if options == nil {
		options = defaultCompressOptions
	}
	handler := gzipHandler{gzPool: gzPool, options: *options}
	return handler.handler
}

func (g *gzipHandler) handler(r *ghttp.Request) {
	if !g.needCompress(r) {
		r.Middleware.Next()
		return
	}
	var buffer bytes.Buffer
	gz := g.gzPool.Get().(*gzip.Writer)
	defer g.gzPool.Put(gz)
	defer gz.Reset(ioutil.Discard)
	gz.Reset(&buffer)
	defer func() {
		gz.Close()
	}()
	r.Middleware.Next()
	gz.Write(r.Response.Buffer())
	gz.Flush()
	r.Response.SetBuffer(buffer.Bytes())
	r.Response.Header().Add("Content-Encoding", "gzip")
	r.Response.Header().Add("Vary", "Accept-Encoding")
}

func (g *gzipHandler) needCompress(r *ghttp.Request) bool {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") ||
		strings.Contains(r.Header.Get("Connection"), "Upgrade") {
		return false
	}

	if g.options.UseGzip != nil && !g.options.UseGzip(r) {
		return false
	}
	if g.options.ExcludeExt != nil {
		extension := filepath.Ext(r.URL.Path)
		if g.options.ExcludeExt.Contains(extension) {
			return false
		}
	}

	return true
}

func Decompress(r *ghttp.Request) {
	if r.Request.Body == nil {
		r.Middleware.Next()
		return
	}
	if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		r.Middleware.Next()
		return
	}
	raw, err := gzip.NewReader(r.Request.Body)
	if err != nil {
		r.Response.WriteStatus(500)
		return
	}
	r.Request.Header.Del("Content-Encoding")
	r.Request.Header.Del("Content-Length")
	r.Body = raw
	r.Middleware.Next()
}
