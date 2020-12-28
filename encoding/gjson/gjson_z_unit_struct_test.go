// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func Test_GetScan(t *testing.T) {
	type User struct {
		Name  string
		Score float64
	}
	j := gjson.New(`[{"name":"john", "score":"100"},{"name":"smith", "score":"60"}]`)
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := j.GetScan("1", &user)
		t.Assert(err, nil)
		t.Assert(user, &User{
			Name:  "smith",
			Score: 60,
		})
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := j.GetScan(".", &users)
		t.Assert(err, nil)
		t.Assert(users, []User{
			{
				Name:  "john",
				Score: 100,
			},
			{
				Name:  "smith",
				Score: 60,
			},
		})
	})
}

func Test_GetScanDeep(t *testing.T) {
	type User struct {
		Name  string
		Score float64
	}
	j := gjson.New(`[{"name":"john", "score":"100"},{"name":"smith", "score":"60"}]`)
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := j.GetScanDeep("1", &user)
		t.Assert(err, nil)
		t.Assert(user, &User{
			Name:  "smith",
			Score: 60,
		})
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := j.GetScanDeep(".", &users)
		t.Assert(err, nil)
		t.Assert(users, []User{
			{
				Name:  "john",
				Score: 100,
			},
			{
				Name:  "smith",
				Score: 60,
			},
		})
	})
}

func Test_ToScan(t *testing.T) {
	type User struct {
		Name  string
		Score float64
	}
	j := gjson.New(`[{"name":"john", "score":"100"},{"name":"smith", "score":"60"}]`)
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := j.ToScan(&users)
		t.Assert(err, nil)
		t.Assert(users, []User{
			{
				Name:  "john",
				Score: 100,
			},
			{
				Name:  "smith",
				Score: 60,
			},
		})
	})
}

func Test_ToScanDeep(t *testing.T) {
	type User struct {
		Name  string
		Score float64
	}
	j := gjson.New(`[{"name":"john", "score":"100"},{"name":"smith", "score":"60"}]`)
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := j.ToScanDeep(&users)
		t.Assert(err, nil)
		t.Assert(users, []User{
			{
				Name:  "john",
				Score: 100,
			},
			{
				Name:  "smith",
				Score: 60,
			},
		})
	})
}

func Test_ToStruct1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type BaseInfoItem struct {
			IdCardNumber        string `db:"id_card_number" json:"idCardNumber" field:"id_card_number"`
			IsHouseholder       bool   `db:"is_householder" json:"isHouseholder" field:"is_householder"`
			HouseholderRelation string `db:"householder_relation" json:"householderRelation" field:"householder_relation"`
			UserName            string `db:"user_name" json:"userName" field:"user_name"`
			UserSex             string `db:"user_sex" json:"userSex" field:"user_sex"`
			UserAge             int    `db:"user_age" json:"userAge" field:"user_age"`
			UserNation          string `db:"user_nation" json:"userNation" field:"user_nation"`
		}

		type UserCollectionAddReq struct {
			BaseInfo []BaseInfoItem `db:"_" json:"baseInfo" field:"_"`
		}
		jsonContent := `{
	"baseInfo": [{
		"idCardNumber": "520101199412141111",
		"isHouseholder": true,
		"householderRelation": "户主",
		"userName": "李四",
		"userSex": "男",
		"userAge": 32,
		"userNation": "苗族",
		"userPhone": "13084183323",
		"liveAddress": {},
		"occupationInfo": [{
			"occupationType": "经商",
			"occupationBusinessInfo": [{
				"occupationClass": "制造业",
				"businessLicenseNumber": "32020000012300",
				"businessName": "土灶柴火鸡",
				"spouseName": "",
				"spouseIdCardNumber": "",
				"businessLicensePhotoId": 125,
				"businessPlace": "租赁房产",
				"hasGoodsInsurance": true,
				"businessScopeStr": "柴火鸡;烧烤",
				"businessAddress": {},
				"businessPerformAbility": {
					"businessType": "服务业",
					"businessLife": 5,
					"salesRevenue": 8000,
					"familyEquity": 6000
				}
			}],
			"occupationWorkInfo": {
				"occupationClass": "",
				"companyName": "",
				"companyType": "",
				"workYearNum": 0,
				"spouseName": "",
				"spouseIdCardNumber": "",
				"spousePhone": "",
				"spouseEducation": "",
				"spouseCompanyName": "",
				"workLevel": "",
				"workAddress": {},
				"workPerformAbility": {
					"familyAnnualIncome": 0,
					"familyEquity": 0,
					"workCooperationState": "",
					"workMoneyCooperationState": ""
				}
			},
			"occupationAgricultureInfo": []
		}],
		"assetsInfo": [],
		"expenditureInfo": [],
		"incomeInfo": [],
		"liabilityInfo": []
	}]
}`
		data := new(UserCollectionAddReq)
		j, err := gjson.LoadJson(jsonContent)
		t.Assert(err, nil)
		err = j.ToStruct(data)
		t.Assert(err, nil)
	})
}

func Test_ToStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Item struct {
			Title string `json:"title"`
			Key   string `json:"key"`
		}

		type M struct {
			Id    string                 `json:"id"`
			Me    map[string]interface{} `json:"me"`
			Txt   string                 `json:"txt"`
			Items []*Item                `json:"items"`
		}

		txt := `{
		  "id":"88888",
		  "me":{"name":"mikey","day":"20009"},
		  "txt":"hello",
		  "items":null
		 }`

		j, err := gjson.LoadContent(txt)
		t.Assert(err, nil)
		t.Assert(j.GetString("me.name"), "mikey")
		t.Assert(j.GetString("items"), "")
		t.Assert(j.GetBool("items"), false)
		t.Assert(j.GetArray("items"), nil)
		m := new(M)
		err = j.ToStruct(m)
		t.Assert(err, nil)
		t.AssertNE(m.Me, nil)
		t.Assert(m.Me["day"], "20009")
		t.Assert(m.Items, nil)
	})
}

func Test_ToStruct_Complicated(t *testing.T) {
	type CertInfo struct {
		UserRealName        string `json:"userRealname,omitempty"`
		IdentType           string `json:"identType,omitempty"`
		IdentNo             string `json:"identNo,omitempty"`
		CompanyName         string `json:"companyName,omitempty"`
		Website             string `json:"website,omitempty"`
		RegisterNo          string `json:"registerNo,omitempty"`
		AreaCode            string `json:"areaCode,omitempty"`
		Address             string `json:"address,omitempty"`
		CommunityCreditCode string `json:"communityCreditCode,omitempty"`
		PhoneNumber         string `json:"phoneNumber,omitempty"`
		AreaName            string `json:"areaName,omitempty"`
		PhoneAreaCode       string `json:"phoneAreaCode,omitempty"`
		OperateRange        string `json:"operateRange,omitempty"`
		Email               string `json:"email,omitempty"`
		LegalPersonName     string `json:"legalPersonName,omitempty"`
		OrgCode             string `json:"orgCode,omitempty"`
		BusinessLicense     string `json:"businessLicense,omitempty"`
		FilePath1           string `json:"filePath1,omitempty"`
		MobileNo            string `json:"mobileNo,omitempty"`
		CardName            string `json:"cardName,omitempty"`
		BankMobileNo        string `json:"bankMobileNo,omitempty"`
		BankCode            string `json:"bankCode,omitempty"`
		BankCard            string `json:"bankCard,omitempty"`
	}

	type CertList struct {
		StatusCode uint     `json:"statusCode,string"`
		SrcType    uint     `json:"srcType,string"`
		CertID     string   `json:"certId"`
		CardType   string   `json:"cardType,omitempty"`
		CertInfo   CertInfo `json:"certInfo"`
	}

	type Response struct {
		UserLevel uint       `json:"userLevel,string,omitempty"`
		CertList  []CertList `json:"certList"`
	}

	gtest.C(t, func(t *gtest.T) {
		jsonContent := `{
"certList":[
{"certId":"2023313","certInfo":"{\"address\":\"xxxxxxx\",\"phoneNumber\":\"15084890\",\"companyName\":\"dddd\",\"communityCreditCode\":\"91110111MBE1G2B\",\"operateRange\":\"fff\",\"registerNo\":\"91110111MA00G2B\",\"legalPersonName\":\"rrr\"}","srcType":"1","statusCode":"2"},
{"certId":"2023314","certInfo":"{\"identNo\":\"342224196507051\",\"userRealname\":\"xxxx\",\"identType\":\"01\"}","srcType":"8","statusCode":"0"},
{"certId":"2023322","certInfo":"{\"businessLicense\":\"91110111MA00BE1G\",\"companyName\":\"sssss\",\"communityCreditCode\":\"91110111MA00BE1\"}","srcType":"2","statusCode":"0"}
]
}`
		j, err := gjson.LoadContent(jsonContent)
		t.Assert(err, nil)
		var response = new(Response)
		err = j.ToStruct(response)
		t.Assert(err, nil)
		t.Assert(len(response.CertList), 3)
		t.Assert(response.CertList[0].CertID, 2023313)
		t.Assert(response.CertList[1].CertID, 2023314)
		t.Assert(response.CertList[2].CertID, 2023322)
		t.Assert(response.CertList[0].CertInfo.PhoneNumber, "15084890")
		t.Assert(response.CertList[1].CertInfo.IdentNo, "342224196507051")
		t.Assert(response.CertList[2].CertInfo.BusinessLicense, "91110111MA00BE1G")
	})
}
