package gjson_test

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
)

func ExampleJson_SetSplitChar() {
	data :=
		`{
        "users" : {
            "count" : 2,
            "list"  : [
                {"name" : "Ming",  "score" : 60},
                {"name" : "John", "score" : 99.5}
            ]
        }
    }`
	if j, err := gjson.DecodeToJson(data); err != nil {
		panic(err)
	} else {
		j.SetSplitChar('#')
		fmt.Println("John Score:", j.Get("users#list#1#score").Float32())
	}
	// Output:
	// John Score: 99.5
}

func ExampleJson_SetViolenceCheck() {
	data :=
		`{
        "users" : {
            "count" : 2,
            "list"  : [
                {"name" : "Ming",  "score" : 60},
                {"name" : "John", "score" : 99.5}
            ]
        }
    }`
	if j, err := gjson.DecodeToJson(data); err != nil {
		panic(err)
	} else {
		j.SetViolenceCheck(false)
		fmt.Println("John Score:", j.Get("users.list.1.score").Float32())
	}
	// Output:
	// John Score: 99.5
}
