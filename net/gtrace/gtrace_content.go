// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtrace

import (
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/encoding/gcompress"

	"github.com/gogf/gf/v2/text/gstr"
)

// SafeContentForHttp cuts and returns given content by `MaxContentLogSize`.
// It appends string `...` to the tail of the result if the content size is greater than `MaxContentLogSize`.
func SafeContentForHttp(data []byte, header http.Header) (string, error) {
	var err error
	if gzipAccepted(header) {
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

// gzipAccepted returns whether the client will accept gzip-encoded content.
func gzipAccepted(header http.Header) bool {
	a := header.Get("Content-Encoding")
	parts := strings.Split(a, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "gzip" || strings.HasPrefix(part, "gzip;") {
			return true
		}
	}

	return false
}
