package gjson

import (
	"fmt"
	"testing"
)

type TestStruct struct {
	Result []map[string]interface{} `json:"result"`
}

func TestA(t *testing.T) {
	ts := &TestStruct{
		Result: []map[string]interface{}{
			{
				"Name": nil,
			},
		},
	}
	a := New(ts)
	fmt.Println(a)
}
