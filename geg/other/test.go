package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	//b, e := gjson.MarshalOrdered(g.Map{
	//	"a" : 1,
	//	"b" : 2,
	//	"c" : 3,
	//})
	//fmt.Println(e)
	//fmt.Println(string(b))

	//m := map[string]interface{}{
	//	"facet_is_special_price":[]string{"1"},
	//	"score_outlet":"0",
	//	"skus":[]string{"DI139BE71WDWDFMX", "DI139BE71WDWDFMX-519406"},
	//	"facet_novelty_two_days":[]string{"0"},
	//	"facet_brand":[]string{"139"},
	//	"sku":[]string{"DI139BE71WDWDFMX"},
	//}

	for {
		m := make(map[string]interface{})
		m["facet_is_special_price"] = []string{"1"}
		m["score_outlet"]           = "0"
		m["skus"]                   = []string{"DI139BE71WDWDFMX", "DI139BE71WDWDFMX-519406"}
		m["facet_novelty_two_days"] = []string{"0"}
		m["facet_brand"]            = []string{"139"}
		m["sku"]                    = []string{"DI139BE71WDWDFMX"}
		b, _ := json.Marshal(m)
		fmt.Println(string(b))
		time.Sleep(100*time.Millisecond)
	}

}