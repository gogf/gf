package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

func main() {
	data := []byte(`
m:
 k: v
    `)
	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		panic(err)
	}
	b, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
