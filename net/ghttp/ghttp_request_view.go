// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "github.com/gogf/gf/os/gview"

// SetView sets template view engine object for this request.
func (r *Request) SetView(view *gview.View) {
	r.viewObject = view
}

// GetView returns the template view engine object for this request.
func (r *Request) GetView() *gview.View {
	view := r.viewObject
	if view == nil {
		view = r.Server.config.View
	}
	if view == nil {
		view = gview.Instance()
	}
	return view
}

// Assigns binds multiple template variables to current request.
func (r *Request) Assigns(data gview.Params) {
	if r.viewParams == nil {
		r.viewParams = make(gview.Params, len(data))
	}
	for k, v := range data {
		r.viewParams[k] = v
	}
}

// Assign binds a template variable to current request.
func (r *Request) Assign(key string, value interface{}) {
	if r.viewParams == nil {
		r.viewParams = make(gview.Params)
	}
	r.viewParams[key] = value
}
