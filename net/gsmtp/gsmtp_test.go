// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gsmtp_test

import (
	"testing"

	"github.com/gogf/gf/net/gsmtp"
)

func TestGsmtp(t *testing.T) {
	addressErrMessage := "address is either empty or incorrect"

	addressErrorValues := []string{
		"",
		":",
		":25",
		"localhost:",
		"local.host:25:28",
	}

	for _, errValue := range addressErrorValues {
		smtpConnection := gsmtp.New(errValue, "smtpUser@smtp.exmail.qq.com", "smtpPassword")
		res := smtpConnection.SendMail("sender@local.host", "recipient1@domain.com;recipientN@anotherDomain.cn", "This is subject", "Hi! <br><br> This is body")
		if res.Error() != addressErrMessage {
			t.Errorf("Test failed on Address: %s", errValue)
		}
	}

}
