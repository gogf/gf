package main

import (
	"fmt"

	"github.com/jin502437344/gf/container/gring"
)

type Player struct {
	position int  // 位置
	alive    bool // 是否存活
}

const (
	playerCount = 41 // 玩家人数
	startPos    = 1  // 开始报数位置
)

var (
	deadline = 3
)

func main() {
	// 关闭并发安全，当前场景没有必要
	r := gring.New(playerCount, false)

	// 设置所有玩家初始值
	for i := 1; i <= playerCount; i++ {
		r.Put(&Player{i, true})
	}

	// 如果开始报数的位置不为1，则设置开始位置
	if startPos > 1 {
		r.Move(startPos - 1)
	}

	counter := 1   // 报数从1开始，因为下面的循环从第二个开始计算
	deadCount := 0 // 死亡人数，初始值为0

	// 直到所有人都死亡，否则循环一直执行
	for deadCount < playerCount {
		// 跳到下一个人
		r.Next()

		// 如果是活着的人，则报数
		if r.Val().(*Player).alive {
			counter++
		}

		// 如果报数为deadline，则此人淘汰出局
		if counter == deadline {
			r.Val().(*Player).alive = false
			fmt.Printf("Player %d died!\n", r.Val().(*Player).position)
			deadCount++
			counter = 0
		}
	}
}
