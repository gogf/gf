package grand

import (
    "time"
    "math/rand"
    "fmt"
)
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var digits  = []rune("0123456789")
const size  = 62

// 获得一个 min, max 之间的随机数
func Rand (min, max int) int {
    //fmt.Printf("min: %d, max: %d\n", min, max)
    if min >= max {
        return min
    }
    rand.Seed(time.Now().UnixNano())
    n := rand.Intn(max)
    if n < min {
        return Rand(min, max)
    }
    return n
}

// 获得指定长度的随机序列号
func RandSeq(n int) string {

    seed := time.Now().UnixNano()
    rand.Seed(seed)

    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(size)]
    }

    return fmt.Sprintf("%d_%s", seed, string(b))

}

// 获得指定长度的随机字符串
func RandStr(n int) string {
    seed := time.Now().UnixNano()
    rand.Seed(seed)

    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(size)]
    }

    return fmt.Sprintf("%s", string(b))
}

// 获得指定长度的随机数字字符串
func RandDigits(n int) string {
    seed := time.Now().UnixNano()
    rand.Seed(seed)

    b := make([]rune, n)
    for i := range b {
        b[i] = digits[rand.Intn(10)]
    }
    return fmt.Sprintf("%s", string(b))
}
