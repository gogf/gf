package main

import (
	"github.com/jin502437344/gf/encoding/gjson"
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/os/glog"
)

type GameUser struct {
	Uid        int         `json:"uid"`
	Account    string      `json:"account"`
	Tel        string      `json:"tel"`
	Role       string      `json:"role"`
	Vip        int         `json:"vip"`
	GameLevel  int         `json:"gamelevel"`
	Diamond    int         `json:"diamond"`
	Coin       int         `json:"coin"`
	Value      int         `json:"value"`
	Area       string      `json:"area"`
	ServerName string      `json:"servername"`
	Time       int         `json:"time"`
	ClientInfo *ClientInfo `json:"client_info"`
}

type ClientInfo struct {
	ClientGuid       string `json:"client_guid"`
	ClientType       int    `json:"client_type"`
	ClientSDKVersion string `json:"client_sdk_version"`
	ClientVersion    string `json:"client_version"`
	PackageId        string `json:"packageid"`
	PhoneType        string `json:"phone_type"`
	DevicesId        string `json:"devices_id"`
	ClientMac        string `json:"client_mac"`
}

func main() {
	s := `{
    "uid":9527,          			    
    "account":"zhangsan",          	    
    "tel":"15248787",          	     
    "role":"test",          		   
    "vip":7,          		   			
    "gamelevel":59,          		  
    "diamond ":59,          		   
    "coin ":59,          		      
    "value ":99,          		   
    "area":"s",          		       
    "servername":"灵动",          	
    "time":15454878787,          	
    "client_info": {
		"client_guid": "aaaa", 
		"client_type": 1,      
		"client_sdk_version": "1.0.1",  
		"client_version": "1.0.1",     
		"packageid":"",                 
		"phone_type": "vivi",
		"devices_id":"",   
		"client_mac":""      
    }
}`

	gameUser := &GameUser{}
	err := gjson.DecodeTo(s, gameUser)
	if err != nil {
		glog.Error(err)
	}
	g.Dump(gameUser)

}
