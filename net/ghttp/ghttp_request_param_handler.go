// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "github.com/gogf/gf/v2/util/gmeta"

// GetHandlerResponse retrieves and returns the handler response object and its error.
func (r *Request) GetHandlerResponse() interface{} {
	return r.handlerResponse
}

// GetServeHandler retrieves and returns the user defined handler used to serve this request.
func (r *Request) GetServeHandler() *HandlerItemParsed {
	return r.serveHandler
}

// GetMetaTag retrieves and returns the metadata value associated with the given key from the request struct.
// The meta value is from struct tags from g.Meta/gmeta.Meta type.
// For example:
//
//	type GetMetaTagReq struct {
//	    g.Meta `path:"/test" method:"post" summary:"meta_tag" tags:"meta"`
//	    // ...
//	}
//
// r.GetServeHandler().GetMetaTag("summary") // returns "meta_tag"
// r.GetServeHandler().GetMetaTag("method")  // returns "post"
func (h *HandlerItemParsed) GetMetaTag(key string) string {
	if h == nil || h.Handler == nil {
		return ""
	}
	metaValue := gmeta.Get(h.Handler.Info.Type.In(1), key)
	if metaValue != nil {
		return metaValue.String()
	}
	return ""
}
