// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmvc provides basic object classes for MVC.
// Deprecated, no longer suggested.
package gmvc

import (
	"github.com/gogf/gf/net/ghttp"
)

// Controller is used for controller register of ghttp.Server.
// Deprecated, no longer suggested.
type Controller struct {
	Request  *ghttp.Request
	Response *ghttp.Response
	Server   *ghttp.Server
	Cookie   *ghttp.Cookie
	Session  *ghttp.Session
	View     *View
}

// Init is the callback function for each request initialization.
func (c *Controller) Init(r *ghttp.Request) {
	c.Request = r
	c.Response = r.Response
	c.Server = r.Server
	c.View = NewView(r.Response)
	c.Cookie = r.Cookie
	c.Session = r.Session
}

// Shut is the callback function for each request close.
func (c *Controller) Shut() {

}

// Exit equals to function Request.Exit().
func (c *Controller) Exit() {
	c.Request.Exit()
}
