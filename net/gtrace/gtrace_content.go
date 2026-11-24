// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtrace

import (
	"mime"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/text/gstr"
)

// SafeContentForHttp cuts and returns given content by `MaxContentLogSize`.
// It appends string `...` to the tail of the result if the content size is greater than `MaxContentLogSize`.
func SafeContentForHttp(data []byte, header http.Header) (string, error) {
	var err error
	if isGzipEncoded(header) {
		if data, err = gcompress.UnGzip(data); err != nil {
			return string(data), err
		}
	}

	return SafeContent(data), nil
}

// SafeContent cuts and returns given content by `MaxContentLogSize`.
// It appends string `...` to the tail of the result if the content size is greater than `MaxContentLogSize`.
func SafeContent(data []byte) string {
	content := string(data)
	if gstr.LenRune(content) > MaxContentLogSize() {
		content = gstr.StrLimitRune(content, MaxContentLogSize(), "...")
	}

	return content
}

// isGzipEncoded returns whether the content is gzip-encoded according to the
// `Content-Encoding` header. Content codings are case-insensitive and may
// carry parameters, eg. "GZip" or "gzip; level=1".
func isGzipEncoded(header http.Header) bool {
	for part := range strings.SplitSeq(header.Get("Content-Encoding"), ",") {
		mediaType, _, err := mime.ParseMediaType(strings.TrimSpace(part))
		if err != nil {
			continue
		}
		if strings.EqualFold(mediaType, "gzip") {
			return true
		}
	}
	return false
}
