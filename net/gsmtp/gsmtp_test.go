// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsmtp_test

import (
	"strings"
	"testing"

	"github.com/gogf/gf/net/gsmtp"
)

func TestAddress(t *testing.T) {
	errMessage := "address is either empty or incorrect"

	errValues := []string{
		"",
		":",
		":25",
		"localhost:",
		"local.host:25:28",
	}

	for _, errValue := range errValues {
		smtpConnection := gsmtp.New(errValue, "smtpUser@smtp.exmail.qq.com", "smtpPassword")
		res := smtpConnection.SendMail("sender@local.host", "recipient1@domain.com;recipientN@anotherDomain.cn", "This is subject", "Hi! <br><br> This is body")
		if !strings.Contains(res.Error(), errMessage) {
			t.Errorf("Test failed on Address: %s", errValue)
		}
	}
}

func TestFrom(t *testing.T) {
	errMessage := "from is invalid"

	errValues := []string{
		"",
		"qwerty",
		// "qwe@rty@com",
		// "@rty",
		// "qwe@",
	}

	for _, errValue := range errValues {
		smtpConnection := gsmtp.New("smtp.exmail.qq.com", "smtpUser@smtp.exmail.qq.com", "smtpPassword")
		res := smtpConnection.SendMail(errValue, "recipient1@domain.com;recipientN@anotherDomain.cn", "This is subject", "Hi! <br><br> This is body")
		if !strings.Contains(res.Error(), errMessage) {
			t.Errorf("Test failed on From: %s", errValue)
		}
	}

}

func TestTos(t *testing.T) {
	errMessage := "tos if invalid"

	errValues := []string{
		"",
		"qwerty",
		"qwe;rty",
		"qwe;rty;com",
		// "qwe@rty@com",
		// "@rty",
		// "qwe@",
	}

	for _, errValue := range errValues {
		smtpConnection := gsmtp.New("smtp.exmail.qq.com", "smtpUser@smtp.exmail.qq.com", "smtpPassword")
		res := smtpConnection.SendMail("from@domain.com", errValue, "This is subject", "Hi! <br><br> This is body")
		if !strings.Contains(res.Error(), errMessage) {
			t.Errorf("Test failed on Tos: %s", errValue)
		}
	}

}
