// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsmtp provides a simple SMTP client to access remote mail server.
//
// Eg:
// s := smtp.New("smtp.exmail.qq.com:25", "notify@a.com", "password")
// glog.Println(s.SendMail("notify@a.com", "ulric@b.com;rain@c.com", "subject", "body, <font color=red>red</font>"))
package gsmtp

import (
	"encoding/base64"
	"fmt"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"net/smtp"
	"strings"
)

// SMTP is the structure for smtp connection.
type SMTP struct {
	Address  string
	Username string
	Password string
}

// New creates and returns a new SMTP object.
func New(address, username, password string) *SMTP {
	return &SMTP{
		Address:  address,
		Username: username,
		Password: password,
	}
}

var (
	// contentEncoding is the BASE64 encoding object for mail content.
	contentEncoding = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
)

// SendMail connects to the server at addr, switches to TLS if
// possible, authenticates with the optional mechanism an if possible,
// and then sends an email from address <from>, to addresses <to>, with
// message msg.
//
// The parameter <contentType> specifies the content type of the mail, eg: html.
func (s *SMTP) SendMail(from, tos, subject, body string, contentType ...string) error {
	var (
		server  = ""
		address = ""
		hp      = strings.Split(s.Address, ":")
	)
	if s.Address == "" || len(hp) > 2 {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"server address is either empty or incorrect: %s",
			s.Address,
		)
	} else if len(hp) == 1 {
		server = s.Address
		address = server + ":25"
	} else if len(hp) == 2 {
		if (hp[0] == "") || (hp[1] == "") {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				"server address is either empty or incorrect: %s",
				s.Address,
			)
		}
		server = hp[0]
		address = s.Address
	}
	var (
		tosArr []string
		arr    = strings.Split(tos, ";")
	)
	for _, to := range arr {
		// TODO: replace with regex
		if strings.Contains(to, "@") {
			tosArr = append(tosArr, to)
		}
	}
	if len(tosArr) == 0 {
		return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid parameter "tos": %s`, tos)
	}

	if !strings.Contains(from, "@") {
		return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid parameter "from": %s`, from)
	}

	header := map[string]string{
		"From":                      from,
		"To":                        strings.Join(tosArr, ";"),
		"Subject":                   fmt.Sprintf("=?UTF-8?B?%s?=", contentEncoding.EncodeToString([]byte(subject))),
		"MIME-Version":              "1.0",
		"Content-Type":              "text/plain; charset=UTF-8",
		"Content-Transfer-Encoding": "base64",
	}
	if len(contentType) > 0 && contentType[0] == "html" {
		header["Content-Type"] = "text/html; charset=UTF-8"
	}
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + contentEncoding.EncodeToString([]byte(body))
	return smtp.SendMail(
		address,
		smtp.PlainAuth("", s.Username, s.Password, server),
		from,
		tosArr,
		[]byte(message),
	)
}
