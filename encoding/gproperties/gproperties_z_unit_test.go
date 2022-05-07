package gproperties_test

import (
	"encoding/json"
	"fmt"
	"testing"

	_ "github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/encoding/gproperties"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

var pStr string = `
# 模板引擎目录
data = "/home/www/templates/"
# MySQL数据库配置
sql.disk.0  = 127.0.0.1:6379,0
sql.cache.0 = 127.0.0.1:6379,1=
sql.cache.1=0
sql.disk.a = 10
`

func TestDecode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {

		decodeStr, err := gproperties.Decode(([]byte)(pStr))
		if err != nil {
			t.Errorf("decode failed. %v", err)
			return
		}
		fmt.Printf("%v\n", decodeStr)
		v, _ := json.Marshal(decodeStr)
		fmt.Printf("%v\n", string(v))
	})
}
func TestEncode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {

		decodeStr, err := gproperties.Encode(map[string]interface{}{
			"sql": g.Map{
				"userName": "admin",
				"password": "123456",
			},
			"user": g.Slice{"admin", "aa"},
			"no":   123,
		})
		if err != nil {
			t.Errorf("decode failed. %v", err)
			return
		}
		fmt.Printf("%v\n", string(decodeStr))
	})
}
