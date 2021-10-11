// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package main

import (
	"fmt"

	"github.com/gogf/gf/v2/net/gsmtp"
)

func main() {

	// create the SMTP connection
	smtpConnection := gsmtp.New("smtp.exmail.qq.com", "smtpUser@smtp.exmail.qq.com", "smtpPassword")
	// or you can specify the port explicitly
	// smtpConnection := smtp.New("smtp.exmail.qq.com:25", "smtpUser@smtp.exmail.qq.com", "smtpPassword")

	// send the Email
	fmt.Println(smtpConnection.SendMail("sender@local.host", "recipient1@domain.com;recipientN@anotherDomain.cn", "This is subject", "Hi! <br><br> This is body"))

}
