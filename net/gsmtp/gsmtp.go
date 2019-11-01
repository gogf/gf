// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsmtp provides a SMTP client to access remote mail server.
//
// Eg:
// s := smtp.New("smtp.exmail.qq.com:25", "notify@a.com", "password")
// glog.Println(s.SendMail("notify@a.com", "ulric@b.com;rain@c.com", "subject", "body, <font color=red>red</font>"))
package gsmtp

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
)

// SMTP is the structure for smtp connection
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

// SendMail connects to the server at addr, switches to TLS if
// possible, authenticates with the optional mechanism a if possible,
// and then sends an email from address from, to addresses to, with
// message msg.
func (s *SMTP) SendMail(from, tos, subject, body string, contentType ...string) error {
	server := ""
	address := ""

	hp := strings.Split(s.Address, ":")
	if (s.Address == "") || (len(hp) > 2) {
		return fmt.Errorf("Server address is either empty or incorrect: %s", s.Address)
	} else if len(hp) == 1 {
		server = s.Address
		address = server + ":25"
	} else if len(hp) == 2 {
		if (hp[0] == "") || (hp[1] == "") {
			return fmt.Errorf("Server address is either empty or incorrect: %s", s.Address)
		}
		server = hp[0]
		address = s.Address
	}

	tosArr := []string{}
	arr := strings.Split(tos, ";")
	for _, to := range arr {
		// TODO: replace with regex
		if strings.Contains(to, "@") {
			tosArr = append(tosArr, to)
		}
	}

	if len(tosArr) == 0 {
		return fmt.Errorf("tos if invalid: %s", tos)
	}

	if !strings.Contains(from, "@") {
		return fmt.Errorf("from is invalid: %s", from)
	}

	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	header := make(map[string]string)
	header["From"] = from
	header["To"] = strings.Join(tosArr, ";")
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"

	ct := "text/plain; charset=UTF-8"
	if len(contentType) > 0 && contentType[0] == "html" {
		ct = "text/html; charset=UTF-8"
	}

	header["Content-Type"] = ct
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString([]byte(body))

	auth := smtp.PlainAuth("", s.Username, s.Password, server)
	return smtp.SendMail(address, auth, from, tosArr, []byte(message))
}
