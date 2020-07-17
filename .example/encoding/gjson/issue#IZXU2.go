package main

import (
	"encoding/json"
	"fmt"

	"github.com/jin502437344/gf/encoding/gjson"
)

type XinYanModel struct {
	Success   bool        `json:"success"`
	Data      Data        `json:"data"`
	ErrorCode interface{} `json:"errorCode"`
	ErrorMsg  interface{} `json:"errorMsg"`
}
type ApplyReportDetail struct {
	ApplyScore        string `json:"apply_score"`
	ApplyCredibility  string `json:"apply_credibility"`
	QueryOrgCount     string `json:"apply_query_org_count"`
	QueryFinanceCount string `json:"apply_query_finance_count"`
	QueryCashCount    string `json:"apply_query_cash_count"`
	QuerySumCount     string `json:"apply_query_sum_count"`
	LatestQueryTime   string `json:"apply_latest_query_time"`
	LatestOneMonth    string `json:"apply_latest_one_month"`
	LatestThreeMonth  string `json:"apply_latest_three_month"`
	LatestSixMonth    string `json:"apply_latest_six_month"`
}
type BehaviorReportDetail struct {
	LoansScore         string `json:"behavior_report_detailloans_score"`
	LoansCredibility   string `json:"behavior_report_detailloans_credibility"`
	LoansCount         string `json:"behavior_report_detailloans_count"`
	LoansSettleCount   string `json:"behavior_report_detailloans_settle_count"`
	LoansOverdueCount  string `json:"behavior_report_detailloans_overdue_count"`
	LoansOrgCount      string `json:"behavior_report_detailloans_org_count"`
	ConsfinOrgCount    string `json:"behavior_report_detailconsfin_org_count"`
	LoansCashCount     string `json:"behavior_report_detailloans_cash_count"`
	LatestOneMonth     string `json:"behavior_report_detaillatest_one_month"`
	LatestThreeMonth   string `json:"behavior_report_detaillatest_three_month"`
	LatestSixMonth     string `json:"behavior_report_detaillatest_six_month"`
	HistorySucFee      string `json:"behavior_report_detailhistory_suc_fee"`
	HistoryFailFee     string `json:"behavior_report_detailhistory_fail_fee"`
	LatestOneMonthSuc  string `json:"behavior_report_detaillatest_one_month_suc"`
	LatestOneMonthFail string `json:"behavior_report_detaillatest_one_month_fail"`
	LoansLongTime      string `json:"behavior_report_detailloans_long_time"`
	LoansLatestTime    string `json:"behavior_report_detailloans_latest_time"`
}
type CurrentReportDetail struct {
	LoansCreditLimit    string `json:"current_report_detailloans_credit_limit"`
	LoansCredibility    string `json:"current_report_detailloans_credibility"`
	LoansOrgCount       string `json:"current_report_detailloans_org_count"`
	LoansProductCount   string `json:"current_report_detailloans_product_count"`
	LoansMaxLimit       string `json:"current_report_detailloans_max_limit"`
	LoansAvgLimit       string `json:"current_report_detailloans_avg_limit"`
	ConsfinCreditLimit  string `json:"current_report_detailconsfin_credit_limit"`
	ConsfinCredibility  string `json:"current_report_detailconsfin_credibility"`
	ConsfinOrgCount     string `json:"current_report_detailconsfin_org_count"`
	ConsfinProductCount string `json:"current_report_detailconsfin_product_count"`
	ConsfinMaxLimit     string `json:"current_report_detailconsfin_max_limit"`
	ConsfinAvgLimit     string `json:"current_report_detailconsfin_avg_limit"`
}
type ResultDetail struct {
	ApplyReportDetail    ApplyReportDetail    `json:"apply_report_detail"`
	BehaviorReportDetail BehaviorReportDetail `json:"behavior_report_detail"`
	CurrentReportDetail  CurrentReportDetail  `json:"current_report_detail"`
}
type Data struct {
	Code         string       `json:"code"`
	Desc         string       `json:"desc1"`
	TransID      string       `json:"trans_id"`
	TradeNo      string       `json:"trade_no"`
	Fee          string       `json:"fee"`
	IDNo         string       `json:"id_no"`
	IDName       string       `json:"id_name"`
	Versions     string       `json:"versions"`
	ResultDetail ResultDetail `json:"result_detail"`
}

var data = `{
    "success": true,
    "data": {
        "code": "0",
        "desc": "查询成功",
        "trans_id": "14910304379231213",
        "trade_no": "201704011507240100057329",
        "fee": "Y",
        "id_no": "0783231bcc39f4957e99907e02ae401c",
        "id_name": "dd67a5943781369ddd7c594e231e9e70 ",
        "versions": "1.0.0",
        "result_detail":{
            "apply_report_detail": {
                "apply_score": "189",
                "apply_credibility": "84",
                "query_org_count": "7",
                "query_finance_count": "2",
                "query_cash_count": "2",
                "query_sum_count": "13",
                "latest_query_time": "2017-09-03",
                "latest_one_month": "1",
                "latest_three_month": "5",
                "latest_six_month": "12"
            },
            "behavior_report_detail": {
                "loans_score": "199",
                "loans_credibility": "90",
                "loans_count": "300",
                "loans_settle_count": "280",
                "loans_overdue_count": "20",
                "loans_org_count": "5",
                "consfin_org_count": "3",
                "loans_cash_count": "2",
                "latest_one_month": "3",
                "latest_three_month": "20",
                "latest_six_month": "23",
                "history_suc_fee": "30",
                "history_fail_fee": "25",
                "latest_one_month_suc": "5",
                "latest_one_month_fail": "20",
                "loans_long_time": "130",
                "loans_latest_time": "2017-09-16"
            },
            "current_report_detail": {
                "loans_credit_limit": "1400",
                "loans_credibility": "80",
                "loans_org_count": "7",
                "loans_product_count": "8",
                "loans_max_limit": "2000",
                "loans_avg_limit": "1000",
                "consfin_credit_limit": "1500",
                "consfin_credibility": "90",
                "consfin_org_count": "8",
                "consfin_product_count": "5",
                "consfin_max_limit": "5000",
                "consfin_avg_limit": "3000"
            }
        }
    },
    "errorCode": null,
    "errorMsg": null
}`

func main() {
	struct1 := new(XinYanModel)
	err := json.Unmarshal([]byte(data), struct1)
	fmt.Println(err)
	fmt.Println(struct1)

	fmt.Println()

	struct2 := new(XinYanModel)
	j, err := gjson.DecodeToJson(data)
	fmt.Println(err)
	fmt.Println(j.Get("data.desc"))
	err = j.ToStruct(struct2)
	fmt.Println(err)
	fmt.Println(struct2)
}
