package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/util/gconv"
)

func main() {
	type Score struct {
		Name   string
		Result int
	}
	type User1 struct {
		Scores Score
	}
	type User2 struct {
		Scores *Score
	}

	user1 := new(User1)
	user2 := new(User2)
	scores := g.Map{
		"Scores": g.Map{
			"Name":   "john",
			"Result": 100,
		},
	}

	if err := gconv.Struct(scores, user1); err != nil {
		fmt.Println(err)
	} else {
		g.Dump(user1)
	}
	if err := gconv.Struct(scores, user2); err != nil {
		fmt.Println(err)
	} else {
		g.Dump(user2)
	}
}
