// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsmtp provides a SMTP client to access remote mail server.
//
// eg:
// s := smtp.New("smtp.exmail.qq.com:25", "notify@a.com", "password")
// glog.Println(s.SendMail("notify@a.com", "ulric@b.com;rain@c.com", "subject", "body, <font color=red>red</font>"))
package gsmtp

import (
    "encoding/base64"
    "fmt"
    "net/smtp"
    "strings"
)

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
    if s.Address == "" {
        return fmt.Errorf("address is necessary")
    }

    hp := strings.Split(s.Address, ":")
    if len(hp) != 2 {
        return fmt.Errorf("address format error")
    }

    arr := strings.Split(tos, ";")
    count := len(arr)
    safeArr := make([]string, 0, count)
    for i := 0; i < count; i++ {
        if arr[i] == "" {
            continue
        }
        safeArr = append(safeArr, arr[i])
    }

    if len(safeArr) == 0 {
        return fmt.Errorf("tos invalid")
    }

    tos  = strings.Join(safeArr, ";")
    b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

    header                := make(map[string]string)
    header["From"]         = from
    header["To"]           = tos
    header["Subject"]      = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(subject)))
    header["MIME-Version"] = "1.0"

    ct := "text/plain; charset=UTF-8"
    if len(contentType) > 0 && contentType[0] == "html" {
        ct = "text/html; charset=UTF-8"
    }

    header["Content-Type"]              = ct
    header["Content-Transfer-Encoding"] = "base64"

    message := ""
    for k, v := range header {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + b64.EncodeToString([]byte(body))

    auth := smtp.PlainAuth("", s.Username, s.Password, hp[0])
    return smtp.SendMail(s.Address, auth, from, strings.Split(tos, ";"), []byte(message))
}