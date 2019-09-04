package main

import (
	"fmt"

	"github.com/gogf/gf/net/gsmtp"
)

func main() {

	smtpConnection := gsmtp.New("smtp.exmail.qq.com", "smtpUser@smtp.exmail.qq.com", "smtpPassword")
	// or you can specify the port explicitly
	// smtpConnection := smtp.New("smtp.exmail.qq.com:25", "smtpUser@smtp.exmail.qq.com", "smtpPassword")
	fmt.Println(smtpConnection.SendMail("sender@local.host", "recipient1@domain.com;recipientN@anotherDomain.cn", "This is subject", "Hi! <br><br> This is body"))

}
