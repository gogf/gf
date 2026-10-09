// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// HTTP header names used by the ghttp package.
//
// Note that the constant values are in canonical format as defined by
// net/textproto.CanonicalMIMEHeaderKey. Using canonical names avoids the
// re-canonicalization that http.Header.Get/Set/Del performs for any other format.
const (
	// Request headers.
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderAuthorization   = "Authorization"
	HeaderContentType     = "Content-Type"
	HeaderOrigin          = "Origin"
	HeaderReferer         = "Referer"
	HeaderXForwardedFor   = "X-Forwarded-For"
	HeaderXForwardedProto = "X-Forwarded-Proto"
	HeaderXRealIp         = "X-Real-Ip"
	HeaderXRequestedWith  = "X-Requested-With"

	// Response headers.
	HeaderAcceptRanges       = "Accept-Ranges"
	HeaderContentDisposition = "Content-Disposition"
	HeaderContentEncoding    = "Content-Encoding"
	HeaderContentLength      = "Content-Length"
	HeaderLocation           = "Location"
	HeaderServer             = "Server"
	HeaderTraceId            = "Trace-Id"
	HeaderWwwAuthenticate    = "Www-Authenticate"

	// Response headers for the CORS feature.
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"

	// HeaderXUrlPath is used for custom route handler, which does not change URL.Path.
	HeaderXUrlPath = "X-Url-Path"
)

// Private HTTP header names for internal usage.
// Note that the values are canonicalized as the above constants are, so there is no need
// in using the original format like `HTTP_CLIENT_IP` for these non-standard header names.
const (
	// Non-standard header names for client ip retrieving.
	headerHttpClientIp      = "Http_client_ip"
	headerHttpXForwardedFor = "Http_x_forwarded_for"
	headerProxyClientIp     = "Proxy-Client-Ip"
	headerWlProxyClientIp   = "Wl-Proxy-Client-Ip"

	// Header names referenced by the default CORS allowed headers.
	headerAccept     = "Accept"
	headerCookie     = "Cookie"
	headerUserAgent  = "User-Agent"
	headerXAuthToken = "X-Auth-Token"
)
