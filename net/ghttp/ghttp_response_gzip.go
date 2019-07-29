// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

// 默认的gzip压缩文件类型
var defaultGzipContentTypes = []string{
	"application/atom+xml",
	"application/font-sfnt",
	"application/javascript",
	"application/json",
	"application/ld+json",
	"application/manifest+json",
	"application/rdf+xml",
	"application/rss+xml",
	"application/schema+json",
	"application/vnd.geo+json",
	"application/vnd.ms-fontobject",
	"application/x-font-ttf",
	"application/x-javascript",
	"application/x-web-app-manifest+json",
	"application/xhtml+xml",
	"application/xml",
	"font/eot",
	"font/opentype",
	"image/bmp",
	"image/svg+xml",
	"image/vnd.microsoft.icon",
	"image/x-icon",
	"text/cache-manifest",
	"text/css",
	"text/html",
	"text/javascript",
	"text/plain",
	"text/vcard",
	"text/vnd.rim.location.xloc",
	"text/vtt",
	"text/x-component",
	"text/x-cross-domain-policy",
	"text/xml",
}

//// 返回内容gzip检查处理
//func (r *Response) handleGzip() {
//    // 如果客户端支持gzip压缩，并且服务端设置开启gzip压缩特性，那么执行压缩
//    encoding := r.request.Header.Get("Accept-Encoding")
//    if encoding != "" && strings.Contains(encoding, "gzip") {
//        mimeType := ""
//        ext := gfile.Ext(r.request.URL.Path)
//        if ext != "" {
//            mimeType = strings.Split(mime.TypeByExtension(ext), ";")[0]
//        }
//        if mimeType == "" {
//            contentType := r.Header().Get("Content-Type")
//            if contentType != "" {
//                mimeType = strings.Split(contentType, ";")[0]
//            }
//        }
//
//        if _, ok := r.Server.gzipMimesMap[mimeType]; ok {
//            r.SetBuffer(gcompress.Gzip(r.buffer))
//            r.Header().Set("Content-Length",   gconv.String(len(r.buffer)))
//            r.Header().Set("Content-Encoding", "gzip")
//        }
//    }
//}
