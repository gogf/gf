package main

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	tplContent := `
eq:
eq "a" "a": {{eq "a" "a"}}
eq "1" "1": {{eq "1" "1"}}
eq  1  "1": {{eq  1  "1"}}

ne:
ne  1  "1": {{ne  1  "1"}}
ne "a" "a": {{ne "a" "a"}}
ne "a" "b": {{ne "a" "b"}}

lt:
lt  1  "2": {{lt  1  "2"}}
lt  2   2 : {{lt  2   2 }}
lt "a" "b": {{lt "a" "b"}}

le:
le  1  "2": {{le  1  "2"}}
le  2   1 : {{le  2   1 }}
le "a" "a": {{le "a" "a"}}

gt:
gt  1  "2": {{gt  1  "2"}}
gt  2   1 : {{gt  2   1 }}
gt "a" "a": {{gt "a" "a"}}

ge:
ge  1  "2": {{ge  1  "2"}}
ge  2   1 : {{ge  2   1 }}
ge "a" "a": {{ge "a" "a"}}
`
	content, err := g.View().ParseContent(tplContent, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
}
