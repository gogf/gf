package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
)

type ModifyFieldInfoType struct {
	Id  int64  `json:"id"`
	New string `json:"new"`
}
type ModifyFieldInfosType struct {
	Duration ModifyFieldInfoType `json:"duration"`
	OMLevel  ModifyFieldInfoType `json:"om_level"`
}

type MediaRequestModifyInfo struct {
	Modify ModifyFieldInfosType `json:"modifyFieldInfos"`
	Field  ModifyFieldInfosType `json:"fieldInfos"`
	FeedID string               `json:"feed_id"`
	Vid    string               `json:"id"`
}

var processQueue chan MediaRequestModifyInfo

func main() {

	jsonContent := `{"dataSetId":2001,"fieldInfos":{"duration":{"id":80079,"value":"59"},"om_level":{"id":2409,"value":"4"}},"id":"g0936lt1u0f","modifyFieldInfos":{"om_level":{"id":2409,"new":"4","old":""}},"timeStamp":1584599734}`
	var t MediaRequestModifyInfo
	err := gjson.DecodeTo(jsonContent, &t)
	fmt.Println(err)
	fmt.Printf("%+v\n", t)
	fmt.Println(gjson.New(t).MustToJsonString())

	b, _ := json.Marshal(t)
	fmt.Println(string(b))
}
