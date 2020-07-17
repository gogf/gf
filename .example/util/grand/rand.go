package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	mrand "math/rand"
	"os"
	"time"

	"github.com/jin502437344/gf/os/gtime"
)

// int 随机
func a1() {
	s1 := mrand.NewSource(time.Now().UnixNano())
	r1 := mrand.New(s1)
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", r1.Intn(100))
	}
	fmt.Printf("\n")
}

// 0/1 true/false  随机
func a2() {
	// Go编程这本书上例子.
	ch := make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case ch <- 0:
		case ch <- 1:
		}
		r := <-ch
		fmt.Printf("%d ", r)
	}
	fmt.Printf("\n")
}

//真随机 -- 用标准库封装好的
func a3() {
	b := make([]byte, 16)
	// On Unix-like systems, Reader reads from /dev/urandom.
	// On Windows systems, Reader uses the CryptGenRandom API.
	_, err := crand.Read(b) //返回长度为0 - 32 的值
	if err != nil {
		fmt.Println("[a3] ", err)
		return
	}
	fmt.Println("[a3] b:", b)
}

//真随机 -- 我们直接调真随机文件生成了事。 但注意，它是阻塞式的。
func a4() {
	f, err := os.Open("/dev/random")
	if err != nil {
		fmt.Println("[a4] ", err)
		return
	}
	defer f.Close()

	b1 := make([]byte, 16)
	_, err = f.Read(b1)
	if err != nil {
		fmt.Println("[a4] ", err)
		return
	}
	fmt.Println("[a4] Read /dev/random:", b1)
}

// a3 的另一种实现方式
func a5() {
	var ret int32
	binary.Read(crand.Reader, binary.LittleEndian, &ret)
	fmt.Println("[a5] ret:", ret)
}

func main() {
	fmt.Println("a1:", gtime.FuncCost(a1))
	fmt.Println("a2:", gtime.FuncCost(a2))
	fmt.Println("a3:", gtime.FuncCost(a3))
	fmt.Println("a4:", gtime.FuncCost(a4))
	fmt.Println("a5:", gtime.FuncCost(a5))
}
