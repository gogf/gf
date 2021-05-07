package main

import (
	"database/sql"

	"github.com/gogf/gf/os/gfile"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
)

func main() {
	db := g.DB()
	table := "medicine_clinics_upload_yinchuan"
	list, err := db.Table(table).All()
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	content := ""
	for _, item := range list {
		if j, err := gjson.DecodeToJson(item["upload_data"].String()); err != nil {
			panic(err)
		} else {
			s, _ := j.ToJsonIndentString()
			content += item["id"].String() + "\t" + item["medicine_clinic_id"].String() + "\t"
			content += s
			content += "\n\n"
			//if _, err := db.Table(table).Data("data_decode", s).Where("id", item["id"].Int()).Update(); err != nil {
			//	panic(err)
			//}
		}
	}
	gfile.PutContents("/Users/john/Temp/medicine_clinics_upload_yinchuan.txt", content)
}
