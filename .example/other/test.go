package main

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/text/gstr"
)

func main() {

	type Sender struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	type To struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	type SendReq struct {
		Sender *Sender `json:"sender"`
		//Name        string  `json:"name"`
		HtmlContent string `json:"htmlContent"`
		Subject     string `json:"subject"`
		To          []*To  `json:"to"`
	}

	//url := "emailCampaigns"
	sendreq := SendReq{
		Sender: &Sender{
			Name:  "123",
			Email: "globalclienthelp@gmail.com",
		},
		To: []*To{{
			Name: "456",
			//Email: order.Email,
			Email: "jinmao88@hotmail.com",
		}},
	}
	subject := ""
	htmlcontent := ""

	subject = " Your order：" + gstr.Split("11111", "-")[1] + " ，  Shipment is on its way."
	htmlcontent = "test"
	sendreq.Subject = subject
	sendreq.HtmlContent = htmlcontent
	j, err := gjson.DecodeToJson(sendreq)
	fmt.Println(err)
	fmt.Println(j)
}
