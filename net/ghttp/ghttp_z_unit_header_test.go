// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests for the canonical spelling of ghttp header name constants.

package ghttp

import (
	"net/textproto"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// Test_HeaderNames_Canonical checks that header name constants are already in the form
// stored by net/http. A non-canonical value makes Header.Get and Header.Set allocate a new key.
func Test_HeaderNames_Canonical(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var headerNames = []string{
			HeaderAcceptEncoding,
			HeaderAuthorization,
			HeaderContentType,
			HeaderOrigin,
			HeaderReferer,
			HeaderXForwardedFor,
			HeaderXForwardedProto,
			HeaderXRealIp,
			HeaderXRequestedWith,
			HeaderAcceptRanges,
			HeaderContentDisposition,
			HeaderContentEncoding,
			HeaderContentLength,
			HeaderLocation,
			HeaderServer,
			HeaderTraceId,
			HeaderWwwAuthenticate,
			HeaderAccessControlAllowCredentials,
			HeaderAccessControlAllowHeaders,
			HeaderAccessControlAllowMethods,
			HeaderAccessControlAllowOrigin,
			HeaderAccessControlExposeHeaders,
			HeaderAccessControlMaxAge,
			HeaderAccessControlRequestHeaders,
			HeaderAccessControlRequestMethod,
			HeaderXUrlPath,
			headerHttpClientIp,
			headerHttpXForwardedFor,
			headerProxyClientIp,
			headerWlProxyClientIp,
			headerAccept,
			headerCookie,
			headerUserAgent,
			headerXAuthToken,
		}
		for _, headerName := range headerNames {
			t.Assert(headerName, textproto.CanonicalMIMEHeaderKey(headerName))
		}
	})
}
